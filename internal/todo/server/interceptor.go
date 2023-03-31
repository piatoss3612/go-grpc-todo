package server

import (
	"context"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

const (
	UserAgentKey        = "user-agent"
	GatewayUserAgentKey = "grpcgateway-user-agent"
	XForwardedForKey    = "x-forwarded-for"
	ReceivedMsgCountKey = "received-msg-count"
	SentMsgCountKey     = "sent-msg-count"
)

type TodoServerInterceptor interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
}

type todoServerInterceptor struct{}

func NewTodoServerInterceptor() TodoServerInterceptor {
	return &todoServerInterceptor{}
}

func (i *todoServerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md := extractMetadata(ctx)

		slog.Info("Request received", "method", info.FullMethod, "user-agent", md.userAgent, "client-ip", md.clientIp)

		resp, err = handler(ctx, req)
		if err != nil {
			slog.Error("Request failed", "method", info.FullMethod, "error", err)
		} else {
			slog.Info("Request handled successfully", "method", info.FullMethod)
		}

		return
	}
}

func (i *todoServerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md := extractMetadata(ss.Context())

		slog.Info("Request received", "method", info.FullMethod, "user-agent", md.userAgent, "client-ip", md.clientIp)

		wss := newWrappedServerStream(ss)

		err := handler(srv, wss)
		if err != nil {
			slog.Error("Request failed", "method", info.FullMethod, "error", err)
			return err
		}

		ctx := wss.Context()
		rmc, smc := ctx.Value(ReceivedMsgCountKey).(int), ctx.Value(SentMsgCountKey).(int)

		slog.Info("Request handled successfully", "method", info.FullMethod, ReceivedMsgCountKey, rmc, SentMsgCountKey, smc)

		return nil
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	recvMsgCount    int
	sentMsgCountKey int
}

func newWrappedServerStream(ss grpc.ServerStream) grpc.ServerStream {
	return &wrappedServerStream{
		ServerStream:    ss,
		recvMsgCount:    0,
		sentMsgCountKey: 0,
	}
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	w.sentMsgCountKey++
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	w.recvMsgCount++
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedServerStream) Context() context.Context {
	ctx := context.WithValue(w.ServerStream.Context(), ReceivedMsgCountKey, w.recvMsgCount)
	ctx = context.WithValue(ctx, SentMsgCountKey, w.sentMsgCountKey)
	return ctx
}
