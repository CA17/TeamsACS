package models

var Tables = []interface{}{
	&SysConfig{},
	&Syslog{},
	&NetRegion{},
	&NetEdge{},
	&NetEdgeEventLog{},
	&NetDevice{},
	&NetDeviceParam{},
	&NetDeviceEventLog{},
	&CwmpVendorConfig{},
	&CwmpVendorConfigHistory{},
	&CwmpDownloadTask{},
	&CwmpSetParamTask{},
	&CwmpParamNameList{},
	&NetLocalScript{},
	&NetIpaddress{},
	&TsDeviceLoad{},
}
