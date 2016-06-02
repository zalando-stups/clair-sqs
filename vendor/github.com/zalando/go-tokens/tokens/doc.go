/*
Package tokens implements a thread safe token manager. Your application starts by creating a
manager with requests for management of named tokens. Those tokens can later be obtained using
their ID. The manager will refresh the tokens in the background when needed.

Usage:

Create a new token manager with the Manage() function

	tokenManager := Manage("http://oauth-endpoint", mgmtRequests)

This creates a new Manager which will obtain OAuth tokens from the http://oauth-endpoint

You can set some options for the manager with the extra variadic argument. The available options are:

	RefreshPercentageThreshold(float64)
	WarningPercentageThreshold(float64)
	UserCredentialsProvider(user.CredentialsProvider)
    ClientCredentialsProvider(client.CredentialsProvider)

The warning threshold should be higher than the refresh threshold. The default values for
these options are 60% (0.60) for the refresh threshold and 80% (0.80) for the warning threshold.
The default credentials providers, for both user and client, read those credentials from JSON files

They can be used when the token manager is created

	tokenManager := Manage(
		"http://oauth-endpoint",
		mgmtRequests,
		RefreshPercentageThreshold(0.90),
		WarningPercentageThreshold(0.95),
	)

Or they can also be used at a later stage. The refresh threshold  will only affect the next scheduling
while the warning threshold takes effect immediately. Usually you don't need to change any of them at
runtime.

You application should only need to get the named tokens for whichever purpose, usually, to access
OAuth2 protected endpoints using the tokens in the Authorization header.

The Get() function returns the AccessToken with a specific name (from the ManagementRequest)

The call can fail with 2 specific errors:

	ErrTokenNotAvailable
	ErrTokenExpired

*/
package tokens
