package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/leporo/sqlf"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
)

type OrganizationRepository struct {
	db DatabaseWrapper
}

func NewOrganizationRepository(db DatabaseWrapper) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, o *model.OrganizationCreate) (*model.Organization, error) {
	res := &model.Organization{
		Name:         o.Name,
		Address:      o.Address,
		ContactEmail: o.ContactEmail,
		ContactPhone: o.ContactPhone,
	}

	err := sqlf.InsertInto("organizations").
		Set("name", o.Name).
		Set("address", o.Address).
		Set("contact_email", o.ContactEmail).
		Set("contact_phone", o.ContactPhone).
		Returning("organization_id").To(&res.OrganizationID).
		QueryRow(ctx, r.db)

	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *OrganizationRepository) GetById(ctx context.Context, orgId int64) (*model.Organization, error) {
	res := &model.Organization{}

	err := sqlf.From("organizations").
		Select("organization_id").To(&res.OrganizationID).
		Select("name").To(&res.Name).
		Select("address").To(&res.Address).
		Select("contact_phone").To(&res.ContactPhone).
		Select("contact_email").To(&res.ContactEmail).
		Where("organization_id = ?", orgId).
		QueryRow(ctx, r.db)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%v: organization with provided id does not exist", ErrOrganizationNotFound)
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, orgId int64, updates map[string]interface{}) (*model.Organization, error) {
	fields := map[string]struct{}{
		"address": {}, "contact_email": {}, "contact_phone": {}, "name": {},
	}
	o := &model.Organization{}
	query := sqlf.Update("organizations").
		Where("organization_id = ?", orgId).
		Returning("organization_id").To(&o.OrganizationID).
		Returning("name").To(&o.Name).
		Returning("address").To(&o.Address).
		Returning("contact_phone").To(&o.ContactPhone).
		Returning("contact_email").To(&o.ContactEmail)

	for key, upd := range updates {
		if _, ok := fields[key]; !ok {
			return nil, fmt.Errorf("%v: could not update field '%r'", ErrLogicError, upd)
		}
		query = query.Set(key, upd)
	}

	err := query.QueryRow(ctx, r.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%v: organization with provided id does not exist", ErrOrganizationNotFound)
	} else if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *OrganizationRepository) Delete(ctx context.Context, orgId int64) error {
	res, err := sqlf.DeleteFrom("organizations").
		Where("organization_id = ?", orgId).
		Exec(ctx, r.db)

	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("%v: organization with provided id does not exist", ErrOrganizationNotFound)
	}
	return nil
}
