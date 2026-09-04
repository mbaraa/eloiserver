package main

import (
	"log"
	"net/http"
	"os"

	"codeberg.org/lordbaraa/eloiserver/apis"
	"codeberg.org/lordbaraa/eloiserver/config"
	"codeberg.org/lordbaraa/eloiserver/utils/overlays"
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
