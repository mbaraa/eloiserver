package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mbaraa/eloiserver/apis"
	"github.com/mbaraa/eloiserver/config"
	"github.com/mbaraa/eloiserver/globals"
	"github.com/mbaraa/eloiserver/utils/gposcrapper"
	"github.com/mbaraa/eloiserver/utils/overlays"
)

func main() {
	var err error
	globals.Overlays, err = overlays.LoadOverlays()
	if os.IsNotExist(err) {
		fmt.Println("Downloading Overlays")
		globals.Overlays = gposcrapper.GetOverlays()

		f, err := os.Create("overlays.json")
		if err != nil {
			panic(err)
		}

		json.NewEncoder(f).Encode(globals.Overlays)
	}

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
