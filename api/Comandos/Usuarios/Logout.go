package Usuarios

import (
	"MIA_1S2025_P1_201708997/Structs"
	"fmt"
)

func Logout() string {
	var respuesta string
	if Structs.UsuarioActual.Status {
		Structs.SalirUsuario()
		fmt.Println("Se ha cerrado la sesion")
		respuesta += "Se ha cerrado la sesion"
	} else {
		respuesta += "ERROR LOGUT: NO HAY SECION INICIADA"
	}

	return respuesta
}
