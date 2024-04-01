package utils

import (
	"encoding/binary"

	"github.com/metis-seq/themis/proto/themis"
)

func ConvertH160toAddress(h160 *themis.H160) [20]byte {
	var addr [20]byte

	binary.BigEndian.PutUint64(addr[0:], h160.Hi.Hi)
	binary.BigEndian.PutUint64(addr[8:], h160.Hi.Lo)
	binary.BigEndian.PutUint32(addr[16:], h160.Lo)

	return addr
}

func ConvertAddressToH160(addr [20]byte) *themis.H160 {
	return &themis.H160{
		Lo: binary.BigEndian.Uint32(addr[16:]),
		Hi: &themis.H128{Lo: binary.BigEndian.Uint64(addr[8:]), Hi: binary.BigEndian.Uint64(addr[0:])},
	}
}

func ConvertH256ToHash(h256 *themis.H256) [32]byte {
	var hash [32]byte

	binary.BigEndian.PutUint64(hash[0:], h256.Hi.Hi)
	binary.BigEndian.PutUint64(hash[8:], h256.Hi.Lo)
	binary.BigEndian.PutUint64(hash[16:], h256.Lo.Hi)
	binary.BigEndian.PutUint64(hash[24:], h256.Lo.Lo)

	return hash
}

func ConvertHashToH256(hash [32]byte) *themis.H256 {
	return &themis.H256{
		Lo: &themis.H128{Lo: binary.BigEndian.Uint64(hash[24:]), Hi: binary.BigEndian.Uint64(hash[16:])},
		Hi: &themis.H128{Lo: binary.BigEndian.Uint64(hash[8:]), Hi: binary.BigEndian.Uint64(hash[0:])},
	}
}
