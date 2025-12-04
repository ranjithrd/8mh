package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func RateLimiter(requestsPerMinute int) echo.MiddlewareFunc {
	go cleanupVisitors()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()

			mu.Lock()
			v, exists := visitors[ip]
			if !exists {
				limiter := rate.NewLimiter(rate.Limit(requestsPerMinute)/60, requestsPerMinute)
				visitors[ip] = &visitor{limiter, time.Now()}
				v = visitors[ip]
			}
			v.lastSeen = time.Now()
			mu.Unlock()

			if !v.limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceeded",
				})
			}

			return next(c)
		}
	}
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
