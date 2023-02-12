package cwmp

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/ca17/teamsacs/common/xmlx"
)

// ScheduleInform  cpe
type ScheduleInform struct {
	ID           string
	Name         string
	NoMore       int
	CommandKey   string
	DelaySeconds int
}

type scheduleInformBodyStruct struct {
	Body scheduleInformStruct `xml:"cwmp:ScheduleInform"`
}

type scheduleInformStruct struct {
	CommandKey   string
	DelaySeconds int
}

// GetID get msg id
func (msg *ScheduleInform) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

// GetName get msg name
func (msg *ScheduleInform) GetName() string {
	return "ScheduleInform"
}

// CreateXML encode into xml
func (msg *ScheduleInform) CreateXML() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IDStruct{Attr: "1", Value: msg.GetID()}
	env.Header = HeaderStruct{ID: id, NoMore: msg.NoMore}
	body := scheduleInformStruct{CommandKey: msg.CommandKey, DelaySeconds: msg.DelaySeconds}
	env.Body = scheduleInformBodyStruct{body}
	// output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

// Parse decode from xml
func (msg *ScheduleInform) Parse(doc *xmlx.Document) {
	// TODO
}
