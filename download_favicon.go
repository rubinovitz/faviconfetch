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
	DumpResponse(resp)
	if err != nil {
		fmt.Println(err)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	matchHtml, err := regexp.MatchString("<html/i", string(contents))
	if !matchHtml {
		contentType := http.DetectContentType(contents)
		matchZip, err := regexp.MatchString("zip/i", contentType)
		if err != nil {
			panic(err)
		}
		if !matchZip {
			favicon := contents
			return favicon
		} else {
			// question: will it only ever be gzip?
			faviconGzipReader, err := gzip.NewReader(req.Body)
			faviconContents, err := gzip.NewReader(faviconGzipReader)
			favicon, err := ioutil.ReadAll(faviconContents)
			return favicon
			if err != nil {
				panic(err)
			}
		}
	} else {
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
