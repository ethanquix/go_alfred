package globals

import (
	"os"
	"strings"
)

// Only primitive types here. NO IMPORTS!!!

var is_prod bool
var is_prod_ready bool = false
var INTERNAL_SERVICE_TOKEN string

func init() {
	is_prod = strings.ToLower(os.Getenv("PROD")) == "true"
	is_prod_ready = true

	if is_prod {
		INTERNAL_SERVICE_TOKEN = os.Getenv("internal_services_token")
		if INTERNAL_SERVICE_TOKEN == "" {
			panic("no env var $internal_services_token set!!!")
		}
	}
}

func IS_PROD() bool {
	if !is_prod_ready {
		panic("import github.com/ethanquix/go_alfred/pkg/globals/globals first!")
	}
	return is_prod
}

func SetIsProd(val bool) {
	is_prod = val
}
