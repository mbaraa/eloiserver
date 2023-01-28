package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mbaraa/eloiserver/apis"
	"github.com/mbaraa/eloiserver/config"
	"github.com/mbaraa/eloiserver/utils/overlays"
)

func main() {
	err := overlays.LoadOverlays()
	if !os.IsNotExist(err) {
		panic(err)
	}

	err = overlays.ScheduleScrapper()
	if err != nil {
		panic(err)
	}

	////

	fmt.Println("Starting HTTP server on port", config.PortNumber())
	http.Handle("/overlays/", apis.NewOverlaysAPI())
	err = http.ListenAndServe(":"+config.PortNumber(), nil)
	if err != nil {
		panic(err)
	}
}
