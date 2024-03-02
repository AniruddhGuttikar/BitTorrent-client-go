package torrentfile

import (
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
)

// type TorrentFile struct {
// 	Announce string
// 	Info
// }

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func ReadFile(f io.Reader) (error) {

	bto := bencodeTorrent{}
	err := bencode.Unmarshal(f, &bto)

	if (err != nil) {
		return err
	}
	fmt.Println(bto.Info.PieceLength)
	return nil
}
