package main

import (
	"flag"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

const defaultPort = 9099

func main() {

	port := flag.Int("p", defaultPort, "PORT")
	flag.Parse()

	var handlers []handler
	handlers = append(handlers, todoHandler()...)
	handlers = append(handlers, todoAlarmHandler()...)
	handlers = append(handlers, alarmTargetHandler()...)

	sort.SliceStable(handlers, func(i, j int) bool {
		return len(handlers[i].path) > len(handlers[j].path)
	})

	fmt.Println(handlers)

	router := mux.NewRouter()
	makeHandlers(handlers, router, "/api")

	fmt.Println("Listening on Port :", *port)

	http.ListenAndServe(":"+strconv.Itoa(*port), router)
}
