package server

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type todoServerMetadata struct {
	userAgent string
	clientIp  string
}

func extractMetadata(ctx context.Context) todoServerMetadata {
	var userAgent, clientIp string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get("user-agent"); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}

		peer, ok := peer.FromContext(ctx)
		if ok {
			clientIp = peer.Addr.String()
		}

		if gwUserAgents := md.Get("grpcgateway-user-agent"); len(gwUserAgents) > 0 {
			userAgent = gwUserAgents[0]
		}

		if clientIps := md.Get("x-forwarded-for"); len(clientIps) > 0 {
			clientIp = clientIps[0]
		}
	}

	return todoServerMetadata{userAgent: userAgent, clientIp: clientIp}
}
