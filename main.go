package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AniruddhGuttikar/BitTorrent-client-go/torrentfile"
)

func main() {
	fmt.Println("Enter the input path of the torrentfile and name/location for the downloaded file")
	inPath := os.Args[1]
	outPath := os.Args[1]

	inFile, err := os.Open(inPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	tf, err := torrentfile.ReadFile(inFile)
	if err != nil {
		log.Fatalf("error in reading the file: %v", err)
		return
	}
	err = tf.DownloadFile(outPath)
	if err != nil {
		log.Fatalf("error in the main program: %+v", err)
	}
	fmt.Println("file has been downloaded successfully huhu :3")
}
