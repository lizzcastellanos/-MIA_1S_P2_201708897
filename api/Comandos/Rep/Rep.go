package rep

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	toolsinodos "MIA_1S2025_P1_201708997/ToolsInodos"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Rep(entrada []string) string {
	var respuesta string
	var name string     //obligatorio Nombre del reporte a generar
	var path string     //obligatorio Nombre que tendrá el reporte
	var id string       //obligatorio sera el del disco o el de la particion
	var rutaFile string //nombre del archivo o carpeta reporte file/IS
	Valido := true

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR REP, valor desconocido de parametros ", valores[1])
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return "ERROR REP, valor desconocido de parametros " + valores[1]
		}

		if strings.ToLower(valores[0]) == "name" {
			name = strings.ToLower(valores[1])
		} else if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1], "\"", "")
		} else if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		} else if strings.ToLower(valores[0]) == "path_file_ls" {
			rutaFile = strings.ReplaceAll(valores[1], "\"", "")
		} else {
			fmt.Println("REP Error: Parametro desconocido: ", valores[0])
			respuesta += "REP Error: Parametro desconocido: " + valores[0]
			Valido = false
			break //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if Valido {
		if name != "" && id != "" && path != "" {
			switch name {
			case "mbr":
				fmt.Println("reporte mbr")
				respuesta += Rmbr(path, id)
			case "disk":
				fmt.Println("reporte disk")
				respuesta += disk(path, id)
			case "inode":
				fmt.Println("reporte inode")
			case "block":
				fmt.Println("reporte block")
			case "bm_inode":
				fmt.Println("reporte bm_inode")
				respuesta += BM_inode(path, id)
			case "bm_block":
				fmt.Println("reporte bm_block")
				respuesta += BM_Bloque(path, id)
			case "sb":
				fmt.Println("reporte sb")
				respuesta += superBloque(path, id)
			case "file":
				fmt.Println("reporte file")
				respuesta += FILE(path, id, rutaFile)
			case "ls":
				respuesta += LS(path, id, rutaFile)
				fmt.Println("reporte ls")
			case "tree":
				fmt.Println("reporte tree")
				respuesta += GenerateTreeReport(path, id)
			default:
				fmt.Println("REP Error: Reporte ", name, " desconocido")
				respuesta += "REP Error: Reporte " + name + " desconocido"
			}
		} else {
			fmt.Println("REP Error: Faltan parametros")
			respuesta += "REP Error: Faltan parametros"
		}
	}
	return respuesta
}

