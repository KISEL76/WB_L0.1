package model

type Item struct {
	ChrtID      int64   `json:"chrt_id"`
	Price       float64 `json:"price"`
	RID         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int16   `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float64 `json:"total_price"`
	NmID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int16   `json:"status"`
	TrackNumber string  `json:"track_number"`
}
