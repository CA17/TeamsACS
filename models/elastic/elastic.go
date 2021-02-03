/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package elastic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
)

const IndexTeamslogPrefix = "teamsacs-log"
const IndexTeamsDnslogPrefix = "teamsdns-log"

type Elastic struct {
	Client *elastic.Client
}

func NewElastic(client *elastic.Client) *Elastic {
	e := &Elastic{Client: client}
	go e.InitTemplate()
	go e.InitDnslogTemplate()
	return e
}

const _defaultMapping = `{
  "mappings": {
      "date_detection": false
  }
}`

const _mapping = `
{
  "order": 0,
  "index_patterns": [
    "teamsacs-*"
  ],
  "settings": {
    "number_of_shards": "3",
    "number_of_replicas": "0",
    "auto_expand_replicas": "0-1",
    "codec": "best_compression",
    "index.refresh_interval": "5s"
  },
  "mappings": {
    "_default_": {
      "date_detection": false
    },
    "properties": {
      "@timestamp": {
        "type": "date"
      },
      "source": {
        "type": "keyword"
      },
      "sn": {
        "type": "keyword"
      },
      "name": {
        "type": "keyword"
      },
      "tags": {
        "type": "keyword"
      },
      "model": {
        "type": "keyword"
      },
      "version": {
        "type": "keyword"
      },
      "devtype": {
        "type": "keyword"
      },
      "sysstat": {
        "properties": {
		  "stattime": {
	        "type": "date"
	      },
          "memPercent": {
            "type": "integer"
          },
          "cpuPercent": {
            "type": "integer"
          },
          "upTime": {
            "type": "long"
          }
        }
      },
      "netstat": {
        "properties": {
          "interface": {
            "type": "keyword"
          },
		  "stattime": {
	        "type": "date"
	      },
          "sendBytes": {
            "type": "long"
          },
          "recvBytes": {
            "type": "long"
          },
          "sendPackets": {
            "type": "long"
          },
          "recvPackets": {
            "type": "long"
          }
        }
      },
      "radiuslog": {
        "properties": {
          "username": {
            "type": "keyword"
          },
          "acctSessionId": {
            "type": "keyword"
          },
          "nasId": {
            "type": "keyword"
          },
          "nasAddr": {
            "type": "keyword"
          },
          "framedIpaddr": {
            "type": "keyword"
          },
          "framedNetmask": {
            "type": "keyword"
          },
          "macAddr": {
            "type": "keyword"
          },
          "nasPort": {
            "type": "keyword"
          },
          "nasClass": {
            "type": "keyword"
          },
          "nasPortId": {
            "type": "keyword"
          },
          "nasPortType": {
            "type": "keyword"
          },
          "serviceType": {
            "type": "keyword"
          },
          "acctSessionTime": {
            "type": "long"
          },
          "acctInputTotal": {
            "type": "long"
          },
          "acctOutputTotal": {
            "type": "long"
          },
          "acctInputPackets": {
            "type": "long"
          },
          "acctOutputPackets": {
            "type": "long"
          },
          "sessionTimeout": {
            "type": "long"
          },
          "acctStartTime": {
            "type": "date"
          },
          "acctStopTime": {
            "type": "date"
          }
        }
      }
    }
  }
}`

const _dnslogmapping = `
{
  "order": 0,
  "index_patterns": [
    "teamsdns-log-*"
  ],
  "settings": {
    "number_of_shards": "3",
    "number_of_replicas": "0",
    "auto_expand_replicas": "0-1",
    "codec": "best_compression",
    "index.refresh_interval": "5s"
  },
  "mappings": {
    "_default_": {
      "date_detection": false
    },
    "properties": {
      "@timestamp": {
        "type": "date"
      },
      "time": {
        "type": "long"
      },
      "cpe_name": {
        "type": "keyword"
      },
      "cpe_sn": {
        "type": "keyword"
      },
      "Ecs": {
        "type": "keyword"
      },
      "src": {
        "type": "keyword"
      },
      "dest": {
        "type": "keyword"
      },
      "result": {
        "type": "nested"
      },
      "tags": {
        "type": "keyword"
      },
      "error": {
        "type": "keyword"
      }
    }
  }
}`

