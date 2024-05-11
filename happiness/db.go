package happiness

import (
	"fmt"
	"github.com/Lyearn/mgod"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func InitConnect() {
	_ = godotenv.Load()
	url := os.Getenv("MONGO_URL")
	dbName := os.Getenv("MONGO_DB_NAME")
	cfg := &mgod.ConnectionConfig{
		Timeout: 5000 * time.Second,
	}
	opts := options.Client().ApplyURI(url)

	err := mgod.ConfigureDefaultConnection(cfg, dbName, opts)
	if err != nil {
		fmt.Println(err)
	}
}
