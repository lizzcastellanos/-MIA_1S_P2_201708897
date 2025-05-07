package adminpermisos

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	ToolsInodos "MIA_1S2025_P1_201708997/ToolsInodos"
	"fmt"
	"strings"
)

func Mkdir(entrada []string) string {
	respuesta := "Comando mkdir"
	var path string
	p := false
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		fmt.Println("ERROR MKFILE: SESION NO INICIADA")
		respuesta += "ERROR MKFILE: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if strings.ToLower(valores[0]) == "path" {
			if len(valores) != 2 {
				fmt.Println("ERROR MKDIR, valor desconocido de parametros ", valores[1])
				respuesta += "ERROR MKDIR, valor desconocido de parametros " + valores[1]
				//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
				return respuesta
			}
			path = strings.ReplaceAll(valores[1], "\"", "")
		} else if strings.ToLower(valores[0]) == "r" {
			if len(tmp) != 1 {
				fmt.Println("MKDIR Error: Valor desconocido del parametro ", valores[0])
				return "MKDIR Error: Valor desconocido del parametro " + valores[0]
			}
			p = true

			//ERROR
		} else {
			fmt.Println("MKFILE ERROR: Parametro desconocido: ", valores[0])
			return "MKFILE ERROR: Parametro desconocido: " + valores[0]
		}
	}

	if path == "" {
		fmt.Println("MKDIR ERROR NO SE INGRESO PARAMETRO PATH")
		return "MKDIR ERROR NO SE INGRESO PARAMETRO PATH"
	}

	//Abrimos el disco
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "MKFILE ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "MKFILE ERROR READ FILE " + err.Error() + "\n"
	}

	// Close bin file
	defer Disco.Close()

	//Encontrar la particion correcta
	agregar := false
	part := -1 //particion a utilizar y modificar
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			agregar = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if agregar {
		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("MKFILE ERROR. Particion sin formato")
			return "MKFILE ERROR. Particion sin formato" + "\n"
		}

		//Validar que exista la ruta
		stepPath := strings.Split(path, "/")
		idInicial := int32(0)
		idActual := int32(0)
		crear := -1
		for i, itemPath := range stepPath[1:] {
			idActual = ToolsInodos.BuscarInodo(idInicial, "/"+itemPath, superBloque, Disco)
			if idInicial != idActual {
				idInicial = idActual
			} else {
				crear = i + 1 //porque estoy iniciando desde 1 e i inicia en 0
				break
			}
		}

		//crear carpetas padre si se tiene permiso
		if crear != -1 {
			if crear == len(stepPath)-1 {
				ToolsInodos.CreaCarpeta(idInicial, stepPath[crear], int64(mbr.Partitions[part].Start), Disco)
			} else {
				if p {
					for _, item := range stepPath[crear:] {
						idInicial = ToolsInodos.CreaCarpeta(idInicial, item, int64(mbr.Partitions[part].Start), Disco)
						if idInicial == 0 {
							fmt.Println("MKDIR ERROR: No se pudo crear carpeta")
							return "MKDIR ERROR: No se pudo crear carpeta"
						}
					}
				} else {
					fmt.Println("MKDIR ERROR: Sin permiso de crear carpetas padre")
				}
			}
			return "Carpeta(s) creada"
		} else {
			fmt.Println("MKDIR ERROR: LA CARPETA YA EXISTE")
			return "MKDIR ERROR: LA CARPETA YA EXISTE"
		}
	}
	return respuesta
}
