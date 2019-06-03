# HTTP Check
Performs an HTTP GET request and returns a result based on the HTTP response code.

`OK`: Successful connection and HTTP response code was not >= 300

`WARNING`: Successful connection and HTTP response code was >= 300 and < 400

`CRITICAL`: Connectiion failed or HTTP response code was >= 400

## Options
- `--url (-u)`: The URL to check. Required.
- `--redirect (-r)`: If set, follow redirects. Default is do not follow redirects.
- `--timeout (-t)`: Timeout in seconds to wait for HTTP server response. Default is 15 seconds.

## Examples
Check `example.com` for good response
```
check_http --url http://www.example.com
```

Check `example.com` for good response within 2 seconds
```
check_http --url http://www.example.com --timeout 2
```
