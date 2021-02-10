Welcome to the TeamsACS Project!

      _____                                    ___     ___     ___   
     |_   _|   ___    __ _    _ __     ___    /   \   / __|   / __|  
       | |    / -_)  / _` |  | '  \   (_-<    | - |  | (__    \__ \  
      _|_|_   \___|  \__,_|  |_|_|_|  /__/_   |_|_|   \___|   |___/  
    _|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""| 
    "`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-' 

[中文](README_CN.md)

# TeamsACS

TeamsACS is committed to providing work teams with exceptional ease of network management. We use Mikrotik's network products as the core foundation, while extending the system's capabilities to a wider range of network devices, such as OpenWrt.

The core of the system is based on Golang technology, providing excellent performance and ease of deployment.

## Systems Architecture

![image](https://user-images.githubusercontent.com/377938/97301570-e28b3d80-1892-11eb-85a8-5cc5f80449a4.png)

## System Features

### TR069 ACS integration

Preferring GenieACS open source ACS system integration， GenieACS can work with any device that supports the TR-069 protocol.

It auto-discovers the device’s parameter tree (including vendor-specific parameters) making no assumptions about the device’s data model.

It’s been tested with a wide range of devices (DSL, cable, fiber optics, LTE CPEs, VoIP phones) from many different manufacturers.

It is also the officially recommended system by Mikrotik, which has been tested extensively and is safe to use.

### Northbound Interface

- Provide a unified API for various third-party management systems, based on HTTPS Json protocol.
- Provide the query API for basic device information and status data, and data maintenance API.
- Provide a variety of policy management APIs, such as firewall rules, routing tables and so on.


## Quick start

### Deploy using Docker-compose

Deployment profile: docker-compose.yml

For specific operations, please refer to [docker-compose](https://docs.docker.com/compose/)

##  Links

- [TeamsACS WIKI Documentation](https://github.com/CA17/TeamsACS/wiki)<!-- @IGNORE PREVIOUS: link -->
- [TLS encryption for RADIUS over TCP (RadSec) 00)](https://tools.ietf.org/id/draft-ietf-radext-radsec-00.html)
- [TLS encryption for RADIUS over TCP (RadSec) 05](https://tools.ietf.org/html/draft-ietf-radext-radsec-05)
- [mikrotik RADIUS Client](https://wiki.mikrotik.com/wiki/Manual:RADIUS_Client)

##  License

    Licensed to the Apache Software Foundation (ASF) under one or more
    contributor license agreements.  See the NOTICE file distributed with
    this work for additional information regarding copyright ownership.
    The ASF licenses this file to You under the Apache License, Version 2.0
    (the "License"); you may not use this file except in compliance with
    the License.  You may obtain a copy of the License at
        http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.