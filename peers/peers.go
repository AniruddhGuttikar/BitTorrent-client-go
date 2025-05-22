package peers

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func UnmarshalPeers(peersBin []byte) ([]Peer, error) {
	//4 for IP + 2 for PORT
	const peerLength = 6
	if len(peersBin)%peerLength != 0 {
		err := fmt.Errorf("received malformed peers")
		return nil, err
	}

	numPeers := len(peersBin) / peerLength
	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		index := peerLength * i
		peers[i].IP = net.IP(peersBin[index : index+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[index+4 : index+6])
	}
	return peers, nil
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}