// =============================== MBR ===============================
func Rmbr(path string, id string) string {
	var Respuesta string
	var pathDico string
	Valido := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido {
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp = strings.Split(pathDico, "/")
		NOmbreDis := strings.Split(tmp[len(tmp)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP MBR Open " + err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			Respuesta += "ERROR REP MBR Read " + err.Error()
		}

		// Close bin file
		defer file.Close()

		//Crea reporte
		cad := `digraph {
			node [shape=none fontname="Arial"]
			TablaReportNodo [label=<
			<table border="0" cellborder="1" cellspacing="0" cellpadding="4">
				<tr>
					<td bgcolor="#2c3e50" colspan="2"><font color="white"><b>REPORTE MBR</b></font></td>
				</tr>
				<tr>
					<td bgcolor="#ecf0f1"><b>mbr_tamano</b></td>
					<td bgcolor="#bdc3c7">` + fmt.Sprintf("%d", mbr.MbrSize) + `</td>
				</tr>
				<tr>
					<td bgcolor="#ecf0f1"><b>mbr_fecha_creacion</b></td>
					<td bgcolor="#bdc3c7">` + string(mbr.FechaC[:]) + `</td>
				</tr>
				<tr>
					<td bgcolor="#ecf0f1"><b>mbr_disk_signature</b></td>
					<td bgcolor="#bdc3c7">` + fmt.Sprintf("%d", mbr.Id) + `</td>
				</tr>`

		cad += Structs.RepGraphviz(mbr, file)
		cad += `</table>>]}`

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
		Respuesta += "Reporte de MBR del disco " + NOmbreDis + " creado con el nombre " + nombre + ".png"
	} else {
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}

	return Respuesta
}

// =============================== DISK ===============================
func disk(path string, id string) string {
	var Respuesta string
	var pathDico string
	Valido := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido {
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp = strings.Split(pathDico, "/")
		NOmbreDis := strings.Split(tmp[len(tmp)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP DISK Open " + err.Error()
			return Respuesta
		}

		var TempMBR Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
			Respuesta += "ERROR REP READ Open " + err.Error()
			return Respuesta
		}

		defer file.Close()

		//inicia contenido del reporte graphviz del disco
		cad := `digraph {
			node [shape=none fontname="Arial"]
			TablaReportNodo [label=<
			<table border="0" cellborder="1" cellspacing="0" cellpadding="4">
				<tr>
					<td bgcolor="#9b59b6" rowspan="3"><font color="white"><b>MBR</b></font></td>`

		cad += Structs.RepDiskGraphviz(TempMBR, file)
		cad += `</table>>]}`

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		fmt.Println("RP ", rutaReporte, " name ", nombre)

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
		Respuesta += "Reporte Disk del disco " + NOmbreDis + " creado con el nombre " + nombre + ".png"
	} else {
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}

	return Respuesta

}

// =============================== SB ===============================
func superBloque(path string, id string) string {
	var respuesta string
	var pathDico string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == "" {
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar {
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE " + err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE " + err.Error()
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		cad := `digraph {
			node [shape=none fontname="Arial"]
			TablaReportNodo [label=<
			<table border="0" cellborder="1" cellspacing="0" cellpadding="4">
				<tr>
					<td bgcolor="#9b59b6" colspan="2"><font color="white"><b>REPORTE SUPERBLOQUE</b></font></td>
				</tr>`

		cad += Structs.RepSB(mbr.Partitions[part], file)
		cad += `</table>>]}`

		//reporte requerido
		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"
		respuesta += "Reporte BM_Bloque " + nombre + " creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	}

	return respuesta
}

// =============================== BM INODE ===============================
func BM_inode(path string, id string) string {
	var respuesta string
	var pathDico string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == "" {
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar {
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE " + err.Error()
		}

		//Obtener mbr
		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE " + err.Error()
		}

		// Close bin file
		defer file.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("REP Error. Particion sin formato")
			return "REP Error. Particion sin formato"
		}

		cad := ""
		inicio := superBloque.S_bm_inode_start
		fin := superBloque.S_bm_block_start
		count := 1 //para contar el numero de caracteres por linea (maximo 20)

		//objeto para leer un byte decodificado
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			//cargo el byte (struct de [1]byte) decodificado como las demas estructuras
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += string("0 ")
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}

		//reporte requerido
		carpeta := filepath.Dir(path) //DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte BM Inode " + nombre + " creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco
	}

	return respuesta
}

// =============================== BM BLOQUE ===============================
func BM_Bloque(path string, id string) string {
	var respuesta string
	var pathDico string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == "" {
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar {
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE " + err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE " + err.Error()
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
			return "REP Error. Particion sin formato"
		}

		cad := ""
		inicio := superBloque.S_bm_block_start
		fin := superBloque.S_inode_start
		count := 1 //para contar el numero de caracteres por linea (maximo 20)

		//objeto para leer un byte decodificado
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			//cargo el byte (struct de [1]byte) decodificado como las demas estructuras
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += string("0 ")
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}

		//reporte requerido
		carpeta := filepath.Dir(path) //DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte BM_Bloque " + nombre + " creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco
	}
	return respuesta
}

// =============================== FILE ===============================
func FILE(path string, id string, rutaFile string) string {
	var respuesta string
	var pathDico string
	var contenido string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == "" {
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar {
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		Disco, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP OPEN FILE " + err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
			return "ERROR REP READ FILE " + err.Error()
		}

		// Close bin file
		defer Disco.Close()

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
		var fileBlock Structs.Fileblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("REP Error. Particion sin formato")
			return "REP Error. Particion sin formato"
		}

		//buscar el inodo que contiene el archivo buscado
		idInodo := toolsinodos.BuscarInodo(0, rutaFile, superBloque, Disco)
		var inodo Structs.Inode

		//idInodo: solo puede existir archivos desde el inodo 1 en adelante (-1 no existe, 0 es raiz)
		if idInodo > 0 {
			contenido += "Contenido del archivo: '" + rutaFile + "'\n"
			Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
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
			fmt.Println("REP ERROR: No se encontro el archivo ", rutaFile)
			return "REP ERROR: No se encontro el archivo " + rutaFile
		}

		//reporte requerido
		carpeta := filepath.Dir(path) //DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, contenido)
		respuesta += "Reporte BM_Bloque " + nombre + " creado \n"
		respuesta += "Pertenece al disco: " + nombreDisco
	}
	return respuesta
}

