package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipRateLimiter struct {
	mu         sync.Mutex
	visitors   map[string]*visitor
	rate       rate.Limit
	burst      int
	entryTTL   time.Duration
	identifier func(*gin.Context) string
}

func newIPRateLimiter(rpm int, burst int, entryTTL time.Duration, identifier func(*gin.Context) string) *ipRateLimiter {
	if rpm <= 0 {
		rpm = 60
	}
	if burst <= 0 {
		burst = 10
	}
	if entryTTL <= 0 {
		entryTTL = 10 * time.Minute
	}
	if identifier == nil {
		identifier = func(c *gin.Context) string { return c.ClientIP() }
	}

	rl := &ipRateLimiter{
		visitors:   make(map[string]*visitor),
		rate:       rate.Limit(float64(rpm) / 60.0),
		burst:      burst,
		entryTTL:   entryTTL,
		identifier: identifier,
	}

	go rl.cleanupLoop()
	return rl
}

func (l *ipRateLimiter) getLimiter(key string) *rate.Limiter {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if v, ok := l.visitors[key]; ok {
		v.lastSeen = now
		return v.limiter
	}

	limiter := rate.NewLimiter(l.rate, l.burst)
	l.visitors[key] = &visitor{limiter: limiter, lastSeen: now}
	return limiter
}

func (l *ipRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-l.entryTTL)

		l.mu.Lock()
		for ip, v := range l.visitors {
			if v.lastSeen.Before(cutoff) {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

func (l *ipRateLimiter) middleware(limitName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.identifier(c)
		if key == "" {
			key = "unknown"
		}

		if !l.getLimiter(key).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests",
				"error":   limitName + " rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}

var (
	apiRateLimitOnce  sync.Once
	apiRateLimitMW    gin.HandlerFunc
	authRateLimitOnce sync.Once
	authRateLimitMW   gin.HandlerFunc
)

func APIRateLimit() gin.HandlerFunc {
	apiRateLimitOnce.Do(func() {
		rpm := envInt("RATE_LIMIT_API_RPM", 120)
		burst := envInt("RATE_LIMIT_API_BURST", 30)
		enabled := envBool("RATE_LIMIT_ENABLED", true)

		if !enabled {
			apiRateLimitMW = func(c *gin.Context) { c.Next() }
			return
		}

		rl := newIPRateLimiter(rpm, burst, 10*time.Minute, nil)
		apiRateLimitMW = rl.middleware("api")
	})

	return apiRateLimitMW
}

func AuthRateLimit() gin.HandlerFunc {
	authRateLimitOnce.Do(func() {
		rpm := envInt("RATE_LIMIT_AUTH_RPM", 20)
		burst := envInt("RATE_LIMIT_AUTH_BURST", 5)
		enabled := envBool("RATE_LIMIT_ENABLED", true)

		if !enabled {
			authRateLimitMW = func(c *gin.Context) { c.Next() }
			return
		}

		rl := newIPRateLimiter(rpm, burst, 20*time.Minute, nil)
		authRateLimitMW = rl.middleware("auth")
	})

	return authRateLimitMW
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
