package Admindiscos

import (
	"MIA_1S2025_P1_201708997/Structs"
	"fmt"
)

func Mounted() string {
	var respuesta string

	if len(Structs.Montadas) == 0 {
		respuesta = "No hay particiones montadas actualmente"
		fmt.Println(respuesta)
		return respuesta
	}

	respuesta = "\n══════════════════════════════════════════════════\n"
	respuesta += "            PARTICIONES MONTADAS\n"
	respuesta += "══════════════════════════════════════════════════\n"

	// Mostrar información detallada de cada partición montada
	for _, montada := range Structs.Montadas {
		respuesta += fmt.Sprintf("ID: %s\n", montada.Id)
		respuesta += fmt.Sprintf("Disco: %s\n", montada.PathM)
		respuesta += "----------------------------------------------\n"
	}

	respuesta += fmt.Sprintf("\nTotal: %d particiones montadas\n", len(Structs.Montadas))
	respuesta += "══════════════════════════════════════════════════\n"

	fmt.Println(respuesta)
	return respuesta
}
