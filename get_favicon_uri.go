package faviconfetch

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// Main access point. Given a uri, return the favicon
func Fetch(uri string) []byte {
	faviconUri := Detect(uri)
	favicon := GetFavicon(uri, faviconUri)
	return favicon
}

// Attempt to get the url's HTML
func Detect(uri string) string {
	req := SetHTTPHeaders(uri)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		// get contents of page
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		return FindFaviconUriInHTML(uri, string(contents))

	} else {
		if strings.Contains(uri, "://www") != true {
			uriStruct, err := url.Parse(uri)
			if err != nil {
				panic(err)
			}
			scheme := uriStruct.Scheme
			host := uriStruct.Host
			newUri := scheme + "://www." + host
			return Detect(newUri)
		} else {
			return ""
		}
	}
}

// Look for <link rel="icon"
func FindFaviconUriInHTML(uri string, contents string) string {
	base, iconUrl := HTMLParserHandler(contents)
	if base != "" {
		iconUrl = base + iconUrl
	} else {
		iconUrl = uri + iconUrl
	}
	if iconUrl == "" && base != "" {
		iconUrl = uri + "/favicon.ico"
		return iconUrl
	} else {
		urlMatch, _ := regexp.MatchString("^https?://", iconUrl)
		if urlMatch {
			return iconUrl
		} else {
			iconUrl = "http://" + iconUrl
			return iconUrl
		}
	}
}

// parse the HTML to get the favicon url
func HTMLParserHandler(contents string) (string, string) {
	base := ""
	uri := ""
	href := ""
	d := html.NewTokenizer(strings.NewReader(contents))

	for {
		if base != "" && uri != "" {
			return base, uri
		}
		// token type
		tokenType := d.Next()
		if tokenType == html.ErrorToken {

		}
		token := d.Token()
		switch tokenType {
		case html.StartTagToken:
			// get link rel tag href value if not shortcut icon
			if (token.Data) == "link" {
				isRel := false
				for i := range token.Attr {
					key := token.Attr[i].Key
					val := token.Attr[i].Val
					shortcutIcon, _ := regexp.MatchString("^(shortcut )?icon$/i", val)
					if key == "rel" && shortcutIcon == false {
						isRel = true
					}
					if key == "href" {
						href = val
					}
				}
				if isRel {
					uri = href
				}
			}

			// get base url if exists
			if (token.Data) == "base" {
				for i := range token.Attr {
					key := token.Attr[i].Key
					val := token.Attr[i].Val
					if key == "href" {
						base = val
					}
				}
			}
		case html.ErrorToken:
			return base, uri
		}
	}

	return base, uri
}
