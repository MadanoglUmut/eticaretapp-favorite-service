package handlers

import (
	"context"
	"favorite_service/internal/models"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type favoriteListService interface {
	GetUserFavoriteListsWithItems(token string, ctx context.Context) ([]models.FavoriteListResponse, error)
	CreateFavoriteList(list *models.FavoriteList, token string, ctx context.Context) error
	UpdateFavoriteList(listId int, list models.UpdateFavoriteList, token string, ctx context.Context) (models.FavoriteList, error)
	DeleteFavoriteList(listId int, token string, ctx context.Context) error
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

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErorResponse{
			Error:   "Token Authorization Hatasi",
			Details: "Token"},
		)
	}

	ctx := c.UserContext()

	favoriteList, err := h.favoriteListService.GetUserFavoriteListsWithItems(authHeader, ctx)

	if err != nil {
		if err == models.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErorResponse{
				Error:   "Kullanıcı bulunamadı",
				Details: err.Error()},
			)

		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{
			Error:   "Servis hatası",
			Details: err.Error()},
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) CreateFavoriteListHandle(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErorResponse{
			Error:   "Token Authorization Hatasi",
			Details: "Token"},
		)

	}

	list := models.CreateFavoriteList{}
	if err := c.BodyParser(&list); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Body Parse Hatasi",
			Details: err.Error()},
		)
	}
	err := list.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Validate Hatasi",
			Details: err.Error()},
		)
	}

	favoriteList := models.FavoriteList{
		ListName: list.ListName,
	}

	ctx := c.UserContext()

	err = h.favoriteListService.CreateFavoriteList(&favoriteList, authHeader, ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{
			Error:   "Servis Hatasi",
			Details: err.Error()},
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) UpdateFavoriteListHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "List Id Bulunamadi",
			Details: err.Error()},
		)
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErorResponse{
			Error:   "Token Authorization Hatasi",
			Details: "Token"},
		)

	}

	list := models.UpdateFavoriteList{}

	if err := c.BodyParser(&list); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Body Parse Hatasi",
			Details: err.Error()},
		)
	}

	err = list.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Validate Hatasi",
			Details: err.Error()},
		)
	}

	ctx := c.UserContext()

	favoriteList, err := h.favoriteListService.UpdateFavoriteList(listId, list, authHeader, ctx)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{
			Error:   "Servis Hatasi",
			Details: err.Error()},
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) DeleteFavoriteListHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "List Id Bulunamadi",
			Details: err.Error()},
		)
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErorResponse{
			Error:   "Token Authorization Hatasi",
			Details: "Token"},
		)

	}

	ctx := c.UserContext()

	err = h.favoriteListService.DeleteFavoriteList(listId, authHeader, ctx)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{
			Error:   "Servis Hatasi",
			Details: err.Error()},
		)
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: "Silindi"})

}

func (h *FavoriteListHandler) SetRoutes(app *fiber.App) {
	listGroup := app.Group("/lists")

	listGroup.Get("/", h.GetUserFavoriteListsWithItemsHandle)
	listGroup.Post("", h.CreateFavoriteListHandle)
	listGroup.Put("/:listId", h.UpdateFavoriteListHandle)
	listGroup.Delete("/:listId", h.DeleteFavoriteListHandle)
}

func ListGetEndpoints() []*endpoint.EndPoint {

	return []*endpoint.EndPoint{
		endpoint.New(
			endpoint.GET,
			"/lists",
			endpoint.WithTags("lists"),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New([]models.FavoriteListResponse{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),

		endpoint.New(
			endpoint.POST,
			"/lists",
			endpoint.WithTags("lists"),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithBody(models.CreateFavoriteList{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.FavoriteList{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),

		endpoint.New(
			endpoint.PUT,
			"/lists/{listId}",
			endpoint.WithTags("lists"),
			endpoint.WithParams(parameter.IntParam("listId", parameter.Path, parameter.WithRequired())),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithBody(models.UpdateFavoriteList{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.FavoriteList{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),
		endpoint.New(
			endpoint.DELETE,
			"/lists/{listId}",
			endpoint.WithTags("lists"),
			endpoint.WithParams(parameter.IntParam("listId", parameter.Path, parameter.WithRequired())),
			endpoint.WithParams(parameter.StrParam("Authorization", parameter.Header, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccesResponse{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),
	}

}
