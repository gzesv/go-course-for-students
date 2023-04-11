package httpfiber

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"homework6/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, er := a.CreateAd(c.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)
		if er != nil {
			switch er {
			case app.ErrValidate:
				c.Status(http.StatusBadRequest)
			default:
				c.Status(http.StatusInternalServerError)
			}
			return c.JSON(AdErrorResponse(er))
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}
		return c.JSON(AdSuccessResponse(&ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, er := a.ChangeAdStatus(c.Context(), int64(adID), reqBody.UserID, reqBody.Published)
		if er != nil {
			switch er {
			case app.ErrValidate:
				c.Status(http.StatusBadRequest)
			case app.ErrAccessDenied:
				c.Status(http.StatusForbidden)
			default:
				c.Status(http.StatusInternalServerError)
			}
			return c.JSON(AdErrorResponse(er))
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(&ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, er := a.UpdateAd(c.Context(), int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if er != nil {
			switch er {
			case app.ErrValidate:
				c.Status(http.StatusBadRequest)
			case app.ErrAccessDenied:
				c.Status(http.StatusForbidden)
			default:
				c.Status(http.StatusInternalServerError)
			}
			return c.JSON(AdErrorResponse(er))
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(&ad))
	}
}
