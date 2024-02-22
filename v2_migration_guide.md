# Migration guide for v2

The `v2` version of the Clerk Go SDK is a complete rewrite and introduces
a lot of breaking changes.

## Minimum Go version

The minimum supported Go version for the `v2` version of the Clerk Go SDK is `1.19`.

## New API and package layout

API operations in `v1` of the Clerk Go SDK are organized by service. There are
different services for each API and all services are properties of a single
`clerk.Client`.

Let's see a concrete example. Assume that we want to create an organization
and then list all available organizations.

Here's how we would do it in `v1`.

```go
import (
    "github.com/clerkinc/clerk-sdk-go"
)


func main() {
    client, err := clerk.NewClient("sk_live_XXX")
    if err != nil {
        // handle error
    }

    // Create an organization
    org, err := client.Organizations().Create(clerk.CreateOrganizationParams{
        Name: "Acme Inc",
    })
    if err != nil {
        if errResp, ok := err.(*clerk.ErrorResponse); ok {
            // Access the API errors
            errResp.Errors
        }
    }
    // List all organizations, limit results to one.
    limit := 1
    orgs, err := client.Organizations().ListAll(clerk.ListAllOrganizationsParams{
        Limit: &limit,
    })
    if err != nil {
        // handle the error
    }
    if orgs.TotalCount > 0 {
        // Get the first organization in the list
        org = orgs.Data[0]
    }
}
```

The `v2` version breaks away of the client and services pattern seen above.
In `v2`, every API has its own package, following a resource-based structure.

There are two ways to call API operations in `v2`: with or without a client. Let's see both approaches.

## Usage without a client

In most cases, you'll only have to deal with a single API key in your project.

```go
import (
    "context"

    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/clerk/clerk-sdk-go/v2/organization"
)

func main() {
    ctx := context.Background()
    clerk.SetKey("sk_live_XXX")

    // Create an organization
    org, err := organization.Create(ctx, &organization.CreateParams{
        Name: clerk.String("Acme Inc"),
    })
    if err != nil {
        if apiErr, ok := err.(*clerk.APIErrorResponse); ok {
            // Access the API errors and additional information
            apiErr.TraceID
            apiErr.Error()
            apiErr.Response.RawJSON
        }
    }
    // List all organizations, limit results to one.
    params := &organization.ListParams{}
    params.Limit = clerk.Int64(1)
    list, err := organization.List(ctx, params)
    if err != nil {
        // handle the error
    }
    if list.TotalCount > 0 {
        // Get the first organization in the list
        org = list.Organizations[0]
    }
}
```

## Usage with a client

When you have to deal with more than one API keys in your project, or need more flexibility.

```go
import (
    "context"

    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/clerk/clerk-sdk-go/v2/organization"
)

func main() {
    ctx := context.Background()
    config := &clerk.ClientConfig{}
    config.Key = "sk_live_XXX"
    client := organization.NewClient(config)

    // Create an organization
    org, err := client.Create(ctx, &organization.CreateParams{
        Name: clerk.String("Acme Inc"),
    })
    if err != nil {
        if apiErr, ok := err.(*clerk.APIErrorResponse); ok {
            // Access the API errors and additional information
            apiErr.TraceID
            apiErr.Error()
            apiErr.Response.RawJSON
        }
    }
    // List all organizations, limit results to one.
    params := &organization.ListParams{}
    params.Limit = clerk.Int64(1)
    list, err := organization.List(ctx, params)
    if err != nil {
        // handle the error
    }
    if list.TotalCount > 0 {
        // Get the first organization in the list
        org = list.Organizations[0]
    }
}
```

## List operation responses

A lot of APIs support operations where a resource list is returned. In `v2` of the Clerk Go SDK,
the return types of list operations are similar. They always contain the total count of resources
available in the server, while the (sometimes filtered or paginated) results are returned in a slice.

Here's an example. Replace `$resource$` with a resource (package) name.

```go
import (
    "github.com/clerk/clerk-sdk-go/v2"
    "github.com/clerk/clerk-sdk-go/v2/$resource$"
)

ctx := context.Background()
list, err := $resource$.List(ctx, &$resource$.ListParams{})
// If $resource$ was user, the following line would read
// fmt.Println(list.TotalCount, list.Users)
fmt.Println(list.TotalCount, list.$resource$s)
```

## Different types

In the library's `v1` version, some operations would return the same type, while others wouldn't.

The `v2` version moves to a more resource-oriented approach. Each API resource has a type and all types
are defined in the `clerk` package.

This means there are types for `clerk.User` and `clerk.Domain` and these types are returned by
the Users and Domains API operations respectively.

