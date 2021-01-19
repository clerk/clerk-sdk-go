# clerk-sdk-go #

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/clerkinc/clerk-sdk-go/clerk)
[![Test Status](https://github.com/clerkinc/clerk-sdk-go/workflows/tests/badge.svg)](https://github.com/clerkinc/clerk-sdk-go/actions?query=workflow%3Atests)

Go client library for accessing the [Clerk Server API v1](https://docs.clerk.dev/server-api/).

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

client := clerk.NewClient(apiKey)

// List all users for current application
users, err := client.Users().ListAll()
```

The services exposed in the `clerk.Client` divide the API into logical chunks and 
follow the same structure that can be found in the [server-side API documentation](https://docs.clerk.dev/server-api/).

For more examples on how to use the client, refer to the [example](https://github.com/clerkinc/clerk-sdk-go/tree/main/example)

## License ##

This SDK is licensed under the MIT license found in the [LICENSE](./LICENSE) file.
