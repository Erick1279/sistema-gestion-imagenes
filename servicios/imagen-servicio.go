package servicios

import (
	"fmt"
	"sistema-imagenes/modelos"
	"sistema-imagenes/repositorio"
)

type ImagenServicio struct {
	repo  repositorio.ImagenRepositorio
	canal chan string
}

// Constructor del servicio
func NuevoImagenServicio(repo repositorio.ImagenRepositorio) *ImagenServicio {
	servicio := &ImagenServicio{
		repo:  repo,
		canal: make(chan string, 10),
	}
	// Goroutine que escucha notificaciones del canal
	go servicio.escucharNotificaciones()
	return servicio
}

// Goroutine - escucha eventos del sistema
func (s *ImagenServicio) escucharNotificaciones() {
	for mensaje := range s.canal {
		fmt.Println("[NOTIFICACION]:", mensaje)
	}
}

func (s *ImagenServicio) AgregarImagen(nombre, descripcion, ruta, formato string, tamanio int64) (*modelos.Imagen, error) {
	imagenTemp, err := modelos.NuevaImagen(0, nombre, descripcion, ruta, formato, tamanio)
	if err != nil {
		return nil, err
	}

	imagenGuardada, err := s.repo.Guardar(imagenTemp)
	if err != nil {
		return nil, err
	}

	s.canal <- fmt.Sprintf("Imagen '%s' agregada con ID %d", nombre, imagenGuardada.GetId())
	return imagenGuardada, nil
}

func (s *ImagenServicio) ObtenerImagenes() []*modelos.Imagen {
	return s.repo.ObtenerTodos()
}

func (s *ImagenServicio) ObtenerImagenPorId(id int) (*modelos.Imagen, error) {
	return s.repo.ObtenerPorId(id)
}

func (s *ImagenServicio) BuscarPorNombre(nombre string) (*modelos.Imagen, error) {
	return s.repo.ObtenerPorNombre(nombre)
}

func (s *ImagenServicio) ActualizarImagen(id int, nombre, descripcion string) (*modelos.Imagen, error) {
	imagen, err := s.repo.ObtenerPorId(id)
	if err != nil {
		return nil, err
	}
	if nombre != "" {
		err = imagen.SetNombre(nombre)
		if err != nil {
			return nil, err
		}
	}
	imagen.SetDescripcion(descripcion)
	err = s.repo.Actualizar(imagen)
	if err != nil {
		return nil, err
	}
	s.canal <- fmt.Sprintf("Imagen ID %d actualizada", id)
	return imagen, nil
}

func (s *ImagenServicio) EliminarImagen(id int) error {
	err := s.repo.Eliminar(id)
	if err != nil {
		return err
	}
	s.canal <- fmt.Sprintf("Imagen ID %d eliminada", id)
	return nil
}