type DeviceSysstat struct {
	UpTime     int64 `json:"upTime"`
	MemPercent int64 `json:"memPercent"`
	CpuPercent int64 `json:"cpuPercent"`
}

type DeviceNetstat struct {
	Interface   string `json:"interface"`
	SendBytes   int64  `json:"sendBytes"`
	RecvBytes   int64  `json:"recvBytes"`
	SendPackets int64  `json:"sendPackets"`
	RecvPackets int64  `json:"recvPackets"`
}

type Radiuslog struct {
	Username          string `json:"username"`
	AcctSessionId     string `json:"acctSessionId"`
	NasId             string `json:"nasId"`
	NasAddr           string `json:"nasAddr"`
	FramedIpaddr      string `json:"framedIpaddr"`
	FramedNetmask     string `json:"framedNetmask"`
	MacAddr           string `json:"macAddr"`
	NasPort           string `json:"nasPort"`
	NasClass          string `json:"nasClass"`
	NasPortId         string `json:"nasPortId"`
	NasPortType       string `json:"nasPortType"`
	ServiceType       string `json:"serviceType"`
	AcctSessionTime   string `json:"acctSessionTime"`
	AcctInputTotal    int64  `json:"acctInputTotal"`
	AcctOutputTotal   int64  `json:"acctOutputTotal"`
	AcctInputPackets  int64  `json:"acctInputPackets"`
	AcctOutputPackets int64  `json:"acctOutputPackets"`
	SessionTimeout    int64  `json:"sessionTimeout"`
	AcctStartTime     string `json:"acctStartTime"`
	AcctStopTime      string `json:"acctStopTime"`
}

type TeamsacsLog struct {
	Timestamp string         `json:"@timestamp"`
	Source    string         `json:"source"`
	Sn        string         `json:"sn"`
	Name      string         `json:"name,omitempty"`
	Tags      string         `json:"tags,omitempty"`
	Model     string         `json:"model,omitempty"`
	Version   string         `json:"version,omitempty"` // ros ver
	Devtype   string         `json:"devtype,omitempty"` // cpe | vpe
	Sysstat   *DeviceSysstat `json:"sysstat,omitempty"`
	Netstat   *DeviceNetstat `json:"netstat,omitempty"`
	Radiuslog *Radiuslog     `json:"radiuslog,omitempty"`
}

type TeamsDnsLog struct {
	Timestamp string   `json:"@timestamp"`
	Time      int64    `json:"time"`
	CpeName string `json:"cpe_name,omitempty"`
	CpeSn string `json:"cpe_sn,omitempty"`
	Src       string   `json:"src"`
	Dest      []string `json:"dest"`
	Ecs       string   `json:"ecs,omitempty"`
	Tags      []string `json:"tags"`
	Result    []map[string]interface{} `json:"result"`
	Error string `json:"error,omitempty"`
}


func GetCurrentTeamslogIndexName() string {
	suffix := time.Now().Format("20060102")
	indexName := fmt.Sprintf("%s-%s", IndexTeamslogPrefix, suffix)
	return indexName
}

func GetCurrentTeamsDnslogIndexName() string {
	suffix := time.Now().Format("20060102")
	indexName := fmt.Sprintf("%s-%s", IndexTeamsDnslogPrefix, suffix)
	return indexName
}

// func GetSearchIndexOfMonth() string {
// 	suffix := time.Now().Format("200601")
// 	indexName := fmt.Sprintf("%s-%s", IndexTeamslogPrefix, suffix)
// 	return indexName
// }

func (e *Elastic) checkClient() error {
	switch {
	case e.Client == nil:
		return errors.New("elasticsearch Client is nil")
	}
	return nil
}

