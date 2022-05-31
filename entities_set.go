package queryplanner

func newFieldNameSet(cap int) fieldNameSet {
	return fieldNameSet{
		data: make(map[FieldName]struct{}, cap),
	}
}

type fieldNameSet struct {
	data map[FieldName]struct{}
}

func (f *fieldNameSet) Add(fieldName FieldName) {
	f.data[fieldName] = struct{}{}
}

func (f *fieldNameSet) Exists(fieldName FieldName) bool {
	_, ok := f.data[fieldName]
	return ok
}

func (f *fieldNameSet) ToStrings() []string {
	fields := make([]string, 0, len(f.data))
	for field := range f.data {
		fields = append(fields, string(field))
	}
	return fields
}

func (f *fieldNameSet) Diff(set *fieldNameSet) *fieldNameSet {
	difference := newFieldNameSet(len(f.data))
	for field := range f.data {
		if !set.Exists(field) {
			difference.Add(field)
		}
	}
	return &difference
}

func (f *fieldNameSet) Length() int {
	return len(f.data)
}

func newFieldProviderSet(cap int) fieldProviderSet {
	return fieldProviderSet{
		data: make(map[FieldProvider]struct{}, cap),
	}
}

type fieldProviderSet struct {
	data map[FieldProvider]struct{}
}

func (f *fieldProviderSet) Add(fieldType FieldProvider) {
	f.data[fieldType] = struct{}{}
}

func (f *fieldProviderSet) Exists(fieldType FieldProvider) bool {
	_, ok := f.data[fieldType]
	return ok
}
