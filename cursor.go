package mgorm

import "github.com/globalsign/mgo"

// Iter wrap original iter, add close function
type Iter struct {
	*mgo.Iter
	closeSession CloseFunc
}

// Close ...
func (it *Iter) Close() error {
	err := it.Iter.Close()
	it.closeSession()
	return err
}
