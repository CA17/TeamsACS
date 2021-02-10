欢迎来到 TeamsACS 项目!

      _____                                    ___     ___     ___   
     |_   _|   ___    __ _    _ __     ___    /   \   / __|   / __|  
       | |    / -_)  / _` |  | '  \   (_-<    | - |  | (__    \__ \  
      _|_|_   \___|  \__,_|  |_|_|_|  /__/_   |_|_|   \___|   |___/  
    _|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""|_|"""""| 
    "`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-'"`-0-0-' 

[English](README.md)

# TeamsACS

TeamsACS 致力于为工作团队提供卓越的网络管理便捷性。我们以Mikrotik的网络产品为基础核心， 同时将系统能力扩展到更广泛的网络设备， 比如 OpenWrt。

系统核心基于 Golang 技术开发， 提供了卓越的性能和便捷的部署体验。

## 系统架构

![image](https://user-images.githubusercontent.com/377938/97301570-e28b3d80-1892-11eb-85a8-5cc5f80449a4.png)

## 系统功能特性

### TR069 ACS 集成

首选 GenieACS 开源ACS系统集成， GenieACS可以与任何支持TR-069协议的设备一起工作。它能自动发现设备的参数树（包括特定于厂商的参数），不对设备的数据模型做任何预定义。

它已经与许多不同制造商的各种设备（DSL、电缆、光纤、LTE CPE、VoIP电话）进行了测试。

同时也是 Mikrotik 官方推荐的系统，经过了大量的测试， 可以很放心的使用。

### 北向接口 Northbound Interface

- 面向各类第三方管理系统提供统一的API， 基于 HTTPS Json 协议。
- 提供基础设备信息以及状态数据的查询 API， 以及数据维护 API。
- 提供各种策略管理 API， 比如防火墙规则，路由表等。

## 快速开始

### 使用 Docker-compose 部署

部署描述文件 docker-compose.yml

具体操作请参考 [docker-compose](https://docs.docker.com/compose/) 文档

##  链接

- [TeamsACS WIKI 文档](https://github.com/CA17/TeamsACS/wiki)<!-- @IGNORE PREVIOUS: link -->
- [TLS encryption for RADIUS over TCP (RadSec) 00)](https://tools.ietf.org/id/draft-ietf-radext-radsec-00.html)
- [TLS encryption for RADIUS over TCP (RadSec) 05](https://tools.ietf.org/html/draft-ietf-radext-radsec-05)
- [mikrotik RADIUS Client](https://wiki.mikrotik.com/wiki/Manual:RADIUS_Client)
- [Manual:TR069-client](https://wiki.mikrotik.com/wiki/Manual:TR069-client)

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