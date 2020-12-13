package main

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"

	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
)

//Message es estructura
type Message struct {
	Dato string `json:"dato"`
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func GetHardwareData(w http.ResponseWriter, r *http.Request) {
	// var salidaT []byte
	runtimeOS := runtime.GOOS

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	dealwithErr(err)
	percentage, err := cpu.Percent(0, true)
	dealwithErr(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)
	html := "<html>OS : " + runtimeOS + "<br>"

	// since my machine has one CPU, I'll use the 0 index
	// if your machine has more than 1 CPU, use the correct index
	// to get the proper data
	html = html + "CPU index number: " + strconv.FormatInt(int64(cpuStat[0].CPU), 10) + "<br>"
	html = html + "VendorID: " + cpuStat[0].VendorID + "<br>"
	html = html + "Family: " + cpuStat[0].Family + "<br>"
	html = html + "Number of cores: " + strconv.FormatInt(int64(cpuStat[0].Cores), 10) + "<br>"
	html = html + "Model Name: " + cpuStat[0].ModelName + "<br>"
	html = html + "Speed: " + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz <br>"

	for idx, cpupercent := range percentage {
		html = html + "Current CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%<br>"
		if idx == 1 {
			salidaJI := &Message{
				//Dato: value + "_aca ya se junto papu"}
				Dato: string(strconv.FormatFloat(cpupercent, 'f', 2, 64))}
			salidaJ, _ := json.Marshal(salidaJI)
			// salidaT = salidaJ
			fmt.Println(string(salidaJ))
		}
	}

	html = html + "Hostname: " + hostStat.Hostname + "<br>"
	html = html + "Uptime: " + strconv.FormatUint(hostStat.Uptime, 10) + "<br>"
	html = html + "Number of processes running: " + strconv.FormatUint(hostStat.Procs, 10) + "<br>"

	html = html + "</html>"

	w.Write([]byte(html))
	// errW := client.WriteJSON(string(salidaT))

}

func SayName(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, I'm a machine and my name is [whatever]"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", SayName)
	mux.HandleFunc("/gethwdata", GetHardwareData)

	http.ListenAndServe(":3080", mux)

}