All types were consolidated and updated to match the currently supported request parameters and responses
of the [Clerk Backend API](https://clerk.com/docs/reference/backend-api). Many fields have been renamed and a lot
of new fields have been added. Deprecated struct fields have been dropped.

## Operation parameters

The `v2` version of Clerk SDK Go introduces another important change about types that can be used as API operation
parameters. Every field for these structs is a pointer.

The `v2` version of the library introduces helper functions to cast basic type values to pointers. These are:
- `clerk.String`
- `clerk.Bool`
- `clerk.Int64`

Using the helpers above, here's how you can invoke an API operation with a `*Params` struct.

```go
domain.Create(context.Background(), &domain.CreateParams{
    Name: clerk.String("clerk.com"),
    IsSatellite: clerk.Bool(true),
})
```

You can explicitly pass zero values with `clerk.String("")` or `clerk.Int64(0)`.

All API operations receive a `context.Context` as the first argument.

## HTTP middleware

The `v1` version of the Clerk Go SDK supports two HTTP middleware functions that can handle authentication with Clerk.
These are `WithSessionV2` and `RequireSessionV2`.

The middleware in `v1` supports authentication with a bearer token in the HTTP request headers, falling back to cookie
based authentication if the "Authorization" header is missing.

In a similar way, the `v2` version of the library also provides two middleware functions that can handle authentication
with Clerk; `WithHeaderAuthorization` and `RequireHeaderAuthorization`.

As the name implies, the new middleware support only authentication with a bearer token, that needs to be present in the
"Authorization" header of the HTTP request.

Cookie based authentication is not supported at all by the `v2` version of the library.

Usage has also changed between `v1` and `v2`, as have the different options that the middleware support.

In `v1`, here's how you would use the `RequireSessionV2` middleware to get access to the session token claims in HTTP handlers.

```go
import (
    "net/http"

    "github.com/clerkinc/clerk-sdk-go"
)

func main() {
    client, err := clerk.NewClient("sk_live_XXX")
    if err != nil {
        panic(err)
    }
    mux := http.NewServeMux()
    mux.Handle("/session", clerk.RequireSessionV2(client)(http.HandlerFunc(handleSession)))
    http.ListenAndServe(":3000", mux)
}

func handleSession(w http.ResponseWriter, r *http.Request) {
    sessionClaims, ok := clerk.SessionFromContext(r.Context())
    if ok {
        // claims contain session information
    } else {
        // there is no active session (non-authenticated user)
    }
}
```

Here's the same HTTP server written for `v2`. Please note that only header based authentication with a bearer token is
supported by the `RequireHeaderAuthorization` middleware. Cookie based authentication is not supported.

```go
import (
    "net/http"

    "github.com/clerk/clerk-sdk-go/v2"
    clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

func main() {
    clerk.SetKey("sk_live_XXX")
    mux := http.NewServeMux()
    mux.Handle("/session", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handleSession)))
    http.ListenAndServe(":3000", mux)
}

func handleSession(w http.ResponseWriter, r *http.Request) {
    sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
    if ok {
        // claims contain session information
    } else {
        // there is no active session (non-authenticated user)
    }
}
```

If you're dealing with multiple Clerk API keys, you can pass a `jwks.Client` to the `RequireHeaderAuthorization` middleware.

```go
import (
    "net/http"

    "github.com/clerk/clerk-sdk-go/v2"
    clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
    "github.com/clerk/clerk-sdk-go/v2/jwks"
)

func main() {
    config := &clerk.ClientConfig{}
    config.Key = "sk_live_XXX"
    client := jwks.NewClient(config)
    mux := http.NewServeMux()
    withAuth := clerkhttp.RequireHeaderAuthorization(
        clerkhttp.JWKSClient(client),
    )
    mux.Handle("/session", withAuth(http.HandlerFunc(handleSession)))
    http.ListenAndServe(":3000", mux)
}

func handleSession(w http.ResponseWriter, r *http.Request) {
    sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
    if ok {
        // claims contain session information
    } else {
        // there is no active session (non-authenticated user)
    }
}
```

### Available middleware options

All available middleware options are preserved in the `v2` version of the library, but they have been renamed.

Name in v1 | Name in v2
-----------|------------
`WithAuthorizedParty` | `AuthorizedParty` and `AuthorizedPartyMatches`
`WithLeeway` | `Leeway`
`WithJWTVerificationKey` | `JSONWebKey`
`WithSatelliteDomain` | `Satellite`
`WithProxyURL` | `ProxyURL`
`WithCustomClaims` | `CustomClaimsConstructor`
n/a | `Clock`
n/a | `JWKSClient`

## Verify tokens

The `clerk.VerifyToken` method in version `v1` of the Clerk Go SDK has been renamed to `jwt.Verify` in `v2`.

The method accepts the same parameters, with two important differences.

- The JSON web key with which the token will be verified is a required parameter.
- The method will not cache the JSON web key.

In the `v1` version, the `clerk.VerifyToken` method would trigger an HTTP request to the Clerk Backend API to
fetch the JSON web key and would cache its value for one hour.

The new `jwt.Verify` method that is included in `v2` accepts the JSON web key as a required parameter. It is up
to the caller to get access to the key and pass it to `jwt.Verify`.

Please note that both HTTP middleware functions, `WithHeaderAuthorization` and `RequireHeaderAuthorization` will cache
the  JSON web key by default.

## Feedback and omissions

Please let us know about your experience upgrading to the `v2` version of the Clerk Go SDK.

You can reach us via [various support channels](https://clerk.com/support).

If you notice any bugs or omissions, feel free to open an [issue on Github](https://github.com/clerk/clerk-sdk-go/issues/new).
