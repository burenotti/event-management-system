package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/jmoiron/sqlx"
	"github.com/leporo/sqlf"
)

const (
	EventsOrgIdFkeyName     = "events_organization_id_fkey"
	EventsCreatorIdFkeyName = "events_creator_id_fkey"
	EventsPkeyName          = "events_pkey"
)

var (
	ErrEventNotFount = errors.New("event does not exist")
)

var EventUpdatesValidator = NewUpdatesValidator([]string{
	"name", "description", "begins_at", "ends_at",
	"registration_needed", "registration_begin", "registration_end",
})

type EventRepository struct {
	db DatabaseWrapper
}

func (r *EventRepository) Create(ctx context.Context, create *model.EventCreate) (*model.Event, error) {
	e := &model.Event{}
	err := sqlf.InsertInto("events").
		Set("organization_id", create.OrganizationID).
		Set("creator_id", create.CreatorID).
		Set("name", create.Name).
		Set("description", create.Description).
		Set("registration_needed", create.RegistrationNeeded).
		Set("registration_begin", create.RegistrationBegin).
		Set("registration_end", create.RegistrationEnd).
		Set("begins_at", create.BeginsAt).
		Set("ends_at", create.EndsAt).
		Returning("event_id").To(&e.EventID).
		Returning("organization_id, creator_id, name, description").
		To(&e.OrganizationID, &e.CreatorID, &e.Name, &e.Description).
		Returning("registration_needed, registration_begin, registration_end").
		To(&e.RegistrationNeeded, &e.RegistrationBegin, &e.RegistrationEnd).
		Returning("begins_at, ends_at").To(&e.BeginsAt, &e.EndsAt).
		QueryRowAndClose(ctx, r.db)

	if getViolatedConstraint(err) == EventsOrgIdFkeyName {
		return nil, ErrOrganizationNotFound
	} else if getViolatedConstraint(err) == EventsCreatorIdFkeyName {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EventRepository) GetById(ctx context.Context, eventId int64) (*model.Event, error) {
	e := &model.Event{}

	err := sqlf.From("events").
		Where("event_id = ? ", eventId).
		Select("event_id").To(&e.EventID).
		Select("organization_id, creator_id, name, description").
		To(&e.OrganizationID, &e.CreatorID, &e.Name, &e.Description).
		Select("registration_needed, registration_begin, registration_end").
		To(&e.RegistrationNeeded, &e.RegistrationBegin, &e.RegistrationEnd).
		Select("begins_at, ends_at").To(&e.BeginsAt, &e.EndsAt).
		QueryRowAndClose(ctx, r.db)

	if getViolatedConstraint(err) == EventsPkeyName {
		return nil, ErrEventNotFount
	} else if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *EventRepository) SelectBy(ctx context.Context, filter EventFilter) ([]model.Event, error) {
	where, args := filter.whereClause()
	limit := filter.limit()
	orderBy := filter.orderByClause()

	builder := sqlf.From("events").
		Select("event_id, organization_id, creator_id, name, description").
		Select("registration_needed, registration_begin, registration_end").
		Select("begins_at, ends_at").
		Where(where, args...)

	if orderBy != "" {
		builder = builder.OrderBy(orderBy)
	}

	if limit >= 0 {
		builder = builder.Limit(limit)
	}
	var events []model.Event
	var scanErr error
	err := builder.QueryAndClose(ctx, r.db, func(rows *sql.Rows) {
		e := model.Event{}
		scanErr = sqlx.StructScan(rows, &e)
		events = append(events, e)
	})
	if err != nil {
		return nil, err
	}
	if errors.Is(scanErr, sql.ErrNoRows) {
		return nil, nil
	} else {
		return events, scanErr
	}
}

func (r *EventRepository) UpdateEvent(ctx context.Context, eventId int64, updates map[string]interface{}) (*model.Event, error) {
	if err := EventUpdatesValidator.Validate(updates); err != nil {
		return nil, err
	}
	e := model.Event{}
	builder := sqlf.Update("events").
		Where("event_id = ?", eventId).
		Returning("event_id").To(&e.EventID).
		Returning("organization_id, creator_id, name, description").
		To(&e.OrganizationID, &e.CreatorID, &e.Name, &e.Description).
		Returning("registration_needed, registration_begin, registration_end").
		To(&e.RegistrationNeeded, &e.RegistrationBegin, &e.RegistrationEnd).
		Returning("begins_at, ends_at").To(&e.BeginsAt, &e.EndsAt)

	for field, val := range updates {
		builder = builder.Set(field, val)
	}
	err := builder.QueryRow(ctx, r.db)
	if getViolatedConstraint(err) == EventsPkeyName {
		return nil, ErrEventNotFount
	} else if getViolatedConstraint(err) == EventsOrgIdFkeyName {
		return nil, ErrOrganizationNotFound
	} else if getViolatedConstraint(err) == EventsCreatorIdFkeyName {
		return nil, ErrUserNotFound
	}
	return &e, err
}

func (r *EventRepository) DeleteEvent(ctx context.Context, eventId int64) error {
	res, err := sqlf.DeleteFrom("events").
		Where("event_id = ?", eventId).
		ExecAndClose(ctx, r.db)

	if err != nil {
		return err
	}

	// RowsAffected returns error only if the method is not supported by the driver.
	// So, since the other such getViolatedConstraint are depends on pgx driver,
	// using another driver is not possible.
	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrEventNotFount
	}
	return nil
}
