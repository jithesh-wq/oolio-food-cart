package models

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Image    image   `json:"image"`
}

type image struct {
	Mobile    string `json:"mobile"`
	Thumbnail string `json:"thumbnail"`
	Desktop   string `json:"desktop"`
	Tablet    string `json:"tablet,omitempty"` // Optional field for tablet image
}
