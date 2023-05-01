package repositories

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEventSinceFilter(t *testing.T) {
	since := time.Now()
	f := NewEventSinceFilter(since)
	query, args := f.whereClause()
	assert.Equal(t, "(created_at >= ?)", query)
	assert.Equal(t, []interface{}{since}, args)
	assert.Equal(t, "", f.orderByClause())
}

func TestEventNameLikeFilter(t *testing.T) {
	f := NewNameLikeFilter("Test")
	query, args := f.whereClause()
	assert.Equal(t, "(name SIMILAR TO ?)", query)
	assert.Equal(t, []interface{}{"Test"}, args)
	assert.Equal(t, "", f.orderByClause())
}

func TestEventAndFilter(t *testing.T) {
	since := time.Now()
	and := NewEventAndFilter(
		NewEventSinceFilter(since),
		NewNameLikeFilter("Test"),
	)
	query, args := and.whereClause()
	assert.Equal(t, "((created_at >= ?) AND (name SIMILAR TO ?))", query)
	assert.Equal(t, []interface{}{since, "Test"}, args)
}

func TestEventOrFilter(t *testing.T) {
	since := time.Now()
	and := NewEventOrFilter(
		NewEventSinceFilter(since),
		NewNameLikeFilter("Test"),
	)
	query, args := and.whereClause()
	assert.Equal(t, "((created_at >= ?) OR (name SIMILAR TO ?))", query)
	assert.Equal(t, []interface{}{since, "Test"}, args)
}

func TestEventFilterComposition(t *testing.T) {
	since := time.Now()
	and := NewEventAndFilter(
		NewEventSinceFilter(since),
		NewEventOrFilter(
			NewNameLikeFilter("John"),
			NewNameLikeFilter("Doe"),
		),
	)
	query, args := and.whereClause()
	assert.Equal(t, "((created_at >= ?) AND ((name SIMILAR TO ?) OR (name SIMILAR TO ?)))", query)
	assert.Equal(t, []interface{}{since, "John", "Doe"}, args)
}

func TestLimitFilter(t *testing.T) {
	since := time.Now()
	withLimit := NewLimitFilter(5,
		NewEventAndFilter(
			NewEventSinceFilter(since),
			NewEventOrFilter(
				NewNameLikeFilter("John"),
				NewNameLikeFilter("Doe"),
			),
		),
	)
	query, args := withLimit.whereClause()
	assert.Equal(t, "((created_at >= ?) AND ((name SIMILAR TO ?) OR (name SIMILAR TO ?)))", query)
	assert.Equal(t, []interface{}{since, "John", "Doe"}, args)
	assert.Equal(t, int64(5), withLimit.limit_)
}
