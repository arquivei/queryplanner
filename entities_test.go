package queryplanner

import (
	"fmt"
	"testing"

	"github.com/arquivei/foundationkit/errors"
	"github.com/stretchr/testify/assert"
)

func Test_cache_ExecuteOrRetrieve_CachingError(t *testing.T) {
	ctx := ExecutionContext{}
	cache := ctx.Cache()
	calledTimes := 0

	for i := 0; i < 10; i++ {
		_, _ = cache.GetOrLoad("biscoito", func() (interface{}, error) {
			calledTimes++
			return nil, errors.E("err biscoito")
		})
	}

	result, err := cache.GetOrLoad("biscoito", func() (interface{}, error) { return nil, nil })

	assert.Equal(t, calledTimes, 1)
	assert.Nil(t, result)
	assert.EqualError(t, err, "err biscoito")
}

func Test_cache_ExecuteOrRetrieve_CachingResult(t *testing.T) {
	ctx := ExecutionContext{}
	cache := ctx.Cache()
	calledTimes := 0

	type customStructure struct {
		A string
		B string
	}

	for i := 0; i < 10; i++ {
		_, _ = cache.GetOrLoad("biscoito", func() (interface{}, error) {
			calledTimes++
			return customStructure{
				A: fmt.Sprintf("A_%d", calledTimes),
				B: fmt.Sprintf("B_%d", calledTimes),
			}, nil
		})
	}

	cached, err := cache.GetOrLoad("biscoito", nil)
	result := cached.(customStructure)

	assert.Nil(t, err)
	assert.Equal(t, "A_1", result.A)
	assert.Equal(t, "B_1", result.B)
	assert.Equal(t, calledTimes, 1)
}
