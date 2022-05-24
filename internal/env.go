package internal

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func GetShopBasePath() string {
	env := os.Getenv("SHOP_BASE_PATH")
	if env == "" {
		log.Error(fmt.Errorf("env 'SHOP_BASE_PATH' is not set"))
	}
	return env
}

func GetApiBasePath() string {
	env := os.Getenv("SHOP_BASE_PATH")
	if env == "" {
		log.Error(fmt.Errorf("env 'SHOP_BASE_PATH' is not set"))
	}
	return env + "/api.php/v2"
}

func GetShopIdentifier() string {
	env := os.Getenv("SHOP_IDENTIFIER")
	if env == "" {
		log.Error(fmt.Errorf("env 'SHOP_IDENTIFIER' is not set"))
	}
	return env
}

func GetUserAgent() string {
	env := os.Getenv("GAMBIO_USER_AGENT")
	if env == "" {
		return "happyann-gambio"
	}
	return env
}

func GetApiUser() string {
	env := os.Getenv("GAMBIO_API_USER")
	if env == "" {
		log.Error(fmt.Errorf("env 'GAMBIO_API_USER' is not set"))
	}
	return env
}

func GetApiPassword() string {
	env := os.Getenv("GAMBIO_API_PASSWORD")
	if env == "" {
		log.Error(fmt.Errorf("env 'GAMBIO_API_PASSWORD' is not set"))
	}
	return env
}
