package queryplanner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fieldProviderByName_AddAndGetByName(t *testing.T) {
	t.Parallel()
	fp := newFieldProviderByName()
	provider1 := &fieldProviderMock{}
	fp.Add("Field1", provider1)
	fp.Add("Field2", provider1)

	for _, fieldName := range []string{"Field1", "Field2"} {
		provider, ok := fp.GetByName(FieldName(fieldName))
		assert.True(t, ok)
		assert.Equal(t, provider, provider1)
	}

	provider2 := &fieldProviderMock{}
	fp.Add("Field2", provider2)
	fp.Add("Field3", provider2)
	for _, fieldName := range []string{"Field2", "Field3"} {
		provider, ok := fp.GetByName(FieldName(fieldName))
		assert.True(t, ok)
		assert.Equal(t, provider, provider2)
	}

	provider, ok := fp.GetByName("Bla")
	assert.False(t, ok)
	assert.Equal(t, nil, provider)
}

func Test_fieldProviderByName_Length(t *testing.T) {
	t.Parallel()
	fp := newFieldProviderByName()
	assert.Equal(t, 0, fp.Length())
	for i := 1; i < 10; i++ {
		fp.Add(FieldName(fmt.Sprintf("Field_%d", i)), &fieldProviderMock{})
		assert.Equal(t, i, fp.Length())
	}
}

func Test_fieldProviderByName_GetFieldNames(t *testing.T) {
	t.Parallel()
	fp := newFieldProviderByName()
	for i := 1; i < 5; i++ {
		fp.Add(FieldName(fmt.Sprintf("Field_%d", i)), &fieldProviderMock{})
	}

	expected := []string{"Field_1", "Field_2", "Field_3", "Field_4"}
	assert.ElementsMatch(t, expected, fp.GetFieldNames())
}
