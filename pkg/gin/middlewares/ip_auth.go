package middlewares

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

var defaultRanges = []string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"::1/128",
	"fc00::/7",
}

// NewIPAuthMiddleware ...
func NewIPAuthMiddleware(ipRanges []string) gin.HandlerFunc {
	if len(ipRanges) == 0 {
		ipRanges = defaultRanges
	}

	var masks []*net.IPNet
	for _, r := range ipRanges {
		_, mask, err := net.ParseCIDR(r)
		if err != nil {
			panic(err)
		}
		masks = append(masks, mask)
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if clientIP == "" {
			c.Next()
			return
		}

		ip := net.ParseIP(clientIP)
		contains := false
		for _, m := range masks {
			if m.Contains(ip) {
				contains = true
			}
		}

		if !contains {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized ip addr",
			})
			return
		}

		c.Next()
	}
}
