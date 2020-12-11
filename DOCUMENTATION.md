# API documentation for REST architecture
## <span style="color:#fc326a">**Table of contents**</span>

- [Authentication](#<span-style="color:#fc326a">**Authentication**</span>)
- [Registration](#<span-style="color:#fc326a">**Registration**</span>)
- [User Router](#<span-style="color:#fc326a">**User-Router**</span>)
---
## <span style="color:#fc326a">**Authentication**</span>

### Description
Current **Authentication** section includes authentication via http request. After *authentication*, as a response client recieve access and refresh tokens in order to use the service.

#### <span style="color:#6f42f5">**HTTP request**</span>
*Request-Line*
````
POST /v1/login HTTP/1.1
````
*Required Headers*
````
Content-Type: application/json
````
#### *HTTP body request params*
```json
{
        "login": "your_login",
        "password":"your_pass"
}
```

#### <span style="color:#6f42f5">**HTTP response**</span>
- Status Created 201

#### *Errors*
- Status Bad Request 400
- Status Unprocessable Entity 422
- Status Not Found 404
- Status Internal Server Error 500

#### *Sample body response*
```json
{
    "Error":false,
    "Message":{
        "accessToken":"accessToken",
        "refreshToken":"refreshToken"
    }
}
```
---
## <span style="color:#fc326a">**Registration**</span>

### Description
**Registration** part creates user and redirect the http request to **Login** section.

#### <span style="color:#6f42f5">**HTTP request**</span>
*Request-Line*
````
POST /v1/register HTTP/1.1
````
*Required Headers*
````
Content-Type: application/json
````
#### *HTTP body request params*
```json
{
    "login": "your_login",
    "password":"your_pass",
    "first_name":"your name",
    "sur_name":"your second name",
    "email":"your_email@example.com"
}
```

#### <span style="color:#6f42f5">**HTTP response**</span>
- Status Created 201

#### *Errors*
- Status Bad Request 400
- Status Unprocessable Entity 422
- Status Internal Server Error 500

#### *Sample body response*
```json
{
    "Error":false,
    "Message":"User created"
}
```
---
## <span style="color:#fc326a">**User Router**</span>

### Description
**User Router** contains methods to *save, read, and delete* user's content.

### <span style="color:#f760c5">**User Delete**</span>
#### <span style="color:#6f42f5">**HTTP request**</span>
*Request-Line*
````
DELETE /v1/{id:[0-9]+} HTTP/1.1
````
*Required Headers*
````
Content-Type: application/json
````
#### *HTTP body request params*
```json
{
    "login":"your_login" // optional: "id":"your_id"
}
```

#### <span style="color:#6f42f5">**HTTP response**</span>
- Status OK 200

#### *Errors*
- Status Bad Request 400
- Status Internal Server Error 500

#### *Sample body response*
```json
{
    "Error":false,
    "Message":"User deleted"
}
```
### <span style="color:#f760c5">**User Update**</span>
#### <span style="color:#6f42f5">**HTTP request**</span>
*Request-Line*
````
POST/PUT /v1/{id:[0-9]+} HTTP/1.1
````
*Required Headers*
````
Content-Type: application/json
````
#### *HTTP body request params*
```json
{
    "first_name":"your first_name",
    "sur_name":"your sur_name",
    "email":"your email"
}
```

#### <span style="color:#6f42f5">**HTTP response**</span>
- Status OK 200

#### *Errors*
- Status Internal Server Error 500
- Status Unprocessable Entity 422

#### *Sample body response*
```json
{
    "Error":false,
    "Message":"Successfully updated"
}
```
### <span style="color:#f760c5">**User Task Read**</span>
#### <span style="color:#6f42f5">**HTTP request**</span>
*Request-Line*
````
GET /v1/{id:[0-9]+}/tasks HTTP/1.1
````
*Required Headers*
````
Content-Type: application/json
````
#### *HTTP body request params*
```json
{
    // optional: "login": "your_login",
    //           "id": 
}
```

#### <span style="color:#6f42f5">**HTTP response**</span>
- Status OK 200

#### *Errors*
- Status Internal Server Error 500
- Status Unprocessable Entity 422

#### *Sample body response*
```json
{
    "Error":false,
    "Message":"Successfully updated"
}
```
---