package client

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func TodoClientUnaryInterceptor(ctx context.Context, method string, req, resp interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	slog.Info("Request sent", "method", method)

	err := invoker(ctx, method, req, resp, cc, opts...)
	if err != nil {
		slog.Error("Request failed", "method", method, "error", err)
		return err
	}

	slog.Info("Request succeeded", "method", method, "message type", fmt.Sprintf("%T", resp))
	return nil
}
