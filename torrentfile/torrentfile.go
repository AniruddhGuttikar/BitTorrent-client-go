package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/AniruddhGuttikar/BitTorrent-client-go/p2p"
	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string     // url to which request has to be made
	InfoHash    [20]byte   // unique identifier of the file (hash of info object)
	PieceLength int        // size of one piece
	PieceHashes [][20]byte // list of 20 byte SHA1 of corresponding pieces (len(PieceHashes) == Length / PieceLength)
	Length      int        // size of the file to be downloaded
	Name        string     // name of the file to be downloaded
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

	fmt.Printf("The .torrent file has been read: Announce=%s,\n Name=%s,\n Length=%d,\n PieceLength=%d\n", bto.Announce, bto.Info.Name, bto.Info.Length, bto.Info.PieceLength)

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
	peerID := createPeerId()
	fmt.Println("Your PeerID is: ", string(peerID))

	peers, err := t.requestPeers([20]byte(peerID))
	if err != nil {
		return err
	}
	for _, peer := range peers {
		fmt.Println("peer data: ", peer.String())
	}

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      [20]byte(peerID),
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
	}

	buf, err := torrent.Download()
	if err != nil {
		return err
	}

	outFile, err := os.Create(f)
	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

// find the hash of the info key of a parsed bencode
// bencode the info and than hash it with SHA1 (20 bytes)
func (bi *bencodeInfo) hash() ([20]byte, error) {
	var buff bytes.Buffer
	err := bencode.Marshal(&buff, *bi)
	if err != nil {
		return [20]byte{}, err
	}
	hash := sha1.Sum(buff.Bytes())
	return hash, nil
}

// converts from type bencodeTorrent to type TorrentFile
func (bto *bencodeTorrent) toTorrent() (TorrentFile, error) {
	var tf = TorrentFile{}

	// infoHash uniquely identifies the .torrent file
	// it is the SHA1 of the info of bencodeTorrent
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	// pieces is a string with length multiple of 20
	// each blob of 20 bytes is the SHA1 hash of the corresponding piece
	const hashSize = 20

	if rem := len(bto.Info.Pieces) % hashSize; rem != 0 {
		err := fmt.Errorf("malformed pieces of length: %d", len(bto.Info.Pieces))
		return TorrentFile{}, err
	}

	// pieces hashes is an array containing these SHA1 hashes
	numHashes := len(bto.Info.Pieces) / hashSize

	// fmt.Println("num hashes: ", numHashes)
	// fmt.Println("my theory: ", bto.Info.Length/bto.Info.PieceLength)
	// if numHashes == bto.Info.Length/bto.Info.PieceLength {
	// 	fmt.Println("HEHEHE ITS RIGHT")
	// }

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

func createPeerId() []byte {
	peerID := make([]byte, 20)

	if _, err := rand.Read(peerID); err != nil {
		fmt.Println("Error in creating the peer Id")
	}
	copy(peerID[:2], "BZ")
	return peerID
}
