package adminpermisos

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	TI "MIA_1S2025_P1_201708997/ToolsInodos"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func MKfile(entrada []string) string {
	respuesta := "Comando mkfile"
	parametrosDesconocidos := false
	var path string
	var cont string //path del archivo que esta en nuestra maquina y se copiara en el usuario utilizado
	size := 0       //opcional, si no viene toma valor 0
	r := false
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

		if len(valores) == 2 {
			// --------------- PAHT ------------------
			if strings.ToLower(valores[0]) == "path" {
				path = strings.ReplaceAll(valores[1], "\"", "")
				//-------------- SIZE ---------------------
			} else if strings.ToLower(valores[0]) == "size" {
				//convierto a tipo int
				var err error
				size, err = strconv.Atoi(valores[1]) //se convierte el valor en un entero
				if err != nil {
					fmt.Println("MKFILE Error: Size solo acepta valores enteros. Ingreso: ", valores[1])
					return "MKFILE Error: Size solo acepta valores enteros. Ingreso: " + valores[1]
				}

				//valido que sea mayor a 0
				if size < 0 {
					fmt.Println("MKFILE Error: Size solo acepta valores positivos. Ingreso: ", valores[1])
					return "MKFILE Error: Size solo acepta valores positivos. Ingreso: " + valores[1]
				}
				//-------------- CONT ---------------------
			} else if strings.ToLower(valores[0]) == "cont" {
				cont = strings.ReplaceAll(valores[1], "\"", "")
				_, err := os.Stat(cont)
				if os.IsNotExist(err) {
					fmt.Println("MKFILE Error: El archivo cont no existe")
					respuesta += "MKFILE Error: El archivo cont no existe" + "\n"
					return respuesta // Terminar el bucle porque encontramos un nombre único
				}
			} else {
				parametrosDesconocidos = true
			}
		} else if len(valores) == 1 {
			if strings.ToLower(valores[0]) == "r" {
				r = true
			} else {
				parametrosDesconocidos = true
			}
		} else {
			parametrosDesconocidos = true
		}

		if parametrosDesconocidos {
			fmt.Println("MKFILE Error: Parametro desconocido: ", valores[0])
			respuesta += "MKFILE Error: Parametro desconocido: " + valores[0]
			return respuesta //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if path == "" {
		fmt.Println("MKFIEL ERROR NO SE INGRESO PARAMETRO PATH")
		return "MKFIEL ERROR NO SE INGRESO PARAMETRO PATH"
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
		finRuta := len(stepPath) - 1 //es el archivo -> stepPath[finRuta] = archivoNuevo.txt
		idInicial := int32(0)
		idActual := int32(0)
		crear := -1
		//No incluye a finRuta, es decir, se queda en el aterior. EJ: Tamaño=5, finRuta=4. El ultimo que evalua es stepPath[3]
		for i, itemPath := range stepPath[1:finRuta] {
			idActual = TI.BuscarInodo(idInicial, "/"+itemPath, superBloque, Disco)
			//si el actual y el inicial son iguales significa que no existe la carpeta
			if idInicial != idActual {
				idInicial = idActual
			} else {
				crear = i + 1 //porque estoy iniciando desde 1 e i inicia en 0
				break
			}
		}

		//crear carpetas padre si se tiene permiso
		if crear != -1 {
			if r {
				for _, item := range stepPath[crear:finRuta] {
					idInicial = TI.CreaCarpeta(idInicial, item, int64(mbr.Partitions[part].Start), Disco)
					if idInicial == 0 {
						fmt.Println("MKDIR ERROR: No se pudo crear carpeta")
						return "MKFILE ERROR: No se pudo crear carpeta"
					}
				}
			} else {
				fmt.Println("MKDIR ERROR: Carpeta ", stepPath[crear], " no existe. Sin permiso de crear carpetas padre")
				return "MKFILE ERROR: Carpeta " + stepPath[crear] + " no existe. Sin permiso de crear carpetas padre"
			}

		}

		//verificar que no exista el archivo (recordar que BuscarInodo busca de la forma /nombreBuscar)
		idNuevo := TI.BuscarInodo(idInicial, "/"+stepPath[finRuta], superBloque, Disco)
		if idNuevo == idInicial {
			if cont == "" {
				digito := 0
				var content string

				//Crea el contenido del archivo con digitos del 0 al 9
				for i := 0; i < size; i++ {
					if digito == 10 {
						digito = 0
					}
					content += strconv.Itoa(digito)
					digito++
				}
				respuesta = crearArchivo(idInicial, stepPath[finRuta], size, content, int64(mbr.Partitions[part].Start), Disco)
			} else {
				archivoC, err := Herramientas.OpenFile(cont)
				if err != nil {
					return "MKFILE ERROR OPEN FILE " + err.Error() + "\n"
				}

				//lee el contenido del archivo
				content, err := ioutil.ReadFile(cont)
				if err != nil {
					fmt.Println(err)
					return "ERROR MKFILE " + err.Error()
				}
				// Close bin file
				defer archivoC.Close()
				respuesta = crearArchivo(idInicial, stepPath[finRuta], size, string(content), int64(mbr.Partitions[part].Start), Disco)
			}
		} else {
			fmt.Println("El archivo ya existe")
			return "ERROR: El archivo ya existe"
		}
	}
	return respuesta
}

