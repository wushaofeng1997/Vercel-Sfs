package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Handle Serverless Func
func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s path:%s query:%v remote:%s", r.Method, r.RequestURI, r.URL.Query(), r.RemoteAddr)
	log.Printf("%+v", *r.URL)
	r.URL.Query()
	var start = time.Now()
	defer func() {
		log.Printf("method: %s path:%s remote:%s spent:%v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	}()
	err := r.ParseMultipartForm(1024 * 1024 * 64)
	if err != nil {
		//log.Println(err)
	}
	urlStr := os.Getenv("PROXY_URL")
	parse, err := url.Parse(urlStr)
	if err != nil {
		write503(w, err)
		return
	}
	log.Println(r.URL.RawQuery)
	//log.Println(parse.Opaque)
	request, err := http.NewRequest(r.Method, parse.Scheme+"://"+path.Join(parse.Host, r.URL.Path+"?"+r.URL.RawQuery), r.Body)
	if err != nil {
		write503(w, err)
		return
	}
	request.Header = r.Header.Clone()
	request.Form = r.Form
	request.PostForm = r.PostForm
	request.MultipartForm = r.MultipartForm
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		write503(w, err)
		return
	}
	for s, vs := range resp.Header {
		if corsIncludes(s) {
			w.Header().Add(s, strings.Join(vs, "; "))
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func corsIncludes(headerKey string)bool{
	switch headerKey {
		case "Cross-Origin-Opener-Policy":
		    return true
		default: 
		    return false
	}
}

func write503(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(err.Error()))
}
