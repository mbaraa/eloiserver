package overlays

import (
	"encoding/json"
	"os"

	"github.com/mbaraa/eloiserver/config"
	"github.com/mbaraa/eloiserver/globals"
	"github.com/mbaraa/eloiserver/models"
)

var overlaysFilePath = config.BackupDirectory() + "/_overlays.json"

func LoadOverlays() error {
	overlaysFile, err := os.Open(overlaysFilePath)
	defer overlaysFile.Close()
	if err != nil {
		return err
	}

	overlays := make(map[string]*models.Overlay)

	err = json.NewDecoder(overlaysFile).Decode(&overlays)
	if err != nil {
		return err
	}

	globals.Overlays = overlays
	globals.Ebuilds = ExtractEbuilds(overlays)
	return nil
}

func SaveOverlays(overlays map[string]*models.Overlay) error {
	overlaysFile, err := os.Create(overlaysFilePath)
	if err != nil {
		return err
	}
	defer overlaysFile.Close()

	return json.NewEncoder(overlaysFile).Encode(overlays)
}
