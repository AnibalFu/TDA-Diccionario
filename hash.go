package diccionario

import (
	"fmt"
)

type estado int

const (
	_VACIO = iota
	_BORRADO
	_OCUPADO
	_TAMAÑO_INICIAL            = 13
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

	indice := buscarIndex(hash, clave)
	// Caso clave ya existe.
	if hash.tabla[indice].estado != _OCUPADO {
		hash.cantidad++
	}

	hash.tabla[indice] = crearCeldaHash(clave, valor)

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
			if !visitar(celda.clave, celda.dato) {
				break
			}
		}
	}
}

func (hash *hashCerrado[K, V]) Iterador() IterDiccionario[K, V] {
	return crearIteradorExterno(hash)
}

func (iterador *iteradorExterno[K, V]) HaySiguiente() bool {
	// Si el indice supera al tamaño de la tabla hash significa que estoy al final
	// y si la cant de elementos es 0 entonces no hay siguiente para ver.
	return iterador.iterIndice < iterador.iterHash.tamaño && iterador.iterHash.cantidad > 0
}

func (iterador *iteradorExterno[K, V]) VerActual() (K, V) {
	if !iterador.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	celda := iterador.iterHash.tabla[iterador.iterIndice]
	return celda.clave, celda.dato
}

func (iterador *iteradorExterno[K, V]) Siguiente() {
	if !iterador.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	iterador.iterIndice = buscarOcupado(iterador.iterHash, iterador.iterIndice+1)
}

// ///////////////////////////////////
//
// Funciones y metodos auxiliares.
//
// ///////////////////////////////////

// Crear iterador establece como atributos el hash y la primera aparicion de un elemento (celda ocupada)
func crearIteradorExterno[K comparable, V any](hash *hashCerrado[K, V]) *iteradorExterno[K, V] {
	return &iteradorExterno[K, V]{iterHash: hash, iterIndice: buscarOcupado(hash, 0)}
}

func crearCeldaHash[K comparable, V any](clave K, valor V) celdaHash[K, V] {
	return celdaHash[K, V]{clave: clave, dato: valor, estado: _OCUPADO}
}

// Primitiva de redimension, vuelve a hashear toda la tabla si se requiere una
// redimensionde la misma, ignorando los vacios y los borrados.
func (hash *hashCerrado[K, V]) redimension(nuevoTam int) {
	viejaTabla := hash.tabla

	//Establecer nuevos valores.
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
	// Busco el primer ocupado o la celda que este ocupado y tenga la misma clave que la que estoy recibiendo.
	for hash.tabla[indice].estado != _VACIO && (hash.tabla[indice].estado == _BORRADO || hash.tabla[indice].clave != clave) {
		indice = (indice + 1) % hash.tamaño // Me aseguro de que el indice siempre este dentro del len de la tabla.

	}

	return indice
}

// Busco el primer ocupado apartir del indice pasado como parametro.
func buscarOcupado[K comparable, V any](hash *hashCerrado[K, V], inicio int) int {
	for i := inicio; i < hash.tamaño; i++ {
		if hash.tabla[i].estado == _OCUPADO {
			break
		}
		inicio++
	}

	return inicio
}
