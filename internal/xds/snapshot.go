package xds

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"example.com/api-traffic-governance/internal/domain"
	"fmt"
	"sync"
	"time"
)

type Manager struct {
	mu      sync.RWMutex
	current domain.Snapshot
	secret  []byte
	clients map[string]string
}

func New(secret string) *Manager {
	return &Manager{secret: []byte(secret), clients: map[string]string{}}
}
func (m *Manager) Compile(routes []domain.Route, quotas []domain.Quota) domain.Snapshot {
	b, _ := json.Marshal(struct {
		R []domain.Route
		Q []domain.Quota
	}{routes, quotas})
	h := sha256.Sum256(b)
	v := hex.EncodeToString(h[:])
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(v))
	n := hex.EncodeToString(mac.Sum(nil))
	routesCopy := append([]domain.Route(nil), routes...)
	quotasCopy := append([]domain.Quota(nil), quotas...)
	s := domain.Snapshot{Version: v, Nonce: n, Routes: routesCopy, Quotas: quotasCopy, CreatedAt: time.Now()}
	m.mu.Lock()
	m.current = s
	m.mu.Unlock()
	return s
}
func (m *Manager) Current() domain.Snapshot { m.mu.RLock(); defer m.mu.RUnlock(); return m.current }
func (m *Manager) Ack(client, version string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if client == "" || version == "" {
		return fmt.Errorf("ack fields required")
	}
	if m.current.Version == "" || version != m.current.Version {
		return fmt.Errorf("unknown snapshot version")
	}
	m.clients[client] = version
	return nil
}
func (m *Manager) ClientVersion(client string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[client]
}
func (m *Manager) VerifyNonce(nonce string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return hmac.Equal([]byte(nonce), []byte(m.current.Nonce))
}
