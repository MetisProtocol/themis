package helper

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPowerFromAmount(t *testing.T) {
	t.Parallel()

	scenarios1 := map[string]string{
		"48000000000000000000000": "48000",
		"10000000000000000000000": "10000",
		"1000000000000000000000":  "1000",
		"4800000000000000000000":  "4800",
		"480000000000000000000":   "480",
		"20000000000000000000":    "20",
		"10000000000000000000":    "10",
		"1000000000000000000":     "1",
	}

	for k, v := range scenarios1 {
		bv, _ := big.NewInt(0).SetString(k, 10)
		p, err := GetPowerFromAmount(bv)
		require.Nil(t, err, "Error must be null for input %v, output %v", k, v)
		require.Equal(t, p.String(), v, "Power must match")
	}
}
