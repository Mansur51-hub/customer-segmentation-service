package model

import "time"

type SegmentMembership struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	SegmentSlug string    `json:"segment_slug"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}
