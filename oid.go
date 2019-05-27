package mgorm

import (
	"encoding/json"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

type OID bson.ObjectId

func (o OID) String() string {
	return bson.ObjectId(o).Hex()
}

func (o OID) ObjectID() bson.ObjectId {
	return bson.ObjectId(o)
}

func ParseOIDFromString(s string) (OID, error) {
	if !bson.IsObjectIdHex(s) {
		return "", errors.Errorf("invalid object id '%s'", s)
	}
	return OID(bson.ObjectIdHex(s)), nil
}

func (o *OID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	oid, err := ParseOIDFromString(s)
	if err != nil {
		return err
	}
	*o = oid
	return nil
}

func (o OID) MarshalJSON() ([]byte, error) {
	s := bson.ObjectId(o).Hex()
	return json.Marshal(s)
}
