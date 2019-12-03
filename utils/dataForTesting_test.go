package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadIndividuals(t *testing.T) {
	v, err := GetUsersForTesting()
	assert.NotEmpty(t,v,"data loaded")
	assert.NotNil(t,v[0],"data loaded")
	assert.Equal(t, err, nil, "OK individual loaded")

}