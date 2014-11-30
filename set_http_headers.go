package faviconfetch

import (
	"net/http"
)

func SetHTTPHeaders(uri string) (*http.Request, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return req, err
	}
	// add headers to request
	req.Header.Add("Accept-Language", "en-US,en;q=0.8,zh;q=0.6,es;q=0.4")
	req.Header.Add("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; DuckDuckGoBot/1.0; +http://duckduckgo.com)")
	return req, err
}
