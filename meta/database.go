package meta

type DatabaseInfo struct {
	Name    string                `json:"name"`
	Text    string                `json:"text"`
	Comment string                `json:"comment,omitempty"`
	Charset string                `json:"charset"`
	Tables  map[string]*TableInfo `json:"-"`
}
