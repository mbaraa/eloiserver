package globals

import (
	"codeberg.org/lordbaraa/eloiserver/models"
)

var (
	Overlays       map[string]*models.Overlay
	SimpleOverlays map[string]*models.Overlay
	Ebuilds        map[string]map[string]*models.Ebuild
)
