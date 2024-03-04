package torrentfile

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/AniruddhGuttikar/BitTorrent-client-go/peers"
	"github.com/jackpal/bencode-go"
)

type bencodedTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

const PORT = 6881

func (t *TorrentFile) buildTrackerURL(peerID [20]byte) (string, error) {
	baseLink, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("info_hash", string(t.InfoHash[:]))
	params.Set("peer_id", string(peerID[:]))
	params.Set("port", strconv.Itoa(PORT))
	params.Set("uploaded", "0")
	params.Set("downloaded", "0")
	params.Set("compact", "1")
	params.Set("left", strconv.Itoa(t.Length))

	baseLink.RawQuery = params.Encode()
	return baseLink.String(), nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte) ([]peers.Peer, error) {
	trackerUrl, err := t.buildTrackerURL(peerID)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	benRes, err := client.Get(trackerUrl)
	if err != nil {
		return nil, err
	}
	// response is a bencoded value
	response := bencodedTrackerResponse{}
	err = bencode.Unmarshal(benRes.Body, &response)
	if err != nil {
		return nil, err
	}
	p, err := peers.UnmarshalPeers([]byte(response.Peers))
	if err != nil {
		return nil, err
	}
	return p, nil
}
