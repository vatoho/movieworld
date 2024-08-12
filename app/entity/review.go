package entity

type Review struct {
	ID      uint64
	Mark    uint32
	Comment string
	Author  *User
}
