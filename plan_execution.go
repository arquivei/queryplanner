package queryplanner

import (
	"context"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/trace"
)

type planExecution struct {
	plan         *plan
	data         *Payload
	filledFields fieldNameSet
}

func (e *planExecution) start(ctx context.Context) error {
	const op = errors.Op("planExecution.start")

	ctx, span := trace.StartSpan(ctx, op.String())
	defer span.End(nil)

	for _, provider := range e.plan.providers {
		err := e.executeProvider(ctx, provider)
		if err != nil {
			return errors.E(op, err)
		}
	}

	e.clearNonRequestedFields()
	return nil
}

func (e *planExecution) clearNonRequestedFields() {
	requestedFields := e.plan.request.GetRequestedFields()

	for _, document := range e.data.Documents {
		e.clearNonRequestedFieldsFromDocument(document, requestedFields)
	}
}

func (e *planExecution) executeProvider(ctx context.Context, provider FieldProvider) error {
	const op = errors.Op("planExecution.executeProvider")

	executionContext := ExecutionContext{
		Context: ctx,
		Request: e.plan.request,
		Payload: e.data,
		cache:   newCache(),
	}

	for _, field := range provider.Provides() {
		if e.filledFields.Exists(field.Name) {
			continue
		}

		for index := range e.data.Documents {
			err := field.Fill(index, executionContext)
			if err != nil {
				return errors.E(op, err)
			}
			e.filledFields.Add(field.Name)
		}
	}
	return nil
}

func (e *planExecution) clearNonRequestedFieldsFromDocument(document Document, requestedFields []string) {
	for _, field := range e.plan.indexProvider.Provides() {
		isRequestedField := isInArray(string(field.Name), requestedFields)
		if !isRequestedField {
			field.Clear(document)
		}
	}

	for _, provider := range e.plan.providers {
		for _, field := range provider.Provides() {
			isRequestedField := isInArray(string(field.Name), requestedFields)
			if !isRequestedField {
				field.Clear(document)
			}
		}
	}
}
