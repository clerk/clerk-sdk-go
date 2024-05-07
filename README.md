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

Would you like to work on Open Source software and help maintain this repository? [Apply today!](https://jobs.ashbyhq.com/clerk)

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
users, err := client.Users().ListAll(clerk.ListAllUsersParams{})
```

The services exposed in the `clerk.Client` divide the API into logical chunks and
follow the same structure that can be found in the [Backend API documentation](https://clerk.com/docs/reference/backend-api).

For more examples on how to use the client, refer to the [examples](https://github.com/clerkinc/clerk-sdk-go/tree/main/examples/operations)

### Options

The SDK `Client` constructor can also accept additional options defined [here](https://github.com/clerk/clerk-sdk-go/blob/main/clerk/clerk_options.go).

A common use case is injecting your own [`http.Client` object](https://pkg.go.dev/net/http#Client) for testing or automatically retrying requests.
An example using [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp/#getting-a-stdlib-httpclient-with-retries) is shown below:

```go
retryClient := retryablehttp.NewClient()
retryClient.RetryMax = 5
standardClient := retryClient.StandardClient() // *http.Client

clerkSDKClient := clerk.NewClient(token, clerk.WithHTTPClient(standardClient))
```

## Middleware

The SDK provides the [`WithSessionV2`](https://pkg.go.dev/github.com/clerkinc/clerk-sdk-go/v2/clerk#WithSessionV2) middleware that injects the active session into the request's context.

The active session's claims can then be accessed using [`SessionFromContext`](https://pkg.go.dev/github.com/clerkinc/clerk-sdk-go/v2/clerk#SessionFromContext).

```go
mux := http.NewServeMux()
injectActiveSession := clerk.WithSessionV2(client)
mux.Handle("/your-endpoint", injectActiveSession(yourEndpointHandler))
```

Additionally, there's [`RequireSessionV2`](https://pkg.go.dev/github.com/clerkinc/clerk-sdk-go/v2/clerk#RequireSessionV2) that will halt the request and respond with 403 if the user is not authenticated. This can be used to restrict access to certain routes unless the user is authenticated.

For more info on how to use the middleware, refer to the
[example](https://github.com/clerkinc/clerk-sdk-go/tree/main/examples/middleware).

### Additional options

The middleware supports the following options:

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
