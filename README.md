<p align="center">
  <a href="https://clerk.com?utm_source=github&utm_medium=sdk_go" target="_blank" rel="noopener noreferrer">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="https://images.clerk.com/static/logo-dark-mode-400x400.png">
      <img src="https://images.clerk.com/static/logo-light-mode-400x400.png" height="64">
    </picture>
  </a>
  <br />
</p>

# Clerk Go SDK

**This is still a work in progress. The current stable release is v1. See the main branch for the stable release.**

The official [Clerk](https://clerk.com) Go client library for accessing the [Clerk Backend API](https://clerk.com/docs/reference/backend-api).

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2)
[![chat on Discord](https://img.shields.io/discord/856971667393609759.svg?logo=discord)](https://discord.com/invite/b5rXHjAg7A)
[![documentation](https://img.shields.io/badge/documentation-clerk-green.svg)](https://clerk.com/docs)
[![twitter](https://img.shields.io/twitter/follow/ClerkDev?style=social)](https://twitter.com/intent/follow?screen_name=ClerkDev)

## Requirements

- Go 1.19 or later.

## Installation

If you are using Go Modules and have a `go.mod` file in your project's root, you can import clerk-sdk-go directly.

```go
import (
    "github.com/clerk/clerk-sdk-go/v2"
)
```

Alternatively, you can `go get` the package explicitly.

```
go get -u github.com/clerk/clerk-sdk-go/v2
```

## Usage

For details on how to use this module, see the [Go Documentation](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2).

The library has a resource-based structure which follows the way the [Clerk Backend API](https://clerk.com/docs/reference/backend-api) resources are organized.
Each API supports specific operations, like Create or List. While operations for each resource vary, a similar pattern is applied throughout the library.

In order to start using API operations the library needs to be configured with your Clerk API secret key. Depending on your use case,
there's two ways to use the library; with or without a client.

For most use cases, the API without a client is a better choice.

On the other hand, if you need to set up multiple Clerk
API keys, using clients for API operations provides more flexibility.

Let's see both approaches in detail.

### Usage without a client

If you only use one API key, you can import the packages required for the APIs you need to interact with and call functions for API operations.
```go
import (
    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/clerk/clerk-sdk-go/v2/$resource$"
)

// Each operation requires a context.Context as the first argument.
ctx := context.Background()

// Set the API key
clerk.SetKey("sk_live_XXX")

// Create
resource, err := $resource$.Create(ctx, &$resource$.CreateParams{})

// Get
resource, err := $resource$.Get(ctx, id)

// Update
resource, err := $resource$.Update(ctx, id, &$resource$.UpdateParams{})

// Delete
resource, err := $resource$.Delete(ctx, id)

// List
list, err := $resource$.List(ctx, &$resource$.ListParams{})
for _, resource := range list.$Resource$s {
    // do something with the resource
}
```

### Usage with a client

If you're dealing with multiple keys, it is recommended to use a client based approach. The API packages for each
resource export a Client, which supports all the API's operations.
This way you can create as many clients as needed, each with their own API key.

```go
import (
    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/clerk/clerk-sdk-go/v2/$resource$"
)

// Each operation requires a context.Context as the first argument.
ctx := context.Background()

// Initialize a client with an API key
config := &$resource$.ClientConfig{}
config.Key = "sk_live_XXX"
client := $resource$.NewClient(config)

// Create
resource, err := client.Create(ctx, &$resource$.CreateParams{})

// Get
resource, err := client.Get(ctx, id)

// Update
resource, err := client.Update(ctx, id, &$resource$.UpdateParams{})

// Delete
resource, err := client.Delete(ctx, id)

// List
list, err := client.List(ctx, &$resource$.ListParams{})
for _, resource := range list.$Resource$s {
    // do something with the resource
}
```

Here's an example of how the above operations would look like for specific APIs.

```go
import (
    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/clerk/clerk-sdk-go/v2/organization"
    "github.com/clerk/clerk-sdk-go/v2/organizationmembership"
    "github.com/clerk/clerk-sdk-go/v2/user"
)

func main() {
    // Each operation requires a context.Context as the first argument.
    ctx := context.Background()

    // Set the API key
    clerk.SetKey("sk_live_XXX")

    // Create an organization
    org, err := organization.Create(ctx, &organization.CreateParams{
        Name: clerk.String("Clerk Inc"),
    })

    // Update the organization
    org, err = organization.Update(ctx, org.ID, &organization.UpdateParams{
        Slug: clerk.String("clerk"),
    })

    // List organization memberships
    listParams := organizationmembership.ListParams{}
    listParams.Limit = clerk.Int64(10)
    memberships, err := organizationmembership.List(ctx, params)
    if memberships.TotalCount < 0 {
        return
    }
    membership := memberships[0]

    // Get a user
    usr, err := user.Get(ctx, membership.UserID)
}
```

### Accessing API responses

Each resource that is returned by an API operation has a `Response` field which
contains information about the response that was sent from the Clerk Backend API.

The `Response` contains fields like the the raw HTTP response's headers,
the status and the raw response body. See the [APIResponse](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#APIResponse)
documentation for available fields and methods.

```go
dmn, err := domain.Create(context.Background(), &domain.CreateParams{})
if !dmn.Response.Success() {
    dmn.Response.TraceID
}
```

### Errors

For cases where an API operation returns an error, the library will try to return an `APIErrorResponse`.
The `APIErrorResponse` type provides information such as the HTTP status code of the response, a list of errors
and a trace ID that can be used for debugging.

The `APIErrorResponse` is an `APIResource`. You can [access the API response](#accessing-api-responses) for errors as well.

See the [APIErrorResponse](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#APIErrorResponse) documentation for available fields and methods.

```go
_, err := user.List(context.Background(), &user.ListParams{})
if apiErr, ok := err.(*clerk.APIErrorResponse); ok {
    apiErr.TraceID
    apiErr.Error()
    apiErr.Response.RawJSON
}
```

### HTTP Middleware

The library provides two functions that can be used for adding authentication with Clerk to HTTP handlers.

Both middleware functions support header based authentication with a bearer token. The token is parsed, verified and
its claims are extracted as [SessionClaims](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#SessionClaims).

The claims will then be made available in the `http.Request.Context` for the next handler in the chain. The library
provides a helper for accessing the claims from the context, [SessionClaimsFromContext](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#SessionClaimsFromContext).

The two middleware functions are [WithHeaderAuthorization](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http#WithHeaderAuthorization)
and [RequireHeaderAuthorization](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http#RequireHeaderAuthorization).
Their difference is that the `RequireHeaderAuthorization` calls `WithHeaderAuthorization` under the hood, but responds
with HTTP 403 Forbidden if it fails to detect valid session claims.

Let's see an example of how the middleware can be used.

```go
import (
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

func main() {
	clerk.SetKey("sk_live_XXX")

	mux := http.NewServeMux()
	mux.HandleFunc("/", publicRoute)
	protectedHandler := http.HandlerFunc(protectedRoute)
	mux.Handle(
		"/protected",
		clerkhttp.WithHeaderAuthorization()(protectedHandler),
	)

	http.ListenAndServe(":3000", mux)
}

func publicRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"access": "public"}`))
}

func protectedRoute(w http.ResponseWriter, r *http.Request) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"access": "unauthorized"}`))
		return
	}
	fmt.Fprintf(w, `{"user_id": "%s"}`, claims.Subject)
}
```

Both `WithHeaderAuthorization` and `RequireHeaderAuthorization` can be
customized. They accept various options as functional arguments.

For a comprehensive list of available options check the
[AuthorizationParams](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http#AuthorizationParams) documentation.

### Testing

There are various ways to mock the library in your test suite.

#### Usage without client

If you're using the library without instantiating clients for APIs, you can
set the package's `Backend` with a custom configuration.

1. Stub out the HTTP client's transport

```go
func TestWithCustomTransport(t *testing.T) {
    clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
        HTTPClient: &http.Client{
            Transport: mockRoundTripper,
        },
    }))
}

type mockRoundTripper struct {}
// Implement the http.RoundTripper interface.
func (r *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
    // Construct and return the http.Response.
}
```

2. Use a httptest.Server

Similar to the custom http.Transport approach, you can use the net/http/httptest package's utilities and
provide the http.Client to the package's Backend.

```go
func TestWithMockServer(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Write the response.
    }))
    defer ts.Close()
    clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
        HTTPClient: ts.Client(),
        URL: &ts.URL,
    }))
}
```

3. Implement your own Backend

You can implement your own [Backend](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#Backend) and set it as the package's default Backend.

```go
func TestWithCustomBackend(t *testing.T) {
    clerk.SetBackend(&customBackend{})
}

type customBackend struct {}
// Implement the Backend interface
func (b *customBackend) Call(ctx context.Context, r *clerk.APIRequest, reader *clerk.ResponseReader) error {
    // Construct a clerk.APIResponse and use the reader's Read method.
    reader.Read(&clerk.APIResponse{})
}
```

#### Usage with client

If you're already using the library by instantiating clients for API operations,
or you need to ensure your test suite can safely run in parallel, you can simply
pass a custom http.Client to your clients.

1. Stub out the HTTP client's transport

```go
func TestWithCustomTransport(t *testing.T) {
    config := &clerk.ClientConfig{}
    config.HTTPClient = &http.Client{
        Transport: mockRoundTripper,
    }
    client := user.NewClient(config)
}

type mockRoundTripper struct {}
// Implement the http.RoundTripper interface.
func (r *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
    // Construct and return the http.Response.
}
```

2. Use a httptest.Server

Similar to the custom http.Transport approach, you can use the net/http/httptest package's utilities and
provide the http.Client to the API client.

```go
func TestWithMockServer(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Write the response.
    }))
    defer ts.Close()
    config := &clerk.ClientConfig{}
    config.HTTPClient = ts.Client()
    config.URL = &ts.URL
    client := user.NewClient(config)
}
```

3. Implement your own Backend

You can implement your own [Backend](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#Backend) and set it as the API client's Backend.

```go
func TestWithCustomBackend(t *testing.T) {
    client := user.NewClient(&clerk.ClientConfig{})
    client.Backend = &customBackend{}
}

type customBackend struct {}
// Implement the Backend interface
func (b *customBackend) Call(ctx context.Context, r *clerk.APIRequest, reader *clerk.ResponseReader) error {
    // Construct a clerk.APIResponse and use the reader's Read method.
    reader.Read(&clerk.APIResponse{})
}
```

## Development

Contributions are welcome. If you submit a pull request please keep in mind that

1. Code must be `go fmt` compliant.
2. All packages, types and functions should be documented.
3. Ensure that `go test ./...` succeeds. Ideally, your pull request should include tests.
4. If your pull request introduces a new API or API operation, run `go generate ./...` to generate the necessary API functions.

## License

This SDK is licensed under the MIT license found in the [LICENSE](./LICENSE) file.
