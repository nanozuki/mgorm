package mgorm

import (
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// Model for collection
type Model struct {
	MongoAlias     string
	DBName         string
	CollectionName string
}

// NewModel make new model
func NewModel(alias, db, collection string) Model {
	return Model{
		MongoAlias:     alias,
		DBName:         db,
		CollectionName: collection,
	}
}

func (m *Model) wrapErr(err error, action string) error {
	return errors.Wrapf(err, "%s %s", m.CollectionName, action)
}

// CloseFunc use for close session
type CloseFunc func()

// C get original colletion for Model
func (m *Model) C() (*mgo.Collection, CloseFunc) {
	session := GetSession(m.MongoAlias)
	col := session.DB(m.DBName).C(m.CollectionName)
	return col, session.Close
}

// CreateIndex for Collection
func (m *Model) CreateIndex(index mgo.Index) error {
	c, done := m.C()
	defer done()
	return m.wrapErr(c.EnsureIndex(index), "create index")
}

// Insert documents
func (m *Model) Insert(docs ...interface{}) error {
	c, done := m.C()
	defer done()
	return m.wrapErr(c.Insert(docs), "insert")
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
		return m.wrapErr(err, "update")
	}
	return m.wrapErr(c.Update(selector, update), "update")
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
	return m.wrapErr(err, "update all")
}

// Delete document matching selector
func (m *Model) Delete(selector bson.M) error {
	c, done := m.C()
	defer done()
	return m.wrapErr(c.Remove(selector), "delete")
}

// DeleteAll find all document matching selector and delete them
func (m *Model) DeleteAll(selector bson.M) error {
	c, done := m.C()
	defer done()
	_, err := c.RemoveAll(selector)
	return m.wrapErr(err, "delete all")
}

// FindOne document by selector
func (m *Model) FindOne(result, selector interface{}) error {
	c, done := m.C()
	defer done()
	return m.wrapErr(c.Find(selector).One(result), "find one")
}

// FindAll documents matching selector
func (m *Model) FindAll(result, selector interface{}, opts ...FindOpt) error {
	c, done := m.C()
	defer done()
	query := c.Find(selector)
	query = applyFindOpts(query, opts...)
	return m.wrapErr(query.All(result), "find all")
}

// FindID find document by id
func (m *Model) FindID(result, id interface{}) error {
	return m.FindOne(result, bson.M{"_id": id})
}

// Count documents matching selector
func (m *Model) Count(selector interface{}) (int, error) {
	c, done := m.C()
	defer done()
	count, err := c.Find(selector).Count()
	return count, m.wrapErr(err, "count")
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

// IsErrNotFound check if not-found error
func IsErrNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}

// IsErrDuplicate check if duplicate-key error
func IsErrDuplicate(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}
