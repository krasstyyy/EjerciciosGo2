# Go TDD Kata

A self-paced coding dojo for learning **Test-Driven Development** in Go. You will build a **Todo List / Task Manager** from scratch, progressing from pure functions to full service layers with test doubles.

Each level introduces new TDD concepts. You write all the code — both tests and production code — yourself.

## Table of Contents

- [What is TDD?](#what-is-tdd)
- [TDD is Not a Silver Bullet](#tdd-is-not-a-silver-bullet)
- [TDD Approaches](#tdd-approaches)
- [Anatomy of a Go Test](#anatomy-of-a-go-test)
- [Test Doubles](#test-doubles)
- [The Test Pyramid](#the-test-pyramid)
- [Design Principles](#design-principles)
- [How to Use This Repo](#how-to-use-this-repo)
- [Levels](#levels)

---

## What is TDD?

**Test-Driven Development (TDD)** is a software development discipline where you write a failing test before writing the production code that makes it pass.

The cycle has three steps:

1. **Red** — Write a test that describes the behavior you want. Run it. It fails.
2. **Green** — Write the simplest production code that makes the test pass. Nothing more.
3. **Refactor** — Clean up both test and production code while keeping all tests green.

Then repeat. Every feature, every behavior, every edge case goes through this cycle.

```text
    ┌───────────┐
    │   RED     │  Write a failing test
    │  (fail)   │
    └─────┬─────┘
          │
          v
    ┌───────────┐
    │  GREEN    │  Write the minimum code
    │  (pass)   │  to make the test pass
    └─────┬─────┘
          │
          v
    ┌───────────┐
    │ REFACTOR  │  Clean up, keeping
    │  (clean)  │  all tests green
    └─────┬─────┘
          │
          └──────> Back to RED
```

### Why TDD?

- **Design feedback** — If something is hard to test, the design is telling you something. TDD forces you to think about interfaces before implementation.
- **Confidence** — A comprehensive test suite lets you refactor fearlessly.
- **Documentation** — Tests describe what the code does. They are living documentation that never goes stale.
- **Small steps** — TDD prevents you from writing code you don't need (YAGNI). You only write code that a test demands.

### The Three Laws (Robert C. Martin)

1. You must not write production code unless it is to make a failing test pass.
2. You must not write more of a test than is sufficient to fail (and not compiling counts as failing).
3. You must not write more production code than is sufficient to pass the currently failing test.

---

## TDD is Not a Silver Bullet

TDD is a powerful discipline, but it requires practice to apply effectively. Your first attempts will feel slow and awkward — that's normal. The productivity gains come after you internalize the rhythm.

### Common criticisms

- **Slower initial development** — Writing tests first takes more time upfront. Studies show TDD can increase development time by 15-35%. The trade-off is significantly fewer defects in production.
- **False sense of security** — Passing tests don't guarantee correct software. Tests are only as good as the scenarios you think to cover.
- **Test-induced design damage** — Over-testing or testing the wrong things can lead to code shaped around testability rather than clarity. This is why principles like deep modules and tell-don't-ask matter.
- **Not always the right tool** — Exploratory prototyping, UI layout work, or trivial glue code may not benefit from TDD. Knowing when *not* to TDD is part of the skill.

### What the research says

Studies (Nagappan et al., Microsoft/IBM 2008) found that teams using TDD produced code with 40-90% fewer defects, at the cost of 15-35% longer development time. Whether that trade-off is worth it depends on your context — a medical device and a weekend hackathon have different quality requirements.

### TDD in the AI era

With AI coding assistants generating code at unprecedented speed, you might wonder: is TDD still relevant?

It's **more relevant than ever**.

**A real-world example:** A single engineer at Cloudflare rebuilt Next.js on top of Vite ([Vinext](https://blog.cloudflare.com/vinext/)) in one week. Almost every line of code was written by AI. The result? 4.4x faster builds, 57% smaller bundles, 94% API coverage of Next.js 16. How was this possible?

**1,700+ unit tests and 380 end-to-end tests.**

The workflow was: define behavior through tests, let AI write the implementation, run the suite, iterate on failures. Tests were the specification the AI coded against. Without them, a single engineer could never have trusted AI-generated code at that scale.

The key lessons from Vinext apply directly to what you'll learn here:

- **Tests are specifications.** When you write a test first, you define *what* the code should do before *how*. This is the best possible prompt for an AI — a precise, executable specification. TDD + AI is a powerful combination: you define behavior, AI helps implement it, tests verify the result.
- **AI generates code, but who verifies it?** AI produces plausible code that looks correct but can have subtle bugs. Tests are your safety net. Without them, you're trusting code you don't fully understand.
- **AI amplifies your skills, it doesn't replace them.** An engineer who understands TDD can use AI to move 10x faster while maintaining quality. An engineer without testing discipline will just produce bugs 10x faster.
- **The code you don't understand will break.** AI-generated code shipped without tests becomes legacy code the moment it's merged. You can't refactor what you can't verify.

As the Vinext author puts it: *"Every line passes the same quality gates you'd expect from human-written code."* The quality gates were the tests.

**The bottom line:** In a world where code is cheap to produce, the ability to *verify and trust* that code becomes the real engineering skill. TDD is that skill. This repo is your training ground.

---

## TDD Approaches

There are two main schools of thought on how to practice TDD. They are not opposites — they are complementary tools. Skilled engineers use both depending on the situation.

### Inside-Out (Chicago School / Classicist)

You start from the innermost building blocks and work your way outward.

1. Build and test `Task` (value object, pure logic)
2. Build and test `TaskList` (uses `Task`, manages a collection)
3. Build and test `TaskService` (uses `TaskList` or `TaskRepository`)

**You discover the design bottom-up.** Each layer is built on top of tested, working components. You rarely need test doubles because you test with real collaborators.

**When to use it:**

- The domain is well understood
- You can start with concrete, isolated logic
- You want high confidence that each piece works independently

**Risk:** You might build components that don't compose well at the top, or build things you don't end up needing.

### Outside-In (London School / Mockist)

You start from what the user wants and work your way inward.

1. Write a test for `TaskService.CompleteTask(title)` — it doesn't exist yet
2. Discover that `TaskService` needs a `TaskRepository` — define the interface
3. Stub the repository in the test, make the service test pass
4. Now implement `TaskRepository` for real, tested separately

**You discover the design top-down.** Interfaces emerge from what the outer layer needs, not from what the inner layer decides to expose. This naturally produces deep modules with narrow interfaces.

**When to use it:**

- The feature requirements are clear but the internal design is not
- You want to discover interfaces and collaborators as you go
- You're building a new feature in an existing system

**Risk:** Over-mocking can couple your tests to implementation details, making them brittle. If every refactor breaks tests, you're testing the wrong things.

---

## Anatomy of a Go Test

### The basics

In Go, tests live in `_test.go` files next to the code they test. Test functions start with `Test` and receive a `*testing.T`:

```go
func TestNewTask(t *testing.T) {
    task, err := NewTask("Buy groceries", time.Now().Add(24*time.Hour))
    if err != nil {
        t.Fatal(err)
    }

    if task.Title() != "Buy groceries" {
        t.Errorf("Title() = %q, want %q", task.Title(), "Buy groceries")
    }
}
```

### Table-driven tests

Go's idiomatic way to test multiple scenarios. Instead of writing separate test functions, you define a table of inputs and expected outputs:

```go
func TestIsValidTitle(t *testing.T) {
    tests := []struct {
        name  string
        title string
        want  bool
    }{
        {name: "valid title", title: "Buy groceries", want: true},
        {name: "empty title", title: "", want: false},
        {name: "whitespace only", title: "   ", want: false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := IsValidTitle(tt.title)
            if got != tt.want {
                t.Errorf("IsValidTitle(%q) = %v, want %v", tt.title, got, tt.want)
            }
        })
    }
}
```

`t.Run` creates **subtests** — each row runs independently and shows its name on failure.

### AAA pattern (Arrange, Act, Assert)

Structure every test in three phases. Also known as **Given-When-Then**:

```go
func TestTask_Complete(t *testing.T) {
    // Arrange: set up the scenario
    task, err := NewTask("Write README", time.Now().Add(24*time.Hour))
    if err != nil {
        t.Fatal(err)
    }

    // Act: execute the behavior
    task.Complete()

    // Assert: verify the outcome
    if !task.IsCompleted() {
        t.Error("task should be completed")
    }
}
```

### Useful tools

- **`t.Helper()`** — Mark a function as a test helper. When it fails, Go reports the caller's line, not the helper's line.
- **`t.Fatal` vs `t.Error`** — `Fatal` stops the test immediately. `Error` records the failure and continues. Use `Fatal` when continuing makes no sense (e.g., nil check before using a value).
- **`go-cmp`** — Use `cmp.Diff` when comparing structs or slices. It produces clear, readable diffs on failure:

```go
if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("mismatch (-want +got):\n%s", diff)
}
```

---

## Test Doubles

A **test double** is any object that stands in for a real dependency during testing. The term comes from "stunt double" in movies — it looks like the real thing but is used in controlled situations.

Understanding the different types matters because each serves a different purpose.

### Dummy

A dummy is passed around but never actually used. It satisfies a function signature.

```go
// We need a Notifier to create the service, but this test
// doesn't exercise notification behavior.
type dummyNotifier struct{}

func (d *dummyNotifier) Notify(string, string) error { return nil }

func TestTaskService_Add(t *testing.T) {
    repo := NewInMemoryTaskRepository()
    service := NewTaskService(repo, &dummyNotifier{})

    err := service.Add("Buy groceries", time.Now().Add(24*time.Hour))
    if err != nil {
        t.Fatal(err)
    }
}
```

### Stub

A stub returns canned answers. No real logic — just predetermined responses.

```go
type stubTaskRepository struct {
    tasks []*Task
    err   error
}

func (s *stubTaskRepository) Save(*Task) error           { return s.err }
func (s *stubTaskRepository) FindAll() ([]*Task, error)  { return s.tasks, s.err }

func TestTaskService_ListOverdue(t *testing.T) {
    overdue, _ := NewTask("Old task", time.Now().Add(-48*time.Hour))
    repo := &stubTaskRepository{tasks: []*Task{overdue}}
    service := NewTaskService(repo, &dummyNotifier{})

    tasks, err := service.ListOverdue()
    if err != nil {
        t.Fatal(err)
    }
    if len(tasks) != 1 {
        t.Errorf("got %d overdue tasks, want 1", len(tasks))
    }
}
```

### Fake

A fake has working logic but takes shortcuts. The classic example is an in-memory repository instead of a real database.

```go
type InMemoryTaskRepository struct {
    tasks map[string]*Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
    return &InMemoryTaskRepository{tasks: make(map[string]*Task)}
}

func (r *InMemoryTaskRepository) Save(task *Task) error {
    r.tasks[task.Title()] = task
    return nil
}

func (r *InMemoryTaskRepository) FindAll() ([]*Task, error) {
    all := make([]*Task, 0, len(r.tasks))
    for _, t := range r.tasks {
        all = append(all, t)
    }
    return all, nil
}
```

A fake is a real implementation — it stores and retrieves data. It just uses a simpler storage mechanism.

### Spy

A spy records what happened so you can verify interactions later.

```go
type spyNotifier struct {
    calls []struct {
        title   string
        message string
    }
}

func (s *spyNotifier) Notify(title, message string) error {
    s.calls = append(s.calls, struct {
        title   string
        message string
    }{title, message})
    return nil
}

func TestTaskService_CompleteTask_NotifiesWhenOverdue(t *testing.T) {
    overdue, _ := NewTask("Old task", time.Now().Add(-48*time.Hour))
    repo := &stubTaskRepository{tasks: []*Task{overdue}}
    notifier := &spyNotifier{}
    service := NewTaskService(repo, notifier)

    _ = service.CompleteTask("Old task")

    if len(notifier.calls) != 1 {
        t.Fatalf("got %d notifications, want 1", len(notifier.calls))
    }
    if notifier.calls[0].title != "Old task" {
        t.Errorf("notified about %q, want %q", notifier.calls[0].title, "Old task")
    }
}
```

### Mock

A mock is pre-programmed with expectations and verifies that specific interactions occurred. In Go without mocking frameworks, a mock is essentially a spy with built-in assertions.

```go
type mockNotifier struct {
    t             *testing.T
    expectedTitle string
    called        bool
}

func (m *mockNotifier) Notify(title, message string) error {
    m.called = true
    if title != m.expectedTitle {
        m.t.Errorf("Notify() called with title %q, want %q", title, m.expectedTitle)
    }
    return nil
}

func (m *mockNotifier) verify() {
    m.t.Helper()
    if !m.called {
        m.t.Error("expected Notify() to be called, but it was not")
    }
}
```

### Summary

| Type  | Has logic? | Records calls? | Verifies behavior? | When to use                                  |
| ----- | ---------- | -------------- | ------------------ | -------------------------------------------- |
| Dummy | No         | No             | No                 | Satisfy a parameter you don't care about     |
| Stub  | No         | No             | No                 | Control what a dependency returns            |
| Fake  | Yes        | No             | No                 | Need a working but simplified implementation |
| Spy   | No         | Yes            | No (you assert)    | Verify that an interaction happened          |
| Mock  | No         | Yes            | Yes (self-asserts) | Verify specific expected interactions        |

---

## The Test Pyramid

```text
        /  E2E  \          Slow, expensive, high confidence
       /----------\        Few tests
      / Integration \      Medium speed, test wiring
     /----------------\    Some tests
    /    Unit Tests     \  Fast, isolated, low cost
   /--------------------\  Many tests
```

- **Unit tests** — Test a single function, method, or type in isolation. Fast, focused, most of your tests. This is where TDD lives.
- **Integration tests** — Test how multiple components work together. Verify wiring, real implementations, actual I/O.
- **End-to-end (E2E) tests** — Test the full system from the user's perspective. Slow and expensive, but highest confidence.

### Where test doubles fit

- **Unit tests** use stubs, spies, and mocks to isolate the system under test.
- **Integration tests** use fakes (e.g., in-memory database) or real implementations.
- **E2E tests** use real everything.

The pyramid shape matters: many fast unit tests at the base, fewer slow integration tests in the middle, very few E2E tests at the top. If your pyramid is inverted (mostly E2E tests), your test suite will be slow, flaky, and expensive to maintain.

---

## Design Principles

This kata follows two key principles that make code both testable and well-designed.

### Deep Modules (John Ousterhout)

A **deep module** has a simple interface that hides significant complexity. The opposite — a **shallow module** — has a complex interface but does very little work.

```text
Deep module:                    Shallow module:

  ┌──────────┐ simple           ┌──────────────────────┐ complex
  │ interface │ interface        │      interface       │ interface
  ├──────────┤                  ├──────────────────────┤
  │          │                  │                      │
  │          │ lots of          │                      │ little
  │          │ functionality    └──────────────────────┘ functionality
  │          │
  └──────────┘
```

In this kata, `TaskService.CompleteTask(title)` is a deep method: one call does a lot (find, validate, complete, save, notify). The caller doesn't need to know any of that.

### Tell, Don't Ask

Instead of querying an object's state and making decisions based on it, **tell the object what to do** and let it decide how.

```go
// BAD: Ask, then decide externally
task, _ := repo.FindByTitle("Buy milk")
if task.Status() == StatusPending {
    task.SetStatus(StatusDone)
    task.SetCompletedAt(time.Now())
    repo.Save(task)
    if task.IsOverdue(time.Now()) {
        notifier.Notify(task.Title(), "completed late")
    }
}

// GOOD: Tell the service what you want
service.CompleteTask("Buy milk")
// The service handles finding, completing, saving, and notifying.
```

The "tell" version keeps behavior with the data. The service owns the completion logic. Callers don't need to understand the internals — they just say what they want done.

These principles make code easier to test (fewer dependencies to set up), easier to change (behavior is encapsulated), and easier to understand (simple interfaces, one place to look).

---

## How to Use This Repo

### Prerequisites

- [Go](https://go.dev/dl/) 1.22 or later installed
- A text editor or IDE (VS Code with the Go extension works well)
- Basic Go syntax: variables, functions, structs, slices, error handling

### Getting started

```bash
git clone https://github.com/jferrl/go-tdd-kata.git
cd go-tdd-kata
```

### Running tests

From any level folder:

```bash
go test -v ./...
```

### The workflow for each level

1. Read the level's `README.md` — understand the concepts and the exercise
2. Create your `.go` and `_test.go` files in the level folder
3. Follow TDD: write a failing test first, then make it pass, then refactor
4. Move to the next level when you're comfortable

---

## Levels

| Level | Folder                                                | What You Learn                                          |
|-------|-------------------------------------------------------|---------------------------------------------------------|
| 1     | [01-pure-functions](01-pure-functions/)               | Red-Green-Refactor, table-driven tests, AAA pattern     |
| 2     | [02-value-objects](02-value-objects/)                 | Structs with behavior, constructors, encapsulation      |
| 3     | [03-inside-out](03-inside-out/)                       | Inside-out TDD, composing tested units bottom-up        |
| 4     | [04-interfaces-and-stubs](04-interfaces-and-stubs/)   | Interfaces, dependency injection, hand-rolled stubs     |
| 5     | [05-fakes-and-spies](05-fakes-and-spies/)             | Fakes, spies, asserting on interactions                 |
| 6     | [06-outside-in](06-outside-in/)                       | Outside-in TDD, mocks, discovering design top-down      |

Each level builds on the previous one. By the end, you will have a fully tested, layered Todo application — and a solid understanding of TDD in Go.

---

## References

- Kent Beck, *Test-Driven Development: By Example* (2002)
- Robert C. Martin, *Clean Code* (2008)
- John Ousterhout, *A Philosophy of Software Design* (2018)
- Martin Fowler, [Mocks Aren't Stubs](https://martinfowler.com/articles/mocksArentStubs.html)
- Cloudflare, [Vinext: Rebuilding Next.js on Vite](https://blog.cloudflare.com/vinext/)
- Nagappan et al., *Realizing Quality Improvement Through Test-Driven Development* (2008)
- Google, [Go Style Guide](https://google.github.io/styleguide/go/)
