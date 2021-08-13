package data

import (
	"sync"

	"github.com/gocql/gocql"
)

// CassandraStore is a CassandraDB data store which implements the Store interface
type CassandraStore struct {
	session *gocql.Session
	nextid  func() int
	sync.Mutex
}

// NewCassandraStore creates an instance of CassandraStore with the given connection string
func NewCassandraStore(connection string) (*CassandraStore, error) {
	cluster := gocql.NewCluster(connection)
	cluster.Consistency = gocql.One
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "cassandra", Password: "cassandra"}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS kittenserver WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};").Exec()
	if err != nil {
		return nil, err
	}
	err = session.Query("CREATE TABLE IF NOT EXISTS kittenserver.kittens (id int, name text , weight float, PRIMARY KEY(id, name))").Exec()
	if err != nil {
		return nil, err
	}
	err = session.Query("CREATE INDEX IF NOT EXISTS ON kittenserver.kittens(name)").Exec()
	if err != nil {
		return nil, err
	}
	return &CassandraStore{session: session, nextid: GeneraId()}, nil
}

// Search returns Kittens from the CassandraDB instance which have the name name
func (c *CassandraStore) SearchName(name string) Kittens {
	kittens := Kittens{}
	kitten := Kitten{}
	sc := c.session.Query("SELECT id, name, weight FROM kittenserver.kittens WHERE name = ?;").Bind(name).Iter().Scanner()
	for sc.Next() {
		err := sc.Scan(&kitten.Id, &kitten.Name, &kitten.Weight)
		if err != nil {
			return nil
		}
		kittens = append(kittens, kitten)
	}
	return kittens
}

// Search returns Kitten from the CassandraDB instance which have the id id
func (c *CassandraStore) SearchId(id int) Kitten {
	kitten := Kitten{}
	_ = c.session.Query("SELECT id, name, weight FROM kittenserver.kittens WHERE id = ?;").Bind(id).Scan(&kitten.Id, &kitten.Name, &kitten.Weight)
	return kitten
}

// DeleteAll deletes all the kittens from the datastore
func (c *CassandraStore) DeleteAll() {
	c.session.Query("TRUNCATE kittenserver.kittens").Exec()
}

// InsertKittens inserts a slice of kittens into the datastore
func (c *CassandraStore) Insert(kitten Kitten) int {
	c.Lock()
	i := c.nextid()
	c.Unlock()
	err := c.session.Query("INSERT INTO kittenserver.kittens (id, name, weight) VALUES (?,?,?)").Bind(i, kitten.Name, kitten.Weight).Exec()
	if err != nil {
		return 0
	}
	return i
}

// Search returns all Kittens from the CassandraDB instance
func (c *CassandraStore) List() Kittens {
	kittens := Kittens{}
	kitten := Kitten{}
	sc := c.session.Query("SELECT id, name, weight FROM kittenserver.kittens;").Iter().Scanner()
	for sc.Next() {
		err := sc.Scan(&kitten.Id, &kitten.Name, &kitten.Weight)
		if err != nil {
			return nil
		}
		kittens = append(kittens, kitten)
	}
	return kittens
}

// DeleteName deletes all the kittens which have the name name
func (c *CassandraStore) DeleteName(name string) bool {
	var id int
	ok := false
	sc := c.session.Query("SELECT id FROM kittenserver.kittens WHERE name = ?;").Bind(name).Iter().Scanner()
	for sc.Next() {
		err := sc.Scan(&id)
		if err != nil {
			return false
		}
		err = c.session.Query("DELETE FROM kittenserver.kittens WHERE id = ?;").Bind(id).Exec()
		if err != nil {
			return false
		}
		ok = true
	}
	return ok
}

// DeleteId deletes all the kittens which have the the id id
func (c *CassandraStore) DeleteId(id int) bool {
	idc := id
	err := c.session.Query("SELECT id FROM kittenserver.kittens WHERE id = ?;").Bind(idc).Scan(&idc)
	if err != nil {
		return false
	}
	err = c.session.Query("DELETE FROM kittenserver.kittens WHERE id = ?;").Bind(id).Exec()
	return err == nil
}
