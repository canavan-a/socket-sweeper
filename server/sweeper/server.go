package sweeper

import (
	"fmt"
	"sync"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Game struct {
	Username     string
	Publisher    *websocket.Conn
	Subscribers  map[*websocket.Conn]bool
	Mutex        sync.Mutex
	StartTime    time.Time
	PublicSecret string
	Board        GameBoard
}

type ServerGames struct {
	Games map[string]*Game
	Mutex sync.Mutex
}

func InitServerGames() (sg *ServerGames) {
	sg = &ServerGames{
		Games: make(map[string]*Game),
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

	// upgrade the route
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid upgrader function"})
		return
	}

	sg.Mutex.Lock()
	// find a game with publisher secret
	_, exists := sg.Games[publisherSecret]
	if exists {
		sg.Games[publisherSecret].Publisher = conn
	} else {

		user := c.Query("user")
		if user == "" {
			user = "anon"
		}

		x, err := strconv.ParseInt(c.Query("x"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "no x param included"})
			return
		}

		y, err := strconv.ParseInt(c.Query("y"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "no y param included"})
			return
		}

		bombs, err := strconv.ParseInt(c.Query("bombs"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "no bombs param included"})
			return
		}

		// check for standard board sizes here....
		// custom board sizes??

		// if not game with that publisher secret then create one
		game := Game{
			Username:     user,
			Publisher:    conn,
			Subscribers:  make(map[*websocket.Conn]bool),
			PublicSecret: uuid.NewString(),
			Mutex:        sync.Mutex{},
			StartTime:    time.Now(),
			Board:        NewGameBoard(int(x), int(y), int(bombs)),
		}
		// add game
		sg.Games[publisherSecret] = &game
	}
	// add publisher connection to subscribers
	sg.Games[publisherSecret].Mutex.Lock()
	sg.Games[publisherSecret].Subscribers[conn] = true
	sg.Games[publisherSecret].Mutex.Unlock()

	// remove publisher from subscribers
	defer func() {
		sg.Games[publisherSecret].Mutex.Lock()
		delete(sg.Games[publisherSecret].Subscribers, conn)
		sg.Games[publisherSecret].Mutex.Unlock()
		conn.Close()
	}()

	sg.Mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// read msg here and do something aka play game
		fmt.Println(msg)

	}
}

func (g *Game) BroadcastToSubs(msgType int, msg []byte) (err error) {
	for sub := range g.Subscribers {
		err = sub.WriteMessage(msgType, msg)
		if err != nil {
			return err
		}
	}
	return
}
