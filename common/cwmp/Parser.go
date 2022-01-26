package cwmp

import (
	"errors"

	"github.com/ca17/teamsacs/common/xmlx"
)

// ParseXML parse xml msg
func ParseXML(data []byte) (msg Message, err error) {
	doc := xmlx.New()
	doc.LoadBytes(data, nil)
	bodyNode := doc.SelectNode("*", "Body")
	if bodyNode != nil {
		var name string
		if len(bodyNode.Children) > 1 {
			name = bodyNode.Children[1].Name.Local
		} else {
			name = bodyNode.Children[0].Name.Local
		}
		switch name {
		case "Inform":
			msg = NewInform()
		case "GetParameterValuesResponse":
			msg = &GetParameterValuesResponse{}
		case "SetParameterValuesResponse":
			msg = &SetParameterValuesResponse{}
		case "GetParameterNames":
			msg = &GetParameterNames{}
		case "GetParameterNamesResponse":
			msg = &GetParameterNamesResponse{}
		case "DownloadResponse":
			msg = &DownloadResponse{}
		case "TransferComplete":
			msg = &TransferComplete{}
		case "GetRPCMethodsResponse":
			msg = &GetRPCMethodsResponse{}
		case "RebootResponse":
			msg = &RebootResponse{}
		default:
			return nil, errors.New("no msg type match")
		}
		if msg != nil {
			msg.Parse(doc)
		}
	}
	return msg, err
}
