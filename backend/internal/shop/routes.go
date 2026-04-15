package shop

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryShopItemsRoute(c *gin.Context) {
	data, err := QueryShopItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": "",
			"err": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"err":  "",
	})
}

func PurchaseItemRoute(c *gin.Context) {
	userID, ex := c.Get("userID")
	if !ex {
		c.JSON(http.StatusNotFound, gin.H{
			"data": "",
			"err":  "userID not found",
		})
		return
	}

	itemKey := c.Param("itemKey")
	item, updatedCoins, err := PurchaseItem(userID.(string), itemKey)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errMessage := "unexpected server error"
		if errors.Is(err, ErrItemNotFound) {
			statusCode = http.StatusNotFound
			errMessage = "shop item not found"
		}
		if errors.Is(err, ErrInsufficientCoins) {
			statusCode = http.StatusBadRequest
			errMessage = "not enough coins"
		}
		c.JSON(statusCode, gin.H{
			"data": "",
			"err":  errMessage,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"item":          item,
			"updated_coins": updatedCoins,
		},
		"err": "",
	})
}
