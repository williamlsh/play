package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	IndexPage = "contents/index.html"
	MediaRoot = "testdata"
	M3u8Name  = "index.m3u8"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex).Methods(http.MethodGet)
	router.HandleFunc("/media/stream/", handleStream).Methods(http.MethodGet)
	router.HandleFunc("/media/stream/{segName:index[0-9]+.ts}", handleStream).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, IndexPage)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	segName, ok := vars["segName"]

	var mediaFile string
	if !ok {
		mediaFile = fmt.Sprintf("%s/%s", MediaRoot, M3u8Name)
		w.Header().Set("Content-Type", "application/x-mpegURL")
	} else {
		mediaFile = fmt.Sprintf("%s/%s", MediaRoot, segName)
		w.Header().Set("Content-Type", "video/MP2T")
	}
	fmt.Printf("Incomming request for media file: %s\n", mediaFile)

	http.ServeFile(w, r, mediaFile)
}
