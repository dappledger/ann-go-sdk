package sdk

import (
	"fmt"
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/common"
)

const MAX_BATCH_PARAMS = 10000

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func ParseData(data ...interface{}) (ret []byte) {
	for _, item := range data {
		switch t := item.(type) {
		case string:
			var str []byte
			if IsHex(t) {
				str = common.Hex2Bytes(t[2:])
			} else {
				str = []byte(t)
			}

			ret = append(ret, common.RightPadBytes(str, 32)...)
		case []byte:
			ret = append(ret, common.LeftPadBytes(t, 32)...)
		}
	}

	return
}

func IsHex(str string) bool {
	l := len(str)
	return l >= 4 && l%2 == 0 && str[0:2] == "0x"
}

func SanitizeHex(hex string) string {
	return strings.TrimPrefix(strings.ToLower(hex), "0x")
}

func bytesN2Slice(value interface{}, m int) ([]byte, error) {
	switch m {
	case 0:
		v := value.([]byte)
		return v, nil
	case 8:
		v := value.([8]byte)
		return v[:], nil
	case 16:
		v := value.([16]byte)
		return v[:], nil
	case 32:
		v := value.([32]byte)
		return v[:], nil
	case 64:
		v := value.([64]byte)
		return v[:], nil
	}
	return nil, fmt.Errorf("type(bytes%d) not support", m)
}
