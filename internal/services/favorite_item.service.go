package services

import (
	"favorite_service/internal/models"
)

type favoriItemRepository interface {
	GetFavoriteItem(listId int) ([]models.FavoriteItem, error)
	CreateFavoriteItem(favoriteItem models.CreateFavoriteItem) (models.FavoriteItem, error)
	DeleteFavoriteItem(listId int, itemId int) error
}

type listRepository interface {
	GetListOwner(listId int) (int, error)
}

type FavoriItemService struct {
	favoriItemRepository favoriItemRepository
	listRepository       listRepository
}

func NewFavoriItemService(favoriItemRepository favoriItemRepository, lislistRepository listRepository) *FavoriItemService {
	return &FavoriItemService{
		favoriItemRepository: favoriItemRepository,
		listRepository:       lislistRepository,
	}
}

func (s *FavoriItemService) GetFavoriteItem(listId int, userId int) ([]models.FavoriteItem, error) {

	ownerId, err := s.listRepository.GetListOwner(listId)

	if err != nil {
		return nil, err
	}

	if ownerId != userId {

		return nil, models.ErrunaUthorizedAction

	}

	return s.favoriItemRepository.GetFavoriteItem(listId)
}

func (s *FavoriItemService) CreateFavoriteItem(item models.CreateFavoriteItem, userId int) (models.FavoriteItem, error) {

	ownerId, err := s.listRepository.GetListOwner(item.ListId)

	if err != nil {
		return models.FavoriteItem{}, err
	}

	if ownerId != userId {
		return models.FavoriteItem{}, models.ErrunaUthorizedAction
	}

	return s.favoriItemRepository.CreateFavoriteItem(item)
}

func (s *FavoriItemService) DeleteFavoriteItem(listId int, itemId int) error {
	return s.favoriItemRepository.DeleteFavoriteItem(listId, itemId)
}
