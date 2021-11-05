package main

//para las variables de entorno en esta ocasion utlize gototenv pero una alternativa a implementar seria VIPER ya que este ademas del ".env" admite JSON, TOML, YAML, HCL, envfile and Java properties config files
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connectDB() (conexion *sql.DB) {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connStr := os.Getenv("DB_POSTGRESQL_URI")

	conexion, err2 := sql.Open("postgres", connStr)

	if err2 != nil {
		panic(err.Error())
	}

	return conexion
}

//leemos los archivos que estan el la carpeta platillas
var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", startApp)
	http.HandleFunc("/create", createPage)
	http.HandleFunc("/insertar", insertNew)
	http.HandleFunc("/delete", deleteEmpleado)
	http.HandleFunc("/edit", editPage)
	http.HandleFunc("/update", updateEmpleado)
	log.Println("Servidor Corriendo")
	http.ListenAndServe(":5000", nil)
}

type Empleado struct {
	ID     int
	Nombre string
	Correo string
}

//con w respondo la peticon
//en r esta la informacion que me estan enviando por ejemplo en el r.Body
func startApp(w http.ResponseWriter, r *http.Request) {
	conexionEstablecida := connectDB()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	empleadosArray := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.ID = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		empleadosArray = append(empleadosArray, empleado)
	}

	fmt.Println(empleadosArray)

	//fmt.Fprintf(w, "Hola develoteca")
	//Si el tercer  parametro es nil significa no se le pasa nada la template "index"
	templates.ExecuteTemplate(w, "index", empleadosArray)

}

func createPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "createPerson", nil)
}

func editPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "editPerson", nil)
}

func insertNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := connectDB()

		sqlStatement := `INSERT INTO empleados(nombre, correo) VALUES($1, $2 )`

		insertarRegistros, err := conexionEstablecida.Prepare(sqlStatement)

		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(nombre, correo)

		http.Redirect(w, r, "/", 301)

	}
}

func deleteEmpleado(w http.ResponseWriter, r *http.Request) {

}

func updateEmpleado(w http.ResponseWriter, r *http.Request) {

}
