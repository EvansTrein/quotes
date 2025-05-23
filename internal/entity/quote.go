package entity

type Quote struct {
	Author string `json:"author"`
	Text   string `json:"text"`
	ID     uint32 `json:"id"`
}
