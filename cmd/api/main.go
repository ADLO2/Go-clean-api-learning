package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/thienkb1123/go-clean-arch/config"
	"github.com/thienkb1123/go-clean-arch/internal/server"
	"github.com/thienkb1123/go-clean-arch/pkg/cache/redis"
	"github.com/thienkb1123/go-clean-arch/pkg/database/mongodb"
	"github.com/thienkb1123/go-clean-arch/pkg/database/mysql"
	"github.com/thienkb1123/go-clean-arch/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	rand.Seed(time.Now().Unix())
	log.Println("Starting api server")
	f, err := os.Create("myprogram.prof")
	if err != nil {
	
	fmt.Println(err)
	return
	
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// token, _ := utils.GenerateJWTToken(&models.User{
	// 	UserID: uuid.New(),
	// }, cfg)

	// fmt.Println("token: ", token)

	ctx := context.Background()
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof(ctx, "AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Repository
	mysqlDB, err := mysql.New(&cfg.MySQL)
	if err != nil {
		appLogger.Fatalf(ctx, "MySQL init: %s", err)
	}

	rdb, err := redis.NewClient(&cfg.Redis)
	if err != nil {
		appLogger.Fatalf(ctx, "RedisCluster init: %s", err)
	}

	mongodbClient, err := mongodb.New(&cfg.MySQL)
	if err != nil {
		appLogger.Fatalf(ctx, "Mongodb client init: %s", err)
	}

	if err = mongodbClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println()
		fmt.Println("PING")
		fmt.Println()
	}

	s := server.NewServer(
		cfg,
		mysqlDB,
		server.Redis(rdb),
		server.Logger(appLogger),
		server.Mongodb(mongodbClient),
	)
	
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
