package mgorm

import (
	"sync"

	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
)

var client = struct {
	sessions map[string]*mgo.Session
	mutex    sync.RWMutex
}{
	sessions: make(map[string]*mgo.Session),
}

// Init load configs:
// usls: <mongo-alias:url>
func Init(urls map[string]string) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	for alias, url := range urls {
		session, err := mgo.Dial(url)
		if err != nil {
			return errors.Wrapf(err, "connect to mongodb: %s", url)
		}
		client.sessions[alias] = session
	}
	return nil
}

func InitWithInfo(infos map[string]*mgo.DialInfo) error {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	for alias, info := range infos {
		session, err := mgo.DialWithInfo(info)
		if err != nil {
			return errors.Wrapf(err, "connect to mongodb: %v", info)
		}
		client.sessions[alias] = session
	}
	return nil
}

// GetSession return an copy of original session for specified server alias
func GetSession(alias string) *mgo.Session {
	client.mutex.RLock()
	defer client.mutex.RUnlock()

	session, ok := client.sessions[alias]
	if !ok {
		panic(errors.Errorf("can't find alias '%s'", alias))
	}
	return session.Copy()
}

// GetSession return original session for specified server alias
func GetRawSession(alias string) *mgo.Session {
	client.mutex.RLock()
	defer client.mutex.RUnlock()

	session, ok := client.sessions[alias]
	if !ok {
		panic(errors.Errorf("can't find alias '%s'", alias))
	}
	return session
}
