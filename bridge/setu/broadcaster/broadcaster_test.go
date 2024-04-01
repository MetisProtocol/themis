package broadcaster

import (
	"os"
	"testing"

	"github.com/metis-seq/themis/helper"
	"github.com/spf13/viper"
)

// Parallel test - to check BroadcastToThemis synchronisation
func TestBroadcastToThemis(t *testing.T) {
	t.Parallel()

	tendermintNode := "tcp://localhost:26657"
	viper.Set(helper.TendermintNodeFlag, tendermintNode)
	viper.Set("log_level", "info")
	helper.InitThemisConfig(os.ExpandEnv("$HOME/.themisd"))
}
