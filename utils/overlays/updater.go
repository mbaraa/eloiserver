package overlays

import (
	"fmt"

	"github.com/mbaraa/eloiserver/globals"
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
	fmt.Println("Updating overlays cache from gpo.zugania.org...")

	globals.Overlays = gposcrapper.GetOverlays()
	globals.Ebuilds = ExtractEbuilds(globals.Overlays)

	err := SaveOverlays(globals.Overlays)
	if err != nil {
		return err
	}

	fmt.Println("All done ✓")
	return nil
}
