package main

import (
	"go-simple-osm-tile-cache-server/internal/config"
	"go-simple-osm-tile-cache-server/internal/db"
	"go-simple-osm-tile-cache-server/internal/tile"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limit *rate.Limiter

type body struct {
	Detail string `json:"detail"`
}

func main() {
	rateLimit := config.Get().SERVER.RATE_LIMIT
	limit = rate.NewLimiter(rate.Every(time.Second), rateLimit)
	//
	route := gin.Default()
	// Set up rate limiter if rate limit is over 0
	if rateLimit > 0 {
		route.Use(rateLimiter())
	}
	route.GET("/tile/:z/:x/:y", getTile)
	gin.SetMode(gin.ReleaseMode)
	route.Run(":3123")
}

func rateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limit.Allow() == false {
			// Rate limit exceeded
			c.File("ratelimit.png")
			c.Abort()
			return
		}
	}
}

func getTile(c *gin.Context) {
	var r string
	z, err := strconv.Atoi(c.Param("z"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	x, err := strconv.Atoi(c.Param("x"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	temp := strings.Split(c.Param("y"), "@")
	y, err := strconv.Atoi(temp[0])
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if len(temp) != 2 {
		r = ""
	} else {
		r = temp[1]
	}
	tileData, err := db.GetTile(x, y, z, r)
	if err != nil {
		// No cache found, fetch from OSM
		log.Printf("Cache miss! fetching from Origin (X=%d, Y=%d, Z=%d, R=%s)\n", x, y, z, r)
		tileData, err := tile.GetTileFromOrigin(x, y, z, r)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		// Save to cache
		err = db.AddTile(x, y, z, r, tileData)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Header("tile-from", "origin")
		c.Data(http.StatusOK, "image/png", tileData)
	} else {
		// Return from cache
		log.Printf("Cache hit! (X=%d, Y=%d, Z=%d, R=%s)\n", x, y, z, r)
		c.Header("tile-from", "cache")
		c.Data(http.StatusOK, "image/png", tileData)
	}
}
