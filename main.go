package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/kfortney/tippingpointbeat/beater"
)

func main() {
	err := beat.Run("tippingpointbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
