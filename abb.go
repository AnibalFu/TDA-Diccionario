package diccionario

import (
	TDAPila "tdas/pila"
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
	iterAbb  *abb[K, V]
	iterPila TDAPila.Pila[nodoAbb[K, V]]
	desde    *K
	hasta    *K
}

func CrearABB[K comparable, V any](funcion_cmp func(K, K) int) DiccionarioOrdenado[K, V] {
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
	continuar := true
	_iterar(ab, ab.raiz, nil, nil, visitar, &continuar)
}

func (ab *abb[K, V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool) {
	continuar := true
	_iterar(ab, ab.raiz, desde, hasta, visitar, &continuar)
}
func (ab *abb[K, V]) Iterador() IterDiccionario[K, V] {
	return ab.IteradorRango(nil, nil)
}

func (ab *abb[K, V]) IteradorRango(desde *K, hasta *K) IterDiccionario[K, V] {
	return &iteradorRangoAbb[K, V]{iterAbb: ab, desde: desde, hasta: hasta, iterPila: ab.crearIterador(desde, hasta)}
}

func (iterador *iteradorRangoAbb[K, V]) HaySiguiente() bool {
	return !iterador.iterPila.EstaVacia()
}

func (iterador *iteradorRangoAbb[K, V]) VerActual() (K, V) {
	if !iterador.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	return iterador.iterPila.VerTope().clave, iterador.iterPila.VerTope().dato
}

func (iterador *iteradorRangoAbb[K, V]) Siguiente() {
	if !iterador.HaySiguiente() {
		panic("El iterador termino de iterar")
	}

	nodo := iterador.iterPila.Desapilar()
	iterApilar(iterador.iterAbb, iterador.iterPila, nodo.der, iterador.desde, iterador.hasta)

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

// La funcion recibe un nodo con 0 hijos o 1 hijo y borra el nodo de forma correcta.
func _borrar[K comparable, V any](nodo **nodoAbb[K, V]) {
	hijoIzq, hijoDer := (*nodo).izq, (*nodo).der

	if hijoDer != nil {
		*nodo = hijoDer

	} else if hijoIzq != nil {
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

// Iterar de forma inorder dado un rango.
func _iterar[K comparable, V any](abb *abb[K, V], nodo *nodoAbb[K, V], desde *K, hasta *K, visitar func(clave K, dato V) bool, continuar *bool) {
	if nodo == nil || !*continuar {
		return
	}

	// Si desde es nil, iterar desde la primera clave.
	if desde == nil || abb.cmp(nodo.clave, *desde) >= 0 {
		_iterar(abb, nodo.izq, desde, hasta, visitar, continuar)
	}

	// Si estamos en el rango se debe visitar el nodo actual pero si continual es false no debo visitar.
	if *continuar && (desde == nil || abb.cmp(nodo.clave, *desde) >= 0) && (hasta == nil || abb.cmp(nodo.clave, *hasta) <= 0) {
		// Visitar el nodo y actualizar continuar si es necesario.
		*continuar = visitar(nodo.clave, nodo.dato)
	}

	// Si hasta es nil, iterar hasta la ultima clave.
	if hasta == nil || abb.cmp(nodo.clave, *hasta) <= 0 {
		_iterar(abb, nodo.der, desde, hasta, visitar, continuar)
	}
}

func (ab *abb[K, V]) crearIterador(desde, hasta *K) TDAPila.Pila[nodoAbb[K, V]] {
	pila := TDAPila.CrearPilaDinamica[nodoAbb[K, V]]()
	iterApilar(ab, pila, ab.raiz, desde, hasta)
	return pila
}

// La funcion
func iterApilar[K comparable, V any](abb *abb[K, V], pila TDAPila.Pila[nodoAbb[K, V]], nodo *nodoAbb[K, V], desde, hasta *K) {
	if nodo == nil {
		return
	}

	// Si estamos en el rango se debe apilar el nodo actual e ir a la izquierda.
	if (desde == nil || (abb.cmp(nodo.clave, *desde) >= 0)) && (hasta == nil || abb.cmp(nodo.clave, *hasta) <= 0) {
		pila.Apilar(*nodo)
		iterApilar(abb, pila, nodo.izq, desde, hasta)

		// Ir a la derecha si la clave del actual es mayor al desde.
	} else if desde != nil && abb.cmp(nodo.clave, *desde) < 0 {
		iterApilar(abb, pila, nodo.der, desde, hasta)

		// Ir a la izquierda si la clave del actual es mejor al hasta.
	} else if hasta != nil && abb.cmp(nodo.clave, *hasta) > 0 {
		iterApilar(abb, pila, nodo.izq, desde, hasta)

	}
}
