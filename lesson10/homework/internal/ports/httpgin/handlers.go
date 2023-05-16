package httpgin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"homework10/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, er := a.CreateAd(c, reqBody.Title, reqBody.Text, reqBody.UserID)
		if errors.Is(er, app.ErrWrongFormat) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(er))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		strAdID := c.Param("ad_id")
		adID, err := strconv.Atoi(strAdID)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, er := a.ChangeAdStatus(c, int64(adID), reqBody.UserID, reqBody.Published)
		if er != nil {
			if errors.Is(er, app.ErrWrongFormat) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(er))
				return
			}
			if errors.Is(er, app.ErrAccessDenied) {
				c.JSON(http.StatusForbidden, AdErrorResponse(er))
				return
			}
		}
		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		strAdID := c.Param("ad_id")
		adID, err := strconv.Atoi(strAdID)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, er := a.UpdateAd(c, int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if er != nil {
			if errors.Is(er, app.ErrWrongFormat) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(er))
				return
			}
			if errors.Is(er, app.ErrAccessDenied) {
				c.JSON(http.StatusForbidden, AdErrorResponse(er))
				return
			}
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody universalUser
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, er := a.CreateUser(c, reqBody.Nickname, reqBody.Email, reqBody.ID)

		if errors.Is(er, app.ErrWrongFormat) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(er))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&u))
	}
}

func listAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter, err := a.NewFilter(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		strAuthorID := c.Query("author_id")
		if strAuthorID != "" {
			authorID, _ := strconv.Atoi(strAuthorID)
			filter, _ = filter.FilterByAuthor(c, int64(authorID))
		}

		f, err := filter.GetFilter(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		ads, e := a.GetAllAdsByFilter(c, f)
		if e != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponseList(&ads))
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		strAdID := c.Param("ad_id")
		adID, err := strconv.Atoi(strAdID)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		var reqBody deleteAdRequest
		err = c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, er := a.DeleteAd(c, int64(adID), reqBody.UserID)
		if errors.Is(er, app.ErrWrongFormat) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(er))
			return
		}
		if errors.Is(er, app.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, AdErrorResponse(er))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ad))
	}
}

func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		strUserID := c.Param("user_id")
		userID, err := strconv.Atoi(strUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, err := a.DeleteUser(c, int64(userID))
		if errors.Is(err, app.ErrWrongFormat) {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&u))
	}
}

func changeUserInfo(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeUserStatusRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		strUserID := c.Param("user_id")
		userID, err := strconv.Atoi(strUserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		u, er := a.ChangeUserInfo(c, int64(userID), reqBody.Nickname, reqBody.Email)
		if er != nil {
			if errors.Is(er, app.ErrWrongFormat) {
				c.JSON(http.StatusBadRequest, AdErrorResponse(er))
				return
			}
			if errors.Is(er, app.ErrAccessDenied) {
				c.JSON(http.StatusForbidden, AdErrorResponse(er))
				return
			}
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&u))
	}
}

func getAdsByTitle(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Query("title")
		ads, err := a.GetAdsByTitle(c, title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponseList(&ads))
	}
}
