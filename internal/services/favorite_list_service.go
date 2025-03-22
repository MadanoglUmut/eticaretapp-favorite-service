package services

import (
	"context"
	"favorite_service/internal/models"
)

type favoriteListRepository interface {
	GetFavoriteList(ctx context.Context, userId int) ([]models.FavoriteList, error)
	CreateFavoriteList(ctx context.Context, favoriteList *models.FavoriteList) error
	UpdateFavoriteList(ctx context.Context, listId int, updtList models.UpdateFavoriteList) (models.FavoriteList, error)
	DeleteFavoriteList(ctx context.Context, listId int) error
	GetListOwner(ctx context.Context, listId int) (models.FavoriteList, error)
}

type favoriteItemRepository interface {
	GetFavoriteItem(ctx context.Context, listId int) ([]models.FavoriteItem, error)
	DeleteFavoriteItemsByListId(ctx context.Context, listId int) error
}

type favoriteListProductClient interface {
	VerifyProduct(ctx context.Context, productId int) (*models.Product, error)
}

type favoriteListUserClient interface {
	VerifyUser(token string, ctx context.Context) (*models.Users, error)
}

type FavoriteListService struct {
	listRepo                  favoriteListRepository
	itemRepo                  favoriteItemRepository
	favoriteListProductClient favoriteListProductClient
	favoriteListUserClient    favoriteListUserClient
}

func NewFavoriteListService(
	listRepo favoriteListRepository,
	itemRepo favoriteItemRepository,
	favoriteListProductClient favoriteListProductClient,
	favoriteListUserClient favoriteListUserClient) *FavoriteListService {

	return &FavoriteListService{
		listRepo:                  listRepo,
		itemRepo:                  itemRepo,
		favoriteListProductClient: favoriteListProductClient,
		favoriteListUserClient:    favoriteListUserClient,
	}
}

func (s *FavoriteListService) GetUserFavoriteListsWithItems(token string, ctx context.Context) ([]models.FavoriteListResponse, error) {

	//ctx, cancel := context.WithTimeout(ctx, time.Second*5)

	//defer cancel()

	user, err := s.favoriteListUserClient.VerifyUser(token, ctx)

	if err != nil {
		return nil, err
	}

	var response []models.FavoriteListResponse

	lists, err := s.listRepo.GetFavoriteList(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	for _, list := range lists {
		items, err := s.itemRepo.GetFavoriteItem(ctx, list.Id)
		if err != nil {
			return nil, err
		}

		products, err := GetProductInfo(ctx, items, s.favoriteListProductClient)
		if err != nil {
			return nil, err
		}

		response = append(response, models.FavoriteListResponse{
			ListId:   list.Id,
			ListName: list.ListName,
			Items:    products,
		})
	}

	return response, nil

}

func (s *FavoriteListService) CreateFavoriteList(list *models.FavoriteList, token string, ctx context.Context) error {

	user, err := s.favoriteListUserClient.VerifyUser(token, ctx)

	if err != nil {
		return err
	}

	list.UserId = user.ID

	return s.listRepo.CreateFavoriteList(ctx, list)
}

func (s *FavoriteListService) UpdateFavoriteList(listId int, list models.UpdateFavoriteList, token string, ctx context.Context) (models.FavoriteList, error) {

	user, err := s.favoriteListUserClient.VerifyUser(token, ctx)

	if err != nil {
		return models.FavoriteList{}, err
	}

	ownerFavoriteList, err := s.listRepo.GetListOwner(ctx, listId)

	if err != nil {
		return models.FavoriteList{}, models.ErrunaUthorizedAction
	}

	if ownerFavoriteList.UserId != user.ID {
		return models.FavoriteList{}, models.ErrunaUthorizedAction
	}

	return s.listRepo.UpdateFavoriteList(ctx, listId, list)
}

func (s *FavoriteListService) DeleteFavoriteList(listId int, token string, ctx context.Context) error {

	user, err := s.favoriteListUserClient.VerifyUser(token, ctx)

	if err != nil {
		return err
	}

	ownerFavoriteList, err := s.listRepo.GetListOwner(ctx, listId)

	if err != nil {
		return err
	}

	if ownerFavoriteList.UserId != user.ID {

		return models.ErrUserUnauthorized

	}

	if err := s.itemRepo.DeleteFavoriteItemsByListId(ctx, listId); err != nil {
		return err
	}

	return s.listRepo.DeleteFavoriteList(ctx, listId)
}