// =============================== LS ===============================
func LS(path string, id string, rutaFile string) string {
	var respuesta string
	var contenido string
	var pathDico string
	reportar := false

	//BUsca en struck de particiones montadas el id ingresado
	for _, montado := range Structs.Montadas {
		if montado.Id == id {
			pathDico = montado.PathM
			reportar = true
		}
	}

	//Verifica que se encontro el ID y la Path del disco
	if pathDico == "" {
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar {
		Color := "#9b59b6"
		contenido = `digraph {
			node [shape=none fontname="Arial"]
			TablaReportNodo [label=<
			<table border="0" cellborder="1" cellspacing="0" cellpadding="4">
				<tr>
					<td bgcolor="` + Color + `"><font color="white"><b>PERMISOS</b></font></td>
					<td bgcolor="` + Color + `"><font color="white"><b>USUARIO</b></font></td>
					<td bgcolor="` + Color + `"><font color="white"><b>GRUPO</b></font></td>
					<td bgcolor="` + Color + `"><font color="white"><b>SIZE</b></font></td>
					<td bgcolor="` + Color + `"><font color="white"><b>FECHA/HORA</b></font></td>
					<td bgcolor="` + Color + `"><font color="white"><b>NOMBRE</b></font></td>
					<td bgcolor="` + Color + `"><font color="white"><b>TIPO</b></font></td>
				</tr>`
		//Obtenermos el nombre del reporte que vamos a crear
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		Disco, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP OPEN FILE " + err.Error()
		}

		var mbr Structs.MBR
		// Read object from bin file
		if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
			return "ERROR REP READ FILE " + err.Error()
		}

		// Close bin file
		defer Disco.Close()

		//Encontrar la particion correcta
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
			}
		}

		//var fileBlock Structs.Fileblock
		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("CAT ERROR. Particion sin formato")
			return "CAT ERROR. Particion sin formato" + "\n"
		}

		var FstInodo Structs.Inode
		//Le agrego una structura de inodo para ver el user.txt que esta en el primer inodo del sb
		Herramientas.ReadObject(Disco, &FstInodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

		var contUs string
		var FistfileBlock Structs.Fileblock
		for _, item := range FstInodo.I_block {
			if item != -1 {
				Herramientas.ReadObject(Disco, &FistfileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
				contUs += string(FistfileBlock.B_content[:])
			}
		}
		lineaID := strings.Split(contUs, "\n")

		idInodo := toolsinodos.BuscarInodo(0, rutaFile, superBloque, Disco)
		var inodo Structs.Inode

		if idInodo > 0 {
			Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
			var folderBlock Structs.Folderblock
			for _, idBlock := range inodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))
					for k := 2; k < 4; k++ {
						apuntador := folderBlock.B_content[k].B_inodo
						if apuntador != -1 {
							pathActual := Structs.GetB_name(string(folderBlock.B_content[k].B_name[:]))

							contenido += InodoLs(pathActual, lineaID, apuntador, superBloque, Disco)
						}
					}
				}
			}

		} else {
			respuesta = "REP ERROR NO SE ENCONTRO LA PATH INGRESADA"
		}

		contenido += "\n</table> > ]\n}"
		cad := Herramientas.EliminartIlegibles(contenido)

		//reporte requerido
		carpeta := filepath.Dir(path) //DIr es para obtener el directorio
		rutaReporte := carpeta + "/" + nombre + ".dot"
		Herramientas.Reporte(rutaReporte, contenido)
		respuesta += "Reporte BM_Bloque " + nombre + " creado \n"
		respuesta += "Pertenece al disco: " + nombreDisco
		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	}
	return respuesta
}

