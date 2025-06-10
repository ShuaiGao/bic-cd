package gin_recover

import (
	"bic-cd/pkg/gen/api"
	"bic-cd/pkg/log"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack, _ := httputil.DumpRequest(c.Request, false)
				// stack := stack(3) //gin官方recovery
				reset := string([]byte{27, 91, 48, 109}) //重置颜色
				log.Errorf(c, "[Recovery] panic recovered:\n%s\n\n%s%s", err, stack, reset)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":   api.ECServerError.Code(),
					"detail": api.ECServerError.String(),
				})
			}
		}()
		c.Next()
	}
}
