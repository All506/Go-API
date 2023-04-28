package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Information related to task number one",
	},
	{
		ID:      2,
		Name:    "Task Two",
		Content: "Information related to task number two",
	},
	{
		ID:      3,
		Name:    "Task Three",
		Content: "Information related to task number three",
	},
}

func main() {
	// modo estricto: /tasks/ no es valido, debe escribir url correcta si o si
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)

	//crea servidor http y muestra posibles errores en ejecucion
	log.Fatal(http.ListenAndServe(":3000", router))

}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola to my API")
}
