# Golang Tokens Library

** WARNING - WORKN IN PROGRESS - CONSIDER THIS ALPHA STAGE **

This is a library very similar to [tokens](https://github.com/zalando-stups/tokens) and [python-tokens](https://github.com/zalando-stups/python-tokens).

In a nutshell, you provide the OAuth2 token endpoint and which tokens and scopes to have managed.
 
The library will make sure that the managed tokens are always valid by refreshing them before they expire.

## Users Guide

The library will fetch credentials from JSON files (client.json and user.json) from the folder defined in the `CREDENTIALS_DIR` environment variable.

The threshold for refresh is around 60% of the expiration time.

## Example

	url := "https://example.com/oauth2/access_token"
	// You can manage multiple tokens with different scopes
	reqs := []tokens.ManagementRequest{
		tokens.NewRequest("test1", "password", "foo.read"),
		tokens.NewRequest("test2", "password", "bar.write"),
	}
	// Manager creation tries to obtain all tokens synchronously initially
	tokensManager, err := tokens.Manage(url, reqs)
	if err != nil {
		log.Fatal(err)
	}

	// You can use any of the above tokens for as long as you want
	for {
		accessToken, err := tokensManager.Get("test1") // or "test2"
		if err != nil {
			log.Println(err)
		} else {
			// Do something with accessToken
		}
	}
