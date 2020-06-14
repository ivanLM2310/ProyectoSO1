package main

import (
	"C"
	"log"
	"net/http"
)

/*func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>hola mundo</h1>")

}*/

func calcSingleCoreUsage(curr, prev linuxproc.CPUStat) float32 {

	PrevIdle := prev.Idle + prev.IOWait
	Idle := curr.Idle + curr.IOWait

	PrevNonIdle := prev.User + prev.Nice + prev.System + prev.IRQ + prev.SoftIRQ + prev.Steal
	NonIdle := curr.User + curr.Nice + curr.System + curr.IRQ + curr.SoftIRQ + curr.Steal

	PrevTotal := PrevIdle + PrevNonIdle
	Total := Idle + NonIdle
	// fmt.Println(PrevIdle, Idle, PrevNonIdle, NonIdle, PrevTotal, Total)

	//  differentiate: actual value minus the previous one
	totald := Total - PrevTotal
	idled := Idle - PrevIdle

	CPU_Percentage := (float32(totald) - float32(idled)) / float32(totald)

	return CPU_Percentage
}

func main() {
	fs := http.FileServer(http.Dir("./css-Proyecto"))
	http.Handle("/", fs)

	log.Println("Listening on :3000..." + CPU_Percentage)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
	//http.HandleFunc("/", indexHandler)
	//http.ListenAndServe(":8000", nil)
}
