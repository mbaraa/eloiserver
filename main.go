package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mbaraa/eloiserver/apis"
	"github.com/mbaraa/eloiserver/config"
	"github.com/mbaraa/eloiserver/utils/overlays"
)

func main() {
	err := overlays.LoadOverlays()
	if os.IsNotExist(err) {
		err = overlays.ScrapeOverlays()
		if err != nil {
			panic(err)
		}
	}

	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	err = overlays.ScheduleScrapper()
	if err != nil {
		panic(err)
	}

	////

	log.Println("Starting HTTP server on port", config.PortNumber())
	http.Handle("/overlays/", apis.NewOverlaysAPI())
	err = http.ListenAndServe(":"+config.PortNumber(), nil)
	if err != nil {
		panic(err)
	}
}
