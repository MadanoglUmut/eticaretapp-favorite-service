package services

import (
	"context"
	"favorite_service/internal/models"
	"sync"
)

type favoriItemRepository interface {
	GetFavoriteItem(ctx context.Context, listId int) ([]models.FavoriteItem, error)
	CreateFavoriteItem(ctx context.Context, favoriteItem models.CreateFavoriteItem) (models.FavoriteItem, error)
	DeleteFavoriteItem(ctx context.Context, listId int, itemId int) error
}

type listRepository interface {
	GetListOwner(ctx context.Context, listId int) (models.FavoriteList, error)
}

type favoriteItemProductClient interface {
	VerifyProduct(ctx context.Context, productId int) (*models.Product, error)
}

type favoriteItemUserClient interface {
	VerifyUser(token string, ctx context.Context) (*models.Users, error)
}

type FavoriItemService struct {
	favoriItemRepository favoriItemRepository
	listRepository       listRepository
	productClient        favoriteItemProductClient
	userClient           favoriteItemUserClient
}

func NewFavoriItemService(favoriItemRepository favoriItemRepository, lislistRepository listRepository,
	productClient favoriteItemProductClient, userClient favoriteItemUserClient) *FavoriItemService {
	return &FavoriItemService{
		favoriItemRepository: favoriItemRepository,
		listRepository:       lislistRepository,
		productClient:        productClient,
		userClient:           userClient,
	}
}

func (s *FavoriItemService) GetFavoriteItem(listId int, token string, ctx context.Context) ([]models.Product, error) {

	user, err := s.userClient.VerifyUser(token, ctx)

	if err != nil {
		return nil, err
	}

	ownerFavoriteList, err := s.listRepository.GetListOwner(ctx, listId)

	if err != nil {
		return nil, err
	}

	if ownerFavoriteList.UserId != user.ID {

		return nil, models.ErrunaUthorizedAction

	}

	favoriteItems, err := s.favoriItemRepository.GetFavoriteItem(ctx, listId)
	if err != nil {
		return nil, err
	}

	return GetProductInfo(ctx, favoriteItems, s.productClient)

}

func (s *FavoriItemService) CreateFavoriteItem(item models.CreateFavoriteItem, token string, ctx context.Context) (models.FavoriteItem, error) {

	user, err := s.userClient.VerifyUser(token, ctx)

	if err != nil {

		return models.FavoriteItem{}, err

	}

	ownerFavoriteList, err := s.listRepository.GetListOwner(ctx, item.ListId)

	if err != nil {
		return models.FavoriteItem{}, err

	}

	if ownerFavoriteList.UserId != user.ID {
		return models.FavoriteItem{}, models.ErrunaUthorizedAction

	}

	return s.favoriItemRepository.CreateFavoriteItem(ctx, item)
}

func (s *FavoriItemService) DeleteFavoriteItem(listId int, itemId int, token string, ctx context.Context) error {

	user, err := s.userClient.VerifyUser(token, ctx)

	if err != nil {
		return err
	}

	ownerFavoriteList, err := s.listRepository.GetListOwner(ctx, listId)

	if err != nil {
		return err
	}

	if ownerFavoriteList.UserId != user.ID {

		return models.ErrunaUthorizedAction

	}

	return s.favoriItemRepository.DeleteFavoriteItem(ctx, listId, itemId)
}

func GetProductInfo(ctx context.Context, items []models.FavoriteItem, productClient favoriteItemProductClient) ([]models.Product, error) {
	var products []models.Product
	var wg sync.WaitGroup
	productChan := make(chan *models.Product, 2)
	errChan := make(chan error, 2)
	sm := make(chan struct{}, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, item := range items {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				wg.Add(1)
				sm <- struct{}{}
				go func(itemId int) {
					defer wg.Done()
					defer func() { <-sm }()
					product, err := productClient.VerifyProduct(ctx, item.ItemId)
					if err != nil {
						errChan <- err
						return
					}
					productChan <- product
				}(item.ItemId)
			}

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
