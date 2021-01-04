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

package elk

import (
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)


// GetClient
// urls: elasticsearch service address, with multiple service addresses separated by commas
// // account and password based on http base auth authentication mechanism
func NewElasticClient(user, pwd string, urls ...string) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(urls...),
		elastic.SetBasicAuth(user, pwd),
		elastic.SetGzip(true),
		elastic.SetHealthcheckInterval(10*time.Second),
		// elastic.SetMaxRetries(5),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewConstantBackoff(5 * time.Second))),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		return nil, err
	}
	return client, nil
}


// GetClient
// urls: elasticsearch service address, with multiple service addresses separated by commas
// // account and password based on http base auth authentication mechanism
func NewSimpleElasticClient(user, pwd string, urls ...string) (*elastic.Client, error) {
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(urls...),
		elastic.SetBasicAuth(user, pwd),
		elastic.SetGzip(true),
		// elastic.SetHealthcheckInterval(10*time.Second),
		// elastic.SetMaxRetries(5),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewConstantBackoff(5 * time.Second))),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		return nil, err
	}
	return client, nil
}


