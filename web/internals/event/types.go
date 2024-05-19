package event

type Message struct {
	Msg string
}

type SentimentalResult struct {
	Label string
	Score float64
}
