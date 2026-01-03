package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleUserInfo struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	Picture       string
}

type OAuthService struct {
	googleConfig *oauth2.Config
}

func NewOAuthService(clientID, clientSecret, redirectURL string) *OAuthService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthService{
		googleConfig: config,
	}
}

// VerifyGoogleToken verifies a Google ID token and returns user info
func (s *OAuthService) VerifyGoogleToken(ctx context.Context, idToken string) (*GoogleUserInfo, error) {
	if idToken == "" {
		return nil, fmt.Errorf("id token is empty")
	}

	// Use tokeninfo endpoint to verify token
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", idToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token verification failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenInfo struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified interface{} `json:"email_verified"` // Can be bool or string
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Error         string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to parse token info: %w", err)
	}

	// Check for errors in response
	if tokenInfo.Error != "" {
		return nil, fmt.Errorf("token verification error: %s - %s", tokenInfo.Error, tokenInfo.ErrorDescription)
	}

	// Parse email_verified (can be bool or string "true"/"false")
	verifiedEmail := false
	switch v := tokenInfo.EmailVerified.(type) {
	case bool:
		verifiedEmail = v
	case string:
		verifiedEmail = v == "true"
	}
	
	return &GoogleUserInfo{
		ID:            tokenInfo.Sub,
		Email:         tokenInfo.Email,
		VerifiedEmail: verifiedEmail,
		Name:          tokenInfo.Name,
		Picture:       tokenInfo.Picture,
	}, nil
}

// GetGoogleUserInfo fetches user info using the OAuth2 service
func (s *OAuthService) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	service, err := googleoauth2.NewService(ctx, option.WithTokenSource(s.googleConfig.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	verified := false
	if userInfo.VerifiedEmail != nil {
		verified = *userInfo.VerifiedEmail
	}

	return &GoogleUserInfo{
		ID:            userInfo.Id,
		Email:         userInfo.Email,
		VerifiedEmail: verified,
		Name:          userInfo.Name,
		Picture:       userInfo.Picture,
	}, nil
}