// InitTemplate
// Initialize the global template and call it only at system startup
func (e *Elastic) InitTemplate() error {
	if err := e.checkClient(); err != nil {
		return err
	}
	ctx := context.Background()
	_, err := e.Client.IndexPutTemplate("teamsacs-template").BodyJson(_mapping).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *Elastic) InitDnslogTemplate() error {
	if err := e.checkClient(); err != nil {
		return err
	}
	ctx := context.Background()
	_, err := e.Client.IndexPutTemplate("teamsdns-log-template").BodyJson(_dnslogmapping).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// BulkTeamslog
func (e *Elastic) BulkTeamslog(logs ...TeamsacsLog) (*elastic.BulkResponse, error) {
	if err := e.checkClient(); err != nil {
		return nil, err
	}
	bulkRequest := e.Client.Bulk()
	for _, log := range logs {
		req := elastic.NewBulkIndexRequest().Index(GetCurrentTeamslogIndexName()).Id(common.UUID()).Doc(log)
		bulkRequest = bulkRequest.Add(req)
	}
	ctx := context.Background()
	bulkResponse, err := bulkRequest.Do(ctx)
	if err != nil {
		return nil, err
	}
	return bulkResponse, nil
}

// BulkTeamsdnsLog
// Batch send dnslog
func (e *Elastic) BulkTeamsDnslog(logs ...TeamsDnsLog) (*elastic.BulkResponse, error) {
	if err := e.checkClient(); err != nil {
		return nil, err
	}
	bulkRequest := e.Client.Bulk()
	for _, _log := range logs {
		req := elastic.NewBulkIndexRequest().Index(GetCurrentTeamsDnslogIndexName()).Id(common.UUID()).Doc(_log)
		bulkRequest = bulkRequest.Add(req)
	}
	ctx := context.Background()
	bulkResponse, err := bulkRequest.Do(ctx)
	if err != nil {
		return nil, err
	}
	return bulkResponse, nil
}


// BulkData
// sync base data
func (e *Elastic) BulkData(indexName string, data []map[string]interface{}, deleteIndex bool) (*elastic.BulkResponse, error) {
	if err := e.checkClient(); err != nil {
		return nil, err
	}
	if deleteIndex {
		// first delete exists index
		_, err := e.Client.DeleteIndex(indexName).Do(context.Background())
		if err != nil {
			log.Error(err)
		}

		_, err = e.Client.CreateIndex(indexName).BodyString(_defaultMapping).Do(context.Background())
		if err != nil {
			log.Error(err)
		}
	}

	bulkRequest := e.Client.Bulk()
	for _, item := range data {
		_id, ok := item["id"]
		if !ok {
			_id = common.UUID()
		}
		delete(item, "_id")
		item["doc_type"] = indexName
		req := elastic.NewBulkIndexRequest().Index(indexName).Id(_id.(string)).Doc(item)
		bulkRequest = bulkRequest.Add(req)
	}
	ctx := context.Background()
	bulkResponse, err := bulkRequest.Do(ctx)
	if err != nil {
		return nil, err
	}
	return bulkResponse, nil
}

// AddData
func (e *Elastic) AddData(indexName string, data map[string]interface{}) error {
	if err := e.checkClient(); err != nil {
		return err
	}
	_id, ok := data["id"]
	if !ok {
		_id = common.UUID()
	}
	delete(data, "_id")
	data["doc_type"] = indexName
	_, err := e.Client.Index().Index(indexName).Id(_id.(string)).BodyJson(data).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// UpdateData
func (e *Elastic) UpdateData(indexName string, data map[string]interface{}) error {
	if err := e.checkClient(); err != nil {
		return err
	}
	_id, ok := data["id"]
	if !ok {
		return errors.New("data _id is empty")
	}
	delete(data, "_id")
	data["doc_type"] = indexName
	_, err := e.Client.Update().Index(indexName).Id(_id.(string)).Doc(data).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// DeleteData
func (e *Elastic) DeleteData(indexName string, _id string) error {
	if err := e.checkClient(); err != nil {
		return err
	}
	_, err := e.Client.Delete().Index(indexName).Id(_id).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
