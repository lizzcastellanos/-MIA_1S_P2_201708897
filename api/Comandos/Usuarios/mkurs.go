package Usuarios

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Mkusr(entrada []string) string {
	var respuesta string
	var user string //obligatorio (max 10 caracteres)
	var pass string //obligatorio (max 10 caracteres)
	var grp string  //obligatorio (max 10 caracteres)
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
			Valido = false
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//******************** GRUP *****************
		if strings.ToLower(valores[0]) == "grp" {
			grp = strings.ReplaceAll(valores[1], "\"", "")
			//validar maximo 10 caracteres
			if len(grp) > 10 {
				Valido = false
				fmt.Println("MKGRP ERROR: grp debe tener maximo 10 caracteres")
				return "ERROR MKGRP: grp debe tener maximo 10 caracteres"
			}
			//********************  USER *****************
		} else if strings.ToLower(valores[0]) == "user" {
			user = strings.ReplaceAll(valores[1], "\"", "")
			tmp1 := strings.TrimRight(user, " ")
			user = tmp1
			//validar maximo 10 caracteres
			if len(user) > 10 {
				Valido = false
				fmt.Println("MKGRP ERROR: user debe tener maximo 10 caracteres")
				return "ERROR MKGRP: user debe tener maximo 10 caracteres"
			}
			//******************** PASS *****************
		} else if strings.ToLower(valores[0]) == "pass" {
			pass = valores[1]
			//validar maximo 10 caracteres
			if len(pass) > 10 {
				Valido = false
				fmt.Println("MKGRP ERROR: pass debe tener maximo 10 caracteres")
				return "ERROR MKGRP: pass debe tener maximo 10 caracteres"
			}
			//******************* ERROR EN LOS PARAMETROS *************
		} else {
			Valido = false
			fmt.Println("MKUSR ERROR: Parametro desconocido: ", valores[0])
			//por si en el camino reconoce algo invalido de una vez se sale
			return "MKUSR ERROR: Parametro desconocido: " + valores[0] + "\n"
		}
	}

	// ------------ COMPROBACIONES DE PARAMETROS OBLIGATORIOS---------------
	if pass == "" {
		Valido = false
		fmt.Println("MKUSR ERROR: FALTO EL PARAMETRO PASS ")
		return "MKUSR ERROR: FALTO EL PARAMETRO PASS " + "\n"
	}

	if user == "" {
		Valido = false
		fmt.Println("MKUSR ERROR: FALTO EL PARAMETRO USER ")
		return "MKUSR ERROR: FALTO EL PARAMETRO USER " + "\n"
	}

	if grp == "" {
		Valido = false
		fmt.Println("MKUSR ERROR: FALTO EL PARAMETRO GRP ")
		return "MKUSR ERROR: FALTO EL PARAMETRO GRP " + "\n"
	}

	if Valido {
		file, err := Herramientas.OpenFile(UsuarioA.PathD)
		if err != nil {
			return "ERROR MKUSR ERROR SB OPEN FILE " + err.Error() + "\n"
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR MKUSR ERROR SB READ FILE " + err.Error() + "\n"
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		AddNewUser := false
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == UsuarioA.IdPart {
				part = i
				AddNewUser = true
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		if AddNewUser {
			var superBloque Structs.Superblock
			errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
			if errREAD != nil {
				fmt.Println("MKUSR Error. Particion sin formato")
				return "MKUSR Error. Particion sin formato" + "\n"
			}

			var inodo Structs.Inode
			//Le agrego una structura de inodo para ver el user.txt que esta en el primer inodo del sb
			Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

			//leer los datos del user.txt
			var contenido string
			var fileBlock Structs.Fileblock
			var idFb int32 //id/numero de ultimo fileblock para trabajar sobre ese
			for _, item := range inodo.I_block {
				if item != -1 {
					Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
					contenido += string(fileBlock.B_content[:])
					idFb = item
				}
			}

			lineaID := strings.Split(contenido, "\n")

			ExGrupo := false
			for _, registro := range lineaID[:len(lineaID)-1] {
				datos := strings.Split(registro, ",")
				//verificamos que el grupo exista
				if len(datos) == 3 {
					if datos[2] == grp {
						ExGrupo = true
					}

				}
				//Verificamos que el usuario no exista
				if len(datos) == 5 {
					if datos[3] == user {
						fmt.Println("MKUSR ERROR: ESTE USUARIO YA EXISTE")
						return "MKUSR ERROR: ESTE USUARIO YA EXISTE"
					}
				}
			}

			if !ExGrupo {
				fmt.Println("NO EXISTE EL GRUPO EN MKURS")
				return "MKURS ERROR, NO EXISTE EL GRUPO, POR FAVOR INGRESE UN GRUPO QUE SI EXISTA"
			}

			//Buscar el ultimo ID activo desde el ultimo hasta el primero (ignorando los eliminado (0))
			//desde -2 porque siempre se crea un salto de linea al final generando una linea vacia al final del arreglo
			id := -1        //para guardar el nuevo ID
			var errId error //para la conversion a numero del ID
			for i := len(lineaID) - 2; i >= 0; i-- {
				registro := strings.Split(lineaID[i], ",")
				//valido que sea un usuario
				if registro[1] == "U" {
					//valido que el id sea distinto a 0 (eliminado)
					if registro[0] != "0" {
						//convierto el id en numero para sumarle 1 y crear el nuevo id
						id, errId = strconv.Atoi(registro[0])
						if errId != nil {
							fmt.Println("MKUSR ERROR: NO SE PUDO OBTENER UN NUEVO ID PARA EL NUEVO GRUPO")
							return "MKUSR ERROR: NO SE PUDO OBTENER UN NUEVO ID PARA EL NUEVO GRUPO"
						}
						id++
						break
					}
				}
			}

			//valido que se haya encontrado un nuevo id
			if id != -1 {
				contenidoActual := string(fileBlock.B_content[:])
				posicionNulo := strings.IndexByte(contenidoActual, 0)
				data := fmt.Sprintf("%d,U,%s,%s,%s\n", id, grp, user, pass)

				//Aseguro que haya al menos un byte libre
				if posicionNulo != -1 {
					libre := 64 - (posicionNulo + len(data))
					if libre > 0 {
						copy(fileBlock.B_content[posicionNulo:], []byte(data))
						//Escribir el fileblock con espacio libre
						Herramientas.WriteObject(file, fileBlock, int64(superBloque.S_block_start+(idFb*int32(binary.Size(Structs.Fileblock{})))))
					} else {
						//Si es 0 (quedó exacta), entra aqui y crea un bloque vacío que podrá usarse para el proximo registro
						data1 := data[:len(data)+libre]
						//Ingreso lo que cabe en el bloque actual
						copy(fileBlock.B_content[posicionNulo:], []byte(data1))
						Herramientas.WriteObject(file, fileBlock, int64(superBloque.S_block_start+(idFb*int32(binary.Size(Structs.Fileblock{})))))

						//Creo otro fileblock para el resto de la informacion
						guardoInfo := true
						for i, item := range inodo.I_block {
							//i es el indice en el arreglo inodo.Iblock
							if item == -1 {
								guardoInfo = false
								//agrego el apuntador del bloque al inodo
								inodo.I_block[i] = superBloque.S_first_blo
								//actualizo el superbloque
								superBloque.S_free_blocks_count -= 1
								superBloque.S_first_blo += 1
								data2 := data[len(data)+libre:]
								//crear nuevo fileblock
								var newFileBlock Structs.Fileblock
								copy(newFileBlock.B_content[:], []byte(data2))

								//escribir las estructuras para guardar los cambios
								// Escribir el superbloque
								Herramientas.WriteObject(file, superBloque, int64(mbr.Partitions[part].Start))

								//escribir el bitmap de bloques (se uso un bloque). inodo.I_block[i] contiene el numero de bloque que se uso
								Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_block_start+inodo.I_block[i]))

								//escribir inodes (es el inodo 1, porque es donde esta users.txt)
								Herramientas.WriteObject(file, inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

								//Escribir bloques
								Herramientas.WriteObject(file, newFileBlock, int64(superBloque.S_block_start+(inodo.I_block[i]*int32(binary.Size(Structs.Fileblock{})))))
								break
							}
						}

						if guardoInfo {
							fmt.Println("MKUSR ERROR: ESPACIO INSUFICIENTE PARA EL NUEVO USUARIO")
							return "MKUSR ERROR: ESPACIO INSUFICIENTE PARA EL NUEVO USUARIO. "
						}
					}
					fmt.Println("Se ha agregado el usuario '" + user + "' al grupo '" + grp + "' exitosamente. ")
					respuesta += "Se ha agregado el usuario '" + user + "' al grupo '" + grp + "' exitosamente. "
					for k := 0; k < len(lineaID)-1; k++ {
						fmt.Println(lineaID[k])
					}
					return respuesta
				} else {
					respuesta += "ERROR MKUSR NO HAY ESPACIO SUFICIENTE"
					fmt.Println("ERROR MKUSR NO HAY ESPACIO SUFICIENTE")
				}
			}

			return respuesta

		} //FIn Add new Usuario
	}

	return respuesta
}
