package contact

type Contact struct {
	WhatsApp  WhatsAppInfo `json:"whatsapp" firestore:"whatsapp"`
	Email     string       `json:"email" firestore:"email"`
	Instagram string       `json:"instagram" firestore:"instagram"`
	Address   string       `json:"address" firestore:"address"`
	Phone     string       `json:"phone" firestore:"phone"`
}

type WhatsAppInfo struct {
	Number string `json:"number" firestore:"number"`
	Link   string `json:"link" firestore:"link"`
}
