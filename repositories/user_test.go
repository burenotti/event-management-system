package repositories

import (
	"context"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestNewUserRepository(t *testing.T) {
	db := NewDatabase(&sqlx.DB{})
	repo := NewUserRepository(db)
	assert.Equal(t, db, repo.db)
}

type UserRepositoryTestSuite struct {
	DBTestSuite
}

func (s *UserRepositoryTestSuite) TearDownTest(t *testing.T) {
	_, err := s.db.Exec("TRUNCATE users")
	require.NoError(s.T(), err, "should teardown test normally")
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &UserRepositoryTestSuite{*DBTestSuiteFromEnv()})
}

func (s *UserRepositoryTestSuite) TestCreate() {
	repo := NewUserRepository(NewDatabase(s.db))

	users := []model.UserCreate{
		{"John", "Doe", "Jr.", "johndoe@example.com"},
		{"Michel", "Smith", "", "smith@example.com"},
		{"John", "Doe", "Jr.", "johndoe@example.com"},
		{"'; DELETE * FROM users; --", "", "", ""},
	}

	for _, u := range users {
		name := fmt.Sprintf("Create user: %s %s", u.FirstName, u.LastName)
		s.T().Run(name, func(t *testing.T) {
			ctx := context.Background()
			user, err := repo.Create(ctx, &u)
			assert.NoError(t, err, "Should create user without error")
			assert.Equal(t, u.FirstName, user.FirstName)
			assert.Equal(t, u.LastName, user.LastName)
			assert.Equal(t, u.MiddleName, user.MiddleName)
			assert.Equal(t, u.Email, user.Email)
			assert.Equal(t, false, user.IsActive, "created user should be deactivated")
		})
	}
}

func (s *UserRepositoryTestSuite) TestGetById() {
	ctx := context.Background()
	repo := NewUserRepository(NewDatabase(s.db))

	users := []model.UserCreate{
		{"John", "Doe", "Jr.", "johndoe@example.com"},
		{"Michel", "Smith", "", "smith@example.com"},
		{"John", "Doe", "Jr.", "johndoe@example.com"},
		{"'; DELETE * FROM users; --", "", "", ""},
	}

	q := `INSERT INTO users (first_name, last_name, middle_name, email) VALUES ($1, $2, $3, $4) RETURNING user_id`
	for _, u := range users {
		userId := 0
		err := s.db.Get(&userId, q, u.FirstName, u.LastName, u.MiddleName, u.Email)
		require.NoError(s.T(), err, "Should normally insert row")
		user, err := repo.GetById(ctx, int64(userId))
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), int64(userId), user.UserID)
		assert.Equal(s.T(), u.FirstName, user.FirstName)
		assert.Equal(s.T(), u.LastName, user.LastName)
		assert.Equal(s.T(), u.MiddleName, user.MiddleName)
		assert.Equal(s.T(), u.Email, user.Email)
	}

	user, err := repo.GetById(ctx, -1)
	assert.Nil(s.T(), user)
	assert.ErrorAs(s.T(), err, &ErrUserNotFound, "user should not exist")
}

func (s *UserRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()
	repo := NewUserRepository(NewDatabase(s.db))

	u := model.UserCreate{FirstName: "John", LastName: "Doe", MiddleName: "Jr.", Email: "johndoe@example.com"}
	q := `INSERT INTO users (first_name, last_name, middle_name, email) VALUES ($1, $2, $3, $4) RETURNING user_id`

	userId := 0
	err := s.db.Get(&userId, q, u.FirstName, u.LastName, u.MiddleName, u.Email)
	require.NoError(s.T(), err, "Should normally insert row")
	user, err := repo.Update(ctx, int64(userId), map[string]interface{}{
		"first_name": "Michel",
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(userId), user.UserID)
	assert.Equal(s.T(), "Michel", user.FirstName)
	assert.Equal(s.T(), u.LastName, user.LastName)
	assert.Equal(s.T(), u.MiddleName, user.MiddleName)
	assert.Equal(s.T(), u.Email, user.Email)

	user, err = repo.Update(ctx, -1, map[string]interface{}{
		"first_name": "Michel",
	})
	assert.Nil(s.T(), user, "user should not be returned if does not exist")
	assert.ErrorAs(s.T(), err, &ErrUserNotFound, "user should not exist")

	user, err = repo.Update(ctx, int64(userId), map[string]interface{}{
		"random_key": "random_value",
	})

	assert.Nil(s.T(), user, "user should not be returned if query is invalid")
	assert.ErrorAs(s.T(), err, &ErrLogicError, "user should not exist")
}

func (s *UserRepositoryTestSuite) TestDelete() {
	ctx := context.Background()
	repo := NewUserRepository(NewDatabase(s.db))

	u := model.UserCreate{FirstName: "John", LastName: "Doe", MiddleName: "Jr.", Email: "johndoe@example.com"}
	q := `INSERT INTO users (first_name, last_name, middle_name, email) VALUES ($1, $2, $3, $4) RETURNING user_id`

	userId := 0
	err := s.db.Get(&userId, q, u.FirstName, u.LastName, u.MiddleName, u.Email)
	require.NoError(s.T(), err, "Should normally insert row")
	err = repo.Delete(ctx, int64(userId))
	assert.NoError(s.T(), err)
	q = `SELECT count(1) FROM users WHERE user_id = $1`
	count := -1
	err = s.db.Get(&count, q, userId)
	require.NoError(s.T(), err, "should correctly count remaining rows")
	assert.Equal(s.T(), 0, count, "should truly delete user")

	err = repo.Delete(ctx, -1)
	assert.ErrorAs(s.T(), err, &ErrUserNotFound, "should return right error")
}

func (s *UserRepositoryTestSuite) TestGetByEmail() {
	ctx := context.Background()
	repo := NewUserRepository(NewDatabase(s.db))
	for i := 1; i < 5; i++ {
		expected := CreateRandomUser(ctx, NewDatabase(s.db), s.T())
		s.T().Run(expected.Email, func(t *testing.T) {
			actual, err := repo.GetByEmail(ctx, expected.Email)
			assert.NoError(s.T(), err, "should find user")
			assert.Equal(s.T(), expected.UserID, actual.UserID)
			assert.Equal(s.T(), expected.Email, actual.Email)
			assert.Equal(s.T(), expected.FirstName, actual.FirstName)
			assert.Equal(s.T(), expected.LastName, actual.LastName)
			assert.Equal(s.T(), expected.MiddleName, actual.MiddleName)
		})
	}
}
