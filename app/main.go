package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_taskHttpDelivery "github.com/fakecodes/gosample/task/delivery/http"
	_taskHttpDeliveryMiddleware "github.com/fakecodes/gosample/task/delivery/http/middleware"
	_taskRepo "github.com/fakecodes/gosample/task/repository/mongo"
	_taskUsecase "github.com/fakecodes/gosample/task/usecase"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	val := url.Values{}
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	dbConnection := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	// Get Client, Context, CancelFunc and err from connect method.
	client, ctx, cancel, err := connect(dbConnection)
	if err != nil {
		panic(err)
	}
	// Release resource when the main function is returned.
	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)

	e := echo.New()
	middL := _taskHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	tr := _taskRepo.NewMongoTaskRepository(client.Database(dbName))
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	tu := _taskUsecase.NewTaskUsecase(tr, timeoutContext)
	_taskHttpDelivery.NewTaskHandler(e, tu)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint
}

func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("context.timeout"))*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}
