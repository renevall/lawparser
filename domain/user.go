package domain

//User Entity
type User struct {
	ID            uint32 `db:"user_id" json:"id"` // Don't use Id, use UserID() instead for consistency with MongoDB
	FirstName     string `db:"first_name" json:"first_name"`
	LastName      string `db:"last_name" json:"last_name"`
	Email         string `db:"email" json:"email"`
	Password      string `db:"password" json:"password"`
	StatusID      uint8  `db:"status_id" json:"status_id"`
	Address       string `db:"address" json:"address"`
	ContactNumber string `db:"contact_number" json:"contact_number"`
	GenderID      string `db:"gender_id" json:"gender_id"`
	PicURL        string `db:"pic_url" json:"pic_url"`
	UserLevel     string `db:"user_level" json:"user_level"`

	Customer Customer
	Token    string `json:"token"`
}

//UserRepository interface to be persisted/retrieved
type UserRepository interface {
	//Save(user User)
	FindByID(id uint32) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
}
