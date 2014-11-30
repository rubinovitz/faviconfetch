package faviconfetch

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// download favicon
func GetFavicon(uri string, faviconUri string) []byte {
<<<<<<< HEAD
=======
	fmt.Printf("inside GetFavicon with uri %s, faviconUri %s\n", uri, faviconUri)
>>>>>>> parent of 47d8628... url parsing fixes
	if faviconUri == "" {
		return nil
	}
	var favicon string
	// remove any new line characters from url
	faviconUri = strings.Replace(faviconUri, "%20", "", -1)
	client := &http.Client{}
	// get request for favicon
	req := SetHTTPHeaders(faviconUri)
	resp, err := client.Do(req)
	if os.Getenv("DEBUG") != "" {
		DumpResponse(resp)
	}
	if err != nil {
		fmt.Println(err)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("contents error: %s", err)
	}
	matchHtml, err := regexp.MatchString("(?i)<html>", string(contents))
	if !matchHtml && string(contents) != "" {
		contentType := http.DetectContentType(contents)
		matchZip, err := regexp.MatchString("(?i)zip", contentType)
		if err != nil {
			panic(err)
		}
		if !matchZip {
			favicon := contents
			return favicon
		} else {
			faviconGzipReader, err := gzip.NewReader(req.Body)
			faviconContents, err := gzip.NewReader(faviconGzipReader)
			favicon, err := ioutil.ReadAll(faviconContents)
			return favicon
			if err != nil {
				panic(err)
			}
		}
	} else {
		return nil
	}
	if favicon == "" {
		return nil
	}
	new_uri := faviconUri + "/favicon.ico"
	if faviconUri != new_uri {
		return GetFavicon(uri, new_uri)
	}
	fmt.Println("Failed to download favicon")
	// failed to download favicon, give up.
	return nil
}

func DumpResponse(resp *http.Response) {
	fmt.Println("\n")
	fmt.Printf("Url: %s\n", resp.Request.URL)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers: %s\n", resp.Header)
	//fmt.Printf("Body: %s\n", resp.Body)

}
