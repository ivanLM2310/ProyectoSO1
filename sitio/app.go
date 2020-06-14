package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string)

//Message es estructura
type Message struct {
	Dato string `json:"dato"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}
		fmt.Println(string(p))
		clients[conn] = string(p)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("1-------------------------------------------")
	fmt.Println(r.Host)

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("2-------------------------------------------")
		log.Println(err)
	}
	defer ws.Close()

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)

}

func envioInfo() {
	for {
		// Grab the next message from the broadcast channel

		// Send it out to every client that is currently connected

		for client := range clients {

			var value string = clients[client]
			log.Println(value)

			data, err := ioutil.ReadFile("/proc/memo_201122826")
			if err != nil {
				fmt.Println("File reading error", err)
				return
			}
			fmt.Print(string(data))

			salidaJI := &Message{
				Dato: value + "_aca ya se junto papu"}

			salidaJ, _ := json.Marshal(salidaJI)
			fmt.Println(string(salidaJ))

			errW := client.WriteJSON(string(salidaJ))
			if errW != nil {
				log.Printf("error: %v", errW)
				client.Close()
				delete(clients, client)
			}

		}
		fmt.Println(len(clients))
		log.Printf("---------")
		time.Sleep(2000 * time.Millisecond)
	}
}

//Compile templates on start
var templates = template.Must(template.ParseFiles("header.html", "footer.html", "main.html", "about.html"))

//A Page structure
type Page struct {
	Title string
}

//Display the named template
func display(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}

//The handlers.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "main", &Page{Title: "Home"})
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	display(w, "about", &Page{Title: "About"})
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/ws", serveWs)
	go envioInfo()
	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
