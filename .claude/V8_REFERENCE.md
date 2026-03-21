# v8 SDK Reference

## Purpose

This file documents patterns from the v8 SDK (gtoggl-api) that worked well
and should be maintained or improved in v9.

## Where is the v8 SDK?

GitHub: <https://github.com/shoekstra/gtoggl-api>

## What to Look At

### Authentication Pattern

- How does v8 handle Bearer tokens?
- Does it validate tokens on client creation?
- How are credentials managed?

### Error Handling

- What custom error types does v8 use?
- How are API errors distinguished from network errors?
- Does it retry on certain errors?

### HTTP Layer

- What HTTP library does v8 use (stdlib or third-party)?
- How does it handle timeouts?
- Are there retry mechanisms?

### Request/Response Patterns

- How does v8 structure request bodies?
- How does it parse responses?
- Are there any quirks in the API that need special handling?

### Services Currently Implemented

- Which endpoints are implemented in v8?
- Are there any that need special handling?
- Which methods have proven useful?

### Known Issues

- Are there any bugs in v8 that v9 should fix?
- Any API endpoints that are unreliable?
- Performance bottlenecks?

## How v9 Improves on v8

The v9 SDK maintains good patterns from v8 but improves:

1. **Better organization**: Service-oriented design vs monolithic
2. **Clearer testing**: Table-driven tests with mocks instead of ???
3. **Modern patterns**: Context throughout, functional options
4. **Type safety**: Stricter typing for better compile-time safety
5. **Documentation**: Better godoc comments and examples

## Migration Notes

When implementing v9 services:

1. **Reference v8 for method names** - Keep familiar API
2. **Check v8 for gotchas** - Avoid known issues
3. **Improve error handling** - More explicit
4. **Better tests** - More comprehensive coverage
5. **Clear examples** - Show how to use each method
