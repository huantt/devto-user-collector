package github

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	httpClient *resty.Client
}

const Endpoint = "https://api.github.com"

func NewService(APIEndpoint string) *Service {
	httpClient := resty.New()
	httpClient.
		SetRetryCount(12).
		SetRetryWaitTime(5 * time.Second).
		SetBaseURL(APIEndpoint).AddRetryCondition(func(response *resty.Response, err error) bool {
		if err != nil {
			return true
		}
		if response.StatusCode() == http.StatusInternalServerError ||
			response.StatusCode() == http.StatusBadGateway ||
			response.StatusCode() == http.StatusGatewayTimeout ||
			response.StatusCode() == http.StatusServiceUnavailable {
			slog.Warn(fmt.Sprintf("Response status code is %d - Request: %s - Body: %s - Retrying...", response.StatusCode(), response.Request.URL, response.Body()))
			return true
		}

		return false
	})
	return &Service{httpClient}
}

func (s *Service) GetUserInfo(ctx context.Context, username string) (*User, error) {
	var userInfo User
	resp, err := s.httpClient.R().SetContext(ctx).
		SetResult(&userInfo).
		SetPathParam("username", username).
		Get("/users/{username}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(fmt.Sprintf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body()))
	}
	return &userInfo, nil
}

const MaxLimit = 100

func (s *Service) GetFollowers(ctx context.Context, username string, page, limit int) ([]User, error) {
	if limit > MaxLimit {
		return nil, fmt.Errorf("limit must be less than %d", MaxLimit)
	}
	var followers []User
	resp, err := s.httpClient.R().SetContext(ctx).
		SetResult(&followers).
		SetPathParam("username", username).
		SetQueryParam("page", fmt.Sprintf("%d", page)).
		SetQueryParam("per_page", fmt.Sprintf("%d", limit)).
		Get("/users/{username}/followers")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(fmt.Sprintf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body()))
	}
	return followers, nil
}