func crearArchivo(idInodo int32, file string, size int, contenido string, initSuperBloque int64, disco *os.File) string {
	//cargar el superBloque actual
	var superB Structs.Superblock
	Herramientas.ReadObject(disco, &superB, initSuperBloque)
	// cargo el inodo de la carpeta que contendra el archivo
	var inodoFile Structs.Inode
	Herramientas.ReadObject(disco, &inodoFile, int64(superB.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	//recorro el inodo de la carpeta para ver donde guardar el archivo (si hay espacio)
	for i := 0; i < 12; i++ {
		idBloque := inodoFile.I_block[i]
		if idBloque != -1 {
			//Existe un folderblock con idBloque que se debe revisar si tiene espacio para el nuevo archivo
			var folderBlock Structs.Folderblock
			Herramientas.ReadObject(disco, &folderBlock, int64(superB.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

			//Recorrer el bloque para ver si hay espacio y si hay crear el archivo
			for j := 2; j < 4; j++ {
				apuntador := folderBlock.B_content[j].B_inodo
				//Hay espacio en el bloque
				if apuntador == -1 {
					//modifico el bloque actual
					copy(folderBlock.B_content[j].B_name[:], file)
					ino := superB.S_first_ino //primer inodo libre
					folderBlock.B_content[j].B_inodo = ino
					//ACTUALIZAR EL FOLDERBLOCK ACTUAL (idBloque) EN EL ARCHIVO
					Herramientas.WriteObject(disco, folderBlock, int64(superB.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

					//creo el nuevo inodo archivo
					var newInodo Structs.Inode
					newInodo.I_uid = Structs.UsuarioActual.IdUsr
					newInodo.I_gid = Structs.UsuarioActual.IdGrp
					newInodo.I_size = int32(size) //Size es el tamaño del archivo
					//Agrego las fechas
					ahora := time.Now()
					date := ahora.Format("02/01/2006 15:04")
					copy(newInodo.I_atime[:], date)
					copy(newInodo.I_ctime[:], date)
					copy(newInodo.I_mtime[:], date)
					copy(newInodo.I_type[:], "1") //es archivo
					copy(newInodo.I_perm[:], "664")

					//apuntadores iniciales
					for i := int32(0); i < 15; i++ {
						newInodo.I_block[i] = -1
					}

					//El apuntador a su primer bloque (el primero disponible)
					fileblock := superB.S_first_blo

					//division del contenido en los fileblocks de 64 bytes
					inicio := 0
					fin := 0
					sizeContenido := len(contenido)
					if sizeContenido < 64 {
						fin = len(contenido)
					} else {
						fin = 64
					}

					//crear el/los fileblocks con el contenido del archivo0
					for i := int32(0); i < 12; i++ {
						newInodo.I_block[i] = fileblock
						//Guardar la informacion del bloque
						data := contenido[inicio:fin]
						var newFileBlock Structs.Fileblock
						copy(newFileBlock.B_content[:], []byte(data))
						//escribo el nuevo bloque (fileblock)
						Herramientas.WriteObject(disco, newFileBlock, int64(superB.S_block_start+(fileblock*int32(binary.Size(Structs.Fileblock{})))))

						//modifico el superbloque (solo el bloque usado por iteracion)
						superB.S_free_blocks_count -= 1
						superB.S_first_blo += 1

						//escribir el bitmap de bloques (se usa un bloque por iteracion).
						Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+fileblock))

						//validar si queda data que agregar al archivo para continuar con el ciclo o detenerlo
						calculo := len(contenido[fin:])
						if calculo > 64 {
							inicio = fin
							fin += 64
						} else if calculo > 0 {
							inicio = fin
							fin += calculo
						} else {
							//detener el ciclo de creacion de fileblocks
							break
						}
						//Aumento el fileblock
						fileblock++
					}

					//escribo el nuevo inodo (ino)
					Herramientas.WriteObject(disco, newInodo, int64(superB.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

					//modifico el superbloque por el inodo usado
					superB.S_free_inodes_count -= 1
					superB.S_first_ino += 1
					//Escribir en el archivo los cambios del superBloque
					Herramientas.WriteObject(disco, superB, initSuperBloque)

					//escribir el bitmap de inodos (se uso un inodo).
					Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_inode_start+ino))

					return "Archivo creado exitosamente"
				} //Fin if apuntadores
			} //fin For bloques
		} else {
			//No hay bloques con espacio disponible
			//modificar el inodo actual (por el nuevo apuntador)
			block := superB.S_first_blo //primer bloque libre
			inodoFile.I_block[i] = block
			//Escribir los cambios del inodo inicial
			Herramientas.WriteObject(disco, &inodoFile, int64(superB.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

			//cargo el primer bloque del inodo actual para tomar los datos de actual y padre (son los mismos para el nuevo)
			var folderBlock Structs.Folderblock
			bloque := inodoFile.I_block[0] //cargo el primer folderblock para obtener los datos del actual y su padre
			Herramientas.ReadObject(disco, &folderBlock, int64(superB.S_block_start+(bloque*int32(binary.Size(Structs.Folderblock{})))))

			//creo el primer bloque que va a apuntar al nuevo archivo
			var newFolderBlock1 Structs.Folderblock
			newFolderBlock1.B_content[0].B_inodo = folderBlock.B_content[0].B_inodo //actual
			copy(newFolderBlock1.B_content[0].B_name[:], ".")
			newFolderBlock1.B_content[1].B_inodo = folderBlock.B_content[1].B_inodo //padre
			copy(newFolderBlock1.B_content[1].B_name[:], "..")
			ino := superB.S_first_ino                          //primer inodo libre
			newFolderBlock1.B_content[2].B_inodo = ino         //apuntador al inodo nuevo
			copy(newFolderBlock1.B_content[2].B_name[:], file) //nombre del inodo nuevo
			newFolderBlock1.B_content[3].B_inodo = -1
			//escribo el nuevo bloque (block)
			Herramientas.WriteObject(disco, newFolderBlock1, int64(superB.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

			//escribir el bitmap de bloques
			Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+block))

			//modifico el superbloque porque mas adelante lo necesito con estos cambios
			superB.S_first_blo += 1
			superB.S_free_blocks_count -= 1

			//creo el nuevo inodo archivo
			var newInodo Structs.Inode
			newInodo.I_uid = Structs.UsuarioActual.IdUsr
			newInodo.I_gid = Structs.UsuarioActual.IdGrp
			newInodo.I_size = int32(size) //Size es el tamaño del archivo
			//Agrego las fechas
			ahora := time.Now()
			date := ahora.Format("02/01/2006 15:04")
			copy(newInodo.I_atime[:], date)
			copy(newInodo.I_ctime[:], date)
			copy(newInodo.I_mtime[:], date)
			copy(newInodo.I_type[:], "1") //es archivo
			copy(newInodo.I_mtime[:], "664")

			//apuntadores iniciales
			for i := int32(0); i < 15; i++ {
				newInodo.I_block[i] = -1
			}

			//El apuntador a su primer bloque (el primero disponible)
			fileblock := superB.S_first_blo

			//division del contenido en los fileblocks de 64 bytes
			inicio := 0
			fin := 0
			sizeContenido := len(contenido)
			if sizeContenido < 64 {
				fin = len(contenido)
			} else {
				fin = 64
			}

			//crear el/los fileblocks con el contenido del archivo0
			for i := int32(0); i < 12; i++ {
				newInodo.I_block[i] = fileblock
				//Guardar la informacion del bloque
				data := contenido[inicio:fin]
				var newFileBlock Structs.Fileblock
				copy(newFileBlock.B_content[:], []byte(data))
				//escribo el nuevo bloque (fileblock)
				Herramientas.WriteObject(disco, newFileBlock, int64(superB.S_block_start+(fileblock*int32(binary.Size(Structs.Fileblock{})))))

				//modifico el superbloque (solo el bloque usado por iteracion)
				superB.S_free_blocks_count -= 1
				superB.S_first_blo += 1

				//escribir el bitmap de bloques (se usa un bloque por iteracion).
				Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+fileblock))

				//validar si queda data que agregar al archivo para continuar con el ciclo o detenerlo
				calculo := len(contenido[fin:])
				if calculo > 64 {
					inicio = fin
					fin += 64
				} else if calculo > 0 {
					inicio = fin
					fin += calculo
				} else {
					//detener el ciclo de creacion de fileblocks
					break
				}
				//Aumento el fileblock
				fileblock++
			}

			//escribo el nuevo inodo (ino)
			Herramientas.WriteObject(disco, newInodo, int64(superB.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

			//modifico el superbloque por el inodo usado
			superB.S_free_inodes_count -= 1
			superB.S_first_ino += 1
			//Escribir en el archivo los cambios del superBloque
			Herramientas.WriteObject(disco, superB, initSuperBloque)

			//escribir el bitmap de inodos (se uso un inodo).
			Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_inode_start+ino))

			return "Archivo creado exitosamente"
		}
	}

	return "ERROR MKFILE: OCURRIO UN ERROR INESPERADO AL CREAR EL ARCHIVO"
}
