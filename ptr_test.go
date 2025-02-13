package Gadget

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToPtr(t *testing.T) {
	i := 12
	res := ToPtr[int](i)
	assert.Equal(t, &i, res)
}
