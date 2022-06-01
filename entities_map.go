package queryplanner

func newFieldProviderByName() fieldProviderByName {
	return fieldProviderByName{
		providersByName: make(map[FieldName]FieldProvider),
	}
}

type fieldProviderByName struct {
	providersByName map[FieldName]FieldProvider
}

func (f *fieldProviderByName) Length() int {
	return len(f.providersByName)
}

func (f *fieldProviderByName) GetFieldNames() []string {
	fieldNames := make([]string, 0, len(f.providersByName))
	for k := range f.providersByName {
		fieldNames = append(fieldNames, string(k))
	}

	return fieldNames
}

func (f *fieldProviderByName) GetByName(field FieldName) (FieldProvider, bool) {
	provider, ok := f.providersByName[field]
	return provider, ok
}

func (f *fieldProviderByName) Add(field FieldName, provider FieldProvider) {
	f.providersByName[field] = provider
}
