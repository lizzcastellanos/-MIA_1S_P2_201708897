package Admindiscos

import (
	"fmt"
	"os"
	"strings"
)

func Rmdisk(entrada []string) string{
	//Quitar espacios en blanco
	tmp := strings.TrimRight(entrada[1]," ")
	valores := strings.Split(tmp,"=")
	var path string

	if len(valores)!=2{
		fmt.Println("ERROR RMDISK, valor desconocido de parametros ",valores[1])
		return "ERROR RMDISK, valor desconocido de parametros "+ valores[1]
	}else{		
		path = strings.ReplaceAll(valores[1],"\"","")
	}

	//validar si existe el archivo a eliminar
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("RMDISK Error: El disco ", path, " no existe")
		return "RMDISK Error: El disco "+ path + " no existe"
	}

	//Eliminar disco
	err2 := os.Remove(path)
	if err2 != nil {
		fmt.Println("RMDISK Error: No pudo removerse el disco ")
		return "RMDISK Error: No pudo removerse el disco "
	}
	fmt.Println("Disco ", path, "eliminado correctamente:")

	disco := strings.Split(path,"/")
	return "Disco " + disco[len(disco)-1] + " eliminado "
}

