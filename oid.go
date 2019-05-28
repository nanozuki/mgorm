package mgorm

import (
	"encoding/json"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type OID bson.ObjectId

func (o OID) String() string {
	return bson.ObjectId(o).Hex()
}

func (o OID) ObjectID() bson.ObjectId {
	return bson.ObjectId(o)
}

func ParseOIDFromString(s string) (OID, error) {
	if s == "" {
		return OID(s), nil
	}
	if !bson.IsObjectIdHex(s) {
		return "", errors.Errorf("invalid object id '%s'", s)
	}
	return OID(bson.ObjectIdHex(s)), nil
}

// BSON

func (o *OID) GetBSON() (interface{}, error) {
	s := o.ObjectID()
	return s, nil
}

func (o *OID) SetBSON(raw bson.Raw) error {
	var bid bson.ObjectId
	if err := raw.Unmarshal(&bid); err != nil {
		return err
	}
	*o = OID(bid)
	return nil
}

// JSON

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

// Msgpack

func (o *OID) EncodeMsgpack(enc *msgpack.Encoder) error {
	s := o.ObjectID().Hex()
	return enc.Encode(s)
}

func (o *OID) DecodeMsgpack(dec *msgpack.Decoder) error {
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	oid, err := ParseOIDFromString(s)
	if err != nil {
		return err
	}
	*o = oid
	return nil
}
