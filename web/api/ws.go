package api

import (
	"html/template"
	"log"
	"net/http"

	"github.com/abielalejandro/web/config"
	"github.com/abielalejandro/web/internals/event"
	"github.com/abielalejandro/web/pkg/logger"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var emojis map[string][]string = map[string][]string{
	"POSITIVE": {"&#128512", "&#128513", "&#128522", "&#128525"},
	"NEGATIVE": {"&#128530", "&#128544", "&#128548", "&#128545"},
}

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
	chR    <-chan event.SentimentalResult
}

type Api interface {
	Run()
}

func NewHttpApi(
	config *config.Config,
	chW chan<- string,
	chR <-chan event.SentimentalResult,
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
	httpApi.Router.HandleFunc("/", httpApi.homeHandler).Methods(http.MethodGet)
	httpApi.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	httpApi.readLoopMsgs()
	log.Fatal(http.ListenAndServe(httpApi.config.HTTP.Port, httpApi.Router))
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

func (httpApi *HttpApi) render(w http.ResponseWriter, tpl string, data interface{}) {
	tmpl, err := template.ParseFiles(tpl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (httpApi *HttpApi) homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (httpApi *HttpApi) handleConn(conn *websocket.Conn) {
	httpApi.conns[conn] = true
	httpApi.readLoop(conn)
}

func (httpApi *HttpApi) readLoop(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			httpApi.log.Error(err)
			return
		}

		t := &event.Message{
			Msg: string(p),
		}

		id, _ := uuid.NewRandom()
		event := cloudevents.NewEvent()
		event.SetID(id.String())
		event.SetDataContentType("application/json")
		event.SetSource("sentimental/ws")
		event.SetType(httpApi.config.RabbitEventBus.ProducerMasterRoutingKey)
		event.SetData(cloudevents.ApplicationJSON, t)
		httpApi.broadcastMsg(&event)

		httpApi.chW <- string(p[:])
		httpApi.log.Info(string(p[:]))
	}

}

func (api *HttpApi) broadcastMsg(evt *cloudevents.Event) {
	for conn, _ := range api.conns {
		if err := conn.WriteJSON(evt); err != nil {
			api.log.Error(err)
			conn.Close()
			delete(api.conns, conn)
		}
	}
}

func (api *HttpApi) readLoopMsgs() {
	go func() {
		for elem := range api.chR {

			e := api.searchEmoji(&elem)

			id, _ := uuid.NewRandom()
			event := cloudevents.NewEvent()
			event.SetID(id.String())
			event.SetDataContentType("text/plain")
			event.SetSource("sentimental/ws")
			event.SetType(api.config.RabbitEventBus.ConsumerMasterRoutingKey)
			event.SetData(cloudevents.TextPlain, e)
			api.broadcastMsg(&event)
		}
	}()
}

func (api *HttpApi) searchEmoji(s *event.SentimentalResult) string {
	ss, ok := emojis[s.Label]

	if !ok {
		return "&#128533"
	}

	var i int
	if s.Score >= 0.0 && s.Score < 0.25 {
		i = 0
	}
	if s.Score >= 0.25 && s.Score < 0.50 {
		i = 1
	}
	if s.Score >= 0.50 && s.Score < 0.75 {
		i = 2
	}
	if s.Score >= 0.75 {
		i = 3
	}

	return ss[i]
}

func NewApi(config *config.Config, chW chan<- string, chR <-chan event.SentimentalResult) Api {
	return NewHttpApi(config, chW, chR)
}
