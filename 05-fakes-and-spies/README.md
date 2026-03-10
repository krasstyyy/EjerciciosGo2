# Level 5: Fakes and Spies

## Learning Objectives

- Understand the difference between stubs, fakes, and spies
- Build a **fake** — a working simplified implementation
- Build a **spy** — a test double that records interactions
- Test that your service **tells** collaborators the right things
- Know when to use each type of test double

## Concepts

### Stubs vs fakes vs spies

In Level 4, you used **stubs** — test doubles that return canned answers. They're simple but limited: they don't behave like real implementations, and they can't tell you *how* they were called.

This level introduces two more test doubles:

**Fake**: A working implementation that takes shortcuts. It has real logic — it stores and retrieves data — but uses a simple mechanism (an in-memory map instead of a database). Fakes are useful when you need realistic behavior without external dependencies.

**Spy**: A test double that records every call made to it. After the test runs, you inspect the recordings to verify that the right interactions happened. Spies answer the question: "Did my service tell this collaborator the right thing?"

### When to use each

- **Stub** — You need to control what a dependency *returns*. "When the repo returns these tasks, does the service filter correctly?"
- **Fake** — You need a dependency that *behaves realistically*. "Does the service correctly save and then retrieve a task?"
- **Spy** — You need to verify what the service *told* a collaborator. "When an overdue task is completed, did the service send a notification?"

### Notifications — a new collaborator

Your `TaskService` will now depend on a `Notifier` interface. When a task that's overdue gets completed, the service tells the notifier. The service doesn't ask — it tells.

## Exercise

Build on your types from previous levels. Create the files inside this folder.

### Step 1: InMemoryTaskRepository (fake)

Build a fake implementation of `TaskRepository` that stores tasks in a `map[string]*Task`.

This is a real, working implementation — it just uses memory instead of a database:

- `Save` stores the task
- `FindAll` returns all stored tasks

Write tests for the fake itself. Yes, test doubles can have their own tests when they contain real logic.

Think about:

- Saving a task and retrieving it
- Saving multiple tasks
- Overwriting a task with the same ID

### Step 2: Use the fake in service tests

Rewrite (or add) service tests using `InMemoryTaskRepository` instead of stubs.

Notice the difference: with a stub, you had to pre-configure what `FindAll` returns. With a fake, you call `Add` through the service, and the fake remembers. The tests become more natural — they exercise the full save-then-retrieve flow.

Think about when this is better and when a stub is still simpler.

### Step 3: Notifier interface

Define a `Notifier` interface:

```go
type Notifier interface {
    Notify(taskTitle string, message string) error
}
```

Update `TaskService` to accept a `Notifier` via its constructor.

### Step 4: Spy notifier

Build a spy that records calls:

```go
type spyNotifier struct {
    calls []notification
}

type notification struct {
    taskTitle string
    message   string
}
```

The spy implements `Notifier` by appending to its `calls` slice. After the test, you inspect `calls` to verify the right notifications were sent.

### Step 5: Notify on overdue completion

Update `TaskService.CompleteTask` so that when a task is overdue at the time of completion, the service tells the `Notifier`.

Write tests using the spy:

- Complete an overdue task — verify the spy received one notification with the correct title and message
- Complete a non-overdue task — verify the spy received no notifications
- Complete a task when the notifier returns an error — what should the service do?

### Step 6: Dummy notifier

For service tests that don't care about notifications (like `Add` or `ListOverdue`), you still need a `Notifier` to satisfy the constructor.

Create a **dummy** — the simplest possible implementation that does nothing:

```go
type dummyNotifier struct{}

func (d *dummyNotifier) Notify(string, string) error { return nil }
```

Use it in tests where notifications are irrelevant. This keeps those tests focused on what they're actually testing.

## Hints

<details>
<summary>Fake or stub?</summary>

If your test says "given these tasks exist, when I do X, then Y happens" — a stub works fine.

If your test says "when I add a task and then list all tasks, the added task appears" — you need a fake, because you're testing the save-then-retrieve flow.

</details>

<details>
<summary>Testing the spy</summary>

After calling the service method, inspect the spy's recorded calls:

```go
if len(spy.calls) != 1 {
    t.Fatalf("got %d notifications, want 1", len(spy.calls))
}
if spy.calls[0].taskTitle != "expected title" {
    t.Errorf("notified about %q, want %q", spy.calls[0].taskTitle, "expected title")
}
```

</details>

## Reflection

After completing this level, think about:

- How did the fake change the feel of your tests compared to stubs? Which felt more natural?
- The spy lets you verify *interactions*, not just *outcomes*. When is that important? When is it overkill?
- Did the dummy feel too simple to be a "real" concept? That's the point — it exists so you don't need to think about irrelevant dependencies.
- Look at your test file: can you tell which test double each test uses, and why? If not, consider renaming them to make intent clear.
