package diccionario

import (
	"fmt"
)

type estado int

const (
	_VACIO = iota
	_BORRADO
	_OCUPADO
	_TAMAÑO_INICIAL            = 10
	_FACTOR_REDIMENSION        = 2
	_FACTOR_CAPACIDAD          = 4
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
	iterHash   *hashCerrado[K, V]
	iterIndice int
}

func CrearHash[K comparable, V any]() *hashCerrado[K, V] {
	return &hashCerrado[K, V]{tabla: make([]celdaHash[K, V], _TAMAÑO_INICIAL), tamaño: _TAMAÑO_INICIAL}
}

func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	indice := buscarIndex(hash, clave)
	return hash.tabla[indice].estado == _OCUPADO

}

func (hash *hashCerrado[K, V]) Obtener(clave K) V {
	if !hash.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}

	indice := buscarIndex(hash, clave)
	return hash.tabla[indice].dato
}

func (hash *hashCerrado[K, V]) Guardar(clave K, valor V) {
	// Veo si necesita redimension.
	if _FACTOR_CARGA <= float64(hash.borrados+hash.cantidad)/float64(hash.tamaño) {
		hash.redimension(hash.tamaño * _FACTOR_REDIMENSION)
	}

	celda := crearCeldaHash(clave, valor)
	indice := buscarIndex(hash, clave)

	// Caso clave ya existe.
	if hash.tabla[indice].estado == _OCUPADO {
		hash.tabla[indice] = celda

	} else {
		hash.tabla[indice] = celda
		hash.cantidad++

	}
}

func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	if !hash.Pertenece(clave) {
		panic("La clave no pertenece al diccionario")
	}

	// Veo si la posible redimencion no me redimensione menos del tamaño inicial y si requiere redimension
	if hash.cantidad*_FACTOR_CAPACIDAD <= hash.tamaño && hash.tamaño > _TAMAÑO_INICIAL {
		hash.redimension(hash.tamaño / _FACTOR_REDIMENSION)
	}

	indice := buscarIndex(hash, clave)
	hash.tabla[indice].estado = _BORRADO
	hash.borrados++
	hash.cantidad--
	return hash.tabla[indice].dato

}

func (hash *hashCerrado[K, V]) Cantidad() int {
	return hash.cantidad
}

func (hash *hashCerrado[K, V]) Iterar(visitar func(clave K, valor V) bool) {
	for _, celda := range hash.tabla {
		if celda.estado == _OCUPADO {
			visito := visitar(celda.clave, celda.dato)
			if !visito {
				break
			}
		}
	}
}

func (hash *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	return crearIteradorExterno[K, V](hash)
}

func (iterador *iteradorExterno[K, V]) HaySiguiente() bool {
	return iterador.iterHash.tamaño > iterador.iterIndice && iterador.iterHash.cantidad > 0
}

func (iterador *iteradorExterno[K, V]) VerActual() (K, V) {
	if !iterador.HaySiguiente() {
		panic("El iterador termino de iterar")
	}
	return iterador.iterHash.tabla[iterador.iterIndice].clave, iterador.iterHash.tabla[iterador.iterIndice].dato
}

func (iterador *iteradorExterno[K, V]) Siguiente() {
	if !iterador.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	for iterador.HaySiguiente() {
		iterador.iterIndice++
		if iterador.HaySiguiente() && iterador.iterHash.tabla[iterador.iterIndice].estado == _OCUPADO {
			break
		}
	}
}

// ///////////////////////////////////
//
// Funciones y metodos auxiliares.
//
// ///////////////////////////////////
func crearIteradorExterno[K comparable, V any](hash *hashCerrado[K, V]) *iteradorExterno[K, V] {
	return &iteradorExterno[K, V]{iterHash: hash, iterIndice: buscarPrimerOcupado(hash)}
}

func crearCeldaHash[K comparable, V any](clave K, valor V) celdaHash[K, V] {
	return celdaHash[K, V]{clave: clave, dato: valor, estado: _OCUPADO}
}

// Primitiva de redimension, vuelve a hashear toda la tabla si se requiere una
// redimensionde la misma, ignorando los vacios y los borrados.
func (hash *hashCerrado[K, V]) redimension(nuevoTam int) {
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

// Funcion de hassing, fuente: https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
func convertirABytes[K comparable](clave K) []byte {
	return []byte(fmt.Sprintf("%v", clave))
}

func hashing[K comparable](clave K, tamaño int) int {
	hash := _FNVOffset_Basis
	bytes := convertirABytes(clave)
	for _, b := range bytes {
		hash ^= uint32(b)
		hash *= _FNVPrime
	}

	return int(hash) % tamaño
}

// Busco el indice correspondiente si es necesario para evitar colisiones.
func buscarIndex[K comparable, V any](hash *hashCerrado[K, V], clave K) int {
	indice := hashing(clave, hash.tamaño)
	for hash.tabla[indice].estado != _VACIO && (hash.tabla[indice].estado == _BORRADO || hash.tabla[indice].clave != clave) {
		indice = (indice + 1) % hash.tamaño // Me aseguro de que el indice siempre este dentro del len de la tabla.

	}

	return indice
}

// Busco el primer ocupado
func buscarPrimerOcupado[K comparable, V any](hash *hashCerrado[K, V]) int {
	indice := 0
	for i := 0; i < hash.tamaño; i++ {
		if hash.tabla[i].estado == _OCUPADO {
			break
		}
		indice++
	}

	return indice
}
