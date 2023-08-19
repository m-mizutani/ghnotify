package model

type RegoInput struct {
	Name  string      `json:"name"`
	Event interface{} `json:"event"`
}

type RegoResult struct {
	Notify []*Notify `json:"notify"`
}

type Notify struct {
	Channel string         `json:"channel"`
	Text    string         `json:"text"`
	Body    string         `json:"body"`
	Color   string         `json:"color"`
	Fields  []*NotifyField `json:"fields"`
}

type NotifyField struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
	URL   string `json:"url"`
}
