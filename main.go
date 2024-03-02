package main

import (
	"fmt"

	"github.com/AniruddhGuttikar/BitTorrent-client-go/torrentfile"
)

func main() {
	str, err := torrentfile.ReadFile("./archlinux-2019.12.01-x86_64.iso.torrent")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)
	
}
