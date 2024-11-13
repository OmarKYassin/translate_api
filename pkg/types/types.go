package types

type entry struct {
	Speaker  string `json:"speaker" binding:"required"`
	Time     string `json:"time" binding:"required"`
	Sentence string `json:"sentence" binding:"required"`
}

type Transcript []entry
