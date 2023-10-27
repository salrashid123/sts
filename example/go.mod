module main

go 1.19

require (
	github.com/salrashid123/sts/grpc v0.0.0
	github.com/salrashid123/sts/http v0.0.0
	github.com/salrashid123/sts_server/echo v0.0.0
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094
	google.golang.org/grpc v1.49.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/salrashid123/sts/http => ../http

replace github.com/salrashid123/sts/grpc => ../grpc

replace github.com/salrashid123/sts_server/echo => ./echo
