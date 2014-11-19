package faviconfetch

import (
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

//Test root url/test one
func templateHandlerOne(w http.ResponseWriter, r *http.Request) {
	s1, err := template.ParseFiles("templates/base.tmpl", "templates/test1.tmpl")
	if err != nil {
		panic(err)
	}
	s1.ExecuteTemplate(w, "base", nil)
}

func serveFileHandlerOne(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./favicons/ddgfavicon.ico")
}

func HandlerOne() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", templateHandlerOne)
	r.HandleFunc("/favicon.ico", serveFileHandlerOne)
	return r

}

func HandlerTwo() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", templateHandlerTwo)
	r.HandleFunc("/static/favicon.ico", serveFileHandlerTwo)
	return r
}

type Path struct {
	Url string
}

//Test base url/test two
func templateHandlerTwo(w http.ResponseWriter, r *http.Request) {
	t := template.New("test2")
	t, err := template.ParseFiles("templates/test2.tmpl")
	if err != nil {
		panic(err)
	}
	data := &Path{Url: "http://" + r.Host}
	t.Execute(w, data)
}

func serveFileHandlerTwo(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./favicons/ddgfavicon.ico")
}
