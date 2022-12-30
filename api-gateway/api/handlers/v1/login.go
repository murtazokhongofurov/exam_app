package v1

import (
	"context"
	"exam/api-gateway/api/models"
	"exam/api-gateway/genproto/customer"
	"exam/api-gateway/pkg/etc"
	"exam/api-gateway/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// customer login
// @Summary 		Login Customer
// @Description 	This function get login customer
// @Tags 			Customer
// @Accept 			json
// @Produce			json
// @Param 			email 		path string true "email"
// @Param 			password 	path string true "password"
// @Success 		200 {object} 	customer.LoginResponse
// @Failure			500 {object} 	models.Error
// @Failure			400 {object} 	models.Error
// @Router			/v1/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		email    = c.Param("email")
		password = c.Param("password")
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.CustomerService().GetByEmail(ctx, &customer.EmailReq{
		Email: email,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Error:       err,
			Description: "Couln't find matching information, Have you registered before?",
		})
		h.log.Error("Error while getting customer by email", logger.Any("post", err))
		return
	}

	if !etc.CheckPasswordHash(password, res.Password) {
		c.JSON(http.StatusNotFound, models.Error{
			Description: "Password or Email error",
			Code:        http.StatusBadRequest,
		})
		return
	}

	h.jwthandler.Iss = "customer"
	h.jwthandler.Sub = res.Id
	h.jwthandler.Role = "authorized"
	h.jwthandler.Aud = []string{"exam-app"}
	h.jwthandler.SigninKey = h.cfg.SignInKey
	h.jwthandler.Log = h.log
	tokens, err := h.jwthandler.GenerateAuthJWT()
	accesstoken := tokens[0]
	refreshToken := tokens[1]
	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}

	res.AccessToken = accesstoken
	res.RefreshToken = refreshToken
	res.Password = ""
	
	response := models.CustomerLogin{
		FullName:     res.FullName,
		Bio:          res.Bio,
		Email:        res.Email,
		Password:     res.Password,
		RefreshToken: res.RefreshToken,
		AccessToken:  res.AccessToken,
	}
	for _, add := range res.Addresses {
		response.Addresses = append(response.Addresses, models.Address{
			Id:      add.Id,
			OwnerId: add.OwnerId,
			Country: add.Country,
			Street:  add.Street,
		})
	}

	c.JSON(http.StatusOK, response)
}
