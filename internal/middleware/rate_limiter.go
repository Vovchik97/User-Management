package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
	"userManagement/internal/dto"
	"userManagement/internal/utils"
)

type clientData struct {
	RequestCount int
	ResetTime    time.Time
}

var clients sync.Map

const (
	LimitRequestsPerMinute = 60
)

// RateLimiter — middleware для ограничения количества запросов
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		now := time.Now()
		value, _ := clients.LoadOrStore(ip, &clientData{
			RequestCount: 0,
			ResetTime:    now.Add(time.Minute),
		})
		client := value.(*clientData)

		if now.After(client.ResetTime) {
			client.RequestCount = 0
			client.ResetTime = now.Add(time.Minute)
		}

		client.RequestCount++

		if client.RequestCount > LimitRequestsPerMinute {
			utils.Log.Warnf("IP %s превысил лимит запросов (%d в минуту)", ip, LimitRequestsPerMinute)
			c.JSON(http.StatusTooManyRequests, dto.ResponseError{
				Message: "Слишком много запросов. Попробуйте позже",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
