package repositorio

import (
	"errors"
	"sistema-imagenes/modelos"
)

type ImagenRepositorio interface {
	Guardar(imagen *modelos.Imagen) (*modelos.Imagen, error)
	ObtenerTodos() []*modelos.Imagen
	ObtenerPorId(id int) (*modelos.Imagen, error)
	ObtenerPorNombre(nombre string) (*modelos.Imagen, error)
	Actualizar(imagen *modelos.Imagen) error
	Eliminar(id int) error
}

type ImagenRepositorioSQL struct{}

func NuevoImagenRepositorio() *ImagenRepositorioSQL {
	return &ImagenRepositorioSQL{}
}

func (r *ImagenRepositorioSQL) Guardar(imagen *modelos.Imagen) (*modelos.Imagen, error) {
	query := `INSERT INTO imagenes (nombre, descripcion, ruta, formato, tamanio) 
	          VALUES (?, ?, ?, ?, ?)`

	result, err := DB.Exec(query,
		imagen.GetNombre(),
		imagen.GetDescripcion(),
		imagen.GetRuta(),
		imagen.GetFormato(),
		imagen.GetTamanio(),
	)
	if err != nil {
		return nil, err
	}

	nuevoId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return modelos.NuevaImagen(int(nuevoId), imagen.GetNombre(), imagen.GetDescripcion(),
		imagen.GetRuta(), imagen.GetFormato(), imagen.GetTamanio())
}

func (r *ImagenRepositorioSQL) ObtenerTodos() []*modelos.Imagen {
	query := `SELECT id, nombre, descripcion, ruta, formato, tamanio FROM imagenes`
	rows, err := DB.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var lista []*modelos.Imagen
	for rows.Next() {
		var id int
		var nombre, descripcion, ruta, formato string
		var tamanio int64
		rows.Scan(&id, &nombre, &descripcion, &ruta, &formato, &tamanio)
		img, err := modelos.NuevaImagen(id, nombre, descripcion, ruta, formato, tamanio)
		if err == nil {
			lista = append(lista, img)
		}
	}
	return lista
}

func (r *ImagenRepositorioSQL) ObtenerPorId(id int) (*modelos.Imagen, error) {
	query := `SELECT id, nombre, descripcion, ruta, formato, tamanio FROM imagenes WHERE id = ?`
	var nombre, descripcion, ruta, formato string
	var tamanio int64

	err := DB.QueryRow(query, id).Scan(&id, &nombre, &descripcion, &ruta, &formato, &tamanio)
	if err != nil {
		return nil, errors.New("imagen no encontrada")
	}
	return modelos.NuevaImagen(id, nombre, descripcion, ruta, formato, tamanio)
}

func (r *ImagenRepositorioSQL) ObtenerPorNombre(nombre string) (*modelos.Imagen, error) {
	query := `SELECT id, nombre, descripcion, ruta, formato, tamanio FROM imagenes WHERE nombre = ?`
	var id int
	var descripcion, ruta, formato string
	var tamanio int64

	err := DB.QueryRow(query, nombre).Scan(&id, &nombre, &descripcion, &ruta, &formato, &tamanio)
	if err != nil {
		return nil, errors.New("imagen no encontrada")
	}
	return modelos.NuevaImagen(id, nombre, descripcion, ruta, formato, tamanio)
}

func (r *ImagenRepositorioSQL) Actualizar(imagen *modelos.Imagen) error {
	query := `UPDATE imagenes SET nombre = ?, descripcion = ? WHERE id = ?`
	result, err := DB.Exec(query,
		imagen.GetNombre(),
		imagen.GetDescripcion(),
		imagen.GetId(),
	)
	if err != nil {
		return err
	}
	filas, _ := result.RowsAffected()
	if filas == 0 {
		return errors.New("imagen no encontrada")
	}
	return nil
}

func (r *ImagenRepositorioSQL) Eliminar(id int) error {
	query := `DELETE FROM imagenes WHERE id = ?`
	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	filas, _ := result.RowsAffected()
	if filas == 0 {
		return errors.New("imagen no encontrada")
	}
	return nil
}
