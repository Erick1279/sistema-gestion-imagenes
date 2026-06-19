package modelos

import "errors"

// Struct principal - POO: encapsulación con campos en minúscula
type Imagen struct {
	id          int
	nombre      string
	descripcion string
	ruta        string
	formato     string
	tamanio     int64
}

// Constructor
func NuevaImagen(id int, nombre string, descripcion string, ruta string, formato string, tamanio int64) (*Imagen, error) {
	if nombre == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if ruta == "" {
		return nil, errors.New("la ruta no puede estar vacía")
	}
	return &Imagen{
		id:          id,
		nombre:      nombre,
		descripcion: descripcion,
		ruta:        ruta,
		formato:     formato,
		tamanio:     tamanio,
	}, nil
}

// Getters
func (i *Imagen) GetId() int             { return i.id }
func (i *Imagen) GetNombre() string      { return i.nombre }
func (i *Imagen) GetDescripcion() string { return i.descripcion }
func (i *Imagen) GetRuta() string        { return i.ruta }
func (i *Imagen) GetFormato() string     { return i.formato }
func (i *Imagen) GetTamanio() int64      { return i.tamanio }

// Setters
func (i *Imagen) SetNombre(nombre string) error {
	if nombre == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	i.nombre = nombre
	return nil
}

func (i *Imagen) SetDescripcion(descripcion string) {
	i.descripcion = descripcion
}

func (i *Imagen) SetRuta(ruta string) error {
	if ruta == "" {
		return errors.New("la ruta no puede estar vacía")
	}
	i.ruta = ruta
	return nil
}
