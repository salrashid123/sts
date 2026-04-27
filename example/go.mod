module main

go 1.26.2

require (
	github.com/salrashid123/sts/grpc v0.0.0
	github.com/salrashid123/sts/http v0.0.0
	github.com/salrashid123/sts_server/echo v0.0.0
	golang.org/x/oauth2 v0.36.0
	google.golang.org/grpc v1.80.0
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20260420184626-e10c466a9529 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/salrashid123/sts/http => ../http

replace github.com/salrashid123/sts/grpc => ../grpc

replace github.com/salrashid123/sts_server/echo => ./echo
