# Radius user management

## Create a new user

Request

```
POST http://{{nbi_url}}/nbi/data/subscribe/add
Content-Type: application/json
authorization: Bearer {{nbi_token}}

{
  "username": "account01",
  "email": "myacsaccount@gmail.com",
  "password": "123456",
  "active_num": 3,
  "addr_pool": "J8009",
  "up_rate": 20000,
  "down_rate": 20000,
  "ip_addr": "",
  "expire_time": "2022-12-31 00:00:00",
  "status": "enabled",
  "remark": "demo user"
}
```

Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Wed, 03 Mar 2021 15:00:11 GMT
Content-Length: 73

{
  "code": 0,
  "msgtype": "info",
  "msg": "Success",
  "data": null
}
```

## Update user

Request

```
POST http://{{nbi_url}}/nbi/data/subscribe/update
Content-Type: application/json
authorization: Bearer {{nbi_token}}

{
    "_id": "603f59c5-3ebb-5f34-a30b-1e8a-1003950a",
    "username": "account01",
    "email": "myacsaccount@gmail.com",
    "password": "123456",
    "active_num": 3,
    "addr_pool": "J8009",
    "up_rate": 20000,
    "down_rate": 20000,
    "ip_addr": "",
    "expire_time": "2022-12-31 00:00:00",
    "status": "enabled",
    "remark": "demo user"
}
```

Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Wed, 03 Mar 2021 15:00:11 GMT
Content-Length: 73

{
  "code": 0,
  "msgtype": "info",
  "msg": "Success",
  "data": null
}
```

## Query Users

Request

```
GET http://{{nbi_url}}/nbi/data/subscribe/query
Content-Type: application/json
authorization: Bearer {{nbi_token}}
```

Query parameter

> Reference [API Param Rules](ApiParamRules)<!-- @IGNORE PREVIOUS: link -->
>
> Such as user matching  `?equal[username]=account01`


Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Wed, 03 Mar 2021 15:21:56 GMT
Transfer-Encoding: chunked

[
    {
        "_id": "603f59c5-3ebb-5f34-a30b-1e8a-1003950a",
        "username": "account01",
        "email": "myacsaccount@gmail.com",
        "password": "123456",
        "active_num": 3,
        "addr_pool": "J8009",
        "up_rate": 20000,
        "down_rate": 20000,
        "ip_addr": "",
        "expire_time": "2022-12-31 00:00:00",
        "status": "enabled",
        "remark": "demo user"
    },
    ...
]
```