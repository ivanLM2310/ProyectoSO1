package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

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
		// read in a message

		// print out that message for clarity
		//fmt.Println(string(p))

		err := conn.WriteJSON("msg de salida")
		if err != nil {
			log.Printf("error: %v", err)

		}
		log.Printf("---------")
		time.Sleep(2000 * time.Millisecond)
	}
}
func lectura(conn *websocket.Conn) {
	for {
		i := 0

		err := conn.WriteJSON("msg" + strconv.Itoa(i))
		if err != nil {
			log.Printf("error: %v", err)
			conn.Close()
		}
		i++
		time.Sleep(1000 * time.Millisecond)
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
	clients[ws] = true

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)

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

	//Listen on port 8080
	http.ListenAndServe(":8080", nil)
}
