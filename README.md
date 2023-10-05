<p align="center">
  <a href="https://www.clerk.com/?utm_source=github&utm_medium=starter_repos&utm_campaign=sdk_go" target="_blank" align="center">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="./docs/clerk-logo-dark.png">
      <img src="./docs/clerk-logo-light.png" height="64">
    </picture>
  </a>
  <br />
</p>

# Clerk Go SDK

Go client library for accessing the [Clerk Backend API](https://clerk.com/docs/reference/backend-api).

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/clerkinc/clerk-sdk-go/clerk)
[![Test Status](https://github.com/clerkinc/clerk-sdk-go/workflows/tests/badge.svg)](https://github.com/clerkinc/clerk-sdk-go/actions?query=workflow%3Atests)

[![chat on Discord](https://img.shields.io/discord/856971667393609759.svg?logo=discord)](https://discord.com/invite/b5rXHjAg7A)
[![documentation](https://img.shields.io/badge/documentation-clerk-green.svg)](https://clerk.com/docs)
[![twitter](https://img.shields.io/twitter/follow/ClerkDev?style=social)](https://twitter.com/intent/follow?screen_name=ClerkDev)

---

**Clerk is Hiring!**

Would you like to work on Open Source software and help maintain this repository? [Apply today!](https://apply.workable.com/clerk-dev/)

---

## Usage

First, add the Clerk SDK as a dependency to your project.

```
$ go get github.com/clerkinc/clerk-sdk-go
```

Add the following import to your Go files.

```go
import "github.com/clerkinc/clerk-sdk-go/clerk"
```

Now, you can create a Clerk client by calling the `clerk.NewClient` function.
This function requires your Clerk API key.
You can get this from the dashboard of your Clerk application.

Once you have a client, you can use the various services to access different parts of the API.

```go
apiKey := os.Getenv("CLERK_API_KEY")

client, err := clerk.NewClient(apiKey)
if err != nil {
    // handle error
}

// List all users for current application
users, err := client.Users().ListAll()
```

The services exposed in the `clerk.Client` divide the API into logical chunks and
follow the same structure that can be found in the [Backend API documentation](https://clerk.com/docs/reference/backend-api).

For more examples on how to use the client, refer to the [examples](https://github.com/clerkinc/clerk-sdk-go/tree/main/examples/operations)

## Middleware

In addition to the API operations, the SDK also provides a middleware that can be used to inject the active session into the request's context.
The Clerk middleware expects a `clerk.Client` and resolves the active session using the incoming session cookie.

The active session object will be added in the request's context using the key `clerk.ActiveSession`.

```go
mux := http.NewServeMux()
injectActiveSession := clerk.WithSession(client)
mux.Handle("/your-endpoint", injectActiveSession(yourEndpointHandler))
```

For a full example of how to use the middleware, refer to
[this](https://github.com/clerkinc/clerk-sdk-go/tree/main/examples/middleware).

### Auth v2

If you're using the newly-introduced [Auth v2](https://clerk.com/docs/upgrade-guides/auth-v2) scheme, you'll have to use the
`clerk.WithSessionV2()` middleware, instead of `clerk.WithSession()`.

Additionally, there's also `clerk.RequireSessionV2()` that will halt the request
and respond with 403 if the user is not authenticated.

Finally, to retrieve the authenticated session's claims you can use
`clerk.SessionFromContext()`.

### Additional options

The new middlewares (`clerk.WithSessionV2()` & `clerk.RequireSessionV2()`) also support the ability to pass some additional options.

- clerk.WithAuthorizedParty() to set the authorized parties to check against the azp claim of the token
- clerk.WithLeeway() to set a custom leeway that gives some extra time to the token to accommodate for clock skew
- clerk.WithJWTVerificationKey() to set the JWK to use for verifying tokens without the need to fetch or cache any JWKs at runtime
- clerk.WithCustomClaims() to pass a type (e.g. struct), which will be populated with the token claims based on json tags.
- clerk.WithSatelliteDomain() to skip the JWT token's "iss" claim verification.
- clerk.WithProxyURL() to verify the JWT token's "iss" claim against the proxy url.

For example

```golang
customClaims := myCustomClaimsStruct{}

clerk.WithSessionV2(
	clerkClient,
	clerk.WithAuthorizedParty("my-authorized-party"),
	clerk.WithLeeway(5 * time.Second),
	clerk.WithCustomClaims(&customClaims),
	clerk.WithSatelliteDomain(true),
	clerk.WithProxyURL("https://example.com/__clerk"),
	)
```

## License

This SDK is licensed under the MIT license found in the [LICENSE](./LICENSE) file.
