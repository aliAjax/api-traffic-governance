package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type AuditEntry struct {
	ID           string
	Action       string
	Actor        string
	Resource     string
	Before       string
	After        string
	PreviousHash string
	Hash         string
	CreatedAt    time.Time
}

func (e AuditEntry) canonical() []byte {
	return []byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", e.ID, e.Action, e.Actor, e.Resource, e.Before, e.After, e.PreviousHash))
}

func (e *AuditEntry) Seal() {
	s := sha256.Sum256(e.canonical())
	e.Hash = hex.EncodeToString(s[:])
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
}
func (e AuditEntry) Verify() bool {
	s := sha256.Sum256(e.canonical())
	return e.Hash == hex.EncodeToString(s[:])
}
