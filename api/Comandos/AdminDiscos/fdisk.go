package Admindiscos

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Fdisk(entrada []string) string {
	var respuesta string
	//size, path, name obligatorios
	//unit, type, fit
	unit := 1024      //Valores B,K,M; por defauld es en k
	tipe := "P"       //Valores P(primaria) E(extendida) L(Logica); por defauld P
	fit := "W"        //Puede ser FF, BF, WF, por default es FF
	var size int      //Obligatorio
	var pathE string  //Obligatorio
	var name string   //Obligatorio
	Valido := true    //Para validar que los parametros cumplen con los requisitos
	InitSize := false //Sirve para saber si se inicializo size (por si no viniera el parametro por ser opcional) false -> no inicializado
	InitPath := false
	var sizeValErr string //Para reportar el error si no se pudo convertir a entero el size

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR FDISK, valor desconocido de parametros ", valores[1])
			respuesta += "ERROR FDISK, valor desconocido de parametros " + valores[1] + "\n"
			Valido = false
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		//********************  SIZE *****************
		if strings.ToLower(valores[0]) == "size" {
			InitSize = true
			var err error
			//size, err = strconv.Atoi(tmp) //se convierte el valor en un entero
			//if err != nil || size <= 0 { //Se manejaria como un solo error
			size, err = strconv.Atoi(valores[1]) //se convierte el valor en un entero
			if err != nil {
				sizeValErr = valores[1] //guarda para el reporte del error si es necesario validar size
			}

			//*************** UNIT ***********************
		} else if strings.ToLower(valores[0]) == "unit" {
			//si la unidad es k
			if strings.ToLower(valores[1]) == "b" {
				unit = 1
				//si la unidad no es k ni m es error (si fuera m toma el valor con el que se inicializo unit al inicio del metodo)
			} else if strings.ToLower(valores[1]) == "m" {
				unit = 1048576 //1024*1024
			} else if strings.ToLower(valores[1]) != "k" {
				fmt.Println("FDISK Error en -unit. Valores aceptados: b, k, m. ingreso: ", valores[1])
				Valido = false
				respuesta += "FDISK Error en -unit. Valores aceptados: b, k, m. ingreso: " + valores[1] + "\n"
				return respuesta
			}

			//******************* PATH *************
		} else if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1], "\"", "")
			InitPath = true

			_, err := os.Stat(pathE)
			if os.IsNotExist(err) {
				fmt.Println("FDISK Error: El disco no existe")
				respuesta += "FDISK Error: El disco no existe" + "\n"
				Valido = false
				return respuesta // Terminar el bucle porque encontramos un nombre único
			}

			//******************* Type *************
		} else if strings.ToLower(valores[0]) == "type" {
			//p esta predeterminado
			if strings.ToLower(valores[1]) == "e" {
				tipe = "E"
			} else if strings.ToLower(valores[1]) == "l" {
				tipe = "L"
			} else if strings.ToLower(valores[1]) != "p" {
				fmt.Println("FDISK Error en -type. Valores aceptados: e, l, p. ingreso: ", valores[1])
				respuesta += "FDISK Error en -type. Valores aceptados: e, l, p. ingreso: " + valores[1] + "\n"
				Valido = false
				return respuesta
			}

			//********************  Fit *****************
		} else if strings.ToLower(valores[0]) == "fit" {
			if strings.ToLower(strings.TrimSpace(valores[1])) == "bf" {
				fit = "B"
			} else if strings.ToLower(valores[1]) == "ff" {
				fit = "F"
			} else if strings.ToLower(valores[1]) != "wf" { //por defecto
				fmt.Println("EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO:", valores[1])
				respuesta += "EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO:" + valores[1] + "\n"
				return respuesta
			}

			//********************  NAME *****************
		} else if strings.ToLower(valores[0]) == "name" {
			// Eliminar comillas
			name = strings.ReplaceAll(valores[1], "\"", "")
			// Eliminar espacios en blanco al final
			name = strings.TrimSpace(name)

			//******************* ERROR EN LOS PARAMETROS *************
		} else {
			fmt.Println("FDISK Error: Parametro desconocido: ", valores[0])
			respuesta += "FDISK Error: Parametro desconocido: " + valores[0] + "\n"
			return respuesta //por si en el camino reconoce algo invalido de una vez se sale
		}

	}

	if InitPath {
		if InitSize {
			if sizeValErr == "" { //Si es un numero (si es numero la variable sizeValErr sera una cadena vacia)
				if size <= 0 { //se valida que sea mayor a 0 (positivo)
					fmt.Println("FDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo ", size)
					respuesta += "FDISK Error: -size debe ser un valor positivo mayor a cero (0). se leyo " + strconv.Itoa(size) + "\n"
					Valido = false
					return respuesta
				}
			} else { //Si sizeValErr es una cadena (por lo que no se pudo dar valor a size)
				fmt.Println("FDISK Error: -size debe ser un valor numerico. se leyo ", sizeValErr)
				respuesta += "FDISK Error: -size debe ser un valor numerico. se leyo " + sizeValErr + "\n"
				Valido = false
				return respuesta
			}
		} else {
			fmt.Println("ERROR: FALTO PARAMETRO SIZE")
			respuesta += "ERROR: FALTO PARAMETRO SIZE" + "\n"
			Valido = false
		}
	} else {
		fmt.Println("ERROR: FALTO PARAMETRO PATH")
		respuesta += "ERROR: FALTO PARAMETRO PATH" + "\n"
		Valido = false
	}

	if Valido {
		if name != "" {
			//Parametros correctos, se puede comenzar a crear las particiones
			disco, err := Herramientas.OpenFile(pathE)
			if err != nil {
				fmt.Println("FDisk Error: No se pudo leer el disco")
				respuesta += "FDisk Error: No se pudo leer el disco" + "\n"
				return respuesta
			}

			//Se crea un mbr para cargar el mbr del disco
			var mbr Structs.MBR
			//Guardo el mbr leido
			if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
				respuesta += "Error Read " + err.Error() + "\n"
				return respuesta
			}

			//Si la particion es tipo extendida validar que no exista alguna extendida
			isPartExtend := false //Indica si se puede usar la particion extendida
			isName := true        //Valida si el nombre no se repite (true no se repite)
			if tipe == "E" {
				for i := 0; i < 4; i++ {
					tipo := string(mbr.Partitions[i].Type[:])

					if tipo != "E" {
						isPartExtend = true
					} else {
						isPartExtend = false
						isName = false //Para que ya no evalue el nombre ni intente hacer nada mas
						fmt.Println("FDISK Error. Ya existe una particion extendida")
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
						respuesta += "FDISK Error. Ya existe una particion extendida \nFDISK Error. No se puede crear la nueva particion con nombre:  " + name + "\n"
						return respuesta
					}
				}
			}

			//verificar si  el nombre existe en las particiones primarias o extendida
			if isName {
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						isName = false
						fmt.Println("FDISK Error. Ya existe la particion : ", name)
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
						respuesta += "FDISK Error. Ya existe la particion : " + name + "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name + "\n"
						return respuesta
					}
				}
			}

			if isName {
				//Buscar en las logicas si ya existe
				var partExtendida Structs.Partition
				//buscar en que particion esta la particion extendida y guardarla en partExtend
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtendida = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtendida = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtendida = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtendida = mbr.Partitions[3]
				}

				if partExtendida.Size != 0 {
					var actual Structs.EBR
					if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
						respuesta += "Error Read " + err.Error() + "\n"
						return respuesta
					}

					//Evaluo la primer ebr
					if Structs.GetName(string(actual.Name[:])) == name {
						isName = false
					} else {
						for actual.Next != -1 {
							//actual = actual.next
							if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
								respuesta += "Error Read " + err.Error() + "\n"
								return respuesta
							}
							if Structs.GetName(string(actual.Name[:])) == name {
								isName = false
								break
							}
						}
					}

					if !isName {
						fmt.Println("FDISK Error. Ya existe la particion : ", name)
						fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
						respuesta += "FDISK Error. Ya existe la particion : " + name
						respuesta += "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name + "\n"
						return respuesta
					}
				}
			}

			//INGRESO DE PARTICIONES PRIMARIAS Y/O EXTENDIDA (SIN LOGICAS)
			sizeNewPart := size * unit //Tamaño de la nueva particion (tamaño * unidades)
			guardar := false           //Indica si se debe guardar la particion, es decir, escribir en el disco
			var newPart Structs.Partition
			if (tipe == "P" || isPartExtend) && isName { //para que  isPartExtend sea true, typee tendra que ser "E"
				sizeMBR := int32(binary.Size(mbr)) //obtener el tamaño del mbr (el que ocupa fisicamente: 165)
				//Para manejar los demas ajustes hacer un if del fit para llamar a la funcion adecuada
				//F = primer ajuste; B = mejor ajuste; else -> peor ajuste

				//INSERTAR PARTICION (Primer ajuste)
				var resTem string
				mbr, newPart, resTem = primerAjuste(mbr, tipe, sizeMBR, int32(sizeNewPart), name, fit) //int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
				respuesta += resTem
				guardar = newPart.Size != 0

				//escribimos el MBR en el archivo. Lo que no se llegue a escribir en el archivo (aqui) se pierde, es decir, los cambios no se guardan
				if guardar {
					//sobreescribir el mbr
					if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
						respuesta += "Error Write " + err.Error() + "\n"
						return respuesta
					}

					//Se agrega el ebr de la particion extendida en el disco
					if isPartExtend {
						var ebr Structs.EBR
						ebr.Start = newPart.Start
						ebr.Next = -1
						if err := Herramientas.WriteObject(disco, ebr, int64(ebr.Start)); err != nil {
							respuesta += "Error Write " + err.Error() + "\n"
							return respuesta
						}
					}
					//para verificar que lo guardo
					var TempMBR2 Structs.MBR
					// Read object from bin file
					if err := Herramientas.ReadObject(disco, &TempMBR2, 0); err != nil {
						respuesta += "Error Read " + err.Error() + "\n"
						return respuesta
					}
					fmt.Println("\nParticion con nombre " + name + " creada exitosamente")
					respuesta += "\nParticion con nombre " + name + " creada exitosamente" + "\n"
					Structs.PrintMBR(TempMBR2)
				} else {
					//Lo podría eliminar pero tendria que modificar en el metodo del ajuste todos los errores para que aparezca el nombre que se intento ingresar como nueva particion
					fmt.Println("FDISK Error. No se puede crear la nueva particion con nombre: ", name)
					respuesta += "FDISK Error. No se puede crear la nueva particion con nombre: " + name
					return respuesta
				}

				// -------------------- INGRESO PARTICIONES LOGICAS----------------
			} else if tipe == "L" && isName {
				var partExtend Structs.Partition
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtend = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtend = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtend = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtend = mbr.Partitions[3]
				} else {
					fmt.Println("FDISK Error. No existe una particion extendida en la cual crear un particion logica")
					respuesta += "FDISK Error. No existe una particion extendida en la cual crear un particion logica" + "\n"
					return respuesta
				}

				//valido que la particion extendida si exista (podría haber entrado al error que no existe extendida)
				if partExtend.Size != 0 {
					//si tuviera los demas ajustes con un if del fit y uso el metodo segun ajuste
					respuesta += primerAjusteLogicas(disco, partExtend, int32(sizeNewPart), name, fit) + "\n" //int32(sizeNewPart) es para castear el int a int32 que es el tipo que tiene el atributo en el struct Partition
					//repLogicas(partExtend, disco)
				}
			}
			return respuesta

		} else {
			respuesta := "ERROR: FALTA PARAMETRO NAME" + "\n"
			fmt.Println("ERROR: FALTA PARAMETRO NAME")
			return respuesta

		}
	}
	return respuesta
}

