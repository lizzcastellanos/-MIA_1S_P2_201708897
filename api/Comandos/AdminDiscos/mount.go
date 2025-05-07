package Admindiscos

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Mount(entrada []string) string {
	var respuesta string
	var name string  // Nombre de la particion a montar
	var pathE string // Path del Disco
	Valido := true

	// Validación de parámetros
	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			respuesta = "ERROR MOUNT, valor desconocido de parametros " + valores[1]
			fmt.Println(respuesta)
			return respuesta
		}

		if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1], "\"", "")
			if _, err := os.Stat(pathE); os.IsNotExist(err) {
				respuesta = "ERROR MOUNT: El disco no existe"
				fmt.Println(respuesta)
				return respuesta
			}
		} else if strings.ToLower(valores[0]) == "name" {
			name = strings.TrimSpace(strings.ReplaceAll(valores[1], "\"", ""))
		} else {
			respuesta = "ERROR MOUNT: Parametro desconocido: " + valores[0]
			fmt.Println(respuesta)
			return respuesta
		}
	}

	if !Valido || pathE == "" || name == "" {
		respuesta = "ERROR: Faltan parámetros obligatorios (path y name)"
		fmt.Println(respuesta)
		return respuesta
	}

	// Procesamiento del montaje
	disco, err := Herramientas.OpenFile(pathE)
	if err != nil {
		respuesta = "ERROR NO SE PUEDE LEER EL DISCO " + err.Error()
		fmt.Println(respuesta)
		return respuesta
	}

	var mbr Structs.MBR
	if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
		respuesta = "ERROR Read " + err.Error()
		fmt.Println(respuesta)
		return respuesta
	}
	defer disco.Close()

	montar := true
	reportar := false

	// Buscar la partición por nombre
	for i := 0; i < 4; i++ {
		nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
		if nombre == name {
			montar = false

			// Validar tipo de partición
			if string(mbr.Partitions[i].Type[:]) == "E" {
				respuesta = "ERROR MOUNT. No se puede montar una particion extendida"
				fmt.Println(respuesta)
				return respuesta
			}

			// Validar si ya está montada
			if string(mbr.Partitions[i].Status[:]) == "A" {
				respuesta = "ERROR MOUNT. Esta particion ya fue montada previamente"
				fmt.Println(respuesta)
				return respuesta
			}

			// Generar ID según reglas
			var id string
			var nuevaLetra byte = 'A'
			contador := 1
			discoYaMontado := false

			// Buscar si el disco ya tiene particiones montadas
			for k := 0; k < len(Structs.Pmontaje); k++ {
				if Structs.Pmontaje[k].MPath == pathE {
					discoYaMontado = true
					Structs.Pmontaje[k].Cont++
					contador = int(Structs.Pmontaje[k].Cont)
					nuevaLetra = Structs.Pmontaje[k].Letter
					break
				}
			}

			if !discoYaMontado {
				// Asignar nueva letra para disco no montado
				if len(Structs.Pmontaje) > 0 {
					ultimaLetra := Structs.Pmontaje[len(Structs.Pmontaje)-1].Letter
					if ultimaLetra >= 'E' {
						respuesta = "ERROR: Límite de discos alcanzado (solo se permiten A-E)"
						fmt.Println(respuesta)
						return respuesta
					}
					nuevaLetra = ultimaLetra + 1
				}
				Structs.AddPathM(pathE, nuevaLetra, 1)
			}

			// Formato del ID: 97 + número + letra
			id = "97" + strconv.Itoa(contador) + string(nuevaLetra)
			fmt.Printf("ID generado: %s (Letra: %c, Contador: %d)\n", id, nuevaLetra, contador)

			// Registrar montaje
			Structs.AddMontadas(id, pathE)

			// Actualizar MBR
			copy(mbr.Partitions[i].Status[:], "A")
			copy(mbr.Partitions[i].Id[:], id)
			mbr.Partitions[i].Correlative = int32(contador)

			if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
				respuesta = "Error al escribir MBR: " + err.Error()
				fmt.Println(respuesta)
				return respuesta
			}

			reportar = true
			respuesta = fmt.Sprintf("Partición '%s' montada correctamente. ID: %s", name, id)
			fmt.Println(respuesta)
			break
		}
	}

	if montar {
		respuesta = fmt.Sprintf("ERROR MOUNT. No se encontró la partición '%s'", name)
		fmt.Println(respuesta)
		return respuesta
	}

	if reportar {
		// Reporte de particiones montadas
		reporte := "\n══════════════════════════════════════════════════\n"
		reporte += "            PARTICIONES MONTADAS ACTUALES\n"
		reporte += "══════════════════════════════════════════════════\n"

		// Particiones en el disco actual
		for i := 0; i < 4; i++ {
			if string(mbr.Partitions[i].Status[:]) == "A" {
				reporte += fmt.Sprintf("Partición %d: %s (ID: %s, Tipo: %s, Size: %dMB)\n",
					i,
					Structs.GetName(string(mbr.Partitions[i].Name[:])),
					string(mbr.Partitions[i].Id[:]),
					string(mbr.Partitions[i].Type[:]),
					mbr.Partitions[i].Size)
			}
		}

		// Lista global de montajes
		reporte += "\nLISTA COMPLETA DE MONTADOS:\n"
		for _, m := range Structs.Montadas {
			reporte += fmt.Sprintf("- ID: %s | Disco: %s\n", m.Id, m.PathM)
		}

		reporte += "══════════════════════════════════════════════════\n"
		respuesta += reporte
		fmt.Println(reporte)
	}

	return respuesta
}
