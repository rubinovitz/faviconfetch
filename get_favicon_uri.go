package faviconfetch

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"regexp"
)

// Main access point. Given a uri, return the favicon
func Fetch(uri string) []byte {
	faviconUri := Detect(uri)
	if faviconUri != "" {
		favicon := GetFavicon(uri, faviconUri)
		return favicon
	} else {
		var favicon []byte
		favicon = nil
		return favicon
	}
}

// Attempt to get the url's HTML
func Detect(uri string) string {
	// remove trailing .ico
	re := regexp.MustCompile(".ico")
	uri = re.ReplaceAllString(uri, "")
	// add http scheme if needed so go url doesn't throw error
	hasHttp, _ := regexp.MatchString("^http", uri)
	if hasHttp != true {
		uri = "http://" + uri
	}
	urlObj, parseErr := url.Parse(uri)
	doc, queryErr := goquery.NewDocument(urlObj.String())
	if parseErr == nil && queryErr == nil {
		log.Print("Parse and query work on first try")
		return FindFaviconUriInHTML(urlObj, doc)
	} else {
		if parseErr != nil {
			log.Print("parseErr != nil")
			scheme := urlObj.Scheme
			if scheme == "" {
				scheme = "http"
			}
			host := urlObj.Host
			newUri := scheme + "://" + host
			urlObj, parseErr := url.Parse(newUri)
			if parseErr == nil {
				doc, queryErr := goquery.NewDocument(urlObj.String())
				if queryErr == nil {
					return FindFaviconUriInHTML(urlObj, doc)
				} else {
					return ""
				}
			} else {
				log.Print("New url did not parse")
				return ""
			}
		} else {
			log.Print("Cannot parse old url")
			return ""
		}
	}
}

// Look for <link rel="icon" and any base url
func FindFaviconUriInHTML(uri *url.URL, doc *goquery.Document) string {
	base, iconUrl := HTMLParserHandler(doc)
	log.Printf("Inside FindFaviconUriInHTML with base %s, iconUrl %s", base, iconUrl)
	// replace urls that start with  // path since go http cannot retrieve them
	re := regexp.MustCompile("^(//)")
	iconUrl = re.ReplaceAllString(iconUrl, uri.Scheme+"://")
	base = re.ReplaceAllString(base, uri.Scheme+"://")
	// if base url and icon url
	if base != "" && iconUrl != "" {
		if base == "/" {
			base = uri.String()
		}
		notRel, _ := regexp.MatchString("^([^/])", iconUrl)
		// make icon url relative pathed if not
		if notRel {
			iconUrl = "/" + iconUrl
		}
		iconUrl = base + iconUrl
		return iconUrl
	}
	// if no base use uri
	if base == "" {

		base = uri.Scheme + "://" + uri.Host

	}

	// if no icon default to checking /favicon.ico
	if iconUrl == "" {
		trailingSlash, _ := regexp.MatchString("/$", base)
		if trailingSlash {
			iconUrl = base + "favicon.ico"
		} else {
			iconUrl = base + "/favicon.ico"
			return iconUrl
		}
	} else {
		// if iconUrl check to make sure its valid
		iconUrlParse, err := url.Parse(iconUrl)
		// if valid, return it
		if err == nil && iconUrlParse.Scheme != "" {
			return iconUrl
			// if invalid, try base and iconUrl
		} else {
			iconUrl = base + iconUrl
			return iconUrl
		}
	}
	return iconUrl
}

// parse the HTML to get the favicon url
func HTMLParserHandler(doc *goquery.Document) (string, string) {
	base := ""
	uri := ""

	// goroutine and channel to look for favicon
	uriChannel := make(chan string)
	go func() {
		doc.Find("link").Each(func(i int, s *goquery.Selection) {
			rel, relExists := s.Attr("rel")
			if relExists == true {
				shortcutIcon, _ := regexp.MatchString("(?i)^(shortcut )?icon$", rel)
				if shortcutIcon == true {
					tagUri, uriExists := s.Attr("href")
					if uriExists == true {
						uriChannel <- tagUri
					}
				}
			}
		})
		if len(uriChannel) == 0 {
			uriChannel <- ""
		}

	}()
	baseChannel := make(chan string)
	go func() {
		doc.Find("base").Each(func(i int, s *goquery.Selection) {
			baseUri, hrefExists := s.Attr("href")
			if hrefExists == true {
				baseChannel <- baseUri
			}
		})
		if len(baseChannel) == 0 {
			baseChannel <- ""
		}
	}()
	base = <-baseChannel
	uri = <-uriChannel
	return base, uri
}
