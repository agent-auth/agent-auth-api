package main

import (
	"github.com/agent-auth/agent-auth-api/cmd"
)

// @host localhost:8002
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cmd.Execute()
}
