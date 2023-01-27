package main

import (
	"fmt"
	"net/http"

	"github.com/mbaraa/eloiserver/apis"
	"github.com/mbaraa/eloiserver/config"
	"github.com/mbaraa/eloiserver/globals"
	"github.com/mbaraa/eloiserver/utils/overlays"
)

func main() {
	var err error
	globals.Overlays, err = overlays.LoadOverlays()
	if err != nil {
		panic(err)
	}

	globals.Ebuilds = overlays.ExtractEbuilds(globals.Overlays)

	fmt.Println("Starting HTTP server on port", config.PortNumber())
	http.Handle("/overlays/", apis.NewOverlaysAPI())
	err = http.ListenAndServe(":"+config.PortNumber(), nil)
	if err != nil {
		panic(err)
	}
}
