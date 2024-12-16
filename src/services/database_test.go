package services

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestStoreUserToken(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db = mockDB

	// Mock token insertion
	mock.ExpectExec(`INSERT INTO user_tokens`).
		WithArgs(123, "sample-token").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = StoreUserToken(123, "sample-token")
	assert.NoError(t, err)

	// Verify expectations
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserToken(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db = mockDB

	// Mock token retrieval
	mock.ExpectQuery(`SELECT token FROM user_tokens WHERE user_id = \$1`).
		WithArgs(123).
		WillReturnRows(sqlmock.NewRows([]string{"token"}).AddRow("sample-token"))

	token, err := GetUserToken(123)
	assert.NoError(t, err)
	assert.Equal(t, "sample-token", token)

	// Verify expectations
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetMute(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db = mockDB

	// Mock is_muted retrieval
	mock.ExpectQuery(`SELECT is_muted FROM user_tokens WHERE user_id = \$1`).
		WithArgs(123).
		WillReturnRows(sqlmock.NewRows([]string{"is_muted"}).AddRow(true))

	isMuted, err := GetMute(123)
	assert.NoError(t, err)
	assert.True(t, isMuted)

	// Verify expectations
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSetMute(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	db = mockDB

	// Mock is_muted update/insert
	mock.ExpectExec(`INSERT INTO user_tokens`).
		WithArgs(123, true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = SetMute(123, true)
	assert.NoError(t, err)

	// Verify expectations
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
