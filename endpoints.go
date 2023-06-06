package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)


func (h *APIHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []task

	iter := h.session.Query("SELECT \"ID\",\"Name\",\"Content\" FROM tareas").Iter()
	defer iter.Close()

	var id, name, content string

	for iter.Scan(&id, &name, &content) {
		t := task{ID: id, Name: name, Content: content}
		tasks = append(tasks, t)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error al obtener las tareas:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al obtener las tareas"))
		return
	}
	println("Lista de tareas obtenida")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *APIHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskName := vars["name"]

	var tasks []task

	query := "SELECT * FROM control_tareas.tareas WHERE \"Name\" = ?"
	iter := h.session.Query(query, taskName).Iter()
	defer iter.Close()

	var id, name, content string

	for iter.Scan(&id, &content, &name) {
		t := task{ID: id, Name: name, Content: content}
		tasks = append(tasks, t)
	}

	if err := iter.Close(); err != nil {
		log.Println("Error al obtener las tareas:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al obtener las tareas"))
		return
	}
	
	println("Lista filtrada de tareas obtenida")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *APIHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var t task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Println("Error al decodificar la solicitud:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Solicitud inválida"))
		return
	}

	uuid := gocql.TimeUUID()
 	t.ID = uuid.String()

	query := "INSERT INTO tareas (\"ID\", \"Name\", \"Content\") VALUES (?, ?, ?)"
	err = h.session.Query(query, t.ID, t.Name, t.Content).Exec()
	if err != nil {
		log.Println("Error al crear la tarea:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al crear la tarea"))
		return
	}

	println("Tarea con nombre " + t.Name + " creada")
	w.WriteHeader(http.StatusCreated)
}

func (h *APIHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	var t task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Println("Error al decodificar la solicitud:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Solicitud inválida"))
		return
	}

	query := "UPDATE control_tareas.tareas SET \"Content\" = ?, \"Name\" = ? WHERE \"ID\" = ?"
	err = h.session.Query(query, t.Content, t.Name, taskID).Exec()
	if err != nil {
		log.Println("Error al actualizar la tarea:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al actualizar la tarea"))
		return
	}

	println("Tarea con nombre " + t.Name + " actualizada")
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	query := "DELETE FROM control_tareas.tareas WHERE \"ID\" = ?;"
	err := h.session.Query(query, taskID).Exec()
	if err != nil {
		log.Println("Error al eliminar la tarea:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al eliminar la tarea"))
		return
	}

	println("Tarea con ID " + taskID + " eliminada")
	w.WriteHeader(http.StatusOK)
}