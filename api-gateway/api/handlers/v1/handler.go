package v1

import (
	t "exam/api-gateway/api/tokens"
	"exam/api-gateway/config"
	"exam/api-gateway/pkg/logger"
	"exam/api-gateway/services"
	"exam/api-gateway/storage/repo"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redis          repo.RedisRepo
	jwthandler     t.JWTHandler
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RedisRepo
	JWTHandler     t.JWTHandler
}

func New(c *HandlerV1Config) handlerV1 {
	return handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
		jwthandler:     c.JWTHandler,
	}
}

func GetClaims(h handlerV1, c *gin.Context) (*t.CustomClaims, error) {
	var (
		claims = t.CustomClaims{}
	)
	strToken := c.GetHeader("Authorization")
	token, err := jwt.Parse(strToken, func(t *jwt.Token) (interface{}, error) { return []byte(h.cfg.SignInKey), nil })

	if err != nil {
		fmt.Println(err)
		h.log.Error("invalid access token")
		return nil, err
	}

	rawClaims := token.Claims.(jwt.MapClaims)

	claims.Sub = rawClaims["sub"].(string)
	claims.Exp = rawClaims["exp"].(float64)
	aud := cast.ToStringSlice(rawClaims["aud"])
	claims.Aud = aud
	claims.Role = rawClaims["role"].(string)
	claims.Sub = rawClaims["sub"].(string)
	claims.Token = token
	return &claims, nil
}
