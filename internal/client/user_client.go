package client

import (
	"favorite_service/internal/models"
	"fmt"
	"net/http"
)

type UserClient struct {
	userServiceURL string
}

func NewUserClient(userServiceURL string) *UserClient {
	return &UserClient{
		userServiceURL: userServiceURL,
	}
}

func (c *UserClient) CheckUserId(userId int) error {

	userServiceURL := fmt.Sprintf("%s/%d", c.userServiceURL, userId)

	resp, err := http.Get(userServiceURL)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {

		return models.ErrUserNotFound
	}

	return nil

}