// ============================================= PRIMER AJUSTE ==================================
func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition, string) {
	var respuesta string
	var newPart Structs.Partition
	var noPart Structs.Partition //para revertir el set info (simula volverla null)

	//PARTICION 1 (libre) - (size = 0 no se ha creado)
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				//caso particion 4 (no existe)
				if mbr.Partitions[3].Size == 0 {
					//859 <= 1024 - 165
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente" + "\n"
					}
				} else {
					//particion 4 existe
					// 600 < 765 - 165 (600 maximo aceptado)
					if sizeNewPart <= mbr.Partitions[3].Start-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						//Si cabe despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno el correlativo para que coincida con el numero de particion en que se guardo
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					}
				}
				//Fin no existe particion 4
			} else {
				// 3 existe
				//entre mbr y 3 -> 300 <= 465 -165
				if sizeNewPart <= mbr.Partitions[2].Start-sizeMBR {
					mbr.Partitions[0] = newPart
				} else {
					//si no cabe entre el mbr y 3 debe ser despues de 3, es decir, en 4
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					} else {
						//4 existe
						//hay espacio entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordenando los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3 //new part traia 4 y quedo en la tercer particion por eso tambien se modifica aqui
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Hay espacio despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//reconfiguro los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					} //fin si hay espacio entre 3 y 4
				} //fin si no cabe antes de 3
			} //fin 3 existe
		} else {
			//2 existe
			//Si la nueva particion se puede guardar antes de 2
			if sizeNewPart <= mbr.Partitions[1].Start-sizeMBR {
				mbr.Partitions[0] = newPart
			} else {
				//Si no cabe entre mbr y 2
				//Validar si existen 3 y 4
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[2].Size == 0 {
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					} else {
						//4 existe (estamos entre 2 y 4)
						//62 < 69-6 (62 maximo aceptado)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							//Si no cabe entre 2 y 4, ver si cabe despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							if sizeNewPart <= mbr.MbrSize-newPart.Start { //1 <= 100-99
								mbr.Partitions[2] = mbr.Partitions[3]
								mbr.Partitions[3] = newPart
								//reordeno correlativos
								mbr.Partitions[2].Correlative = 3
							} else {
								newPart = noPart
								fmt.Println("FDISK Error. Espacio insuficiente")
								respuesta += "FDISK Error. Espacio insuficiente" + "\n"
							}
						} //Fin si cabe antes o despues de 4
					} //fin de 4 existe o no existe
				} else {
					//3 existe
					//entre 2 y 3
					if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
						mbr.Partitions[0] = mbr.Partitions[1]
						mbr.Partitions[1] = newPart
						//Reordeno correlativos
						mbr.Partitions[0].Correlative = 1
						mbr.Partitions[1].Correlative = 2
					} else if mbr.Partitions[3].Size == 0 {
						//entre 3 y el final
						//cambiamos el inicio de la nueva particion porque 3 existe y no cabe antes de 3
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					} else {
						//si 4 existe
						//hay espacio entre 3 y 4
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//entre 4 y el final
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Reordeno correlativos
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					} //Fin si 4 existe o no (3 activa)
				} //Fin 3 existe o no existe
			} //Fin entre 2 y final (antes de 2 o depues de 2)
		} //Fin 2 existe o no existe
		//Fin de 1 no existe

		//PARTICION 2 (no existe)
	} else if mbr.Partitions[1].Size == 0 {
		//Si hay espacio entre el mbr y particion 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start { //particion 1 ya existe (debe existir para entrar a este bloque)
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno correlativo
			mbr.Partitions[1].Correlative = 2
		} else {
			//Si no hay espacio antes de particion 1
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2) //el nuevo inicio es donde termina 1
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[1] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente" + "\n"
					}
				} else {
					//4 existe
					//entre 1 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[1] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente" + "\n"
					}
				} //Fin 4 existe o no existe
			} else {
				//3 Activa
				//entre 1 y 3
				if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					//despues de 3
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
							//corrijo el correlativo
							mbr.Partitions[3].Correlative = 4
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					} else {
						//4 existe
						//entre 3 y 4
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							//Corrijo el correlativo
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							//Despues de 4
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							//Corrijo los correlativos
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							fmt.Println("FDISK Error. Espacio insuficiente")
							respuesta += "FDISK Error. Espacio insuficiente" + "\n"
						}
					} //fin 4 existe o no existe
				} //Fin para entre 1 y 3, y despues de 3
			} //Fin 3 existe o no existe
		} //Fin antes o despues de particion 1
		//Fin particion 2 no existe

		//PARTICION 3
	} else if mbr.Partitions[2].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativo
				mbr.Partitions[2].Correlative = 3
			} else {
				//despues de 2
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[2] = newPart
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente" + "\n"
					}
				} else {
					//4 existe
					//entre 2 y 4
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[2] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						//despues de 4
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						//Reordeno correlativo
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						fmt.Println("FDISK Error. Espacio insuficiente")
						respuesta += "FDISK Error. Espacio insuficiente" + "\n"
					}
				} //Fin de 4 existe o no existe
			} //Fin espacio entre 1 y 2 o despues de 2
		} //Fin espacio antes de 1
		//Fin particion 3

		//PARTICION 4
	} else if mbr.Partitions[3].Size == 0 {
		//antes de 1
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[3] = mbr.Partitions[2]
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			//Reordeno los correlativos
			mbr.Partitions[3].Correlative = 4
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			//si no cabe antes de 1
			//entre 1 y 2
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				//Reordeno correlativos
				mbr.Partitions[3].Correlative = 4
				mbr.Partitions[2].Correlative = 3
			} else if sizeNewPart <= mbr.Partitions[2].Start-mbr.Partitions[1].GetEnd() {
				//entre 2 y 3
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = newPart
				//Reordeno correlativo
				mbr.Partitions[3].Correlative = 4
			} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
				//despues de 3
				newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
				mbr.Partitions[3] = newPart
			} else {
				newPart = noPart
				fmt.Println("FDISK Error. Espacio insuficiente")
				respuesta += "FDISK Error. Espacio insuficiente" + "\n"
			}
		} //Fin antes y despues de 1
		//Fin particion 4
	} else {
		newPart = noPart
		fmt.Println("FDISK Error. Particiones primarias y/o extendidas ya no disponibles")
		respuesta += "FDISK Error. Particiones primarias y/o extendidas ya no disponibles" + "\n"
	}

	return mbr, newPart, respuesta
}

