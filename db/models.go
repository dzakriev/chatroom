package db

type User struct {
	Name     string
	Username string
	Id       int
}

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Room struct {
	Name string
	Id   int
}

type UserInRoom struct {
	UserId int
	RoomId int
}

type IMessageClient interface {
	Create(id string) Message
	Read(id string) (Message, bool)
	ReadAll() []Message
	Delete(id string) (Message, bool)
	DeleteAll() int
}
