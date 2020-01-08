# HTTP Check
Performs an HTTP GET request and returns a result based on the HTTP response code and if requested, an expected value or expression.

- `OK`: HTTP response code was not >= 300 and if requested, there was a match on the expected value or expression
- `WARNING`: HTTP response code was >= 300 and < 400
- `CRITICAL`: Connection failed or HTTP response code was >= 400 or on a failed match if an expected value or expression was supplied

## Options
- `--url` (`-u`): The URL to check. Required.
- `--insecure` (`-k`): Do not validate the server's certificate
- `--redirect` (`-r`): If set, follow redirects. Default is do not follow redirects.
- `--timeout` (`-t`): Timeout in seconds to wait for HTTP server response. Default is 15 seconds.
- `--path` (`-p`) and `--expression`: Used together. A json path and expression value to compare. Use this rather than `--path` and `--expectedValue` for making comparisons.
- `--path` (`-p`) and `--expectedValue` (`-e`): Used together. `--path` is the json path for retrieving a value and `--expectedValue` is the value to expected at the path. Use `--expression` instead for a more consistent interface.

## Examples
Check `example.com` for good response
```
check_http --url http://www.example.com
```

Check `example.com` for good response within 2 seconds
```
check_http --url http://www.example.com --timeout 2
```

## Using Expressions
Use expressions (`--expression`) for the ability to make comparisons other than simple string equality to a json field.

This option implements the [Gval expression evaluation package](https://github.com/PaesslerAG/gval) to compare against a json field, yielding the ability to use `==`, `!=`, `>`, `<`, `>=`, `<=` and probably more. These comparisons will work strings and numbers. The string or number comparison made depends on the use of double-quotes after the comparison operand.

Some examples. `== "comparetothistring"` gives a string comparison while `== 1337` gives a number comparison. Note the use of these quotes is important because `"1337" < "500"` (string comparison) will yield `true` while `1337 < 500` (number comparison) will yield `false`.

Consider using `--expression` rather than `--expectedValue`.

Some examples of how to use this `--expression` option...

The body returned for the example URL used below is

```
$ curl -s -H "Accept: application/json" https://icanhazdadjoke.com/j/HeaFdiyIJe | jq '.'
{
  "id": "HeaFdiyIJe",
  "joke": "What kind of magic do cows believe in? MOODOO.",
  "status": 200
}
```

For starters, here's how the `--expectedValue` and `--expression` options can be equally used.

```
check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path id --expectedValue HeaFdiyIJe

check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path id --expression '== "HeaFdiyIJe"'
```

Note the use of double-quotes above. That's because strings are being compared. If strings were not used, it would give:

```
check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path id --expression '== HeaFdiyIJe'
CheckHttp CRITICAL - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at id with value HeaFdiyIJe does not match expression "== HeaFdiyIJe"
```

More examples using other operators (continue to note the use of double-quotes for strings):

```
check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path id --expression '<= "IeaFdiyIJe"'
CheckHttp OK - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at id with value HeaFdiyIJe and expression "<= "IeaFdiyIJe"" yields true

check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path id --expression '>= "HeaFdiyIJd"'
CheckHttp OK - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at id with value HeaFdiyIJe and expression ">= "HeaFdiyIJd"" yields true

check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path id --expression '!= "notequal"'
CheckHttp OK - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at id with value HeaFdiyIJe and expression "!= "notequal"" yields true
```

And examples using number comparisons (double-quotes not used)

```
check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path status --expression '== 200'
CheckHttp OK - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at status with value 200 and expression "== 200" yields true

check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path status --expression '>= 200'
CheckHttp OK - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at status with value 200 and expression ">= 200" yields true

check_http --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path status --expression '< 400'
CheckHttp OK - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at status with value 200 and expression "< 400" yields true
```

Finally, an example using string comparison especially for Windows users at the Command Prompt shell. Notice the `--expression` option is contained in double-quotes and the double-quotes inside the option are given by using two double-quotes.

```
check_http.exe --url https://icanhazdadjoke.com/j/HeaFdiyIJe --format json --path id --expression "== ""HeaFdiyIJe"""
CheckHttp OK - Url https://icanhazdadjoke.com/j/HeaFdiyIJe responded with 200. The value found at id with value HeaFdiyIJe and expression "== "HeaFdiyIJe"" yields true
```
