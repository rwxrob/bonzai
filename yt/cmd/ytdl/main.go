package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rwxrob/bonzai/yt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <youtube-url>\n", os.Args[0])
		os.Exit(1)
	}

	url := os.Args[1]

	res, err := yt.Download(yt.DownloadOptions{
		URL:       url,
		OutputDir: ".",
		Timeout:   30 * time.Minute,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("downloaded: %s\n", res.FilePath)
}
