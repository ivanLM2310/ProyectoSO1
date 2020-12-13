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
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/cpu"
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

//ListProceso Lista de procesos
type ListProceso struct {
	Lista []Proceso `json:"lista"`
}

//UtilizacionR es estructura para guardar la ram utilizada
type UtilizacionR struct {
	Total       string `json:"total"`
	Utilizacion string `json:"utilizacion"`
	Porcentaje  string `json:"porcentaje"`
}

//UtilizacionCPU es estructura para guardar porcentajes cpu
type UtilizacionCPU struct {
	CPU []string `json:"cpu"`
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
			//var salidaT []byte
			if value == "PRINCIPAL" {
				//----------------------------PAGINA PRINCIPAL---------------------------
				data, err := ioutil.ReadFile("/proc/cpu_201122826")
				if err != nil {
					fmt.Println("File reading error", err)
					return
				}
				strData := string(data)
				//fmt.Print(string(data))

				arrLineas := strings.Split(strData, "\n")
				arrS := []Proceso{}

				for _, linea := range arrLineas {
					if linea != "" {
						arrElement := strings.Split(linea, "\t\t")
						if len(arrElement) == 4 {
							var e Proceso
							e.PID = strings.Trim(strings.Split(arrElement[0], ":")[1], " ")
							e.Nombre = strings.Trim(strings.Split(arrElement[1], ":")[1], " ")
							e.Porcen = strings.Trim(strings.Split(arrElement[2], ":")[1], " ")
							e.Estado = strings.Trim(strings.Split(arrElement[3], ":")[1], " ")
							e.Padre = ""
							e.Usuario = ""
							arrS = append(arrS, e)
						} else if len(arrElement) == 5 {
							var e Proceso
							e.PID = strings.Trim(strings.Split(arrElement[1], ":")[1], " ")
							e.Nombre = strings.Trim(strings.Split(arrElement[2], ":")[1], " ")
							e.Porcen = strings.Trim(strings.Split(arrElement[3], ":")[1], " ")
							e.Estado = strings.Trim(strings.Split(arrElement[4], ":")[1], " ")
							e.Padre = strings.Trim(strings.Split(arrElement[0], ":")[1], " ")
							e.Usuario = ""
							arrS = append(arrS, e)
						} else {
							//ERROR
						}
					}
				}

				/*
					salidaJI := &Message{
						//Dato: value + "_aca ya se junto papu"}
						Dato: string(data)}
					salidaJ, _ := json.Marshal(salidaJI)
					salidaT = salidaJ
					fmt.Println(string(salidaJ))
				*/
				salidaJI := &ListProceso{
					Lista: arrS}
				errW := client.WriteJSON(salidaJI)
				if errW != nil {
					log.Printf("error: %v", errW)
					client.Close()
					delete(clients, client)
				}

			} else if value == "CPU" {
				//----------------------------PAGINA CPU---------------------------
				salidaJ := getInfoCPU()
				errW := client.WriteJSON(&salidaJ)
				if errW != nil {
					log.Printf("error: %v", errW)
					client.Close()
					delete(clients, client)
				}
			} else if value == "RAM" {
				//----------------------------PAGINA RAM---------------------------
				data, err := ioutil.ReadFile("/proc/memo_201122826")
				if err != nil {
					fmt.Println("File reading error", err)
					return
				}
				strData := string(data)
				fmt.Print(string(data))
				arrLineas := strings.Split(strData, "\n")
				elem := arrLineas[len(arrLineas)-2]
				numR := strings.Replace(strings.Trim(strings.Split(elem, ":")[1], " "), " %", "", -1)
				var r UtilizacionR
				r.Total = strings.Split(arrLineas[len(arrLineas)-4], ":")[1]
				r.Utilizacion = strings.Split(arrLineas[len(arrLineas)-3], ":")[1]
				r.Porcentaje = strings.Replace(numR, "\"", "", -1)
				errW := client.WriteJSON(&r)
				if errW != nil {
					log.Printf("error: %v", errW)
					client.Close()
					delete(clients, client)
				}
			} else {
				clients[client] = "PRINCIPAL"
				if i, err := strconv.Atoi(value); err == nil {
					proc, err := os.FindProcess(i)
					if err != nil {
						log.Println(err)
					}
					// Kill el proceso
					proc.Kill()
					log.Println("proceso eliminado ...")
				}
				continue
			}

		}
		fmt.Println(len(clients))
		log.Printf("---------")
		time.Sleep(2000 * time.Millisecond)
	}
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

//getInfoCPU funcion para sacar rendimiento cpu
func getInfoCPU() UtilizacionCPU {
	percentage, err := cpu.Percent(0, true)
	dealwithErr(err)

	listC := []string{}
	for _, cpupercent := range percentage {
		listC = append(listC, strconv.FormatFloat(cpupercent, 'f', 2, 64))
	}
	var objC UtilizacionCPU
	objC.CPU = listC
	return objC
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
