package configure

import (
    "fmt"
    "gopkg.in/ini.v1"
    "strings"
    "errors"
)

const (
    SectCore        = "CORE"
    SectTimer       = "TIMER"
    SectNet         = "NET"
    SectLog         = "LOG"
    SectDatabase    = "DATABASE"

    SectKredigo     = "Kredigo" // Kredigo: insert your section name string
)

var (
    Config ConfigInfo
)

type ConfigInfo struct {
    File  *ini.File
    Path  string
    Value *Values
}

type Values struct {
    // Core
    Core *ValueCore
    Time *ValueTime
    Net  *ValueNet
    Log  *ValueLog
    Db   *ValueDatabase

    // Application
    Kredigo *ValueKredigo // Kredigo: insert your structure variable
}

type ValueCore struct {
    sectionName string
}

type ValueTime struct {
    sectionName string

    IdleTimeout uint64
}

type ValueNet struct {
    sectionName string

    EnableHttp   bool
    EnableWebsocket bool
    EnableGrpc   bool
    EnableTcp    bool
    EnableSerial bool
    EnableMqueue bool
}

type ValueLog struct {
    sectionName string

    EnableDebug bool
    EnableInfo  bool
    EnableError bool
    LogFile     string
    MaxSize     int
    MaxBackups  int
    MaxAge      int
    LocalTime   bool
    Compress    bool
}

type ValueDatabase struct {
    sectionName string

    Type      string
    UserId    string
    Password  string
    IpAddress string
    Port      int
    DbName    string
    DbNameLog string
}

func NewValues() *Values {
    if Config.Value == nil {
        Config.Value = &Values{}

        Config.Value.Core = &ValueCore{}
        Config.Value.Time = &ValueTime{}
        Config.Value.Net = &ValueNet{}
        Config.Value.Log = &ValueLog{}
        Config.Value.Db = &ValueDatabase{}

        Config.Value.Kredigo = &ValueKredigo{} // Kredigo: insert your new structure
    }
    return Config.Value
}

func GetConfig() *Values {
    return Config.Value
}

func (c *Values) chkCfgFile(filePath string) (*Values, error) {
    if Config.File == nil {
        var err error
        Config.File, err = ini.Load(filePath)
        if err != nil {
            return nil, err
        }
        Config.Path = filePath
        // fmt.Printf("Load INI file successed [%s]\n", filePath)
    }
    return c, nil
}

func (c *Values) GetValueALL(filePath string) (*Values, error) {
    var err error
    _, err = c.chkCfgFile(filePath)
    if err != nil {
        return nil, err
    }

    _, err = c.GetValueCore(filePath)
    if err != nil {
        return nil, err
    }
    _, err = c.GetValueTimer(filePath)
    if err != nil {
        return nil, err
    }
    _, err = c.GetValueNet(filePath)
    if err != nil {
        return nil, err
    }
    _, err = c.GetValueLog(filePath)
    if err != nil {
        return nil, err
    }
    _, err = c.GetValueDatabase(filePath)
    if err != nil {
        return nil, err
    }

    _, err = c.GetValueKredigo(filePath) // Kredigo: insert your option parsing fuction
    if err != nil {
        return nil, err
    }

    return c, nil
}

func (c *Values) GetValueCore(filePath string) (*Values, error) {
    _, err := c.chkCfgFile(filePath)
    if err != nil {
        return nil, err
    }

    c.Core.sectionName = SectCore

    // fmt.Printf("Core Config Values: %+v\n", c.Log)

    return c, nil
}

func (c *Values) GetValueTimer(filePath string) (*Values, error) {
    _, err := c.chkCfgFile(filePath)
    if err != nil {
        return nil, err
    }

    c.Time.sectionName = SectTimer
    c.Time.IdleTimeout = Config.File.Section(SectTimer).Key("IdleTimeout").MustUint64(10)

    // fmt.Printf("Time Config Values: %+v\n", c.Log)

    return c, nil
}

