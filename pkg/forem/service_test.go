package forem

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUser(t *testing.T) {
	service := NewService(DevToEndpoint, 3, "")
	user, err := service.GetUser(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
func TestGetFollowers(t *testing.T) {
	service := NewService(DevToEndpoint, 3, "")
	followers, err := service.GetFollowers(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, followers)
}

func TestGetIP(t *testing.T) {
	service := NewService(DevToEndpoint, 3, "http://localhost:8081")
	for i := 0; i < 20; i++ {
		service.GetIP()
	}
}
