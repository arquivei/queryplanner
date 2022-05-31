package providers

import (
	"fmt"

	"github.com/arquivei/queryplanner"
)

type CovidDatabaseProvider struct {
	covidRepository map[string]bool
}

// NewCovidDatabaseProvider You can use a constructor to pass in the dependencies to your provider, such as database connections.
func NewCovidDatabaseProvider(repo map[string]bool) queryplanner.FieldProvider {
	return &CovidDatabaseProvider{
		covidRepository: repo,
	}
}

// Provides return an []queryplanner.Field.
// Each field has the function to calculate itself and populate the entries with the information.
func (p *CovidDatabaseProvider) Provides() []queryplanner.Field {
	return []queryplanner.Field{
		{
			// Name is used to know what was provided.
			// When a provider depends on information from other providers,
			// the library will use the names to match
			Name: "HadCovid",
			// The fill function receives the index of the document i and the execution context.
			// It must fill the document with the information.
			Fill: func(i int, ec queryplanner.ExecutionContext) error {
				fmt.Println("CovidDatabaseProvider being executed")
				doc, ok := ec.Payload.Documents[i].(*Person) // Cast to a pointer so we can change the underlying struct
				if !ok {
					return fmt.Errorf("Error casting document to person struct.")
				}
				hadCovid, ok := p.covidRepository[*doc.CPF] // Note that the CPF is needed to query the covid database.
				if ok {
					doc.HadCovid = &hadCovid
				}
				return nil
			},
			// Clear function should remove the information from the document.
			// It is used to remove information that was not requested by the query.
			Clear: func(d queryplanner.Document) {
				doc, _ := d.(*Person)
				doc.HadCovid = nil
			},
		},
	}
}

// DependsOn informs the library which information is needed before the provider can be executed. The names must match the ones defined in other providers.
func (p *CovidDatabaseProvider) DependsOn() []queryplanner.FieldName {
	return []queryplanner.FieldName{
		"CPF", // Informs the library that this provider needs the CPF field to execute the query.
	}
}
