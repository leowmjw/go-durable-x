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

- <<ADD HERE>>