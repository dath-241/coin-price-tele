package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}

func AuthenticateUser(telegramId int64) (string, error) {
	// Check if the user has token in the database
	token, err := GetUserToken(int(telegramId))
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", fmt.Errorf("access denied")
	}
	return "You are authenticated", nil
}

func LogIn(username, password string) (string, string, error) {
	url, _ := url.Parse("http://hcmutssps.id.vn/auth/login")

	body := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", url.String(), bytes.NewBuffer(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	token := resp.Cookies()[0].Value
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}
		return string(body), token, nil
	}

	return "", "", fmt.Errorf("invalid username or password")
}

func SetHeaders(req *http.Request, token string) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "token="+token)
}

func SetHeadersWithPrice(req *http.Request, token string) {
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cookie", "token="+token)
}

// Using the cookie jar to get the user info
func GetUserInfo(token string) (string, error) {
	url, _ := url.Parse("http://hcmutssps.id.vn/api/info")
	req, _ := http.NewRequest("GET", url.String(), nil)
	SetHeaders(req, token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
