package postgres
import (
	"fmt"
	log "github.com/sirupsen/logrus"	
	"time"
	//"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

/*****************************************************************************/
/*  States                                                                   */ 
/*****************************************************************************/
type Detail struct {
	TaskName string `yaml:"TaskName" json:"TaskName"`
	Status string `yaml:"Status" json:"Status"`
	TimeStart *time.Time `yaml:"TimeStart" json:"TimeStart"`
	TimeEnd   *time.Time `yaml:"TimeEnd" json:"TimeEnd"`
}

type Task struct {
	Sequence  int `yaml:"Sequence"`
	CheckDetail []Detail `yaml:"CheckDetail"`
}

type Step struct {
	StepName string `yaml:"StepName"`
	Tasks    []Task `yaml:"Tasks"`
}

type Status []Step

type State struct {
	gorm.Model
	BugreportID       int	
	//Step string
	//Sequence string
	TaskName string
	State string
	TimeStart *time.Time
	TimeEnd *time.Time
}

type CheckPoints struct {
	TaskName string
	State string
	TimeStart string
	TimeEnd string
}

type ProgressStatusStruct struct {
	StatusTemplate Status
	CheckPoints []CheckPoints
}

type ProgressStatus interface {
	InitStatus( int )
	SetStatus( int, string, string )
	GetCheckPoints()
	LoadCheckPoints( int ) []string
}

const (
   STARTED = "STARTED"
   FINISHED = "FINISHED"
)

//load template
func ( c ProgressStatusStruct ) InitStatus( id int ) {

	var bugreport Bugreport
	result := Db.Unscoped().First( &bugreport, id )
	if result.Error != nil {
		err1 := fmt.Errorf("postgres: CreateStatus bugreport: %d %w", id, result.Error)
		log.Error(err1)
		return
	}
	
	yamlFile, err := ioutil.ReadFile( fmt.Sprintf("./configs/%s",bugreport.ConfigFile) )
	if err != nil {
        	log.Fatalf("yamlFile.Get err   #%v ", err)
	}
	
	var yamlContent []Step
	err = yaml.Unmarshal(yamlFile, &yamlContent )
 	if err != nil {
        	log.Fatalf("Unmarshal: %v", err)
	}
	c.StatusTemplate= yamlContent
}

//return list of tasks of the template
func ( c ProgressStatusStruct ) GetCheckPoints() []string {

var tasks []string
for _,s := range c.StatusTemplate {
	for _,t := range s.Tasks {
		for _,td := range t.CheckDetail {
			tasks = append( tasks, td.TaskName )
		}
	}
}
return tasks
}

//load list of tasks in the status template
func ( c ProgressStatusStruct ) LoadCheckPoints( id int )  {

var states []State
result := Db.Where( " bugreport_id = ? ", id ).Find( states )
if result.Error != nil {
	err1 := fmt.Errorf("postgres: CreateStatus bugreport: %d %w", id, result.Error)
	log.Error(err1)
	return
}

for _,s := range c.StatusTemplate {
	for _,t := range s.Tasks {
		for _,td := range t.CheckDetail {
			for _,v := range states {
				if td.TaskName == v.TaskName {
					d:=CheckPoints{}
					d.TaskName=v.TaskName
					d.TimeStart=v.TimeStart.Format("2006-01-02T15:04:05Z")
					d.TimeEnd=v.TimeEnd.Format("2006-01-02T15:04:05Z")
					c.CheckPoints = append( c.CheckPoints, d )
				}
			}
		}
	}
}

return
}
/*****************************************************************************/
//get all states for a bugreport without template, time formated 
func GetStates( id int ) []CheckPoints {
	var results []CheckPoints
	if err := Db.Table("states").Select("task_name,state,TO_CHAR(time_start,'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') time_start,TO_CHAR(time_end,'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') time_end").Where("bugreport_id=?", id).Order("id").Find(&results).Error; err != nil {
		log.Error(err)
		return nil
	}
return results
}


func ( s *State ) GetState( bugReportId int, taskName string ) {

	result := Db.Where(" bugreport_id = ? and task_name= ? ", bugReportId, taskName ).First( s )
	fmt.Println( "recupero0 ",*s)	
	if result.Error != nil {
		err1 := fmt.Errorf("postgres: GetState bugreport: %d %w", bugReportId, result.Error)
		log.Error(err1)
		return
	}
}

//set status and update database
func ( s *State ) SetState ( bugReportId int, taskName string, state string ) {

	s.GetState( bugReportId, taskName )
	
	fmt.Println( "recupero ",*s)
	
	t:=time.Now()
	
	s.BugreportID = bugReportId
	s.TaskName = taskName
	s.State = state
	if state == "STARTED" {
		s.TimeStart = &t
	} else if state == "FINISHED" {
		s.TimeEnd = &t
	}

	result := Db.Save(s)
	if result.Error != nil {
		err1 := fmt.Errorf("postgres: SetState bugreport: %d %w", bugReportId, result.Error)
		log.Error(err1)
		return
	}
	
}

