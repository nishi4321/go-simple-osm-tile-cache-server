package tile

import (
	"fmt"
	"go-simple-osm-tile-cache-server/internal/config"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func GetTileFromOrigin(x, y, z int, r string) ([]byte, error) {
	// Fetch tile from OSM
	url := generateOriginURL(x, y, z, r)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch tile: %s", resp.Status)
	}

	tile, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return tile, nil
}

func generateOriginURL(x, y, z int, r string) string {
	config := config.Get()
	url := config.TILE.ORIGIN
	// Set subdomain character
	subdomainChar := string([]rune("abc")[rand.Intn(3)])
	url = strings.ReplaceAll(url, "{s}", subdomainChar)
	// Set tile values
	url = strings.ReplaceAll(url, "{x}", strconv.Itoa(x))
	url = strings.ReplaceAll(url, "{y}", strconv.Itoa(y))
	url = strings.ReplaceAll(url, "{z}", strconv.Itoa(z))
	if r != "" {
		url = strings.ReplaceAll(url, "{r}", "@"+r)
	} else {
		url = strings.ReplaceAll(url, "{r}", "")
	}
	// Set optional values
	for _, v := range config.TILE.VALUES {
		url = strings.ReplaceAll(url, "{"+v.Key+"}", v.Value)
	}
	return url
}
