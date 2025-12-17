package queryplanner

import (
	"context"
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/arquivei/foundationkit/ref"
	"github.com/stretchr/testify/assert"
)

// nolint
func TestQueryPlan_New(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                               string
		providers                          []FieldProvider
		indexProvider                      IndexProvider
		request                            Request
		expectedFieldsToBeFetchedFromIndex []string
		expectedNewError                   string
	}{
		{
			name:      "[Success] empty - no Field providers",
			providers: []FieldProvider{},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name:  "a",
						Clear: func(d Document) {},
					},
					{
						Name:  "b",
						Clear: func(d Document) {},
					},
					{
						Name:  "c",
						Clear: func(d Document) {},
					},
				},
				data: &Payload{},
			},
			request:                            &requestMock{[]string{"a", "b", "c"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a", "b", "c"},
			expectedNewError:                   "",
		},
		{
			name: "[Success] With field providers",
			providers: []FieldProvider{
				&fieldProviderMock{
					name:      "b-provider",
					dependsOn: []FieldName{"a"},
					provides: []Field{
						{
							Name: "b",
							Fill: func(index int, executionContext ExecutionContext) error {
								return nil
							},
							Clear: func(d Document) {},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name:  "a",
						Clear: func(d Document) {},
					},
				},
				data: &Payload{},
			},
			request:                            &requestMock{[]string{"a"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a"},
			expectedNewError:                   "",
		},
		{
			name: "[Success] With field providers overwritting index-provided field",
			providers: []FieldProvider{
				&fieldProviderMock{
					name:      "a-provider",
					dependsOn: []FieldName{"_a"},
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								return nil
							},
							Clear: func(d Document) {},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name:  "a",
						Clear: func(d Document) {},
					},
				},
				data: &Payload{},
			},
			request:                            &requestMock{[]string{"a"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a"},
			expectedNewError:                   "",
		},
		{
			name: "[Error] Dependency cicle with index provided field. DependsOn field missing underline.",
			providers: []FieldProvider{
				&fieldProviderMock{
					name:      "a-provider",
					dependsOn: []FieldName{"a"},
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								return nil
							},
							Clear: func(d Document) {},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name:  "a",
						Clear: func(d Document) {},
					},
				},
				data: &Payload{},
			},
			request:                            &requestMock{},
			expectedFieldsToBeFetchedFromIndex: nil,
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkCycle: cycle found in field dependency [cycle= -> a -> a]",
		},
		{
			name: "[Error] FieldProvider should implement `Fill` methods",
			providers: []FieldProvider{
				&fieldProviderMock{
					name:      "a-provider",
					dependsOn: []FieldName{"b"},
					provides: []Field{
						{
							Name:  "a",
							Fill:  nil,
							Clear: func(_ Document) {},
						},
					},
				},
			},
			indexProvider:                      &indexProviderMock{},
			request:                            &requestMock{[]string{}},
			expectedFieldsToBeFetchedFromIndex: nil,
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkIfFieldProvidersAreDeclaredCorrectly: checkMethodsFromFieldProvider: there is no `fill' method for field [fieldName=a]",
		},
		{
			name: "[Error] FieldProvider should implement `Clear` methods",
			providers: []FieldProvider{
				&fieldProviderMock{
					name:      "a-provider",
					dependsOn: []FieldName{"b"},
					provides: []Field{
						{
							Name:  "a",
							Fill:  func(_ int, _ ExecutionContext) error { return nil },
							Clear: nil,
						},
					},
				},
			},
			indexProvider:                      &indexProviderMock{},
			request:                            &requestMock{[]string{}},
			expectedFieldsToBeFetchedFromIndex: nil,
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkIfFieldProvidersAreDeclaredCorrectly: checkMethodsFromFieldProvider: there is no `clear` method for field [fieldName=a]",
		},
		{
			name:      "[Error] IndexProvider should implement `Clear` methods",
			providers: []FieldProvider{},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name: "a",
					},
				},
			},
			request:                            &requestMock{[]string{}},
			expectedFieldsToBeFetchedFromIndex: nil,
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkIfIndexProviderIsDeclaredCorrectly: fieldprovider has no `clear` method defined [fieldName=a]",
		},
		{
			name:                               "[Error]  We should provide an IndexProvider - not initialized",
			providers:                          []FieldProvider{},
			indexProvider:                      nil,
			request:                            &requestMock{[]string{}},
			expectedFieldsToBeFetchedFromIndex: nil,
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkIfIndexProviderIsDeclaredCorrectly: indexProvider should not be nil",
		},
		{
			name:                               "[Error]  We should provide an IndexProvider - initialized",
			providers:                          []FieldProvider{},
			indexProvider:                      func() *indexProviderMock { return nil }(),
			request:                            &requestMock{[]string{}},
			expectedFieldsToBeFetchedFromIndex: nil,
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkIfIndexProviderIsDeclaredCorrectly: indexProvider should not be nil",
		},
		{
			name: "[Error] FieldProvider should not be nil - not initialized",
			providers: []FieldProvider{
				nil, nil,
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name:  "a",
						Clear: func(d Document) {},
					},
				},
				data: &Payload{},
			},
			request:                            &requestMock{[]string{"a"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a"},
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkIfFieldProvidersAreDeclaredCorrectly: fieldprovider should not be nil",
		},
		{
			name: "[Error] FieldProvider should not be nil - initialized",
			providers: []FieldProvider{
				func() *fieldProviderMock { return nil }(),
				func() *fieldProviderMock { return nil }(),
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name:  "a",
						Clear: func(d Document) {},
					},
				},
				data: &Payload{},
			},
			request:                            &requestMock{[]string{"a"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a"},
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkIfFieldProvidersAreDeclaredCorrectly: fieldprovider should not be nil",
		},
		{
			name: "[Error] Dependency cycle",
			providers: []FieldProvider{
				&fieldProviderMock{
					name:      "a-provider",
					dependsOn: []FieldName{"b"},
					provides: []Field{
						{
							Name:  "a",
							Fill:  func(i int, executionContext ExecutionContext) error { return nil },
							Clear: func(d Document) {},
						},
					},
				},
				&fieldProviderMock{
					name:      "b-provider",
					dependsOn: []FieldName{"c"},
					provides: []Field{
						{
							Name:  "b",
							Fill:  func(i int, executionContext ExecutionContext) error { return nil },
							Clear: func(d Document) {},
						},
					},
				},
				&fieldProviderMock{
					name:      "c-provider",
					dependsOn: []FieldName{"a"},
					provides: []Field{
						{
							Name:  "c",
							Fill:  func(i int, executionContext ExecutionContext) error { return nil },
							Clear: func(d Document) {},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{},
				data: &Payload{
					Documents: wrapDocuments([]*document{}),
				},
			},
			request:                            &requestMock{[]string{}},
			expectedFieldsToBeFetchedFromIndex: []string{},
			expectedNewError:                   "queryplanner.NewQueryPlanner: checkCycle: cycle found in field dependency [cycle= -> a -> b -> c -> a]",
		},
		{
			name: "[Error] Multiple FieldProviders for same Field",
			providers: []FieldProvider{
				&fieldProviderMock{
					name:      "a-provider",
					dependsOn: []FieldName{"b"},
					provides: []Field{
						{
							Name:  "a",
							Fill:  func(i int, executionContext ExecutionContext) error { return nil },
							Clear: func(d Document) {},
						},
					},
				},
				&fieldProviderMock{
					name:      "a-provider-2",
					dependsOn: []FieldName{"c"},
					provides: []Field{
						{
							Name:  "a",
							Fill:  func(i int, executionContext ExecutionContext) error { return nil },
							Clear: func(d Document) {},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{},
				data: &Payload{
					Documents: wrapDocuments([]*document{}),
				},
			},
			request:                            &requestMock{[]string{}},
			expectedFieldsToBeFetchedFromIndex: []string{},
			expectedNewError:                   "queryplanner.NewQueryPlanner: queryPlannerImpl.registerProviders: queryPlannerImpl.registerProvider: two providers for the same field [field=a]",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			planner, err := NewQueryPlanner(test.indexProvider, test.providers...)

			if test.expectedNewError != "" {
				assert.EqualError(t, err, test.expectedNewError)
				return
			}

			assert.NotNil(t, planner)
			assert.NoError(t, err)
		})
	}
}

//nolint:forcetypeassert
func TestQueryPlan_Execute_DependecyChain(t *testing.T) {
	t.Parallel()
	providerA := fieldProviderMock{
		name:      "a-provider",
		dependsOn: []FieldName{"b"},
		provides: []Field{
			{
				Name: "a",
				Fill: func(index int, executionContext ExecutionContext) error {
					doc := (executionContext.Payload.Documents[index]).(*document)
					doc.a = ref.Of("a")

					customPayload := (executionContext.Payload.CustomData).(*payload)
					customPayload.calledProviders = append(customPayload.calledProviders, "a-provider")
					return nil
				},
				Clear: func(d Document) {
					doc := (d).(*document)
					doc.a = nil
				},
			},
		},
	}

	providerB := fieldProviderMock{
		name:      "b-provider",
		dependsOn: []FieldName{"c"},
		provides: []Field{
			{
				Name: "b",
				Fill: func(index int, executionContext ExecutionContext) error {
					doc := (executionContext.Payload.Documents[index]).(*document)
					doc.b = ref.Of("b")

					customPayload := (executionContext.Payload.CustomData).(*payload)
					customPayload.calledProviders = append(customPayload.calledProviders, "b-provider")
					return nil
				},
				Clear: func(d Document) {
					doc := (d).(*document)
					doc.b = nil
				},
			},
		},
	}

	providerC := fieldProviderMock{
		name:      "c-provider",
		dependsOn: []FieldName{"d"},
		provides: []Field{
			{
				Name: "c",
				Fill: func(index int, executionContext ExecutionContext) error {
					doc := (executionContext.Payload.Documents[index]).(*document)
					doc.c = ref.Of("c")

					customPayload := (executionContext.Payload.CustomData).(*payload)
					customPayload.calledProviders = append(customPayload.calledProviders, "c-provider")
					return nil
				},
				Clear: func(d Document) {
					doc := (d).(*document)
					doc.c = nil
				},
			},
		},
	}

	indexProvider := &indexProviderMock{
		provides: []Index{
			{
				Name: "d",
				Clear: func(d Document) {
					doc := (d).(*document)
					doc.d = nil
				},
			},
			{
				Name: "e",
				Clear: func(d Document) {
					doc := (d).(*document)
					doc.e = nil
				},
			},
			{
				Name: "f",
				Clear: func(d Document) {
					doc := (d).(*document)
					doc.f = nil
				},
			},
		},
		data: &Payload{
			Documents: func() []Document {
				docs := []*document{
					{d: ref.Of("d1"), e: ref.Of("e1"), f: ref.Of("f1")},
					{d: ref.Of("d2"), e: ref.Of("e2"), f: ref.Of("f2")},
					{d: ref.Of("d3"), e: ref.Of("e3"), f: ref.Of("f3")},
				}
				return wrapDocuments(docs)
			}(),
			CustomData: &payload{calledProviders: []string{}},
		},
	}
	request := &requestMock{[]string{"a", "f"}}
	allProviders := wrapProviders([]*fieldProviderMock{&providerA, &providerB, &providerC})

	planner, err := NewQueryPlanner(indexProvider, allProviders...)
	assert.NoError(t, err)

	p := planner.NewPlan(request)

	myplan := p.(plan)
	expectedFieldsToBeFetchedFromIndex := []string{"d", "f"}
	assert.ElementsMatch(t, expectedFieldsToBeFetchedFromIndex, myplan.fieldsToBeFetchedFromIndex.ToStrings())

	data, err := p.Execute(context.Background())
	assert.NoError(t, err)
	payloadData := (data.CustomData).(*payload)

	expectedCallOrder := []string{
		"c-provider", "c-provider", "c-provider",
		"b-provider", "b-provider", "b-provider",
		"a-provider", "a-provider", "a-provider",
	}
	assert.Equal(t, expectedCallOrder, payloadData.calledProviders)

	documents := unwrapDocuments(data.Documents)
	expectedDocuments := []*document{
		{
			a: ref.Of("a"),
			f: ref.Of("f1"),
		},
		{
			a: ref.Of("a"),
			f: ref.Of("f2"),
		},
		{
			a: ref.Of("a"),
			f: ref.Of("f3"),
		},
	}
	assert.EqualValues(t, expectedDocuments, documents)
}

// nolint
func TestQueryPlan_Execute(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                               string
		providers                          []*fieldProviderMock
		indexProvider                      *indexProviderMock
		request                            *requestMock
		expectedFieldsToBeFetchedFromIndex []string
		expectedExecutionError             string
		expectedDocuments                  []*document
	}{
		{
			name: "[Success] Fields with underline",
			providers: []*fieldProviderMock{
				{
					name:      "a-provider",
					dependsOn: []FieldName{"_a"},
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.a = ref.Of("field_provider_modified_a")
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.a = nil
							},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name: "a",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.d = nil
						},
					},
				},
				data: &Payload{
					Documents: []Document{
						&document{
							a: ref.Of("index-provided-a"),
						},
					},
				},
			},
			request:                            &requestMock{[]string{"a"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a"},
			expectedExecutionError:             "",
			expectedDocuments: []*document{
				{
					a: ref.Of("field_provider_modified_a"),
				},
			},
		},
		{
			name: "[Success] Force index fields",
			providers: []*fieldProviderMock{
				{
					name:      "a-provider",
					dependsOn: []FieldName{"b", "_a"},
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.a = ref.Of("a___from_field_provider")
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.a = nil
							},
						},
					},
				},
				{
					name:      "c-provider",
					dependsOn: []FieldName{"a"},
					provides: []Field{
						{
							Name: "c",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.c = ref.Of("c___from_field_provider")
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.c = nil
							},
						},
					},
				},
				{
					name:      "e-provider",
					dependsOn: []FieldName{"a"},
					provides: []Field{
						{
							Name: "e",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.e = ref.Of("e___from_field_provider")
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.e = nil
							},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name: "a",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.a = nil
						},
					},
					{
						Name: "b",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.b = nil
						},
					},
					{
						Name: "c",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.c = nil
						},
					},
					{
						Name: "d",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.d = nil
						},
					},
					{
						Name: "e",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.e = nil
						},
					},
				},
				data: &Payload{
					Documents: func() []Document {
						docs := []*document{
							{
								a: ref.Of("a___from_index_provider"),
								b: ref.Of("b___from_index_provider"),
								c: ref.Of("c___from_index_provider"),
								d: ref.Of("d___from_index_provider"),
								e: ref.Of("e___from_index_provider"),
							},
						}
						return wrapDocuments(docs)
					}(),
					CustomData: nil,
				},
			},
			request:                            &requestMock{[]string{"a", "b", "c", "d", "_e"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a", "b", "d", "e"},
			expectedExecutionError:             "",
			expectedDocuments: []*document{
				{
					a: ref.Of("a___from_field_provider"),
					b: ref.Of("b___from_index_provider"),
					c: ref.Of("c___from_field_provider"),
					d: ref.Of("d___from_index_provider"),
					e: nil,
				},
			},
		},
		{
			name: "[Error] Execute should fail because there was a failure filling a document",
			providers: []*fieldProviderMock{
				{
					name:      "a-provider",
					dependsOn: []FieldName{"b"},
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.a = ref.Of("a")
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.a = nil
							},
						},
					},
				},
				{
					name:      "b-provider",
					dependsOn: []FieldName{"c"},
					provides: []Field{
						{
							Name: "b",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.b = ref.Of("b")
								return errors.E("problem filling b")
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.b = nil
							},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name: "c",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.c = nil
						},
					},
				},
				data: &Payload{
					Documents: func() []Document {
						docs := []*document{
							{
								a: ref.Of("a___from_index_provider"),
							},
						}
						return wrapDocuments(docs)
					}(),
					CustomData: nil,
				},
			},
			request:                            &requestMock{[]string{"a"}},
			expectedFieldsToBeFetchedFromIndex: []string{"c"},
			expectedExecutionError:             "queryplanner.Plan.Execute: planExecution.start: planExecution.executeProvider: problem filling b",
			expectedDocuments:                  []*document{},
		},

		{
			name: "[Error] Requested fields are not defined on IndexProvider",
			providers: []*fieldProviderMock{
				{
					name:      "a-provider",
					dependsOn: []FieldName{"b"},
					provides: []Field{
						{
							Name:  "a",
							Fill:  func(i int, executionContext ExecutionContext) error { return nil },
							Clear: func(d Document) {},
						},
					},
				},
				{
					name:      "b-provider",
					dependsOn: []FieldName{},
					provides: []Field{
						{
							Name:  "b",
							Fill:  func(i int, executionContext ExecutionContext) error { return nil },
							Clear: func(d Document) {},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{},
				data: &Payload{
					Documents: wrapDocuments([]*document{}),
				},
			},
			request:                            &requestMock{[]string{"a", "b", "c", "d"}},
			expectedFieldsToBeFetchedFromIndex: []string{"c", "d"},
			expectedExecutionError:             "queryplanner.Plan.Execute: checkIfIndexHasTheNecessaryFields: unsupported fields by index [fields=c,d]",
		},
		{
			name: "[Success] Document should have only asked fields",
			providers: []*fieldProviderMock{
				{
					name:      "a-provider",
					dependsOn: []FieldName{"b"},
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.a = ref.Of("a")
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.a = nil
							},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name: "a",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.a = nil
						},
					},
					{
						Name: "b",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.b = nil
						},
					},
					{
						Name: "c",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.c = nil
						},
					},
					{
						Name: "d",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.d = nil
						},
					},
					{
						Name: "e",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.e = nil
						},
					},
					{
						Name: "f",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.f = nil
						},
					},
					{
						Name: "g",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.g = nil
						},
					},
					{
						Name: "h",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.h = nil
						},
					},
					{
						Name: "i",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.i = nil
						},
					},
				},
				data: &Payload{
					Documents: func() []Document {
						docs := []*document{
							{
								a: ref.Of("a"),
								b: ref.Of("b"),
								c: ref.Of("c"),
								d: ref.Of("d"),
								e: ref.Of("e"),
								f: ref.Of("f"),
								g: ref.Of("g"),
								h: ref.Of("h"),
								i: ref.Of("i"),
							},
						}
						return wrapDocuments(docs)
					}(),
					CustomData: nil,
				},
			},
			request:                            &requestMock{[]string{"a", "d", "i"}},
			expectedFieldsToBeFetchedFromIndex: []string{"b", "d", "i"},
			expectedExecutionError:             "",
			expectedDocuments: []*document{
				{
					a: ref.Of("a"),
					d: ref.Of("d"),
					i: ref.Of("i"),
				},
			},
		},
		{
			name: "[Success] One provider use the same cache for all fields",
			providers: []*fieldProviderMock{
				{
					name: "provider",
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								docFromCache, _ := executionContext.Cache().GetOrLoad("ab", func() (interface{}, error) {
									return document{
										a: ref.Of("a"),
										b: ref.Of("b"),
									}, nil
								})
								doc.a = docFromCache.(document).a
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.a = nil
							},
						},
						{
							Name: "b",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								docFromCache, _ := executionContext.Cache().GetOrLoad("ab", func() (interface{}, error) {
									return document{
										b: ref.Of("should use cached value from previous field"),
									}, nil
								})
								doc.b = docFromCache.(document).b
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.b = nil
							},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{},
				data: &Payload{
					Documents:  []Document{&document{}},
					CustomData: nil,
				},
			},
			request:                            &requestMock{[]string{"a", "b"}},
			expectedFieldsToBeFetchedFromIndex: []string{},
			expectedExecutionError:             "",
			expectedDocuments: []*document{
				{
					a: ref.Of("a"),
					b: ref.Of("b"),
				},
			},
		},
		{
			name: "[Error] Index Provider returns an error in Execution",
			providers: []*fieldProviderMock{
				{
					name:      "a-provider",
					dependsOn: []FieldName{"_a"},
					provides: []Field{
						{
							Name: "a",
							Fill: func(index int, executionContext ExecutionContext) error {
								doc := (executionContext.Payload.Documents[index]).(*document)
								doc.a = ref.Of("a")
								return nil
							},
							Clear: func(d Document) {
								doc := (d).(*document)
								doc.a = nil
							},
						},
					},
				},
			},
			indexProvider: &indexProviderMock{
				provides: []Index{
					{
						Name: "a",
						Clear: func(d Document) {
							doc := (d).(*document)
							doc.d = nil
						},
					},
				},
				data: &Payload{},
				execute: func() func(indexProvider *indexProviderMock, ctx context.Context, request Request, fields []string) (*Payload, error) {
					return func(indexProvider *indexProviderMock, ctx context.Context, request Request, fields []string) (*Payload, error) {
						return nil, errors.New("index provider error")
					}
				}(),
			},
			request:                            &requestMock{[]string{"a"}},
			expectedFieldsToBeFetchedFromIndex: []string{"a"},
			expectedExecutionError:             "queryplanner.Plan.Execute: index provider error",
			expectedDocuments:                  []*document{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			planner, err := NewQueryPlanner(test.indexProvider, wrapProviders(test.providers)...)
			assert.NoError(t, err)

			p := planner.NewPlan(test.request)

			myplan := p.(plan)
			assert.ElementsMatch(t, test.expectedFieldsToBeFetchedFromIndex, myplan.fieldsToBeFetchedFromIndex.ToStrings())

			data, err := p.Execute(context.Background())
			if test.expectedExecutionError != "" {
				assert.EqualError(t, err, test.expectedExecutionError)
			} else {
				assert.NoError(t, err)
				documents := unwrapDocuments(data.Documents)
				assert.Equal(t, test.expectedDocuments, documents)
				assert.NoError(t, err)
			}
		})
	}
}

type payload struct {
	calledProviders []string
}

type document struct {
	a *string
	b *string
	c *string
	d *string
	e *string
	f *string
	g *string
	h *string
	i *string
}

func (d *document) GetDocument() interface{} {
	return d
}

// nolint
func unwrapDocuments(documents []Document) []*document {
	docs := make([]*document, 0, len(documents))
	for _, doc := range documents {
		docs = append(docs, (doc).(*document))
	}
	return docs
}

func wrapDocuments(documents []*document) []Document {
	docs := make([]Document, 0, len(documents))
	for _, doc := range documents {
		docs = append(docs, doc)
	}
	return docs
}

func wrapProviders(mocks []*fieldProviderMock) []FieldProvider {
	providers := make([]FieldProvider, 0, len(mocks))
	for _, p := range mocks {
		providers = append(providers, p)
	}

	return providers
}

type indexProviderMock struct {
	provides []Index
	data     *Payload
	execute  func(indexProvider *indexProviderMock, ctx context.Context, request Request, fields []string) (*Payload, error)
}

func (i *indexProviderMock) Provides() []Index {
	return i.provides
}

func (i *indexProviderMock) Execute(ctx context.Context, request Request, fields []string) (*Payload, error) {
	if i.execute == nil {
		return i.data, nil
	}
	return i.execute(i, ctx, request, fields)
}

type fieldProviderMock struct {
	name      string
	dependsOn []FieldName
	provides  []Field
}

func (m *fieldProviderMock) DependsOn() []FieldName {
	return m.dependsOn
}

func (m *fieldProviderMock) Provides() []Field {
	return m.provides
}

type requestMock struct{ Fields []string }

func (r *requestMock) GetRequestedFields() []string {
	return r.Fields
}

func (r *requestMock) GetRequest() interface{} {
	return nil
}
