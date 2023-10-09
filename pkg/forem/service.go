package forem

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// Service docs: https://developers.forem.com/api
type Service struct {
	httpClient  *resty.Client
	rateLimiter *time.Ticker
	authToken   *string
}

func NewService(APIEndpoint string, rps int, proxy string) *Service {
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
		if response.StatusCode() == http.StatusTooManyRequests {
			retryAfterStr := response.Header().Get("Retry-After")
			retryAfter, err := strconv.ParseInt(retryAfterStr, 10, 64)
			if err != nil {
				slog.Error(fmt.Sprintf("Failed to parse Retry-After header - Err: %s", err.Error()))
				return false
			}
			slog.Warn(fmt.Sprintf("Rate limit exceeded, retry after %d seconds", retryAfter))
			time.Sleep(time.Duration(retryAfter) * time.Second)
			return true
		}

		return false
	})
	if proxy != "" {
		httpClient.
			SetCloseConnection(true). //To rotate proxies
			SetProxy(proxy)
	}
	return &Service{httpClient: httpClient, rateLimiter: time.NewTicker(time.Second / time.Duration(rps))}
}

func NewAuthenticatedService(APIEndpoint string, rps int, authToken string, proxy string) *Service {
	service := NewService(APIEndpoint, rps, proxy)
	service.authToken = &authToken
	return service
}

func (s *Service) GetUser(ctx context.Context, userID int64) (*User, error) {
	<-s.rateLimiter.C
	slog.Debug(fmt.Sprintf("Requesting user by id: %d", userID))
	var user *User
	resp, err := s.httpClient.R().SetContext(ctx).
		SetResult(&user).
		SetPathParam("id", fmt.Sprintf("%d", userID)).
		Get("/api/users/{id}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		if resp.StatusCode() == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body())
	}
	return user, nil
}

// GetFollowers authentication is required
func (s *Service) GetFollowers(ctx context.Context, userID int, page int) ([]Follower, error) {
	<-s.rateLimiter.C
	slog.Debug(fmt.Sprintf("Requesting user followers by id: %d", userID))
	var followers []Follower
	resp, err := s.httpClient.R().SetContext(ctx).
		SetResult(&followers).
		SetQueryParam("page", strconv.Itoa(page)).
		SetPathParam("id", strconv.Itoa(userID)).
		Get("/api/followers/users")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("Request: %s - Response code: %d - Response body: %s", resp.Request.URL, resp.StatusCode(), resp.Body())
	}
	return followers, nil
}

func (s *Service) GetIP() {
	resp, err := s.httpClient.R().Get("https://api.myip.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.String())
}
