package api

import "net/http"

// Handle Serverless Func
func Handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %s path:%s remote:%s", r.Method, r.RequestURI, r.RemoteAddr)
	var start = time.Now()
	defer func() {
		log.Printf("method: %s path:%s remote:%s spent:%v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
	}()
	getenv := os.Getenv("PROXY_URL")
	parse, err := url.Parse(getenv)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(err.Error()))
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parse)
	proxy.ServeHTTP(w, r)
}
