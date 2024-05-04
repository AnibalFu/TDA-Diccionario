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

func CrearABB[K comparable, V any](funcion_cmp func(K, K) int) DiccionarioOrdenado[K, V] {
	return &abb[K, V]{cmp: funcion_cmp}
}

func (ab *abb[K, V]) Pertenece(clave K) bool {
	nodoAct, _ := buscarActYAnt(ab, ab.raiz, ab.raiz, clave)
	return nodoAct != nil
}

func (ab *abb[K, V]) Obtener(clave K) V {
	nodoAct, _ := buscarActYAnt(ab, ab.raiz, ab.raiz, clave)
	if nodoAct == nil {
		panic("La clave no pertenece al diccionario")
	}

	return nodoAct.dato
}

func (ab *abb[K, V]) Guardar(clave K, valor V) {
	nuevoNodo := crearNodoAbb(clave, valor)
	_, nodoAnt := buscarActYAnt(ab, ab.raiz, ab.raiz, clave)

	// Caso arbol sin nodos
	if nodoAnt == nil {
		ab.raiz = nuevoNodo
	} else {
		resCmp := ab.cmp(nodoAnt.clave, clave)

		// Caso clave ya existe
		if resCmp == 0 {
			nodoAnt.dato = valor
			return

		} else if resCmp > 0 {
			nodoAnt.izq = nuevoNodo

		} else {
			nodoAnt.der = nuevoNodo

		}
	}

	ab.cantidad++

}

func (ab *abb[K, V]) Borrar(clave K) V {
	return ab.raiz.dato
}

func (ab *abb[K, V]) Cantidad() int {
	return ab.cantidad
}

func (ab *abb[K, V]) Iterar(visitar func(clave K, valor V) bool) {

}

func (ab *abb[K, V]) IterarRango(desde *K, hasta *K, visitar func(clave K, dato V) bool) {

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

// Busca el nodo actual y su anterior dado una clave pasada por parametro, devuelve el nodo actual y anterior
func buscarActYAnt[K comparable, V any](ab *abb[K, V], nodoAct, nodoAnt *nodoAbb[K, V], clave K) (*nodoAbb[K, V], *nodoAbb[K, V]) {
	if nodoAct == nil {
		return nodoAct, nodoAnt
	}

	resCmp := ab.cmp(nodoAct.clave, clave)

	if resCmp == 0 {
		return nodoAct, nodoAnt

	} else if resCmp > 0 {
		return buscarActYAnt(ab, nodoAct.izq, nodoAct, clave)

	} else {
		return buscarActYAnt(ab, nodoAct.der, nodoAct, clave)
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
