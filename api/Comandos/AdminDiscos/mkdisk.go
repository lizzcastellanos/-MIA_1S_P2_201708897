package Admindiscos

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// recibe los parametros de mkdisk
func Mkdisk(entrada []string) string {

	var size int      //Obligatorio
	var pathE string  //Obligatorio
	fit := "F"        //Puede ser FF, BF, WF, por default es FF
	unit := 1048576   //PUede ser megas(1048576) o kilos (1024), por default es megas
	Valido := true    //Valida los parametros correcto
	InitSize := false //Valida el ingreso del parametro size
	InitPath := false

	/*
		Se recorren todos los parametros
		_ seria el indice, pero se omite.
		El [1:] indica que se inicializa en el primer parametro de mkdisk
	*/
	for _, parametro := range entrada[1:] {
		//Quitar espacios en blanco
		tmp := strings.TrimRight(parametro, " ")

		//Dividir los parametros entre parametro y valor
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR MKDIS, valor desconocido de parametros ", valores[1])
			Valido = false
			return "ERROR MKDIS, valor desconocido de parametros " + valores[1]
		}

		//********************  SIZE *****************
		if strings.ToLower(valores[0]) == "size" {

			InitSize = true
			var err error
			size, err = strconv.Atoi(valores[1]) //se convierte el valor en un entero
			//if err != nil || size <= 0 { //Se manejaria como un solo error
			if err != nil {
				fmt.Println("MKDISK Error: -size debe ser un valor numerico. se leyo ", valores[1])
				Valido = false
				return "MKDISK Error: -size debe ser un valor numerico. se leyo " + valores[1]
			} else if size <= 0 { //se valida que sea mayor a 0 (positivo)
				fmt.Println("MKDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", valores[1])
				Valido = false
				return "MKDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo " + valores[1]
			}

			//********************  Fit *****************
		} else if strings.ToLower(valores[0]) == "fit" {
			if strings.ToLower(valores[0]) == "bf" {
				fit = "B"
			} else if strings.ToLower(valores[0]) == "wf" {
				fit = "W"
			} else if strings.ToLower(valores[0]) != "ff" {
				fmt.Println("EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO: ", valores[1])
			}

			//*************** UNIT ***********************
		} else if strings.ToLower(valores[0]) == "unit" {
			//si la unidad es k
			if strings.ToLower(valores[1]) == "k" {
				//asigno el valor del parametro en su respectiva variable
				unit = 1024
				//si la unidad no es k ni m es error (si fuera m toma el valor con el que se inicializo unit al inicio del metodo)
			} else if strings.ToLower(valores[1]) != "m" {
				fmt.Println("MKDISK Error en -unit. Valores aceptados: k, m. ingreso: ", valores[1])
				Valido = false
				return "MKDISK Error en -unit. Valores aceptados: k, m. ingreso: " + valores[1]
			}

			//******************* PATH *************
		} else if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1], "\"", "")
			InitPath = true

			//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("MKDISK Error: Parametro desconocido: ", valores[0])
			Valido = false
			return "MKDISK Error: Parametro desconocido: " + valores[0] //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if Valido {
		if InitSize {
			if InitPath {
				tam := size * unit
				// Create file
				err := Herramientas.CrearDisco(pathE)
				if err != nil {
					fmt.Println("MKDISK Error: ", err)
					return "MKDISK Error: " + err.Error()
				}
				// Open bin file
				file, err := Herramientas.OpenFile(pathE)
				if err != nil {
					return "MKDISK Error: " + err.Error()
				}

				datos := make([]byte, tam)
				newErr := Herramientas.WriteObject(file, datos, 0)
				if newErr != nil {
					fmt.Println("MKDISK Error: ", newErr)
					return "MKDISK Error: " + newErr.Error()
				}

				//obtener hora para el id
				ahora := time.Now()
				//obtener los segundos y minutos
				//segundos := ahora.Second()
				minutos := ahora.Minute()

				//genera un numero aleario de 1 a 100
				rand.Seed(time.Now().Unix())
				num := rand.Intn(100)

				//concateno los segundos y minutos como una cadena (de 4 digitos)
				cad := fmt.Sprintf("%02d%02d", num, minutos)
				//convierto la cadena a numero en un id temporal
				idTmp, err := strconv.Atoi(cad)
				if err != nil {
					fmt.Println("MKDISK Error: no converti fecha en entero para id")
				}
				//fmt.Println("id guardado actual ", idTmp)
				// Create a new instance of MBR
				var newMBR Structs.MBR
				newMBR.MbrSize = int32(tam)
				newMBR.Id = int32(idTmp)
				copy(newMBR.Fit[:], fit)
				copy(newMBR.FechaC[:], ahora.Format("02/01/2006 15:04"))
				// Write object in bin file
				if err := Herramientas.WriteObject(file, newMBR, 0); err != nil {
					return "ERROR"
				}

				// Close bin file
				defer file.Close()

				fmt.Println("\n Se creo el disco de forma exitosa")

				//imprimir el disco creado para validar que todo este correcto
				var TempMBR Structs.MBR
				if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
					return "ERROR"
				}
				Structs.PrintMBR(TempMBR)

				fmt.Println("\n======End MKDISK======")

				disco := strings.Split(pathE, "/")
				return "Disco " + disco[len(disco)-1] + " creado "

			} else {
				fmt.Println("ERROR: Debe ingresar el parametro Path")
				return "ERROR: Debe ingresar el parametro Path"
			}
		} else {
			fmt.Println("ERROR: Debe ingresar el parametro Size")
			return "ERROR: Debe ingresar el parametro Size"
		}
	}

	return ""

}
