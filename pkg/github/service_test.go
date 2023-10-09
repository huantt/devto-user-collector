package github

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUser(t *testing.T) {
	username := "huantt"
	service := NewService(Endpoint)
	user, err := service.GetUserInfo(context.Background(), username)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Login)
}

func TestGetFollowers(t *testing.T) {
	username := "huantt"
	service := NewService(Endpoint)
	followers, err := service.GetFollowers(context.Background(), username, 1, 30)
	assert.NoError(t, err)
	assert.NotEmpty(t, followers)
}
