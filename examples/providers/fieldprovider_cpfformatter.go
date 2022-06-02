package providers

import (
	"fmt"

	"github.com/arquivei/queryplanner"
)

type CPFFormatterProvider struct{}

// Provides return an []queryplanner.Field.
// Each field has the function to calculate itself and populate the entries with the information.
func (p *CPFFormatterProvider) Provides() []queryplanner.Field {
	return []queryplanner.Field{
		{
			Name: "CPF", // Modifies an already existing field
			Fill: func(i int, ec queryplanner.ExecutionContext) error {
				fmt.Println("CPFFormatterProvider being executed")
				doc, ok := ec.Payload.Documents[i].(*Person) // Cast to a pointer so we can change the underlying struct
				if !ok {
					return fmt.Errorf("error casting document to person struct")
				}
				fmt.Println(ref((*doc.CPF)[0:9] + "-" + (*doc.CPF)[9:]))
				doc.CPF = ref((*doc.CPF)[0:9] + "-" + (*doc.CPF)[9:])
				return nil
			},
			// Clear function should remove the information from the document.
			// It is used to remove information that was not requested by the query.
			Clear: func(d queryplanner.Document) {
				doc, _ := d.(*Person)
				doc.CPF = nil
			},
		},
	}
}

// DependsOn informs the library which information is needed before the provider can be executed. The names must match the ones defined in other providers.
func (p *CPFFormatterProvider) DependsOn() []queryplanner.FieldName {
	return []queryplanner.FieldName{
		// Informs the library that this provider needs the CPF field to execute the query.
		// The "_" allows us to modify the same field we are providing, else it would be caught as a cyclic dependency by the library.
		"_CPF",
	}
}
