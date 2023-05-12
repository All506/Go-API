package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	/*{
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
	},*/
}

func getSession() *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "control_tareas"
	session, _ := cluster.CreateSession()
	return session
}

func main() {
	// modo estricto: /tasks/ no es valido, debe escribir url correcta si o si

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	println("API is running...")
	//crea servidor http y muestra posibles errores en ejecucion
	log.Fatal(http.ListenAndServe(":3000", router))

}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	session := getSession()
	defer session.Close()

	iter := session.Query("SELECT \"ID\",\"Name\",\"Content\" FROM tareas").Iter()
	var id int
	var name, content string

	tasks := allTasks{}
	for iter.Scan(&id, &name, &content) {
		tasks = append(tasks, task{ID: id, Name: name, Content: content})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	session := getSession()
	taskID, err := strconv.Atoi(vars["id"])
	println(taskID)

	if err != nil {
		fmt.Fprintf(w, "Invalid Id")
		return
	}

	defer session.Close()
	iter := session.Query("SELECT \"ID\",\"Name\",\"Content\" FROM tareas WHERE \"ID\"= ?", taskID).Iter()
	var name, content string

	tasks := allTasks{}
	for iter.Scan(&taskID, &name, &content) {
		print(&taskID, &name, &content)
		tasks = append(tasks, task{ID: taskID, Name: name, Content: content})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	//Separa el JSON a un objeto
	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid Id")
		return
	}

	// i representa el indice
	for i, task := range tasks {
		if task.ID == taskID {
			//Conserva todo lo que este detras y delante del indice
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "Task with id %v has been remove succesfully", taskID)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "Invalid Id")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please insert valid data")
		return
	}

	json.Unmarshal(reqBody, &updatedTask)

	for i, task := range tasks {
		if task.ID == taskID {
			//Conserva todo lo que este detras y delante del indice
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			//Se añade la nueva tarea
			tasks = append(tasks, updatedTask)
			fmt.Fprintf(w, "Task with id %v has been updated", taskID)
		}
	}

}