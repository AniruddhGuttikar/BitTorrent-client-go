package main

import (
	"os"
	"log"
	"github.com/AniruddhGuttikar/BitTorrent-client-go/torrentfile"
)

func main() {
	inPath := os.Args[1]
	// outPath := os.Args[1]
	
	inFile, err := os.Open(inPath)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = torrentfile.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	
}
