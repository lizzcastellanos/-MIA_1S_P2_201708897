package Usuarios

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	"encoding/binary"
	"fmt"
	"strings"
)

func Rmusr(entrada []string) string {
	var respuesta string
	var user string
	Valido := true
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		Valido = false
		respuesta += "ERROR MKUSR: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	if UsuarioA.Nombre != "root" {
		Valido = false
		fmt.Println("ERROR FALTA DE PERMISOS, NO ES EL USUARIO ROOT")
		respuesta += "ERROR MKGRO: ESTE USUARIO NO CUENTA CON LOS PERMISOS PARA REALIZAR ESTA ACCION"
		return respuesta
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR MKGRP, valor desconocido de parametros ", valores[1])
			respuesta += "ERROR MKGRP, valor desconocido de parametros " + valores[1] + "\n"
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//******************** USER *****************
		if strings.ToLower(valores[0]) == "user" {
			user = (valores[1])
			//validar maximo 10 caracteres
			if len(user) > 10 {
				fmt.Println("MKGRP ERROR: user debe tener maximo 10 caracteres")
				return "ERROR MKGRP: user debe tener maximo 10 caracteres"
			}
			//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("LOGIN ERROR: Parametro desconocido: ", valores[0])
			//por si en el camino reconoce algo invalido de una vez se sale
			return "LOGIN ERROR: Parametro desconocido: " + valores[0] + "\n"
		}
	}
	// ------------ COMPROBACIONES DE PARAMETROS OBLIGATORIOS---------------

	if user == "" {
		Valido = false
		fmt.Println("MKUSR ERROR: FALTO EL PARAMETRO USER ")
		return "MKUSR ERROR: FALTO EL PARAMETRO USER " + "\n"
	}

	//Abrimos el disco
	file, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "RMUSR ERRORSB OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return "RMUSR ERRORSB READ FILE " + err.Error() + "\n"
	}

	// Close bin file
	defer file.Close()

	//Encontrar la particion correcta
	delete := false
	part := -1 //particion a utilizar y modificar
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			delete = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if Valido && delete {
		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("RMGRP ERROR. Particion sin formato")
			return "RMGRP ERROR. Particion sin formato" + "\n"
		}

		var inodo Structs.Inode
		//Le agrego una structura de inodo para ver el user.txt que esta en el primer inodo del sb
		Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

		//leer los datos del user.txt
		var contenido string
		var fileBlock Structs.Fileblock
		for _, item := range inodo.I_block {
			if item != -1 {
				Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
				contenido += string(fileBlock.B_content[:])
			}
		}

		lineaID := strings.Split(contenido, "\n")
		modificadoUsr := false
		for k := 0; k < len(lineaID); k++ {
			datos := strings.Split(lineaID[k], ",")
			if len(datos) == 5 {
				if datos[3] == user {
					if datos[0] != "0" {
						datos[0] = "0"
						modificadoUsr = true
						lineaID[k] = datos[0] + "," + datos[1] + "," + datos[2] + "," + datos[3] + "," + datos[4]
					} else {
						fmt.Println("ERROR RMUSR, EL USUARIO '" + user + "' ya fue eliminado previamente")
						return "ERROR RMUSR, EL USUARIO '" + user + "' ya fue eliminado previamente"
					}
				}
			}
		}

		if modificadoUsr {
			mod := ""
			for _, reg := range lineaID {
				mod += reg + "\n"
			}

			inicio := 0
			var fin int
			if len(mod) > 64 {
				//si el contenido es mayor a 64 bytes. la primera vez termina en 64
				fin = 64
			} else {
				//termina en el tamaño del contenido. Solo habra un fileblock porque ocupa menos de la capacidad de uno
				fin = len(mod)
			}

			for _, newItem := range inodo.I_block {
				if newItem != -1 {
					//tomo 64 bytes de la cadena o los bytes que queden
					data := mod[inicio:fin]
					//Modifico y guardo el bloque actual
					var newFileBlock Structs.Fileblock
					copy(newFileBlock.B_content[:], []byte(data))
					Herramientas.WriteObject(file, newFileBlock, int64(superBloque.S_block_start+(newItem*int32(binary.Size(Structs.Fileblock{})))))
					//muevo a los siguientes 64 bytes de la cadena (o los que falten)
					inicio = fin
					calculo := len(mod[fin:]) //tamaño restante de la cadena
					//else if
					if calculo > 64 {
						fin += 64
					} else {
						fin += calculo
					}
				}
			}

			fmt.Println("El grupo '" + user + "' fue eliminado con extiso")
			respuesta += "El grupo '" + user + "' fue eliminado con extiso"
			for k := 0; k < len(lineaID)-1; k++ {
				fmt.Println(lineaID[k])
			}
			return respuesta
		} else {
			fmt.Println("ERROR RMUSR NO SE PUDO ENCONTRAR EL USUARIO")
			return "ERROR RMUSR NO SE PUDO ENCONTRAR EL USUARIO"
		}
	} else {
		fmt.Println("ERROR RMUSR ALGO SALIO MAL CON LAS PARTICIONES")
		respuesta += "ERROR RMUSR ALGO SALIO MAL CON LAS PARTICIONES"
	}
	return respuesta
}
