package queryplanner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fieldNameSet_AddAndExists(t *testing.T) {
	fieldSet := newFieldNameSet(2)

	field1 := FieldName("Field1")
	fieldSet.Add(field1)

	field2 := FieldName("Field2")

	assert.True(t, fieldSet.Exists(field1))
	assert.False(t, fieldSet.Exists(field2))
}

func Test_fieldNameSet_Diff(t *testing.T) {
	A := newFieldNameSet(0)
	B := newFieldNameSet(0)

	A.Add("Field1")
	A.Add("Field2")
	A.Add("Field3")

	B.Add("Field3")
	B.Add("Field4")
	B.Add("Field5")

	AminusB := A.Diff(&B)
	assert.ElementsMatch(t, []string{"Field1", "Field2"}, AminusB.ToStrings())

	BminusA := B.Diff(&A)
	assert.ElementsMatch(t, []string{"Field4", "Field5"}, BminusA.ToStrings())

	BminusB := B.Diff(&B)
	assert.ElementsMatch(t, []string{}, BminusB.ToStrings())

	AminusA := A.Diff(&A)
	assert.ElementsMatch(t, []string{}, AminusA.ToStrings())
}

func Test_fieldNameSet_Length(t *testing.T) {
	A := newFieldNameSet(0)
	assert.Equal(t, 0, A.Length())
	for i := 0; i < 10; i++ {
		A.Add(FieldName(fmt.Sprintf("Field %d", i)))
		assert.Equal(t, i+1, A.Length())

	}
}

func Test_fieldProviderSet_AddAndExists(t *testing.T) {
	set := newFieldProviderSet(0)
	provider1 := &fieldProviderMock{}
	provider2 := &fieldProviderMock{}
	assert.False(t, set.Exists(provider1))

	set.Add(provider1)
	assert.True(t, set.Exists(provider1))
	assert.False(t, set.Exists(provider2))

	set.Add(provider2)
	assert.True(t, set.Exists(provider1))
	assert.True(t, set.Exists(provider2))
}
