package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/e-commerce/common/config"
	"github.com/e-commerce/common/constant"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DBConnection  *sqlx.DB
	DBString      string
	RetryInterval int
	MaxConn       int
	doneChannel   chan bool
}

type MasterSlave struct {
	Master *DB
	Slave  *DB
}

var DBConnMap map[string]*MasterSlave

func InitDatabase(dbConfigs map[string]*config.DatabaseConfig) {

	DBConnMap = make(map[string]*MasterSlave)

	for k, cfg := range dbConfigs {
		masterDsn := cfg.Master
		slaveDsn := cfg.Slave

		Master := &DB{
			DBString:      masterDsn,
			RetryInterval: 10,
			MaxConn:       140,
			doneChannel:   make(chan bool),
		}
		Master.ConnectAndMonitor(constant.DriverMysql)

		Slave := &DB{
			DBString:      slaveDsn,
			RetryInterval: 10,
			MaxConn:       110,
			doneChannel:   make(chan bool),
		}
		Slave.ConnectAndMonitor(constant.DriverMysql)

		DBConnMap[k] = &MasterSlave{
			Master: Master,
			Slave:  Slave,
		}
	}
}

// GetOpenConnections return open connections to db
func (d *DB) GetOpenConnections() int64 {
	return d.GetOpenConnections()
}

// Connect to database
func (d *DB) Connect(driver string) error {
	var db *sqlx.DB
	var err error

	db, err = sqlx.Open(driver, d.DBString)

	if err != nil {
		log.Println("[Error]: DB open connection error", err.Error())
		return err
	}

	d.DBConnection = db
	err = db.Ping()

	if err != nil {
		log.Println("[Error]: DB connection error", err.Error())
		return err
	}

	db.SetMaxOpenConns(d.MaxConn)

	return err
}

// ConnectAndMonitor to database
func (d *DB) ConnectAndMonitor(driver string) {
	err := d.Connect(driver)

	if err != nil {
		log.Println("Not connected to database %s, trying", d.DBString)
	} else {
		log.Println("Success connecting to database %s", d.DBString)
	}

	ticker := time.NewTicker(time.Duration(d.RetryInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if d.DBConnection == nil {
					d.Connect(driver)
				} else {
					err := d.DBConnection.Ping()
					if err != nil {
						log.Println("[Error]: DB reconnect error", err.Error())
					}
				}
			case <-d.doneChannel:
				return
			}
		}
	}()
}

// DoneConnectAndMonitor to exit connect and monitor
func (d *DB) DoneConnectAndMonitor() {
	d.doneChannel <- true
}

//Prepare query for database queries
func (d *DB) Prepare(query string) *sql.Stmt {
	statement, err := d.DBConnection.Prepare(query)

	if err != nil {
		log.Fatalf("Failed to prepare query: %s. Error: %s", query, err.Error())
	}

	return statement
}

//Preparex query for database queries
func (d *DB) Preparex(query string) (*sqlx.Stmt, error) {
	if d == nil {
		log.Println("Failed to prepare query, database object is nil. Query: %s", query)
		return nil, nil
	}

	statement, err := d.DBConnection.Preparex(query)

	if err != nil {
		log.Println("Failed to preparex query: %s. Error: %s", query, err.Error())
		return nil, err
	}

	return statement, nil
}
