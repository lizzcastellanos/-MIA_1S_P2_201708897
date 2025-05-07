package Structs


type UserInfo struct {
	IdPart  string //identificar la particion del usuario
	IdGrp  	int32  //id del grupo al que pertenece el usuario
	IdUsr  	int32  //id del usuario
	Nombre 	string //saber que usuario es (identifica si es root o cualquir otro)
	Status 	bool   //si esta iniciada la sesion
	PathD	string	//Path del disco
}

var UsuarioActual UserInfo

func SalirUsuario(){
	UsuarioActual.IdGrp = 0
	UsuarioActual.IdPart = ""
	UsuarioActual.IdUsr = 0
	UsuarioActual.Nombre = ""
	UsuarioActual.Status = false
	UsuarioActual.PathD = ""
}


//Para almacenar la informacion del usuario con sesion iniciada

//Valores por defecto al crear un objeto de esta estructura
//Id = ""
//Status = false -> false no hay sesion iniciada. true sesion iniciada