package cwmp

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/ca17/teamsacs/common/xmlx"
)

// FactoryResetResponse
type FactoryResetResponse struct {
	ID   string
	Name string
}

type factoryResetResponseBodyStruct struct {
	Body factoryResetResponseStruct `xml:"cwmp:FactoryResetResponse"`
}

type factoryResetResponseStruct struct {
}

// GetID get msg id
func (msg *FactoryResetResponse) GetID() string {
	if len(msg.ID) < 1 {
		msg.ID = fmt.Sprintf("ID:intrnl.unset.id.%s%d.%d", msg.GetName(), time.Now().Unix(), time.Now().UnixNano())
	}
	return msg.ID
}

// GetName get msg type
func (msg *FactoryResetResponse) GetName() string {
	return "FactoryResetResponse"
}

// CreateXML encode into xml
func (msg *FactoryResetResponse) CreateXML() []byte {
	env := Envelope{}
	env.XmlnsEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	env.XmlnsEnc = "http://schemas.xmlsoap.org/soap/encoding/"
	env.XmlnsXsd = "http://www.w3.org/2001/XMLSchema"
	env.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	env.XmlnsCwmp = "urn:dslforum-org:cwmp-1-0"
	id := IDStruct{Attr: "1", Value: msg.GetID()}
	env.Header = HeaderStruct{ID: id}
	body := factoryResetResponseStruct{}
	env.Body = factoryResetResponseBodyStruct{body}
	// output, err := xml.Marshal(env)
	output, err := xml.MarshalIndent(env, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output
}

// Parse decode from xml
func (msg *FactoryResetResponse) Parse(doc *xmlx.Document) {
	msg.ID = getDocNodeValue(doc, "*", "ID")
}
