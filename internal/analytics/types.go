package analytics

type Event struct {
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties"`
}
