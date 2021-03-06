package data

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestHashJSON(t *testing.T) {
	hashIn := LiteralHash("abc")
	data, err := json.Marshal(hashIn)
	if err != nil {
		t.Fatalf("Marshal() failed: %s", err)
	}
	hashOut := &Hash{}
	if err := json.Unmarshal(data, hashOut); err != nil {
		t.Fatalf("Unmarshal() failed: %s", err)
	}
	if !hashIn.Equal(*hashOut) {
		t.Fatalf("Round-trip failed: got %s want %s", hashOut, hashIn)
	}
	if !reflect.DeepEqual(hashIn, *hashOut) {
		t.Fatalf("Round-trip DeepEqual failed: got %#v want %#v", hashOut, hashIn)
	}
}
func TestHashJSONEmpty(t *testing.T) {
	var hashIn Hash
	data, err := json.Marshal(hashIn)
	if err != nil {
		t.Fatalf("Marshal() failed: %s", err)
	}
	hashOut := &Hash{}
	if err := json.Unmarshal(data, hashOut); err != nil {
		t.Fatalf("Unmarshal() failed: %s", err)
	}
	if !hashIn.Equal(*hashOut) {
		t.Fatalf("Round-trip failed: got %s want %s", hashOut, hashIn)
	}
	if !reflect.DeepEqual(hashIn, *hashOut) {
		t.Fatalf("Round-trip DeepEqual failed: got %#v want %#v", hashOut, hashIn)
	}
}
func TestHashDatabase(t *testing.T) {
	hashIn := LiteralHash("abc")
	val, err := hashIn.Value()
	if err != nil {
		t.Fatalf("Value() failed: %s", err)
	}
	_, ok := val.([]byte)
	if !ok {
		t.Fatalf("Value() did not return bytes")
	}
	hashOut := &Hash{}
	if err := hashOut.Scan(val); err != nil {
		t.Fatalf("Scan() failed: %s", err)
	}
	if !hashIn.Equal(*hashOut) {
		t.Fatalf("Round-trip Equal failed: got %s want %s", hashOut, hashIn)
	}
	if !reflect.DeepEqual(hashIn, *hashOut) {
		t.Fatalf("Round-trip DeepEqual failed: got %#v want %#v", hashOut, hashIn)
	}
}
func TestHashDatabaseEmpty(t *testing.T) {
	var hashIn Hash
	val, err := hashIn.Value()
	if err != nil {
		t.Fatalf("Value() failed: %s", err)
	}
	_, ok := val.([]byte)
	if !ok {
		t.Fatalf("Value() did not return bytes")
	}
	hashOut := &Hash{}
	if err := hashOut.Scan(val); err != nil {
		t.Fatalf("Scan() failed: %s", err)
	}
	if !hashIn.Equal(*hashOut) {
		t.Fatalf("Round-trip Equal failed: got %s want %s", hashOut, hashIn)
	}
	if !reflect.DeepEqual(hashIn, *hashOut) {
		t.Fatalf("Round-trip DeepEqual failed: got %#v want %#v", hashOut, hashIn)
	}
}
func TestStoredJSON(t *testing.T) {
	storedIn := Stored{Type: JPG, Encoding: GZip}
	data, err := json.Marshal(storedIn)
	if err != nil {
		t.Fatalf("Marshal() failed: %s", err)
	}
	if string(data) != `"jpg.gz"` {
		t.Fatalf("want stringified representation got: %s", data)
	}
	storedOut := &Stored{}
	if err := json.Unmarshal(data, storedOut); err != nil {
		t.Fatalf("Unmarshal() failed: %s", err)
	}
	if storedIn != *storedOut {
		t.Fatalf("Round-trip failed: got %s want %s", storedOut, storedIn)
	}
	if !reflect.DeepEqual(storedIn, *storedOut) {
		t.Fatalf("Round-trip DeepEqual failed: got %#v want %#v", storedOut, storedIn)
	}
}
func TestStoredDatabase(t *testing.T) {
	storedIn := Stored{Type: JPG, Encoding: GZip}
	val, err := storedIn.Value()
	if err != nil {
		t.Fatalf("Value() failed: %s", err)
	}
	b, ok := val.([]byte)
	if !ok {
		t.Fatalf("Value() did not return bytes")
	}
	if string(b) != `jpg.gz` {
		t.Fatalf("want stringified representation got: %s", b)
	}
	storedOut := &Stored{}
	if err := storedOut.Scan(val); err != nil {
		t.Fatalf("Scan() failed: %s", err)
	}
	if storedIn != *storedOut {
		t.Fatalf("Round-trip Equal failed: got %s want %s", storedOut, storedIn)
	}
	if !reflect.DeepEqual(storedIn, *storedOut) {
		t.Fatalf("Round-trip DeepEqual failed: got %#v want %#v", storedOut, storedIn)
	}
}
