<!-- edda-333:condition-output -->
Decision #333; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1 PDF178 (Charmed), 182 (Frightened) and186 (Prone), as consumed by Fey Ancestry on85, Hill's Tumble on85 and Brave on86; SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

Qualification only, one PR. This task owns no runtime, condition rule, species trait or production admission. #488 keeps the Charmed, Frightened and Prone contract and this task changes none of it.

Why it exists: #488 has 97 outgoing blocking relations and GitHub allows 100. The Elf, Goliath and Halfling tasks of decision #333 are blocked by this one output instead of taking the last three slots, and this output is blocked by #488.

Prove against the merged #488 implementation, without reimplementing it, that an origin consumer can rely on each of these: a saving throw made to avoid or to end the Charmed condition, and one made to avoid or to end the Frightened condition, is identifiable as such on its receipt, so that a named Advantage source can attach to it and to no other save; the Prone condition can be applied to a target by a named source inside a confirmed hit; and each condition keeps one store, its source and its ending.

Positive fixtures: an isolated test Advantage source attached to saves that avoid or end Charmed applies to a save against a declared charm effect and to the repeat save that ends it, and not to a save against Frightened. The same holds with the two conditions exchanged. An isolated test source applies Prone to a Large target on a confirmed hit and the ordinary #488 Prone effects follow. Reject: Advantage on a save that neither avoids nor ends the named condition; a second condition store; Prone applied on a miss or after the hit is committed; closing this task from #488's issue state or this decision in place of the merged implementation and its evidence; any new condition rule, species trait or eligibility added by this PR.

Evidence is deterministic source and state tests plus isolated PostgreSQL replay, reload, rollback and concurrency checks run through the existing owner's own prepare and confirm path; previews stay inert and nothing here adds a command, cost or die. Shared web and trusted TUI show the same terms and explicit unavailable reasons the owner already exposes; this PR adds no new surface. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission. Unknown or unqualified behavior stays unavailable.

Native parent #333; directly block #337. Hard prerequisites: #488. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
