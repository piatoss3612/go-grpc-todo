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

	slog.Info("Send a message", "type", fmt.Sprintf("%T", resp))
	return nil
}

func TodoClientStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	slog.Info("Request sent", "method", method)

	cs, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		slog.Error("Request failed", "method", method, "error", err)
		return nil, err
	}

	return newWrappedClientStream(cs), nil
}

type wrappedClientStream struct {
	grpc.ClientStream
}

func newWrappedClientStream(cs grpc.ClientStream) grpc.ClientStream {
	return &wrappedClientStream{cs}
}

func (w *wrappedClientStream) SendMsg(m interface{}) error {
	slog.Info("Send a message", "type", fmt.Sprintf("%T", m))
	return w.ClientStream.SendMsg(m)
}

func (w *wrappedClientStream) RecvMsg(m interface{}) error {
	slog.Info("Receive a message", "type", fmt.Sprintf("%T", m))
	return w.ClientStream.RecvMsg(m)
}
