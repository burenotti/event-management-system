package repositories

import (
	"context"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OrganizationRepositoryTestSuite struct {
	DBTestSuite
}

func AssertEqualPtrValues[T comparable](t *testing.T, expected, actual *T, msgAndArgs ...interface{}) {
	assert.True(t, (expected == nil && actual == nil) || (*actual == *expected), msgAndArgs...)
}

func (s *OrganizationRepositoryTestSuite) TearDownTest(t *testing.T) {
	_, err := s.db.Exec("TRUNCATE users, organzizations")
	require.NoError(s.T(), err, "should teardown test normally")
}

func TestOrganizationRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &OrganizationRepositoryTestSuite{
		*DBTestSuiteFromEnv(),
	})
}

func (s *OrganizationRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	r := NewOrganizationRepository(NewDatabase(s.db))
	addr := "Moscow"
	email := "org@example.com"
	phone := "71111111111"
	name := "My organization"
	expectedOrg := model.OrganizationCreate{
		OrganizationID: 0,
		Name:           name,
		Address:        &addr,
		ContactEmail:   &email,
		ContactPhone:   &phone,
	}
	org, err := r.Create(ctx, &expectedOrg)
	assert.NoError(s.T(), err, "should create user without errors")
	assert.Equal(s.T(), name, org.Name)
	assert.Equal(s.T(), &addr, org.Address)
	assert.Equal(s.T(), &phone, org.ContactPhone)
	assert.Equal(s.T(), &email, org.ContactEmail)

}

func (s *OrganizationRepositoryTestSuite) TestGetById() {
	ctx := context.Background()
	r := NewOrganizationRepository(NewDatabase(s.db))
	org, err := r.GetById(ctx, 9999)
	assert.Nil(s.T(), org, "should return nil if organization does not exist")
	assert.ErrorAsf(s.T(), err, &ErrOrganizationNotFound, "should return wrapped ErrOgranizationNotFound")

	addr := "Moscow"
	email := "org@example.com"
	phone := "71111111111"
	name := "My organization"
	orgId := -1
	q := `
		INSERT INTO organizations
		    (name, address, contact_email, contact_phone)
		VALUES 
			($1, $2, $3, $4)
		RETURNING organization_id
	`
	err = s.db.GetContext(ctx, &orgId, q, name, addr, email, phone)
	require.NoError(s.T(), err)

	org, err = r.GetById(ctx, int64(orgId))
	assert.NoError(s.T(), err, "should normally get organization by id")
	assert.Equal(s.T(), name, org.Name)
	assert.Equal(s.T(), &addr, org.Address)
	assert.Equal(s.T(), &email, org.ContactEmail)
	assert.Equal(s.T(), &phone, org.ContactPhone)
}

func (s *OrganizationRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	r := NewOrganizationRepository(NewDatabase(s.db))
	org, err := r.Update(ctx, 1, map[string]interface{}{
		"random_key": "random_value",
	})
	assert.Nil(s.T(), org, "should not return organization if error occurred")

	org, err = r.Update(ctx, -1, map[string]interface{}{
		"address": "Updated my organization",
	})

	assert.Nil(s.T(), org, "should not return organization if error occurred")
	assert.ErrorAs(s.T(), err, &ErrLogicError, "should return error if object does not exist")

	addr := "Moscow"
	email := "org@example.com"
	phone := "71111111111"
	name := "My organization"
	orgId := -1
	q := `
		INSERT INTO organizations
		    (name, address, contact_email, contact_phone)
		VALUES 
			($1, $2, $3, $4)
		RETURNING organization_id
	`
	err = s.db.GetContext(ctx, &orgId, q, name, addr, email, phone)
	require.NoError(s.T(), err)

	newAddr := "SPB"
	org, err = r.Update(ctx, int64(orgId), map[string]interface{}{
		"address": &newAddr,
	})
	assert.NoError(s.T(), err, "should normally update address")
	assert.Equal(s.T(), name, org.Name)
	AssertEqualPtrValues(s.T(), &newAddr, org.Address)
	assert.Equal(s.T(), &email, org.ContactEmail)
	assert.Equal(s.T(), &phone, org.ContactPhone)
}

func (s *OrganizationRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	r := NewOrganizationRepository(NewDatabase(s.db))
	err := r.Delete(ctx, -1)
	assert.ErrorAs(s.T(), err, &ErrUserNotFound, "should return error if there was a try to delete organization that does not exist")

	addr := "Moscow"
	email := "org@example.com"
	phone := "71111111111"
	name := "My organization"
	orgId := -1
	q := `
		INSERT INTO organizations
		    (name, address, contact_email, contact_phone)
		VALUES 
			($1, $2, $3, $4)
		RETURNING organization_id
	`
	err = s.db.GetContext(ctx, &orgId, q, name, addr, email, phone)
	require.NoError(s.T(), err)

	err = r.Delete(ctx, int64(orgId))
	assert.NoError(s.T(), err, "should normally delete organization")
}
