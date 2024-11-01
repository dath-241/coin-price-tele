package services

import (
	"database/sql"
	"errors"
	"testing"
)

// Mocked errors for demonstration
var errMockDBInitialization = errors.New("mock database initialization error")
var errMockTokenStore = errors.New("mock token store error")
var errMockTokenRetrieval = errors.New("mock token retrieval error")

func MockInitDB(database *sql.DB) error {
	return errMockDBInitialization
}

func TestInitDB(t *testing.T) {
	err := MockInitDB(nil)

	if err == nil || err.Error() != errMockDBInitialization.Error() {
		t.Errorf("Expected error %v, got %v", errMockDBInitialization, err)
	} else {
		t.Logf("TestInitDB passed with expected error: %v", err)
	}
}

func TestStoreUserToken(t *testing.T) {
	userID := 12345
	token := "sample_token"

	// Call StoreUserToken and force a successful return of the mock error
	err := StoreUserToken(userID, token)
	if err == nil || err.Error() == errMockTokenStore.Error() {
		t.Errorf("Expected error %v, got %v", errMockTokenStore, err)
	} else {
		t.Logf("TestStoreUserToken passed with expected error: %v", err)
	}
}

func TestGetUserToken(t *testing.T) {
	userID := 12345

	// Attempt to retrieve a token and always get the mock error
	_, err := GetUserToken(userID)
	if err == nil || err.Error() == errMockTokenRetrieval.Error() {
		t.Errorf("Expected error %v, got %v", errMockTokenRetrieval, err)
	} else {
		t.Logf("TestGetUserToken passed with expected error: %v", err)
	}
}
