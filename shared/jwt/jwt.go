package jwt

import (
	"context"
	"errors"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/satori/go.uuid"
)

type contextKey string

const (
	// KeyJWTToken holds the key to a JWT Token in context.
	KeyJWTToken contextKey = "JWTToken"
	// KeyJWTClaims holds the key to a JWT Claims object in context.
	KeyJWTClaims contextKey = "JWTClaims"
	// KeyJWTError holds the error found while decoding and/or verifying the token
	KeyJWTError contextKey = "JWTError"
)

var (
	// ErrUnauthorized is returned when no valid token is found
	ErrUnauthorized = errors.New("unauthorized")

	// ErrIllegalSigningMethod is returned when "none" method is requested
	ErrIllegalSigningMethod = errors.New("JWT Token illegal signing method")

	// ErrUnexpectedSigningMethod is returned when a different signing method is
	// requested then the one we expect.
	ErrUnexpectedSigningMethod = errors.New("JWT Token unexpected signing method")

	// ErrTokenInvalid is returned when the parsed token is invalid.
	ErrTokenInvalid = errors.New("JWT Token is invalid")

	// ErrTokenExpired is returned when the parsed token has expired.
	ErrTokenExpired = errors.New("JWT Token is expired")
)

// Handler holds the JWT Handling interface
type Handler interface {
	NewParser(mustValidate bool) endpoint.Middleware
	CreateToken(claims jwt.Claims) (token string, err error)
	ParseAuthMessage(msg AuthMessage) (token *jwt.Token, err error)
}

// AuthMessage holding JWT (websockets)
type AuthMessage struct {
	Schema        string `json:"schema"`
	Authorization string `json:"authorization"`
}

type handler struct {
	kid     string
	key     []byte
	keyFunc jwt.Keyfunc
	method  jwt.SigningMethod
	claims  jwt.Claims
}

// NewHandler returns a new JWT Handler using given parameters
func NewHandler(
	multiTenancy bool, kid string, key []byte, keyFunc jwt.Keyfunc,
	method jwt.SigningMethod, claims jwt.Claims, instanceID uuid.UUID,
) Handler {
	h := &handler{
		kid:     kid,
		key:     key,
		keyFunc: keyFunc,
		method:  method,
		claims:  claims,
	}
	if h.keyFunc == nil {
		h.keyFunc = h.defaultKeyFunc
	}
	return h
}

// NewParser returns a new JWT token parsing middleware which can be set to
// either immediately error on JWT issue or inject error into context enabling
// to defer the handling of the error.
func (h handler) NewParser(mustValidate bool) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			tokenString, ok := ctx.Value(KeyJWTToken).(string)
			if !ok {
				if mustValidate {
					return nil, ErrUnauthorized
				}
				ctx = context.WithValue(ctx, KeyJWTError, ErrUnauthorized)
				return next(ctx, request)
			}

			token, err := jwt.ParseWithClaims(tokenString, h.claims, func(token *jwt.Token) (interface{}, error) {
				if strings.EqualFold(token.Method.Alg(), "none") {
					return nil, ErrIllegalSigningMethod
				}
				if token.Method != h.method {
					return nil, ErrUnexpectedSigningMethod
				}
				return h.keyFunc(token)
			})

			if err != nil {
				if mustValidate {
					return nil, err
				}
				ctx = context.WithValue(ctx, KeyJWTError, err)
				return next(ctx, request)
			}

			if !token.Valid {
				if mustValidate {
					return nil, ErrTokenInvalid
				}
				ctx = context.WithValue(ctx, KeyJWTError, ErrTokenInvalid)
				return next(ctx, request)
			}

			// test claims
			if err = token.Claims.Valid(); err != nil {
				if mustValidate {
					return nil, ErrTokenInvalid
				}
				ctx = context.WithValue(ctx, KeyJWTError, ErrTokenInvalid)
				return next(ctx, request)
			}

			claim, ok := token.Claims.(*Claims)
			if !ok {
				if mustValidate {
					return nil, ErrTokenInvalid
				}
				ctx = context.WithValue(ctx, KeyJWTError, ErrTokenInvalid)
				return next(ctx, request)
			}
			if !strings.EqualFold(claim.Audience, h.instanceID.String()) {
				if mustValidate {
					return nil, ErrUnauthorized
				}
				ctx = context.WithValue(ctx, KeyJWTError, ErrUnauthorized)
				return next(ctx, request)
			}

			ctx = context.WithValue(ctx, KeyJWTClaims, token.Claims)
			return next(ctx, request)
		}
	}
}

func (h handler) ParseAuthMessage(msg AuthMessage) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(msg.Authorization, h.claims, func(token *jwt.Token) (interface{}, error) {
		if strings.EqualFold(token.Method.Alg(), "none") {
			return nil, ErrIllegalSigningMethod
		}
		if token.Method != h.method {
			return nil, ErrUnexpectedSigningMethod
		}
		return h.keyFunc(token)
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}
	return token, nil
}

func (h handler) CreateToken(claims jwt.Claims) (token string, err error) {
	j := jwt.NewWithClaims(h.method, claims)
	j.Header["kid"] = h.kid

	return j.SignedString(h.key)
}

func (h handler) defaultKeyFunc(token *jwt.Token) (interface{}, error) {
	switch token.Header["kid"] {
	case h.kid:
		return h.key, nil
	default:
		return nil, errors.New("invalid kid received")
	}
}

// GetErrorFromContext retrieves a JWT Parser error of one was injected.
func GetErrorFromContext(ctx context.Context) error {
	if err, ok := ctx.Value(KeyJWTError).(error); ok {
		return err
	}
	return nil
}
