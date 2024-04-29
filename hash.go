package diccionario

const (
	VACIO        = -1
	BORRADO      = 0
	OCUPADO      = 1
	TAMANIOTABLA = 10
)

type celdaHash[K comparable, V any] struct {
	clave  K
	dato   V
	estado int
}

type hashCerrado[K comparable, V any] struct {
	tabla    []celdaHash[K, V]
	cantidad int
	tam      int
	borrados int
}

type iteradorExterno[K comparable, V any] struct {
	tablaHash *hashCerrado[K, V]
}

func CrearHash[K comparable, V any]() *hashCerrado[K, V] {
	return &hashCerrado[K, V]{tabla: make([]celdaHash[K, V], TAMANIOTABLA)}
}

func crearIteradorExterno[K comparable, V any](hash *hashCerrado[K, V]) *iteradorExterno[K, V] {
	return &iteradorExterno[K, V]{tablaHash: hash}
}

func (hash *hashCerrado[K, V]) Guardar(clave K, dato V) {

}

func (hash *hashCerrado[K, V]) Pertenece(clave K) bool {
	return true
}

func (hash *hashCerrado[K, V]) Obtener(clave K) V {
	return hash.tabla[0].dato
}

func (hash *hashCerrado[K, V]) Borrar(clave K) V {
	return hash.tabla[0].dato
}

func (hash *hashCerrado[K, V]) Cantidad() int {
	return 0
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
