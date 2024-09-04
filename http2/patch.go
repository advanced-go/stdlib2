package http2

const (
	OpTest    = "test"
	OpRemove  = "remove"
	OpAdd     = "add"
	OpReplace = "replace"
	OpMove    = "move"
	OpCopy    = "copy"
)

type Operation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type Patch struct {
	Updates []Operation
}
