package controladores

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sistema-imagenes/servicios"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ImagenControlador struct {
	servicio *servicios.ImagenServicio
}

func NuevoImagenControlador(servicio *servicios.ImagenServicio) *ImagenControlador {
	return &ImagenControlador{servicio: servicio}
}

func responderJSON(w http.ResponseWriter, codigo int, dato interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(dato)
}

// POST /imagenes/subir - Subida real de archivo
func (c *ImagenControlador) SubirImagen(w http.ResponseWriter, r *http.Request) {
	// Limitar tamaño máximo a 10MB
	r.ParseMultipartForm(10 << 20)

	archivo, header, err := r.FormFile("imagen")
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "no se encontró el archivo"})
		return
	}
	defer archivo.Close()

	descripcion := r.FormValue("descripcion")
	nombre := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	formato := strings.TrimPrefix(filepath.Ext(header.Filename), ".")
	tamanio := header.Size

	// Crear carpeta imagenes si no existe
	os.MkdirAll("imagenes", os.ModePerm)

	// Guardar archivo físico
	rutaArchivo := fmt.Sprintf("imagenes/%s", header.Filename)
	destino, err := os.Create(rutaArchivo)
	if err != nil {
		responderJSON(w, http.StatusInternalServerError, map[string]string{"error": "no se pudo guardar el archivo"})
		return
	}
	defer destino.Close()
	io.Copy(destino, archivo)

	// Guardar metadatos en base de datos
	imagen, err := c.servicio.AgregarImagen(nombre, descripcion, rutaArchivo, formato, tamanio)
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// URL donde se puede acceder a la imagen
	urlImagen := fmt.Sprintf("http://localhost:8080/imagenes/archivo/%s", header.Filename)

	responderJSON(w, http.StatusCreated, map[string]interface{}{
		"id":          imagen.GetId(),
		"nombre":      imagen.GetNombre(),
		"descripcion": imagen.GetDescripcion(),
		"ruta":        imagen.GetRuta(),
		"formato":     imagen.GetFormato(),
		"tamanio":     imagen.GetTamanio(),
		"url":         urlImagen,
	})
}

// GET /imagenes
func (c *ImagenControlador) ObtenerImagenes(w http.ResponseWriter, r *http.Request) {
	imagenes := c.servicio.ObtenerImagenes()
	var resultado []map[string]interface{}
	for _, img := range imagenes {
		urlImagen := fmt.Sprintf("http://localhost:8080/imagenes/archivo/%s.%s",
			img.GetNombre(), img.GetFormato())
		resultado = append(resultado, map[string]interface{}{
			"id":          img.GetId(),
			"nombre":      img.GetNombre(),
			"descripcion": img.GetDescripcion(),
			"ruta":        img.GetRuta(),
			"formato":     img.GetFormato(),
			"tamanio":     img.GetTamanio(),
			"url":         urlImagen,
		})
	}
	if resultado == nil {
		resultado = []map[string]interface{}{}
	}
	responderJSON(w, http.StatusOK, resultado)
}

// GET /imagenes/{id}
func (c *ImagenControlador) ObtenerImagenPorId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	imagen, err := c.servicio.ObtenerImagenPorId(id)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	urlImagen := fmt.Sprintf("http://localhost:8080/imagenes/archivo/%s.%s",
		imagen.GetNombre(), imagen.GetFormato())
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"id":          imagen.GetId(),
		"nombre":      imagen.GetNombre(),
		"descripcion": imagen.GetDescripcion(),
		"ruta":        imagen.GetRuta(),
		"formato":     imagen.GetFormato(),
		"tamanio":     imagen.GetTamanio(),
		"url":         urlImagen,
	})
}

// GET /imagenes/buscar?nombre=xxx
func (c *ImagenControlador) BuscarPorNombre(w http.ResponseWriter, r *http.Request) {
	nombre := r.URL.Query().Get("nombre")
	imagen, err := c.servicio.BuscarPorNombre(nombre)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	urlImagen := fmt.Sprintf("http://localhost:8080/imagenes/archivo/%s.%s",
		imagen.GetNombre(), imagen.GetFormato())
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"id":          imagen.GetId(),
		"nombre":      imagen.GetNombre(),
		"descripcion": imagen.GetDescripcion(),
		"ruta":        imagen.GetRuta(),
		"formato":     imagen.GetFormato(),
		"tamanio":     imagen.GetTamanio(),
		"url":         urlImagen,
	})
}

// PUT /imagenes/{id}
func (c *ImagenControlador) ActualizarImagen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	var datos struct {
		Nombre      string `json:"nombre"`
		Descripcion string `json:"descripcion"`
	}
	json.NewDecoder(r.Body).Decode(&datos)
	imagen, err := c.servicio.ActualizarImagen(id, datos.Nombre, datos.Descripcion)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, map[string]interface{}{
		"id":          imagen.GetId(),
		"nombre":      imagen.GetNombre(),
		"descripcion": imagen.GetDescripcion(),
		"ruta":        imagen.GetRuta(),
		"formato":     imagen.GetFormato(),
		"tamanio":     imagen.GetTamanio(),
	})
}

// DELETE /imagenes/{id}
func (c *ImagenControlador) EliminarImagen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responderJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}
	err = c.servicio.EliminarImagen(id)
	if err != nil {
		responderJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"mensaje": "imagen eliminada correctamente"})
}

// RegistrarRutas registra todas las rutas con gorilla/mux
func (c *ImagenControlador) RegistrarRutas(r *mux.Router) {
	r.HandleFunc("/imagenes/subir", c.SubirImagen).Methods("POST")
	r.HandleFunc("/imagenes/buscar", c.BuscarPorNombre).Methods("GET")
	r.HandleFunc("/imagenes", c.ObtenerImagenes).Methods("GET")
	r.HandleFunc("/imagenes/{id}", c.ObtenerImagenPorId).Methods("GET")
	r.HandleFunc("/imagenes/{id}", c.ActualizarImagen).Methods("PUT")
	r.HandleFunc("/imagenes/{id}", c.EliminarImagen).Methods("DELETE")

	// Servir archivos estáticos desde la carpeta imagenes
	r.PathPrefix("/imagenes/archivo/").Handler(
		http.StripPrefix("/imagenes/archivo/", http.FileServer(http.Dir("imagenes/"))),
	)
}
