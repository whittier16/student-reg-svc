package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/whittier16/student-reg-svc/internal/pkg/exception"
)

// MySQL exposes the different methods supported by the storage application
type MySQL struct {
	DBClient *sql.DB
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// New gets new instance of the MySQL database
func New(user string, pass string, host string, port string, DBName string) (*MySQL, error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		pass,
		host,
		port,
		DBName,
	)
	dbClient, err := sql.Open(`mysql`, connection)
	if err != nil {
		fmt.Println("Error occured: ", err)
		exception.PanicIfNeeded(err)
	}

	return &MySQL{
		DBClient: dbClient,
	}, nil
}

// GetDb gets MySQL database connection
func (m *MySQL) GetDb() *sql.DB {
	return m.DBClient
}

// CloseDb closes MySQL database connection
func (m *MySQL) CloseDb() {
	err := m.DBClient.Close()
	if err != nil {
		log.Error("Error occurred: ", err)
		return
	}
}
