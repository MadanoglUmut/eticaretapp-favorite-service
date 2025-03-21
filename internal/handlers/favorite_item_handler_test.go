package handlers

import (
	"bytes"
	"encoding/json"
	"favorite_service/internal/models"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteItemHandler(t *testing.T) {

	t.Run("", func(t *testing.T) {

		request := httptest.NewRequest("GET", "/items/1", nil)

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

	})

	t.Run("TestGetFavoriteItemHandleNotFound", func(t *testing.T) {

		request := httptest.NewRequest("GET", "/items/999", nil)

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	})

	t.Run("TestCreateFavoriteItemHandle", func(t *testing.T) {

		newItems := models.CreateFavoriteItem{
			ItemId: 10,
			ListId: 1,
		}

		body, err := json.Marshal(newItems)

		assert.Nil(t, err)

		request := httptest.NewRequest("POST", "/items", bytes.NewReader(body))

		request.Header.Set("Content-Type", "application/json")

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

	})

	t.Run("TestCreateFavoriteItemHandleFail", func(t *testing.T) {

		newItems := models.CreateFavoriteItem{
			ItemId: 10,
			ListId: 100,
		}

		body, err := json.Marshal(newItems)

		assert.Nil(t, err)

		request := httptest.NewRequest("POST", "/items", bytes.NewReader(body))

		request.Header.Set("Content-Type", "application/json")

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

	})

	t.Run("TestDeleteFavoriteItemHandle", func(t *testing.T) {

		request := httptest.NewRequest("DELETE", "/items/1/item?itemId=1", nil)

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

	})

	t.Run("TestDeleteFavoriteItemHandleError", func(t *testing.T) {
		request := httptest.NewRequest("DELETE", "/items/1/item?itemId=99", nil)

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
	})

}
