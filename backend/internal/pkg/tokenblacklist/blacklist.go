package tokenblacklist

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

// Blacklist stores revoked tokens in memory with automatic expiry cleanup.
type Blacklist struct {
	entries map[string]time.Time // token hash -> expiry
	mu      sync.RWMutex
}

func New() *Blacklist {
	bl := &Blacklist{
		entries: make(map[string]time.Time),
	}
	go bl.cleanup()
	return bl
}

func hash(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// Revoke adds a token to the blacklist until its expiry time.
func (bl *Blacklist) Revoke(token string, expiry time.Time) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	bl.entries[hash(token)] = expiry
}

// IsRevoked returns true if the token has been blacklisted.
func (bl *Blacklist) IsRevoked(token string) bool {
	bl.mu.RLock()
	defer bl.mu.RUnlock()
	exp, exists := bl.entries[hash(token)]
	if !exists {
		return false
	}
	return time.Now().Before(exp)
}

func (bl *Blacklist) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		bl.mu.Lock()
		now := time.Now()
		for k, exp := range bl.entries {
			if now.After(exp) {
				delete(bl.entries, k)
			}
		}
		bl.mu.Unlock()
	}
}
