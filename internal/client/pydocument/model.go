package client

type Chunk struct {
	Page int    `json:"page"`
	Text string `json:"text"`
}
type DocumentSplitResponse struct {
	Filename string  `json:"filename"`
	Chunks   []Chunk `json:"chunks"`
}
