package postgres
import (
	"fmt"
	"uploader/pkg/encoding"
	log "github.com/sirupsen/logrus"	
	"time"
	//"github.com/spf13/viper"	
)
func OpenSession( fileName string, userName string ) string {
	bugreport := &Bugreport{
		FileName:         fileName,
		UserName:         userName,
		Reporter:         userName,
		UploadSize:       0,
		DateUploaded:     time.Now().Format("2006-01-02 3:4:5 pm") ,
		SerialNumber:     fileName,
		BuildFingerprint: fileName,
	}

	result := Db.Create(bugreport)
	if result.Error != nil {
		err1 := fmt.Errorf("postgres: UploadBugreport: %w", result.Error)
		log.Error(err1)
		return ""
	}

	result = Db.Delete(&bugreport)
	if result.Error != nil {
		err1 := fmt.Errorf("postgres: UploadBugreport: %w", result.Error)
		log.Error(err1)
		return ""
	}
	id := int(bugreport.ID)
	
	//mySecret:=viper.GetString("session.secret")
	hash := encoding.Encode([]byte(fmt.Sprintf("%d|%s",id,userName)))

	return hash
}


func FinishTransfer( id int, filename string, upload_size int ) {

	fmt.Println( "id ", id )
	var bugreport Bugreport
	result := Db.Unscoped().First( &bugreport, id )
	if result.Error != nil {
		err1 := fmt.Errorf("postgres: Finishtransfer bugreport: %d %w", id, result.Error)
		log.Error(err1)
		return
	}

	bugreport.FileName=filename
	bugreport.UploadSize=int64(upload_size)
	bugreport.DateUploaded=time.Now().Format("2006-01-02 3:4:5 pm")
	
	fmt.Println( "DateUploaded ", bugreport.DateUploaded, bugreport.ID )

	result = Db.Save(&bugreport)
	if result.Error != nil {
		err1 := fmt.Errorf("postgres: Finishtransfer bugreport: %d %w", id, result.Error)
		log.Error(err1)
		return
	}

	return
}



