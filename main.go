package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
  "strings"
)

type handler struct{}

func serve(filename string, w http.ResponseWriter) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	rdr := bufio.NewReader(file)
	firstChunk, err := rdr.ReadString(0)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Fprint(w, firstChunk)
}

func serveChunks(w http.ResponseWriter) {
	serve("first_chunk.txt", w)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	time.Sleep(500 * time.Millisecond)
	serve("second_chunk.txt", w)
}


func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if strings.HasSuffix(r.URL.Path, ".css") {
    w.Header()["Content-Type"] = []string{"text/css; charset=utf-8"}
    time.Sleep(500 * time.Millisecond)
    fmt.Fprint(w, "body { background-color: green }")
  } else {
    serveChunks(w)
  }
}

func main() {
	var myHandler handler
	s := &http.Server{
		Addr:           ":8080",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 0,
	}
	log.Fatal(s.ListenAndServe())
}
