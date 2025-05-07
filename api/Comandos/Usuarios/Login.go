package Usuarios

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Login(entrada []string) string {
	var respuesta string
	var user string //obligatorio. Nombre
	var pass string //obligatorio
	var id string   //obligatorio. Id de la particion en la que quiero iniciar sesion
	Valido := true
	var pathDico string

	if Structs.UsuarioActual.Status {
		Valido = false
		return "LOGIN ERROR: Ya existe una sesion iniciada, cierre sesion para iniciar otra"
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR LOGIN, valor desconocido de parametros ", valores[1])
			respuesta += "ERROR LOGIN, valor desconocido de parametros " + valores[1] + "\n"
			Valido = false
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//********************  ID *****************
		if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])

			//********************  USER *****************
		} else if strings.ToLower(valores[0]) == "user" {
			user = valores[1]

			//******************** PASS *****************
		} else if strings.ToLower(valores[0]) == "pass" {
			pass = valores[1]

			//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("LOGIN ERROR: Parametro desconocido: ", valores[0])
			Valido = false
			//por si en el camino reconoce algo invalido de una vez se sale
			return "LOGIN ERROR: Parametro desconocido: " + valores[0] + "\n"
		}
	}

	// Se valida que se haya ingresado los diferentes parametros
	if id != "" {
		//BUsca en struck de particiones montadas el id ingresado
		for _, montado := range Structs.Montadas {
			if montado.Id == id {
				pathDico = montado.PathM
			}
		}
		if pathDico == "" {
			Valido = false
			return "ERROR LOGIN: ID NO ENCONTRADO" + "\n"
		}
	} else {
		fmt.Println("LOGIN ERROR: FALTO EL PARAMETRO ID ")
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO ID " + "\n"
	}

	if pass == "" {
		fmt.Println("LOGIN ERROR: FALTO EL PARAMETRO PASS ")
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO PASS " + "\n"
	}

	if user == "" {
		fmt.Println("LOGIN ERROR: FALTO EL PARAMETRO USER ")
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO USER " + "\n"
	}

	if Valido {
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE " + err.Error() + "\n"
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE " + err.Error() + "\n"
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("REP Error. Particion sin formato")
			return "REP Error. Particion sin formato" + "\n"
		}

		var inodo Structs.Inode
		//Le agrego una structura de inodo para ver el user.txt que esta en el primer inodo del sb
		Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

		//leer los datos del user.txt
		var contenido string
		var fileBLock Structs.Fileblock
		for _, item := range inodo.I_block {
			if item != -1 {
				Herramientas.ReadObject(file, &fileBLock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
				contenido += string(fileBLock.B_content[:])
			}
		}

		//separa el contenido de user.txt
		linea := strings.Split(contenido, "\n")
		//UID, TIPO, GRUPO, USUARIO, CONTRASEÑA

		logeado := false
		for _, reglon := range linea {
			Usuario := strings.Split(reglon, ",")

			//Verifica que cuente con todos los datos
			if len(Usuario) == 5 {
				if Usuario[0] != "0" {
					if Usuario[3] == user {
						if Usuario[4] == pass {
							Structs.UsuarioActual.IdPart = id
							Structs.UsuarioActual.Nombre = user
							Structs.UsuarioActual.Status = true
							Structs.UsuarioActual.PathD = pathDico
							Add_idUsr(Usuario[0])
							logeado = true
							Search_IdGrp(linea, Usuario[2])
						} else {
							fmt.Println("ERROR CONTRASEÑA INCORRECTA")
							return "ERROR LOGIN: LA CONTRASEÑA ES INCORRECTA"
						}
						break
					}
				}
			}
		}

		if logeado {
			respuesta += "EL ususario '" + user + "' ha iniciado sesion exitosamente! \n"
			fmt.Println("IdPart: ", Structs.UsuarioActual.IdPart, " IdGr: ", Structs.UsuarioActual.IdGrp, " User: ", Structs.UsuarioActual.IdUsr, " Nombre: ", Structs.UsuarioActual.Nombre, " Status: ", Structs.UsuarioActual.Status)
		} else {
			fmt.Println("ERROR AL INTENTAR INGRESAR, NO SE ENCONTRO EL USUARIO \nPOR FAVOR INGRESE LOS DATOS CORRECTOS")
			respuesta += "ERROR AL INTENTAR INGRESAR, NO SE ENCONTRO EL USUARIO \n"
			respuesta += "POR FAVOR INGRESE LOS DATOS CORRECTOS \n"
		}
	}

	return respuesta
}

func Add_idUsr(id string) string {
	idU, errId := strconv.Atoi(id)
	if errId != nil {
		fmt.Println("LOGIN ERROR: Error desconcocido con el idUsr")
		return "LOGIN ERROR: Error desconcocido con el idUsr"
	}
	Structs.UsuarioActual.IdUsr = int32(idU)
	return ""
}

func Search_IdGrp(lineaID []string, grupo string) {
	for _, registro := range lineaID[:len(lineaID)-1] {
		datos := strings.Split(registro, ",")
		if len(datos) == 3 {
			if datos[2] == grupo {
				//convertir a numero
				id, errId := strconv.Atoi(datos[0])
				if errId != nil {
					fmt.Println("LOGIN ERROR: Error desconcocido con el idGrp")
					return
				}
				Structs.UsuarioActual.IdGrp = int32(id)
				return
			}
		}
	}
}
