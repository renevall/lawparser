package models


//AuthenticatedUser hold the data once a user is signed
type AuthenticatedUser struct {
	ID int `json:"id"`
	FirstName  string `json:"first_name"`
	LastName  string `json:"last_name"`	
	Status string `json:"status"`
	UserLevel string `json:"user_level"`
	Email string `json:"email"`
	Address string `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Gender string `json:"gender"`
	PicURL string `json:"pic_url"`
	Token string `json:"token"`
}
