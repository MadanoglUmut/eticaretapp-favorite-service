package services

import "favorite_service/internal/models"

type favoriItemRepository interface {
	GetFavoriteItem(listId int) ([]models.FavoriteItem, error)
	CreateFavoriteItem(favoriteItem models.CreateFavoriteItem) (models.FavoriteItem, error)
	DeleteFavoriteItem(listId int, itemId int) error
}

type FavoriItemService struct {
	favoriItemRepository favoriItemRepository
}

func NewFavoriItemService(favoriItemRepository favoriItemRepository) *FavoriItemService {
	return &FavoriItemService{
		favoriItemRepository: favoriItemRepository,
	}
}

func (s *FavoriItemService) GetFavoriteItem(listId int) ([]models.FavoriteItem, error) {
	return s.favoriItemRepository.GetFavoriteItem(listId)
}

func (s *FavoriItemService) CreateFavoriteItem(item models.CreateFavoriteItem) (models.FavoriteItem, error) {
	return s.favoriItemRepository.CreateFavoriteItem(item)
}

func (s *FavoriItemService) DeleteFavoriteItem(listId int, itemId int) error {
	return s.favoriItemRepository.DeleteFavoriteItem(listId, itemId)
}