func (c *Values) GetValueNet(filePath string) (*Values, error) {
    _, err := c.chkCfgFile(filePath)
    if err != nil {
        return nil, err
    }

    c.Net.sectionName = SectNet
    c.Net.EnableHttp = Config.File.Section(SectNet).Key("EnableHttp").MustBool(false)
    c.Net.EnableWebsocket = Config.File.Section(SectNet).Key("EnableWebsocket").MustBool(false)
    c.Net.EnableGrpc = Config.File.Section(SectNet).Key("EnableGrpc").MustBool(false)
    c.Net.EnableTcp = Config.File.Section(SectNet).Key("EnableTcp").MustBool(false)
    c.Net.EnableSerial = Config.File.Section(SectNet).Key("EnableSerial").MustBool(false)
    c.Net.EnableMqueue = Config.File.Section(SectNet).Key("EnableMqueue").MustBool(false)

    // fmt.Printf("Net Config Values: %+v\n", c.Log)

    return c, nil
}

func (c *Values) GetValueLog(filePath string) (*Values, error) {
    _, err := c.chkCfgFile(filePath)
    if err != nil {
        return nil, err
    }

    c.Log.sectionName = SectLog
    c.Log.EnableDebug = Config.File.Section(SectLog).Key("EnableDebug").MustBool(true)
    c.Log.EnableInfo = Config.File.Section(SectLog).Key("EnableInfo").MustBool(true)
    c.Log.EnableError = Config.File.Section(SectLog).Key("EnableError").MustBool(true)
    c.Log.LogFile = Config.File.Section(SectLog).Key("LogFile").MustString("sample.log")
    c.Log.MaxSize = Config.File.Section(SectLog).Key("MaxSize").MustInt(10)
    c.Log.MaxBackups = Config.File.Section(SectLog).Key("MaxBackups").MustInt(1)
    c.Log.MaxAge = Config.File.Section(SectLog).Key("MaxAge").MustInt(10)
    c.Log.LocalTime = Config.File.Section(SectLog).Key("LocalTime").MustBool(true)
    c.Log.Compress = Config.File.Section(SectLog).Key("Compress").MustBool(false)

    // fmt.Printf("Log Config Values: %+v\n", c.Log)

    return c, nil
}

func (c *Values) GetValueDatabase(filePath string) (*Values, error) {
    _, err := c.chkCfgFile(filePath)
    if err != nil {
        return nil, err
    }

    c.Db.sectionName = SectDatabase
    c.Db.Type = Config.File.Section(SectDatabase).Key("Type").MustString("mysql")
    c.Db.UserId = Config.File.Section(SectDatabase).Key("UserId").MustString("admin")
    c.Db.Password = Config.File.Section(SectDatabase).Key("Password").MustString("admin")
    c.Db.IpAddress = Config.File.Section(SectDatabase).Key("IpAddress").MustString("127.0.0.1")
    c.Db.Port = Config.File.Section(SectDatabase).Key("Port").MustInt(3306)
    c.Db.DbName = Config.File.Section(SectDatabase).Key("DbName").MustString("kredigo")
    c.Db.DbNameLog = Config.File.Section(SectDatabase).Key("DbNameLog").MustString("logs")

    // fmt.Printf("Dayabase Config Values: %+v\n", c.Db)

    return c, nil
}

func (c *Values) PrintValues(sectName string) (prtStr string) {
    isAll := false

    sectName = strings.ToUpper(sectName)
    switch sectName {
    case SectCore:
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", sectName, c.Core)
    case SectTimer:
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", sectName, c.Time)
    case SectLog:
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", sectName, c.Log)
    case SectNet:
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", sectName, c.Net)
    case SectDatabase:
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", sectName, c.Db)

    case SectKredigo: // application configure value
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", sectName, c.Kredigo)
    default:
        isAll = true
    }

    if isAll {
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", SectCore, c.Core)
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", SectTimer, c.Time)
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", SectLog, c.Log)
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", SectNet, c.Net)
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", SectDatabase, c.Db)
        prtStr += fmt.Sprintf(" * %s Config Values: [%+v]\n", SectKredigo, c.Kredigo) // application configure value
    }
    return
}

// SetConfigValue 설정값을 변경한다
func (c *Values) SetConfigValue(section string, key string, value string) error {
    if Config.File == nil {
        return errors.New("error: Config file is nil")
    }
    Config.File.Section(section).Key(key).SetValue(value)
    return nil
}

// SaveConfigFile 변경한 설정값을 파일로 저장한다
func (c *Values) SaveConfigFile() error {
    err := Config.File.SaveTo(Config.Path)
    if err != nil {
        return err
    }
    return nil
}
