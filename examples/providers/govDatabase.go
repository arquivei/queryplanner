package providers

import (
	"fmt"

	"github.com/arquivei/queryplanner"
)

type PersonalInfo struct {
	Name string
}

type GovDatabaseProvider struct {
	govDatabase map[string]PersonalInfo
}

func NewGovDatabaseProvider(repo map[string]PersonalInfo) queryplanner.FieldProvider {
	return &GovDatabaseProvider{
		govDatabase: repo,
	}
}

// Provides return an []queryplanner.Field.
// Each field has the function to calculate itself and populate the entries with the information.
func (p *GovDatabaseProvider) Provides() []queryplanner.Field {
	return []queryplanner.Field{
		{
			// Name is used to know what was provided.
			// When a provider depends on information from other providers,
			// the library will use the names to match
			Name: "Name",
			// The fill function receives the index of the document i and the execution context.
			// It must fill the document with the information.
			Fill: func(i int, ec queryplanner.ExecutionContext) error {
				fmt.Println("GovDatabaseProvider being executed")
				doc, ok := ec.Payload.Documents[i].(*Person) // Cast to a pointer so we can change the underlying struct
				if !ok {
					return fmt.Errorf("Error casting document to person struct.")
				}
				hadCovid, ok := p.govDatabase[*doc.CPF] // Note that the CPF is needed to query the gov database.
				if ok {
					doc.Name = &hadCovid.Name
				}
				return nil
			},
			// Clear function should remove the information from the document.
			// It is used to remove information that was not requested by the query.
			Clear: func(d queryplanner.Document) {
				doc, _ := d.(*Person)
				doc.Name = nil
			},
		},
	}
}

// DependsOn informs the library which information is needed before the provider can be executed. The names must match the ones defined in other providers.
func (p *GovDatabaseProvider) DependsOn() []queryplanner.FieldName {
	return []queryplanner.FieldName{
		"CPF", // Informs the library that this provider needs the CPF field to execute the query.
	}
}
