package repositories

import (
	"favorite_service/internal/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemRepository(t *testing.T) {

	favoriteItemRepository := NewFavoriteItemRepository(db)

	cFavoriteItem := models.CreateFavoriteItem{
		ItemId: 6,
		ListId: 1,
	}

	t.Run("TestCreateFavoriteItem", func(t *testing.T) {

		favoriteItem, err := favoriteItemRepository.CreateFavoriteItem(ctx, cFavoriteItem)

		assert.Nil(t, err)

		assert.NotEmpty(t, favoriteItem.CreatedDate)
	})

	t.Run("TestGetFavoriteItem", func(t *testing.T) {

		fetchedFavoriteItemList, err := favoriteItemRepository.GetFavoriteItem(ctx, cFavoriteItem.ListId)

		assert.Nil(t, err)

		check := false

		for _, item := range fetchedFavoriteItemList {

			fmt.Println("FAVORİ:", cFavoriteItem.ItemId, "ITEM:", item.ItemId)

			if item.ItemId == cFavoriteItem.ItemId && item.ListId == cFavoriteItem.ListId {

				check = true
				break

			}

		}

		assert.True(t, check)

	})

	t.Run("TestGetFavoriteItemLıstId", func(t *testing.T) {

		fetchedFavoriteItemList, err := favoriteItemRepository.GetFavoriteItem(ctx, cFavoriteItem.ListId)

		assert.Nil(t, err)

		for _, item := range fetchedFavoriteItemList {

			fmt.Println("FAVORİ:", cFavoriteItem.ItemId, "ITEM:", item.ItemId)

			assert.Equal(t, cFavoriteItem.ListId, item.ListId)

		}

	})

	t.Run("TestDeleteFavoriteItem", func(t *testing.T) {

		err := favoriteItemRepository.DeleteFavoriteItem(ctx, cFavoriteItem.ListId, cFavoriteItem.ItemId)

		assert.Nil(t, err)

		newItemList, err := favoriteItemRepository.GetFavoriteItem(ctx, cFavoriteItem.ListId)

		assert.Nil(t, err)

		for _, item := range newItemList {

			fmt.Println("FAVORİ:", cFavoriteItem.ItemId, "ITEM:", item.ItemId)

			assert.NotEqual(t, cFavoriteItem.ItemId, item.ItemId)

		}

	})

	t.Run("TestDeleteFavoriteItemsByListId", func(t *testing.T) {

		err := favoriteItemRepository.DeleteFavoriteItemsByListId(ctx, cFavoriteItem.ListId)

		assert.Nil(t, err)

	})

	t.Run("Deneme", func(t *testing.T) {
		err := favoriteItemRepository.DeleteFavoriteItemsByListId(ctx, cFavoriteItem.ListId)

		fmt.Println(models.ErrRecordNotFound, err)

		assert.Equal(t, models.ErrRecordNotFound, err)
	})

}
