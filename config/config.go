package config

import (
	"os"

	// "github.com/spf13/viper"
)

type configuration struct {
	env string
}

// Config -
var Config configuration


func init() {
	if os.Getenv("ENVIRONMENT") == "DEV" || os.Getenv("ENVIRONMENT") == "TEST" {

	} else {
		
	}

}