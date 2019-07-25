package goex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloatToString(t *testing.T) {
	assert.Equal(t, "1", FloatToString(1.10231000, 0))
	assert.Equal(t, "0.102", FloatToString(0.10231000, 3))
	assert.Equal(t, "1.10231", FloatToString(1.10231000, 8))
	assert.NotEqual(t, "1.10231000", FloatToString(1.10231000, 8))
}
