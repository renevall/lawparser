package domain

//Customer domain struc
type Customer struct {
	ID    int
	Name  string
	Email string
}

//CustomerRepository interface to be persisted/retrieved
type CustomerStore interface {
	Save(customer Customer)
	Find(id int) *Customer
}
