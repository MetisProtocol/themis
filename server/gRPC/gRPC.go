package gRPC

import (
	"context"
	"net"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	proto "github.com/metis-seq/themis/proto/themis"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	spanURL         = "/metis/span/%v"
	eventRecordList = "/clerk/event-record/list"
)

var logger tmLog.Logger

type ThemisGRPCServer struct {
	proto.UnimplementedThemisServer
	cdc *codec.Codec
}

func SetupGRPCServer(shutDownCtx context.Context, cdc *codec.Codec, addr string, lggr tmLog.Logger) error {
	logger = lggr
	grpcServer := grpc.NewServer(withLoggingUnaryInterceptor())
	proto.RegisterThemisServer(grpcServer,
		&ThemisGRPCServer{
			cdc: cdc,
		})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed to serve grpc server", "err", err)
		}

		<-shutDownCtx.Done()
		grpcServer.Stop()
		lis.Close()
		logger.Info("GRPC Server stopped", "addr", addr)
	}()

	logger.Info("GRPC Server started", "addr", addr)

	return nil
}

func withLoggingUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(loggingServerInterceptor)
}

func loggingServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)
	if err != nil {
		err = status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Request", "method", info.FullMethod, "duration", time.Since(start), "error", err)

	return h, err
}
