package repositorio

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func ConectarDB() {
	var err error
	DB, err = sql.Open("sqlite", "sistema_imagenes.db")
	if err != nil {
		log.Fatal("Error al abrir base de datos:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error al conectar con SQLite:", err)
	}

	// Crear tabla si no existe
	crearTabla()

	fmt.Println("Conexión a base de datos exitosa!")
}

func crearTabla() {
	query := `CREATE TABLE IF NOT EXISTS imagenes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		descripcion TEXT,
		ruta TEXT NOT NULL,
		formato TEXT,
		tamanio INTEGER,s
		fecha_creacion DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error al crear tabla:", err)
	}
}
