package server

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	UserAgentKey        = "user-agent"
	GatewayUserAgentKey = "grpcgateway-user-agent"
	XForwardedForKey    = "x-forwarded-for"
)

type todoServiceServerMetadata struct {
	userAgent string
	clientIp  string
}

func extractMetadata(ctx context.Context) todoServiceServerMetadata {
	var userAgent, clientIp string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(UserAgentKey); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}

		peer, ok := peer.FromContext(ctx)
		if ok {
			clientIp = peer.Addr.String()
		}

		if gwUserAgents := md.Get(GatewayUserAgentKey); len(gwUserAgents) > 0 {
			userAgent = gwUserAgents[0]
		}

		if clientIps := md.Get(XForwardedForKey); len(clientIps) > 0 {
			clientIp = clientIps[0]
		}
	}

	return todoServiceServerMetadata{userAgent: userAgent, clientIp: clientIp}
}
