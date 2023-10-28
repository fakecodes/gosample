package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_taskHttpDelivery "github.com/fakecodes/gosample/task/delivery/http"
	_taskHttpDeliveryMiddleware "github.com/fakecodes/gosample/task/delivery/http/middleware"
	_taskRepo "github.com/fakecodes/gosample/task/repository/mysql"
	_taskUsecase "github.com/fakecodes/gosample/task/usecase"

	_ "github.com/go-sql-driver/mysql"

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
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	e := echo.New()
	middL := _taskHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	rr := _taskRepo.NewMysqlTaskRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	tu := _taskUsecase.NewTaskUsecase(rr, timeoutContext)
	_taskHttpDelivery.NewTaskHandler(e, tu)

	log.Fatal(e.Start(viper.GetString("server.address"))) //nolint
}