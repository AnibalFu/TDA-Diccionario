package diccionario

import (
	cola "tdas/cola"
)

type funcCmp[K comparable] func(K, K) int

type nodoAbb[K comparable, V any] struct {
	izq   *nodoAbb[K, V]
	der   *nodoAbb[K, V]
	clave K
	dato  V
}

type abb[K comparable, V any] struct {
	raiz     *nodoAbb[K, V]
	cantidad int
	cmp      funcCmp[K]
}

type iteradorRangoAbb[K comparable, V any] struct {
	iterAbb *abb[K, V]
	desde   *K
	hasta   *K
}

func CrearAbb[K comparable, V any](funcion_cmp func(K, K) int) DiccionarioOrdenado[K, V] {
	return &abb[K, V]{cmp: funcion_cmp}
}

func (ab *abb[K, V]) Pertenece(clave K) bool {
	nodo := buscarPunteroNodo(ab, &ab.raiz, clave)
	return (*nodo) != nil
}

func (ab *abb[K, V]) Obtener(clave K) V {
	nodo := buscarPunteroNodo(ab, &ab.raiz, clave)
	if (*nodo) == nil {
		panic("La clave no pertenece al diccionario")
	}

	return (*nodo).dato
}

func (ab *abb[K, V]) Guardar(clave K, valor V) {
	// Busco la referencia a la referencia donde deberia de ir la clave.
	nodo := buscarPunteroNodo(ab, &ab.raiz, clave)

	// Caso guardo elemento nuevo, si ya esta no deberia sumar 1 a cantidad
	if *nodo == nil {
		ab.cantidad++
		*nodo = crearNodoAbb(clave, valor)

	} else {
		(*nodo).dato = valor

	}
}

func (ab *abb[K, V]) Borrar(clave K) V {
	nodo := buscarPunteroNodo(ab, &ab.raiz, clave)

	if *nodo == nil {
		panic("La clave no pertenece al diccionario")
	}

	ab.cantidad--
	dato := (*nodo).dato

	// Veo en que caso de borrado estoy, 0 hijos, 1 hijo, 2 hijos.
	// Caso 2 hijos, busco el reemplazo, reemplazo los datos y borro el reemplazo.
	if (*nodo).izq != nil && (*nodo).der != nil {
		reemplazo := buscarReemplazo(&(*nodo).izq)
		(*nodo).clave, (*nodo).dato = (*reemplazo).clave, (*reemplazo).dato
		_borrar(reemplazo)

		// Caso 0 hijos o 1 hijo.
	} else {
		_borrar(nodo)

	}

	return dato

}

func (ab *abb[K, V]) Cantidad() int {
	return ab.cantidad
}

func (ab *abb[K, V]) Iterar(visitar func(clave K, valor V) bool) {
	_iterar(ab, ab.raiz, nil, nil, visitar)
}

func (ab *abb[K, V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool) {
	_iterar(ab, ab.raiz, desde, hasta, visitar)
}

func (ab *abb[K, V]) Iterador() IterDiccionario[K, V] {
	return &iteradorExterno[K, V]{}
}

func (ab *abb[K, V]) IteradorRango(desde *K, hasta *K) IterDiccionario[K, V] {
	return &iteradorExterno[K, V]{}
}

func (iterador *iteradorRangoAbb[K, V]) HaySiguiente() bool {
	return true
}

func (iterador *iteradorRangoAbb[K, V]) VerActual() (K, V) {
	return iterador.iterAbb.raiz.clave, iterador.iterAbb.raiz.dato
}

func (iterador *iteradorRangoAbb[K, V]) Siguiente() {

}

/////////////////////////////////////
//
// Funciones y metodos auxiliares
//
/////////////////////////////////////

func crearNodoAbb[K comparable, V any](clave K, valor V) *nodoAbb[K, V] {
	return &nodoAbb[K, V]{clave: clave, dato: valor}
}

// Busca la direccion y devuelve el puntero del puntero que le correnponderia a la clave pasada por parametro.
func buscarPunteroNodo[K comparable, V any](ab *abb[K, V], nodo **nodoAbb[K, V], clave K) **nodoAbb[K, V] {
	if *nodo == nil {
		return nodo
	}

	// Veo para que lado realizar la sig busqueda o si ya estoy en el correcto.
	resCmp := ab.cmp((*nodo).clave, clave)
	if resCmp == 0 {
		return nodo

	} else if resCmp > 0 {
		// Llamo recursivamente pasando la direccion de memoria del nodo izq del nodo actual.
		return buscarPunteroNodo(ab, &(*nodo).izq, clave)

	} else {
		// Llamo recursivamente pasando la direccion de memoria del nodo der del nodo actual.
		return buscarPunteroNodo(ab, &(*nodo).der, clave)

	}
}

// //////////////////////////////////////////////
func (ab *abb[K, V]) RecorrerPorNiveles() []K {
	res := []K{}
	if ab.raiz == nil {
		return res
	}

	c := cola.CrearColaEnlazada[*nodoAbb[K, V]]()
	c.Encolar(ab.raiz)

	for !c.EstaVacia() {
		nodo := c.Desencolar() // Desencolamos el nodo actual

		// Ejecutamos la función visitar en el nodo actual
		res = append(res, nodo.clave)

		// Encolamos los hijos izquierdo y derecho del nodo actual si existen
		if nodo.izq != nil {
			c.Encolar(nodo.izq)
		}
		if nodo.der != nil {
			c.Encolar(nodo.der)
		}
	}

	return res
}

// La funcion recibe un nodo con 0 hijos o 1 hijo y borra el nodo de forma correcta.
func _borrar[K comparable, V any](nodo **nodoAbb[K, V]) {
	hijoIzq, hijoDer := (*nodo).izq, (*nodo).der

	if hijoIzq == nil && hijoDer != nil {
		*nodo = hijoDer

	} else if hijoDer == nil && hijoIzq != nil {
		*nodo = hijoIzq

	} else {
		*nodo = nil

	}
}

// Dado un nodo que tiene 2 hijos y lo quiero borrar, busco el reeemplazo correcto,
// devolviendo el puntero del puntero del nodo que corresponde.
func buscarReemplazo[K comparable, V any](nodo **nodoAbb[K, V]) **nodoAbb[K, V] {
	if (*nodo).der == nil {
		return nodo
	}

	return buscarReemplazo(&(*nodo).der)

}

func _iterar[K comparable, V any](abb *abb[K, V], nodo *nodoAbb[K, V], desde *K, hasta *K, visitar func(clave K, dato V) bool) {
	if nodo == nil {
		return
	}

	// Si desde es nil, iterar desde la primera clave.
	if desde == nil || abb.cmp(nodo.clave, *desde) >= 0 {
		_iterar(abb, nodo.izq, desde, hasta, visitar)
	}

	// Si estamos en el rango se debe visitar el nodo actual.
	if (desde == nil || abb.cmp(nodo.clave, *desde) >= 0) && (hasta == nil || abb.cmp(nodo.clave, *hasta) <= 0) {
		if !visitar(nodo.clave, nodo.dato) {
			return
		}
	}

	// Si hasta es nil, iterar hasta la última clave.
	if hasta == nil || abb.cmp(nodo.clave, *hasta) <= 0 {
		_iterar(abb, nodo.der, desde, hasta, visitar)
	}
}
