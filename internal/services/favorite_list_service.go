package services

import (
	"favorite_service/internal/models"
)

type favoriteListRepository interface {
	GetFavoriteList(userId int) ([]models.FavoriteList, error)
	CreateFavoriteList(favoriteList *models.FavoriteList) error
	UpdateFavoriteList(id int, updtList models.UpdateFavoriteList) (models.FavoriteList, error)
	DeleteFavoriteList(listId int) error
	GetListOwner(listId int) (models.FavoriteList, error)
}

type favoriteItemRepository interface {
	GetFavoriteItem(listId int) ([]models.FavoriteItem, error)
	DeleteFavoriteItemsByListId(listId int) error
}

type favoriteListProductClient interface {
	VerifyProduct(productId int) (*models.Product, error)
}

type favoriteListUserClient interface {
	VerifyUser(token string) (*models.Users, error)
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

func (s *FavoriteListService) GetUserFavoriteListsWithItems(token string) ([]models.FavoriteListResponse, error) {

	user, err := s.favoriteListUserClient.VerifyUser(token)

	if err != nil {
		return nil, err
	}

	var response []models.FavoriteListResponse

	lists, err := s.listRepo.GetFavoriteList(user.ID)
	if err != nil {
		return nil, err
	}

	for _, list := range lists {
		items, err := s.itemRepo.GetFavoriteItem(list.Id)
		if err != nil {
			return nil, err
		}

		products, err := GetProductInfo(items, s.favoriteListProductClient)
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

func (s *FavoriteListService) CreateFavoriteList(list *models.FavoriteList, token string) error {

	user, err := s.favoriteListUserClient.VerifyUser(token)

	if err != nil {
		return err
	}

	list.UserId = user.ID

	return s.listRepo.CreateFavoriteList(list)
}

func (s *FavoriteListService) UpdateFavoriteList(listId int, list models.UpdateFavoriteList, token string) (models.FavoriteList, error) {

	user, err := s.favoriteListUserClient.VerifyUser(token)

	if err != nil {
		return models.FavoriteList{}, err
	}

	ownerFavoriteList, err := s.listRepo.GetListOwner(listId)

	if err != nil {
		return models.FavoriteList{}, models.ErrunaUthorizedAction
	}

	if ownerFavoriteList.UserId != user.ID {
		return models.FavoriteList{}, models.ErrunaUthorizedAction
	}

	return s.listRepo.UpdateFavoriteList(listId, list)
}

func (s *FavoriteListService) DeleteFavoriteList(listId int, token string) error {

	user, err := s.favoriteListUserClient.VerifyUser(token)

	if err != nil {
		return err
	}

	ownerFavoriteList, err := s.listRepo.GetListOwner(listId)

	if err != nil {
		return err
	}

	if ownerFavoriteList.UserId != user.ID {

		return models.ErrUserUnauthorized

	}

	if err := s.itemRepo.DeleteFavoriteItemsByListId(listId); err != nil {
		return err
	}

	return s.listRepo.DeleteFavoriteList(listId)
}
