package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/metis-seq/themis/sidechannel/types"
)

func TestCodec(t *testing.T) {
	t.Parallel()
	require.NotNil(t, types.ModuleCdc, "ModuleCdc shouldn't be nil")
}
