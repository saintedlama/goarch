# ArchScout Agents Instructions

## Keep the README up to date

After each feature implemented, make sure to update the README.md file with relevant information about the new feature. This includes:

- A brief description of the feature.
- Code snippets demonstrating how to use the feature.
- Any important notes or caveats related to the feature.

## Testing

- For every new feature implemented, write comprehensive tests to ensure the feature works as expected.
- Tests are in *_test.go files alongside source code
- Archictecture tests use `arch_` prefix such as `arch_test.go`.
- Use `testify` for assertions.

## Development Commands

```bash
# Run go vet to check for any issues in the code
make vet

# Run staticcheck to perform static analysis and catch potential issues
make lint

# Run go fmt to format the code according to Go standards
make fmt

# Build the project to ensure all code is compiled and ready for testing
make build

# Run tests
make test
```
