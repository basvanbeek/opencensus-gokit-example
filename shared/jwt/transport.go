package jwt

import (
	"context"
	"fmt"
	stdhttp "net/http"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/kit/transport/http"
)

const (
	bearer       string = "bearer"
	bearerFormat string = "Bearer %s"
)

// ToHTTPContext moves JWT token from HTTP Request Header to context.
func ToHTTPContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		token, ok := extractTokenFromAuthHeader(r.Header.Get("Authorization"))
		if !ok {
			return ctx
		}

		return context.WithValue(ctx, KeyJWTToken, token)
	}
}

// FromHTTPContext moves JWT token from context to HTTP Request Header.
func FromHTTPContext() http.RequestFunc {
	return func(ctx context.Context, r *stdhttp.Request) context.Context {
		token, ok := ctx.Value(KeyJWTToken).(string)
		if ok {
			r.Header.Add("Authorization", generateAuthHeaderFromToken(token))
		}
		return ctx
	}
}

// ToGRPCContext moves JWT token from grpc metadata to context.
func ToGRPCContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		// capital "Key" is illegal in HTTP/2.
		authHeader, ok := md["authorization"]
		if !ok {
			return ctx
		}

		token, ok := extractTokenFromAuthHeader(authHeader[0])
		if ok {
			ctx = context.WithValue(ctx, KeyJWTToken, token)
		}

		return ctx
	}
}

// FromGRPCContext moves JWT token from context to grpc metadata.
func FromGRPCContext() grpc.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		token, ok := ctx.Value(KeyJWTToken).(string)
		if ok {
			(*md)["authorization"] = []string{generateAuthHeaderFromToken(token)}
		}

		return ctx
	}
}

func extractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != bearer {
		return "", false
	}

	return authHeaderParts[1], true
}

func generateAuthHeaderFromToken(token string) string {
	return fmt.Sprintf(bearerFormat, token)
}
