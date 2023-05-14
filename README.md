<!---DIVISION-->

# Investigación Ingeniería de Software
<h2 align="center">
Golang + Cassandra
</h2>
<p align="center">   
<br> Estudiantes: <br> <br>
<table align="center">
	<tbody>
		<tr>
			<td>Allán Fabián Trejos Salazar</td>
			<td>C07870</td>
		</tr>
		<tr>
			<td>Luis Arguedas Villalobos</td>
			<td>C00648</td>
		</tr>
	</tbody>
</table>
</p>
<p align="center">   
Profesor: <br> <br>
MSc. Leonardo Camacho Navarro <br><br>
</p>

<h2 align="center"> 
Universidad de Costa Rica, 2023
</h2>

***

<!---DIVISION-->

**Índice**
1. [Golang y Cassandra](#Ensayo)
2. [Instalaciones necesarias](#Instalacion)
3. [Extensiones Golang para APIs](#Extensiones)
4. [Referencias](#Referencias)

<!--DIVISION-->

<div id='Ensayo'>

# ¿Qué es Golang y Cassandra?
	
Golang fue creado por Google en el año 2007 y desde entonces una variedad de productos y servicios de consumo masivo han sido implementados exitosamente con este lenguaje. La implementación de Golang en su gran mayoría se debe a la facilidad de compilar el código con gran velocidad y precisión, donde algunas empresas han convertido segundos en milisegundos al migrar sus plataformas a este lenguaje. Otra de las grandes ventajas que Go provee sobre otros lenguajes utilizados actualmente como C# es el reducido consumo de recursos y la facilidad de escalabilidad que provee. Tras la publicación de Golang los casos de éxito no han hecho más que aumentar exponencialmente, empresas como Meta han creado ORMs desde cero con el lenguaje, Dropbox migró sus backend de Python a Go lo que mejoró la concurrencia y tiempos de ejecución, Netflix al ocupar un lenguaje que generara menos latencia que Java comenzaron a utilizar Go y finalmente, Riot Games implementa el lenguaje para la mayoría de sus operaciones backend debido a que según Aaron Torres, Administrador de Ingeniería de esta compañía, Go es relativamente sencillo en comparación a otros lenguajes, el código se construye tan rápido que se puede editar y volver a construir sin mayor complicación facilitando las pruebas, Go tiene una amplia biblioteca de producción de web server, raramente rompe la retrocompatibilidad y cuando lo hace generalmente se debe a las bibliotecas y no al lenguaje y concluye con que Go es bastante popular, por lo que hay un gran soporte por compañías de terceros.
	
En cuanto a Cassandra, es una base de datos distribuida NoSQL de código abierto que miles de compañías alrededor del mundo suelen utilizar debido a su facilidad de escalar y alta disponibilidad sin nunca comprometer el rendimiento de las consultas CQL (Cassandra Query Language) que se realizan. En el caso de que una base de datos Cassandra necesite más potencia, simplemente se pueden crear nuevos nodos que se interconectan entre sí a través del protocolo P2P por lo que cada nodo tendrá un token que ayudará a identificar donde está almacenada la información. Algunas de las empresas que suelen utilizar Cassandra Db son: Facebook, Spotify, Reddit e Instagram.

Entre las razones por las que muchas compañias actualmente implementan Cassandra destaca que esta base de datos es excelente para almacenar datos de series temporales donde no se necesitan actualizar datos antiguos, la sencilla expansión geográfica que propone ya que debido a su modelo bajo funcionamiento P2P no existe un clúster padre del que los demás dependan, además de que el envio de información a nodos lejanos puede no ser tan caro por el mismo modelo anteriormente mencionado.

</div>
<!---DIVISION-->


<!---DIVISION-->

<div id='Instalacion'>

# Instalación Golang y Cassandra

## Instalación Golang

Para descargar Golang es necesario acceder al [link](https://go.dev/dl/) y descargar la versión más reciente con respecto al Sistema Operativo que actualmente se utiliza. Posteriormente a la descarga, se abre el archivo exe, se selecciona la ubicación donde se instalará el compilador y finalmente se instalara sin solicitar ninguna información demás.

## Instalación Cassandra

La última versión con soporte para Windows de Cassandra es Apache Cassandra 3.0, por lo tanto, al ingresar a la sección de [descargas](https://cassandra.apache.org/_/download.html) de Apache Cassandra seleccionaremos la opción 3.0 de Octubre 23, 2022. 

```
Es necesario tener Python y un JDK instalado para ejecutar Cassandra localmente
```

En sistemas de producción donde el uso de Cassandra es sí o sí necesario, se recomienda la instalación a través de imagenes de Docker o directamente utilizar una máquina con cualquier sistema Linux.

Tras extraer los archivos descargados y moverlos a un lugar seguro, es necesario añadir una variable de sistema en la sección temp de Windows, que rediriga a la carpeta bin dentro del lugar donde Cassandra fue instalado

<p align="center">
<img src="./resources/install1.PNG">
</p>

Una vez agregada la variable de entorno en el sistema, es necesario hacer modificaciones en el archivo cassandra.bat. Se debera de añadir la dirección donde se encuentra el jdk de la computadora que actualmente estamos utilizando. La línea de código se deberá añadir desde cero, por lo que la estructura de esta será:

```
set JAVA_HOME=DIRECCIÓN_DE_JDK
```

<p align="center">
<img src="./resources/install2.PNG">
</p>


Finalmente, para verificiar que la instalación de Cassandra Db fue realizada con éxito, se deberá de ingresar a través de CMD de Windows a la carpeta donde se guardo el Cassandra, especificamente a la carpeta bin y ejecutar el commando:

```
Cassandra
```

<p align="center">
<img src="./resources/install3.PNG"> <p>Node localhost/127.0.0.1 state jump to NORMAL</p>
</p>


</div>
<!---DIVISION-->
<div id='Extensiones'>

# Extensiones Golang

## Instalar Go Mux

Go Mux es una extensión que ayuda a manejar direcciones HTTP para la creación de APIs en Golang.

Para su instalación es necesario ejecutar el siguiente comando en consola: 

```
go get -u github.com/gorilla/mux
```

Para abrir un servidor local que ejecute la API se debe de correr el comando:

```
go run .
```

Este comando ejecutará todos los archivos con extensión .go que se encuentren dentro de la dirección.

## Instalar Compile Demon

En el camino se harán muchos cambios en los endpoints, por lo que cancelar la ejecución del servidor para después volverlo a ejecutar no es una opción. Compile Demon facilita esta tarea a través de guardar los cambios en los endpoints y re ejecutar el servidor. Para instalar esta herramienta se debe de correr el siguiente comando en consola:

```
go install github.com/githubnemo/CompileDaemon
```

Para asignar el .exe creado a un comando, se debe ejecutar:

```
CompileDaemon -command="Nombre del EXE"
```

## Controlador para conectar con Cassandra

En el caso de Cassandra, es necesario instalar un controlador (gocql) que sea capaz de enviar las queries para que posteriormente sean ejecutadas por el gestor de bases de datos. El comando para instalar el conector es:

```
go get github.com/gocql/gocql
```
</div>

<!---DIVISION-->
<div id='Referencias'>
	
# Referencias

[Casos de Estudio Go](https://go.dev/solutions/case-studies)
	
[Golang Game Development & Operations in Riot Games](https://technology.riotgames.com/news/leveraging-golang-game-development-and-operations)
	
[Introduccion a Apache Cassandra](https://aprenderbigdata.com/introduccion-apache-cassandra/)

[Cassandra Basics](https://cassandra.apache.org/_/cassandra-basics.html)
	
[Cassandra Top Benefits - Canonical Post](https://ubuntu.com/blog/apache-cassandra-top-benefits)
	

</div>
