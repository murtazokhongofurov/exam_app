package v1

import (
	"context"
	"exam/api-gateway/api/models"
	ps "exam/api-gateway/genproto/post"
	l "exam/api-gateway/pkg/logger"
	"exam/api-gateway/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// Create 		post...
// @Summary	 	Create 	Post
// @Description post 	service create
// @Tags 		Post
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		post 	body  	 models.CreatePost true "Post"
// @Success 	200 	{object} post.PostResp
// @Failure 	400 	{object} models.Error
// @Router 		/v1/post [post]
func (h *handlerV1) CreatePost(c *gin.Context) {
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
		body        ps.PostReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true
	body.OwnerId = claims.Sub
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:        http.StatusBadRequest,
			Error:       err,
			Description: "Check your data",
		})
		h.log.Error("Error while binding json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	postToBeSaved := &ps.PostReq{
		Id:          uuid.New().String(),
		OwnerId:     body.OwnerId,
		Name:        body.Name,
		Description: body.Description,
	}
	postId := postToBeSaved.Id
	for _, m := range body.Medias {
		postToBeSaved.Medias = append(postToBeSaved.Medias, &ps.Media{
			Id:   uuid.New().String(),
			Name: m.Name,
			Link: m.Link,
			Type: m.Type,
		})
	}
	for _, r := range body.Reviews {
		postToBeSaved.Reviews = append(postToBeSaved.Reviews, &ps.Review{
			Id:          uuid.New().String(),
			PostId:      postId,
			OwnerId:     body.OwnerId,
			Name:        r.Name,
			Rating:      r.Rating,
			Description: r.Description,
		})
	}
	response, err := h.serviceManager.PostService().CreatePost(ctx, postToBeSaved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:        http.StatusInternalServerError,
			Error:       err,
			Description: "Check your data",
		})
		h.log.Error("Error while creating post", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// get 			post 	review
// @Summary 	Get 	Post Review
// @Description this 	will display the post review information
// @Tags 		Post
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		id 		path 	 string true "id"
// @Success 	200 	{object} post.PostInfo
// @Failure 	400 	{object} models.Error
// @Router /v1/post/get/{id} [get]
func (h *handlerV1) GetPostReview(c *gin.Context) {
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
		jspbpMarshal protojson.MarshalOptions
	)
	jspbpMarshal.UseEnumNumbers = true
	postId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	body := &ps.ID{
		PostID: postId,
	}
	responsePost, err := h.serviceManager.PostService().GetPostReview(ctx, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post review", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, responsePost)
}

// list post
// @Summary 		list 	post
// @Description 	lists 	posts
// @Tags 			Post
// @Security		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			page 	query 	string false "page"
// @Param 			limit 	query 	string false "limit"
// @Param			search	query 	string false "search"
// @Success 		200 	{object} post.ListPostResponse
// @Failure 		400 	{object} models.Error
// @Router /v1/post/list/{page}/{limit}/{search} 	[get]
func (h *handlerV1) ListPost(c *gin.Context) {
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

	var jspbpMarshal protojson.MarshalOptions
	jspbpMarshal.UseEnumNumbers = true
	queryParams := c.Request.URL.Query()

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

	response, errs := h.serviceManager.PostService().ListPost(ctx, &ps.ListPostRequest{
		Page:  params.Page,
		Limit: params.Limit,
		Value: params.Search,
	})
	if errStr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errs.Error(),
		})
		h.log.Error("failed to get post review", l.Error(errs))
		return
	}

	c.JSON(http.StatusOK, response)
}

// update 			post
// @Summary 		Update 		Post
// @Description 	this 		updating information of post
// @Tags 			Post
// @Security		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			postbody   body     models.PostUpdate true "Post"
// @Success 		200 	   {object} post.PostResp
// @Failure 		400 	   {boject} models.Error
// @Failure 		500 	   {object} models.Error
// @Router 			/v1/post/update [put]
func (h *handlerV1) UpdatePost(c *gin.Context) {
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
		postbody     models.PostUpdate
		jspbpMarshal protojson.MarshalOptions
	)
	jspbpMarshal.UseEnumNumbers = true
	err = c.ShouldBindJSON(&postbody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update post", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	postSaved := &ps.PostUp{
		Id:          postbody.Id,
		Name:        postbody.Name,
		Description: postbody.Description,
	}
	for _, m := range postbody.Medias {
		postSaved.Medias = append(postSaved.Medias, &ps.Media{
			Id:   m.Id,
			Name: m.Name,
			Link: m.Link,
			Type: m.Type,
		})
	}
	_, err = h.serviceManager.PostService().UpdatePost(ctx, postSaved)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update post service", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "Post Succesfully updated",
		StatusCode: http.StatusOK,
	})
}

// delete 			post
// @Summary 		Delete 	Post
// @Description 	this 	deleting information of post
// @Tags 			Post
// @Security		BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id 		path string true "id"
// @Success 		200 	{object} post.Empty
// @Failure 		400 	{object} models.Error
// @Router 			/v1/post/delete/{id} [delete]
func (h *handlerV1) DeletePost(c *gin.Context) {
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
		jspbpMarshal protojson.MarshalOptions
	)
	jspbpMarshal.UseEnumNumbers = true

	postId := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	body := &ps.ID{
		PostID: postId,
	}
	response, err := h.serviceManager.PostService().DeletePost(ctx, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete post service")
	}
	c.JSON(http.StatusOK, response)
}
