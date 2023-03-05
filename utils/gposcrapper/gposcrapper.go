/*
Package gposcrapper scrapes through http://gpo.zugaina.org/Overlays
and retrieves all available Gentoo ebuilds from different overlays
*/
package gposcrapper

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/mbaraa/eloiserver/models"

	"github.com/gocolly/colly"
)

var (
	baseURL                 = "http://gpo.zugaina.org/Overlays"
	allowedCharsInURLRegExp = `[a-zA-Z0-9-\._~:\?#\[\]@!\$&'\(\)\*\+,;%=]+`
	overlayLinkRegExp       = regexp.MustCompile(`http:\/\/gpo\.zugaina\.org\/Overlays\/` + allowedCharsInURLRegExp + "$")
	ebuildGroupLinkRegExp   = regexp.MustCompile(fmt.Sprintf(`http:\/\/gpo\.zugaina\.org\/Overlays\/%s\/%s$`, allowedCharsInURLRegExp, allowedCharsInURLRegExp))
	ebuildLinkRexExp        = regexp.MustCompile(fmt.Sprintf(`http:\/\/gpo\.zugaina\.org\/Overlays\/%s\/%s\/%s$`, allowedCharsInURLRegExp, allowedCharsInURLRegExp, allowedCharsInURLRegExp))
)

func GetOverlays() (map[string]*models.Overlay, error) {
	overlays, err := getOverlays(baseURL)
	if err != nil {
		return nil, err
	}

	return mergeMetadataWithOriginalEbuildsData(overlays)
}

func GetOverlay(overlayName string) (*models.Overlay, error) {
	overlays, err := getOverlays(baseURL + "/" + overlayName)
	if err != nil {
		return nil, err
	}

	overlays2, err := mergeMetadataWithOriginalEbuildsData(overlays)
	if err != nil {
		return nil, err
	}

	return overlays2[overlayName], nil
}

func GetOverlaysMetadata() (map[string]*models.Overlay, error) {
	return getOverlaysMetadata()
}

type ConcurrentMap struct {
	mu       sync.RWMutex
	overlays map[string]*models.Overlay
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		overlays: make(map[string]*models.Overlay),
	}
}

func (c *ConcurrentMap) CreateOverlay(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.overlays[name]; ok {
		return
	}
	c.overlays[name] = &models.Overlay{EbuildGroups: make(map[string]*models.EbuildGroup)}
}

func (c *ConcurrentMap) CreateEbuildGroup(overlayName, groupName string) {
	c.CreateOverlay(overlayName)
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.overlays[overlayName].EbuildGroups[groupName]; ok {
		return
	}
	c.overlays[overlayName].EbuildGroups[groupName] = &models.EbuildGroup{Name: groupName, Ebuilds: make(map[string]*models.Ebuild)}
}

func (c *ConcurrentMap) CreateEbuild(overlayName, groupName, ebuildName string, ebuild *models.Ebuild) {
	c.CreateEbuildGroup(overlayName, groupName)
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.overlays[overlayName].EbuildGroups[groupName].Ebuilds[ebuildName]; ok {
		return
	}

	c.overlays[overlayName].EbuildGroups[groupName].Ebuilds[ebuildName] = ebuild
}

func (c *ConcurrentMap) GetOverlays() map[string]*models.Overlay {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.overlays
}

