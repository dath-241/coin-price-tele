package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	// "strings"
)

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}
type RespondBody struct {
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Path      string `json:"path"`
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
	url, _ := url.Parse(apiUrl + "/auth/login")

	body := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", url.String(), bytes.NewBuffer(jsonBody))
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	// If the response is 401, return the error message
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("invalid username or password")
	}
	token := resp.Cookies()[0].Value
	defer resp.Body.Close()

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %v", err)
	}

	var respond RespondBody
	if err := json.Unmarshal(message, &respond); err != nil {
		return "", "", fmt.Errorf("error unmarshalling response: %v", err)
	}
	if resp.StatusCode == http.StatusBadRequest {
		return "", "", fmt.Errorf("%s", respond.Message)
	}

	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}
		return respond.Message, token, nil
	}

	return "", "", fmt.Errorf("invalid username or password")
}

func Regsiter(email, name, username, password string) (string, error) {
	url := apiUrl + "/auth/register"
	body := map[string]string{
		"email":    email,
		"name":     name,
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	// token := resp.Cookies()[0].Value
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var respond RespondBody
	if err := json.Unmarshal(message, &respond); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}
	if resp.StatusCode == http.StatusBadRequest {
		return "", fmt.Errorf("%s", respond.Message)
	}

	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return respond.Message, nil
	}

	return "", fmt.Errorf("something wrong?")
}

func ForgotPassword(username string) (string, error) {
	//?OTP
	url := apiUrl + "/auth/forgotPassword?username=" + username
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var respond RespondBody
	if err := json.Unmarshal(message, &respond); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}
	if resp.StatusCode == http.StatusBadRequest {
		return "", fmt.Errorf("%s", respond.Message)
	}
	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return respond.Message, nil
	}
	return "", fmt.Errorf("something wrong?")
}

func Testadmin(username, token string) (string, error) {
	url := apiUrl + "/admin/removeUserByUsername?username=" + username
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}
	// Optionally, set headers if needed
	req.Header.Set("Authorization", token) // Example header
	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// var respond RespondBody
	// if err := json.Unmarshal(message, &respond); err != nil {
	// 	return "", fmt.Errorf("error unmarshalling response: %v", err)
	// }
	if resp.StatusCode == http.StatusBadRequest {
		return "", fmt.Errorf("%s", string(message))
	}
	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(message), nil
	}
	return "", fmt.Errorf(string(message))
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
	url, _ := url.Parse(apiUrl + "/api/info")
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
