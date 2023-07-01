package api

import (
	"log"
	"net/http"
	"time"
)

// Handle Serverless Func
func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s path:%s query:%v remote:%s", r.Method, r.RequestURI, r.URL.Query(), r.RemoteAddr)
	var start = time.Now()
	defer func() {
		log.Printf("method: %s path:%s remote:%s spent:%v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	}()
	w.Write([]byte(r.RequestURI))
}

func write503(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(err.Error()))
}
