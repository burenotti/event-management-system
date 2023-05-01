package model

import "time"

type EventCreate struct {
	Name               string
	OrganizationID     int64
	CreatorID          int64
	Description        string
	BeginsAt           time.Time
	EndsAt             time.Time
	RegistrationNeeded bool
	RegistrationBegin  *time.Time
	RegistrationEnd    *time.Time
}

type Event struct {
	EventID            int64      `db:"event_id"`
	Name               string     `db:"name"`
	OrganizationID     int64      `db:"organization_id"`
	CreatorID          int64      `db:"creator_id"`
	Description        string     `db:"description"`
	RegistrationBegin  *time.Time `db:"registration_begin"`
	RegistrationEnd    *time.Time `db:"registration_end"`
	BeginsAt           time.Time  `db:"begins_at"`
	EndsAt             time.Time  `db:"ends_at"`
	RegistrationNeeded bool       `db:"registration_needed"`
	CreatedAt          time.Time  `db:"created_at"`
	PublishedAt        *time.Time `db:"published_at"`
}

func (e *Event) IsPublished() bool {
	return e.PublishedAt != nil
}
