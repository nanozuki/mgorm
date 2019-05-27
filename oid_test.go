package mgorm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
)

func TestMarshalJSON(t *testing.T) {
	type foo struct {
		ID OID `json:"id"`
	}
	bid := bson.NewObjectId()
	foo1 := &foo{ID: OID(bid)}
	bytes, _ := json.Marshal(foo1)
	str := string(bytes)
	want := fmt.Sprintf(`{"id":"%s"}`, bid.Hex())
	if str != want {
		t.Errorf("Marshal() = '%s', want = '%s'", str, want)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	type foo struct {
		ID OID `json:"id"`
	}
	bid := bson.NewObjectId()
	str := fmt.Sprintf(`{"id":"%s"}`, bid.Hex())
	var got foo
	json.Unmarshal([]byte(str), &got)
	want := foo{ID: OID(bid)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal() = %+v, want %+v", got, want)
	}
}

func TestBSON(t *testing.T) {
	type foo struct {
		ID OID `bson:"_id"`
		A  int `bson:"a"`
	}
	foo1 := foo{ID: OID(bson.NewObjectId()), A: 5}
	bytes, err := bson.Marshal(&foo1)
	if err != nil {
		t.Errorf("marshal bson failed: %v", err)
	}
	var foo2 foo
	if err := bson.Unmarshal(bytes, &foo2); err != nil {
		t.Errorf("unmarshal bson failed: %v", err)
	}
	if !reflect.DeepEqual(foo1, foo2) {
		t.Errorf("marshal unmarshal bson failed, before: %v, after: %v", foo1, foo2)
	}
}
