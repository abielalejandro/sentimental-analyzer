package storage

import "time"

type SentimentalResult struct {
	Label string
	Score float64
}

type Message struct {
	Id          string
	Msg         string
	MsgAnalyzed SentimentalResult
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
}
