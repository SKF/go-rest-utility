package main

import (
	"context"
	"fmt"
	"log"
	"time"

	dd_http "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	rest "github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-rest-utility/client/auth"
	"github.com/SKF/go-rest-utility/client/retry"
)

type GetNodeResponse struct {
	Node Node `json:"node"`
}

type Node struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

func main() {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("eu-west-1"),
		config.WithSharedConfigProfile("hierarchy_playground"),
	)
	if err != nil {
		log.Fatal(err)
	}

	provider := &auth.SecretCredentialsTokenProvider{
		SecretID: "arn:aws:secretsmanager:eu-west-1:633888256817:secret:user-credentials/hierarchy_service",
		SecretsClient: auth.SecretsManagerV2Client{
			Client: secretsmanager.NewFromConfig(cfg),
		},
		Retry: &retry.ExponentialJitterBackoff{
			Base:        time.Second,
			Cap:         5 * time.Second, //nolint:gomnd
			MaxAttempts: 10,              //nolint:gomnd
		},
	}

	ctx := context.Background()

	client := rest.NewClient(
		rest.WithBaseURL("https://api.sandbox.hierarchy.enlight.skf.com/"),
		rest.WithTokenProvider(provider),
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
