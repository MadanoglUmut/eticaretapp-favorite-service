package handlers

import (
	"favorite_service/internal/models"
	"strconv"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type favoriteListService interface {
	GetUserFavoriteListsWithItems(userId int) ([]models.FavoriteListResponse, error)
	CreateFavoriteList(list *models.FavoriteList) error
	UpdateFavoriteList(listId int, list models.UpdateFavoriteList, userId int) (models.FavoriteList, error)
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

		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "User Id Parametre Hatasi", Details: err.Error()})

	}

	favoriteList, err := h.favoriteListService.GetUserFavoriteListsWithItems(userId)

	if err != nil {
		if err == models.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(models.ErorResponse{Error: "Kullanici Bulunamadi", Details: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{Error: "Servis Hatasi", Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) CreateFavoriteListHandle(c *fiber.Ctx) error {

	list := models.CreateFavoriteList{}
	if err := c.BodyParser(&list); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "Body Parse Hatasi", Details: err.Error()})
	}
	err := list.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "Validate Hatasi", Details: err.Error()})
	}

	favoriteList := models.FavoriteList{
		ListName: list.ListName,
		UserId:   list.UserId,
	}

	err = h.favoriteListService.CreateFavoriteList(&favoriteList)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{Error: "Servis Hatasi", Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) UpdateFavoriteListHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "List Id Bulunamadi", Details: err.Error()})
	}

	userIdStr := c.Query("userId")

	if userIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "Query Parametre Hatasi", Details: "Query Parametresi Zorunlu"})
	}

	userId, err := strconv.Atoi(userIdStr)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "User Id Tip Dönüşüm Hatasi", Details: err.Error()})
	}

	list := models.UpdateFavoriteList{}

	if err := c.BodyParser(&list); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "Body Parse Hatasi", Details: err.Error()})
	}

	err = list.Validate()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "Validate Hatasi", Details: err.Error()})
	}

	favoriteList, err := h.favoriteListService.UpdateFavoriteList(listId, list, userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErorResponse{Error: "Servis Hatasi", Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: favoriteList})

}

func (h *FavoriteListHandler) DeleteFavoriteListHandle(c *fiber.Ctx) error {

	listId, err := c.ParamsInt("listId")

	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "List Id Bulunamadi", Details: err.Error()})

	}

	err = h.favoriteListService.DeleteFavoriteList(listId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErorResponse{Error: "Servis Hatasi", Details: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{SuccesData: "Silindi"})

}

func (h *FavoriteListHandler) SetRoutes(app *fiber.App) {
	listGroup := app.Group("/lists")

	listGroup.Get("/:userId", h.GetUserFavoriteListsWithItemsHandle)
	listGroup.Post("", h.CreateFavoriteListHandle)
	listGroup.Put("/:listId/user", h.UpdateFavoriteListHandle)
	listGroup.Delete("/:listId", h.DeleteFavoriteListHandle)
}

func ListGetEndpoints() []*endpoint.EndPoint {

	return []*endpoint.EndPoint{
		endpoint.New(
			endpoint.GET,
			"/lists/{userId}",
			endpoint.WithTags("lists"),
			endpoint.WithParams(parameter.IntParam("userId", parameter.Path, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.FavoriteListResponse{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),

		endpoint.New(
			endpoint.POST,
			"/lists",
			endpoint.WithTags("lists"),
			endpoint.WithBody(models.CreateFavoriteList{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.FavoriteList{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),

		endpoint.New(
			endpoint.PUT,
			"/lists/{listId}/user",
			endpoint.WithTags("lists"),
			endpoint.WithParams(parameter.IntParam("listId", parameter.Path, parameter.WithRequired())),
			endpoint.WithParams(parameter.IntParam("userId", parameter.Query, parameter.WithRequired())),
			endpoint.WithBody(models.UpdateFavoriteList{}),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.FavoriteList{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),
		endpoint.New(
			endpoint.DELETE,
			"/lists/{listId}",
			endpoint.WithTags("lists"),
			endpoint.WithParams(parameter.IntParam("listId", parameter.Path, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.FavoriteList{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.ErorResponse{}, "404", "Bad Request")}),
		),
	}

}
