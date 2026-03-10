# Level 3: Inside-Out TDD

## Learning Objectives

- Practice **inside-out TDD**: build from the core, compose outward
- Design a type that manages a collection of value objects
- Apply **tell, don't ask**: callers tell the `TaskList` what to do
- Experience how tested units compose naturally

## Concepts

### Inside-out TDD

In Levels 1 and 2, you built and tested `Task` in isolation. Now you'll build a `TaskList` that manages tasks. This is inside-out TDD in action:

1. **Inner layer** (done): `Task` — validated, tested, trustworthy
2. **Outer layer** (this level): `TaskList` — uses `Task`, adds collection behavior

You build outward from a solid foundation. Each layer trusts the one below because it's already tested.

### Tell, don't ask — at the collection level

The `TaskList` should encapsulate finding and acting on tasks. The caller doesn't fetch a task, modify it, and put it back. The caller says what it wants done:

```go
// Tell: the list handles finding and completing
list.Complete("Buy groceries")

// Don't ask: the caller shouldn't need to do this
task := list.Find("Buy groceries")
task.Complete()
```

The `TaskList` is a **deep module**: few methods, each doing significant work internally.

### When you don't need test doubles

Inside-out TDD uses real collaborators. Your `TaskList` tests will create real `Task` values. No stubs, no mocks — the inner types are fast, deterministic, and already tested.

This is the strength of inside-out: each test exercises real behavior end to end within its layer.

## Exercise

Create `tasklist.go` and `tasklist_test.go` inside this folder. You can copy your `Task` type from Level 2 or import it if you prefer.

Follow TDD for every method. Build the `TaskList` incrementally — one method at a time.

### Step 1: Add

Design a `TaskList` type with an `Add(task *Task) error` method.

Think about:

- What happens when adding a task with a duplicate title?
- Should there be a capacity limit?

### Step 2: Complete

Add a `Complete(title string) error` method. The list finds the task and completes it in one call.

Think about:

- What if the title doesn't exist?
- What if the task is already completed?

### Step 3: Postpone

Add a `Postpone(title string, days int) error` method. Same pattern — the list handles finding and delegating.

### Step 4: Pending

Add a `Pending() []*Task` method that returns only tasks that are not completed.

Use `go-cmp` to compare the results. Think about the order — should it be deterministic?

### Step 5: Overdue

Add an `Overdue(now time.Time) []*Task` method that returns tasks past their due date.

Again, accept `now` as a parameter for testability.

### Step 6: Reflect on the design

Look at your `TaskList` API:

- `Add`, `Complete`, `Postpone` — these are **commands** (tell)
- `Pending`, `Overdue` — these are **queries** (read)

Is the interface simple? Does each method do significant work? If so, you have a deep module.

## Hints

<details>
<summary>Storing tasks</summary>

A slice `[]*Task` works fine. If you need fast lookup by title, consider a `map[string]*Task`. Let the tests guide you — start simple, refactor when needed.

</details>

<details>
<summary>Error design</summary>

Define sentinel errors for common failures:

```go
var (
    ErrTaskNotFound      = errors.New("task not found")
    ErrDuplicateTask     = errors.New("task already exists")
)
```

Return them from methods so tests can check with `errors.Is`.

</details>

## Reflection

After completing this level, think about:

- You built `TaskList` on top of a tested `Task`. Did you ever doubt whether `Task.Complete()` worked? That's the confidence inside-out gives you.
- Did you feel tempted to test `Task` behavior again from within `TaskList` tests? That's redundant — trust the inner layer.
- Look at your `TaskList` interface from the outside. Is it easy to use? Would a caller need to understand the internals?
- Could you swap the internal storage (slice to map) without changing any tests? If yes, your tests are testing behavior, not implementation.
