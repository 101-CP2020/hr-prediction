package io

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_apiImpl_GetPeriods(t *testing.T) {
	api := NewAPI()

	result, err := api.GetPeriods(11196, 0)

	assert.NoError(t, err)
	assert.Equal(t, nil, result)
}
