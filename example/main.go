package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	stsgrpc "github.com/salrashid123/sts/grpc"
	stshttp "github.com/salrashid123/sts/http"

	pb "github.com/salrashid123/sts_server/echo"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	httpAddress = flag.String("httpAddress", "https://httpbin.org/get", "host:port of http server")

	stsaddress  = flag.String("stsaddress", "https://stsserver-6w42z6vi3q-uc.a.run.app/token", "STS Server address")
	stsaudience = flag.String("stsaudience", "stsserver-6w42z6vi3q-uc.a.run.app", "the audience and resource value to send to STS server")
	scope       = flag.String("scope", "https://www.googleapis.com/auth/cloud-platform", "scope to send to STS server")

	grpcAddress = flag.String("grpcAddress", "grpcserver-6w42z6vi3q-uc.a.run.app:443", "host:port of gRPC server")
)

const (
	secret = "iamtheeggman"
)

func main() {
	flag.Parse()

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
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New Token: %s", tok.AccessToken)

	client := oauth2.NewClient(context.TODO(), stsTokenSource)
	resp, err := client.Get(*httpAddress)
	if err != nil {
		log.Printf("Error creating client %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error connecting to server %v", http.StatusText(resp.StatusCode))
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	log.Printf("%s", bodyString)

	// gRPC

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
	if err != nil {
		log.Fatalf("unable to create TokenSource: %v\n", err)
	}

	ctx := context.Background()

	conn, err := grpc.Dial(*grpcAddress,
		grpc.WithTransportCredentials(ce),
		grpc.WithPerRPCCredentials(stscreds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewEchoServerClient(conn)

	r, err := c.SayHello(ctx, &pb.EchoRequest{Name: "unary RPC msg "})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("RPC Response: %s", r)

}
