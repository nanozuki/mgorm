package mgorm

import (
	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
)

// Iter wrap original iter, add close function
type Iter struct {
	*mgo.Iter
	closeSession CloseFunc
}

// Close ...
func (it *Iter) Close() error {
	err := it.Iter.Close()
	it.closeSession()
	return errors.Wrap(err, "close iter")
}
