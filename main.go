package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

const (
	IndexPage = "contents/index.html"
	MediaRoot = "testdata"
	MediaFile = "sample-mp4-file.mp4"
	M3u8Name  = "index.m3u8"
)

func main() {
	if err := MP4ToHLS(MediaRoot, MediaFile); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully converted MP4 to HLS.")

	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex).Methods(http.MethodGet)
	router.HandleFunc("/media/stream/", handleStream).Methods(http.MethodGet)
	router.HandleFunc("/media/stream/{segName:index[0-9]+.ts}", handleStream).Methods(http.MethodGet)

	fmt.Println("Media server started at :8080")
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

func MP4ToHLS(root, file string) error {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i",
		file,
		"-profile:v",
		"baseline",
		"-level",
		"3.0",
		"-s",
		"640x360",
		"-start_number",
		"0",
		"-hls_time",
		"10",
		"-hls_list_size",
		"0",
		"-f",
		"hls",
		M3u8Name,
	)
	cmd.Dir = root
	
	return cmd.Run()
}
