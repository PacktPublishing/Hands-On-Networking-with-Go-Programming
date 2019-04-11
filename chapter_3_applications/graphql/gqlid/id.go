package gqlid

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
)

// Type of the ID.
type Type int

const (
	_ = iota
	// TypeInt64 means that the underlying ID is an int64 value.
	TypeInt64 Type = 1
	// TypeInt32 means that the underlying ID is an int32 value.
	TypeInt32
	// TypeString means that the underlying ID is a string value.
	TypeString
)

// Int64 ID.
func Int64(service, resource string, id int64) *ID {
	return &ID{
		Service:    service,
		Resource:   resource,
		Type:       TypeInt64,
		ValueInt64: id,
	}
}

// Int32 ID.
func Int32(service, resource string, id int32) *ID {
	return &ID{
		Service:    service,
		Resource:   resource,
		Type:       TypeInt32,
		ValueInt32: id,
	}
}

// String ID.
func String(service, resource string, id string) *ID {
	return &ID{
		Service:     service,
		Resource:    resource,
		Type:        TypeString,
		ValueString: id,
	}
}

// ID of an item within a GraphQL service.
type ID struct {
	Service     string
	Resource    string
	Type        Type
	ValueInt64  int64
	ValueInt32  int32
	ValueString string
}

// Int64 gets the value of the ID.
func (id *ID) Int64() (value int64, ok bool) {
	if id.Type == TypeInt64 {
		return id.ValueInt64, true
	}
	return
}

// Int32 gets the value of the ID.
func (id *ID) Int32() (value int32, ok bool) {
	if id.Type == TypeInt32 {
		return id.ValueInt32, true
	}
	return
}

// String gets the value of the ID.
func (id *ID) String() (value string, ok bool) {
	if id.Type == TypeString {
		return id.ValueString, true
	}
	return
}

// ErrNotValidID is returned when the ID if not in the required format (base64 encoded gob).
var ErrNotValidID = errors.New("ID: invalid ID format")

// Parse an ID from a base64 encoded GraphQL ID.
func Parse(s string) (id *ID, err error) {
	ids, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		err = ErrNotValidID
		return
	}
	d := gob.NewDecoder(bytes.NewReader(ids))
	err = d.Decode(&id)
	return
}

// ParseInt64For parses an ID and returns an error if it's not an int64.
func ParseInt64For(service, resource string, s string) (idValue int64, err error) {
	id, err := Parse(s)
	if err != nil {
		return
	}
	if id.Service != service || id.Resource != resource {
		err = fmt.Errorf("expected service/resource of '%s/%s', got '%s/%s'", service, resource, id.Service, id.Resource)
		return
	}
	idValue, ok := id.Int64()
	if !ok {
		err = fmt.Errorf("expected Int64 ID, got %v", id.Type)
		return
	}
	return
}

// Encoded returns the base64 encoded ID or panics if an error occurred.
func (id *ID) Encoded() (s string) {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	if err := e.Encode(id); err != nil {
		panic(err.Error())
	}
	s = base64.StdEncoding.EncodeToString(w.Bytes())
	return
}
