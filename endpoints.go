package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type task struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

func getSession() *gocql.Session {
	//cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster := gocql.NewCluster("35.81.53.7", "35.81.67.246")
	cluster.Keyspace = "control_tareas"
	cluster.Timeout = time.Second * 60
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username:              "iccassandra",
		Password:              "0d3d7fae6c8664f43527531a742582fb",
		AllowedAuthenticators: []string{"com.instaclustr.cassandra.auth.InstaclustrPasswordAuthenticator"},
	}
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		println("Hubo un error: " + err.Error())
	}
	return session
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	session := getSession()
	defer session.Close()

	// Crear un canal para recibir los resultados de la consulta
	taskChan := make(chan task)
	errChan := make(chan error)

	// Ejecutar la consulta en una goroutine
	go func() {
		iter := session.Query("SELECT \"ID\",\"Name\",\"Content\" FROM tareas").Iter()

		var id, name, content string

		for iter.Scan(&id, &name, &content) {
			t := task{ID: id, Name: name, Content: content}
			taskChan <- t
		}

		// Comprobar si hubo algún error en la iteración
		if err := iter.Close(); err != nil {
			errChan <- err
		}

		close(taskChan)
		close(errChan)
	}()

	var tasks allTasks
	var errors []error

	for t := range taskChan {
		tasks = append(tasks, t)
	}

	for err := range errChan {
		errors = append(errors, err)
	}

	// Comprobar si hubo algún error en la consulta
	if len(errors) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al obtener las tareas"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	session := getSession()
	defer session.Close()

	// Goroutine para controlar la concurrencia
	errChan := make(chan error, 1)
	go func() {
		query := "SELECT * FROM control_tareas.tareas WHERE \"Name\" = ?"
		iter := session.Query(query, name).Iter()

		var tasks []task
		var t task
		for iter.Scan(&t.ID, &t.Content, &t.Name) {
			tasks = append(tasks, t)
		}

		if err := iter.Close(); err != nil {
			errChan <- err
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

		errChan <- nil
	}()

	// Control de error en caso de que no se pueda conectar a la base de datos
	select {
	case err := <-errChan:
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case <-time.After(time.Second * 30):
		http.Error(w, "Tiempo de espera agotado para la conexión a la base de datos", http.StatusInternalServerError)
		return
	}
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

	// Crear un canal para recibir errores
	errChan := make(chan error, 1)

	// Goroutine para la inserción en la base de datos
	go func() {
		// Control de error en la inserción
		if err := session.Query("INSERT INTO tareas (\"ID\", \"Name\", \"Content\") VALUES (?, ?, ?)",
			uuid, t.Name, t.Content).Exec(); err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Esperar el resultado de la goroutine
	err = <-errChan
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	println("Task with name " + t.Name + " created")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	session := getSession()
	defer session.Close()

	errChan := make(chan error) // Canal para recibir errores de la goroutine

	go func() {
		err := session.Query("DELETE FROM control_tareas.tareas WHERE \"ID\" = ?;", taskID).Exec()
		errChan <- err // Enviar el error al canal
	}()

	select {
	case err := <-errChan: // Recibir el error de la goroutine
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case <-time.After(time.Second * 30): // Tiempo de espera máximo
		http.Error(w, "Tiempo de espera agotado", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Tarea con ID %s eliminada exitosamente", taskID)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	var task task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := getSession()
	defer session.Close()

	// Utilizamos un canal para recibir el resultado de la goroutine
	resultChan := make(chan error)

	go func() {
		err := session.Query(`
			UPDATE control_tareas.tareas SET "Content" = ?, "Name" = ? WHERE "ID" = ?
		`).Bind(task.Content, task.Name, taskID).Exec()
		resultChan <- err
	}()

	// Esperamos el resultado de la goroutine
	err = <-resultChan

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.ID = taskID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
