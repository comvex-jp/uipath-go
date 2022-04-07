# Running Examples

Examples can be run by modifying the `examples.go` file and executing `examples.Run()` via `gore` or any Go REPL you are familiar with.

You need to modify the following items in `examples.go`
```
	credentials := uipath.Credentials{
		ClientID:   "{{ClientID}}",   // UIPath Client ID
		UserKey:    "{{UserKey}}",    // UIPath UserKey
		TenantName: "{{TenantName}}", // UIPath Tenant Name
	}
    ...
	c := uipath.Client{
		HttpClient:  httpClient,
		Credentials: credentials,
		URLEndpoint: "{{URLEndpoint}}", //UIPath URL endpoint
	}
```

## Running the example via `gore`
```
gore> :import github.com/comvex-jp/uipath-go/examples
gore> examples.Run()
{239722 Asset 39 true Global Text TestValue TestValue false 0    0 false  0 [] []} <nil>

```