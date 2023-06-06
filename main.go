package main

import (
	"log"
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

type APIHandler struct {
	session *gocql.Session
}

func NewAPIHandler(session *gocql.Session) *APIHandler {
	return &APIHandler{
		session: session,
	}
}

func main() {
	cloudCassandra := true
	cluster := gocql.NewCluster("127.0.0.1:9042")

	if cloudCassandra {
		cluster = gocql.NewCluster("35.81.53.7", "35.81.67.246")
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username:              "iccassandra",
			Password:              "0d3d7fae6c8664f43527531a742582fb",
			AllowedAuthenticators: []string{"com.instaclustr.cassandra.auth.InstaclustrPasswordAuthenticator"},
		}
		cluster.Consistency = gocql.Quorum
		cluster.ProtoVersion = 4

	}

	cluster.Keyspace = "control_tareas"
	cluster.Timeout = time.Second * 60

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Error al establecer la conexión a Cassandra:", err)
	}
	defer session.Close()

	apiHandler := NewAPIHandler(session)

	router := mux.NewRouter()
	router.HandleFunc("/tasks", apiHandler.GetTasks).Methods(http.MethodGet)
	router.HandleFunc("/tasks", apiHandler.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id}", apiHandler.UpdateTask).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", apiHandler.DeleteTask).Methods(http.MethodDelete)
	router.HandleFunc("/tasks/{name}", apiHandler.GetTask).Methods(http.MethodGet)

	log.Println("Servidor iniciado en http://localhost:3000")
	err = http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
