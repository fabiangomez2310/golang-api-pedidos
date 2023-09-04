package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type pedido struct { // se crea la estructura de pedido
	ID          int     `json:"id"`
	Descripcion string  `json:"descripcion"`
	Valor       float64 `json:"valor"`
	Estado      string  `json:"estado"`
}

var pedidos = []pedido{} // se crea el slice para almacenar la informacion en la memoria

func getPedidoByID(c *gin.Context) { // Funcion que recibe un id y devuelve los datos correspondientes

	id := c.Param("id")            // el parametro recibido trae el "id" pero con el tipo de dato string
	idInt, err := strconv.Atoi(id) // se realiza la conversion con la funcion strconv.Atoi, esta returna el numero o el error

	if err != nil { // en caso de haber generado un error por motivo de que la funcion no pudo convertir el daton entonces se genera una repuesta serializada
		c.JSON(http.StatusBadRequest, gin.H{
			"ERROR": "El ID proporcionado no es un numero valido",
		})
	}

	for _, p := range pedidos { // caso contrario entonces se empieza a iterar el slice pedidos
		if p.ID == idInt { // se crea un if para evaluar si el dato actual p.ID es igual al dato convertido idInt
			c.JSON(http.StatusOK, p) // en caso de encontrarlo se devuelve una respuesta http 200 y el dato actual que corresponde con el dato buscado.
			return                   // se retorna para evitar que la funcion se siga ejecutando
		}
	}

	// en caso tal que el datos no haya sido encontrado entonces se devuelve una respuesta indicando lo sucedido.
	c.JSON(http.StatusNotFound, gin.H{
		"ERROR": "El pedido no existe",
	})

}

func getPedidos(c *gin.Context) { // funcion que devuelve todos los datos del slice pedidos

	if len(pedidos) == 0 { // se realiza la validacion para saber si el slice esta vacio, en caso afirmativo se devuelve una respuesta informando lo sucedido.
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No se encontraron pedidos",
		})
		return // se retorna para evitar que la funcion se siga ejecutando
	}

	c.IndentedJSON(http.StatusOK, pedidos) // en caso contrario se devuelve una respuesta con la informacion, c.IndentedJSON lo que hace es serializar los datos del slice pedido.

}

func postPedidos(c *gin.Context) { // funcion que recibe unos parametros y crea un pedido

	var newPedido pedido // se crea una variable newPedido de la estructura pedido

	if err := c.BindJSON(&newPedido); err != nil { // se pasan los parametros a la variable newPedido con la funcion c.BindJSON es por eso que le pasamos la referencia de la memoria, en caso tal de suceder un error entonces se retorna la funcion para evitar seguir ejecutandola.
		return //
	}

	pedidos = append(pedidos, newPedido) // se agrega los datos de la variable newPedido al slice pedidos, esto se hace con la funcion append.

	c.IndentedJSON(http.StatusCreated, newPedido) // se devuelve una respuesta con el StatusCreated 201 y los datos del nuevo pedido.

}

func putPedido(c *gin.Context) { // funcion que recibe los parametros de la nuevo informacion del pedido, en estos parametros viene el ID el cual sera usado para realizar la busqueda.
	id := c.Param("id")            // entre los parametros de la nuevo informacion del pedido se obtiene especificamente el "id" el cual se almacena en la variable id..
	idInt, err := strconv.Atoi(id) // por medio de la funcion strconv.Atoi se realiza la conversion de strign a int, esta funcion devuelve un error en caso que no sea un ID compatible para la conversion.

	if err != nil { // se evalua si err trae algun valor diferente de nil
		c.JSON(http.StatusBadRequest, gin.H{ // si el valor es diferente de nil quiere decir que sucedio un error y se ejecuta la respuesta donde se devuelve el status correspondiente y un mensaje de lo sucedido.
			"ERROR": "El ID proporcionado no es un numero valido",
		})

		return // se retorna la funcion para evitar seguir ejecutando.
	}

	for i, p := range pedidos { // en caso de no haber sucedido lo anterior entonces vamos a recorrer el slice pedidos con el animo de encontrar la siguiente condicion.
		if p.ID == idInt { // si la condicion p.ID == idInt entonces quiere decir que se encontro el Id
			var newPedido = pedido{}                       // se crea entonces una instancia de la estructura pedidos
			if err := c.BindJSON(&newPedido); err != nil { // c.bindJSON lo que hace es agregar la informacion a la varible newPedido por esa razon se le pasa la variable por referencia.
				return // se retorna la funcion
			}
			pedidos[i] = newPedido // en el slice pedidos posicion i se almacena la informacion guardada en newPedido para con esto actualizar todos los datos de la ubicacion i

			c.IndentedJSON(http.StatusOK, newPedido) // se devuelve la respueta con el status ok y la informacion del pedido actualizado.
			return                                   // se retorna la funcion.
		}
	}

	c.JSON(http.StatusNotFound, gin.H{ // en caso tal no haber encontrado el id entonces devuelve el siguiente mensaje
		"ERROR": "El pedido no existe",
	})
}

func deletePedido(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"ERROR": "El ID proporcionado no es un numero valido",
		})
		return // SALIR DE LA FUNCION
	}

	for i, p := range pedidos {
		if p.ID == idInt {
			pedidos = append(pedidos[:i], pedidos[i+1:]...)
			c.JSON(http.StatusNoContent, gin.H{
				"message": "Pedido eliminado",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"ERROR": "El pedido no existe",
	})
}

func main() {
	router := gin.Default()

	router.GET("/pedidos", getPedidos)
	router.GET("/pedidos/:id", getPedidoByID)
	router.POST("/pedidos", postPedidos)
	router.PUT("/pedidos/:id", putPedido)
	router.DELETE("/pedidos/:id", deletePedido)

	router.Run(":8080")

}
