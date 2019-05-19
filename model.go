package mgorm

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Model for collection
type Model struct {
	MongoName      string
	DBName         string
	CollectionName string
}

// CloseFunc use for close session
type CloseFunc func()

// C get original colletion for Model
func (m *Model) C() (*mgo.Collection, CloseFunc) {
	session := GetSession(m.MongoName)
	col := session.DB(m.DBName).C(m.CollectionName)
	return col, session.Close
}

// CreateIndex for Collection
func (m *Model) CreateIndex(index mgo.Index) error {
	c, done := m.C()
	defer done()
	return c.EnsureIndex(index)
}

// Create documents
func (m *Model) Create(docs ...interface{}) error {
	c, done := m.C()
	defer done()
	return c.Insert(docs)
}

// Update document matching selector
func (m *Model) Update(selector, update interface{}, opts ...UpdateOpt) error {
	c, done := m.C()
	defer done()

	uo := &updateOpts{}
	for _, opt := range opts {
		opt(uo)
	}
	if uo.Upsert {
		_, err := c.Upsert(selector, update)
		return err
	}
	return c.Update(selector, update)
}

// UpdateID updata document matching ID
func (m *Model) UpdateID(id, update interface{}, opts ...UpdateOpt) error {
	selector := bson.M{"_id": id}
	return m.Update(selector, update, opts...)
}

// UpdateAll finds all documents matching the provided selector document and update them
func (m *Model) UpdateAll(selector, update interface{}) error {
	c, done := m.C()
	defer done()
	_, err := c.UpdateAll(selector, update)
	return err
}

// Delete document matching selector
func (m *Model) Delete(selector bson.M) error {
	c, done := m.C()
	defer done()
	return c.Remove(selector)
}

// DeleteAll find all document matching selector and delete them
func (m *Model) DeleteAll(selector bson.M) error {
	c, done := m.C()
	defer done()
	_, err := c.RemoveAll(selector)
	return err
}

// FindOne document by selector
func (m *Model) FindOne(result, selector interface{}) error {
	c, done := m.C()
	defer done()
	return c.Find(selector).One(result)
}

// FindAll documents matching selector
func (m *Model) FindAll(result, selector interface{}, opts ...FindOpt) error {
	c, done := m.C()
	defer done()
	query := c.Find(selector)
	query = applyFindOpts(query, opts...)
	return query.All(result)
}

// FindID find document by id
func (m *Model) FindID(result, id interface{}) error {
	return m.FindOne(result, bson.M{"_id": id})
}

// Count documents matching selector
func (m *Model) Count(selector interface{}) (int, error) {
	c, done := m.C()
	defer done()
	return c.Find(selector).Count()
}

// GetIter make an iter for selector
func (m *Model) GetIter(selector interface{}) *Iter {
	c, done := m.C()
	ogIter := c.Find(selector).Iter()
	return &Iter{
		closeSession: done,
		Iter:         ogIter,
	}
}
