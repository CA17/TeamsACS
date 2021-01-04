# API User Management


## Initialize the superuser

The superuser needs to be configured in a command line terminal, execute the following command on the TeamsACS server

```
teamsacs --init-admin -admin=<username>

```

The above operation will create or update the administrator user, and print the administrator API Secret on the console after successful execution

## Create a normal API user

Request

```
POST http://{{nbi_url}}/nbi/opr/add
Content-Type: application/json
authorization: Bearer {{nbi_token}}

{
    "id": "5f9ec620-bad9-51d7-a062-e6ac-00f0fa6s",
    "email": "test@teamsacs.com",
    "username": "opr",
    "level": "opr",
    "remark": "opr"
}

```

Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding
Date: Mon, 04 Jan 2021 01:31:16 GMT

{
    "code": 0,
    "msgtype": "info",
    "msg": "Success",
    "data": null
}

```

## Query API users

Request

```
GET http://{{nbi_url}}/nbi/opr/query?filter[username]=opr
authorization: Bearer {{nbi_token}}

```

Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding
Date: Mon, 04 Jan 2021 02:00:14 GMT

[
  {
    "_id": "5f9ec620-bad9-51d7-a062-e6ac-00f0fa61",
    "api_secret": "5ff26fe4-9b57-d15b-65f9-b5ea-1ed761d0",
    "email": "test@teamsacs.com",
    "level": "opr",
    "remark": "opr",
    "status": "enabled",
    "username": "opr"
  }
]
```

## Update API user

Request

```
POST http://{{nbi_url}}/nbi/opr/update
Content-Type: application/json
authorization: Bearer {{nbi_token}}

{
  "email": "test2@teamsacs.com",
  "username": "opr"
}
```

Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding
Date: Mon, 04 Jan 2021 03:47:04 GMT

{
  "code": 0,
  "msgtype": "info",
  "msg": "Success",
  "data": null
}
```

## Delete API user

Request

```
GET http://{{nbi_url}}/nbi/opr/delete?filter[username]=opr
authorization: Bearer {{nbi_token}}
```


Response

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Vary: Accept-Encoding
Date: Mon, 04 Jan 2021 03:53:46 GMT

{
  "code": 0,
  "msgtype": "info",
  "msg": "Success",
  "data": null
}
```




