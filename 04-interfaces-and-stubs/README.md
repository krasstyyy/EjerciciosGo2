# Level 4: Interfaces and Stubs

## Learning Objectives

- Define Go interfaces based on what the consumer needs
- Use **dependency injection** via constructor parameters
- Write hand-rolled **stubs** to control dependency behavior in tests
- Test a service layer in isolation from its dependencies
- Understand the difference between a stub and a real implementation

## Concepts

### Why interfaces?

In Level 3, `TaskList` managed tasks in memory using a slice or map. That works, but what if you want to store tasks in a database? Or read them from a file? Or switch storage without changing the service logic?

**Interfaces** decouple what something does from how it does it.

### Define interfaces where they're used

In Go, interfaces are defined by the **consumer**, not the implementer. Your `TaskService` needs something that can save and retrieve tasks — so it defines what it needs:

```go
type TaskRepository interface {
    Save(*Task) error
    FindAll() ([]*Task, error)
}
```

This is a narrow, deep interface: two methods that handle significant work. The service doesn't care how storage works — it just calls `Save` and `FindAll`.

### Dependency injection

The service receives its dependencies through the constructor:

```go
func NewTaskService(repo TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}
```

In production, you pass a real database repository. In tests, you pass a stub.

### Stubs

A **stub** returns canned answers. It has no real logic — it just gives back what you told it to:

```go
type stubRepo struct {
    tasks []*Task
    err   error
}

func (s *stubRepo) Save(*Task) error           { return s.err }
func (s *stubRepo) FindAll() ([]*Task, error)   { return s.tasks, s.err }
```

Stubs let you control the test scenario precisely: "what should the service do when the repository returns these tasks?" or "what should the service do when the repository returns an error?"

## Exercise

Create `service.go`, `service_test.go`, and `repository.go` inside this folder. You can copy your `Task` type from previous levels.

### Step 1: Define TaskRepository

Define the `TaskRepository` interface in `repository.go`. Keep it narrow — only the methods the service actually needs:

- `Save(*Task) error`
- `FindAll() ([]*Task, error)`

### Step 2: Build TaskService with Add

Create a `TaskService` struct that receives a `TaskRepository` via its constructor.

Implement `Add(title string, dueDate time.Time) error` — the service creates a task and tells the repository to save it.

Write a test using a stub repository. Verify that `Add` succeeds when the repo succeeds, and returns an error when the repo fails.

### Step 3: CompleteTask

Implement `CompleteTask(title string) error`.

The service should:

1. Fetch all tasks from the repository
2. Find the one matching the title
3. Complete it
4. Save it back

Write tests with a stub that returns different task lists. Think about:

- Task found and completed successfully
- Task not found
- Repository returns an error

### Step 4: ListOverdue

Implement `ListOverdue() ([]*Task, error)`.

The service fetches all tasks and filters for overdue ones. The filtering logic lives in the service (or delegates to `Task.IsOverdue`).

Stub the repository with a mix of overdue and non-overdue tasks. Verify the service returns only the right ones.

### Step 5: Error propagation

Go back through your tests. For each service method, make sure you have a test where the stub returns an error. Verify the service propagates it correctly — doesn't swallow it, doesn't panic.

## Hints

<details>
<summary>Stub design</summary>

Keep your stubs simple. A stub for `TaskRepository` needs just two fields: the data to return and the error to return. Configure them per test case.

</details>

<details>
<summary>Testing error paths</summary>

Create a stub that returns an error:

```go
repo := &stubRepo{err: errors.New("connection failed")}
```

Then verify that the service method returns an error wrapping or matching it.

</details>

## Reflection

After completing this level, think about:

- How did the interface change the way you think about the service? The service no longer knows about storage details.
- Were the stubs easy to write? In Go, implementing a two-method interface is trivial — no framework needed.
- Look at your tests: do they test the service's behavior, or the repository's behavior? Each should be tested independently.
- Could you write a real `PostgresTaskRepository` that satisfies the same interface? The service wouldn't change at all. That's the power of interfaces.
