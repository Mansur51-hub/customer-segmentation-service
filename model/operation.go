package model

import "time"

type Operation struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	SegmentSlug string    `json:"segmentSlug"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
}

const (
	CustomerAddToSegmentOperationType      = "adding"
	CustomerRemoveFromSegmentOperationType = "removing"
)
