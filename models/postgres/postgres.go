package postgres

import (
	//"errors"
	"fmt"

	//"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	logdebug "log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"
	//"strconv"
	"time"
	"os"
)

type Bugreport struct {
	//gorm.Model
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index;index:natural_unique_key,unique"`

	FileName             string
	Reporter             string
	UserName             string
	UploadSize           int64
	DateUploaded         string
	DateCreated          string
	SerialNumber         string `gorm:"index:natural_unique_key,unique"`
	HardwareRevision     string
	BuildFingerprint     string `gorm:"index:natural_unique_key,unique"`
	ArgusID              string
	NumberOfBoots        int
	TotalNumberOfCrashes int
	GrafanaDashboardLink string
	ConcatenatedFileLink string
	FilesConcatenated    int
	LogLinesParsed       int
	OriginalFile         bool `gorm:"default:false"`
	JiraLink             string
	IniTime              time.Time       `gorm:"type:timestamp;column:iniTime","index:natural_unique_key,unique"`
	EndTime              time.Time       `gorm:"type:timestamp;column:endTime"`
	Duration             int             `gorm:"index:natural_unique_key,unique"`
	TangoStatus          int             `gorm:"default:0"`
	States		     []State        `gorm:"constraint:OnDelete:CASCADE"`
	ConfigFile	     string
}


var Db *gorm.DB

func NewConnection() error {

	if Db != nil {
		return nil
	}

	name := viper.GetString("postgres.name")
	host := viper.GetString("postgres.host")
	port := viper.GetString("postgres.port")
	user := viper.GetString("postgres.user")
	password := viper.GetString("postgres.password")


	newLogger := logger.New(
	  logdebug.New(os.Stdout, "\r\n", logdebug.LstdFlags), // io writer
	  logger.Config{
	    SlowThreshold:              time.Second,   // Slow SQL threshold
	    LogLevel:                   logger.Info, // Log level
	    IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
	    Colorful:                  false,          // Disable color
	  },
	)

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port)

	var err error
	// open connection
	// Note: if we can't make a connection to the database no err is being returned
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger,})
	//Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true,})

	if err != nil {
		return err
	} else {
		log.Infof("postgres: successfully connected to database on %s %s %s", host, port, name)
	}
	//Setup()
	return nil
}

