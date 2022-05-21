package common

type Admin struct {
	name string
}

func NewAdmin(name string) *Admin {
	return &Admin{name: name}
}

func (a Admin) Recipient() string {
	return a.name
}
