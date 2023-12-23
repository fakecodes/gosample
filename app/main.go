package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_taskHttpDelivery "github.com/fakecodes/gosample/task/delivery/http"
	_taskHttpDeliveryMiddleware "github.com/fakecodes/gosample/task/delivery/http/middleware"
	_taskRepo "github.com/fakecodes/gosample/task/repository/postgres"
	_taskUsecase "github.com/fakecodes/gosample/task/usecase"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
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
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dbHost := viper.GetString(`database.host`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbPort := viper.GetString(`database.port`)
	dbName := viper.GetString(`database.name`)
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbPort, dbName)

	// Get DB Context and err from connect method.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Release resource when the main function is returned.
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _taskHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	tr := _taskRepo.NewPostgresTaskRepository(db)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	tu := _taskUsecase.NewTaskUsecase(tr, timeoutContext)
	_taskHttpDelivery.NewTaskHandler(e, tu)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint
}
