package sqltool

import (
	properties "TMManager/src/Properties"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var sqlDB *sql.DB

func InitialSql() error {
	user, err := properties.Get("SqlUsername")
	if err != nil {
		return err
	}
	password, err := properties.Get("SqlPassword")
	if err != nil {
		return err
	}
	ip, err := properties.Get("SqlIP")
	if err != nil {
		return err
	}
	port, err := properties.Get("SqlPort")
	if err != nil {
		return err
	}
	database, err := properties.Get("Database")
	if err != nil {
		return err
	}
	connectStr := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(connectStr), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err = db.DB()
	if err != nil {
		return err
	}
	MaxIdleConns, err := properties.GetInt("MaxIdleConns")
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(MaxIdleConns)
	MaxOpenConns, err := properties.GetInt("MaxOpenConns")
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(MaxOpenConns)
	MaxLifeTime, err := properties.GetInt("MaxLifeTime")
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(time.Duration(MaxLifeTime))
	return nil
}

func CreateTables() error {
	err := db.AutoMigrate(
		&Problem{},
		&Judge{},
		&Judge_case{},
		&Contest{},
		&Contest_problem{},
		&Contest_register{},
		&Contest_record{},
		&Contest_judge{},
		&Contest_judge_case{},
		&User_info{},
	)
	if err == nil {
		fmt.Println("数据库表初始化成功")
	}
	return err
}

func QueryJudge(JudgeID string) Judge {
	judge := Judge{}
	db.Preload("Problem").Find(&judge, JudgeID)
	return judge
}

func SaveJudge(judge *Judge) {
	db.Save(judge)
}

func CreateJudgeCases(judegCases *[]Judge_case) {
	db.Create(judegCases)
}
