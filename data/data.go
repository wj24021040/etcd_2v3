package data

import (
	"time"
)

type BackupStrategy struct {
	Keys      []string `json:"keys,omitempty"`
	Sorted    bool     `json:"sorted,omitempty"`
	Recursive bool     `json:"recursive,omitempty"`
}

type BackupKey struct {
	Key        string     `json:"key"`
	Value      *string    `json:"value,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
	TTL        int64      `json:"ttl,omitempty"`
}
