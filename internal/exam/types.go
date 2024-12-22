package exam

type Exam struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	TopicIDs []string `json:"topic_ids"`
}
