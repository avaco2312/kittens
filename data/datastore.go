package data

// Store is an interface used for interacting with the backend datastore
type Store interface {
	List() Kittens
	SearchName(name string) Kittens
	DeleteAll()
	Insert(kitten Kitten) int
	DeleteName(name string) bool
	SearchId(id int) Kitten
	DeleteId(id int) bool
}

type Kitten struct {
	Id     int  `json:"id" bson:"id"`
	Name   string  `json:"name" bson:"name"`
	Weight float32 `json:"weight" bson:"weight"`
}

type Kittens []Kitten 

func GeneraId() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}





