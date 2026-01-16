package entity

func NewDomains() []interface{} {
	return []interface{}{
		&Customer{},
		&Account{},
		&Token{},
	}
}
