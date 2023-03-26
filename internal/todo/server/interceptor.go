package server

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func TodoServerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Info("Request received", "method", info.FullMethod, "time", time.Now().Format(time.RFC3339))

	resp, err := handler(ctx, req)

	if err != nil {
		slog.Error("Request failed", "method", info.FullMethod, "error", err, "time", time.Now().Format(time.RFC3339))
		return resp, err
	}

	slog.Info("Request succeeded", "method", info.FullMethod, "type", fmt.Sprintf("%T", resp), "time", time.Now().Format(time.RFC3339))
	return resp, nil
}

func TodoServerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	slog.Info("Request received", "method", info.FullMethod, "time", time.Now().Format(time.RFC3339))

	wss := newWrappedServerStream(ss)

	err := handler(srv, wss)

	if err != nil {
		slog.Error("Request failed", "method", info.FullMethod, "error", err, "time", time.Now().Format(time.RFC3339))
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
	slog.Info("Send a message", "type", fmt.Sprintf("%T", m), "time", time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	slog.Info("Receive a message", "type", fmt.Sprintf("%T", m), "time", time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}