func getOverlays(overlaysURL string) (map[string]*models.Overlay, error) {
	overlays := NewConcurrentMap()

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
	)

	noSSL := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.WithTransport(noSSL)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})

	c.OnHTML("table.usetable", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			currentOverlayName := ""

			tr.ForEach("td", func(i int, td *colly.HTMLElement) {
				// 0 => Name
				// 1 => Description
				// 2 => NumEbuilds
				// 3 => Homepage
				// 4 => Feed
				// 5 => Mail
				// 6 => Source

				switch i {
				case 0:
					link := td.Request.AbsoluteURL(td.ChildAttr("a", "href"))
					c.Visit(link)
					currentOverlayName = td.Text
					overlays.CreateOverlay(currentOverlayName)
					overlays.GetOverlays()[currentOverlayName].Name = td.Text
				}
			})
		})
	})

	c.OnHTML("#contentInner", func(e *colly.HTMLElement) {
		match := overlayLinkRegExp.MatchString(e.Request.URL.String())
		if !match {
			return
		}

		e.ForEach("a", func(i int, a *colly.HTMLElement) {
			link := a.Request.AbsoluteURL(a.Attr("href"))
			c.Visit(link)
		})
	})

	c.OnHTML("#browsebox", func(e *colly.HTMLElement) {
		match := ebuildGroupLinkRegExp.MatchString(e.Request.URL.String())
		if !match {
			return
		}

		e.ForEach("a", func(i int, a *colly.HTMLElement) {
			link := a.Request.AbsoluteURL(a.Attr("href"))
			c.Visit(link)

			url := e.Request.URL.String()
			noHost := url[len("http://gpo.zugaina.org/Overlays/"):]
			overlayName := noHost[:strings.Index(noHost, "/")]
			groupName := noHost[strings.Index(noHost, "/")+1:]

			overlays.CreateEbuildGroup(overlayName, groupName)
			overlays.GetOverlays()[overlayName].EbuildGroups[groupName].Link = url
		})
	})

	c.OnHTML("#contentInner", func(e *colly.HTMLElement) {
		match := ebuildLinkRexExp.MatchString(e.Request.URL.String())
		if !match {
			return
		}

		url := e.Request.URL.String()
		noHost := url[len("http://gpo.zugaina.org/Overlays/"):]
		overlayName := noHost[:strings.Index(noHost, "/")]
		groupName := noHost[strings.Index(noHost, "/")+1 : strings.LastIndex(noHost, "/")]
		ebuildName := noHost[strings.LastIndex(noHost, "/")+1:]

		ebuildDescription := ""
		e.ForEach("h5", func(i int, h5 *colly.HTMLElement) {
			ebuildDescription = h5.Text
		})

		e.ForEach("#"+overlayName, func(i int, div *colly.HTMLElement) {
			div.ForEach("li", func(i int, li *colly.HTMLElement) {
				text := li.Text
				license := text[strings.Index(text, "License"):strings.LastIndex(text, "   ")]
				license = license[len("License: "):]
				ebuild := models.Ebuild{}
				li.ForEach("div,a,br", func(i int, div *colly.HTMLElement) {
					// 0 => Name-Version
					// 1 => Architecture
					// 2 => cpu_flags
					// 3 => IGNORE
					// 4 => ebuild view link
					// 5 => ebuild download link
					// 6 => ebuild browse
					// 7 => overlay name
					switch i {
					case 0:
						version := div.Text[len(ebuildName)+1:]

						ebuild.Name = ebuildName
						ebuild.Version = version
					case 1:
						ebuild.Architecture = div.Text
					case 2:
						ebuild.Flags = div.Text
					case 4:
						ebuild.Homepage = div.ChildAttr("a", "href")
					case 7:
						ebuild.OverlayName = div.Text[len("Overlay: "):]
					}
				})
				ebuild.License = license
				ebuild.GroupName = groupName
				ebuild.Description = ebuildDescription
				overlays.CreateEbuild(overlayName, groupName, ebuild.Name+"-"+ebuild.Version, &ebuild)
			})
		})
	})

	err := c.Visit(overlaysURL)
	if err != nil {
		return nil, err
	}
	c.Wait()

	return overlays.GetOverlays(), nil
}

func mergeMetadataWithOriginalEbuildsData(overlays map[string]*models.Overlay) (map[string]*models.Overlay, error) {
	overlaysWithCorrectMetadata, err := getOverlaysMetadata()
	if err != nil {
		return nil, err
	}

	for name, overlay := range overlaysWithCorrectMetadata {
		if _, ok := overlays[name]; ok {
			overlay.EbuildGroups = overlays[overlay.Name].EbuildGroups
		}
	}

	return overlaysWithCorrectMetadata, nil
}

func getOverlaysMetadata() (map[string]*models.Overlay, error) {
	resp, err := http.Get("http://gpo.zugaina.org/lst/layman-repositories.xml")
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)

	var repos struct {
		Repos []*models.Overlay `xml:"repo"`
	}

	err = xml.Unmarshal(bodyBytes, &repos)
	if err != nil {
		return nil, err
	}

	overlays := make(map[string]*models.Overlay)

	for _, overlay := range repos.Repos {
		name := overlay.Name
		overlays[name] = overlay
		overlays[name].EbuildGroups = make(map[string]*models.EbuildGroup)
	}

	return overlays, nil
}
