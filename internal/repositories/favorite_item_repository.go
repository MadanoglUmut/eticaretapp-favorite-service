package repositories

import (
	"favorite_service/internal/models"

	"gorm.io/gorm"
)

type FavoriteItemRepository struct {
	db *gorm.DB
}

func NewFavoriteItemRepository(db *gorm.DB) *FavoriteItemRepository {
	return &FavoriteItemRepository{
		db: db,
	}

}

func (r *FavoriteItemRepository) GetFavoriteItem(listId int) ([]models.FavoriteItem, error) {

	var favoriteItems []models.FavoriteItem

	if err := r.db.Table("favoriteitem").Debug().Where("listid = ?", listId).Find(&favoriteItems).Error; err != nil {
		return nil, err
	}

	return favoriteItems, nil

}

func (r *FavoriteItemRepository) CreateFavoriteItem(favoriteItem models.CreateFavoriteItem) (models.FavoriteItem, error) {

	newFavoriteItem := models.FavoriteItem{
		ItemId: favoriteItem.ItemId,
		ListId: favoriteItem.ListId,
	}

	if err := r.db.Table("favoriteitem").Create(&newFavoriteItem).Error; err != nil {
		return models.FavoriteItem{}, err
	}

	return newFavoriteItem, nil

}

func (r *FavoriteItemRepository) DeleteFavoriteItem(listId int, itemId int) error {

	result := r.db.Table("favoriteitem").Where("listid = ? AND itemid = ?", listId, itemId).Delete(&models.FavoriteItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil

}

func (r *FavoriteItemRepository) DeleteFavoriteItemsByListId(listId int) error {

	result := r.db.Table("favoriteitem").Where("listid = ?", listId).Delete(&models.FavoriteItem{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil

}
