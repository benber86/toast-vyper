# Side-by-Side Editors — Design & Learnings

Experimental notes on the direction taken for side-by-side editors (issue #49).
This documents *what* we built and *why*, plus the trade-offs and rough edges
we know about. Read it before touching the split code — several "obvious" fixes
have non-obvious interactions.

## Direction

The goal: two editable panes, set up by keyboard **or** drag-and-drop, with
tabs owned by their pane so panes can be managed from the UI.

The model is VS Code's editor groups:

- **Each pane owns its tab strip.** When split, the layout is:

  ```
  row 0:   breadcrumb (full width, shows the focused pane's file)
  row 1:   file tree | left pane tabs | │ | right pane tabs
  rows 2+: file tree | left editor    | │ | right editor
  ```

  The tab strips render *inside* their pane columns, not as a global top
  strip. When not split, the layout is unchanged (tabs on top, breadcrumb
  below). `tabBarRow()`/`breadcrumbRow()` encode the two layouts so mouse
  routing never hardcodes a row.

- **Keyboard split** (`ctrl+alt+\`) duplicates the active file into the right
  pane *with its own tab*. Both panes show the same buffer (same `EditBuffer`
  rope), so edits and undo stay in sync. Pressing the key again cycles focus.

- **Panes are managed from the UI.** Each strip ends in a `✕` that closes that
  pane (`PaneCloseRequestMsg`); `ctrl+alt+w` closes the right pane. Closing or
  moving every tab out of a pane reverts to single view with whichever pane
  still has tabs (`paneEmptied` → `collapseSplit`).

- **Drag-and-drop moves tabs between panes.** Dragging a tab onto the other
  pane's strip or editor moves its tab there (creating the split on demand).
  The drop target pane is the pane under the pointer; the tab-strip row is a
  valid drop target.

## Key implementation decisions

- **Explicit pane routing on tab messages.** `ActiveBufferChangedMsg`,
  `BufferOpenedMsg`, `BufferClosedMsg`, `CloseTabRequestMsg`,
  `TabDragStartMsg` carry a `Pane int` (0 = left, 1 = right). This removes all
  guesswork about which pane a message belongs to — important because a buffer
  can be *tabbed in both panes* after a keyboard split.

- **Global buffer bookkeeping, per-pane tabs.** `openBuffers`,
  `bufferSnapshots`, `closedSnapshots` stay global. Closing a tab runs
  `closeBufferFully` only when **no** pane still tabs the buffer; unsaved
  edits are preserved in `closedSnapshots` for the quit prompt.

- **Shared ropes, never divergent copies.** Panes showing the same file share
  the `EditBuffer`. Switching a pane to a buffer that's active in the *other*
  pane restores the other pane's snapshot instead of reloading from disk —
  otherwise the fresh copy's edits would be silently lost when the pane
  closed.

- **Find/replace targets the pane it opened in** (`findReplacePane`). Search
  and markdown preview still replace the whole content area.

## Rough edges & open questions

- **Move vs duplicate is inconsistent.** Keyboard split *duplicates* the
  active tab (so a single-file split works), drag-drop *moves* tabs. A file
  can therefore be tabbed in both panes at once. Combined with
  collapse-on-empty, moving the last tab out of a pane closes the split —
  which is the requested behavior, but worth confirming it's what we want.
- **One `bufferSnapshots` slot per buffer.** With a buffer open in both panes,
  the two panes' cursors/viewports can diverge (content stays in sync via the
  rope); switching away and back may restore the other pane's cursor.
- **The quit `✕` disappears while split** (replaced by per-pane close
  buttons). `ctrl+q` still works; we may want a small global affordance.
- **Per-pane breadcrumbs** (VS Code shows one per group) were deferred — we
  have one breadcrumb showing the focused pane.
- **Tab reordering within a pane** isn't implemented; middle-click close works
  per pane.
- LSP diagnostics / git diff are keyed by bufferID and reach both panes (each
  filters) — fine, but a duplicated buffer shares one LSP document.

## What we'd plan differently next time

1. **Decide tab ownership up front.** Duplicated vs partitioned tabs changes
   message routing, the collapse rule, and drag semantics. We converged on
   "duplicate on split, move on drop" late.
2. **Pin the layout before wiring mouse routing.** The tab strips moved from
   a global top row into the pane columns after the fact; every mouse handler
   had to be re-checked for hardcoded rows.
3. **Write the data-safety invariants first:** shared ropes, when
   `closeBufferFully` runs, and where unsaved edits go. Two bugs (divergent
   rope copies losing edits, stale pane headers) came from not having these
   written down.
