package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

func AppRouter(r gin.IRoutes, a app.App) {
	r.POST("/ads", createAd(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a))
	r.PUT("/ads/:ad_id", updateAd(a))
	r.GET("/ads", listAds(a))
	r.POST("/users", createUser(a))
	r.PUT("/users/:user_id", changeUserInfo(a))
	r.GET("/ads/by_title", getAdsByTitle(a))
	r.DELETE("/ads/:ad_id", deleteAd(a))
	r.DELETE("/users/:user_id", deleteUser(a))
}
