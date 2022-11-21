package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"im/models"
	"im/pkg/database"
	"im/pkg/jwt"
	"im/pkg/resp"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
)

func NewJwtCheckMiddleware(
	jwt *jwt.Jwt,
	mysqlClient *database.Client,
	cacheClient *cache.Cache,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id uint = 0
		token := c.GetHeader("Authorization")
		if token == "" {
			token, _ = c.GetQuery("t")
		}
		if token != "" {
			id = jwt.ParseToken(token)
		}
		if id == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.Response{
				Code:    resp.ERROR,
				Message: resp.TIMEOUT,
			})
			return
		}

		// 查询并缓存账号
		account := &models.Account{}
		err := cacheClient.Once(&cache.Item{
			Ctx:   c.Request.Context(),
			Key:   fmt.Sprintf("account:%d", id),
			Value: account,
			TTL:   time.Second * 3,
			Do: func(i *cache.Item) (interface{}, error) {
				account := &models.Account{}
				err := mysqlClient.Db().
					Model(&models.Account{}).
					Where("id = ?", id).
					First(account).
					Error
				return account, err
			},
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.Response{
				Code:    resp.ERROR,
				Message: err.Error(),
			})
			return
		}
		if account.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.Response{
				Code:    resp.ERROR,
				Message: resp.TIMEOUT,
			})
			return
		}

		// Context存储用户id
		c.Set("id", id)

		// Token刷新
		newlyToken, _ := jwt.BuildToken(
			id,
			models.LoginExpired,
		)
		c.Header("Authorization", newlyToken)
		c.Header("Access-Control-Expose-Headers", "Authorization")

		c.Next()
	}
}
