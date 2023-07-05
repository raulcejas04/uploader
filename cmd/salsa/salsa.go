package main

import(
	"fmt"
	postgres "uploader/models/postgres"	
	"net/http"	
	//"database/sql"	
	"encoding/json"	
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"io"	
	"strconv"
	"strings"
	"time"	
	"uploader/pkg/encoding"
	"path/filepath"
	"path"
)

func handleRequests() {
	r := mux.NewRouter()

	r.HandleFunc("/uploader/state", HandlerState).Methods("GET")
	r.HandleFunc("/uploader/open", HandlerOpenSession).Methods("GET")
	r.HandleFunc("/uploader/transfer", HandlerTransfer).Methods("POST")
	
	address := viper.GetString("transferServer.address")
	port := viper.GetString("transferServer.port")
	fmt.Println( " address ", address, port )

	srv := &http.Server{
		Addr:    address + ":" + port,
		Handler: r,
	}
	srv.ListenAndServe()

}

func HandlerState(w http.ResponseWriter, r *http.Request) {


	hash:=r.FormValue("hash")
	if hash=="" {
		http.Error(w, "Missing hash",  http.StatusInternalServerError)
		return
	}

	dummy:=string(encoding.Decode(hash))
	res := strings.Split(dummy, "|")	

	id:= res[0]
	//username:=res[1] //to compare with token
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
	}

	results:=postgres.GetStates( idInt )
	
	jsonResp, err := json.Marshal(results)
	if err != nil {
		log.Error(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	return
}

func HandlerOpenSession(httpWriter http.ResponseWriter, httpRequest *http.Request) {
	fileName:= httpRequest.URL.Query().Get("bugReportName")

	//open bugreport in the database and get id
	hash:=postgres.OpenSession( fileName, "fakeUsername" )	
	
	jsonResp, err := json.Marshal(hash)
	if err != nil {
		log.Error(err)
		return
	}

	httpWriter.Header().Set("Content-Type", "application/json")
	httpWriter.Write(jsonResp)
	return
}

func HandlerTransfer(w http.ResponseWriter, r *http.Request) {

	// Maximum upload of 10 MB files legacy from classic
	r.ParseMultipartForm(10 << 20)

	hash:=r.FormValue("hash")
	
	if hash=="" {
		http.Error(w, "Missing hash",  http.StatusInternalServerError)
		return
	}

		
	dummy:=string(encoding.Decode(hash))

	res := strings.Split(dummy, "|")	

	id:= res[0]
	//username:=res[1] //to compare with token
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Error(err)
	}
		
	state := &postgres.State{}
	state.SetState( idInt, "TRANSFER_FILE", postgres.STARTED )

	
	//fmt.Println("dummy ",dummy, id, username )
	//TODO check if the token is from the same username
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file-upload")
	if err != nil {
		log.Error("Error Retrieving the File",err)
		return
	}
	
 	upload_size:=int(handler.Size/1024/1024)
	log.Info("Starting Upload", handler.Filename," size ", upload_size, " Mb")
	start := time.Now()
	
	defer file.Close()
	//fmt.Printf("File Size: %+v\n", handler.Size)
	//fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file

	filePath:=viper.GetString("storage.long_term_storage_path")
	storagePathLongTerm := filePath + "/original/"
	log.Info( "Check path ", storagePathLongTerm )
	if err := os.MkdirAll(filepath.Dir(storagePathLongTerm), 0770); err != nil {
		log.Error(fmt.Errorf("File: %w", err))
		return
	}	
	dstPath := path.Join(storagePathLongTerm, fmt.Sprintf( "%s:%s", id, handler.Filename))
	log.Info( "Creating file ", dstPath )
	dst, err := os.Create(dstPath)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	elapsed := time.Since(start)
	log.Info("Successfully Uploaded File", handler.Filename, " elapsed ", elapsed)
	
	postgres.FinishTransfer( idInt, handler.Filename, upload_size )
	state.SetState ( idInt, "TRANSFER_FILE", postgres.FINISHED )
	
}


func main() {
	log.Info("Sarting salsa")
	initConfig()

	if postgres.Db == nil {
		postgres.NewConnection()
	}

	handleRequests()	

}

func initConfig() {
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

