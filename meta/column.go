package meta

type ColumnInfo struct {
	Name         string      `json:"name"`
	Text         string      `json:"text,omitempty"`
	Comment      string      `json:"comment,omitempty"`
	DataType     string      `json:"dataType,omitempty"`
	Length       int         `json:"length,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Index        int         `json:"index,omitempty"`
	Expression   string      `json:"expr,omitempty"`
}
