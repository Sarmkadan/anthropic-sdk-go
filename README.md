# Claude SDK for Go

<!-- x-release-please-start-version -->

<a href="https://pkg.go.dev/github.com/anthropics/anthropic-sdk-go"><img src="https://pkg.go.dev/badge/github.com/anthropics/anthropic-sdk-go.svg" alt="Go Reference"></a>

<!-- x-release-please-end -->

![Build](https://github.com/sarmkadan/anthropic-sdk-go/actions/workflows/build.yml/badge.svg)

The Claude SDK for Go provides access to the [Claude API](https://docs.anthropic.com/en/api/) from Go applications.

## Documentation

Full documentation is available at **[platform.claude.com/docs/en/api/sdks/go](https://platform.claude.com/docs/en/api/sdks/go)**.

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
)
```

<!-- x-release-please-end -->

Or explicitly add the dependency:

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/anthropics/anthropic-sdk-go@v1.46.0'
```

<!-- x-release-please-end -->

## Getting started

```go
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func main() {
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"), // defaults to os.LookupEnv("ANTHROPIC_API_KEY")
	)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
		Model: anthropic.ModelClaudeOpus4_6,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)
}
```

## Examples

The SDK includes comprehensive usage examples in the [`examples/`](examples/) directory:

- **[BasicUsage.go](examples/BasicUsage.go)** – Minimal setup and first call
- **[AdvancedUsage.go](examples/AdvancedUsage.go)** – Configuration, custom options, error handling, and tool use
- **[IntegrationExample.go](examples/IntegrationExample.go)** – Application integration patterns with dependency injection and context management

## Requirements

Go 1.22+

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## Docker Support


The SDK can be built and run using Docker for development and deployment.



### Building the SDK with Docker


```bash
# Build the SDK binary in a Docker container
docker build -t sarmkadan/anthropic-sdk-go .
```


### Running the SDK with Docker


```bash
# Run the SDK container with your Anthropic API key
docker run --rm -it   -e ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY:-}   sarmkadan/anthropic-sdk-go
```

### Using docker-compose for Development


```bash
# Start development environment with build cache volumes
docker-compose up -d


# Build and run
docker-compose build
docker-compose up

# Run tests
docker-compose up test

# Run example applications
docker-compose up basic-example
```

### Development with Docker


The project includes a docker-compose.yml file for development with:
- Build cache volumes for Go modules and build artifacts
- Environment variable support for ANTHROPIC_API_KEY
- Health checks for the container
- Example service demonstrating SDK usage


### Production Deployment


The Docker setup uses a multi-stage build to create a minimal production image:
- Builder stage compiles the SDK binary
- Final stage uses Alpine Linux for minimal footprint
- Only the compiled binary and runtime dependencies are included
- No build tools or source code in the final image


### Dockerfile Features
- Multi-stage build for optimized image size
- Alpine-based for minimal footprint
- Builds the SDK binary in the builder stage
- Final image contains only runtime dependencies

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
