package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGen(t *testing.T) {
	res := gen(46)
	assert.Equal(t, "1A", res)

	res = gen(35)
	assert.Equal(t, "Z", res)

	res = gen(370)
	assert.Equal(t, "AA", res)
}