// Nombre,   contenia users.txt		no. bloque		superbloque					DIsco
func InodoLs(name string, lineaID []string, idInodo int32, superBloque Structs.Superblock, file *os.File) string {
	var contenido string

	//cargar el inodo a reportar
	var inodo Structs.Inode
	Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	//Busco el grupo y el usuario
	usuario := ""
	grupo := ""
	for m := 0; m < len(lineaID); m++ {
		datos := strings.Split(lineaID[m], ",")
		if len(datos) == 5 {
			us := fmt.Sprintf("%d", inodo.I_uid)
			if us == datos[0] {
				usuario = datos[3]
			}
		}
		if len(datos) == 3 {
			gr := fmt.Sprintf("%d", inodo.I_gid)
			if gr == (datos[0]) {
				grupo = datos[2]
			}
		}

	}

	Color := "#33ffc7"
	tipoArchivo := "Archivo"
	var permisos string

	//r w x
	//r -> lectura
	// w escritura
	// x ejecucion
	//son 3 numeros porque son aplicados a: propierarios   grupos  y  otros
	for i := 0; i < 3; i++ {
		if string(inodo.I_perm[i]) == "0" { //ninun permiso
			permisos += "---"
		} else if string(inodo.I_perm[i]) == "1" { // ejecucion
			permisos += "--x"
		} else if string(inodo.I_perm[i]) == "2" { //	escritura
			permisos += "-w-"
		} else if string(inodo.I_perm[i]) == "3" { // 	ecritura ejecucion
			permisos += "-wx"
		} else if string(inodo.I_perm[i]) == "4" { //lectura
			permisos += "r--"
		} else if string(inodo.I_perm[i]) == "5" { //lectura  	ejecucion
			permisos += "r-x"
		} else if string(inodo.I_perm[i]) == "6" { // lectura escritura
			permisos += "rw-"
		} else if string(inodo.I_perm[i]) == "7" { //lectura escritura ejecucion
			permisos += "rwx"
		}
	}

	if string(inodo.I_type[:]) == "0" {
		Color = "#e67e22"
		tipoArchivo = "Carpeta"
		permisos = "rw-rw-r--"
	}
	permisos = "rw-rw-r--"
	contenido += `
    <tr>
        <td bgcolor="` + Color + `">` + permisos + `</td>
        <td bgcolor="` + Color + `">` + usuario + `</td>
        <td bgcolor="` + Color + `">` + grupo + `</td>
        <td bgcolor="` + Color + `">` + fmt.Sprintf("%d", inodo.I_size) + `</td>
        <td bgcolor="` + Color + `">` + string(inodo.I_ctime[:]) + `</td>
        <td bgcolor="` + Color + `">` + name + `</td>
        <td bgcolor="` + Color + `">` + tipoArchivo + `</td>
    </tr>`
	//reportar el inodo
	return contenido
}

func GenerateTreeReport(path string, id string) string {
	var respuesta string
	var pathDisco string
	encontrado := false

	// 1. Buscar el disco montado con el ID proporcionado
	for _, montada := range Structs.Montadas {
		if montada.Id == id {
			pathDisco = montada.PathM
			encontrado = true
			break
		}
	}

	if !encontrado {
		return "ERROR: No se encontró el ID " + id + " en las particiones montadas"
	}

	// 2. Abrir el disco y leer el MBR
	file, err := Herramientas.OpenFile(pathDisco)
	if err != nil {
		return "ERROR al abrir el disco: " + err.Error()
	}
	defer file.Close()

	var mbr Structs.MBR
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return "ERROR al leer el MBR: " + err.Error()
	}

	// 3. Encontrar la partición correspondiente al ID
	var particion Structs.Partition
	encontrada := false
	for i := 0; i < 4; i++ {
		if Structs.GetId(string(mbr.Partitions[i].Id[:])) == id {
			particion = mbr.Partitions[i]
			encontrada = true
			break
		}
	}

	if !encontrada {
		return "ERROR: No se encontró la partición con ID " + id
	}

	// 4. Leer el superbloque de la partición
	var superbloque Structs.Superblock
	if err := Herramientas.ReadObject(file, &superbloque, int64(particion.Start)); err != nil {
		return "ERROR al leer el superbloque: " + err.Error()
	}

	// 5. Generar el contenido DOT para Graphviz
	dotContent := `digraph G {
        rankdir=TB;
        node [shape=record, fontname="Courier New", fontsize=10];
        edge [arrowhead=vee, arrowsize=0.8];
        
        graph [nodesep=0.5, ranksep=0.5];
    `

	// 6. Recorrer los inodos y bloques para construir el árbol
	// Empezamos con el inodo raíz (generalmente inodo 0)
	dotContent += generarNodoInodo(file, &superbloque, 0)

	// 7. Cerrar el gráfico DOT
	dotContent += "\n}"

	// 8. Guardar el archivo DOT y generar la imagen PNG usando Herramientas.RepGraphizMBR
	nombreReporte := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	rutaDot := filepath.Join(filepath.Dir(path), nombreReporte+".dot")

	// Primero guardamos el contenido DOT
	if err := Herramientas.Reporte(rutaDot, dotContent); err != nil {
		return "ERROR al guardar archivo DOT: " + err.Error()
	}

	// Luego generamos la imagen PNG usando tu función existente
	if err := Herramientas.RepGraphizMBR(rutaDot, dotContent, nombreReporte); err != nil {
		return "ERROR al generar imagen PNG: " + err.Error()
	}

	respuesta = "Reporte tree generado exitosamente:\n"
	respuesta += "- Archivo DOT: " + rutaDot + "\n"
	respuesta += "- Imagen PNG: " + filepath.Join(filepath.Dir(path), nombreReporte+".png")

	return respuesta
}

