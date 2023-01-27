package apis

import (
	"encoding/json"
	"net/http"

	"github.com/mbaraa/eloiserver/globals"
)

type OverlaysAPI struct {
	endpoints map[string]http.HandlerFunc
}

func NewOverlaysAPI() http.Handler {
	return new(OverlaysAPI).initEndPoints()
}

func (o *OverlaysAPI) initEndPoints() *OverlaysAPI {
	o.endpoints = map[string]http.HandlerFunc{
		"GET /all":     o.getOverlays,
		"GET /single":  o.getOverlay,
		"GET /ebuilds": o.getEbuilds,
	}
	return o
}

func (o *OverlaysAPI) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp.Header().Set("Access-Control-Allow-Origin", "*")

	if handler, exists := o.endpoints[req.Method+" "+req.URL.Path[len("/overlays"):]]; exists {
		handler(resp, req)
		return
	}
	if req.Method != http.MethodOptions {
		http.NotFound(resp, req)
	}
}

func (o *OverlaysAPI) getOverlays(resp http.ResponseWriter, req *http.Request) {
	json.NewEncoder(resp).Encode(globals.Overlays)
}

func (o *OverlaysAPI) getOverlay(resp http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	if overlay, ok := globals.Overlays[name]; ok {
		json.NewEncoder(resp).Encode(overlay)
		return
	}

	resp.WriteHeader(http.StatusBadRequest)
}

func (o *OverlaysAPI) getEbuilds(resp http.ResponseWriter, req *http.Request) {
	json.NewEncoder(resp).Encode(globals.Ebuilds)
}
