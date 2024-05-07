package diccionario_test

import (
	TDAAbb "tdas/diccionario"
	"testing"

	"github.com/stretchr/testify/require"
)

func CompararStrings(a, b string) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func CompararEnteros(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func TestDiccionarioAbbVacio(t *testing.T) {
	t.Log("Comprueba que Diccionario vacio no tiene claves")
	dic := TDAAbb.CrearAbb[string, string](CompararStrings)
	require.EqualValues(t, 0, dic.Cantidad())
	require.False(t, dic.Pertenece("A"))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener("A") })
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar("A") })
}

func TestBorrarConCeroHijos(t *testing.T) {
	dic := TDAAbb.CrearAbb[int, int](CompararEnteros)

	dic.Guardar(5, 5)
	dic.Guardar(2, 2)
	dic.Guardar(6, 6)

	require.EqualValues(t, 6, dic.Borrar(6))
	require.False(t, dic.Pertenece(6))

	dic.Guardar(6, 6)
	dic.Guardar(7, 7)
	require.EqualValues(t, 6, dic.Borrar(6))
	require.False(t, dic.Pertenece(6))
	require.EqualValues(t, 3, dic.Cantidad())

}

func TestBorrarConUnHijos(t *testing.T) {
	dic := TDAAbb.CrearAbb[int, int](CompararEnteros)

	dic.Guardar(5, 5)
	dic.Guardar(7, 7)
	dic.Guardar(6, 6)

	require.True(t, dic.Pertenece(6))
	require.EqualValues(t, 7, dic.Borrar(7))
	require.False(t, dic.Pertenece(7))

	dic.Guardar(7, 7)
	require.True(t, dic.Pertenece(6))
	require.EqualValues(t, 6, dic.Borrar(6))
	require.False(t, dic.Pertenece(7))

}

func TestBorrarConDosHijos(t *testing.T) {
	dic := TDAAbb.CrearAbb[int, int](CompararEnteros)

	dic.Guardar(10, 10)
	dic.Guardar(4, 4)
	dic.Guardar(5, 5)
	dic.Guardar(1, 1)
	dic.Guardar(13, 13)
	dic.Guardar(11, 11)
	dic.Guardar(12, 12)
	dic.Guardar(14, 14)

	require.EqualValues(t, 4, dic.Borrar(4))
	require.False(t, dic.Pertenece(4))
	require.EqualValues(t, 13, dic.Borrar(13))
	require.False(t, dic.Pertenece(13))
	require.EqualValues(t, 10, dic.Borrar(10))
	require.False(t, dic.Pertenece(10))
	require.EqualValues(t, 5, dic.Cantidad())

}

func TestDiccionarioAbbClaveDefault(t *testing.T) {
	t.Log("Prueba sobre un Hash vacío que si justo buscamos la clave que es el default del tipo de dato, " +
		"sigue sin existir")
	dic := TDAAbb.CrearAbb[string, string](CompararStrings)
	require.False(t, dic.Pertenece(""))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener("") })
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Borrar("") })

	dicNum := TDAAbb.CrearAbb[int, string](CompararEnteros)
	require.False(t, dicNum.Pertenece(0))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dicNum.Obtener(0) })
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dicNum.Borrar(0) })
}

func TestUnElementAbb(t *testing.T) {
	t.Log("Comprueba que Diccionario con un elemento tiene esa Clave, unicamente")
	dic := TDAAbb.CrearAbb[string, int](CompararStrings)
	dic.Guardar("A", 10)
	require.EqualValues(t, 1, dic.Cantidad())
	require.True(t, dic.Pertenece("A"))
	require.False(t, dic.Pertenece("B"))
	require.EqualValues(t, 10, dic.Obtener("A"))
	require.PanicsWithValue(t, "La clave no pertenece al diccionario", func() { dic.Obtener("B") })
}

func TestDiccionarioAbbGuardar(t *testing.T) {
	t.Log("Guarda algunos pocos elementos en el diccionario, y se comprueba que en todo momento funciona acorde")
	clave1 := "Gato"
	clave2 := "Perro"
	clave3 := "Vaca"
	valor1 := "miau"
	valor2 := "guau"
	valor3 := "moo"
	claves := []string{clave1, clave2, clave3}
	valores := []string{valor1, valor2, valor3}

	dic := TDAAbb.CrearAbb[string, string](CompararStrings)
	require.False(t, dic.Pertenece(claves[0]))
	require.False(t, dic.Pertenece(claves[0]))
	dic.Guardar(claves[0], valores[0])
	require.EqualValues(t, 1, dic.Cantidad())
	require.True(t, dic.Pertenece(claves[0]))
	require.True(t, dic.Pertenece(claves[0]))
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))

	require.False(t, dic.Pertenece(claves[1]))
	require.False(t, dic.Pertenece(claves[2]))
	dic.Guardar(claves[1], valores[1])
	require.True(t, dic.Pertenece(claves[0]))
	require.True(t, dic.Pertenece(claves[1]))
	require.EqualValues(t, 2, dic.Cantidad())
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))
	require.EqualValues(t, valores[1], dic.Obtener(claves[1]))

	require.False(t, dic.Pertenece(claves[2]))
	dic.Guardar(claves[2], valores[2])
	require.True(t, dic.Pertenece(claves[0]))
	require.True(t, dic.Pertenece(claves[1]))
	require.True(t, dic.Pertenece(claves[2]))
	require.EqualValues(t, 3, dic.Cantidad())
	require.EqualValues(t, valores[0], dic.Obtener(claves[0]))
	require.EqualValues(t, valores[1], dic.Obtener(claves[1]))
	require.EqualValues(t, valores[2], dic.Obtener(claves[2]))
}

