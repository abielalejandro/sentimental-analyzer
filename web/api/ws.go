package api

import (
	"net/http"

	"github.com/abielalejandro/web/config"
	"github.com/abielalejandro/web/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type HttpApi struct {
	config *config.Config
	Router *mux.Router
	log    *logger.Logger
	conns  map[*websocket.Conn]bool
	chW    chan<- string
	chR    <-chan string
}

type Api interface {
	Run()
}

func NewHttpApi(
	config *config.Config,
	chW chan<- string,
	chR <-chan string,
) *HttpApi {

	return &HttpApi{
		Router: mux.NewRouter().StrictSlash(true),
		config: config,
		log:    logger.New(config.Log.Level),
		conns:  make(map[*websocket.Conn]bool),
		chW:    chW,
		chR:    chR,
	}
}

func (httpApi *HttpApi) Run() {
	httpApi.Router.HandleFunc("/health", httpApi.health).Methods("GET")
	httpApi.Router.HandleFunc("/ws", httpApi.handlerWs)

	httpApi.log.Fatal(http.ListenAndServe(":8844", httpApi.Router))
}

func (httpApi *HttpApi) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("UP"))
}

func (httpApi *HttpApi) handlerWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpApi.log.Error(err)
		return
	}

	httpApi.handleConn(conn)
}

func (httpApi *HttpApi) handleConn(conn *websocket.Conn) {
	httpApi.conns[conn] = true
	httpApi.readLoop(conn)
}

func (httpApi *HttpApi) readLoop(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			httpApi.log.Error(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			httpApi.log.Error(err)
			conn.Close()
			delete(httpApi.conns, conn)
			return
		}
		httpApi.chW <- string(p[:])
		httpApi.log.Info(string(p[:]))
	}

}

func NewApi(config *config.Config, chW chan<- string, chR <-chan string) Api {
	return NewHttpApi(config, chW, chR)
}
