package client

import (
	"context"
	"encoding/json"
	"favorite_service/internal/models"
	"fmt"
	"net/http"
	"time"
)

type circuitBreaker interface {
	Execute(req func() (interface{}, error)) (interface{}, error)
}

type UserClient struct {
	userServiceURL string
	cb             circuitBreaker
}

func NewUserClient(userServiceURL string, cb circuitBreaker) *UserClient {

	return &UserClient{
		userServiceURL: userServiceURL,
		cb:             cb,
	}
}

func (c *UserClient) VerifyUser(token string, ctx context.Context) (*models.Users, error) {

	result, err := c.cb.Execute(func() (interface{}, error) {

		userServiceURL := fmt.Sprintf("%s/users/me", c.userServiceURL)

		req, err := http.NewRequestWithContext(ctx, "GET", userServiceURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", token)

		client := &http.Client{

			Timeout: 8 * time.Second,
		}

		resp, err := client.Do(req)

		select {
		case <-ctx.Done():
			fmt.Println("Context zaman aşımına uğradı")
		default:
			fmt.Println("Contex zaman aşımına uğramadı")
		}

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {

			return nil, models.ErrUserUnauthorized
		}

		var userResponse models.UserResponse

		if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {

			return nil, err

		}

		return &userResponse.SuccesData, nil

	})

	if err != nil {
		return nil, err
	}

	return result.(*models.Users), nil

}
