package health

import (
	"context"
	"net/http"

	"github.com/agent-auth/agent-auth-api/db/mongodb"
	"github.com/agent-auth/agent-auth-api/db/redisdb"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	_ "github.com/agent-auth/agent-auth-api/web/renderers" // swag
	"github.com/go-chi/render"
)

type health struct {
	redis *redis.Client
	mongo *mongo.Database
}

// NewHealth returns health impl
func NewHealth() Health {
	return &health{
		redis: redisdb.NewRedisClient(),
		mongo: mongodb.NewMongoClient(),
	}
}

// @Summary Get health of the service
// @Description It returns the health of the service
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /health [get]
// GetHealth returns heath of service, can be extended if
// service is running on multile instances
func (h *health) GetHealth(w http.ResponseWriter, r *http.Request) {
	_, err := h.redis.Ping(context.Background()).Result()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError)
	}

	var result bson.M
	err = h.mongo.RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Decode(&result)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError)
	}
}
