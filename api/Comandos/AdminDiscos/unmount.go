package Admindiscos

import (
	"MIA_1S2025_P1_201708997/Herramientas"
	"MIA_1S2025_P1_201708997/Structs"
	"fmt"
	"strings"
)

func Unmoun(entrada []string) string {
	var respuesta string
	var id string

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			fmt.Println("ERROR MOUNT, valor desconocido de parametros ", valores[1])
			respuesta += "ERROR MOUNT, valor desconocido de parametros " + valores[1]
			//Si falta el valor del parametro actual lo reconoce como error e interrumpe el proceso
			return respuesta
		}

		if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		} else {
			fmt.Println("MKFS Error: Parametro desconocido: ", valores[0])
			return "MKFS Error: Parametro desconocido: " + valores[0] //por si en el camino reconoce algo invalido de una vez se sale
		}
	}

	if id != "" {
		var pathDico string
		var registro int //registro a eliminar

		eliminar := false
		//BUsca en struck de particiones montadas el id ingresado
		for i, montado := range Structs.Montadas {
			if montado.Id == id {
				eliminar = true
				pathDico = montado.PathM
				registro = i
			}
		}

		if eliminar {
			Disco, err := Herramientas.OpenFile(pathDico)
			if err != nil {
				return "ERROR REP OPEN FILE " + err.Error()
			}

			var mbr Structs.MBR
			// Read object from bin file
			if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
				return "ERROR REP READ FILE " + err.Error()
			}

			//Encontrar la particion en el disco
			for i := 0; i < 4; i++ {
				identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
				if identificador == id {
					name := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					var unmount Structs.Partition

					//Eliminar el id usando el id de la variable unmount
					mbr.Partitions[i].Id = unmount.Id
					copy(mbr.Partitions[i].Status[:], "I")

					//sobreescribir el mbr para guardar los cambios
					if err := Herramientas.WriteObject(Disco, mbr, 0); err != nil { //Sobre escribir el mbr
						return "ERROR UNMOUNT " + err.Error()
					}
					fmt.Println("Particion con nombre ", name, " desmontada correctamente")
					break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
				}
			}

			//elimina el la particion montada del struck
			Structs.Montadas = append(Structs.Montadas[:registro], Structs.Montadas[registro+1:]...)

			for _, montada := range Structs.Montadas {
				fmt.Println("Id " + string(montada.Id) + ", Disco: " + montada.PathM + "\n")
				//partMontadas += "Id "+ string(montada.Id)+ ", Disco: "+ montada.PathM+"\n"
			}

			for i := 0; i < 4; i++ {
				estado := string(mbr.Partitions[i].Status[:])
				if estado == "A" {
					//tmpMontadas:= "Particion: " + strconv.Itoa(i) + ", name: " +string(mbr.Partitions[i].Name[:]) + ", status: "+string(mbr.Partitions[i].Status[:])+", id: "+string(mbr.Partitions[i].Id[:])+", tipo: "+string(mbr.Partitions[i].Type[:])+", correlativo: "+ strconv.Itoa(int(mbr.Partitions[i].Correlative)) + ", fit: "+string(mbr.Partitions[i].Fit[:])+ ", start: "+strconv.Itoa(int(mbr.Partitions[i].Start))+ ", size: "+strconv.Itoa(int(mbr.Partitions[i].Size))
					//partMontadas += Herramientas.EliminartIlegibles(tmpMontadas)+"\n"
					fmt.Println("patcion: ", i, ", name: ", string(mbr.Partitions[i].Name[:]), ", status: "+string(mbr.Partitions[i].Status[:]))
				}
			}
		} else {
			fmt.Println("ERROR UNMOUNT: ID NO ENCONTRADO")
			return "ERROR UNMOUNT: ID NO ENCONTRADO"
		}

	} else {
		fmt.Println("ERROR UNMOUNT NO SE INGRESO PARAMETRO ID")
		return "ERROR UNMOUNT NO SE INGRESO PARAMETRO ID"
	}
	return respuesta
}
