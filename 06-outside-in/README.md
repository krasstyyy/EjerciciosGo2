# Level 6: Outside-In TDD

## Learning Objectives

- Practice **outside-in TDD**: start from what the user wants, discover the design inward
- Experience how interfaces emerge from need, not from anticipation
- Use **mocks** — test doubles with built-in expectations
- Contrast outside-in with inside-out and understand when to use each

## Concepts

### Outside-in TDD

In Levels 1-5, you built bottom-up: `Task` first, then `TaskList`, then `TaskService`. You knew what pieces you needed because we told you.

In real projects, you rarely know the full design upfront. Outside-in TDD helps you **discover** the design:

1. Start with what the user (or caller) wants: *"I want to tag a task and search by tag"*
2. Write a test for the outermost layer — `TaskService.TagTask(title, tag)`
3. The test won't compile — `TaskService` doesn't have that method yet. Good.
4. To implement it, you realize the service needs something that stores tags. Define `TagRepository` as an interface.
5. Stub or mock `TagRepository` in the test. Make the test pass.
6. Now go one layer deeper: implement a real `TagRepository`, tested separately.

You're building outside-in: the outer layer **tells** you what the inner layer's interface should be. This naturally produces deep modules with narrow interfaces — the interface has exactly the methods the consumer needs, nothing more.

### Mocks

A **mock** is a test double pre-programmed with expectations. Unlike a spy (which records and lets you assert later), a mock knows what it expects and can verify itself:

```go
type mockTagRepo struct {
    t             *testing.T
    expectTitle   string
    expectTag     string
    addTagCalled  bool
}

func (m *mockTagRepo) AddTag(taskID, tag string) error {
    m.addTagCalled = true
    if taskID != m.expectTitle {
        m.t.Errorf("AddTag() taskID = %q, want %q", taskID, m.expectTitle)
    }
    if tag != m.expectTag {
        m.t.Errorf("AddTag() tag = %q, want %q", tag, m.expectTag)
    }
    return nil
}

func (m *mockTagRepo) verify() {
    m.t.Helper()
    if !m.addTagCalled {
        m.t.Error("expected AddTag() to be called")
    }
}
```

Use mocks when you need to verify that the service **tells** its collaborator the right thing with the right arguments. This is the "tell, don't ask" principle verified through tests.

### When mocks become a problem

Mocks couple your tests to implementation details. If you refactor how the service calls its collaborators (even if the end result is the same), mock-based tests break. This is the main criticism of outside-in / London-school TDD.

Use mocks for **interactions that matter** — when the service's job *is* to coordinate collaborators correctly. Don't mock everything.

## Exercise

You will build a **tagging feature** from scratch using outside-in TDD. Start from the service layer and work your way inward.

This time, the exercise doesn't tell you the exact interface or implementation. The outside-in process will guide you.

### Step 1: Start from the outside

Write a test for a method that doesn't exist yet:

```go
func TestTaskService_TagTask(t *testing.T) {
    // What does TagTask need?
    // You'll figure out as you go.
}
```

Think about what `TagTask(title string, tag string) error` should do:

- Find the task by title (needs `TaskRepository`)
- Associate the tag with the task (needs... something new)

Let the compiler errors guide you. When you need a new capability, define an interface.

### Step 2: Discover TagRepository

As you implement `TagTask`, you'll need an interface for tag storage. Define it based on what the service actually needs — nothing more.

Think about:

- What method does the service call to store a tag?
- What arguments does it need?
- Does it need the task ID or the task title?

Stub or mock this new interface in your test. Make the test pass.

### Step 3: SearchByTag

Write a test for `SearchByTag(tag string) ([]*Task, error)`.

Think about the flow:

1. The service asks the tag store for task IDs with this tag
2. The service asks the task repository for those tasks
3. The service returns them

This requires coordination between two collaborators. Use stubs, spies, or mocks — choose based on what you need to verify.

### Step 4: Implement TagRepository for real

Now go one layer deeper. Build an `InMemoryTagRepository` (fake) that satisfies your interface.

Test it independently — this is the inside-out part of the process. Outside-in and inside-out aren't mutually exclusive: you discover the design from the outside, then test the implementation from the inside.

### Step 5: Wire it together

Write an integration-style test that uses real implementations (your fakes) instead of stubs:

1. Create a service with `InMemoryTaskRepository` and `InMemoryTagRepository`
2. Add a task
3. Tag it
4. Search by tag
5. Verify the task is returned

This test exercises the full flow without mocks. It sits higher on the test pyramid.

### Step 6: Compare the approaches

Look at your test file. You should have:

- **Mock/stub-based tests** — Fast, isolated, verify specific interactions
- **Fake-based integration test** — Exercises the real flow, higher confidence

Think about which tests you'd want in a real project. The answer is usually: both, in the right proportion.

## Hints

<details>
<summary>Letting compiler errors guide you</summary>

Outside-in TDD embraces compilation failures as design feedback. When your test doesn't compile:

1. The service method doesn't exist — create it
2. The service needs a dependency — define an interface
3. The interface method doesn't exist — add it

Each error tells you what to build next.

</details>

<details>
<summary>Mock vs spy for this exercise</summary>

For `TagTask`, a mock works well: you want to verify that the service calls `AddTag` with the correct task ID and tag.

For `SearchByTag`, a stub is simpler: you just need to control what `FindByTag` returns.

Choose the simplest double that answers your test's question.

</details>

## Reflection

After completing this level, think about:

- How did it feel to start from a method that doesn't exist? Did the compiler errors help or frustrate you?
- The `TagRepository` interface emerged from what the service needed. Would you have designed it differently if you started bottom-up?
- Look at your mock tests vs your integration test. Which gives you more confidence? Which is easier to maintain?
- When would you choose inside-out over outside-in? Think about situations where you know the domain well vs where you're exploring.
- Could you combine both approaches in a real project? Most experienced developers do — inside-out for well-understood core logic, outside-in for features with unclear boundaries.
