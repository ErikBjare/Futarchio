package db

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"log"
)

// NOTE: This is a test to generalize the Google Datastore API, doesn't look like it'll work though.

// All non-saved fields should be non-exportable
type BaseEntity struct {
	kind string
	key  *datastore.Key
}

type Entity interface {
	Save(appengine.Context) error
	Key() *datastore.Key
}

func (e *BaseEntity) Key() *datastore.Key {
	return e.key
}

func (e *BaseEntity) Log() {
	log.Printf("", e)
}

func (e *BaseEntity) Save(c appengine.Context) error {
	key := e.Key()
	if key == nil {
		return errors.New("entity doesn't have a key")
	}

	key, err := datastore.Put(c, key, e)
	e.key = key
	return err
}

func (e *BaseEntity) Put(c appengine.Context, parent *datastore.Key) error {
	key := e.Key()
	if key != nil {
		return errors.New("entity already exists in datastore")
	}

	e.key = datastore.NewIncompleteKey(c, e.kind, parent)
	return e.Save(c)
}

func RunQuery(c appengine.Context, q *datastore.Query, entities []BaseEntity) ([]BaseEntity, error) {
	// TODO: How can this be done? Skip?
	// Maybe if we register models and
	keys, err := q.GetAll(c, &entities)
	for i := range keys {
		entities[i].key = keys[i]
	}
	return entities, err
}
