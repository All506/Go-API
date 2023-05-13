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

/*type Tarea struct {
	ID      gocql.UUID
	Content string
	Name    string
}*/

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
	router.HandleFunc("/tasks/{name}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}/{name}", updateTask).Methods("PUT")
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
	name := vars["name"]

	session := getSession()
	defer session.Close()

	var t task
	query := "SELECT * FROM control_tareas.tareas WHERE \"Name\" = ? LIMIT 1"
	if err := session.Query(query, name).Consistency(gocql.One).Scan(&t.ID, &t.Name, &t.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)

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

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskName := vars["name"]

	session := getSession()
	defer session.Close()

	var taskID gocql.UUID
	err := session.Query("SELECT \"ID\" FROM control_tareas.tareas WHERE \"Name\" = ?;", taskName).Scan(&taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = session.Query("DELETE FROM control_tareas.tareas WHERE \"ID\" = ?;", taskID).Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Tarea con nombre %s eliminada exitosamente", taskName)

}

func updateTask(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID y el nombre de la tarea del parámetro de la URL
	vars := mux.Vars(r)
	taskID := vars["id"]
	taskName := vars["name"]

	// Decodificar el cuerpo de la solicitud en una tarea
	var task task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Abrir una sesión con la base de datos Cassandra
	session := getSession()
	defer session.Close()

	// Actualizar la tarea en la base de datos
	err = session.Query(`
		UPDATE control_tareas.tareas SET "Content" = ? WHERE "ID" = ? AND "Name" = ?
	`).Bind(task.Content, taskID, taskName).Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver la tarea actualizada
	task.ID = taskID
	task.Name = taskName
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}
