// File: internal/middleware/rate_limiter.go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

type RateLimiter struct {
    visitors map[string]*visitor
    mu      sync.RWMutex
    rate    rate.Limit
    burst   int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
    limiter := &RateLimiter{
        visitors: make(map[string]*visitor),
        rate:    r,
        burst:   burst,
    }

    // Start cleanup routine
    go limiter.cleanupVisitors()
    
    return limiter
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        limiter := rl.getVisitor(ip)

        if !limiter.Allow() {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (rl *RateLimiter) cleanupVisitors() {
    for {
        time.Sleep(time.Hour)

        rl.mu.Lock()
        for ip, v := range rl.visitors {
            if time.Since(v.lastSeen) > 24*time.Hour {
                delete(rl.visitors, ip)
            }
        }
        rl.mu.Unlock()
    }
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    v, exists := rl.visitors[ip]
    if !exists {
        limiter := rate.NewLimiter(rl.rate, rl.burst)
        rl.visitors[ip] = &visitor{
            limiter:  limiter,
            lastSeen: time.Now(),
        }
        return limiter
    }

    v.lastSeen = time.Now()
    return v.limiter
}

// PerRouteRateLimiter allows different rate limits for different routes
type PerRouteRateLimiter struct {
    limiters map[string]*RateLimiter
    mu       sync.RWMutex
}

func NewPerRouteRateLimiter() *PerRouteRateLimiter {
    return &PerRouteRateLimiter{
        limiters: make(map[string]*RateLimiter),
    }
}

// AddRoute adds a new rate limit for a specific route
func (prl *PerRouteRateLimiter) AddRoute(route string, r rate.Limit, burst int) {
    prl.mu.Lock()
    defer prl.mu.Unlock()

    prl.limiters[route] = NewRateLimiter(r, burst)
}

// Middleware returns a middleware function for the specified route
func (prl *PerRouteRateLimiter) Middleware(route string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            prl.mu.RLock()
            limiter, exists := prl.limiters[route]
            prl.mu.RUnlock()

            if !exists {
                next.ServeHTTP(w, r)
                return
            }

            limiter.Limit(next).ServeHTTP(w, r)
        })
    }
}
