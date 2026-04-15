package middleware

import (
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

func CircuitBreaker(command string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hystrix.ConfigureCommand(command, hystrix.CommandConfig{
			Timeout:                5000,  // 执行超时，单位：毫秒
			MaxConcurrentRequests:  10000, // 最大并发量
			RequestVolumeThreshold: 0,     // 一个窗口内触发熔断的最小请求数
			SleepWindow:            1000,  // 熔断打开后，重新检测服务是否恢复前的等待时间，单位：毫秒
			ErrorPercentThreshold:  25,    // 触发熔断的百分比错误率
		})

		err := hystrix.Do(command, func() error {
			c.Next()
			return nil
		}, nil)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable"})
			c.Abort()
		}
	}
}

func CircuitBreakerCustom(command string, conf hystrix.CommandConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		hystrix.ConfigureCommand(command, conf)

		err := hystrix.Do(command, func() error {
			c.Next()
			return nil
		}, nil)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service Unavailable"})
			c.Abort()
		}
	}
}
