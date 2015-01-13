# mirror
A small web server that will just output the body and url of each request to standard output

## Install
`$ go get github.com/fredr/mirror`

## Use
##### Start
`$ mirror`
##### Request
`$ curl "http://localhost:12345/some/url" -d 'Hello World' -XPUT`
