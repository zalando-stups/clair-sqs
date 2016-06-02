Go Tokens Library
=================

.. image:: https://travis-ci.org/zalando/go-tokens.svg?branch=master
    :target: https://travis-ci.org/zalando/go-tokens

.. image:: https://codecov.io/github/zalando/go-tokens/coverage.svg?branch=master
    :target: https://codecov.io/github/zalando/go-tokens?branch=master

.. image:: https://goreportcard.com/badge/github.com/zalando/go-tokens
    :target: https://goreportcard.com/report/github.com/zalando/go-tokens

.. image:: https://godoc.org/github.com/zalando/go-tokens?status.svg
    :target: https://godoc.org/github.com/zalando/go-tokens

This is a library very similar to `tokens`_ and `python-tokens`_.

In a nutshell, you provide the OAuth2 token endpoint and which tokens and scopes to be managed.
 
The library can use any custom credentials provider and it will make sure that the managed tokens are always
valid by refreshing them before they expire.

Credentials Provider
--------------------

The library currently contains implementations of the ``user.CredentialsProvider`` and the ``client.CredentialsProvider``
that you can use out of the box. The simplest providers just returns some static values used on creation. The other
providers fetch credentials from JSON files (user.json and client.json) from a folder defined in the
``CREDENTIALS_DIR`` environment variable.

You can easily use your own credential providers.

User Credentials
~~~~~~~~~~~~~~~~

User credentials are, very simply, a username and a password. Any type that implements the ``user.Credentials`` should
be able to provide them. The ``user.CredentialsProvider`` is the interface to implement for any custom type that is
able to provide ``user.Credentials``.

For a simple example, check the ``user/static.go`` file.

Client Credentials
~~~~~~~~~~~~~~~~~~

Client credentials are very similar to the user credentials. It consists of a client ID and client secret. Any type
that implements the ``client.Credentials`` should be able to provide them. The ``client.CredentialsProvider`` is the
interface to implement for any custom type that is able to provide ``client.Credentials``.

For a simple example, check the ``client/static.go`` file.

User Guide
----------

.. code-block:: go

    url := "https://example.com/oauth2/access_token"

    // You can manage multiple tokens with different scopes
    reqs := []tokens.ManagementRequest{
        tokens.NewPasswordRequest("test1", "foo.read"),
        tokens.NewPasswordRequest("test2", "user.email", "user.name"),
    }

    // Manager creation tries to obtain all tokens synchronously initially
    tokensManager, err := tokens.Manage(url, reqs)
    if err != nil {
        log.Fatal(err)
    }

    if test1, err := tokensManager.Get("test1"); err == nil {
        // Do something with the access token "test1"
    }

    if test2, err := tokensManager.Get("test2"); err == nil {
        // Do something with the access token "test2"
    }

This example would create a token manager using the JSON files credentials providers and a refresh threshold of 60% of the token validity time.

Manager Options
~~~~~~~~~~~~~~~

You can customize the behavior of the Manager with the following options:
    
RefreshPercentageThreshold(float64)
    Accepts any float between 0 and 1 (exclusive) which defines the percentage of token validity when to schedule background refreshing

WarningPercentageThreshold(float64)
    Accepts any float between 0 (exclusive) and 1 (inclusive) which defines when the library starts logging warnings that tokens will, eventually expire.
    This can happen if, for example, the background refresh is failing.
    It should be higher than the refresh threshold.
    
UserCredentialsProvider(user.CredentialsProvider)
    Accepts any user.CredentialsProvider instance that will provide user credentials for the password grant type
    
ClientCredentialsProvider(client.CredentialsProvider)
    Accepts any client.CredentialsProvider instance that will provide client credentials for the OAuth calls

.. _tokens: https://github.com/zalando-stups/tokens
.. _python-tokens: https://github.com/zalando-stups/python-tokens
