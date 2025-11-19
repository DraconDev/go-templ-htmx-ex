package middleware

import (
	"sync"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
)

// SessionCache stores validation results with 15-second TTL
type SessionCache struct {
	sync.RWMutex
	entries map[string]*cacheEntry
}

type cacheEntry struct {
	userInfo  layouts.UserInfo
	expiresAt time.Time
}

// NewSessionCache creates a new session cache
func NewSessionCache() *SessionCache {
	return &SessionCache{
		entries: make(map[string]*cacheEntry),
	}
}

// Get retrieves cached user info if not expired
func (c *SessionCache) Get(sessionID string) (layouts.UserInfo, bool) {
	c.RLock()
	defer c.RUnlock()

	entry, exists := c.entries[sessionID]
	if !exists {
		return layouts.UserInfo{LoggedIn: false}, false
	}

	if time.Now().After(entry.expiresAt) {
		// Expired entry, clean up
		delete(c.entries, sessionID)
		return layouts.UserInfo{LoggedIn: false}, false
	}

	return entry.userInfo, true
}

// Set caches user info with 15-second TTL
func (c *SessionCache) Set(sessionID string, userInfo layouts.UserInfo) {
	c.Lock()
	defer c.Unlock()

	c.entries[sessionID] = &cacheEntry{
		userInfo:  userInfo,
		expiresAt: time.Now().Add(15 * time.Second),
	}
}
