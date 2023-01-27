package overlays

import (
	"github.com/mbaraa/eloiserver/models"
)

func ExtractEbuilds(overlays map[string]*models.Overlay) map[string]map[string]*models.Ebuild {
	ebuilds := make(map[string]map[string]*models.Ebuild)

	for _, overlay := range overlays {
		for _, group := range overlay.EbuildGroups {
			for _, ebuild := range group.Ebuilds {
				fullName := ebuild.GroupName + "/" + ebuild.Name
				if _, ok := ebuilds[fullName]; !ok {
					ebuilds[fullName] = make(map[string]*models.Ebuild)
				}
				ebuilds[fullName][ebuild.Version] = ebuild
			}
		}
	}

	return ebuilds
}
