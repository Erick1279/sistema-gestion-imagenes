package main

import (
	"fmt"
	"log"
	"net/http"
	"sistema-imagenes/controladores"
	"sistema-imagenes/repositorio"
	"sistema-imagenes/servicios"

	"github.com/gorilla/mux"
)

func main() {
	// Conectar a SQL Server
	repositorio.ConectarDB()

	// Inicializar capas del sistema
	repo := repositorio.NuevoImagenRepositorio()
	servicio := servicios.NuevoImagenServicio(repo)
	controlador := controladores.NuevoImagenControlador(servicio)

	// Crear router con gorilla/mux
	router := mux.NewRouter()

	// Middleware de logs
	router.Use(middlewareLogs)

	// Registrar todas las rutas
	controlador.RegistrarRutas(router)

	fmt.Println("  SISTEMA DE GESTIÓN DE IMÁGENES")
	fmt.Println("  Servidor corriendo en puerto 8080")
	fmt.Println("Rutas disponibles:")
	fmt.Println("  POST   /imagenes          - Agregar imagen")
	fmt.Println("  GET    /imagenes          - Listar imágenes")
	fmt.Println("  GET    /imagenes/{id}     - Obtener por ID")
	fmt.Println("  GET    /imagenes/buscar   - Buscar por nombre")
	fmt.Println("  PUT    /imagenes/{id}     - Actualizar imagen")
	fmt.Println("  DELETE /imagenes/{id}     - Eliminar imagen")

	log.Fatal(http.ListenAndServe(":8080", router))
}

// Middleware de logs - registra cada petición HTTP
func middlewareLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[LOG] %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
