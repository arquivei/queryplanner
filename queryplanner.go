package queryplanner

import (
	"reflect"

	"github.com/arquivei/foundationkit/errors"
)

// QueryPlanner is an interface that creates a Plan
type QueryPlanner interface {
	NewPlan(Request) Plan
}

type queryPlanner struct {
	fieldToProviderMap fieldProviderByName
	indexProvider      IndexProvider
}

// NewQueryPlanner returns a new query planner unsing @providers
func NewQueryPlanner(indexProvider IndexProvider, providers ...FieldProvider) (QueryPlanner, error) {
	const op = errors.Op("queryplanner.NewQueryPlanner")

	err := checkIfIndexProviderIsDeclaredCorrectly(indexProvider)
	if err != nil {
		return nil, errors.E(op, err)
	}

	err = checkIfFieldProvidersAreDeclaredCorrectly(providers)
	if err != nil {
		return nil, errors.E(op, err)
	}

	planner := &queryPlanner{
		fieldToProviderMap: newFieldProviderByName(),
		indexProvider:      indexProvider,
	}

	err = planner.registerProviders(providers...)
	if err != nil {
		return nil, errors.E(op, err)
	}

	err = checkCycle(planner.fieldToProviderMap)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return planner, nil
}

func (q *queryPlanner) NewPlan(request Request) Plan {
	p := plan{
		fieldsToBeFetchedFromIndex: newFieldNameSet(0),
		providers:                  make([]FieldProvider, 0),
		indexProvider:              q.indexProvider,
		request:                    request,

		processedFields:    newFieldNameSet(0),
		processedProviders: newFieldProviderSet(0),
	}

	for _, field := range request.GetRequestedFields() {
		p.activateField(FieldName(field), q.fieldToProviderMap)
	}

	return p
}

func (q *queryPlanner) registerProviders(providers ...FieldProvider) error {
	const op = errors.Op("queryPlannerImpl.registerProviders")
	for _, provider := range providers {
		err := q.registerProvider(provider)
		if err != nil {
			return errors.E(op, err)
		}
	}
	return nil
}

func (q *queryPlanner) registerProvider(provider FieldProvider) error {
	const op = errors.Op("queryPlannerImpl.registerProvider")
	for _, field := range provider.Provides() {
		if _, foundProvider := q.fieldToProviderMap.GetByName(field.Name); foundProvider {
			return errors.E(
				op,
				"two providers for the same field",
				errors.KV("field", field.Name),
			)
		}
		q.fieldToProviderMap.Add(field.Name, provider)
	}
	return nil
}

func checkIfIndexProviderIsDeclaredCorrectly(indexProvider IndexProvider) error {
	const op = errors.Op("checkIfIndexProviderIsDeclaredCorrectly")
	if reflect.ValueOf(indexProvider).IsNil() {
		return errors.E(op, "indexProvider should not be nil")
	}
	for _, field := range indexProvider.Provides() {
		if field.Clear == nil {
			return errors.E(op, "fieldprovider has no `clear` method defined", errors.KV("fieldName", field.Name))
		}
	}
	return nil
}

func checkIfFieldProvidersAreDeclaredCorrectly(providers []FieldProvider) error {
	const op = errors.Op("checkIfFieldProvidersAreDeclaredCorrectly")

	for _, fieldProvider := range providers {
		if reflect.ValueOf(fieldProvider).IsNil() {
			return errors.E(op, "fieldprovider should not be nil")
		}
		err := checkMethodsFromFieldProvider(fieldProvider)
		if err != nil {
			return errors.E(op, err)
		}
	}
	return nil
}