func TestReemplazoDatoAbb(t *testing.T) {
	t.Log("Guarda un par de claves, y luego vuelve a guardar, buscando que el dato se haya reemplazado")
	clave := "Gato"
	clave2 := "Perro"
	dic := TDAAbb.CrearAbb[string, string](CompararStrings)
	dic.Guardar(clave, "miau")
	dic.Guardar(clave2, "guau")
	require.True(t, dic.Pertenece(clave))
	require.True(t, dic.Pertenece(clave2))
	require.EqualValues(t, "miau", dic.Obtener(clave))
	require.EqualValues(t, "guau", dic.Obtener(clave2))
	require.EqualValues(t, 2, dic.Cantidad())

	dic.Guardar(clave, "miu")
	dic.Guardar(clave2, "baubau")
	require.True(t, dic.Pertenece(clave))
	require.True(t, dic.Pertenece(clave2))
	require.EqualValues(t, 2, dic.Cantidad())
	require.EqualValues(t, "miu", dic.Obtener(clave))
	require.EqualValues(t, "baubau", dic.Obtener(clave2))
}

func TestReemplazoDatoHopscotchAbb(t *testing.T) {
	t.Log("Guarda bastantes claves, y luego reemplaza sus datos. Luego valida que todos los datos sean " +
		"correctos. Para una implementación Hopscotch, detecta errores al hacer lugar o guardar elementos.")

	dic := TDAAbb.CrearAbb[int, int](CompararEnteros)
	for i := 0; i < 500; i++ {
		dic.Guardar(i, i)
	}
	for i := 0; i < 500; i++ {
		dic.Guardar(i, 2*i)
	}
	ok := true
	for i := 0; i < 500 && ok; i++ {
		ok = dic.Obtener(i) == 2*i
	}
	require.True(t, ok, "Los elementos no fueron actualizados correctamente")
}

func TestIteradorInterno(t *testing.T) {
	t.Log("Valida condicion de corte del iterador interno cuando un elemento no cumple con la funcion visitar")

	dic := TDAAbb.CrearAbb[int, int](CompararEnteros)
	dic.Guardar(7, 7)
	dic.Guardar(6, 6)
	dic.Guardar(2, 2)
	dic.Guardar(3, 3)
	dic.Guardar(4, 4)
	dic.Guardar(5, 5)

	res := 0
	dic.Iterar(func(clave int, dato int) bool {
		if clave <= 5 {
			res += clave
		}
		return true

	})
	require.EqualValues(t, 14, res)

	res = 0
	dic.Iterar(func(_ int, dato int) bool {
		if res > 10 {
			return false
		}
		res += dato
		return true

	})
	require.EqualValues(t, 14, res)

	res = 0
	dic.Iterar(func(clave int, dato int) bool {
		if clave > 3 {
			return false
		}
		res += dato
		return true

	})
	require.EqualValues(t, 5, res)

}

func TestIteradorInternoConBorrados(t *testing.T) {
	t.Log("Valida que los datos sean recorridas correctamente (y una única vez) con el iterador interno, sin recorrer datos borrados")
	clave0 := "Elefante"
	clave1 := "Gato"
	clave2 := "Perro"
	clave3 := "Vaca"
	clave4 := "Burrito"
	clave5 := "Hamster"

	dic := TDAAbb.CrearAbb[string, int](CompararStrings)
	dic.Guardar(clave0, 7)
	dic.Guardar(clave1, 6)
	dic.Guardar(clave2, 2)
	dic.Guardar(clave3, 3)
	dic.Guardar(clave4, 4)
	dic.Guardar(clave5, 5)

	dic.Borrar(clave0)

	factorial := 1
	ptrFactorial := &factorial
	dic.Iterar(func(_ string, dato int) bool {
		*ptrFactorial *= dato
		return true
	})

	require.EqualValues(t, 720, factorial)
}

func TestIterarRangoSuma(t *testing.T) {
	dic := TDAAbb.CrearAbb[int, int](CompararEnteros)

	dic.Guardar(2, 2)
	dic.Guardar(10, 10)
	dic.Guardar(11, 11)
	dic.Guardar(5, 5)
	dic.Guardar(6, 6)
	dic.Guardar(3, 3)
	dic.Guardar(7, 7)

	var tres int = 3
	var siete int = 7

	res := 0
	dic.IterarRango(&tres, &siete, func(clave, dato int) bool {
		res += dato
		return true
	})

	require.EqualValues(t, 21, res)
}

func TestIterarRangoMixtoSuma(t *testing.T) {
	dic := TDAAbb.CrearAbb[int, int](CompararEnteros)

	dic.Guardar(2, 2)
	dic.Guardar(10, 10)
	dic.Guardar(11, 11)
	dic.Guardar(5, 5)
	dic.Guardar(6, 6)
	dic.Guardar(3, 3)
	dic.Guardar(7, 7)

	var tres int = 3
	var siete int = 7

	res := 0
	dic.IterarRango(nil, &siete, func(clave, dato int) bool {
		res += dato
		return true
	})

	require.EqualValues(t, 23, res)

	res = 0
	dic.IterarRango(&tres, nil, func(clave, dato int) bool {
		res += dato
		return true
	})

	require.EqualValues(t, 42, res)
}
