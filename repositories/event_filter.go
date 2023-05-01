package repositories

import (
	"fmt"
	"strings"
	"time"
)

type EventFilter interface {
	whereClause() (string, []interface{})
	orderByClause() string
	limit() int
}

type EventJoinerFilter struct {
	joiner  string
	filters []EventFilter
	limit_  int
}

func (f *EventJoinerFilter) whereClause() (string, []interface{}) {
	queryItems := make([]string, 0, len(f.filters))
	queryArgs := make([]interface{}, 0, len(f.filters))
	for _, filter := range f.filters {
		where, args := filter.whereClause()
		if where != "" {
			queryItems = append(queryItems, where)
		}
		queryArgs = append(queryArgs, args...)
	}
	query := strings.Join(queryItems, f.joiner)
	return fmt.Sprintf("(%s)", query), queryArgs
}

func (f *EventJoinerFilter) orderByClause() string {
	orderBy := make([]string, len(f.filters))
	for i, filter := range f.filters {
		orderBy[i] = filter.orderByClause()
	}
	return strings.Join(orderBy, ", ")
}

func (f *EventJoinerFilter) limit() int {
	return -1
}

type EventAndFilter struct {
	EventJoinerFilter
}

func NewEventAndFilter(filters ...EventFilter) *EventAndFilter {
	return &EventAndFilter{
		EventJoinerFilter{
			joiner:  " AND ",
			filters: filters,
		},
	}
}

type EventOrFilter struct {
	EventJoinerFilter
}

func NewEventOrFilter(filters ...EventFilter) *EventOrFilter {
	return &EventOrFilter{
		EventJoinerFilter{
			joiner:  " OR ",
			filters: filters,
		},
	}
}

type BaseWhereFilter struct {
	query string
	args  []interface{}
}

func (f *BaseWhereFilter) whereClause() (string, []interface{}) {
	return f.query, f.args
}

func (f *BaseWhereFilter) orderByClause() string {
	return ""
}

func (f *BaseWhereFilter) limit() int {
	return -1
}

type EventSinceFilter struct {
	BaseWhereFilter
}

func NewEventSinceFilter(since time.Time) *EventSinceFilter {
	return &EventSinceFilter{
		BaseWhereFilter{
			query: "(created_at >= ?)",
			args:  []interface{}{since},
		},
	}
}

type NameLikeFilter struct {
	BaseWhereFilter
}

func NewNameLikeFilter(pattern string) *NameLikeFilter {
	return &NameLikeFilter{
		BaseWhereFilter{
			query: "(name SIMILAR TO ?)",
			args:  []interface{}{pattern},
		},
	}
}

type LimitFilter struct {
	child  EventFilter
	limit_ int64
}

func NewLimitFilter(limit int64, child EventFilter) *LimitFilter {
	return &LimitFilter{child: child, limit_: limit}
}

func (f *LimitFilter) whereClause() (string, []interface{}) {
	return f.child.whereClause()
}

func (f *LimitFilter) orderByClause() string {
	return f.child.orderByClause()
}

func (f *LimitFilter) limit() int64 {
	return f.limit_
}
