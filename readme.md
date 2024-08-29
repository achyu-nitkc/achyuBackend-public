# achyu Backend
## Working flow
### email authentication
1. **Call /create**
    - Insert data into verify table
    - Generate jwt with short expiration
    - Send confirm Mail with code
2. **Call /verify**
    - Check jwt
    - Check verify code
    - Generate jwt
    - Move data from verify table to user table
3. **Call /login**
   - Check password
   - Generate jwt
### oauth2.0
1. Call /Oauth
2. Redirect to google oauth2.0
3. Generate jwt
   - if needed, add data to user table
## Endpoints
### /login [POST]
request
```go
type requestLoginBody struct {
	Email    string
	Password string
}
```
response
```text
200 statusOK
    Token(jwt,exp 3 month)
500 statusInternalServerError
401 statusUnauthorized
400 statusBadRequest
```
### /create [POST]
request
```go
type requestCreateBody struct {
	Email       string
	Password    string
	DisplayName string
}
```
response
```text
200 statusOK
    Token(jwt,exp 5 min)
500 statusInternalServerError
401 statusUnauthorized
400 statusBadRequest
```
### /verify [POST]
request
```go
type requestVerifyBody struct {
	token string
	code string
}
```
response
```text
200 statusOK
    Token(jwt,exp 3 month)
500 statusInternalServerError
401 statusUnauthorized
400 statusBadRequest
```
### /Oauth [POST]
Oauth redirect function
### /OauthCallback
Oauth callback function
```text
200 statusOK
    Token(jwt,exp 3 month)
```
### /post [POST]
request
```go
type requestPost struct {
Token     string
Latitude  float64
Longitude float64
Content   string
ImageURL  string
}
```
response
```text
200 statusOK
500 statusInternalServerError
401 statusUnauthorized
400 statusBadRequest
```
### /delete [POST]
request
```go
type requestDelete struct {
	Token  string
	PostId string
}
```
response
```text
200 statusOK
500 statusInternalServerError
401 statusUnauthorized
400 statusBadRequest
```
### /get [GET]
request
```go
type requestGet struct {
	Token     string
	Latitude  float64
	Longitude float64
}
```
response
```text
200 statusOK
    []ResponseGet
    
    type ResponseGet struct {
    	PostId           string  `json:"postId"`
    	DisplayName      string  `json:"displayName"`
    	Content          string  `json:"content"`
    	ImageURL         string  `json:"imageUrl"`
    	Latitude         float64 `json:"latitude"`
    	Longitude        float64 `json:"longitude"`
    	Address          string  `json:"address"`
    	ConstructionName string  `json:"constructionName"`
    	roadName         string  `json:"roadName"`
    }
    
500 statusInternalServerError
401 statusUnauthorized
400 statusBadRequest
```
## run
```shell
$ git clone https://github.com/achyu-nitkc/achyuBackend.git
$ cd achyuBackend
# ADD SECRET FILE
$ sudo docker build -t achyu .
$ sudo docker run -p 8080:8080 achyu 
```
## Secrets
auth/jwtSecret.go
```go
package auth

func jwtSecret() []byte {
   return []byte("")
}
```
auth/oAuthSecret.go
```go
package auth
//https://cloud.google.com/apigee/docs/api-platform/security/oauth/access-tokens?hl=ja
func oauthSecret() (ClientID, ClientSecret string) {
   clientID := ""
   clientSecret := ""
   return clientID, clientSecret
}

func oauthStateString() string {
   return ""
}
```
auth/smtpConfig.go
```go
package auth

//Get application password
//https://myaccount.google.com/apppasswords

func config() (hostname string, port int, username string, password string) {
   return "smtp.gmail.com", 587, "Your Email Address", "Your Application Password"
}
```

yolp/yolpSecret.go

[yolp dashbord](https://e.developer.yahoo.co.jp/dashboard)
```go
package yolp

func secret() string {
   return "appid="+"Your Client Id"
}
```