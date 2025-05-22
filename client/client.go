package client

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/AniruddhGuttikar/BitTorrent-client-go/bitfield"
	"github.com/AniruddhGuttikar/BitTorrent-client-go/handshake"
	"github.com/AniruddhGuttikar/BitTorrent-client-go/message"
	"github.com/AniruddhGuttikar/BitTorrent-client-go/peers"
)

// Client is a TCP connection with a Peer
type Client struct {
	PeerId   [20]byte
	Peer     peers.Peer
	Conn     net.Conn
	InfoHash [20]byte
	BitField bitfield.BitField
	Choked   bool
}

func completeHandshake(conn net.Conn, infohash, peerID [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	handshakeRequest := handshake.New(infohash, peerID)

	_, err := conn.Write(handshakeRequest.Serialize())
	if err != nil {
		return nil, err
	}

	handshakeResponse, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(handshakeResponse.InfoHash[:], infohash[:]) {
		return nil, fmt.Errorf("expected infohash %x but fot %x", handshakeResponse.InfoHash, infohash)
	}
	return handshakeResponse, nil
}

func receiveBitfield(conn net.Conn) (bitfield.BitField, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf("expected bitfield but got %s", msg)
	}
	if msg.ID != message.MsgBitfield {
		return nil, fmt.Errorf("expected bitfield (ID %d), got ID %d", message.MsgBitfield, msg.ID)
	}
	return msg.Payload, nil
}

// New connects with a peer and completes the handshake and receives handshake response
func New(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := receiveBitfield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		BitField: bf,
		Peer:     peer,
		PeerId:   peerID,
		InfoHash: infoHash,
	}, nil
}

func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	return msg, err
}

// SendRequest sends a Request message to the peer
func (c *Client) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// SendInterested sends an Interested message to the peer
func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer
func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have message to the peer
func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
