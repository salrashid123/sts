##  Security Token Service (STS) Credentials for HTTP and gRPC (rfc8693)


This library provides `HTTP` and `gRPC` credentials where the final `access_token` is acquired through [STS OAuth 2.0 Token Exchange:  rfc8693](https://www.rfc-editor.org/rfc/rfc8693) 


You can use this library to setup sts credentials for use with either `net/http` Client or  gRPC `RPCCredentials` where you exchange an intermediate credential with an STS server for a final `access_token`.   The final token is then used to access the resource server


![images/sts.png](images/sts.png)


---

The output uses a live sts and grpc server accessible live:

see `examples/` folder uses a sample trivial implementation from a different repo which i've deployed on cloud run

the STS server accepts a bearer token `iammtheeggman` and responds back with a new token `iamthewalrus`
 
see [Serverless Security Token Exchange Server(STS) and gRPC STS credentials](https://github.com/salrashid123/sts_server/blob/main/sts_server/sts_server.go#L44)

```golang
const (
	inboundPassphrase  = "iamtheeggman"
	outboundPassphrase = "iamthewalrus"
)
```

```log
$ go run main.go 
2023/10/27 12:00:27 New Token: iamthewalrus
2023/10/27 12:00:27 {
  "args": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Authorization": "Bearer iamthewalrus", 
    "Host": "httpbin.org", 
    "User-Agent": "Go-http-client/2.0", 
    "X-Amzn-Trace-Id": "Root=1-653bde9b-7b86e72a6e1006f802a9bc80"
  }, 
  "origin": "108.51.25.168", 
  "url": "https://httpbin.org/get"
}
2023/10/27 12:00:28 RPC Response: message:"Hello unary RPC msg   from K_REVISION grpcserver-00006-vsr"
```

##### `http`:

the first output shows the echo response back from httpbin:

`client-->stsserver`  --> `sts server responds back with a token` --> `client sends new token to httpbin`

the output from httpbin's echo shows the bearer token it recieved (whcih is `iamthewalrus`)

##### `grpc`

the second output shows the grpc response from a server on cloud run

`client-->stsserver`  --> `sts server responds back with a token` --> `client sends new token to a grpc server which echo's back some data`

the GRPC server implenentation here only accepts a bearer token of `"iamthewalrus"` (which is what the sts server respond back with )

---

##### References

* [Serverless Security Token Exchange Server(STS)](https://github.com/salrashid123/sts_server)
* [Certificate Bound Tokens using Security Token Exchange Server (STS)](https://github.com/salrashid123/cert_bound_sts_server)


---

### HTTP


```golang
import (
	stshttp "github.com/salrashid123/sts/http"
)


	rootTS := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: secret,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Duration(time.Second * 60)),
	})

	// HTTP
	stsTokenSource, _ := stshttp.STSTokenSource(
		&stshttp.STSTokenConfig{
			TokenExchangeServiceURI: *stsaddress,
			Resource:                *stsaudience,
			Audience:                *stsaudience,
			Scope:                   *scope,
			SubjectTokenSource:      rootTS,
			SubjectTokenType:        "urn:ietf:params:oauth:token-type:access_token",
			RequestedTokenType:      "urn:ietf:params:oauth:token-type:access_token",
			HTTPClient:              http.DefaultClient,
		},
	)

	tok, err := stsTokenSource.Token()
	log.Printf("New Token: %s", tok.AccessToken)

	client := oauth2.NewClient(context.TODO(), stsTokenSource)
	resp, err := client.Get(*httpAddress)
```

---

### gRPC

Note that upstream [google.golang.org/grpc/credentials/sts](https://pkg.go.dev/google.golang.org/grpc/credentials/sts) provides the same credential object except that this variation allows for


* arbitrary HTTPClients  [issue #5611](https://github.com/grpc/grpc-go/pull/5611)

* allowing source tokens from arbitrary `oauth2.TokenSource`:

```golang
	// token source for the subject token
	SubjectTokenSource *oauth2.TokenSource
```    


Example usage:

```golang
import (
	stsgrpc "github.com/salrashid123/sts/grpc"
)


	rootTS := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: secret,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Duration(time.Second * 60)),
	})

	ce := credentials.NewTLS(&tls.Config{})

	// ### test with sts
	stscreds, err := stsgrpc.NewCredentials(stsgrpc.Options{
		TokenExchangeServiceURI: *stsaddress,
		Resource:                *stsaudience,
		Audience:                *stsaudience,
		Scope:                   *scope,
		SubjectTokenSource:      &rootTS,
		SubjectTokenType:        "urn:ietf:params:oauth:token-type:access_token",
		RequestedTokenType:      "urn:ietf:params:oauth:token-type:access_token",
		HTTPClient:              http.DefaultClient,
	})

	ctx := context.Background()

	conn, err := grpc.Dial(*grpcAddress,
		grpc.WithTransportCredentials(ce),
		grpc.WithPerRPCCredentials(stscreds))

	defer conn.Close()
	c := pb.NewEchoServerClient(conn)

	r, err := c.SayHello(ctx, &pb.EchoRequest{Name: "unary RPC msg "})

```