// Las funciones auxiliares generarNodoInodo, generarNodoBloque y getTipoInodo
// se mantienen exactamente igual que en la versión anterior

// Función auxiliar para generar el nodo de un inodo y sus bloques
func generarNodoInodo(file *os.File, superbloque *Structs.Superblock, numInodo int32) string {
	var inodo Structs.Inode
	offset := superbloque.S_inode_start + (numInodo * int32(binary.Size(Structs.Inode{})))

	if err := Herramientas.ReadObject(file, &inodo, int64(offset)); err != nil {
		return ""
	}

	// Generar el nodo del inodo
	dotContent := fmt.Sprintf(`
    inodo%d [label=<
        <table border="0" cellborder="1" cellspacing="0">
            <tr><td colspan="2" bgcolor="#e0e0e0"><b>INODO %d</b></td></tr>
            <tr><td>ID</td><td>%d</td></tr>
            <tr><td>UID</td><td>%d</td></tr>
            <tr><td>Fecha</td><td>%s</td></tr>
            <tr><td>Tipo</td><td>%s</td></tr>
            <tr><td>Size</td><td>%d</td></tr>
        </table>>];`,
		numInodo, numInodo, numInodo, inodo.I_uid,
		string(inodo.I_ctime[:]), getTipoInodo(inodo.I_type), inodo.I_size)

	// Recorrer los apuntadores del inodo
	for i, apuntador := range inodo.I_block {
		if apuntador == -1 {
			continue
		}

		if i < 12 { // Apuntadores directos
			dotContent += generarNodoBloque(file, superbloque, apuntador, numInodo, i)
		} else if i == 12 { // Apuntador indirecto simple
			// Implementar lógica para bloques indirectos si es necesario
		}
		// Puedes agregar más casos para apuntadores indirectos dobles/triples
	}

	return dotContent
}

// Función auxiliar para generar el nodo de un bloque
func generarNodoBloque(file *os.File, superbloque *Structs.Superblock, numBloque int32, numInodoPadre int32, indiceApuntador int) string {
	offset := superbloque.S_block_start + (numBloque * int32(binary.Size(Structs.Fileblock{})))

	var bloque Structs.Fileblock
	if err := Herramientas.ReadObject(file, &bloque, int64(offset)); err != nil {
		return ""
	}

	// Generar el nodo del bloque
	dotContent := fmt.Sprintf(`
    bloque%d [label=<
        <table border="0" cellborder="1" cellspacing="0">
            <tr><td colspan="2" bgcolor="#f0f0f0"><b>BLOQUE %d</b></td></tr>`,
		numBloque, numBloque)

	// Agregar contenido del bloque (dependiendo del tipo)
	contenido := Herramientas.EliminartIlegibles(string(bloque.B_content[:]))
	lineas := strings.Split(contenido, "\n")
	for _, linea := range lineas {
		if linea != "" {
			dotContent += fmt.Sprintf("<tr><td>%s</td></tr>", linea)
		}
	}

	dotContent += `</table>>];`

	// Conectar el inodo padre con este bloque
	dotContent += fmt.Sprintf(`
    inodo%d -> bloque%d [label="Apuntador Directo %d"];`,
		numInodoPadre, numBloque, indiceApuntador+1)

	return dotContent
}

// Función auxiliar para obtener el tipo de inodo como texto
func getTipoInodo(tipo [1]byte) string {
	switch tipo[0] {
	case '0':
		return "Directorio"
	case '1':
		return "Archivo"
	default:
		return "Desconocido"
	}
}
