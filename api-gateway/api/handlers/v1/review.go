package v1

import (
	"context"
	"exam/api-gateway/api/models"
	"exam/api-gateway/genproto/review"
	l "exam/api-gateway/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary 		Get Review
// @Description 	this will display the review information
// @Tags 			Review
// @Security		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path 		string true "id"
// @Success 		200 	{object} 	review.Reviews
// @Failure 		400 	{object}	models.Error
// @Router 			/v1/review/get/{id} [get]
func (h *handlerV1) GetReviewById(c *gin.Context) {
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

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	body := &review.ReviewId{
		Id: id,
	}

	respose, err := h.serviceManager.ReviewService().GetReviewById(ctx, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error id of review")
	}
	c.JSON(http.StatusOK, respose)
}

// update review
// @Summary 		Update Review
// @Description 	this updating information of review
// @Tags 			Review
// @Security		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			body 	body 	 	models.ReviewUpdate true "Review"
// @Success 		200 	{object} 	review.ReviewResponse
// @Failure 		400 	{object}	models.Error
// @Failure 		500 	{object}	models.Error
// @Router 			/v1/review/update 	[put]
func (h *handlerV1) UpdateReview(c *gin.Context) {
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
	var (
		body        models.ReviewUpdate
		jsrsMarshal protojson.MarshalOptions
	)
	jsrsMarshal.UseEnumNumbers = true
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	reviewSaved := &review.ReviewUp{
		Id:          body.Id,
		Name:        body.Name,
		Rating:      body.Rating,
		Description: body.Description,
	}

	res, err := h.serviceManager.ReviewService().UpdateReview(ctx, reviewSaved)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error update review", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, res)
}

// delete review
// @Summary 		Delete Review
// @Description 	this deleting information of review
// @Tags 			Review
// @Security		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 	path 		string true "id"
// @Success 		200 {object} 	review.Empty
// @Failure 		400 {object}	models.Error
// @Router 			/v1/review/delete/{id} [delete]
func (h *handlerV1) DeleteReview(c *gin.Context) {
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
	var jsrsMarshal protojson.MarshalOptions
	jsrsMarshal.UseEnumNumbers = true

	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	body := &review.ReviewId{
		Id: id,
	}
	_, err = h.serviceManager.ReviewService().DeleteReview(ctx, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while deleting review")
	}
	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "Review Succesfully deleted",
		StatusCode: http.StatusOK,
	})

}
