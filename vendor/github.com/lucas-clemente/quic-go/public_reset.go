package quic

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/lucas-clemente/quic-go/handshake"
	"github.com/lucas-clemente/quic-go/protocol"
	"github.com/lucas-clemente/quic-go/utils"
)

type publicReset struct {
	rejectedPacketNumber protocol.PacketNumber
	nonce                uint64
}

func writePublicReset(connectionID protocol.ConnectionID, rejectedPacketNumber protocol.PacketNumber, nonceProof uint64) []byte {
	b := &bytes.Buffer{}
	b.WriteByte(0x0a)
	utils.WriteUint64(b, uint64(connectionID))
	utils.WriteUint32(b, uint32(handshake.TagPRST))
	utils.WriteUint32(b, 2)
	utils.WriteUint32(b, uint32(handshake.TagRNON))
	utils.WriteUint32(b, 8)
	utils.WriteUint32(b, uint32(handshake.TagRSEQ))
	utils.WriteUint32(b, 16)
	utils.WriteUint64(b, nonceProof)
	utils.WriteUint64(b, uint64(rejectedPacketNumber))
	return b.Bytes()
}

func parsePublicReset(r *bytes.Reader) (*publicReset, error) {
	pr := publicReset{}
	tag, tagMap, err := handshake.ParseHandshakeMessage(r)
	if err != nil {
		return nil, err
	}
	if tag != handshake.TagPRST {
		return nil, errors.New("wrong public reset tag")
	}

	rseq, ok := tagMap[handshake.TagRSEQ]
	if !ok {
		return nil, errors.New("RSEQ missing")
	}
	if len(rseq) != 8 {
		return nil, errors.New("invalid RSEQ tag")
	}
	pr.rejectedPacketNumber = protocol.PacketNumber(binary.LittleEndian.Uint64(rseq))

	rnon, ok := tagMap[handshake.TagRNON]
	if !ok {
		return nil, errors.New("RNON missing")
	}
	if len(rnon) != 8 {
		return nil, errors.New("invalid RNON tag")
	}
	pr.nonce = binary.LittleEndian.Uint64(rnon)

	return &pr, nil
}
