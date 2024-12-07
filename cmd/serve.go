package cmd

import (
	"context"

	"github.com/agent-auth/agent-auth-api/db/mongodb"
	"github.com/agent-auth/agent-auth-api/db/redis_dal"
	"github.com/agent-auth/agent-auth-api/db/redisdb"
	"github.com/agent-auth/agent-auth-api/web/server"
	"github.com/agent-auth/common-lib/pkg/logger"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.NewLogger()

		go func() {
			logger.Info("Starting mongo client")
			_ = mongodb.NewMongoClient()
		}()

		go func() {
			logger.Info("Starting redis client")
			_ = redisdb.NewRedisClient()
		}()

		go func() {
			logger.Info("Starting initial sync of roles")
			redis_dal.NewRedisRolesDal().SyncRolesCollection(context.Background())
		}()

		server := server.NewServer()
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
