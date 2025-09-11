package network

type User struct {
	Name string
	Id   uint32
}

func NewUser(name string, id uint32) User {
	return User{
		Name: name,
		Id:   id,
	}
}
