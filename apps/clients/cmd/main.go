package main

import "github.com/owjoel/client-factpack/apps/clients/pkg/web"
import _ "github.com/owjoel/client-factpack/apps/clients/docs"

// 	Swagger
//	@title			client-factpack/clients
//	@version		1.0
//	@description	Client resource module. Manages manually typed and compiled online data of prospective clients
//	@host			localhost:8080
//	@BasePath		/api/v1

func main() {
	web.Run()
}
