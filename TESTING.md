# Testing nagiosfoundation
Unit tests for code submitted to the project are encouraged. Every code submission is subject to unit tests, verified for a successful build and test coverage from the unit tests is submitted for review. The `Makefile` contains a couple of test targets to encourage the use of unit tests and make them easy to execute. Before submitting code, verify unit tests are successful and the code has adequate test coverage.

## Executing Go Tests
Go tests can be launched with the `Makefile` using the `test` target. The output is visual.
```
make test
```

## Code Coverage
Code coverage is launched with the `Makefile` by using the `coverage` target. Test output goes to `coverage.txt` and `coverage.html` files. View the html file using a browser to inspect the coverage.
```
make coverage
```