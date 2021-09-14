package configure

// ValueKredigo structure
type ValueKredigo struct {
	sectionName string
	// Insert configure values here
	Enabled 				bool

	// Scheduler
	SchedulerMinute     bool    // 매 분 00초 실행 스케쥴러 사용여부
	SchedulerHour       bool    // 매 시 00분 실행 스케쥴러 사용여부
	SchedulerDay        bool    // 매 일 00:00 실행 스케쥴러 사용여부

	EnableHttpServer 		bool
	EnableWebsocketServer   bool
	EnableGrpcServer		bool
	EnableTcpServer	 		bool
	EnableSerial	 		bool
	EnableMqueue			bool

	HttpServerPort 	 		string
	GrpcServerPort			string

	TcpServerPort 	        string
	TcpClientPort 	        string

	SerialPort      		string

	EnableSwagger           bool
	SwaggerAddr             string
}

func (c *Values) GetValueKredigo(filePath string) (*Values, error) {
	_, err := c.chkCfgFile(filePath)
	if err != nil {
		return nil, err
	}

	c.Kredigo.sectionName = SectKredigo
	c.Kredigo.Enabled = Config.File.Section(SectKredigo).Key("Enabled").MustBool(false)

	// Scheduler
	c.Kredigo.SchedulerMinute = Config.File.Section(SectKredigo).Key("SchedulerMinute").MustBool(false)
	c.Kredigo.SchedulerHour = Config.File.Section(SectKredigo).Key("SchedulerHour").MustBool(false)
	c.Kredigo.SchedulerDay = Config.File.Section(SectKredigo).Key("SchedulerDay").MustBool(false)

	// Insert configure values here
	c.Kredigo.EnableHttpServer = Config.File.Section(SectKredigo).Key("EnableHttpServer").MustBool(false)
	c.Kredigo.EnableWebsocketServer = Config.File.Section(SectKredigo).Key("EnableWebsocketServer").MustBool(false)
	c.Kredigo.EnableGrpcServer = Config.File.Section(SectKredigo).Key("EnableGrpcServer").MustBool(false)
	c.Kredigo.EnableTcpServer = Config.File.Section(SectKredigo).Key("EnableTcpServer").MustBool(false)
	c.Kredigo.EnableSerial = Config.File.Section(SectKredigo).Key("EnableSerial").MustBool(false)
	c.Kredigo.EnableMqueue = Config.File.Section(SectKredigo).Key("EnableMqueue").MustBool(false)

	c.Kredigo.HttpServerPort = Config.File.Section(SectKredigo).Key("HttpServerPort").MustString("1323")
	c.Kredigo.GrpcServerPort = Config.File.Section(SectKredigo).Key("GrpcServerPort").MustString("1324")

	c.Kredigo.TcpServerPort = Config.File.Section(SectKredigo).Key("TcpServerPort_Kredigo").MustString("9900")
	c.Kredigo.TcpClientPort = Config.File.Section(SectKredigo).Key("TcpClientPort_Kredigo").MustString("9901")

	c.Kredigo.SerialPort = Config.File.Section(SectKredigo).Key("SerialPort_Kredigo").MustString("/dev/ttyUSB0")

	c.Kredigo.EnableSwagger = Config.File.Section(SectKredigo).Key("EnableSwagger").MustBool(false)
	c.Kredigo.SwaggerAddr = Config.File.Section(SectKredigo).Key("SwaggerAddr").MustString("localhost")

	return c, nil
}
