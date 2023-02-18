package overlays

import (
	"log"

	"github.com/mbaraa/eloiserver/globals"
	"github.com/mbaraa/eloiserver/models"
	"github.com/mbaraa/eloiserver/utils/gposcrapper"
	"github.com/robfig/cron"
)

func ScheduleScrapper() error {
	cronie := cron.New()

	var err2 error
	err := cronie.AddFunc("0 0 0 * * *", func() {
		err2 = ScrapeOverlays()
	})
	if err != nil {
		return err
	}

	if err2 != nil {
		return err2
	}

	cronie.Start()

	return nil
}

func ScrapeOverlays() error {
	log.Println("Updating overlays cache from gpo.zugania.org...")

	var err error
	globals.Overlays, err = gposcrapper.GetOverlays()
	if err != nil {
		return err
	}
	globals.Ebuilds = ExtractEbuilds(globals.Overlays)
	globals.SimpleOverlays = getSimpleOverlays(globals.Overlays)

	err = SaveOverlays(globals.Overlays)
	if err != nil {
		return err
	}

	log.Println("All done ✓")
	return nil
}

func getSimpleOverlays(overlays map[string]*models.Overlay) (simple map[string]*models.Overlay) {
	simple = make(map[string]*models.Overlay)

	for name, overlay := range overlays {
		simple[name] = &models.Overlay{
			Name:        name,
			Description: overlay.Description,
			Homepage:    overlay.Homepage,
			Feed:        overlay.Feed,
			Owner:       overlay.Owner,
			Source:      overlay.Source[:],
		}
	}

	return
}
