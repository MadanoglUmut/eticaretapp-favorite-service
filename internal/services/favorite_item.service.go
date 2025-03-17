package services

import (
	"favorite_service/internal/models"
	"sync"
)

type favoriItemRepository interface {
	GetFavoriteItem(listId int) ([]models.FavoriteItem, error)
	CreateFavoriteItem(favoriteItem models.CreateFavoriteItem) (models.FavoriteItem, error)
	DeleteFavoriteItem(listId int, itemId int) error
}

type listRepository interface {
	GetListOwner(listId int) (models.FavoriteList, error)
}

type favoriteItemProductClient interface {
	VerifyProduct(productId int) (*models.Product, error)
}

type FavoriItemService struct {
	favoriItemRepository favoriItemRepository
	listRepository       listRepository
	productClient        favoriteItemProductClient
}

func NewFavoriItemService(favoriItemRepository favoriItemRepository, lislistRepository listRepository,
	productClient favoriteItemProductClient) *FavoriItemService {
	return &FavoriItemService{
		favoriItemRepository: favoriItemRepository,
		listRepository:       lislistRepository,
		productClient:        productClient,
	}
}

func (s *FavoriItemService) GetFavoriteItem(listId int, userId int) ([]models.Product, error) {

	ownerFavoriteList, err := s.listRepository.GetListOwner(listId)

	if err != nil {
		return nil, err
	}

	if ownerFavoriteList.UserId != userId {

		return nil, models.ErrunaUthorizedAction

	}

	favoriteItems, err := s.favoriItemRepository.GetFavoriteItem(listId)
	if err != nil {
		return nil, err
	}

	return GetProductInfo(favoriteItems, s.productClient)

}

func (s *FavoriItemService) CreateFavoriteItem(item models.CreateFavoriteItem, userId int) (models.FavoriteItem, error) {

	ownerFavoriteList, err := s.listRepository.GetListOwner(item.ListId)

	if err != nil {
		return models.FavoriteItem{}, err
	}

	if ownerFavoriteList.UserId != userId {
		return models.FavoriteItem{}, models.ErrunaUthorizedAction
	}

	return s.favoriItemRepository.CreateFavoriteItem(item)
}

func (s *FavoriItemService) DeleteFavoriteItem(listId int, itemId int) error {
	return s.favoriItemRepository.DeleteFavoriteItem(listId, itemId)
}

func GetProductInfo(items []models.FavoriteItem, productClinet favoriteItemProductClient) ([]models.Product, error) {

	var products []models.Product
	var wg sync.WaitGroup
	productChan := make(chan *models.Product, 2)
	errChan := make(chan error, 2)
	sm := make(chan struct{}, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, item := range items {
			wg.Add(1)
			sm <- struct{}{}
			go func(itemId int) {
				defer wg.Done()
				defer func() { <-sm }()
				product, err := productClinet.VerifyProduct(item.ItemId)
				if err != nil {
					errChan <- err
					return
				}
				productChan <- product
			}(item.ItemId)

			if len(errChan) > 0 {
				break

			}
		}

	}()

	go func() {
		wg.Wait()
		close(errChan)
		close(productChan)

	}()

	for product := range productChan {
		products = append(products, *product)
	}

	for err := range errChan {
		return nil, err
	}

	return products, nil

}
