package services

import (
	"favorite_service/internal/models"
)

type favoriteListRepository interface {
	GetFavoriteList(userId int) ([]models.FavoriteList, error)
	CreateFavoriteList(favoriteList *models.FavoriteList) error
	UpdateFavoriteList(id int, updtList models.UpdateFavoriteList) (models.FavoriteList, error)
	DeleteFavoriteList(listId int) error
}

type favoriteItemRepository interface {
	GetFavoriteItem(listId int) ([]models.FavoriteItem, error)
	DeleteFavoriteItemsByListId(listId int) error
}

type userClient interface {
	CheckUserId(userId int) error
}

type FavoriteListService struct {
	listRepo   favoriteListRepository
	itemRepo   favoriteItemRepository
	userClient userClient
}

func NewFavoriteListService(listRepo favoriteListRepository, itemRepo favoriteItemRepository, userClient userClient) *FavoriteListService {
	return &FavoriteListService{
		listRepo:   listRepo,
		itemRepo:   itemRepo,
		userClient: userClient,
	}
}

func (s *FavoriteListService) GetUserFavoriteListsWithItems(userId int) ([]models.FavoriteListResponse, error) {

	if err := s.userClient.CheckUserId(userId); err != nil {
		return nil, models.ErrUserNotFound
	}

	lists, err := s.listRepo.GetFavoriteList(userId)

	if err != nil {
		return nil, err
	}

	var response []models.FavoriteListResponse

	for _, list := range lists {

		items, err := s.itemRepo.GetFavoriteItem(list.Id)

		if err != nil {
			return nil, err
		}

		response = append(response, models.FavoriteListResponse{
			ListId:   list.Id,
			ListName: list.ListName,
			Items:    items,
		})

	}

	return response, nil

}

func (s *FavoriteListService) CreateFavoriteList(list *models.FavoriteList) error {
	return s.listRepo.CreateFavoriteList(list)
}

func (s *FavoriteListService) UpdateFavoriteList(id int, list models.UpdateFavoriteList) (models.FavoriteList, error) {
	return s.listRepo.UpdateFavoriteList(id, list)
}

func (s *FavoriteListService) DeleteFavoriteList(id int) error {

	if err := s.itemRepo.DeleteFavoriteItemsByListId(id); err != nil {
		return err
	}

	return s.listRepo.DeleteFavoriteList(id)
}
