package repositories

import (
	"favorite_service/internal/models"

	"gorm.io/gorm"
)

type FavoriteListRepository struct {
	db *gorm.DB
}

func NewFavoriteListRepository(db *gorm.DB) *FavoriteListRepository {
	return &FavoriteListRepository{
		db: db,
	}
}

func (r *FavoriteListRepository) GetFavoriteList(userId int) ([]models.FavoriteList, error) {
	var favoriteList []models.FavoriteList

	if err := r.db.Table("favoritelist").Debug().Where("userid = ?", userId).Find(&favoriteList).Error; err != nil {
		return nil, err
	}

	return favoriteList, nil

}

func (r *FavoriteListRepository) CreateFavoriteList(favoriteList *models.FavoriteList) error {

	if err := r.db.Debug().Table("favoritelist").Create(&favoriteList).Error; err != nil {
		return err

	}
	return nil
}

func (r *FavoriteListRepository) UpdateFavoriteList(listId int, updtList models.UpdateFavoriteList) (models.FavoriteList, error) {

	var favoriList models.FavoriteList
	if err := r.db.Table("favoritelist").First(&favoriList, listId).Error; err != nil {
		return models.FavoriteList{}, err
	}

	favoriList.ListName = updtList.ListName

	if err := r.db.Table("favoritelist").Save(&favoriList).Error; err != nil {
		return models.FavoriteList{}, err
	}

	return favoriList, nil

}

func (r *FavoriteListRepository) DeleteFavoriteList(listId int) error {

	result := r.db.Table("favoritelist").Where("id = ?", listId).Delete(&models.FavoriteList{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil

}

//BU FONKSİYON TEK DEĞER DÖNMESİN

func (r *FavoriteListRepository) GetListOwner(listId int) (int, error) {

	var favoriteList models.FavoriteList

	if err := r.db.Table("favoritelist").Where("id = ?", listId).First(&favoriteList).Error; err != nil {

		return 0, err

	}

	return favoriteList.UserId, nil

}
