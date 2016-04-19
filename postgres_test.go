package qb

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type PostgresTestSuite struct {
	suite.Suite
	session *Session
}

func (suite *PostgresTestSuite) SetupTest() {

	var err error
	suite.session, err = New("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.session)
}

func (suite *PostgresTestSuite) TestPostgres() {

	type User struct {
		ID          string         `qb:"type:uuid; constraints:primary_key, auto_increment"`
		Email       string         `qb:"constraints:unique, notnull"`
		FullName    string         `qb:"constraints:notnull"`
		Bio         sql.NullString `qb:"type:text; constraints:null"`
		Oscars      int            `qb:"constraints:default(0)"`
		IgnoreField string         `qb:"-"`
	}

	type Session struct {
		ID             int64     `qb:"type:bigserial; constraints:primary_key"`
		UserID         string    `qb:"type:uuid; constraints:ref(user.id)"`
		AuthToken      string    `qb:"type:uuid; constraints:notnull, unique; index"`
		CreatedAt      time.Time `qb:"constraints:notnull"`
		ExpiresAt      time.Time `qb:"constraints:notnull"`
		CompositeIndex `qb:"index:created_at, expires_at"`
	}

	var err error

	suite.session.Metadata().Add(User{})
	suite.session.Metadata().Add(Session{})

	err = suite.session.Metadata().CreateAll()
	assert.Nil(suite.T(), err)

	// add sample user & session
	suite.session.AddAll(
		&User{
			ID:       "b6f8bfe3-a830-441a-a097-1777e6bfae95",
			Email:    "jack@nicholson.com",
			FullName: "Jack Nicholson",
			Bio:      sql.NullString{String: "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", Valid: true},
		}, &Session{
			UserID:    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
			AuthToken: "e4968197-6137-47a4-ba79-690d8c552248",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
	)

	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	// find user
	var user User

	suite.session.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).First(&user)

	assert.Equal(suite.T(), user.Email, "jack@nicholson.com")
	assert.Equal(suite.T(), user.FullName, "Jack Nicholson")
	assert.Equal(suite.T(), user.Bio.String, "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.")

	// select using join
	sessions := []Session{}
	err = suite.session.Select("s.user_id", "s.id", "s.auth_token", "s.created_at", "s.expires_at").
		From("user u").
		InnerJoin("session s", "u.id = s.user_id").
		Where("u.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		All(&sessions)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(sessions), 1)

	assert.Equal(suite.T(), sessions[0].ID, int64(1))
	assert.Equal(suite.T(), sessions[0].UserID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(suite.T(), sessions[0].AuthToken, "e4968197-6137-47a4-ba79-690d8c552248")

	// update user
	update := suite.session.
		Update("user").
		Set(map[string]interface{}{
			"bio": nil,
		}).
		Where("id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Query()

	suite.session.AddQuery(update)
	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	suite.session.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).First(&user)
	assert.Equal(suite.T(), user.Bio, sql.NullString{String: "", Valid: false})

	// delete session
	suite.session.Delete(&Session{AuthToken: "99e591f8-1025-41ef-a833-6904a0f89a38"})
	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.session.Metadata().DropAll())

	// fail model
	type FailModel struct {
		ID int64 `qb:"type:notype"`
	}

	assert.Panics(suite.T(), func() {
		suite.session.Add(FailModel{})
	})
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
