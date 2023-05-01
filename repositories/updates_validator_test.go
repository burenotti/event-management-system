package repositories

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUpdatesValidator(t *testing.T) {
	v := NewUpdatesValidator([]string{"event_id", "org_id"})
	expectedSet := map[string]struct{}{"event_id": {}, "org_id": {}}
	assert.Equal(t, expectedSet, v.allowedFields)
}

func TestUpdatesValidator_Validate(t *testing.T) {
	v := NewUpdatesValidator([]string{"event_id", "org_id", "name"})
	correctMap := UpdatesMap{
		"event_id": 5,
		"org_id":   7,
	}
	err := v.Validate(correctMap)
	assert.NoError(t, err, "should not return error on correct updates")
	incorrectMap := UpdatesMap{
		"event_id":      7,
		"unknown_field": "123",
	}
	err = v.Validate(incorrectMap)
	assert.ErrorIs(t, err, ErrUpdatesValidationError, "should return wrapped ErrUpdatesValidationError")
}
