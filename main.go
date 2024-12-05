package main

import (
	"go-simple-osm-tile-cache-server/internal/db"
	"go-simple-osm-tile-cache-server/internal/tile"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/tile/:z/:x/:y", getTile)
	gin.SetMode(gin.ReleaseMode)
	route.Run(":3123")
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
