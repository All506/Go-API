# Instalar Go Mux

Correr el siguiente comando en consola: 

go get -u github.com/gorilla/mux

Para ejecutar la aplicación Go, se deberá de ejecutar el comando:

go run main.go

# Instalar Compile Demon

En el camino se harán muchos cambios en los endpoints, por lo que cancelar la ejecución del servidor para después volverlo a ejecutar no es una opción. Compile Demon facilita esta tarea a través de guardar los cambios en los endpoints y re ejecutar el servidor. Para instalar esta herramienta se debe de correr el siguiente comando en consola:

go get github.com/githubnemo/CompileDaemon

A partir de este punto, para ejecutar la api desde consola se debera de ejecutar el comando:

CompileDaemon