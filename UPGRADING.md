# Clerk Go SDK upgrade guide

The Clerk Go SDK uses [Semantic Versioning](https://semver.org/)
for tracking all changes to the library.

Semantic versions take the form of `MAJOR.MINOR.PATCH`. Incrementing

- MAJOR version introduces incompatible API changes
- MINOR version introduces functionality in a backwards compatible manner
- PATCH version includes backwards compatible bug fixes

This file contains upgrade notes for all `MAJOR` version bumps.

While there are no plans to sunset older versions of the Go SDK,
the current `MAJOR` version of `clerk-sdk-go` is actively supported.

Bug fixes and security updates may be added to previous versions, but
there is no guarantee.

## 1.x.x to 2.x.x

The `v2` version of the Clerk Go SDK is a complete rewrite and introduces
a lot of breaking changes.

### Minimum Go version

The minimum supported Go version for the `v2` version of the Clerk Go SDK is `1.19`.

### Setting an API key

```diff
- client, err := clerk.NewClient("sk_live_XXX")
+ clerk.SetKey("sk_live_XXX")
```

### Invoking API operations

```diff
- client.$Resource$().Create(clerk.Create$Resource$Params{})
+ $resource$.Create(ctx, $resource$.CreateParams{})
```

API operations in `v1` of the Clerk Go SDK are organized by service. There are
different services for each API and all services are properties of a single
`clerk.Client`.

In `v2` API operations are grouped by API resource. Every API resource is defined in its own package.

```diff
// Create an organization, in v1 and v2. Error handling is omitted.
- client, err := clerk.NewClient("sk_live_XXX")
- org, err := client.Organizations().Create(clerk.CreateOrganizationParams{
-     Name: "Acme Inc",
- })
+ ctx := context.Background()
+ clerk.SetKey("sk_live_XXX")
+ org, err := organization.Create(ctx, &organization.CreateParams{
+     Name: clerk.String("Acme Inc"),
+ })
```

### Support for context.Context

All API operations in `v2` receive a `context.Context` as their first argument.

### List operation responses

In `v2` of the Clerk Go SDK, the return types of list operations are similar. They always contain the total count of resources
available in the server, and a slice with the operation results.

```diff
// Fetch a list of 10 users. Error handling is omitted.
- limit := 10
- users, err := client.Users().ListAll(clerk.ListAllUserParams{
-    Limit: &limit,
- })
- if len(users) > 0 {
-    fmt.Println(users[0].ID)
- }
+ params := &user.ListParams{}
+ params.Limit = clerk.Int64(10)
+ list, err := user.List(context.Background(), &params)
+ if list.TotalLength > 0 {
+     fmt.Println(list.Users[0].ID)
+ }
```

### Every field in API operation parameters is a pointer

The `v2` version of the library introduces helper functions to generate pointers from various type values. These are:
- `clerk.String`
- `clerk.Bool`
- `clerk.Int64`
- `clerk.JSONRawMessage`

Using the helpers above, here's how you can invoke an API operation with a `*Params` struct.

```go
domain.Create(context.Background(), &domain.CreateParams{
    Name: clerk.String("clerk.com"),
    IsSatellite: clerk.Bool(true),
})
```
You can explicitly pass zero values with `clerk.String("")` or `clerk.Int64(0)`.

### The `clerk.ErrorResponse` type changed to `clerk.APIErrorResponse`

```diff
- clerk.ErrorResponse
+ clerk.APIErrorResponse
```

The `v2` version of the library introduces a new type for API responses that contain errors.
The new type is [clerk.APIErrorResponse](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#APIErrorResponse) and it replaces `clerk.ErrorResponse`.

The `clerk.APIErrorResponse` type contains additional fields and provides access to more debugging information.

```diff
- org, err := client.Organizations().Create(clerk.CreateOrganizationParams{
-     Name: "Acme Inc",
- })
- if err != nil {
-     if errResp, ok := err.(*clerk.ErrorResponse); ok {
-         // Access the API errors
-         errResp.Errors
-     }
- }
+ ctx := context.Background()
+ clerk.SetKey("sk_live_XXX")
+ org, err := organization.Create(ctx, &organization.CreateParams{
+     Name: clerk.String("Acme Inc"),
+ })
+ if err != nil {
+     if apiErr, ok := err.(*clerk.APIErrorResponse); ok {
+         // Access the API errors and additional information
+         apiErr.TraceID
+         apiErr.Error()
+         apiErr.Response.RawJSON
+     }
+ }
```

### HTTP middleware

```diff
- clerk.WithSessionV2
+ http.WithHeaderAuthorization

- clerk.RequireSessionV2
+ http.RequireHeaderAuthorization

- clerk.SessionFromContext
+ clerk.SessionClaimsFromContext
```

The `clerk.WithSessionV2` and `clerk.RequireSessionV2` middleware functions from `v1` are replaced by [http.WithHeaderAuthorization](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http#WithHeaderAuthorization) and [http.RequireHeaderAuthorization](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http#RequireHeaderAuthorization) in `v2`

Please note that as the name implies `WithHeaderAuthorization` and `RequireHeaderAuthorization` support only authentication with a bearer token.
The token needs to be provided in the "Authorization" request header.

> [! IMPORTANT]
> Cookie based authentication is not supported at all by the `v2` version of the library.

To get access to the active session claims from the http.Request context, you must replace `clerk.SessionFromContext` with [clerk.SessionClaimsFromContext](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2#SessionClaimsFromContext).

```diff
// Protect a route with Clerk authentication middleware.
// Error handling is omitted.
mux := http.NewServeMux()
- client, err := clerk.NewClient("sk_live_XXX")
- mux.Handle("/session", clerk.RequireSessionV2(client)(http.HandlerFunc(handleSession)))
+ clerk.SetKey("sk_live_XXX")
+ mux.Handle("/session", clerkhttp.RequireHeaderAuthorization()(http.HandlerFunc(handleSession)))
http.ListenAndServe(":3000", mux)

func handleSession(w http.ResponseWriter, r *http.Request) {
-    sessionClaims, ok := clerk.SessionFromContext(r.Context())
+    sessionClaims, ok := clerk.SessionClaimsFromContext(r.Context())
    if ok {
        // claims contain session information
    } else {
        // there is no active session (non-authenticated user)
    }
}
```

#### Available middleware options

All available middleware options are preserved in the `v2` version of the library, but they have been renamed.

```diff
- WithAuthorizedParty(...string)
+ AuthorizedPartyMatches(...string)

- WithLeeway(time.Duration)
+ Leeway(time.Duration)

- WithJWTVerificationKey(string)
+ JSONWebKey(string)

- WithSatelliteDomain(string)
+ Satellite(string)

- WithProxyURL(string)
+ ProxyURL(string)

- WithCustomClaims(interface{})
+ CustomClaimsConstructor(func(context.Context) any)
```

The `v2` version of the Clerk Go SDK provides additional middleware options.

- [AuthorizedParty(func(string) bool)](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http#AuthorizedParty)
- [JWKSClient(\*jwks.Client)](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/http#JWKSClient)

### Verify tokens

The `clerk.VerifyToken` method in version `v1` of the Clerk Go SDK has been renamed to [jwt.Verify](https://pkg.go.dev/github.com/clerk/clerk-sdk-go/v2/jwt#Verify) in `v2`.

The method accepts the same parameters, with two important differences.

- The JSON web key with which the token will be verified is a required parameter.
- The method will not cache the JSON web key.
