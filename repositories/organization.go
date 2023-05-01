package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/leporo/sqlf"
)

const (
	MembersOrgFkeyName  = "organization_members_organization_id_fkey"
	MembersUserFkeyName = "organization_members_user_id_fkey"
	MembersPkeyName     = "organization_members_pkey"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrMemberNotFound       = errors.New("organization member not found")
	ErrUserAlreadyMember    = errors.New("user is already member of organization")
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

func (r *OrganizationRepository) AddMember(ctx context.Context, orgId int64, mem *model.OrganizationMemberCreate) (*model.OrganizationMember, error) {
	_, err := sqlf.InsertInto("organization_members").
		Set("organization_id", orgId).
		Set("user_id", mem.UserID).
		Set("can_manage_members", mem.Can.ManageMembers).
		Set("can_edit_events", mem.Can.EditEvents).
		Set("can_view_events", mem.Can.ViewEvents).
		ExecAndClose(ctx, r.db)

	if getViolatedConstraint(err) == MembersPkeyName {
		return nil, ErrUserAlreadyMember
	} else if getViolatedConstraint(err) == MembersUserFkeyName {
		return nil, ErrUserNotFound
	} else if getViolatedConstraint(err) == MembersOrgFkeyName {
		return nil, ErrOrganizationNotFound
	} else if err != nil {
		return nil, err
	}

	return &model.OrganizationMember{
		UserID: mem.UserID,
		Can:    mem.Can,
	}, nil
}

func (r *OrganizationRepository) ListMembers(ctx context.Context, orgId int64) ([]model.OrganizationMember, error) {
	var members []model.OrganizationMember
	var scanErr error
	err := sqlf.From("organization_members").
		Where("organization_id = ?", orgId).
		Select("user_id, can_manage_members, can_edit_events, can_view_events").
		QueryAndClose(ctx, r.db, func(rows *sql.Rows) {
			m := model.OrganizationMember{}
			err := rows.Scan(&m.UserID, &m.Can.ManageMembers, &m.Can.EditEvents, &m.Can.ViewEvents)
			if errors.Is(err, sql.ErrNoRows) {
				scanErr = ErrOrganizationNotFound
			} else if err != nil {
				scanErr = err
				return
			}
			members = append(members, m)
		})

	if err != nil {
		return nil, err
	}
	if scanErr != nil {
		return nil, scanErr
	}
	return members, nil
}

func (r *OrganizationRepository) SetMemberRights(ctx context.Context, orgId, userId int64, newRights model.MemberPrivileges) (*model.OrganizationMember, error) {
	res, err := sqlf.Update("organization_members").
		Where("organizationId = ? AND user_id = ?", orgId, userId).
		Set("can_view_events", newRights.ViewEvents).
		Set("can_edit_events", newRights.EditEvents).
		Set("can_manage_members", newRights.ManageMembers).
		Exec(ctx, r.db)

	if err != nil {
		return nil, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if count != 1 {
		return nil, ErrMemberNotFound
	}
	return &model.OrganizationMember{
		UserID: userId,
		Can:    newRights,
	}, nil
}

func (r *OrganizationRepository) DeleteMember(ctx context.Context, orgId, userId int64) error {
	res, err := sqlf.DeleteFrom("organization_members").
		Where("organization_id = ? AND user_id = ?", orgId, userId).
		Exec(ctx, r.db)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrMemberNotFound
	}
	return nil
}
