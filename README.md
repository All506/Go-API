# Instalar Go Mux

Correr el siguiente comando en consola: 

go get -u github.com/gorilla/mux

Para ejecutar la aplicación Go, se deberá de ejecutar el comando:

go run main.go

# Instalar Compile Demon

En el camino se harán muchos cambios en los endpoints, por lo que cancelar la ejecución del servidor para después volverlo a ejecutar no es una opción. Compile Demon facilita esta tarea a través de guardar los cambios en los endpoints y re ejecutar el servidor. Para instalar esta herramienta se debe de correr el siguiente comando en consola:

go install github.com/githubnemo/CompileDaemon

Para asignar el .exe creado a un comando, se debe ejecutar:

CompileDaemon -command="Nombre del EXE"

# Controlador para conectar con Cassandra

go get github.com/gocql/gocql