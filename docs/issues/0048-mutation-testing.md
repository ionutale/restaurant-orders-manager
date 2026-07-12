# 0048 — Mutation testing: Go backend code

## What to build

Add mutation testing using go-mutesting (or similar) for the backend Go code. Mutation testing verifies that tests actually catch bugs by introducing code mutations and checking if tests fail. This measures test quality, not just coverage.

## Acceptance criteria

- [ ] Install and configure go-mutesting
- [ ] Run mutation tests on `internal/handler/` package
- [ ] Run mutation tests on `internal/auth/` package
- [ ] Run mutation tests on `internal/domain/` package
- [ ] Achieve mutation score > 70% (at least 70% of mutations are killed)
- [ ] Document which mutations survived (uncovered code paths)
- [ ] Add missing tests for each surviving mutation
- [ ] Integrate mutation testing into CI (optional, may be slow)
- [ ] Add `make mutate` or `go test -mutate` target

## Blocked by

- 0031 — Unit tests: Go backend handlers
- 0032 — Unit tests: Go auth middleware and JWT
