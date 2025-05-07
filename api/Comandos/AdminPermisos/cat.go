package adminpermisos

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	toolsinodos "MIA_1S2025_P1_201708997/ToolsInodos"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Cat(entrada []string) string {
	respuesta := ""
	var filen []string

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR CAT: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR CAT, valor desconocido de parametros ", valores[1])
			respuesta += "ERROR CAT, valor desconocido de parametros " + valores[1] + "\n"
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}
		fileN := valores[0][:4] //toma los primeros 4 caracteres de valores[0]

		//******************** File *****************
		if strings.ToLower(fileN) == "file" {
			numero := strings.Split(strings.ToLower(valores[0]), "file")
			_, errId := strconv.Atoi(numero[1])
			if errId != nil {
				fmt.Println("CAT ERROR: No se pudo obtener un numero de fichero")
				return "CAT ERROR: No se pudo obtener un numero de fichero"
			}
			//eliminar comillas
			tmp1 := strings.ReplaceAll(valores[1], "\"", "")
			filen = append(filen, tmp1)
			//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("CAT ERROR: Parametro desconocido: ", valores[0])
			//por si en el camino reconoce algo invalido de una vez se sale
			return "CAT ERROR: Parametro desconocido: " + valores[0] + "\n"
		}
	}

	//Abrimos el disco
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "CAR ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	// Read object from bin file
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "CAR ERROR READ FILE " + err.Error() + "\n"
	}

	// Close bin file
	defer Disco.Close()

	//Encontrar la particion correcta
	buscar := false
	part := -1 //particion a utilizar y modificar
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			buscar = true
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	if buscar {
		var contenido string
		var fileBlock Structs.Fileblock
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("CAT ERROR. Particion sin formato")
			return "CAT ERROR. Particion sin formato" + "\n"
		}

		//buscar el contenido de cada archivo especificado
		for _, item := range filen {
			//buscar el inodo que contiene el archivo buscado
			idInodo := toolsinodos.BuscarInodo(0, item, superBloque, Disco)
			var inodo Structs.Inode

			//idInodo: solo puede existir archivos desde el inodo 1 en adelante (-1 no existe, 0 es raiz)
			if idInodo > 0 {
				contenido += "\nContenido del archivo: '" + item + "':\n"
				Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

				//Verifica que el usuario logiado sea root(root tiene todos los permisos) o que sea el propietario del archivo
				if inodo.I_uid == UsuarioA.IdUsr || UsuarioA.Nombre == "root" {

					//recorrer los fileblocks del inodo para obtener toda su informacion
					for _, idBlock := range inodo.I_block {
						if idBlock != -1 {
							Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
							tmpConvertir := Herramientas.EliminartIlegibles(string(fileBlock.B_content[:]))
							contenido += tmpConvertir
						}
					}
					contenido += "\n"
				} else {
					contenido += "ERROR CAT: No tiene permisos para visualizar el archivo " + item + "\n"
				}

			} else {
				contenido += "\nCAT ERROR: No se encontro el archivo " + item + "\n"
			}
		}
		respuesta += contenido
		fmt.Println("Contenido ", contenido)
	}

	return respuesta
}
