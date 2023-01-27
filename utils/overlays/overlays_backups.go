package overlays

import (
	"encoding/json"
	"os"

	"github.com/mbaraa/eloiserver/config"
	"github.com/mbaraa/eloiserver/models"
)

var overlaysFilePath = config.BackupDirectory() + "/overlays.json"

func LoadOverlays() (map[string]*models.Overlay, error) {
	overlaysFile, err := os.Open(overlaysFilePath)
	defer overlaysFile.Close()
	if err != nil {
		return nil, err
	}

	overlays := make(map[string]*models.Overlay)

	err = json.NewDecoder(overlaysFile).Decode(&overlays)
	if err != nil {
		return nil, err
	}

	return overlays, nil
}

func SaveOverlays(overlays map[string]*models.Overlay) error {
	overlaysFile, err := os.Create(overlaysFilePath)
	defer overlaysFile.Close()
	if err != nil {
		return err
	}

	return json.NewEncoder(overlaysFile).Encode(overlays)
}
