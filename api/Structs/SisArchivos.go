package Structs

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"fmt"
	"os"
	"strings"
)

// NOTA: Recordar que los atributos de los struct deben iniciar con mayuscula
// SUPERBLOQUES
type Superblock struct {
	S_filesystem_type   int32    //numero que identifica el sistema de archivos usado //0->no formateada; 2->ext2; 3->ext3
	S_inodes_count      int32    //numero total de inods creados
	S_blocks_count      int32    //numero total de bloques creados
	S_free_blocks_count int32    //numero de bloques libres
	S_free_inodes_count int32    //numero de inodos libres
	S_mtime             [16]byte //ultima fecha en que el sistema fue montado "02/01/2006 15:04"
	S_umtime            [16]byte //ultima fecha en que el sistema fue desmontado "02/01/2006 15:04"
	S_mnt_count         int32    //numero de veces que se ha montado el sistema
	S_magic             int32    //valor que identifica el sistema de archivos (Sera 0xEF53)
	S_inode_size        int32    //tamaño de la etructura inodo
	S_block_size        int32    //tamaño de la estructura bloque
	S_first_ino         int32    //primer inodo libre
	S_first_blo         int32    //primer bloque libre
	S_bm_inode_start    int32    //inicio del bitmap de inodos
	S_bm_block_start    int32    //inicio del bitmap de bloques
	S_inode_start       int32    //inicio de la tabla de inodos
	S_block_start       int32    //inicio de la tabla de bloques
}

// INODO
type Inode struct {
	I_uid   int32     //ID del usuario propietario del archivo o carpeta
	I_gid   int32     //ID del grupo al que pertenece el archivo o carpeta
	I_size  int32     //tamaño del archivo en bytes
	I_atime [16]byte  //ultima fecha que se leyó el inodo sin modificarlo "02/01/2006 15:04"
	I_ctime [16]byte  //fecha en que se creo el inodo "02/01/2006 15:04"
	I_mtime [16]byte  //ultima fecha en la que se modifica el inodo "02/01/2006 15:04"
	I_block [15]int32 //-1 si no estan usados. los valores del arreglo son: primeros 12 -> bloques directo;: 13 -> bloque simple indirecto; 14->bloque doble indirecto; 15 -> bloque triple indirecto
	I_type  [1]byte   //1 -> archivo; 0 -> carpeta
	I_perm  [3]byte   //permisos del usuario o carpeta
}

// BLOQUE DE CARPETAS
// tamaño en bytes 4(estructuras content)*12(B_name)*4(B_inodo) = 64
type Folderblock struct {
	B_content [4]Content //contenido de la carpeta
}

type Content struct {
	B_name  [12]byte //nombre de carpeta/archivo
	B_inodo int32    //apuntador a un inodo asociado al archivo/carpeta
}

// BLOQUE DE ARCHIVOS
type Fileblock struct {
	B_content [64]byte //contenido del archivo
}

// BLOQUE DE APUNTADORES INDIRECTOS
type Pointerblock struct {
	B_pointers [16]int32 //apuntadores a bloques (archivo/carpeta)
}

// Metodo que anula bytes nulos para B_name
func GetB_name(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)

	if posicionNulo != -1 {
		if posicionNulo != 0 {
			//tiene bytes nulos
			nombre = nombre[:posicionNulo]
		} else {
			//el  nombre esta vacio
			nombre = "-"
		}

	}
	return nombre //-1 el nombre no tiene bytes nulos
}

// Metodo que anula bytes nulos para B_content PARA EL REPORTE
func GetB_content(nombre string) string {
	// Reemplazar todos los saltos de línea con un guion (-)
	nombre = strings.ReplaceAll(nombre, "\n", "<br/>")
	posicionNulo := strings.IndexByte(nombre, 0)

	if posicionNulo != -1 {
		if posicionNulo != 0 {
			//tiene bytes nulos
			nombre = nombre[:posicionNulo]
		} else {
			//el  nombre esta vacio
			nombre = "-"
		}

	}
	//regreso los saltos de linea ya sin bytes nulos
	//nombre = strings.ReplaceAll(nombre, "-", "\n")
	return nombre //-1 el nombre no tiene bytes nulos
}

// para leer byte por byte los bitmaps (reportes)
type Bite struct {
	Val [1]byte
}

// =========================== REPORTE SuperBLoque =============================================
func RepSB(particion Partition, disco *os.File) string {
	cad := ""
	//cargar el superbloque
	var SuperBloque Superblock
	//para las logicas el superbloque se escribe/lee en particion.start+binary.size(structs.EBR)
	err := Herramientas.ReadObject(disco, &SuperBloque, int64(particion.Start))
	if err != nil {
		fmt.Println("REP Error. Particion sin formato")
		return cad
	}
	//LLenar los campos del reporte
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_filesystem_type </td> \n  <td bgcolor='Azure'> EXT%d </td> \n </tr> \n", SuperBloque.S_filesystem_type)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_inodes_count </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_inodes_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_blocks_count </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_blocks_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_free_inodes_count </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_free_inodes_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_free_blocks_count </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_free_blocks_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_mtime </td> \n  <td bgcolor='#7FC97F'> %s </td> \n </tr> \n", string(SuperBloque.S_mtime[:]))
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_umtime </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(SuperBloque.S_mtime[:]))
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_mnt_count </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_mnt_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_magic </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_magic)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_inode_size </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_inode_size)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_block_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_block_size)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_first_ino </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_first_ino)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_first_blo </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_first_blo)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_bm_inode_start </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_bm_inode_start)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_bm_block_start </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_bm_block_start)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_inode_start </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_inode_start)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_block_start </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_block_start)

	return cad
}
