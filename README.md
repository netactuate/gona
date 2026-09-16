# gona

Go client for the NetActuate API. It is the SDK behind the
[NetActuate Terraform provider](https://github.com/netactuate/terraform-provider-netactuate)
and covers cloud and dedicated compute, VPCs and their gateways, load balancers, BGP and
anycast, object and block storage, DNS, secrets, tags, OIDC, and NKE clusters with their
add-ons.

The full Go API reference is published at
[pkg.go.dev/github.com/netactuate/gona](https://pkg.go.dev/github.com/netactuate/gona/gona).
For the platform itself, its concepts and the REST API behind this client, see the
[NetActuate documentation](https://netactuate.com/docs).

## Requirements

- Go 1.19 or later
- A NetActuate account and an API key

## Installation

```bash
go get github.com/netactuate/gona@latest
```

## Two clients, two API versions

The platform exposes two APIs and this package provides a client for each. Which one you need
depends on the resource, and a program that touches both will construct both.

| client | API | covers |
| --- | --- | --- |
| `gona.Client` | vAPI2 | cloud and dedicated servers, BGP, DNS, packages, images |
| `gona.V3Client` | vAPI3 | VPCs, gateways, load balancers, storage, NKE, secrets, OIDC |

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/netactuate/gona/gona"
)

func main() {
	key := os.Getenv("NETACTUATE_API_KEY")

	servers, err := gona.NewClient(key).GetServers()
	if err != nil {
		log.Fatalf("list servers: %v", err)
	}
	for _, s := range servers {
		fmt.Printf("%d %s\n", s.ID, s.Name)
	}

	// An empty base URL uses the production vAPI3 endpoint.
	vpcs, err := gona.NewV3Client(key, "").ListVPCs()
	if err != nil {
		log.Fatalf("list VPCs: %v", err)
	}
	for _, v := range vpcs {
		fmt.Printf("%d %s\n", v.VPCID, v.Metadata.Label)
	}
}
```

## Authentication

Both clients take the API key as their first argument. Generate one in the portal under your
account's API settings, and keep it out of source control.

`gona.GetKeyFromEnv()` reads `NA_API_KEY` for callers that want that convention. The Terraform
provider uses `NETACTUATE_API_KEY`, so read the environment directly when you want to match it.

## Errors

Every call returns an error rather than a zero value on failure. Two helpers distinguish a
missing object from a failed request, which is what a caller needs in order to tell "this was
deleted" from "this did not work":

```go
server, err := client.GetServer(id)
if gona.IsNotFound(err) {
	// The server no longer exists.
}
```

`IsNotFound` covers vAPI2 and `IsV3NotFound` covers vAPI3.

## Custom endpoints

`NewClientCustom` and the second argument to `NewV3Client` override the base URL when you need
to point at something other than production.

## Contributing

Run the tests before opening a pull request:

```bash
go build ./...
go vet ./...
go test ./...
```

The tests use `httptest` servers and never contact the live API.

## License

See [LICENSE](LICENSE).
