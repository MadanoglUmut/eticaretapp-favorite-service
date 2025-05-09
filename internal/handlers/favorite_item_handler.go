package handlers

import (
	"context"
	"favorite_service/internal/models"
	"strconv"
	"time"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type favoriteItemService interface {
	GetFavoriteItem(listId int, token string, ctx context.Context) ([]models.Product, error)
	CreateFavoriteItem(item models.CreateFavoriteItem, token string, ctx context.Context) (models.FavoriteItem, error)
	DeleteFavoriteItem(listId int, itemId int, token string, ctx context.Context) error
}

type metric interface {
	ObserveHandler(name string, startTime time.Time, status int)
}

type FavoriteItemHandler struct {
	favoriteItemService favoriteItemService
	metric              metric
}

func NewFavoriteItemHandler(favoriteItemRepository favoriteItemService, metric metric) *FavoriteItemHandler {
	return &FavoriteItemHandler{
		favoriteItemService: favoriteItemRepository,
		metric:              metric,
	}
}

func (h *FavoriteItemHandler) GetFavoriteItemHandle(c *fiber.Ctx) error {

	defer func() {
		h.metric.ObserveHandler("FavoriteItemHandler_GetFavoriteItem", time.Now(), c.Response().StatusCode())
	}()

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "List Id Bulunamadi",
			Details: err.Error()},
		)
	}

	autHeader := c.Get("Authorization")

	if autHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErorResponse{
			Error:   "Token Authorization Hatasi",
			Details: "Token",
		})
	}

	ctx := c.UserContext()

	itemList, err := h.favoriteItemService.GetFavoriteItem(listId, autHeader, ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{
			Error:   "Servis Hatası",
			Details: err.Error()},
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: itemList})

}

func (h *FavoriteItemHandler) CreateFavoriteItemHandle(c *fiber.Ctx) error {

	defer func() {
		h.metric.ObserveHandler("FavoriteItemHandler_CreateFavoriteItem", time.Now(), c.Response().StatusCode())
	}()

	newItem := models.CreateFavoriteItem{}

	if err := c.BodyParser(&newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Body Parse Hatasi",
			Details: err.Error()},
		)
	}

	err := newItem.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Validate Hatasi",
			Details: err.Error()},
		)
	}

	autHeader := c.Get("Authorization")

	if autHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErorResponse{
			Error:   "Token Authorization Hatasi",
			Details: "Token"},
		)
	}

	ctx := c.UserContext()

	favoriteItem, err := h.favoriteItemService.CreateFavoriteItem(newItem, autHeader, ctx)

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{
			Error:   "Servis Hatasi",
			Details: err.Error()},
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteItem})

}

func (h *FavoriteItemHandler) DeleteFavoriteItemHandle(c *fiber.Ctx) error {

	defer func() {
		h.metric.ObserveHandler("FavoriteItemHandler_DeleteFavoriteItem", time.Now(), c.Response().StatusCode())
	}()

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Parametre Hatasi",
			Details: err.Error()},
		)
	}

	itemIdStr := c.Query("itemId")
	if itemIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Query Parametre Hatasi",
			Details: "Query Param Zorunlu"},
		)
	}

	itemId, err := strconv.Atoi(itemIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Item Id Tip Dönüşüm Hatasi",
			Details: err.Error()},
		)
	}

	autHeader := c.Get("Authorization")

	if autHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErorResponse{
			Error:   "Token Authorization Hatasi",
			Details: "Token"},
		)
	}

	ctx := c.UserContext()

	err = h.favoriteItemService.DeleteFavoriteItem(listId, itemId, autHeader, ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{
			Error:   "Servis Hatasi",
			Details: err.Error()},
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: nil})

}

func (h *FavoriteItemHandler) SetRoutes(app *fiber.App) {

	itemGroup := app.Group("/items")

	itemGroup.Get("/:listId", h.GetFavoriteItemHandle)
	itemGroup.Post("", h.CreateFavoriteItemHandle)
	itemGroup.Delete("/:listId/item", h.DeleteFavoriteItemHandle)

}

func ItemGetEndpoints() []*endpoint.EndPoint {
	return []*endpoint.EndPoint{
		endpoint.New(
			endpoint.GET,
			"/items/{listId}",
			endpoint.WithTags("item"),
			endpoint.WithParams(parameter.IntParam("listId", parameter.Path, parameter.WithRequired())),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New([]models.Product{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "400", "Bad Request")}),
		),

		endpoint.New(
			endpoint.POST,
			"/items",
			endpoint.WithTags("item"),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithBody(models.CreateFavoriteItem{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.FavoriteItem{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "400", "Bad Request")}),
		),

		endpoint.New(
			endpoint.DELETE,
			"/items/{listId}/item",
			endpoint.WithTags("item"),
			endpoint.WithParams(parameter.IntParam("listId", parameter.Path, parameter.WithRequired())),
			endpoint.WithParams(parameter.IntParam("itemId", parameter.Query, parameter.WithRequired())),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccesResponse{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "400", "Bad Request")}),
		),
	}
}
