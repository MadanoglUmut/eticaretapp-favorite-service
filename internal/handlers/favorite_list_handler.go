package handlers

import (
	"favorite_service/internal/models"

	"github.com/gofiber/fiber/v2"
)

type favoriteListService interface {
	GetUserFavoriteListsWithItems(userId int) ([]models.FavoriteListResponse, error)
	CreateFavoriteList(list *models.FavoriteList) error
	UpdateFavoriteList(id int, list models.UpdateFavoriteList) (models.FavoriteList, error)
	DeleteFavoriteList(id int) error
}

type FavoriteListHandler struct {
	favoriteListService favoriteListService
}

func NewFavoriteListHandler(favoriteListService favoriteListService) *FavoriteListHandler {
	return &FavoriteListHandler{
		favoriteListService: favoriteListService,
	}
}

func (h *FavoriteListHandler) GetUserFavoriteListsWithItemsHandle(c *fiber.Ctx) error {

	userId, err := c.ParamsInt("userId")

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})

	}

	favoriteList, err := h.favoriteListService.GetUserFavoriteListsWithItems(userId)

	if err != nil {
		if err == models.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.FailResponse{FailData: err})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{FailData: err})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) CreateFavoriteListHandle(c *fiber.Ctx) error {

	list := models.CreateFavoriteList{}
	if err := c.BodyParser(&list); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}
	err := list.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}

	favoriteList := models.FavoriteList{
		ListName: list.ListName,
		UserId:   list.UserId,
	}

	err = h.favoriteListService.CreateFavoriteList(&favoriteList)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{FailData: err})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) UpdateFavoriteListHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}

	list := models.UpdateFavoriteList{}

	if err := c.BodyParser(&list); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}

	err = list.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}

	favoriteList, err := h.favoriteListService.UpdateFavoriteList(listId, list)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{FailData: err})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) DeleteFavoriteListHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})

	}

	err = h.favoriteListService.DeleteFavoriteList(listId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: "Silindi"})

}

func (h *FavoriteListHandler) SetRoutes(app *fiber.App) {
	listGroup := app.Group("/lists")

	listGroup.Get("/:userId", h.GetUserFavoriteListsWithItemsHandle)
	listGroup.Post("", h.CreateFavoriteListHandle)
	listGroup.Put("/:listId", h.UpdateFavoriteListHandle)
	listGroup.Delete("/:listId", h.DeleteFavoriteListHandle)
}
