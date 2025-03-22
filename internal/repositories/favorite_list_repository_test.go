package repositories

import (
	"favorite_service/internal/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRepository(t *testing.T) {

	favoriteListRepository := NewFavoriteListRepository(db)

	cFavoriteList := models.FavoriteList{
		ListName: "Deneme",
		UserId:   1,
	}

	t.Run("TestCreateFavoriteList", func(t *testing.T) {

		err := favoriteListRepository.CreateFavoriteList(ctx, &cFavoriteList)

		assert.Nil(t, err)

		var favoriteList models.FavoriteList

		db.Table("favoritelist").Find(&favoriteList, cFavoriteList.Id)

		assert.Equal(t, cFavoriteList.ListName, favoriteList.ListName)

	})

	t.Run("TestGetFavoriteListIf", func(t *testing.T) {

		fetchedFavoriteItemList, err := favoriteListRepository.GetFavoriteList(ctx, cFavoriteList.UserId)

		assert.Nil(t, err)

		check := false

		for _, item := range fetchedFavoriteItemList {

			fmt.Println("FAVORİ:", cFavoriteList.ListName, "ITEM:", item.ListName)

			if item.Id == cFavoriteList.Id && item.UserId == cFavoriteList.UserId {

				check = true
				break

			}

		}

		assert.True(t, check)

	})

	t.Run("TestGetFavoriteList", func(t *testing.T) {

		fetchedFavoriteLists, err := favoriteListRepository.GetFavoriteList(ctx, cFavoriteList.UserId)

		assert.Nil(t, err)

		for _, list := range fetchedFavoriteLists {

			fmt.Println(cFavoriteList.ListName, list.ListName)

			assert.Equal(t, cFavoriteList.UserId, list.UserId)

		}

	})

	t.Run("TestUpdateFavoriteList", func(t *testing.T) {

		updatedList := models.UpdateFavoriteList{
			ListName: "UpdateDeneme",
		}

		updatedFavoriteList, err := favoriteListRepository.UpdateFavoriteList(ctx, cFavoriteList.Id, updatedList)

		assert.Nil(t, err)

		fmt.Println("UpdateDeneme", updatedFavoriteList.ListName)

		assert.Equal(t, updatedList.ListName, updatedFavoriteList.ListName)

	})

	t.Run("TestDeleteFavoriteList", func(t *testing.T) {
		err := favoriteListRepository.DeleteFavoriteList(ctx, cFavoriteList.Id)

		assert.Nil(t, err)

		newFavoriteLists, err := favoriteListRepository.GetFavoriteList(ctx, cFavoriteList.Id)

		assert.Nil(t, err)

		for _, list := range newFavoriteLists {

			fmt.Println(cFavoriteList.Id, list.Id)

			assert.NotEqual(t, cFavoriteList.Id, list.Id)
		}

	})

	t.Run("TestDeleteFavoriteListAgain", func(t *testing.T) {
		err := favoriteListRepository.DeleteFavoriteList(ctx, cFavoriteList.Id)

		fmt.Println(models.ErrRecordNotFound, err)

		assert.Equal(t, models.ErrRecordNotFound, err)

		err = favoriteListRepository.DeleteFavoriteList(ctx, 888)

		fmt.Println(models.ErrRecordNotFound, err)

		assert.Equal(t, models.ErrRecordNotFound, err)

	})

}
