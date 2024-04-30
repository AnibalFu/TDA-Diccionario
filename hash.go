package diccionario

import (
	"fmt"
)

type estado int

const (
	_VACIO = iota
	_BORRADO
	_OCUPADO
	_TAMANIO_INICIAL           = 11
	_FACTOR_REDIMENSION        = 2
	_FACTOR_CARGA              = 0.7
	_FNVOffset_Basis    uint32 = 2166136261
	_FNVPrime           uint32 = 16777619
)

type celdaHash[K comparable, V any] struct {
	clave  K
	dato   V
	estado estado
}

type hashCerrado[K comparable, V any] struct {
	tabla    []celdaHash[K, V]
	cantidad int //cant elementos
	borrados int //cant borrados
	tamaño   int //len de la tabla
}

type iteradorExterno[K comparable, V any] struct {
	tablaHash *hashCerrado[K, V]
}

func CrearHash[K comparable, V any]() *hashCerrado[K, V] {
	return &hashCerrado[K, V]{tabla: make([]celdaHash[K, V], _TAMANIO_INICIAL), tamaño: _TAMANIO_INICIAL}
}

func (hash *hashCerrado[K, V]) Guardar(clave K, valor V) {
	if _FACTOR_CARGA <= float64((hash.borrados+hash.cantidad)/hash.tamaño) {
		hash.redimension()
	}

	celda := *crearCeldaHash(clave, valor)
	indice := hashing(clave, hash.tamaño)

	for true {
		if hash.tabla[indice].estado == _VACIO {
			hash.tabla[indice] = celda
			hash.cantidad++
			break

		} else if hash.tabla[indice].estado == _OCUPADO && hash.tabla[indice].clave == clave {
			hash.tabla[indice].dato = valor
			break
		}

		indice++
		if indice == len(hash.tabla)-1 {
			indice = 0
		}

	}

}

func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	indice := hashing(clave, hash.tamaño)
	for _VACIO != hash.tabla[indice].estado {
		if clave == hash.tabla[indice].clave && _BORRADO != hash.tabla[indice].estado {
			return true
		}

		indice++
		if indice == len(hash.tabla)-1 {
			indice = 0
		}

	}

	return false
}

func (hash *hashCerrado[K, V]) Obtener(clave K) V {
	if !hash.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}

	indice := hashing(clave, hash.tamaño)
	for clave != hash.tabla[indice].clave || _BORRADO == hash.tabla[indice].estado {
		indice++
		if indice == len(hash.tabla)-1 {
			indice = 0
		}
	}

	return hash.tabla[indice].dato
}

func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	if !hash.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}
	indice := hashing(clave, hash.tamaño)

	for clave != hash.tabla[indice].clave || hash.tabla[indice].estado == _BORRADO {
		indice++
		if indice == len(hash.tabla)-1 {
			indice = 0
		}
	}

	hash.tabla[indice].estado = _BORRADO
	hash.borrados++
	hash.cantidad--
	return hash.tabla[indice].dato

}

func (hash *hashCerrado[K, V]) Cantidad() int {
	return hash.cantidad
}

func (hash *hashCerrado[K, V]) Iterar(func(clave K, dato V) bool) {
	return
}

func (hash *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	return crearIteradorExterno[K, V](hash)
}

func (iterador *iteradorExterno[K, V]) HaySiguiente() bool {
	return true
}

func (iterador *iteradorExterno[K, V]) VerActual() (K, V) {
	return iterador.tablaHash.tabla[0].clave, iterador.tablaHash.tabla[0].dato
}

func (iterador *iteradorExterno[K, V]) Siguiente() {
	return
}

func crearIteradorExterno[K comparable, V any](hash *hashCerrado[K, V]) *iteradorExterno[K, V] {
	return &iteradorExterno[K, V]{tablaHash: hash}
}

func crearCeldaHash[K comparable, V any](clave K, valor V) *celdaHash[K, V] {
	return &celdaHash[K, V]{clave: clave, dato: valor, estado: _OCUPADO}
}

func (hash *hashCerrado[K, V]) redimension() {
	nuevoTam := hash.tamaño * _FACTOR_REDIMENSION
	viejaTabla := hash.tabla

	hash.tabla = make([]celdaHash[K, V], nuevoTam)
	hash.tamaño = nuevoTam
	hash.borrados = 0
	hash.cantidad = 0

	for i := 0; i < len(viejaTabla); i++ {
		if viejaTabla[i].estado == _OCUPADO {
			hash.Guardar(viejaTabla[i].clave, viejaTabla[i].dato)
		}
	}
}

// https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func hashing[K comparable](clave K, capacidad int) int {
	hash := _FNVOffset_Basis
	bytes := convertirABytes(clave)
	for _, b := range bytes {
		hash ^= uint32(b)
		hash *= _FNVPrime
	}

	return int(hash) % capacidad
}
