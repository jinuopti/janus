package userdb

import (
    . "github.com/jinuopti/janus/log"
    "github.com/jinuopti/janus/database"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    UserId    string `gorm:"not null; uniqueIndex; size:20"`
    Password  string `gorm:"not null"`
    Name      string `gorm:"not null"`
    Email     string
    UserLevel string
}

var DbInfo *database.DbInfo

func GetDbInfo() (dbInfo *database.DbInfo) {
    return DbInfo
}

func DropTableUser() {
    if DbInfo == nil || DbInfo.Db == nil {
        Loge("Database is nil")
        return
    }

    DbInfo.Db.Exec("DROP TABLE users")
}

func CreateTableUser() {
    var err error

    if DbInfo == nil {
        DbInfo, err = database.ConnNewDbFromConfig()
        if err != nil {
            Loge("ConnNewDbFromConfig Error : %s", err)
            return
        }
    }

    err = DbInfo.Db.AutoMigrate(&User{})
    if err != nil {
        Loge("Failed auto migrate table : %s", err)
        return
    }
}

func InsertUser(user *User) {
    DbInfo.Db.Create(user)
}

func FindByIdUser(id string) *User {
    if DbInfo == nil {
        var err error
        DbInfo, err = database.ConnNewDbFromConfig()
        if err != nil {
            Loge("ConnNewDbFromConfig Error : %s", err)
            return nil
        }
    }
    var user User
    DbInfo.Db.First(&user, "user_id = ?", id)

    Logd("User id=%s, pass=%s, name=%s, email=%s, level=%s, createAt=%s, updateAt=%s",
        user.UserId, user.Password, user.Name, user.Email, user.UserLevel,
        user.Model.CreatedAt.String(), user.Model.UpdatedAt.String())

    return &user
}

func InitUserTable() {
    var err error
    DbInfo, err = database.ConnNewDbFromConfig()
    if err != nil {
        return
    }
    CreateTableUser()
}
