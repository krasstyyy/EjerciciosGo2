# Level 1: Pure Functions

## Learning Objectives

- Understand the Red-Green-Refactor cycle by doing it
- Write table-driven tests in Go
- Structure tests using the AAA pattern (Arrange, Act, Assert)
- Use `t.Run` for subtests
- Use `t.Errorf` with clear failure messages

## Concepts

A **pure function** takes inputs and returns outputs. It has no side effects — it doesn't modify anything outside itself, it doesn't depend on external state. This makes pure functions the easiest thing to test.

This is where you start learning TDD. The domain is simple on purpose — focus on the rhythm, not the complexity.

### The Red-Green-Refactor cycle

1. **Red** — Write a test for behavior that doesn't exist yet. Run it. Watch it fail.
2. **Green** — Write the minimum code to make the test pass.
3. **Refactor** — Clean up. Rename. Simplify. Keep the tests green.

Repeat for every new behavior.

### Table-driven tests

Instead of writing one function per test case, Go developers use a table of inputs and expectations:

```go
tests := []struct {
    name string
    // inputs
    // expected outputs
}{
    {name: "scenario one", ...},
    {name: "scenario two", ...},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Arrange, Act, Assert
    })
}
```

Each row is an independent subtest. When one fails, the name tells you exactly which scenario broke.

## Exercise

Create `task.go` and `task_test.go` inside this folder. For every function, follow the TDD cycle: write a failing test first, make it pass, then refactor.

Use **table-driven tests** for all exercises.

### Step 1: FormatTitle

Write a function `FormatTitle(s string) string` that trims whitespace and capitalizes the first letter.

Think about these scenarios:

- Simple lowercase title
- Title with leading/trailing spaces
- Already capitalized title
- Empty string

### Step 2: IsValidTitle

Write a function `IsValidTitle(title string) bool` that returns whether a title is valid. A title is valid if it is non-empty after trimming whitespace.

Think about edge cases. What about a maximum length? You decide the rules — then write tests that enforce them.

### Step 3: FilterByStatus

First, define a `Task` struct and a `Status` type in `task.go`:

```go
type Status string

const (
    StatusTodo Status = "todo"
    StatusDone Status = "done"
)

type Task struct {
    Title  string
    Status Status
}
```

Write `FilterByStatus(tasks []Task, status Status) []Task` — it returns only the tasks matching the given status.

Use `go-cmp` to compare slices in your tests:

```go
import "github.com/google/go-cmp/cmp"

if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("mismatch (-want +got):\n%s", diff)
}
```

Scenarios to consider: empty slice, no matches, all match, mixed statuses.

### Step 4: CountByStatus

Write `CountByStatus(tasks []Task) map[Status]int` — it returns a count of tasks grouped by status.

By now you should be comfortable with the cycle. Design the test table, fail, pass, refactor.

## Hints

<details>
<summary>FormatTitle hint</summary>

Look at `strings.TrimSpace` and `unicode.ToUpper` for the first rune. Remember Go strings are UTF-8 encoded — be careful with byte vs rune indexing.

</details>

<details>
<summary>go-cmp setup</summary>

Run `go get github.com/google/go-cmp/cmp` from the repo root to add the dependency.

</details>

## Reflection

After completing this level, think about:

- How did it feel to write the test before the code?
- Did the test cases help you think about edge cases you would have missed?
- Did you find yourself wanting to write more code than the test required? That's the discipline — resist it.
- When did you refactor? What did you improve?
