package appmmain

import (
	"github.com/jinuopti/janus/communication/grpc"
	"github.com/jinuopti/janus/communication/http"
	"github.com/jinuopti/janus/communication/socket"
	. "github.com/jinuopti/janus/configure"
	. "github.com/jinuopti/janus/log"
	"sync"
	"github.com/jinuopti/janus/database/gorm/userdb"
	"github.com/jinuopti/janus/database/gorm/logdb"
	utility "github.com/jinuopti/janus/library"
	gormdb "github.com/jinuopti/janus/database/gorm"
)

func Run(conf *Values) {
	if conf.Kredigo.Enabled == false {
		Logi("Kredigo Application Disabled")
		return
	}
	Logi("Kredigo Application Started")

	wg := new(sync.WaitGroup)

	// Init DATABASE
	err := gormdb.InitSingletonDB()
	if err == nil {
		userdb.InitUserTable()
		logdb.InitLoggingTable()
	}

	// HTTP
	if conf.Net.EnableHttp && conf.Kredigo.EnableHttpServer {
		go httpserver.HttpServer(conf.Kredigo.HttpServerPort)
	}

	// Websocket
	if conf.Net.EnableWebsocket && conf.Kredigo.EnableWebsocketServer {

	}

	// gRPC
	if conf.Net.EnableGrpc && conf.Kredigo.EnableGrpcServer {
		go grpcserver.GrpcServer(conf.Kredigo.GrpcServerPort)
	}

	// TCP Socket
	if conf.Net.EnableTcp && conf.Kredigo.EnableTcpServer {
		socket.TcpServerRun("Kredigo", conf.Kredigo.TcpServerPort)
	}

	// Serial
	if conf.Net.EnableSerial && conf.Kredigo.EnableSerial {

	}

	// Message Queue
	if conf.Net.EnableMqueue && conf.Kredigo.EnableMqueue {

	}

	// Start Scheduler
	if conf.Kredigo.SchedulerMinute {
		go utility.SchedulerMin(SchedulerMin)
	}
	if conf.Kredigo.SchedulerHour {
		go utility.SchedulerHour(SchedulerHour)
	}
	if conf.Kredigo.SchedulerDay {
		go utility.SchedulerDay(SchedulerDay)
	}

	wg.Wait()
}