func primerAjusteLogicas(disco *os.File, partExtend Structs.Partition, sizeNewPart int32, name string, fit string) string {
	var respuesta string
	//Se crea un ebr para cargar el ebr desde el disco y la particion extendida
	save := true //false indica que guardo en el primer ebr, true significa que debe seguir buscando
	var actual Structs.EBR
	sizeEBR := int32(binary.Size(actual)) //obtener el tamaño del ebr (el que ocupa fisicamente: 31)
	//fmt.Println("Tamaño fisico del ebr ", sizeEBR)

	//Guardo el ebr leido
	if err := Herramientas.ReadObject(disco, &actual, int64(partExtend.Start)); err != nil {
		respuesta += "Error Read " + err.Error() + "\n"
		return respuesta
	}

	//NOTA: debe caber la particion con el tamaño establecido MAS su EBR
	//NOTA2: Recordar que a la hora de escribir (usar) la particion se inicia donde termina fisicamente la estructura del ebr
	//ej: si el ebr ocupa 5 bytes y la particion es de 10 bytes. los primeros 5 son del ebr entonces uso de 5-15 para escribir en el archivo el contenido de la particion

	//si el primer ebr esta vacio o no existe
	if actual.Size == 0 {
		if actual.Next == -1 {
			//validar si el tamaño de la nueva particion junto al ebr es menor al tamaño de la particion extendida
			if sizeNewPart+sizeEBR <= partExtend.Size {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, -1)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " + err.Error() + "\n"
					return respuesta
				}
				save = false //ya guardo la nueva particion
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre " + name + " creada correctamente" + "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas")
				respuesta += "FDISK Error. Espacio insuficiente logicas" + "\n"
			}
		} else {
			//Para insertar si se elimino la primera particion (primer EBR)
			//Si actual.Next no es -1 significa que hay otra particion despues de la actual y actual.next tiene el inicio de esa particion
			disponible := actual.Next - partExtend.Start //del inicio hasta donde inicia la siguiente
			if sizeNewPart+sizeEBR <= disponible {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, actual.Next)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " + err.Error() + "\n"
					return respuesta
				}
				save = false //ya guardo la nueva particion
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre " + name + " creada correctamente" + "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 2")
				respuesta += "FDISK Error. Espacio insuficiente logicas" + "\n"
			}
		}
		//Si esta despues del primer ebr
	}

	if save {
		//siguiente = actual.next //el valor del siguiente es el inicio de la siguiente particion
		for actual.Next != -1 {
			//si el ebr y la particion caben
			if sizeNewPart+sizeEBR <= actual.Next-actual.GetEnd() {
				break
			}
			//paso al siguiente ebr (simula un actual = actual.next)
			if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
				respuesta += "Error Read " + err.Error() + "\n"
				return respuesta
			}

		}

		//Despues de la ultima particion
		if actual.Next == -1 {
			//ya no es el tamaño porque ya hay espacio ocupado por lo que tomo donde termina la extendida y se resta donde termina la ultima
			if sizeNewPart+sizeEBR <= (partExtend.GetEnd() - actual.GetEnd()) {
				//guardar cambios en el ebr actual (cambio el Next)
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " + err.Error() + "\n"
					return respuesta
				}

				//crea y guarda la nueva particion logica
				newStart := actual.GetEnd()                          //la nueva ebr inicia donde termina la ultima ebr
				actual.SetInfo(fit, newStart, sizeNewPart, name, -1) //cambia actual con los nuevos valores
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " + err.Error() + "\n"
					return respuesta
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre " + name + " creada correctamente" + "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 3")
				respuesta += "FDISK Error. Espacio insuficiente logicas" + "\n"
			}
		} else {
			//Entre dos particiones
			if sizeNewPart+sizeEBR <= (actual.Next - actual.GetEnd()) {
				siguiente := actual.Next //guardo el siguiente de la actual para ponerlo en el siguiente de la nueva particion
				//guardar cambio de siguiente en la actual
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " + err.Error() + "\n"
					return respuesta
				}

				//agrego la nueva particion apuntando a la siguiente de la actual
				newStart := actual.GetEnd()                                 //la nueva ebr inicia donde termina la ultima ebr
				actual.SetInfo(fit, newStart, sizeNewPart, name, siguiente) //cambia actual con los nuevos valores
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					respuesta += "Error Write " + err.Error() + "\n"
					return respuesta
				}
				fmt.Println("Particion con nombre ", name, " creada correctamente")
				respuesta += "Particion con nombre " + name + " creada correctamente" + "\n"
			} else {
				fmt.Println("FDISK Error. Espacio insuficiente logicas 4")
				respuesta += "FDISK Error. Espacio insuficiente logicas " + "\n"
			}
		}
	}
	return respuesta
}
