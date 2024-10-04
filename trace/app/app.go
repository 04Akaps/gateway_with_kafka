package app

import (
	"github.com/04Akaps/gateway_with_kafka.git/config"
	"github.com/04Akaps/gateway_with_kafka.git/global"
	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/kafka"
	"github.com/04Akaps/gateway_with_kafka.git/trace/repository/mysql"
	"github.com/04Akaps/gateway_with_kafka.git/trace/service"
	"log"
	"time"
)

type App struct {
	sql     *mysql.MySQL
	kafka   *kafka.Kafka
	service *service.Service

	cfg config.TraceCfg
}

func NewApp(cfg config.TraceCfg) *App {
	app := &App{
		cfg: cfg,
	}

	log.Println("start initialize Trace Moudle")

	var err error

	if app.sql, err = mysql.NewMySQL(cfg); err != nil {
		log.Println("Failed to connect db")
		panic(err)
	} else if app.kafka, err = kafka.NewKafka(cfg); err != nil {
		log.Println("Failed to connect kafka")
		panic(err)
	} else {
		app.service = service.NewService(cfg, app.kafka, app.sql)
		app.service.PollingAll()
	}

	return app
}

func (a *App) Run() {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("trace module started")
	global.InitializeWork.Done()

	for {
		select {
		case <-ticker.C:
			log.Println("working channels", "count", a.service.TotalChannels())
		}
	}
}
