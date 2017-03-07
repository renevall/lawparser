package profile

import (
	"github.com/ReneVallecillo/office.go/domain"
)

//Domain holds reference to the domain struct
type Domain struct {
	service domain.UserStore
}

//Profile loads the user required to render frontend
func (d Domain) Profile(id uint64) (*domain.User, error) {
	var user *domain.User
	user, err := d.service.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil

}
