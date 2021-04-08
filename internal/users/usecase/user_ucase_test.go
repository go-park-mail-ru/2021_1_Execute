package usecase_test

import (
	"2021_1_Execute/src/api"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPasswordValide(t *testing.T) {
	assert.True(t, api.IsPasswordValid("password"))
	assert.False(t, api.IsPasswordValid("pass"))
}
