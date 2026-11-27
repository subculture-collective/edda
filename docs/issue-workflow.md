# Issue workflow

On 2026-11-27 the 1,765 open task issues that came out of the planning decisions
(#276, #321–#334, #54, #80 and a few smaller ones) were folded into 18 tracking
issues, one per decision. Each tracker is titled `Tracker: …`, carries the
`Planning/Epic` label and holds a checklist of the old issues. Those issues
were closed as "not planned" only to shorten the open list. The work is still
planned, and each closed issue keeps its full specification.

## Find work

List the trackers by title only. Do not list every issue with its body: the old
specifications average 4.5 KB each.

```sh
gh issue list --state open --label Planning/Epic --search 'Tracker in:title' --json number,title
gh issue view <tracker>              # the checklist
gh issue view <item>                 # one item's full specification
```

Read an item's specification only when you take it. Check its `Blocked by`
section; a blocker is done when it is ticked in its tracker or closed as
completed.

## One pull request per batch

- Take about six related unchecked items from one tracker per pull request,
  the way the SRD catalogue waves do. Fewer, larger pull requests mean fewer CI
  runs.
- In the pull request body write `Part of #<tracker>` and list the item
  numbers it completes. Do not write `Closes #<item>`: the items are already
  closed.
- After the pull request merges, tick its items in the tracker. Read the
  tracker again just before editing it, because another pull request may have
  ticked items since you last read it.
- Close a tracker when every item is ticked.
- For a local race run use `task test:race -- ./internal/rules/srd521/...`.
  It sets `-timeout 45m`; that package takes about 15 minutes with `-race`.

## New work

Open new issues as before. When a planning decision splits into many small
tasks, put them in one tracking issue with a checklist instead of one issue
per task.
