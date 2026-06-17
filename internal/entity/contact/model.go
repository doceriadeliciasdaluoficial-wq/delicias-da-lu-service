package contact

type Contact struct {
	WhatsApp  WhatsAppInfo  `json:"whatsapp" firestore:"whatsapp"`
	Email     EmailInfo     `json:"email" firestore:"email"`
	Instagram InstagramInfo `json:"instagram" firestore:"instagram"`
	Address   *string       `json:"address,omitempty" firestore:"address"`
	Phone     *string       `json:"phone,omitempty" firestore:"phone"`
	Location  *Location     `json:"location,omitempty" firestore:"location"`
}

type WhatsAppInfo struct {
	Number  string            `json:"number" firestore:"number"`
	Link    string            `json:"link" firestore:"link"`
	Display *string           `json:"display,omitempty" firestore:"display"`
	Message *WhatsAppMessages `json:"message,omitempty" firestore:"message"`
}

type WhatsAppMessages struct {
	Default *string `json:"default,omitempty" firestore:"default"`
	Order   *string `json:"order,omitempty" firestore:"order"`
	Custom  *string `json:"custom,omitempty" firestore:"custom"`
}

type EmailInfo struct {
	Address string  `json:"address" firestore:"address"`
	Subject *string `json:"subject,omitempty" firestore:"subject"`
}

type InstagramInfo struct {
	Handle   *string `json:"handle,omitempty" firestore:"handle"`
	URL      *string `json:"url,omitempty" firestore:"url"`
	EmbedURL *string `json:"embedUrl,omitempty" firestore:"embedUrl"`
}

type Location struct {
	Name        *string      `json:"name,omitempty" firestore:"name"`
	Address     *string      `json:"address,omitempty" firestore:"address"`
	Coordinates *Coordinates `json:"coordinates,omitempty" firestore:"coordinates"`
	MapsURL     *string      `json:"mapsUrl,omitempty" firestore:"mapsUrl"`
}

type Coordinates struct {
	Lat float64 `json:"lat" firestore:"lat"`
	Lng float64 `json:"lng" firestore:"lng"`
}
