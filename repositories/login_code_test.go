package repositories

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type LoginCodeRepositoryTestSuite struct {
	DBTestSuite
}

func TestLoginCodeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &LoginCodeRepositoryTestSuite{
		*DBTestSuiteFromEnv(),
	})
}

func (s *LoginCodeRepositoryTestSuite) TestCreate() {
	ctx := context.Background()
	repo := LoginCodeRepository{
		Db:      NewDatabase(s.db),
		CodeTTL: 24 * time.Hour,
	}
	user := CreateRandomUser(ctx, NewDatabase(s.db), s.T())
	err := repo.CreateLoginCode(ctx, user.UserID, "1234")
	assert.NoError(s.T(), err, "should create code without errors")

	var insertedCode string
	var insertedUserId int64
	var isUsed bool
	var expiresAt time.Time
	q := `SELECT user_id, code, is_used, expires_at FROM login_code WHERE user_id = $1 AND code = $2`
	row := s.db.QueryRow(q, user.UserID, "1234")
	err = row.Scan(&insertedUserId, &insertedCode, &isUsed, &expiresAt)
	require.NoError(s.T(), err, "row should be fetched without errors")
	assert.Equal(s.T(), false, isUsed)
	assert.Equal(s.T(), "1234", insertedCode)
	assert.Equal(s.T(), user.UserID, insertedUserId)
}

func (s *LoginCodeRepositoryTestSuite) TestMarkCodeUsed() {
	ctx := context.Background()
	repo := LoginCodeRepository{
		Db:      NewDatabase(s.db),
		CodeTTL: 24 * time.Hour,
	}

	user := CreateRandomUser(ctx, NewDatabase(s.db), s.T())
	err := repo.CreateLoginCode(ctx, user.UserID, "1234")
	require.NoError(s.T(), err, "should create code without error")

	err = repo.MarkCodeUsed(ctx, user.UserID, "1234")
	assert.NoError(s.T(), err)

	err = repo.MarkCodeUsed(ctx, user.UserID, "1234")
	assert.ErrorIs(s.T(), err, ErrCodeNotFound)
}
