package v1

import (
	"context"
	"encoding/json"
	"exam/api-gateway/api/models"
	"exam/api-gateway/email"
	"exam/api-gateway/genproto/customer"
	"exam/api-gateway/pkg/etc"
	"exam/api-gateway/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// register customer
// @Summary		register customer
// @Description	this registers customer
// @Tags		Customer
// @Accept		json
// @Produce 	json
// @Param 		body	body  	 models.CustomerRegister true "Register customer"
// @Success		201 	{object} models.Error
// @Failure		500 	{object} models.Error
// @Router		/v1/customer/register 	[post]
func (h *handlerV1) RegisterCustomer(c *gin.Context) {
	var body models.CustomerRegister

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{
			Error: err,
		})
		h.log.Error("Error while binding json", logger.Any("json", err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)

	body.Email = strings.ToLower(body.Email)

	body.Password, err = etc.HashPassword(body.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error: err,
		})
		h.log.Error("couldn't hash the password")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	emailExists, err := h.serviceManager.CustomerService().CheckField(ctx, &customer.CheckFieldReq{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error: err,
		})
		h.log.Error("Error while cheking email uniqeness", logger.Any("check", err))
		return
	}

	if emailExists.Exists {
		c.JSON(http.StatusConflict, models.Error{
			Error:       err,
			Description: "You have already signed up",
		})
		return
	}

	exists, err := h.redis.Exists(body.Email)
	if err != nil {
		h.log.Error("Error while checking email uniqueness")
		c.JSON(http.StatusConflict, models.Error{
			Error: err,
		})
		return
	}
	if emailExists.Exists {
		c.JSON(http.StatusConflict, models.Error{
			Error: err,
		})
		return
	}

	if cast.ToInt(exists) == 1 {
		c.JSON(http.StatusConflict, models.Error{
			Error: err,
		})
		return
	}
	customerToBeSaved := &customer.CustomerRequest{
		Id:       uuid.New().String(),
		FullName: body.FullName,
		Email:    body.Email,
		Bio:      body.Bio,
		Password: body.Password,
	}
	for _, addres := range body.Addresses {
		customerToBeSaved.Addresses = append(customerToBeSaved.Addresses,
			&customer.Address{
				Id:      uuid.New().String(),
				Country: addres.Country,
				Street:  addres.Street,
			},
		)
	}
	customerToBeSaved.Code = etc.GenerateCode(6)
	msg := "Subject: Exam email verification\n Your verification code: " + customerToBeSaved.Code
	err = email.SendEmail([]string{body.Email}, []byte(msg))

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Error:       nil,
			Code:        http.StatusAccepted,
			Description: "Your Email is not valid, Please recheck it",
		})
		return
	}

	bodyByte, err := json.Marshal(customerToBeSaved)
	if err != nil {
		h.log.Error("Error while marshaling to json", logger.Any("json", err))
		return
	}

	err = h.redis.SetWithTTL(customerToBeSaved.Email, string(bodyByte), 600)
	if err != nil {
		h.log.Error("Error while marshaling to json", logger.Any("json", err))
		return
	}

	c.JSON(http.StatusAccepted, models.Error{
		Error:       nil,
		Code:        http.StatusAccepted,
		Description: "Your request successfuly accepted we have send code to your email, Your code is : " + customerToBeSaved.Code,
	})

}
