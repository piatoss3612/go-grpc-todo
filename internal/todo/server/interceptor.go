package server

import (
	"context"
	"time"

	"github.com/piatoss3612/go-grpc-todo/internal/event"
	te "github.com/piatoss3612/go-grpc-todo/internal/todo/event"
	"github.com/piatoss3612/go-grpc-todo/proto/gen/go/todo/v1"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type TodoServiceServerInterceptor interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
	todo.TodoServiceServer
}

type interceptor struct {
	srv todo.TodoServiceServer
	pub event.Publisher
}

func NewInterceptor(srv todo.TodoServiceServer, pub event.Publisher) TodoServiceServerInterceptor {
	return &interceptor{srv: srv, pub: pub}
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

		md := extractMetadata(ss.Context())

		slog.Info("Request received", "method", info.FullMethod, "user-agent", md.userAgent, "client-ip", md.clientIp)

		return handler(srv, ss)
	}
}

func (i *interceptor) Add(ctx context.Context, req *todo.AddRequest) (resp *todo.AddResponse, err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(ctx, err)
			return
		}
		_ = i.publishEvent(ctx, te.EventTopicTodoCreated, resp.String())
	}()
	return i.srv.Add(ctx, req)
}

func (i *interceptor) AddMany(stream todo.TodoService_AddManyServer) (err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(stream.Context(), err)
			return
		}
		_ = i.publishEvent(stream.Context(), te.EventTopicTodoCreated, "Added many todos")
	}()
	return i.srv.AddMany(stream)
}

func (i *interceptor) Get(ctx context.Context, req *todo.GetRequest) (resp *todo.Todo, err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(ctx, err)
		}
	}()
	return i.srv.Get(ctx, req)
}

func (i *interceptor) GetAll(req *todo.Empty, stream todo.TodoService_GetAllServer) (err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(stream.Context(), err)
		}
	}()
	return i.srv.GetAll(req, stream)
}

func (i *interceptor) Update(ctx context.Context, req *todo.UpdateRequest) (resp *todo.UpdateResponse, err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(ctx, err)
			return
		}
		_ = i.publishEvent(ctx, te.EventTopicTodoUpdated, resp.String())
	}()
	return i.srv.Update(ctx, req)
}

func (i *interceptor) UpdateMany(stream todo.TodoService_UpdateManyServer) (err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(stream.Context(), err)
			return
		}
		_ = i.publishEvent(stream.Context(), te.EventTopicTodoUpdated, "Updated many todos")
	}()
	return i.srv.UpdateMany(stream)
}

func (i *interceptor) Delete(ctx context.Context, req *todo.DeleteRequest) (resp *todo.DeleteResponse, err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(ctx, err)
			return
		}
		_ = i.publishEvent(ctx, te.EventTopicTodoDeleted, resp.String())
	}()
	return i.srv.Delete(ctx, req)
}

func (i *interceptor) DeleteAll(ctx context.Context, req *todo.Empty) (resp *todo.DeleteAllResponse, err error) {
	defer func() {
		if err != nil {
			_ = i.publishError(ctx, err)
			return
		}
		_ = i.publishEvent(ctx, te.EventTopicTodoDeleted, "Deleted all todos")
	}()
	return i.srv.DeleteAll(ctx, req)
}

func (i *interceptor) publishEvent(ctx context.Context, topic te.EventTopic, data interface{}) error {
	evt, _ := te.NewTodoEvent(topic, data)
	return i.pub.Publish(ctx, evt)
}

func (i *interceptor) publishError(ctx context.Context, err error) error {
	var msg string

	s, ok := status.FromError(err)
	if ok {
		msg = s.Message()
	} else {
		msg = err.Error()
	}

	return i.publishEvent(ctx, te.EventTopicTodoError, msg)
}
