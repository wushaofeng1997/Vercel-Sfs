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
	log.Printf("method: %s path:%s remote:%s", r.Method, r.RequestURI, r.RemoteAddr)
	var start = time.Now()
	defer func() {
		log.Printf("method: %s path:%s remote:%s spent:%v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	}()
	parse, err := url.Parse(os.Getenv("PROXY_URL"))
	if err != nil {
		write503(w, err)
		return
	}
	log.Println(parse.Scheme)
	//log.Println(parse.Opaque)
	request, err := http.NewRequest(r.Method, parse.Scheme+"://"+path.Join(parse.Host, r.RequestURI), r.Body)
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
		w.Header().Add(s, strings.Join(vs, "; "))
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func write503(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(err.Error()))
}
