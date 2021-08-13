package data

import (
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStore is a MongoDB data store which implements the Store interface
type MongoStore struct {
	session *mgo.Session
	nextid  func() int
	sync.Mutex
}

// NewMongoStore creates an instance of MongoStore with the given connection string
func NewMongoStore(connection string) (*MongoStore, error) {
	session, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &MongoStore{session: session, nextid: GeneraId()}, nil
}

// Search returns Kittens from the MongoDB instance which have the name name
func (m *MongoStore) SearchName(name string) Kittens {
	s := m.session.Clone()
	defer s.Close()
	var results Kittens
	c := s.DB("kittenserver").C("kittens")
	err := c.Find(bson.M{"name": name}).All(&results)
	if err != nil {
		return nil
	}
	return results
}

// Search returns Kitten from the MongoDB instance which have the id id
func (m *MongoStore) SearchId(id int) Kitten {
	s := m.session.Clone()
	defer s.Close()
	var result Kitten
	c := s.DB("kittenserver").C("kittens")
	_ = c.Find(bson.M{"id": id}).One(&result)
	return result
}

// DeleteAll deletes all the kittens from the datastore
func (m *MongoStore) DeleteAll() {
	s := m.session.Clone()
	defer s.Close()
	s.DB("kittenserver").C("kittens").DropCollection()
}

// InsertKittens inserts a slice of kittens into the datastore
func (m *MongoStore) Insert(kitten Kitten) int {
	s := m.session.Clone()
	defer s.Close()
	m.Lock()
	i := m.nextid()
	m.Unlock()
	kitten.Id = i
	s.DB("kittenserver").C("kittens").Insert(kitten)
	return i
}

// Search returns all Kittens from the MongoDB instance
func (m *MongoStore) List() Kittens {
	s := m.session.Clone()
	defer s.Close()
	var results Kittens
	c := s.DB("kittenserver").C("kittens")
	err := c.Find(bson.M{}).All(&results)
	if err != nil {
		return nil
	}
	return results
}

// DeleteName deletes all the kittens which have the name name
func (m *MongoStore) DeleteName(name string) bool {
	s := m.session.Clone()
	defer s.Close()
	c := s.DB("kittenserver").C("kittens")
	err := c.Remove((bson.M{"name": name}))
	return err == nil
}

// DeleteId deletes all the kittens which have the the id id
func (m *MongoStore) DeleteId(id int) bool {
	s := m.session.Clone()
	defer s.Close()
	c := s.DB("kittenserver").C("kittens")
	err := c.Remove((bson.M{"id": id}))
	return err == nil
}
