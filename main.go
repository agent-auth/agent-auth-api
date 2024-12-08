package main

import (
	"github.com/agent-auth/agent-auth-api/cmd"
)

// @title Agent Auth API Documentation
// @version 2.0
// @description Agent Auth API Documentation

// @contact.name API Support
// @contact.url http://xyz.ai
// @contact.email hello@xyz.ai

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8002
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cmd.Execute()
}
