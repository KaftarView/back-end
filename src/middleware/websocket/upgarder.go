package middleware_websocket

import (
	"first-project/src/bootstrap"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketMiddleware struct {
	constants *bootstrap.Constants
}

func NewWebsocketMiddleware(constants *bootstrap.Constants) *WebsocketMiddleware {
	return &WebsocketMiddleware{
		constants: constants,
	}
}

func (wsMiddleware *WebsocketMiddleware) UpgradeToWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
			// return r.Header.Get("Origin") == "https://frontend-domain.com"
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}
	c.Set(wsMiddleware.constants.Context.WSConnection, conn)

	c.Next()
}
