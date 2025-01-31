package main
// added comment to test github actions
import (
	"github.com/owjoel/client-factpack/apps/auth/pkg/web"
	_ "github.com/owjoel/client-factpack/apps/auth/docs"
)

// 	Swagger
//	@title			client-factpack/auth
//	@version		1.0
//	@description	Authentication service for managing auth flows
//	@host			localhost:8080
//	@BasePath		/api/v1


func main() {
	web.Run()
}