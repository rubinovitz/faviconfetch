package faviconfetch

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"regexp"
	"strings"
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
	uri = strings.Replace(uri, ".ico", "", -1)
	// add http:// to url if not there, otherwise go url does not recognize it as an url
	if strings.Contains(uri, "http") != true {
		uri = "http://" + uri
	}
	urlObj, parseErr := url.Parse(uri)
	doc, queryErr := goquery.NewDocument(urlObj.String())
	if parseErr == nil && queryErr == nil {
		return FindFaviconUriInHTML(urlObj, doc)
	}
	// there's an error with the URL format
	if parseErr != nil {
		uriStruct, err := url.Parse(uri)
		if err != nil {
			log.Fatal(err)
		}
		scheme := uriStruct.Scheme
		host := uriStruct.Host
		newUri := scheme + "://" + host
		return Detect(newUri)
	}
	// there's an error reaching the site
	return ""
}

// Look for <link rel="icon" and any base url
func FindFaviconUriInHTML(uri *url.URL, doc *goquery.Document) string {
	base, iconUrl := HTMLParserHandler(doc)
	// replace // path since go HTTP cannot retrieve them
	re := regexp.MustCompile("^(//)")
	iconUrl = re.ReplaceAllString(iconUrl, uri.Scheme+"://")
	base = re.ReplaceAllString(base, uri.Scheme+"://")
	// if iconUrl is not relative, make it so
	notRel, _ := regexp.MatchString("^([^/])", iconUrl)
	notHTTP, _ := regexp.MatchString("^([^http://])", iconUrl)

	if notRel && notHTTP {
		iconUrl = "/" + iconUrl
	}
	if base != "" && iconUrl != "" {
		iconUrl = base + iconUrl
		return iconUrl
	}
	if base == "" {
		base = uri.Host
	}
	if iconUrl == "" {
		iconUrl = uri.Scheme + "://" + base + "/favicon.ico"
		return iconUrl
	} else {
		iconUrlParse, err := url.Parse(iconUrl)
		if err == nil && iconUrlParse.Scheme != "" {
			return iconUrl
		} else {
			iconUrl = uri.Scheme + "://" + base + iconUrl
			return iconUrl
		}
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
