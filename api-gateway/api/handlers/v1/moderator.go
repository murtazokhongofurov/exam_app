package v1

import (
	"context"
	"exam/api-gateway/api/models"
	"exam/api-gateway/genproto/customer"
	"exam/api-gateway/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// Login admin
// @Summary			Login moderator
// @Description		Login moderator
// @Tags			Admins
// @Accept			json
// @Produce			json
// @Param			name		path string true "name"
// @Param 			password 	path string true "password"
// @Success			200 		{object} 	customer.GetModeratorRes
// @Failure			400			{object}	models.Error
// @Failure			500			{object}	models.Error
// @Failure			404			{object}	models.Error
// @Failure			409			{object}	models.Error
// @Router			/v1/moderator/login/{name}/{password} [get]
func (h *handlerV1) LoginModerator(c *gin.Context) {
	var (
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseEnumNumbers = true

	var (
		name     = c.Param("name")
		password = c.Param("password")
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.CustomerService().GetModerator(ctx, &customer.GetModeratorReq{Name: name})

	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Code:        http.StatusNotFound,
			Error:       err,
			Description: "Couln't find matching information, Have you registered before?",
		})
		h.log.Error("Error while getting admin by admin Name", logger.Any("Get", err))
		return
	}
	// if !etc.CheckPasswordHash(password, res.Password) {
	// 	c.JSON(http.StatusConflict, models.Error{
	// 		Description: "Password or adminName error",
	// 		Code:        http.StatusConflict,
	// 	})
	// 	return
	// }

	if res.Password != password {
		fmt.Println("error password")
		return
	}

	h.jwthandler.Iss = "moderator"
	h.jwthandler.Sub = res.Id
	h.jwthandler.Role = "moderator"
	h.jwthandler.Aud = []string{"exam-app"}
	h.jwthandler.SigninKey = h.cfg.SignInKey
	h.jwthandler.Log = h.log
	tokens, err := h.jwthandler.GenerateAuthJWT()
	accessToken := tokens[0]

	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	res.AccessToken = accessToken
	res.Password = ""

	c.JSON(http.StatusOK, res)
}
