package sweeper

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Game struct {
	Username    string
	Publisher   *websocket.Conn
	Subscribers map[*websocket.Conn]bool
	Mutex       sync.Mutex
	StartTime   time.Time
}

type ServerGames struct {
	Games map[string]Game
	Mutex sync.Mutex
}

func InitServerGames() (sg *ServerGames) {
	sg = &ServerGames{
		Games: make(map[string]Game),
		Mutex: sync.Mutex{},
	}
	return
}

// ws handler function to
func (sg *ServerGames) PublisherRoute(c *gin.Context) {

	publisherSecret := c.Query("publisherSecret")
	if publisherSecret == "" {
		c.JSON(400, gin.H{"error": "no game found"})
		return
	}

	// find a game with publisher secret

}
