/*
Package gposcrapper scrapes through http://gpo.zugaina.org/Overlays
and retrieves all available Gentoo ebuilds from different overlays
*/
package gposcrapper

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

func GetOverlays() map[string]*models.Overlay {
	return getOverlays(baseURL)
}

func GetOverlay(overlayName string) *models.Overlay {
	return getOverlays(baseURL + "/" + overlayName)[overlayName]
}

func GetOverlaysMetadata() map[string]*models.Overlay {
	return getOverlaysMetadata(baseURL)
}

func getOverlays(overlaysURL string) map[string]*models.Overlay {
	overlays := make(map[string]*models.Overlay)

	c := getColly()

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
					createOverlay(overlays, currentOverlayName)
					overlays[currentOverlayName].Name = td.Text
					overlays[currentOverlayName].URL = link
				case 1:
					overlays[currentOverlayName].Description = td.Text
				case 2:
					num, _ := strconv.ParseInt(td.Text, 10, 32)
					overlays[currentOverlayName].NumEbuilds = int(num)
				case 3:
					overlays[currentOverlayName].Homepage = td.ChildAttr("a", "href")
				case 4:
					overlays[currentOverlayName].Feed = td.ChildAttr("a", "href")
				case 5:
					overlays[currentOverlayName].Mail = td.ChildAttr("a", "href")
				case 6:
					overlays[currentOverlayName].Source = td.Text
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

			createEbuildGroup(overlays, overlayName, groupName)
			overlays[overlayName].EbuildGroups[groupName].Link = url
		})
	})

	c.OnHTML("#ebuild_list", func(e *colly.HTMLElement) {
		match := ebuildLinkRexExp.MatchString(e.Request.URL.String())
		if !match {
			return
		}

		url := e.Request.URL.String()
		noHost := url[len("http://gpo.zugaina.org/Overlays/"):]
		overlayName := noHost[:strings.Index(noHost, "/")]
		groupName := noHost[strings.Index(noHost, "/")+1 : strings.LastIndex(noHost, "/")]
		ebuildName := noHost[strings.LastIndex(noHost, "/")+1:]

		e.ForEach("#"+overlayName, func(i int, div *colly.HTMLElement) {
			div.ForEach("li", func(i int, li *colly.HTMLElement) {
				text := li.Text
				license := text[strings.Index(text, "License"):strings.LastIndex(text, "   ")]
				license = license[len("License: "):]
				//				fmt.Println("License", license)
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
						ebuild.OverlayName = div.Text
					}
				})
				ebuild.License = license
				ebuild.GroupName = groupName
				createEbuild(overlays, overlayName, groupName, ebuild.Name+"-"+ebuild.Version, &ebuild)
			})
		})
	})

	c.Visit(overlaysURL)

	return overlays
}

func getOverlaysMetadata(overlaysURL string) map[string]*models.Overlay {
	overlays := make(map[string]*models.Overlay)

	c := getColly()

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
					currentOverlayName = td.Text
					createOverlay(overlays, currentOverlayName)
					overlays[currentOverlayName].Name = td.Text
					overlays[currentOverlayName].URL = link
				case 1:
					overlays[currentOverlayName].Description = td.Text
				case 2:
					num, _ := strconv.ParseInt(td.Text, 10, 32)
					overlays[currentOverlayName].NumEbuilds = int(num)
				case 3:
					overlays[currentOverlayName].Homepage = td.ChildAttr("a", "href")
				case 4:
					overlays[currentOverlayName].Feed = td.ChildAttr("a", "href")
				case 5:
					overlays[currentOverlayName].Mail = td.ChildAttr("a", "href")
				case 6:
					overlays[currentOverlayName].Source = td.Text
				}
			})
		})
	})

	c.Visit(overlaysURL)

	return overlays
}

func createOverlay(overlays map[string]*models.Overlay, name string) {
	if _, ok := overlays[name]; ok {
		return
	}
	overlays[name] = &models.Overlay{EbuildGroups: make(map[string]*models.EbuildGroup)}
}

func createEbuildGroup(overlays map[string]*models.Overlay, overlayName, groupName string) {
	createOverlay(overlays, overlayName)

	if _, ok := overlays[overlayName].EbuildGroups[groupName]; ok {
		return
	}
	overlays[overlayName].EbuildGroups[groupName] = &models.EbuildGroup{Name: groupName, Ebuilds: make(map[string]*models.Ebuild)}
}

func createEbuild(overlays map[string]*models.Overlay, overlayName, groupName, ebuildName string, ebuild *models.Ebuild) {
	createEbuildGroup(overlays, overlayName, groupName)

	if _, ok := overlays[overlayName].EbuildGroups[groupName].Ebuilds[ebuildName]; ok {
		return
	}

	overlays[overlayName].EbuildGroups[groupName].Ebuilds[ebuildName] = ebuild
}

func getColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	noSSL := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.WithTransport(noSSL)

	return c
}
