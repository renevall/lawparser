package profile

import (
	"github.com/ReneVallecillo/office.go/domain"
)

//Domain holds reference to the domain struct
type Domain struct {
	UserRepository domain.UserRepository
}

//Profile loads the user required to render frontend
func (d Domain) Profile(id uint32) (*domain.User, error) {
	var user *domain.User
	user, err := d.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil

}
