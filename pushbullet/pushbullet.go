package pushbullet

// User represents the User object for pushbullet
type User struct {
	Iden            string  `json:"iden"`
	Email           string  `json:"email"`
	EmailNormalized string  `json:"email_normalized"`
	Created         float64 `json:"created"`
	Modified        float64 `json:"modified"`
	Name            string  `json:"name"`
	ImageURL        string  `json:"image_url"`
}

type Pushes struct {
	Cursor    string `json:"cursor"`
	PushItems []Push `json:"pushes"`
}

type Push struct {
	Body  string `json:"body"`
	URL   string `json:"url"`
	Title string `json:"title"`
}
