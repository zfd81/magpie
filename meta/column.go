package meta

const (
	DataTypeString  = "STRING"
	DataTypeInteger = "INT"
	DataTypeFloat   = "FLOAT"
	DataTypeBool    = "BOOL"
)

type ColumnInfo struct {
	Name         string      `json:"name"`
	Text         string      `json:"text,omitempty"`
	Comment      string      `json:"comment,omitempty"`
	DataType     string      `json:"dataType,omitempty"`
	Length       int         `json:"length,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Index        int         `json:"-"`
	Expression   string      `json:"expr,omitempty"`
}
