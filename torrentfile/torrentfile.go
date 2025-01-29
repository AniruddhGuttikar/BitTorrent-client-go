package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceLength int
	PieceHashes [][20]byte
	Length      int
	Name        string
}

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

func ReadFile(f io.Reader) (TorrentFile, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(f, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	Torrent, err := bto.toTorrent()
	if err != nil {
		return TorrentFile{}, err
	}
	// dst := make([]byte, hex.EncodedLen(len(Torrent.InfoHash)))
	// hex.Encode(dst, Torrent.InfoHash[:])
	//fmt.Printf("%+v\n", Torrent)
	return Torrent, nil
}

// function which will download the file into desired location
func (t *TorrentFile) DownloadFile(f string) error {
	//create a random peerID
	peerID := make([]byte, 20)
	if _, err := rand.Read(peerID); err != nil {
		return err
	}
	copy(peerID[:2], "BZ")
	fmt.Println("Your PeerID is: ", peerID)
	peers, err := t.requestPeers([20]byte(peerID))
	if err != nil {
		return err
	}
	for _, peer := range peers {
		fmt.Println("peer data: ", peer.String())
	}
	return nil
}

// find the hash of the info key of a parsed bencode
// bencode the info and than hash it with SHA1 (20 bytes)
func (i *bencodeInfo) hash() ([20]byte, error) {
	var buff bytes.Buffer
	err := bencode.Marshal(&buff, *i)
	if err != nil {
		return [20]byte{}, err
	}
	hash := sha1.Sum(buff.Bytes())
	return hash, nil
}

// converts from type bencodeTorrent to type TorrentFile
func (bto *bencodeTorrent) toTorrent() (TorrentFile, error) {
	var tf = TorrentFile{}

	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	//pieces is a string with length multiple of 20
	//each blob of 20 bytes is the SHA1 hash of the corresponding piece
	const hashSize = 20
	if rem := len(bto.Info.Pieces) % hashSize; rem != 0 {
		err := fmt.Errorf("malformed pieces of length: %d", len(bto.Info.Pieces))
		return TorrentFile{}, err
	}

	//pieces hashes is an array containing these SHA1 hashes
	numHashes := len(bto.Info.Pieces) / hashSize
	buf := []byte(bto.Info.Pieces)
	tempHash := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(tempHash[i][:], buf[i*hashSize:(i+1)*hashSize])
	}

	tf.InfoHash = infoHash
	tf.Announce = bto.Announce
	tf.Length = bto.Info.Length
	tf.Name = bto.Info.Name
	tf.PieceLength = bto.Info.PieceLength
	tf.PieceHashes = tempHash

	return tf, nil
}
