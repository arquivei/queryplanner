package queryplanner

import "context"

// FieldProvider is able to load an existing set of []Document with certain
// fields. The `Provides()` returns a list of Field's, which in turn contains
// methods to filling and cleaning a Document.
type FieldProvider interface {
	Provides() []Field
	DependsOn() []FieldName
}

// IndexProvider is similar to FieldProvider in some aspects: it also provides
// specifically defined fields and encapsulates the mechanics of retrieving and
// populating them. However, the IndexProvider should not depend on anyone else,
// as FieldProvider does.  Also, an IndexProvider has an additional responsibility:
// to create the initial base of documents (encapsulated inside Payload) to be enriched
// by the FieldProviders
type IndexProvider interface {
	Execute(ctx context.Context, request Request, fields []string) (*Payload, error)
	Provides() []Index
}
