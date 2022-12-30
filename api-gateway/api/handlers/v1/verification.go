package v1

import (
	"context"
	"encoding/json"
	"exam/api-gateway/api/models"
	"exam/api-gateway/genproto/customer"
	"exam/api-gateway/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
)

// Verify customer
// @Summary      Verify customer
// @Description  Verifys customer
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        email  path string true "email"
// @Param        code   path string true "code"
// @Success      200  {object}  models.VerifyResponse
// @Router      /v1/verify/{email}/{code} [get]
func (h *handlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	var (
		code  = c.Param("code")
		email = c.Param("email")
	)

	customerBody, err := h.redis.Get(email)
	if err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while getting customer from redis", logger.Any("redis", err))
	}
	if customerBody == nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{
			"info": "Your time has expired",
		})
		return
	}
	customerBodys := cast.ToString(customerBody)
	body := customer.CustomerRequest{}

	err = json.Unmarshal([]byte(customerBodys), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while unmarshaling from json to customer body", logger.Any("json", err))
		return
	}

	if body.Code != code {
		fmt.Println(body.Code)
		c.JSON(http.StatusConflict, gin.H{
			"info": "Wrong code",
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	// Genrating refresh and jwt tokens
	h.jwthandler.Iss = "customer"
	h.jwthandler.Sub = body.Id
	h.jwthandler.Role = "authorized"
	h.jwthandler.Aud = []string{"exam-app"}
	h.jwthandler.SigninKey = h.cfg.SignInKey
	h.jwthandler.Log = h.log
	tokens, err := h.jwthandler.GenerateAuthJWT()
	accessToken := tokens[0]
	refreshToken := tokens[1]

	if err != nil {
		h.log.Error("error occured while generating tokens")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong,please try again",
		})
		return
	}
	body.RefreshToken = refreshToken
	res, err := h.serviceManager.CustomerService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Error while creating customer", logger.Any("post", err))
		return
	}
	response := &models.VerifyResponse{
		Id:           res.Id,
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

	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	c.JSON(http.StatusOK, response)
}
