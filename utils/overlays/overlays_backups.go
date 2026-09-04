package overlays

import (
	"encoding/json"
	"os"

	"codeberg.org/lordbaraa/eloiserver/config"
	"codeberg.org/lordbaraa/eloiserver/globals"
	"codeberg.org/lordbaraa/eloiserver/models"
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
	globals.SimpleOverlays = getSimpleOverlays(overlays)
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
