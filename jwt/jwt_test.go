package jwt

import (
	"context"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/stretchr/testify/require"
)

func TestVerify_InvalidParams(t *testing.T) {
	ctx := context.Background()
	token := "eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yOWR6bUdmQ3JydzdSMDRaVFFZRDNKSTB5dkYiLCJ0eXAiOiJKV1QifQ.eyJhenAiOiJodHRwczovL2Rhc2hib2FyZC5wcm9kLmxjbGNsZXJrLmNvbSIsImV4cCI6MTcwNzMwMDMyMiwiaWF0IjoxNzA3MzAwMjYyLCJpc3MiOiJodHRwczovL2NsZXJrLnByb2QubGNsY2xlcmsuY29tIiwibmJmIjoxNzA3MzAwMjUyLCJvcmdzIjp7Im9yZ18ySUlwcVIxenFNeHJQQkhSazNzTDJOSnJUQkQiOiJvcmc6YWRtaW4iLCJvcmdfMllHMlNwd0IzWEJoNUo0ZXF5elFVb0dXMjVhIjoib3JnOmFkbWluIiwib3JnXzJhZzJ6bmgxWGFjTXI0dGRXYjZRbEZSQ2RuaiI6Im9yZzphZG1pbiIsIm9yZ18yYWlldHlXa3VFSEhaRmRSUTFvVjYzMnZWaFciOiJvcmc6YWRtaW4ifSwic2lkIjoic2Vzc18yYm84b2gyRnIyeTNueVoyRVZQYktBd2ZvaU0iLCJzdWIiOiJ1c2VyXzI5ZTBXTnp6M245V1Q5S001WlpJYTBVVjNDNyJ9.6GtQafMBYY3Ij3pKHOyBYKt76LoLeBC71QUY_ho3k5nb0FBSvV0upKFLPBvIXNuF7hH0FK2QqDcAmrhbzAI-2qF_Ynve8Xl4VZCRpbTuZI7uL-tVjCvMffEIH-BHtrZ-QcXhEmNFQNIPyZTu21242he7U6o4S8st_aLmukWQzj_4qir7o5_fmVhm7YkLa0gYG5SLjkr2czwem1VGFHEVEOrHjun-g6eMnDNMMMysIOkZFxeqiCnqpc4u1V7Z7jfoK0r_-Unp8mGGln5KWYMCQyp1l1SkGwugtxeWfSbE4eklKRmItGOdVftvTyG16kDGpzsb22AQGtg65Iygni4PHg"

	// Verifying without providing a key returns an error.
	_, err := Verify(ctx, &VerifyParams{
		Token: token,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing json web key")

	// Verify needs a key.
	_, err = Verify(ctx, &VerifyParams{
		Token: token,
		JWK:   &clerk.JSONWebKey{},
	})
	if err != nil {
		require.NotContains(t, err.Error(), "missing json web key")
	}

	// Try an invalid token.
	_, err = Verify(ctx, &VerifyParams{
		Token: "this-is-not-a-token",
		JWK:   &clerk.JSONWebKey{},
	})
	require.Error(t, err)
}
