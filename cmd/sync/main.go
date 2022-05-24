package main

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/matthiasbruns/gambio-gx3-go/client"
	log "github.com/sirupsen/logrus"

	"github.com/happyann/happyann-gambio/internal"
	"github.com/happyann/happyann-gambio/internal/gambio"
	"github.com/happyann/happyann-gambio/internal/happyann"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// 2022-05-10 15:02:19
const dateFormat = "2006-01-02 15:04:05"

func main() {
	shopBasePath := internal.GetShopBasePath()
	shopIdentifier := internal.GetShopIdentifier()

	categoryMap := make(map[int64]client.GxCategory, 0)

	var page int64 = 0
	products := make([]happyann.ProductData, 0)
	paginate := true

	for paginate {
		importedProducts, productCategoryMap := fetchProducts(shopBasePath, shopIdentifier, page)
		products = append(products, importedProducts...)

		for _, product := range products {
			categoryIds := productCategoryMap[product.Id]
			categories := make([]happyann.CategoryData, 0)
			for _, catId := range categoryIds {
				cat, ok := categoryMap[catId]
				if !ok {
					catFromApi := gambio.FetchCategoryById(catId)
					if catFromApi == nil {
						continue
					}

					categoryMap[catId] = *catFromApi
					cat = *catFromApi

					categoryData := happyann.CategoryData{
						Id:       cat.Id,
						ParentId: cat.ParentId,
					}

					title := cat.Name.De
					if title == "" {
						title = cat.Name.En
					}
					categoryData.Title = title

					lastMod, err := time.Parse(dateFormat, cat.LastModified)
					if err != nil {
						log.WithError(err).Warnf("Could not parse LastModified %s", cat.LastModified)
					} else {
						categoryData.LastModified = lastMod
					}

					dateAdded, err := time.Parse(dateFormat, cat.DateAdded)
					if err != nil {
						log.WithError(err).Warnf("Could not parse DateAdded %s", cat.DateAdded)
					} else {
						categoryData.DateAdded = dateAdded
					}

					categories = append(categories, categoryData)

				}
			}
			product.Categories = categories
		}

		paginate = len(products) > 0
		page++
	}

	fmt.Sprintln(products)
}

func fetchProductDetails(productId int64, detailsChan chan *client.GxProduct) {
	product := gambio.FetchProductDetails(productId)
	detailsChan <- product
}

func fetchProductLinks(productId int64, linksChan chan []int64) {
	links := gambio.FetchProductLink(productId)
	linksChan <- links
}

func fetchProducts(shopBasePath string, shopIdentifier string, page int64) ([]happyann.ProductData, map[int64][]int64) {
	productListings := fetchProductListingsPage(page)
	gxProducts := make([]*client.GxProduct, 0)
	productLinksMap := make(map[int64][]int64, 0)

	for _, listing := range productListings {
		details := gambio.FetchProductDetails(listing.Id)
		gxProducts = append(gxProducts, details)

		// Product links
		links := gambio.FetchProductLink(listing.Id)
		productLinksMap[listing.Id] = links

	}

	updates := make([]happyann.ProductData, 0)
	for _, gxProduct := range gxProducts {
		// Product prices
		prices := gambio.FetchProductPrices(gxProduct.Id)
		if prices != nil {
			log.Warnf("Did not handle custom price for product %d", gxProduct.Id)
			// TODO add more details pricing handling
		}

		images := make([]string, 0)
		for _, image := range gxProduct.Images {
			images = append(images, shopBasePath+"/images/product_images/info_images/"+url.QueryEscape(image.Filename))
		}

		if strings.Contains(gxProduct.Name.De, "Sonder-Special-Paket Aquarell Flowers") {
			println("Found it")
		}

		price := int(math.Round(float64(gxProduct.Price) * 100.0))
		update := happyann.ProductData{
			Id:           gxProduct.Id,
			Source:       shopIdentifier,
			Title:        gxProduct.Name.De,
			Images:       images,
			PriceUnit:    "PIECE",
			UnitCount:    1,
			IsActive:     gxProduct.IsActive,
			PricePerUnit: price,
			ProductTypes: nil,
			Categories:   nil,
		}

		lastMod, err := time.Parse(dateFormat, gxProduct.LastModified)
		if err != nil {
			log.WithError(err).Warnf("Could not parse LastModified %s", gxProduct.LastModified)
		} else {
			update.LastModified = lastMod
		}

		dateAdded, err := time.Parse(dateFormat, gxProduct.DateAdded)
		if err != nil {
			log.WithError(err).Warnf("Could not parse DateAdded %s", gxProduct.DateAdded)
		} else {
			update.DateAdded = dateAdded
		}

		title := gxProduct.Name.De
		if title == "" {
			title = gxProduct.Name.En
		}
		update.Title = title

		desc := gxProduct.Description.De
		if desc == "" {
			desc = gxProduct.Description.En
		}
		update.Description = desc

		var pUrl string
		if gxProduct.Url.De != "" {
			pUrl = gxProduct.Url.De
		} else if gxProduct.UrlKeywords.De != "" {
			pUrl = shopBasePath + gxProduct.UrlKeywords.De + ".html"
		} else if gxProduct.Url.En != "" {
			pUrl = gxProduct.Url.En
		} else if gxProduct.UrlKeywords.En != "" {
			pUrl = shopBasePath + gxProduct.UrlKeywords.En + ".html"
		} else {
			log.Errorf("No pUrl found for %d", gxProduct.Id)
			continue
		}
		update.Url = pUrl

		updates = append(updates, update)
	}

	return updates, productLinksMap
}

func fetchProductListingsPage(page int64) []client.GxProductListing {
	productListings := gambio.FetchProductsFromShop(page)
	return productListings
}
