package issue

import "time"

type ErrorInstance struct {
	RequestBody string    `firestore:"request_body" json:"request_body"`
	RequestDate time.Time `firestore:"request_date" json:"request_date"`
	RequestURL  string    `firestore:"request_url" json:"request_url"`
	Status      int       `firestore:"status" json:"status"`
	Title       string    `firestore:"title" json:"title"`
	TraceID     string    `firestore:"trace_id" json:"trace_id"`
	Type        string    `firestore:"type" json:"type"`
	UserAgent   string    `firestore:"user_agent" json:"user_agent"`
}

type ErrorType struct {
	Html      string    `firestore:"html"`
	UpdatedAt time.Time `firestore:"updatedAt"`
}
