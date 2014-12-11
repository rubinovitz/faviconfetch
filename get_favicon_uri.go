package faviconfetch

import (
	"github.com/PuerkitoBio/goquery"
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
		return FindFaviconUriInHTML(urlObj, doc)
	} else {
		if parseErr != nil {
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
				return ""
			}
		} else {
			return ""
		}
	}
}

// Look for <link rel="icon" and any base url
func FindFaviconUriInHTML(uri *url.URL, doc *goquery.Document) string {
	base, iconUrl := HTMLParserHandler(doc)
	// replace urls that start with  // path since go http cannot retrieve them
	re := regexp.MustCompile("^(//)")
	iconUrl = re.ReplaceAllString(iconUrl, uri.Scheme+"://")
	base = re.ReplaceAllString(base, uri.Scheme+"://")
	relIconUrl, _ := regexp.MatchString("^/", iconUrl)
	if base == "" {
		base = uri.String()
	}
	// remove trailing base slash
	trailingSlash := regexp.MustCompile("/$")
	base = trailingSlash.ReplaceAllString(base, "")

	if iconUrl == "" {
		return base + "/favicon.ico"

	} else {
		parseIconUrl, err := url.Parse(iconUrl)

		if err == nil && parseIconUrl.Host != "" {
			if parseIconUrl.Scheme == "" {
				if relIconUrl {
					iconUrl = uri.Scheme + "://" + base + iconUrl
				} else {
					iconUrl = uri.Scheme + "://" + base + "/" + iconUrl

				}
			} else {
				// iconUrl has a scheme
				return iconUrl

			}
		}

		return ""

	}

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
