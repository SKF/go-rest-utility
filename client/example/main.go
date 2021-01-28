package main

import (
	"context"
	"fmt"
	"log"

	dd_http "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"

	rest "github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-rest-utility/client/auth"
)

type GetNodeResponse struct {
	Node Node `json:"node"`
}

type Node struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

func main() {
	identityToken := "eyj..."

	ctx := context.Background()

	client := rest.NewClient(
		rest.WithBaseURL("https://api.sandbox.hierarchy.enlight.skf.com/"),
		rest.WithTokenProvider(auth.RawToken(identityToken)),
		rest.WithDatadogTracing(dd_http.RTWithServiceName("my-example-service")),
	)

	request := rest.Get("nodes/{id}").
		Assign("id", "df3214a6-2db7-11e8-b467-0ed5f89f718b").
		SetHeader("Accept", "application/json")

	response, err := client.Do(ctx, request)
	if err != nil {
		log.Fatal(err)
	}

	node := GetNodeResponse{}
	if err := response.Unmarshal(&node); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s: %s\n", node.Node.Label, node.Node.Description)
}
