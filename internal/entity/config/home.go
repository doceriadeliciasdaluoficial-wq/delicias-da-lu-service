package config

type HomeConfig struct {
	FeaturedCakes []FeaturedCake `json:"featuredCakes" firestore:"featuredCakes"`
}

type FeaturedCake struct {
	ID            string                 `json:"id" firestore:"id"`
	Name          string                 `json:"name" firestore:"name"`
	DefaultWeight string                 `json:"defaultWeight" firestore:"defaultWeight"`
	DefaultConfig string                 `json:"defaultConfig" firestore:"defaultConfig"`
	BasePrice     float64                `json:"basePrice" firestore:"basePrice"`
	Tag           *string                `json:"tag,omitempty" firestore:"tag"`
	Image         string                 `json:"image" firestore:"image"`
	Description   string                 `json:"description" firestore:"description"`
	Config        map[string]interface{} `json:"config" firestore:"config"`
}
