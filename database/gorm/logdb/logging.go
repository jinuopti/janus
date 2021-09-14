package logdb

import (
    . "github.com/jinuopti/janus/log"
    "gorm.io/gorm"
    "time"
    "github.com/jinuopti/janus/database"
)

const Debug string = "DEBUG"
const Info string = "INFO"
const Warning string = "WARNING"
const Error string = "ERROR"

type Logging struct {
    gorm.Model
    AppName string
    LogType string
    Message string
}

var DbInfo *database.DbInfo

func DropTableLogging() {
    if DbInfo == nil {
        Loge("Database is nil")
        return
    }
    DbInfo.Db.Exec("DROP TABLE logging")
}

func CreateTableLogging() {
    if DbInfo == nil {
        var err error
        DbInfo, err = database.ConnNewDbFromConfig()
        if err != nil {
            Loge("ConnNewDbFromConfig Error : %s", err)
            return
        }
    }

    err := DbInfo.Db.AutoMigrate(&Logging{})
    if err != nil {
        Loge("Failed auto migrate table")
        return
    }
}

func WriteLog(log *Logging) {
    DbInfo.Db.Create(log)
}

func FindByTimeLogging(start time.Time, end time.Time) *Logging {
    if DbInfo == nil {
        var err error
        DbInfo, err = database.ConnNewDbFromConfig()
        if err != nil {
            Loge("ConnNewDbFromConfig Error : %s", err)
            return nil
        }
    }
    var log Logging
    DbInfo.Db.Where("created_at BETWEEN ? AND ?", start, end).Find(&log)

    Logd("AppName=%s, LogType=%s, Message=[%s], createAt=%s, updateAt=%s",
        log.AppName, log.LogType, log.Message, log.Model.CreatedAt.String(), log.Model.UpdatedAt.String())

    return &log
}

func InitLoggingTable() {
    var err error
    DbInfo, err = database.ConnNewDbFromConfig()
    if err != nil {
        return
    }
    CreateTableLogging()
}
