package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type task struct {
	ID      string `json:ID`
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
	cluster.Timeout = time.Second * 60
	session, _ := cluster.CreateSession()
	return session
}

func main() {
	// modo estricto: /tasks/ no es valido, debe escribir url correcta si o si

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{name}", getTask).Methods("GET")
	//router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	//router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
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

	var id, name, content string

	tasks := allTasks{}
	for iter.Scan(&id, &name, &content) {
		tasks = append(tasks, task{ID: id, Name: name, Content: content})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskName := vars["name"]

	session := getSession()
	defer session.Close()

	var content string
	query := fmt.Sprintf("SELECT \"Content\" FROM control_tareas.tareas WHERE \"Name\" = '%s';", taskName)
	fmt.Println(query)

	err := session.Query(query).Scan(&content)
	if err == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task := task{Name: taskName, Content: content}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func createTask(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t task
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	session := getSession()
	defer session.Close()

	uuid := gocql.TimeUUID()
	t.ID = uuid.String()

	if err := session.Query("INSERT INTO tareas (\"ID\", \"Name\", \"Content\") VALUES (?, ?, ?)",
		uuid, t.Name, t.Content).Exec(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)

}

/*
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
*/
