<!-- edda-333:creation-output -->
Decision #333; roadmap #47. One implementation PR; Task / Size S.

Source: English SRD 5.2.1 PDF19–23 (the creation steps whose grants carry provenance: background, species, class, languages and proficiencies); SHA-256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Retain CC-BY-4.0 attribution and adaptation notice in `docs/srd-5.2.1-conformance.md`. Fixture and ownership decisions are Edda additions.

Qualification only, one PR. This task owns no runtime, grant, species trait, background, feat, spell effect or production admission. #430 keeps the whole creation provenance contract and this task changes none of it.

Why it exists: #430 already has 98 outgoing blocking relations and GitHub allows 100, so the thirteen origin tasks of decision #333 that record a creation grant cannot each take a native edge to it. They are blocked by this one output instead, and this output is blocked by #430. It is the same output pattern as #1738 for Warlock.

Prove against the merged #430 implementation, without reimplementing it, that an origin consumer can rely on each of these: a distinct provenance instance for a species, background, class, feat and language grant of the same skill, tool or language; #430's duplicate-proficiency resolution, with no silent replacement and no doubled Proficiency Bonus; rejection of an unsupported species or background identity; removal or replacement of one source changing only that source's grants; and exactly one grant per replayed confirmation.

Positive fixtures: an isolated test species grant and an isolated test background grant of Perception are both retained with their own source and the Proficiency Bonus is applied once. A language recorded by the creation choice and the same language from a class feature remain two source records. Removing the isolated background grant leaves the species grant in place. Reject: a consumer write with no source instance; a second grant on same-ID replay; an unknown species or background key; closing this task from #430's issue state, a catalogue record or this decision in place of the merged implementation and its evidence; any new grant, effect, species, background or eligibility added by this PR.

Evidence is deterministic source and state tests plus isolated PostgreSQL replay, reload, rollback and concurrency checks run through the existing owner's own prepare and confirm path; previews stay inert and nothing here adds a command, cost or die. Shared web and trusted TUI show the same terms and explicit unavailable reasons the owner already exposes; this PR adds no new surface. Keep private GM facts, hidden actors and private choices out of discovery, preview, errors and public history. Tests use explicit isolated entitlements and grant no production admission. Unknown or unqualified behavior stays unavailable.

Native parent #333; directly block #337. Hard prerequisites: #430. Existing bodies, parents and edges remain unchanged. No new #283 child.

Implementation verification: deterministic source/numeric/choice tests with injectable dice, isolated PostgreSQL atomic rollback/replay/concurrency/reload, shared web/TUI disclosure/confirmation checks, `GOMAXPROCS=2 go test -p 2 -short -race ./...` and `GOMAXPROCS=2 go vet -p 2 ./...`. Short tests do not qualify PostgreSQL or devices. No production grant, deployment or Full Rules promotion.
