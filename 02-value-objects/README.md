# Level 2: Value Objects

## Learning Objectives

- Design structs with behavior (methods), not just data
- Use constructor functions to enforce invariants
- Keep struct fields unexported — expose behavior through methods
- Write test helper functions with `t.Helper()`
- Understand the difference between `t.Fatal` and `t.Error`

## Concepts

### Deep modules start with constructors

In Level 1 you used a `Task` struct with exported fields. Anyone could create an invalid task. Now you will design a `Task` that **protects its own invariants**.

A constructor function is the only way to create a valid `Task`. The struct fields are unexported — callers interact through methods:

```go
// The caller's view: simple interface
task, err := NewTask("Buy groceries", time.Now().Add(24*time.Hour))
task.Complete()
fmt.Println(task.IsCompleted()) // true
```

The caller doesn't know how completion is tracked internally. That's a deep module — simple surface, hidden complexity.

### Tell, Don't Ask

Instead of checking a task's status and deciding what to do, you **tell** the task what to do:

```go
// BAD: ask, then decide
if task.Status() == StatusTodo {
    task.SetStatus(StatusDone)
    task.SetCompletedAt(time.Now())
}

// GOOD: tell
task.Complete()
// The task handles its own state transitions.
```

### Test helpers

When multiple tests need the same setup, extract a helper:

```go
func newTestTask(t *testing.T) *Task {
    t.Helper()
    task, err := NewTask("Test task", time.Now().Add(24*time.Hour))
    if err != nil {
        t.Fatal(err)
    }
    return task
}
```

`t.Helper()` tells Go's test runner that this function is a helper. When an assertion fails inside it, the error is reported at the **caller's** line, not inside the helper. This makes failures easy to locate.

Use `t.Fatal` (not `t.Error`) in helpers — if creating a test task fails, continuing the test makes no sense.

## Exercise

Create a new `Task` type with unexported fields and behavior. Create `task.go` and `task_test.go` in this folder.

### Step 1: NewTask constructor

```go
func NewTask(title string, dueDate time.Time) (*Task, error)
```

Rules:

- Title must be valid (non-empty after trimming). Return an error if not.
- The task starts with status "todo".
- `CreatedAt` is set to `time.Now()`.
- An ID is generated internally (use a UUID or a simple counter — your choice).

Write tests for:

- Valid inputs produce a task with correct title and status
- Empty title returns an error
- Whitespace-only title returns an error
- Due date in the past — is this an error? You decide. Document your decision in a test.

### Step 2: Complete

```go
func (t *Task) Complete()
```

Completing a task changes its status and records when it was completed.

Write tests for:

- A new task is not completed
- After calling `Complete()`, the task is completed
- Completing an already completed task — should this panic, be a no-op, or return an error? Pick one. Test it.

You will need accessor methods to verify state in tests: `IsCompleted() bool`, `Title() string`, `Status() Status`.

### Step 3: Postpone

```go
func (t *Task) Postpone(days int) error
```

Postponing extends the due date by the given number of days.

Rules:

- Days must be positive. Return an error otherwise.
- A completed task cannot be postponed. Return an error.

### Step 4: Prioritize

Define a `Priority` type:

```go
type Priority int

const (
    PriorityLow    Priority = iota
    PriorityMedium
    PriorityHigh
)
```

```go
func (t *Task) Prioritize(p Priority) error
```

Rules:

- The priority must be a valid value. Return an error for invalid values.
- A new task starts with `PriorityMedium` by default.

### Step 5: IsOverdue

```go
func (t *Task) IsOverdue(now time.Time) bool
```

A task is overdue if it is not completed and its due date is before `now`.

Inject `now` as a parameter instead of calling `time.Now()` inside the method. This makes the function deterministic and easy to test — you control time.

### Step 6: String

```go
func (t *Task) String() string
```

The task controls its own representation. Format it however you think is clear and useful.

## Hints

<details>
<summary>How to generate an ID</summary>

For simplicity, you can use a package-level counter with `sync/atomic`, or `crypto/rand` for a random string. Don't over-engineer this — the point is that the caller never provides an ID.

</details>

<details>
<summary>Testing time</summary>

Pass `time.Now()` for "now" in tests, or create fixed times:

```go
now := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
dueDate := now.Add(48 * time.Hour)
```

This gives you deterministic, reproducible tests.

</details>

## Reflection

After completing this level, think about:

- How does hiding fields behind methods change the way you write tests?
- Did you feel the urge to export fields "just for testing"? That's a design smell — if you need to expose internal state to test it, the API might be missing a method.
- How does the constructor make invalid states harder to create?
- Did passing `now` as a parameter feel awkward? That pattern (dependency injection for time) is fundamental to testable design.
