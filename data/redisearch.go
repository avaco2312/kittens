package data

import (
	"strconv"
	"sync"

	"github.com/RediSearch/redisearch-go/redisearch"
)

// RediSearchStore is a RediSearchDB data store which implements the Store interface
type RediSearchStore struct {
	session *redisearch.Client
	nextid  func() int
	sync.Mutex
}

// NewRediSearchStore creates an instance of RediSearchStore with the given connection string
func NewRediSearchStore(connection string) (*RediSearchStore, error) {
	session := redisearch.NewClient(connection+":6379", "kittensIndex")
	sc := redisearch.NewSchema(redisearch.DefaultOptions).AddField(redisearch.NewTextField("name"))
	session.Drop()
	if err := session.CreateIndex(sc); err != nil {
		return nil, err
	}
	return &RediSearchStore{session: session, nextid: GeneraId()}, nil
}

// Search returns Kittens from the RediSearchDB instance which have the name name
func (r *RediSearchStore) SearchName(name string) Kittens {
	docs, _, _ := r.session.Search(redisearch.NewQuery(name).SetReturnFields("id", "name", "weight"))
	if len(docs) == 0 {
		return nil
	}
	kittens := Kittens{}
	kitten := Kitten{}
	for _, s := range docs {
		kitten.Id, _ = strconv.Atoi(s.Properties["id"].(string))
		kitten.Name = s.Properties["name"].(string)
		f, _ := strconv.ParseFloat(s.Properties["weight"].(string), 32)
		kitten.Weight = float32(f)
		kittens = append(kittens, kitten)
	}
	return kittens
}

// Search returns Kitten from the RediSearchDB instance which have the id id
func (r *RediSearchStore) SearchId(id int) *Kitten {
	kitten := Kitten{}
	s, _ := r.session.Get(strconv.Itoa(id))
	if s == nil {
		return nil
	}
	kitten.Id, _ = strconv.Atoi(s.Properties["id"].(string))
	kitten.Name = s.Properties["name"].(string)
	f, _ := strconv.ParseFloat(s.Properties["weight"].(string), 32)
	kitten.Weight = float32(f)
	return &kitten
}

// DeleteAll deletes all the kittens from the datastore
func (r *RediSearchStore) DeleteAll() {
	r.session.Drop()
	sc := redisearch.NewSchema(redisearch.DefaultOptions).AddField(redisearch.NewTextField("name"))
	r.session.CreateIndex(sc)
}

// InsertKittens inserts a slice of kittens into the datastore
func (r *RediSearchStore) Insert(kitten Kitten) int {
	r.Lock()
	i := r.nextid()
	r.Unlock()
	kitten.Id = i
	doc := redisearch.NewDocument(strconv.Itoa(kitten.Id), 1.0)
	doc.Set("id", kitten.Id).Set("id", kitten.Id).Set("name", kitten.Name).Set("weight", kitten.Weight)
	err := r.session.Index([]redisearch.Document{doc}...)
	if err != nil {
		return 0
	}
	return i
}

// Search returns all Kittens from the RediSearchDB instance
func (r *RediSearchStore) List() Kittens {
	kittens := Kittens{}
	kitten := Kitten{}
	docs, _, err := r.session.Search(redisearch.NewQuery("*").SetReturnFields("id", "name", "weight"))
	if err != nil {
		return nil
	}
	for _, s := range docs {
		kitten.Id, _ = strconv.Atoi(s.Properties["id"].(string))
		kitten.Name = s.Properties["name"].(string)
		f, _ := strconv.ParseFloat(s.Properties["weight"].(string), 32)
		kitten.Weight = float32(f)
		kittens = append(kittens, kitten)
	}
	return kittens
}

// DeleteName deletes all the kittens which have the name name
func (r *RediSearchStore) DeleteName(name string) bool {
	ok := false
	docs, _, _ := r.session.Search(redisearch.NewQuery(name).SetReturnFields("id"))
	for _, s := range docs {
		err := r.session.DeleteDocument(s.Id)
		if err != nil {
			return false
		}
		ok = true
	}
	return ok
}

// DeleteId deletes all the kittens which have the the id id
func (r *RediSearchStore) DeleteId(id int) bool {
	s, err := r.session.Get(strconv.Itoa(id))
	if err != nil || s == nil {
		return false
	}
	err = r.session.DeleteDocument(strconv.Itoa(id))
	return err == nil
}
