package gambio

import (
	"context"
	"time"

	"github.com/antihax/optional"
	"github.com/matthiasbruns/gambio-gx3-go/client"
	log "github.com/sirupsen/logrus"

	"github.com/happyann/happyann-gambio/internal"
)

func apiClient() *client.APIClient {
	basePath := internal.GetApiBasePath()
	userAgent := internal.GetUserAgent()
	gambioConfig := client.NewConfiguration()
	gambioConfig.BasePath = basePath // "https://www.monalienchen.de/api.php/v2"
	gambioConfig.UserAgent = userAgent
	return client.NewAPIClient(gambioConfig)
}

func ctxWithTimeout() (context.Context, context.CancelFunc) {
	auth := context.WithValue(
		context.Background(), client.ContextBasicAuth, client.BasicAuth{
			UserName: internal.GetApiUser(),
			Password: internal.GetApiPassword(),
		},
	)

	return context.WithTimeout(auth, time.Second*30)
}

func FetchProductsFromShop(page int64) []client.GxProductListing {
	const perPage int64 = 100
	ctx, cancelFunc := ctxWithTimeout()
	defer cancelFunc()

	products, _, err := apiClient().ProductsApi.GetProducts(
		ctx, &client.ProductsApiGetProductsOpts{
			Page:    optional.NewInt64(page),
			PerPage: optional.NewInt64(perPage),
			Sort:    optional.NewInterface([]string{"-lastModified"}),
		},
	)
	if err != nil {
		log.WithError(err).Errorf("Could not fetch products from api")
		return make([]client.GxProductListing, 0)
	}

	return products
}

func FetchProductDetails(id int64) *client.GxProduct {
	ctx, cancelFunc := ctxWithTimeout()
	defer cancelFunc()

	product, _, err := apiClient().ProductsApi.GetProduct(ctx, id)
	if err != nil {
		log.WithError(err).Errorf("Could not fetch product from api")
		return nil
	}

	return &product
}

func FetchProductLink(id int64) []int64 {
	ctx, cancelFunc := ctxWithTimeout()
	defer cancelFunc()

	links, _, err := apiClient().ProductsApi.GetProductLinks(ctx, id)
	if err != nil {
		log.WithError(err).Errorf("Could not fetch category from api")
		return make([]int64, 0)
	}

	return links
}

func FetchCategoryById(id int64) *client.GxCategory {
	ctx, cancelFunc := ctxWithTimeout()
	defer cancelFunc()

	cat, _, err := apiClient().CategoriesApi.GetCategory(ctx, id)
	if err != nil {
		log.WithError(err).Errorf("Could not fetch category from api")
		return nil
	}

	return &cat
}

func FetchProductPrices(id int64) []client.GxProductPrices {
	ctx, cancelFunc := ctxWithTimeout()
	defer cancelFunc()

	prices, _, err := apiClient().ProductPricesApi.GetProductPrices(ctx, id)
	if err != nil {
		log.WithError(err).Errorf("Could not fetch prices from api")
		return nil
	}

	return prices
}
