package server

import (
	"context"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func TodoUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Info("Request received", "method", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		slog.Error("Request failed", "method", info.FullMethod, "error", err)
		return resp, err
	}

	slog.Info("Request succeeded", "method", info.FullMethod, "message", resp)
	return resp, nil
}
