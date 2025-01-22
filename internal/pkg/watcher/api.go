package watcher

type Counter struct {
	Iteration int    `json:"iteration"`
	HexStr    string `json:"value"`
}

type CounterReset struct{}
