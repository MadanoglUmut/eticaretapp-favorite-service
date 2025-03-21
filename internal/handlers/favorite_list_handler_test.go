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

func TestFavoriteListHandler(t *testing.T) {

	t.Run("TestGetUserFavoriteListsWithItemsHandle", func(t *testing.T) {

		request := httptest.NewRequest("GET", "/lists", nil)

		request.Header.Set("Authorization", "1")

		resp, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	})

	t.Run("TestGetUserFavoriteListsWithItemsHandleNotFound", func(t *testing.T) {

		request := httptest.NewRequest("GET", "/lists", nil)

		request.Header.Set("Authorization", "99")

		resp, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	})

	t.Run("TestCreateFavoriteListHandle", func(t *testing.T) {

		newList := models.CreateFavoriteList{
			ListName: "Deneme",
		}

		body, err := json.Marshal(newList)

		assert.Nil(t, err)

		request := httptest.NewRequest("POST", "/lists", bytes.NewReader(body))

		request.Header.Set("Content-Type", "application/json")

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

	})

	t.Run("TestCreateFavoriteListHandleBadRequest", func(t *testing.T) {

		newList := models.CreateFavoriteList{
			ListName: "",
		}

		body, err := json.Marshal(newList)

		assert.Nil(t, err)

		request := httptest.NewRequest("POST", "/lists", bytes.NewReader(body))

		request.Header.Set("Content-Type", "application/json")

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
	})

	t.Run("TestUpdateFavoriteListHandle", func(t *testing.T) {

		updatedList := models.UpdateFavoriteList{
			ListName: "Deneme-2",
		}

		body, err := json.Marshal(updatedList)

		assert.Nil(t, err)

		request := httptest.NewRequest("PUT", "/lists/1", bytes.NewReader(body))

		request.Header.Set("Content-Type", "application/json")

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

	})

	t.Run("TestUpdateFavoriteListHandleBadRequest", func(t *testing.T) {

		updatedList := models.UpdateFavoriteList{
			ListName: "",
		}

		body, err := json.Marshal(updatedList)

		assert.Nil(t, err)

		request := httptest.NewRequest("PUT", "/lists/1", bytes.NewReader(body))

		request.Header.Set("Content-Type", "application/json")

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	})

	t.Run("TestDeleteFavoriteListHandle", func(t *testing.T) {

		request := httptest.NewRequest("DELETE", "/lists/1", nil)

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	})

	t.Run("TestDeleteFavoriteListHandleBadRequest", func(t *testing.T) {

		request := httptest.NewRequest("DELETE", "/lists/100", nil)

		request.Header.Set("Authorization", "1")

		response, err := app.Test(request)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	})

}
