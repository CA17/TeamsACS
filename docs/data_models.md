# Data model

## VPE

VPE (VPNPE) is a special type of PE that is connected to the CPE not by traditional leased line technologies such as DDN/E1/POS/ETH/PVC,  
but by tunneling technologies such as IPSec/L2TP/GRE/UDPVPN.

The following is the data structure of vpe in the teamsacs system，  Some attributes are collected from genieacs, and some attributes are submitted through the form

```
{
    "device_id": "E48D8C-CHR-VBxa82233", // genieacs device id
    "name": "vpe01-chr", // device name
    "sn": "VBxa82233", // device sn
    "manufacturer": "MikroTik", // genieacs attr
    "oui": "E48D8C", // genieacs attr
    "model": "CHR", // genieacs attr
    "sversion": "6.47", device software version
    "version": "v1.0", // device hardware version
    "product_class": "CHR", // genieacs attr
    "memuse": 27, // device memary use from genieacs
    "cpuuse": 10, // device cpu use from genieacs
    "uptime": 15057052, // device uptime from genieacs
    "online_status": "on", //device online status
    "status": "enabled", // device status
    "identifier": "vpe01-chr", // device identfier, radiu attr
    "ipaddr": "10.10.10.1", // device ipaddr
    "secret": "mysecret", // radius secret
    "coa_port": "3799", // radius coa port
    "api_addr": "",
    "api_user": "admin",
    "api_pwd": "xxxxxx",
    "tags": "vpe",
    "remark": "vpe",
    "last_inform": "2021-03-03T13:59:01.232Z", // genieacs last inform time
    "update_time": "2021-03-03 13:59:01 Z UTC" // data update time
}
```

## Subscribe

Subscribe is the RADIUS user model

```
{
    "username": "account01", // radius username
    "email": "myacsaccount@gmail.com",
    "password": "123456", // radius password
    "active_num": 3, // max radius online count
    "addr_pool": "J8009", // radius addr pool
    "up_rate": 20000, // radius Uplink rate bps
    "down_rate": 20000, // radius Downlink rate bps
    "ip_addr": "", // radius user ipaddr
    "expire_time": "2022-12-31 00:00:00", 
    "status": "enabled", // user status
    "remark": "demo user"
}
```