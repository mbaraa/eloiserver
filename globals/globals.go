package globals

import (
	"github.com/mbaraa/eloiserver/models"
)

var (
	Overlays       map[string]*models.Overlay
	SimpleOverlays map[string]*models.Overlay
	Ebuilds        map[string]map[string]*models.Ebuild
)
