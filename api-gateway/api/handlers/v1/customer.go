package v1

import (
	"context"
	"exam/api-gateway/api/models"
	"exam/api-gateway/genproto/customer"
	l "exam/api-gateway/pkg/logger"
	"exam/api-gateway/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// get customer
// @Summary 	this function getting the customers posts
// @Description this function select the customers posts
// @Tags 		Customer
// @Security    BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		id 		path 		string	 	true 	"user id"
// @Success 	200 	{object} 	customer.CustomerInfo
// @Failure 	400 	{object}	models.Error
// @Router 		/v1/customer/get/{id} [get]
func (h *handlerV1) GetCustomerById(c *gin.Context) {
	claims, err := GetClaims(*h, c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Error:       err,
			Description: "You are not authorized",
		})
		h.log.Error("Checking Authorization", l.Error(err))
		return
	}

	if !(claims.Role == "authorized" || claims.Role == "admin") {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Description: "You are not authorized",
		})
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseEnumNumbers = true

	userId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	body := &customer.CustomerID{
		Id: userId,
	}
	responseCus, err := h.serviceManager.CustomerService().GetCustomerInfo(ctx, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get customer", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, responseCus)
}

// Get Customer with search
// @Summary 	Customer list getting
// @Description Customer list getting
// @Tags		Customer
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param		page 		query string false "page"
// @Param 		limit		query string false "limit"
// @Param 		order		query string false "Order format should be 'key-value' key->(column_name)"
// @Param		search		query string false 	"Search format should be 'key-value'"
// @Success		200		{object}	customer.CustomerAll
// @Failure		400 	{object}	models.Error
// @Failure		500 	{object}	models.Error
// @Router		/v1/customers/{page}/{limit}/{order}/{search}	[get]
func (h *handlerV1) GetCustomerBySearchOrder(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Error:       err,
			Description: "You are not authorized",
		})
		h.log.Error("Checking Authorization", l.Error(err))
		return
	}
	if !(claims.Role == "authorized" || claims.Role == "admin" || claims.Role == "moderator") {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Description: "You are not authorized",
		})
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseEnumNumbers = true

	queryParams := c.Request.URL.Query()

	search := strings.Split(c.Query("search"), "-")
	order := strings.Split(c.Query("order"), "-")

	if len(search) != 2 && len(order) != 2 {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:        http.StatusBadRequest,
			Description: "Enter needed params",
		})
		h.log.Error("failed to get all params")
		return
	}
	params, errStr := utils.ParseQueryParams(queryParams)

	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json " + errStr[0])
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, errs := h.serviceManager.CustomerService().GetCustomerBySearchOrder(ctx, &customer.GetListUserRequest{
		Page:   params.Page,
		Limit:  params.Limit,
		Search: &customer.Search{Field: search[0], Value: search[1]},
		Orders: &customer.Order{Field: order[0], Value: order[1]},
	})
	if errs != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:        http.StatusInternalServerError,
			Description: "couldn't Get",
		})
		h.log.Error("Get Customers with search and order", l.Error(errs))
		return
	}
	c.JSON(http.StatusOK, response)
}

// user update
// @Summary 		this function update
// @Description 	this function updating the customers
// @Tags 			Customer
// @Security        BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			body    body 		models.CustomerUpdateReq true "Update Customer"
// @Success 		200 	{object} 	customer.CustomerResponse
// @Failure 		400 	{object}	models.Error
// @Router 			/v1/customer/update [put]
func (h *handlerV1) UpdateCustomer(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Error:       err,
			Description: "You are not authorized",
		})
		h.log.Error("Checking Authorization", l.Error(err))
		return
	}
	if !(claims.Role == "authorized" || claims.Role == "admin") {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Description: "You are not authorized",
		})
		return
	}

	var (
		body        models.CustomerUpdateReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseEnumNumbers = true

	err = c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	customerSaved := &customer.CustomerUpdate{
		Id:       body.Id,
		FullName: body.FullName,
		Bio:      body.Bio,
		Email:    body.Email,
	}
	_, err = h.serviceManager.CustomerService().UpdateCustomer(ctx, customerSaved)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}

	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "Profile Succesfully updated",
		StatusCode: http.StatusOK,
	})
}

// delete 		customer
// @Summary 	this function delete
// @Description this function delting customer
// @Tags 		Customer
// @Security    BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		id 	path 		string	 true "id"
// @Success 	200 {object}  	customer.CustomerResponse
// @Failure 	400 {object}	models.Error
// @Router 		/v1/customer/{id} [delete]
func (h *handlerV1) DeleteCustomer(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Error:       err,
			Description: "You are not authorized",
		})
		h.log.Error("Checking Authorization", l.Error(err))
		return
	}
	if !(claims.Role == "authorized" || claims.Role == "admin") {
		c.JSON(http.StatusUnauthorized, models.Error{
			Code:        http.StatusUnauthorized,
			Description: "You are not authorized",
		})
		return
	}

	var (
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseEnumNumbers = true

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	body := &customer.CustomerID{
		Id: id,
	}
	_, err = h.serviceManager.CustomerService().DeleteCustomer(ctx, body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
	}

	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "customer Succesfully deleted",
		StatusCode: http.StatusOK,
	})

}
