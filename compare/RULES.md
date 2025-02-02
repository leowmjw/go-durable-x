# RULES

## CODE STYLE

- Use Latest Golang features >= v1.23
- Use slog for structured logging; one line per request
- Use only stdlib servermux and use latest http func handler style
- All common shared types will be located under the temporal/types folder
- Ensure there is always unit test for major functionality and coverage at least 70%
- Allowed to use the httpz library for routing and middleware
- Use native testsuite each for Temporal, restate + inngest
- If mocking out external systems; prefer to use function injection method instead of mocking interfaces
- External system events; whether a human approval or 
- Be as thorough to cover all possible edge cases

## LEARNINGS

1. Project Structure
   - Keep module structure flat and simple (temporal, restate, inngest)
   - Each implementation should be self-contained with its own go.mod
   - Shared types should be in a types/ directory within each implementation
   - Activity implementations should be in activities/ directory

2. Temporal Workflow Design
   - Start with basic happy path and simple compensation flows
   - Implement activities as standalone functions for better testability
   - Use structured logging (slog) consistently across activities
   - Keep retry policies configurable at workflow level
   - Register activities explicitly in main() rather than passing interfaces

3. Testing Strategy
   - Use Temporal's testsuite for workflow testing
   - Mock activities to test different failure scenarios
   - Test compensation flows thoroughly
   - Ensure proper cleanup in failure scenarios
   - Use testify for assertions and mocking

4. Code Organization
   - Keep activity implementations simple and focused
   - Use clear naming conventions for activities (Book*, Cancel*)
   - Maintain consistent error handling patterns
   - Use interfaces for external dependencies
   - Follow Go 1.23+ idioms and features

5. Error Handling
   - Always implement compensation logic for failures
   - Log errors with proper context using structured logging
   - Use meaningful error messages that help debugging
   - Consider retry policies carefully for each activity

6. Module Management
   - Keep go.mod dependencies minimal and explicit
   - Use specific versions for stability
   - Maintain module path consistency
   - Handle indirect dependencies properly

7. Development Process
   - Build and test incrementally
   - Start with core functionality before advanced features
   - Document implementation status and pending work
   - Track technical debt and future improvements

8. Best Practices Identified
   - Use function injection for mocking instead of interfaces
   - Keep activities stateless
   - Use constants for configuration
   - Implement proper cleanup for failed operations
   - Follow standard Go project layout