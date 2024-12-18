package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	_ "github.com/agent-auth/agent-auth-api/db/migrations" // import migrations
	"github.com/spf13/cobra"
	migrate "github.com/xakep666/mongo-migrate"
)

var (
	action  string
	message string
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "use migration tool",
	Long:  `migrate uses mongo-migrate migration tool under the hood supporting the same commands and an additional reset command`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrateCmd.Flags().StringVar(&action, "action", "", "Creates a new migration into the migrations folder")
	migrateCmd.Flags().StringVar(&message, "message", "", "Apply migrations up")
}

func run() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(getEnvInt("MONGODB_QUERY_TIMEOUT_SECONDS", 30))*time.Second,
	)
	defer cancel()

	opt := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(opt)
	if err != nil {
		log.Fatal(err.Error())
	}

	db := client.Database(os.Getenv("MONGODB_DATABASE"))

	migrate.SetDatabase(db)
	migrate.SetMigrationsCollection("migrations")

	switch action {
	case "new":
		if len(message) == 0 {
			log.Fatal("Provide message for new migration")
		}
		fName := fmt.Sprintf("./database/migrations/%s_%s.go", time.Now().Format("20060102150405"), strings.ReplaceAll(message, " ", "_"))
		from, err := os.Open("./database/migrations/20240101000000_template.go")
		if err != nil {
			log.Fatal("Migration template not found")
		}
		defer from.Close()

		to, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("New migration created: %s\n", fName)
	case "up":
		err = migrate.Up(ctx, 0)
	case "down":
		err = migrate.Down(ctx, 0)
	default:
		log.Fatal("Invalid operation")
	}

	if err != nil {
		log.Fatal(err.Error())
	}
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
