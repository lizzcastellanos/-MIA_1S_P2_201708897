package main

import (
	Admn "MIA_1S2025_P1_201708997/Comandos/AdminDiscos"
	Per "MIA_1S2025_P1_201708997/Comandos/AdminPermisos"
	AdmnA "MIA_1S2025_P1_201708997/Comandos/AdminSisArchivos"
	Reportes "MIA_1S2025_P1_201708997/Comandos/Rep"
	Usr "MIA_1S2025_P1_201708997/Comandos/Usuarios"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

// Estructura para manejar la entrada de texto del usuario
type Entrada struct {
	Text string `json:"text"`
}

// Estructura para la respuesta del servidor
type StatusResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func main() {
	// Configuración del endpoint principal
	http.HandleFunc("/analizar", procesarComandos)

	// Configuración de CORS para permitir solicitudes cruzadas
	corsHandler := cors.Default()

	// Iniciar el servidor en el puerto 8080
	fmt.Println("[Servidor] Iniciando en http://localhost:8080")
	http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux))
}

func procesarComandos(w http.ResponseWriter, r *http.Request) {
	var salida string
	w.Header().Set("Content-Type", "application/json")

	var estado StatusResponse

	if r.Method != http.MethodPost {
		estado = StatusResponse{Message: "Método HTTP no soportado", Type: "error"}
		json.NewEncoder(w).Encode(estado)
		return
	}

	var cmdEntrada Entrada
	if err := json.NewDecoder(r.Body).Decode(&cmdEntrada); err != nil {
		estado = StatusResponse{Message: "Formato JSON inválido", Type: "error"}
		json.NewEncoder(w).Encode(estado)
		return
	}

	// Procesamiento línea por línea de los comandos
	scanner := bufio.NewScanner(strings.NewReader(cmdEntrada.Text))
	for scanner.Scan() {
		lineaTexto := scanner.Text()
		if lineaTexto != "" {
			// Separar comandos de comentarios
			partes := strings.Split(lineaTexto, "#")
			if len(partes[0]) > 0 {
				fmt.Println("\n══════════════════════════════════════════════════════════════")
				fmt.Println("Procesando comando:", partes[0])
				salida += "══════════════════════════════════════════════════════════════\n"
				salida += "Comando: " + partes[0] + "\n"
				salida += interpretarComando(partes[0]) + "\n"
			}

			// Manejar comentarios si existen
			if len(partes) > 1 && partes[1] != "" {
				fmt.Println("Nota:", partes[1]+"\n")
				salida += "# " + partes[1] + "\n"
			}
		}
	}

	estado = StatusResponse{Message: salida, Type: "éxito"}
	json.NewEncoder(w).Encode(estado)
}

func interpretarComando(comando string) string {
	respuesta := ""
	// Normalizar el comando eliminando espacios redundantes
	comandoLimpio := strings.TrimRight(comando, " ")
	// Dividir comando principal de sus parámetros
	args := strings.Split(comandoLimpio, " -")

	switch strings.ToLower(args[0]) {
	// ------------------- Administración de discos -------------------
	case "mkdisk":
		if len(args) > 1 {
			respuesta = Admn.Mkdisk(args)
		} else {
			fmt.Println("[Error] mkdisk: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para mkdisk"
		}

	case "rmdisk":
		if len(args) > 1 {
			respuesta = Admn.Rmdisk(args)
		} else {
			fmt.Println("[Error] rmdisk: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para rmdisk"
		}

	case "mounted":
		respuesta = Admn.Mounted()

	case "fdisk":
		if len(args) > 1 {
			respuesta = Admn.Fdisk(args)
		} else {
			fmt.Println("[Error] fdisk: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para fdisk"
		}

	case "mount":
		if len(args) > 1 {
			respuesta = Admn.Mount(args)
		} else {
			fmt.Println("[Error] mount: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para mount"
		}

	case "unmount":
		if len(args) > 1 {
			respuesta = Admn.Unmoun(args)
		} else {
			fmt.Println("[Error] unmount: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para unmount"
		}

	// ----------- Administración de sistema de archivos -----------
	case "mkfs":
		if len(args) > 1 {
			respuesta = AdmnA.MKfs(args)
		} else {
			fmt.Println("[Error] mkfs: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para mkfs"
		}

	// ----------- Administración de usuarios y grupos -----------
	case "login":
		if len(args) > 1 {
			respuesta = Usr.Login(args)
		} else {
			fmt.Println("[Error] login: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para login"
		}

	case "logout":
		respuesta = Usr.Logout()

	case "mkgrp":
		if len(args) > 1 {
			respuesta = Usr.Mkgrp(args)
		} else {
			fmt.Println("[Error] mkgrp: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para mkgrp"
		}

	case "rmgrp":
		if len(args) > 1 {
			respuesta = Usr.Rmgrp(args)
		} else {
			fmt.Println("[Error] rmgrp: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para rmgrp"
		}

	case "mkusr":
		if len(args) > 1 {
			respuesta = Usr.Mkusr(args)
		} else {
			fmt.Println("[Error] mkusuario: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para mkusuario"
		}

	case "rmusr":
		if len(args) > 1 {
			respuesta = Usr.Rmusr(args)
		} else {
			fmt.Println("[Error] rmusuario: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para rmusuario"
		}

	case "chgrp":
		if len(args) > 1 {
			respuesta = Usr.Chgrp(args)
		} else {
			fmt.Println("[Error] chgrp: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para chgrp"
		}

	// ----------- Gestión de archivos y permisos -----------
	case "mkfile":
		if len(args) > 1 {
			respuesta = Per.MKfile(args)
		} else {
			fmt.Println("[Error] mkfile: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para mkfile"
		}

	case "cat":
		if len(args) > 1 {
			respuesta = Per.Cat(args)
		} else {
			fmt.Println("[Error] cat: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para cat"
		}

	case "mkdir":
		if len(args) > 1 {
			respuesta = Per.Mkdir(args)
		} else {
			fmt.Println("[Error] mkdir: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para mkdir"
		}

	// ------------------- Reportes -------------------
	case "rep":
		if len(args) > 1 {
			respuesta = Reportes.Rep(args)
		} else {
			fmt.Println("[Error] rep: parámetros insuficientes")
			respuesta = "Error: Faltan parámetros para rep"
		}

	case "":
		// Línea vacía, no hacer nada
		return ""

	default:
		fmt.Println("[Error] Comando no implementado:", args[0])
		return "Error: Comando no reconocido"
	}

	return respuesta
}
