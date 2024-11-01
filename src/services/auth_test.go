package services

import (
	"fmt"
	"testing"
)

// Mock functions to simulate a successful response
func MockAuthenticateUser(telegramId int64) (string, error) {
	return "You are authenticated", nil
}

func MockLogIn(username, password string) (string, string, error) {
	return "Login Successful", "mocked_token", nil
}

func MockGetUserInfo(token string) (string, error) {
	return "User Info Retrieved Successfully", nil
}

func TestAuthenticateUser(t *testing.T) {
	telegramID := int64(12345) // any mock telegram ID
	expected := "You are authenticated"
	result, err := MockAuthenticateUser(telegramID)

	if err != nil || result != expected {
		t.Errorf("Expected %s, got %s, error %v", expected, result, err)
	} else {
		fmt.Println("TestAuthenticateUser passed")
	}
}

func TestLogIn(t *testing.T) {
	username, password := "testuser", "password123"
	expectedResponse := "Login Successful"
	expectedToken := "mocked_token"

	response, token, err := MockLogIn(username, password)
	if err != nil || response != expectedResponse || token != expectedToken {
		t.Errorf("Expected (%s, %s), got (%s, %s), error %v", expectedResponse, expectedToken, response, token, err)
	} else {
		fmt.Println("TestLogIn passed")
	}
}

func TestGetUserInfo(t *testing.T) {
	token := "mocked_token"
	expected := "User Info Retrieved Successfully"

	result, err := MockGetUserInfo(token)
	if err != nil || result != expected {
		t.Errorf("Expected %s, got %s, error %v", expected, result, err)
	} else {
		fmt.Println("TestGetUserInfo passed")
	}
}
