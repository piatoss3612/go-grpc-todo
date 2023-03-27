package server

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func TodoServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Info("Request received", "method", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		slog.Error("Request failed", "method", info.FullMethod, "error", err)
		return resp, err
	}

	slog.Info("Send a message", "type", fmt.Sprintf("%T", resp))
	return resp, nil
}

func TodoServerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	slog.Info("Request received", "method", info.FullMethod)

	wss := newWrappedServerStream(ss)

	err := handler(srv, wss)
	if err != nil {
		slog.Error("Request failed", "method", info.FullMethod, "error", err)
		return err
	}

	return nil
}

type wrappedServerStream struct {
	grpc.ServerStream
}

func newWrappedServerStream(ss grpc.ServerStream) grpc.ServerStream {
	return &wrappedServerStream{ss}
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	slog.Info("Send a message", "type", fmt.Sprintf("%T", m))
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	slog.Info("Receive a message", "type", fmt.Sprintf("%T", m))
	return w.ServerStream.RecvMsg(m)
}
