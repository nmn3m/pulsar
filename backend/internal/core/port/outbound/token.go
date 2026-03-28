package outbound

import "time"

type TokenRevoker interface {
	Revoke(token string, expiry time.Time)
	IsRevoked(token string) bool
}
