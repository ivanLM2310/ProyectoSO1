package main

import (
	"C"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)
import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

var clients = make(map[*websocket.Conn]string)

//Message es estructura
type Message struct {
	Dato string `json:"dato"`
}

//Proceso es la estructura para transformar archivo
type Proceso struct {
	PID     string `json:"pid"`
	Nombre  string `json:"nombre"`
	Usuario string `json:"usuario"`
	Estado  string `json:"estado"`
	Porcen  string `json:"porcen"`
	Padre   string `json:"padre"`
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

// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("1-------------------------------------------")
	fmt.Println(r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("2-------------------------------------------")
		log.Println(err)
	}
	defer ws.Close()

	reader(ws)
}

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

func envioInfo() {
	for {
		for client := range clients {

			var value string = clients[client]
			log.Println(value)
			var salidaT []byte
			if value == "PRINCIPAL" {
				data, err := ioutil.ReadFile("/proc/cpu_201122826")
				if err != nil {
					fmt.Println("File reading error", err)
					return
				}
				strData := string(data)
				fmt.Print(string(data))

				arrLineas := strings.Split(strData, "\n")
				for _, linea := range arrLineas {
					if linea != "" {
						arrElement := strings.Split(linea, "\t\t")
						if len(arrElement) == 4 {
							var e Proceso
							e.PID = strings.Split(arrElement[0], ":")[1]
							e.Nombre = strings.Split(arrElement[1], ":")[1]
							e.Porcen = strings.Split(arrElement[2], ":")[1]
							e.Estado = strings.Split(arrElement[3], ":")[1]
							e.Padre = ""
							e.Usuario = ""
						} else if len(arrElement) == 5 {

						} else {

						}
					}
				}

				salidaJI := &Message{
					//Dato: value + "_aca ya se junto papu"}
					Dato: string(data)}
				salidaJ, _ := json.Marshal(salidaJI)
				salidaT = salidaJ
				fmt.Println(string(salidaJ))
			} else if value == "CPU" {

			} else if value == "RAM" {

			}

			errW := client.WriteJSON(string(salidaT))
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

func main() {

	fs := http.FileServer(http.Dir("./css-Proyecto"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", serveWs)
	go envioInfo()

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
