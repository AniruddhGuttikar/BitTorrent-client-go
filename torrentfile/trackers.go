package torrentfile

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	fmt.Printf("torrent file built: %+v\n", t)
	fmt.Printf("torrent file succsfully %s\n", params.Encode())

	baseLink.RawQuery = params.Encode()
	return baseLink.String(), nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte) ([]peers.Peer, error) {
	trackerUrl, err := t.buildTrackerURL(peerID)
	if err != nil {
		return nil, err
	}
	fmt.Println("tracker URL: ", trackerUrl)

	response := bencodedTrackerResponse{}

	if strings.HasPrefix(trackerUrl, "http://") {
		client := http.Client{
			Timeout: 15 * time.Second,
		}
		benRes, err := client.Get(trackerUrl)
		if err != nil {
			return nil, err
		}
		// response is a bencoded value
		err = bencode.Unmarshal(benRes.Body, &response)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(trackerUrl, "udp://") {
		parsedURL, _ := url.Parse(trackerUrl)

		// resolve the UDP server address
		udpAddr, err := net.ResolveUDPAddr("udp", parsedURL.Hostname()+ ":" +parsedURL.Port())
		if err != nil {
			return nil, err
		}
		// create a UDP connection
		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			log.Fatal("Error in establishing udp connection")
			return nil, err
		}
		defer conn.Close()
		message := []byte("Hello, UDP server!")

		// Send the message to the server
		_, err = conn.Write(message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return nil, err
		}

		fmt.Println("Message sent to", parsedURL.Hostname())
		fmt.Println("trackers expecting UDP protocol is still under development")

		// Receive response from the server
		// buffer := make([]byte, 1024) // Buffer for incoming data
		// n, _, err := conn.ReadFromUDP(buffer)
		// if err != nil {
		// 	fmt.Println("Error receiving response:", err)
		// 	return nil, err
		// }

		// response := buffer[:n]
		// fmt.Println("Response from server:", string(response))

	}

	p, err := peers.UnmarshalPeers([]byte(response.Peers))
	if err != nil {
		return nil, err
	}
	return p, nil
}
