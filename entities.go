package queryplanner

import (
	"context"
)

// Request is an interface for getting the requested fields of a Request
// and also for retrieving an inner Request structure as well.
type Request interface {
	GetRequestedFields() []string
}

// Document is an interface that represents a generic data structure.
type Document interface{}

// Field represents a valid field. It has a name and
// functions for filling and cleaning itself.
type Field struct {
	Name  FieldName
	Fill  func(int, ExecutionContext) error
	Clear func(Document)
}

// Index represents a valid index. It has a name and
// functions for cleaning itself.
type Index struct {
	Name  FieldName
	Clear func(Document)
}

// ExecutionContext is used during the filling process. It stores essential
// data structures to executing the enrichment processes.
type ExecutionContext struct {
	Context context.Context
	Request Request
	cache   *Cache
	Payload *Payload
}

// Cache instatiates a new Cache.
func (e *ExecutionContext) Cache() *Cache {
	if e.cache == nil {
		e.cache = newCache()
	}
	return e.cache
}

func newCache() *Cache {
	return &Cache{cache: make(map[interface{}]*CacheEntry)}
}

// CacheEntryLoader is a function that caches the result of the first call of function.
type CacheEntryLoader func() (interface{}, error)

// Cache caches the result of a function.
type Cache struct {
	cache map[interface{}]*CacheEntry
}

// GetOrLoad tries to retrieve an existing element from the cache by an indexed `key` . If there is already an entry for
// (CacheEntry) the key, the cached content is returned. If there is not a cached value, then `loader` (CacheEntryLeader)
// is executed and its results are cached using the provided `key` as index.
func (c *Cache) GetOrLoad(key interface{}, loader CacheEntryLoader) (interface{}, error) {
	result, ok := c.cache[key]
	if ok {
		return result.data, result.err
	}
	data, err := loader()
	c.cache[key] = &CacheEntry{
		data: data,
		err:  err,
	}
	return data, err
}

// CacheEntry is the stored element in Cache.
type CacheEntry struct {
	data interface{}
	err  error
}

// Payload stores the slice of documents and also supports an arbitrary data
// to be used when necessary.
type Payload struct {
	Documents  []Document
	CustomData interface{}
}

// FieldName is a string representing a valid field.
type FieldName string
