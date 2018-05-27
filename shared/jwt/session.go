package jwt

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
)

// SessionData holds tenant data that will be serialized into the JWT.
type SessionData struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	SessionID uuid.UUID `json:"session_id"`
	Scopes    []string  `json:"scopes"`
	Name      string    `json:"name"`
}

// TokenGenerator is a function which takes a SessionData struct and tries to
// create a JWT token for it.
type TokenGenerator func(session SessionData) (string, error)

// Claims holds the typical claims object for this application
type Claims struct {
	SessionData
	jwt.StandardClaims
}

// NewDefaultClaims returns a newly constructed Claims object.
func NewDefaultClaims(session SessionData, issuerID uuid.UUID) *Claims {
	return &Claims{
		SessionData: session,
		StandardClaims: jwt.StandardClaims{
			Issuer:   issuerID.String(),
			IssuedAt: time.Now().Unix(),
			Subject:  session.TenantID.String(),
			Audience: instanceID.String(),
		},
	}
}

// NewTokenGenerator returns a new TokenGenerator func using provided jwt.Handler
func NewTokenGenerator(h Handler, instanceID uuid.UUID) TokenGenerator {
	return func(session SessionData) (string, error) {
		claims := NewDefaultClaims(session, instanceID)
		return h.CreateToken(claims)
	}
}

// GetClaimsFromContext retrieves a JWT Claims object from context if available.
func GetClaimsFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(KeyJWTClaims).(*Claims)
	if !ok {
		return nil
	}
	return claims
}

//GetTenantIDFromContext retrieves the TenantID from context. Defaults to Nil
// UUID if not found.
func GetTenantIDFromContext(ctx context.Context) uuid.UUID {
	claims, ok := ctx.Value(KeyJWTClaims).(*Claims)
	if !ok {
		return uuid.Nil
	}
	return claims.TenantID
}
