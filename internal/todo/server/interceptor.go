package server

import (
	"context"
	"time"

	"github.com/piatoss3612/go-grpc-todo/proto/gen/go/todo/v1"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	UserAgentKey        = "user-agent"
	GatewayUserAgentKey = "grpcgateway-user-agent"
	XForwardedForKey    = "x-forwarded-for"
	ReceivedMsgCountKey = "received-msg-count"
	SentMsgCountKey     = "sent-msg-count"
)

type TodoServiceServerInterceptor interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
	todo.TodoServiceServer
}

type interceptor struct {
	srv todo.TodoServiceServer
	// TODO: add event bus
}

func NewInterceptor(srv todo.TodoServiceServer) TodoServiceServerInterceptor {
	return &interceptor{srv: srv}
}

func (i *interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		defer func() {
			if err != nil {
				s, ok := status.FromError(err)
				if ok {
					slog.Error("Request failed", "method", info.FullMethod, "error", s.Message(),
						"code", s.Code(), "duration", time.Since(start).String())
				} else {
					slog.Error("Request failed", "method", info.FullMethod, "error", err,
						"duration", time.Since(start).String())
				}
				return
			}
			slog.Info("Request handled successfully", "method", info.FullMethod, "duration", time.Since(start).String())
		}()

		md := extractMetadata(ctx)

		slog.Info("Request received", "method", info.FullMethod, "user-agent", md.userAgent, "client-ip", md.clientIp)

		return handler(ctx, req)
	}
}

func (i *interceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		start := time.Now()
		wss := newWrappedServerStream(ss)

		defer func() {
			if err != nil {
				s, ok := status.FromError(err)
				if ok {
					slog.Error("Request failed", "method", info.FullMethod, "error", s.Message(),
						"code", s.Code(), "duration", time.Since(start).String())
				} else {
					slog.Error("Request failed", "method", info.FullMethod, "error", err,
						"duration", time.Since(start).String())
				}
				return
			}

			ctx := wss.Context()
			rmc, smc := ctx.Value(ReceivedMsgCountKey).(int), ctx.Value(SentMsgCountKey).(int)

			slog.Info("Request handled successfully", "method", info.FullMethod,
				ReceivedMsgCountKey, rmc, SentMsgCountKey, smc, "duration", time.Since(start).String())
		}()

		md := extractMetadata(ss.Context())

		slog.Info("Request received", "method", info.FullMethod, "user-agent", md.userAgent, "client-ip", md.clientIp)

		return handler(srv, wss)
	}
}

func (i *interceptor) Add(ctx context.Context, req *todo.AddRequest) (resp *todo.AddResponse, err error) {
	defer func() {
		// TODO: send event to event bus
	}()
	return i.srv.Add(ctx, req)
}

func (i *interceptor) AddMany(stream todo.TodoService_AddManyServer) (err error) {
	defer func() {
		// TODO: send event to event bus
	}()
	return i.srv.AddMany(stream)
}

func (i *interceptor) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.Todo, err error) {
	return i.srv.Get(ctx, req)
}

func (i *interceptor) GetAll(req *todo.Empty, stream todo.TodoService_GetAllServer) (err error) {
	return i.srv.GetAll(req, stream)
}

func (i *interceptor) Update(ctx context.Context, req *todo.UpdateRequest) (resp *todo.UpdateResponse, err error) {
	defer func() {
		// TODO: send event to event bus
	}()
	return i.srv.Update(ctx, req)
}

func (i *interceptor) UpdateMany(stream todo.TodoService_UpdateManyServer) (err error) {
	defer func() {
		// TODO: send event to event bus
	}()
	return i.srv.UpdateMany(stream)
}

func (i *interceptor) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	defer func() {
		// TODO: send event to event bus
	}()
	return i.srv.Delete(ctx, req)
}

func (i *interceptor) DeleteAll(ctx context.Context, req *todo.Empty) (resp *todo.DeleteAllResponse, err error) {
	defer func() {
		// TODO: send event to event bus
	}()
	return i.srv.DeleteAll(ctx, req)
}

type wrappedServerStream struct {
	grpc.ServerStream
	recvMsgCount int
	sentMsgCount int
}

func newWrappedServerStream(ss grpc.ServerStream) grpc.ServerStream {
	return &wrappedServerStream{
		ServerStream: ss,
		recvMsgCount: 0,
		sentMsgCount: 0,
	}
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	w.sentMsgCount++
	return w.ServerStream.SendMsg(m)
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	w.recvMsgCount++
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedServerStream) Context() context.Context {
	ctx := context.WithValue(w.ServerStream.Context(), ReceivedMsgCountKey, w.recvMsgCount)
	ctx = context.WithValue(ctx, SentMsgCountKey, w.sentMsgCount)
	return ctx
}
