package handlers

import (
	"favorite_service/internal/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type favoriteItemService interface {
	GetFavoriteItem(listId int) ([]models.FavoriteItem, error)
	CreateFavoriteItem(item models.CreateFavoriteItem) (models.FavoriteItem, error)
	DeleteFavoriteItem(listId int, itemId int) error
}

type FavoriteItemHandler struct {
	favoriteItemService favoriteItemService
}

func NewFavoriteItemHandler(favoriteItemRepository favoriteItemService) *FavoriteItemHandler {
	return &FavoriteItemHandler{
		favoriteItemService: favoriteItemRepository,
	}
}

func (h *FavoriteItemHandler) GetFavoriteItemHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}

	itemList, err := h.favoriteItemService.GetFavoriteItem(listId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{FailData: err})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: itemList})

}

func (h *FavoriteItemHandler) CreateFavoriteItemHandle(c *fiber.Ctx) error {

	newItem := models.CreateFavoriteItem{}

	if err := c.BodyParser(&newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: "err1"})
	}

	err := newItem.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: "err2"})
	}

	favoriteItem, err := h.favoriteItemService.CreateFavoriteItem(newItem)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{FailData: "err3"})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteItem})

}

func (h *FavoriteItemHandler) DeleteFavoriteItemHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err})
	}

	itemIdStr := c.Query("itemId")
	if itemIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: "itemId bilgisi zorunlu"})
	}

	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{FailData: err.Error()})
	}

	err = h.favoriteItemService.DeleteFavoriteItem(listId, itemId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{FailData: err})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: nil})

}

func (h *FavoriteItemHandler) SetRoutes(app *fiber.App) {

	itemGroup := app.Group("/items")

	itemGroup.Get("/:listId", h.GetFavoriteItemHandle)
	itemGroup.Post("", h.CreateFavoriteItemHandle)
	itemGroup.Delete("/:listId/item", h.DeleteFavoriteItemHandle)

}
