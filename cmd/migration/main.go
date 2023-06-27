package main

import (
	postgres "uploader/models/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"	
	"os"		
)

func main() {

	initialize()

	postgres.NewConnection()
	Setup()
}

func Setup() {
	allModels := []interface{}{ &postgres.Bugreport{}, &postgres.State{} }

	if err := postgres.Db.Migrator().DropTable(allModels...); err != nil {
		log.Printf("Failed to drop table, got error %v\n", err)
		os.Exit(1)
	}

	if err := postgres.Db.AutoMigrate(allModels...); err != nil {
		log.Printf("Failed to auto migrate, but got error %v\n", err)
		os.Exit(1)
	}

	/*DbLC.Migrator().DropTable(
		&BRPartition{},
		&File{},
		&Label{},
		&Section{},
		&Priority{},
		&BootFolder{},		
		&Content{},
	)

	DbLC.AutoMigrate(
		&BRPartition{},
		&File{},
		&Label{},
		&Section{},
		&Priority{},
		&BootFolder{},		
		&Content{},		
	)


	log.Info("postgres logcat: migrated tables")*/
}
func initialize() {


	// setup logger
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// load configuration
	viper.SetConfigName("argus-uploader-config") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs") // config file path
	err := viper.ReadInConfig()

	if err != nil {
		log.Error("server: failed to read config file")
		log.Fatal(err)
	}

}
