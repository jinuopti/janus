package processdb

import (
    gormdb "github.com/jinuopti/janus/database/gorm"
    . "github.com/jinuopti/janus/log"
    "errors"
    "gorm.io/gorm"
)

const ProcessTableName = "process"

type Processes struct {
    gorm.Model
    Name        string  `gorm:"not null; uniqueIndex; size:40"`  // unique name (ex: "PshmIOC")
    AppName     string  // "HMI", "DT", ...
    Type        string  // "IOC", "Alarm", ...
    Pid         int     // Last pid
    Cmd         string  `gorm:"not null"` // "/path/to/process"
    Args        string  // "-i -s 100"
    Running     bool    `gorm:"not null"` // true: running(start on boot), false: stopped
    AutoRestart bool    // true: autorestart
    Description string  `gorm:"not null; type:varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci"` // "PSHM IOC"
}
func (Processes) TableName() string {
    return ProcessTableName
}

func CreateTableProcess() {
    var err error

    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return
    }

    err = gormdb.MainDB.Db.AutoMigrate(&Processes{})
    if err != nil {
        Loge("Failed auto migrate table : %s", err)
        return
    }
}

func GetProcess(name string) (*Processes, error) {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return nil, errors.New("MainDB is nil")
    }
    var proc Processes
    result := gormdb.MainDB.Db.First(&proc, "name = ?", name)
    if result.Error != nil {
        Logd("DB First result error: %s", result.Error)
        return nil, result.Error
    }

    return &proc, nil
}

func GetAppProcess(name string) (*Processes, error) {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return nil, errors.New("MainDB is nil")
    }
    var proc Processes
    result := gormdb.MainDB.Db.First(&proc, "app_name = ?", name)
    if result.Error != nil {
        Logd("DB First result error: %s", result.Error)
        return nil, result.Error
    }

    return &proc, nil
}

func GetProcessAll() (retProc []*Processes, err error) {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return nil, errors.New("MainDB is nil")
    }
    var procs []*Processes
    gormdb.MainDB.Db.Model(&Processes{}).Find(&procs)
    for _, p := range procs {
        retProc = append(retProc, p)
    }
    return retProc, nil
}

func InsertProcess(proc *Processes) error {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return errors.New("MainDB is nil")
    }
    p, _ := GetProcess(proc.Name)
    if p == nil {
        gormdb.MainDB.Db.Create(proc)
    } else if p.Pid == proc.Pid {
        _ = UpdateProcess(proc)
    }
    return nil
}

func UpdateProcess(proc *Processes) error {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return errors.New("MainDB is nil")
    }
    gormdb.MainDB.Db.Model(&Processes{}).Where("name = ?", proc.Name).Updates(proc)
    return nil
}

func UpdateProcessStatus(name string, running bool) error {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return errors.New("MainDB is nil")
    }
    gormdb.MainDB.Db.Model(&Processes{}).Where("name = ?", name).Update("running", false)
    Logd("[%s] Update running status to %d", name, running)
    return nil
}

func DeleteProcess(name string) error {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return errors.New("MainDB is nil")
    }
    gormdb.MainDB.Db.Model(&Processes{}).Where("name = ?", name).Unscoped().Delete(&Processes{})
    return nil
}

func InitProcessTable() {
    if gormdb.MainDB == nil {
        Logd("gormdb.MainDB is nil")
        return
    }
    CreateTableProcess()
    Logd("DB: Init User Table")
}
