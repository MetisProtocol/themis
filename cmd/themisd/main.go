package main

import (
	"context"

	"github.com/metis-seq/themis/cmd/themisd/service"
)

func main() {
	service.NewThemisService(context.Background(), nil)
}
