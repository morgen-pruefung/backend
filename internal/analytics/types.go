package analytics

type Event struct {
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties"`
}

type PageVisitRequest struct {
	URL     string `json:"url"`
	Referer string `json:"referer"`
}
