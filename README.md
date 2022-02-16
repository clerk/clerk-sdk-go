# clerk-sdk-go #

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/clerkinc/clerk-sdk-go/clerk)
[![Test Status](https://github.com/clerkinc/clerk-sdk-go/workflows/tests/badge.svg)](https://github.com/clerkinc/clerk-sdk-go/actions?query=workflow%3Atests)

Go client library for accessing the [Clerk Backend API v1](https://docs.clerk.dev/reference/backend-api-reference).

---

**Clerk is Hiring!**

Would you like to work on Open Source software and help maintain this repository? Apply today https://apply.workable.com/clerk-dev/.

---

## Usage ##

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
follow the same structure that can be found in the [Backend API documentation](https://docs.clerk.dev/backend/backend-api-reference).

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

If you're using the newly-introduced [Auth v2](https://docs.clerk.dev/main-concepts/auth-v2) scheme, you'll have to use the 
`clerk.WithSessionV2()` middleware, instead of `clerk.WithSession()`.

Additionally, there's also `clerk.RequireSessionV2()` that will halt the request 
and respond with 403 if the user is not authenticated.

Finally, to retrieve the authenticated session's claims you can use 
`clerk.SessionFromContext()`.

## License ##

This SDK is licensed under the MIT license found in the [LICENSE](./LICENSE) file.
