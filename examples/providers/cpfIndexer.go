package providers

import (
	"context"
	"fmt"

	"github.com/arquivei/queryplanner"
)

// An indexer is a special type of provider.
// It must initialize the documents that will be enriched, being the first query to be made.
type CPFIndexer struct{}

// NewCovidDatabaseProvider You can use a constructor to pass in the dependencies to your provider, such as database connections.
func NewCPFIndexer() queryplanner.IndexProvider {
	return &CPFIndexer{}
}

// Provides return an []queryplanner.Index.
// The index field only has the clear method, since it will be populated in the Execute function of the IndexProvider
func (p *CPFIndexer) Provides() []queryplanner.Index {
	return []queryplanner.Index{
		{
			// Name is used to know what was provided.
			// When a provider depends on information from other providers,
			// the library will use the names to match
			Name: "CPF",
			// Clear function should remove the information from the document.
			// It is used to remove information that was not requested by the query.
			Clear: func(d queryplanner.Document) {
				doc, _ := d.(*Person)
				doc.CPF = nil
			},
		},
	}
}

// Execute is the first function to be executed in a query plan. It will return the fields that will pass through each step to the enrichment process.
// The request param can be a struct that might carry needed information for your service, such as pagination, filters, etc.
func (p *CPFIndexer) Execute(ctx context.Context, request queryplanner.Request, fields []string) (*queryplanner.Payload, error) {
	fmt.Println("CPFIndexer being executed")
	return &queryplanner.Payload{
		Documents: []queryplanner.Document{
			&Person{
				CPF: ref("44452427138"),
			},
			&Person{
				CPF: ref("85022625806"),
			},
		},
	}, nil
}
