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

package snmp

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/hallidave/mibtool/smi"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/snmp/mibs/ifmib"
)

func TestWalk(t *testing.T) {

	mib := smi.NewMIB("/usr/share/snmp/mibs")
	mib.Debug = true
	err := mib.LoadModules("IF-MIB")
	if err != nil {
		log.Fatal(err)
	}

	var community = "publicSD2"
	target := "192.168.100.1"

	gosnmp.Default.Target = target
	gosnmp.Default.Community = community
	gosnmp.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking
	err = gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}
	defer gosnmp.Default.Conn.Close()

	rs, err := gosnmp.Default.BulkWalkAll(ifmib.IF_MIB_ifEntry_OID)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
	for _, pdu := range rs {
		fmt.Printf("name:%s = ", pdu.Name)

		switch pdu.Type {
		case gosnmp.OctetString:
			b := pdu.Value.([]byte)
			fmt.Printf("STRING: %s\n", string(b))
		default:
			fmt.Printf(" TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
		}
	}
}

func TestMibParse(t *testing.T) {
	mib := smi.NewMIB("/usr/share/snmp/mibs")
	mib.Debug = true
	err := mib.LoadModules("IF-MIB")
	if err != nil {
		log.Fatal(err)
	}

	// Walk all symbols in MIB
	mib.VisitSymbols(func(sym *smi.Symbol, oid smi.OID) {
		_s := strings.ReplaceAll(sym.String(), "-", "_")
		_s = strings.ReplaceAll(_s, "::", "_")
		fmt.Printf("%-45s = \"%s\"\n", _s+"_OID", oid)
		fmt.Printf("%-45s = \"%s\"\n", _s+"_NAME", sym.String())
	})
}

func TestXXX(t *testing.T) {
	oid := ".1.3.6.1.2.1.2.2.1.2.13"
	fmt.Println(oid[strings.LastIndex(oid, ".")+1:])
}

func TestSnmpV2Client_QueryInterfaces(t *testing.T) {
	sc := NewSnmpV2Client("192.168.100.1", 161, "publicSD2")
	err := sc.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()
	ifcs, err := sc.QueryInterfaces()
	if err != nil {
		t.Fatal(err)
	}
	sc.QueryInterfacesInOctets(ifcs)
	sc.QueryInterfacesOutOctets(ifcs)
	t.Log(ifcs)

}

func Benchmark(b *testing.B) {
	sc := NewSnmpV2Client("192.168.100.1", 161, "publicSD2")
	err := sc.Connect()
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()
	for i := 20; i < b.N; i++ {
		ifcs, err := sc.QueryInterfaces()
		sc.QueryInterfacesInOctets(ifcs)
		sc.QueryInterfacesOutOctets(ifcs)
		if err != nil {
			b.Fatal(err)
		}
	}
}
