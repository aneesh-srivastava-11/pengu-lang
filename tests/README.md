# Pengu-Lang Tests

This directory contains the consolidated test suite for the `pengu-lang` project. 

## Running the Tests

To run the full test suite, simply use the standard Go test command from the root directory of the project:

```bash
go test ./tests/... -v
```

This will run all tests, including:
- **Compiler Layer:** Lexer, Parser, and Code Generation unit tests.
- **CLI Layer:** Integration tests verifying that the `pengu.exe` commands work.
- **E2E Layer:** End-to-End tests that compile an actual microservice, start it, and verify the HTTP endpoints receive responses successfully.

*Note: The E2E tests automatically run the built server on port `8080`. If you already have a server running on that port, the test will gracefully skip the HTTP checks rather than erroring out.*
