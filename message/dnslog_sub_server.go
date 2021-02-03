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
	Start   time.Time
	End     time.Time
	Time    int64
	Src     string
	Dest    []string
	Matcher []string
	Ecs     string
	Tags    []string
	Result  []DnsResult
	Error   string
}

type DnsResult struct {
	Type  string
	Name  string
	Value string
	Ttl   uint32
}

var EmptyDnslog = DnsLog{}


// Remove the log from the channel and put it into the batch processing queue
func (t *PubSubService)  processDnslog(ch chan DnsLog) {
	readWithSelect := func(ch chan DnsLog) (log DnsLog, err error) {
		timeout := time.NewTimer(time.Millisecond * 100)
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
		if err == nil  && len(queue) < 2000 {
			_log, err = readWithSelect(ch)
			tlog := elastic.TeamsDnsLog{
				Timestamp: _log.Start.Format(time.RFC3339),
				Time:      _log.Time,
				CpeName:   "N/A",
				CpeSn:     "N/A",
				Src:       _log.Src,
				Dest:      _log.Dest,
				Ecs:       _log.Ecs,
				Tags:      _log.Tags,
				Result: make([]map[string]interface{}, 0),
				Error:     _log.Error,
			}
			for _, result := range _log.Result {
				tlog.Result = append(tlog.Result, map[string]interface{}{
					"type": result.Type,
					"name": result.Name,
					"value": result.Value,
					"ttl": result.Ttl,
				})
			}
			queue = append(queue, tlog)
		}
		go func() {
			_, err = t.Manager.Elastic.BulkTeamsDnslog(queue...)
			if err != nil {
				log.Error(err.Error())
			}
		}()
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
			if err != nil {
				log.Errorf("Dnslog Subscriber recv Message error %s", err.Error())
				continue
			}
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

	return nil
}


