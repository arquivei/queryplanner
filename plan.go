package queryplanner

import (
	"context"
	"sort"
	"strings"

	"github.com/arquivei/foundationkit/errors"
)

// Plan is the product of the QueryPlanner. It can be executed, returning a
// Payload with the enriched Document and CustomData
type Plan interface {
	Execute(context.Context) (*Payload, error)
}

// plan implements a Plan. It can be executed using an IndexProvider and a
// set of FieldProviders to enrich the Payload
type plan struct {
	fieldsToBeFetchedFromIndex fieldNameSet
	indexProvider              IndexProvider
	providers                  []FieldProvider
	request                    Request

	processedFields    fieldNameSet
	processedProviders fieldProviderSet
}

// Execute runs a Plan and returns the enriched Payload
func (p plan) Execute(ctx context.Context) (*Payload, error) {
	const op = errors.Op("queryplanner.Plan.Execute")

	err := p.checkIfIndexHasTheNecessaryFields()
	if err != nil {
		return nil, errors.E(op, err)
	}

	fieldToBeFetchedFromIndex := p.fieldsToBeFetchedFromIndex.ToStrings()
	sort.Strings(fieldToBeFetchedFromIndex)

	data, err := p.indexProvider.Execute(ctx, p.request, fieldToBeFetchedFromIndex)
	if err != nil {
		return nil, errors.E(op, err)
	}

	filledFields := p.getRequestedFieldsFromIndex()

	execution := planExecution{
		plan:         &p,
		data:         data,
		filledFields: filledFields,
	}

	err = execution.start(ctx)
	if err != nil {
		return nil, errors.E(op, err)
	}

	return execution.data, nil
}

func (p plan) getRequestedFieldsFromIndex() fieldNameSet {
	filledFields := newFieldNameSet(p.fieldsToBeFetchedFromIndex.Length())
	for _, field := range p.indexProvider.Provides() {
		if p.fieldsToBeFetchedFromIndex.Exists(field.Name) {
			filledFields.Add(field.Name)
		}
	}
	return filledFields
}

func (p plan) checkIfIndexHasTheNecessaryFields() error {
	const op = errors.Op("checkIfIndexHasTheNecessaryFields")

	fieldsFromIndexProvider := p.indexProvider.Provides()
	fieldNamesFromIndexProvider := newFieldNameSet(len(fieldsFromIndexProvider))

	for _, field := range fieldsFromIndexProvider {
		fieldNamesFromIndexProvider.Add(field.Name)
	}

	fieldsNotDefinedInIndexProvider := p.fieldsToBeFetchedFromIndex.Diff(&fieldNamesFromIndexProvider)
	if fieldsNotDefinedInIndexProvider.Length() > 0 {
		fields := fieldsNotDefinedInIndexProvider.ToStrings()
		sort.Strings(fields)
		return errors.E(op, "unsupported fields by index", errors.KV("fields", strings.Join(fields, ",")))
	}
	return nil
}

func (p *plan) activateField(
	fieldName FieldName,
	fieldToProviderMap fieldProviderByName,
) {
	if p.processedFields.Exists(fieldName) {
		return
	}
	p.processedFields.Add(fieldName)

	provider, providerExists := fieldToProviderMap.GetByName(fieldName)
	if providerExists {
		p.activateProvider(provider, fieldToProviderMap)
	} else {
		p.fieldsToBeFetchedFromIndex.Add(transformIntoIndexField(fieldName))
	}
}

func (p *plan) activateProvider(
	provider FieldProvider,
	fieldToProviderMap fieldProviderByName,
) {
	if p.processedProviders.Exists(provider) {
		return
	}

	for _, field := range provider.DependsOn() {
		p.activateField(field, fieldToProviderMap)
	}

	p.processedProviders.Add(provider)
	p.providers = append(p.providers, provider)
}

func transformIntoIndexField(field FieldName) FieldName {
	isIndexField := field[0] == '_'
	if isIndexField {
		// Notation: fields starting with '_' are always provided by the index
		return field[1:]
	}
	return field
}
