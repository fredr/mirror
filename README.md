# mirror
A small web server that will output any requests against it

## Install
`$ go get github.com/fredr/mirror`

## Use
##### Start
`$ mirror`
##### Request
`$ curl "http://localhost:12345/some/url" -d 'Hello World' -XPUT`
##### Output
```
PUT /some/url HTTP/1.1
Host: localhost:12345
Accept: */*
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.30.0

Hello World
```
