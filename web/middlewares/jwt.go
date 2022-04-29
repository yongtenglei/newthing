package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/pkg/jwtx"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取token
		token := c.Request.Header.Get("token")
		if token == "" || len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": e.InvalidTokenErr,
			})
			c.Abort()
			return
		}

		// 解析token

		parsedToken, err := jwtx.ParseUserClaims(token)
		if err != nil {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"msg": err.Error(),
				//"msg": e.InvalidTokenErr,
			})
			c.Abort()
			return
		}

		// 解析成功后判断是否需要被刷新 ???
		//now := time.Now().Unix()
		//expiresAt := parsedToken.ExpiresAt
		//oneday := int64(86400) //  60 * 60 * 24

		//if expiresAt > now && expiresAt-now < oneday {
		//j.RefreshToken(token)
		//}

		c.Set("claims", parsedToken)
		c.Next()
	}
}
