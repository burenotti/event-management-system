package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/burenotti/rtu-it-lab-recruit/repositories"
)

var (
	ErrBusinessLogicViolation = errors.New("business logic violation: ")
)

type OrganizationStorage interface {
	Create(ctx context.Context, o *model.OrganizationCreate) (*model.Organization, error)
	GetById(ctx context.Context, orgId int64) (*model.Organization, error)
	Update(ctx context.Context, orgId int64, updates map[string]interface{}) (*model.Organization, error)
	Delete(ctx context.Context, orgId int64) error
	AddMember(ctx context.Context, orgId int64, mem *model.OrganizationMemberCreate) (*model.OrganizationMember, error)
	ListMembers(ctx context.Context, orgId int64) ([]model.OrganizationMember, error)
	SetMemberRights(ctx context.Context, orgId int64, userId int64, newRights model.MemberRights) (*model.OrganizationMember, error)
	GetMember(ctx context.Context, orgId int64, userId int64) (*model.OrganizationMember, error)
	DeleteMember(ctx context.Context, orgId int64, userId int64) error
}

type OrganizationUseCase struct {
	Transactioner       StorageTransactioner
	OrganizationStorage OrganizationStorage
}

func (c *OrganizationUseCase) CreateOrganization(ctx context.Context, userId int64, org *model.OrganizationCreate) (*model.Organization, error) {
	createdOrg := &model.Organization{}
	err := c.Transactioner.Atomic(ctx, func(ctx context.Context) (err error) {
		createdOrg, err = c.OrganizationStorage.Create(ctx, org)
		if err != nil {
			return err
		}

		_, err = c.OrganizationStorage.AddMember(ctx, createdOrg.OrganizationID, &model.OrganizationMemberCreate{
			UserID:  userId,
			IsOwner: true,
			Rights: model.MemberRights{
				EditEvents:    true,
				ManageMembers: true,
			},
		})
		return err
	})
	return createdOrg, err
}

func (c *OrganizationUseCase) GetOrganization(ctx context.Context, orgId int64) (*model.Organization, error) {
	return c.OrganizationStorage.GetById(ctx, orgId)
}

func (c *OrganizationUseCase) DeleteOrganization(ctx context.Context, user *model.AuthPayload, orgId int64) error {
	return c.Transactioner.Atomic(ctx, func(ctx context.Context) error {
		errLogic := fmt.Errorf("%w: only member can delete organization", ErrBusinessLogicViolation)

		mem, err := c.OrganizationStorage.GetMember(ctx, orgId, user.UserID)
		if errors.Is(err, repositories.ErrMemberNotFound) {
			return errLogic
		} else if err != nil {
			return err
		}
		if !mem.IsOwner {
			return errLogic
		}
		return c.OrganizationStorage.Delete(ctx, orgId)
	})
}
