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

package message

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models/elastic"
)

// DNS LOG
type DnsLog struct {
	Start    time.Time
	End      time.Time
	Time     int64
	Src      string
	Question []DnsQuestion
	Ecs      string
	Tags     []string
	Result   []DnsResult
	Error    string
}

type DnsQuestion struct {
	Name   string
	Qtype  string
	Qclass string
}

type DnsResult struct {
	Type  string
	Name  string
	Value string
	Ttl   uint32
}

var EmptyDnslog = DnsLog{}

// Remove the log from the channel and put it into the batch processing queue
func (t *PubSubService) processDnslog(ch chan DnsLog) {
	readWithSelect := func(ch chan DnsLog) (log DnsLog, err error) {
		timeout := time.NewTimer(time.Millisecond * 1000)
		select {
		case log = <-ch:
			return log, nil
		case <-timeout.C:
			return EmptyDnslog, errors.New("read time out")
		}
	}

	for {
		queue := make([]elastic.TeamsDnsLog, 0)
		var err error
		var _log DnsLog
		for {
			if len(queue) >= 2000 {
				break
			}
			_log, err = readWithSelect(ch)
			if err != nil {
				break
			}
			tlog := elastic.TeamsDnsLog{
				Timestamp: _log.Start.Format(time.RFC3339),
				Time:      _log.Time,
				CpeName:   "N/A",
				CpeSn:     "N/A",
				Src:       _log.Src,
				Question:  make([]map[string]interface{}, 0),
				Ecs:       _log.Ecs,
				Tags:      _log.Tags,
				Result:    make([]map[string]interface{}, 0),
			}

			// ECS IP Geographical processing
			if tlog.Ecs != "" && t.Manager.Ipdb != nil {
				ipinfo, err := t.Manager.Ipdb.FindInfo(tlog.Ecs, "CN")
				if err == nil {
					tlog.EcsCity = ipinfo.CityName
					tlog.EcsCountry = ipinfo.CountryName
				}
			}

			// DNS Request data
			for _, qe := range _log.Question {
				_qmap := map[string]interface{}{
					"name":   qe.Name,
					"qtype":  qe.Qtype,
					"qclass": qe.Qclass,
				}
				tlog.Question = append(tlog.Question, _qmap)
			}

			// DNS Response data
			for _, result := range _log.Result {
				_rmap := map[string]interface{}{
					"type":  result.Type,
					"name":  result.Name,
					"value": result.Value,
					"ttl":   result.Ttl,
				}

				// A Recode Geographical processing
				if (result.Type == "A" || result.Type == "AAAA") && t.Manager.Ipdb != nil {
					ipinfo, err := t.Manager.Ipdb.FindInfo(result.Value, "CN")
					if err == nil {
						_rmap["city"] = ipinfo.CityName
						_rmap["country"] = ipinfo.CountryName
					}
				}

				tlog.Result = append(tlog.Result, _rmap)
			}

			queue = append(queue, tlog)
			log.Infof("%+v", tlog)
		}
		if len(queue) > 0 {
			_queue := make([]elastic.TeamsDnsLog, len(queue))
			copy(_queue, queue)
			go func() {
				_, err = t.Manager.Elastic.BulkTeamsDnslog(_queue...)
				if err != nil {
					log.Error(err.Error())
				}
			}()
		}
	}
}

// _appendDnslog Append the received subscription log to the channel
// If the channel is full, the message will be discarded to avoid blocking
func _appendDnslog(ch chan DnsLog, log DnsLog) error {
	timeout := time.NewTimer(time.Microsecond * 500)
	select {
	case ch <- log:
		return nil
	case <-timeout.C:
		return errors.New("channel full, write dnslog time out")
	}
}

// StartDnslogSubServer Start DNS log message subscription service
func (t *PubSubService) StartDnslogSubServer() error {
	var sock mangos.Socket
	var err error
	var topicbyte = []byte(TeamsDNSLOGTopic)
	if sock, err = sub.NewSocket(); err != nil {
		return err
	}
	if err = sock.Listen(t.Manager.Config.Messaged.SubAddress); err != nil {
		log.Errorf("%+v", errors.WithStack(err))
		return err
	}

	if err = sock.SetOption(mangos.OptionSubscribe, topicbyte); err != nil {
		return fmt.Errorf("dnslog subscriber sub error %s", err.Error())
	}

	t.SubSock = sock
	log.Infof("dnslog SubServer listen %s...", t.Manager.Config.Messaged.SubAddress)

	chl := make(chan DnsLog, 10000)

	go t.processDnslog(chl)

	for {
		msg, err := sock.Recv()
		if err != nil {
			log.Errorf("Dnslog Subscriber recv Message error %s", err.Error())
			continue
		}
		var smsg = new(DnsLog)
		err = msgpack.Unmarshal(msg[len(topicbyte):], smsg)
		if err != nil {
			log.Errorf("Dnslog Subscriber Unmarshal dnslog Message error %s", err.Error())
			continue
		}

		if err := _appendDnslog(chl, *smsg); err != nil {
			log.Errorf("Dnslog Subscriber process dnslog Message error %s", err.Error())
		}
	}
}
