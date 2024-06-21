package sweeper

import (
	"encoding/json"
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
	Games           map[string]*Game
	Mutex           sync.Mutex
	MenuConnections map[*websocket.Conn]bool
}

func InitServerGames() (sg *ServerGames) {
	sg = &ServerGames{
		Games:           make(map[string]*Game),
		Mutex:           sync.Mutex{},
		MenuConnections: make(map[*websocket.Conn]bool),
	}
	return
}

func (sg *ServerGames) MenuConnectionRoute(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid upgrader function"})
		return
	}

	sg.Mutex.Lock()
	sg.MenuConnections[conn] = true
	sg.Mutex.Unlock()

	defer func() {
		sg.Mutex.Lock()
		delete(sg.MenuConnections, conn)
		sg.Mutex.Unlock()
		conn.Close()
	}()

	allGames, err := sg.GetMenuListJSON()
	if err != nil {
		c.JSON(400, gin.H{"error": "issue creating json"})
		return
	}

	err = conn.WriteMessage(1, allGames)
	if err != nil {
		c.JSON(400, gin.H{"error": "issue publishing menu entries"})
		return
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (sg *ServerGames) BoradcastToMenuConnections(msgType int, data []byte) (err error) {
	// need to publish menu entries to all conns on new game
	for conn := range sg.MenuConnections {
		err = conn.WriteMessage(msgType, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sg *ServerGames) GetMenuListJSON() (data []byte, err error) {
	type MenuEntry struct {
		User   string `json:"username"`
		GameID string `json:"id"`
	}

	var MenuEntries []MenuEntry

	for _, Game := range sg.Games {
		MenuEntries = append(MenuEntries, MenuEntry{
			User:   Game.Username,
			GameID: Game.PublicSecret,
		})
	}

	data, err = json.Marshal(MenuEntries)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
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
		// new game section
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

		// publish to home connections that new game has been created
		allGames, err := sg.GetMenuListJSON()
		if err != nil {
			c.JSON(400, gin.H{"error": "could not generate game list"})
			return
		}

		err = sg.BoradcastToMenuConnections(1, allGames)
		if err != nil {
			c.JSON(400, gin.H{"error": "could not broadcast to home page conns"})
			return
		}

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

	sg.Mutex.Lock()
	game, exists := sg.Games[publisherSecret]
	if !exists {
		return
	}
	sg.Mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		// read msg here and do something aka play game
		fmt.Println(msg)

		// game.Board.Open()

		// get picked coordinates...

	}
}

func (g *Game) BroadcastToSubs(msgType int, msg []byte) (err error) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	for sub := range g.Subscribers {
		err = sub.WriteMessage(msgType, msg)
		if err != nil {
			return err
		}
	}
	return
}

// route determines all subscriber sockets, even used by the game publisher to view the board
func (sg *ServerGames) SubscriberRoute(c *gin.Context) {
	publicSecret := c.Query("publicSecret")
	if publicSecret == "" {
		c.JSON(400, gin.H{"error": "could not parse game public secret"})
		return
	}

	game := sg.getPublicGame(publicSecret)
	if game == nil {
		c.JSON(400, gin.H{"error": "public game could not be found"})
		return
	}

	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": "public game could not be found"})
		return
	}

	// add conn to subscribers

	game.Mutex.Lock()
	game.Subscribers[conn] = true
	game.Mutex.Unlock()

	defer func() {
		game.Mutex.Lock()
		delete(game.Subscribers, conn)
		game.Mutex.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}

}

func (sg *ServerGames) getPublicGame(publicSecret string) *Game {
	sg.Mutex.Lock()
	defer sg.Mutex.Unlock()

	for _, game := range sg.Games {
		if game.PublicSecret == publicSecret {
			// upgrade con on the game
			return game

		}
	}
	return nil
}
