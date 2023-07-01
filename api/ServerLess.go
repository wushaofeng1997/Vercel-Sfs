package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CommonResp[E any] struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Data E      `json:"data"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	start := time.Now()
	log.Printf("method: %s ,reqUrl: %s ,rAddr:%s", r.Method, r.RequestURI, r.RemoteAddr)
	type Message struct {
		Content string `json:"content"`
		Id      int    `json:"id"`
	}
	var msg Message
	msg.Content = "test serverless func"
	var resp CommonResp[Message]
	resp.Data = msg
	resp.Code = 2000
	marshal, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.Header().Set("Content-Length", strconv.Itoa(len(marshal)))
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(marshal)
	if err != nil {
		log.Println(err)
	}
	log.Printf("spent %v", time.Since(start))
}
