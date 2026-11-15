# Edda story-led play and SRD 5.2.1 gameplay plan

Status: published as a GitHub roadmap following the owner's interview. The decisions below are agreed; layout details and delivery boundaries are recommendations. Publication does not implement gameplay or qualify a release.

Source checkout: `b84d5dc8bde8438920c867c2dab4d64c50c001b5`, reviewed November 14, 2026. Existing review issues are [#28–#46](https://github.com/subculture-collective/edda/issues?q=is%3Aissue+is%3Aopen).

## Published roadmap

[Open the GitHub roadmap](https://github.com/subculture-collective/edda/issues/47). It contains 55 work issues under seven epics: 36 new work issues and all 19 original review issues, reused with their evidence retained. All 63 roadmap/epic/work issues are in the repository's linked Project.

The delivery units below are design checkpoints. The work issues and native blocked-by links are the executable breakdown; an aggregate checkpoint is not an instruction to serialize otherwise independent PRs.

| Epic | Tracker |
| --- | --- |
| Make campaign actions private, durable and recoverable | [#48](https://github.com/subculture-collective/edda/issues/48) |
| Keep story-led gameplay usable within one screen | [#49](https://github.com/subculture-collective/edda/issues/49) |
| Play supported SRD actions from the interactive sheet | [#50](https://github.com/subculture-collective/edda/issues/50) |
| Create and advance valid rules-versioned characters | [#51](https://github.com/subculture-collective/edda/issues/51) |
| Offer useful player-known action suggestions and reminders | [#52](https://github.com/subculture-collective/edda/issues/52) |
| Read and export faithful campaign chapters | [#53](https://github.com/subculture-collective/edda/issues/53) |
| Complete and qualify the declared SRD 5.2.1 coverage | [#54](https://github.com/subculture-collective/edda/issues/54) |

| Design unit | Published work issues |
| --- | --- |
| P01 | [#31](https://github.com/subculture-collective/edda/issues/31), [#32](https://github.com/subculture-collective/edda/issues/32) |
| P02 | [#28](https://github.com/subculture-collective/edda/issues/28) |
| P03 | [#29](https://github.com/subculture-collective/edda/issues/29) |
| P04 | [#30](https://github.com/subculture-collective/edda/issues/30), [#59](https://github.com/subculture-collective/edda/issues/59) |
| P05 | [#45](https://github.com/subculture-collective/edda/issues/45) |
| P06 | [#55](https://github.com/subculture-collective/edda/issues/55), [#46](https://github.com/subculture-collective/edda/issues/46) |
| P07 | [#56](https://github.com/subculture-collective/edda/issues/56) |
| P08 | [#57](https://github.com/subculture-collective/edda/issues/57), [#58](https://github.com/subculture-collective/edda/issues/58) |
| P09 | [#61](https://github.com/subculture-collective/edda/issues/61) |
| P10 | [#63](https://github.com/subculture-collective/edda/issues/63), [#64](https://github.com/subculture-collective/edda/issues/64) |
| P11 | [#60](https://github.com/subculture-collective/edda/issues/60), [#90](https://github.com/subculture-collective/edda/issues/90), [#62](https://github.com/subculture-collective/edda/issues/62), [#71](https://github.com/subculture-collective/edda/issues/71) |
| P12 | [#65](https://github.com/subculture-collective/edda/issues/65), [#79](https://github.com/subculture-collective/edda/issues/79) |
| P13 | [#66](https://github.com/subculture-collective/edda/issues/66), [#67](https://github.com/subculture-collective/edda/issues/67), [#81](https://github.com/subculture-collective/edda/issues/81) |
| P14 | [#68](https://github.com/subculture-collective/edda/issues/68), [#77](https://github.com/subculture-collective/edda/issues/77) |
| P15 | [#69](https://github.com/subculture-collective/edda/issues/69), [#72](https://github.com/subculture-collective/edda/issues/72), [#83](https://github.com/subculture-collective/edda/issues/83) |
| P16 | [#33](https://github.com/subculture-collective/edda/issues/33), [#34](https://github.com/subculture-collective/edda/issues/34), [#35](https://github.com/subculture-collective/edda/issues/35), [#36](https://github.com/subculture-collective/edda/issues/36) |
| P17 | [#74](https://github.com/subculture-collective/edda/issues/74), [#80](https://github.com/subculture-collective/edda/issues/80), [#82](https://github.com/subculture-collective/edda/issues/82), [#73](https://github.com/subculture-collective/edda/issues/73), [#88](https://github.com/subculture-collective/edda/issues/88) |
| P18 | [#75](https://github.com/subculture-collective/edda/issues/75), [#84](https://github.com/subculture-collective/edda/issues/84), [#85](https://github.com/subculture-collective/edda/issues/85) |
| P19 | [#38](https://github.com/subculture-collective/edda/issues/38), [#39](https://github.com/subculture-collective/edda/issues/39), [#40](https://github.com/subculture-collective/edda/issues/40) |
| P20 | [#70](https://github.com/subculture-collective/edda/issues/70), [#76](https://github.com/subculture-collective/edda/issues/76), [#86](https://github.com/subculture-collective/edda/issues/86), [#87](https://github.com/subculture-collective/edda/issues/87) |
| P21 | [#37](https://github.com/subculture-collective/edda/issues/37), [#42](https://github.com/subculture-collective/edda/issues/42), [#43](https://github.com/subculture-collective/edda/issues/43) |
| P22 | [#44](https://github.com/subculture-collective/edda/issues/44), [#46](https://github.com/subculture-collective/edda/issues/46), [#78](https://github.com/subculture-collective/edda/issues/78) |
| P23 | [#41](https://github.com/subculture-collective/edda/issues/41) |
| P24 | [#89](https://github.com/subculture-collective/edda/issues/89) |

The one-screen gameplay requirement is explicit in the viewport concept, shell, bounded composer, responsive assessment and regression issues. It is an acceptance contract, not an optional style preference.

The full SRD coverage target remains the expansion epic's exit criterion. Its three bounded investigations publish concrete class/advancement, spell/item/effect, and monster/exploration PR batches into this same roadmap. Closing those investigations does not complete the content implementation. Large integration slices and the genuine multi-branch qualification dependency join are recorded sizing/graph warnings, not hidden prerequisites.

## Goal and agreed decisions

Make Edda a story-led RPG with a readable transcript, a faithful narrative retelling, useful action assistance, and a detailed interactive character sheet backed by consistent rules. Desktop and phone must support complete play through layouts suited to each. The main gameplay screen must fit within one viewport, retaining the current application's screen-fitting behavior while improving readability and interaction.

The owner selected:

- SRD 5.2.1 as the authoritative reference for detailed rules.
- Interactive character sheets: select an ability, choose targets and options, inspect costs, then confirm.
- Rolls on confirmation, with visible dice, modifiers, advantage/disadvantage and outcomes. No separate Roll prompt.
- Narrative combat with initiative, target selection, ranges and conditions. A grid-based battle view is not required.
- Faithful chapters: flowing prose that compresses repetition while preserving the played plot, without invented events or hidden lore.
- Subtle action completions accepted with Tab, including useful choices based on skills, spells, equipment, known information and situational reminders.
- A refresh of the existing Art Deco interface. Keep the current design language.
- Keep the main gameplay screen within one viewport, similar to the current layout. The primary controls and current context stay in the frame; long content scrolls inside its panel.
- Keep the campaign/character choice-based creation process as a foundation. Its overall process is satisfactory; the interview path needs more testing.
- Treat this as full-stack work. Detailed rules are the development priority; the lighter modes can follow from that foundation.

## Scope and working assumptions

In scope: the campaign library and creation surfaces; transcript and action composer; character sheet and character generation; rules, resources and progression; narrative combat; contextual suggestions; faithful chapters; world/journal navigation; mobile and accessibility behavior; persistence, reconnect and replay contracts; fixtures and verification.

Recommended defaults, which the owner has not individually selected:

- Continue with one player-controlled character per campaign and the existing NPC model. This plan does not introduce cooperative multiplayer or a party-management product.
- Offer full rules as the recommended mode for new campaigns once its supported content is qualified. Keep existing mode identifiers and legacy behavior available during the transition.
- Use a desktop companion panel and a phone character destination/drawer rather than displaying a miniature desktop HUD on a phone.
- Preserve third-person, past-tense chapter prose, matching the existing summarizer's direction.
- Start chapters with explicit on-demand generation from completed play; add automatic generation at persisted scene boundaries once that boundary exists. Existing ten-turn summaries remain a separate recap/memory mechanism.
- Preserve the existing account, model-provider, campaign, replay and export capabilities. Keep the current React/Vite/TypeScript/Tailwind and Go/PostgreSQL stack.

Outside this plan: a tactical grid editor, a new visual identity, non-SRD expansion catalogues, a marketplace, cooperative multiplayer, native apps, a new hosting stack, and unrelated backend maintenance. Existing security and consistency defects that affect these workflows are dependencies, not exclusions.

The target is detailed gameplay against the selected reference, including the SRD's character advancement and content. A first playable slice is a delivery checkpoint, not permission to call partial support complete. Each unsupported rule or content item remains visible in a coverage register until implemented or explicitly excluded by a product decision.

## Evidence from the existing code

| Existing boundary | What to retain or change |
| --- | --- |
| `frontend/src/App.tsx` | Account, library, creation, play and replay routes already exist. Retain their public paths. |
| `frontend/src/pages/CampaignPlayPage.tsx` and `components/layout/AppShell.tsx` | Eight tabs, sidebar and widescreen framing exist. Retain the screen-fitting gameplay frame, improve its reading surface and controls, and adapt the arrangement to phone viewports. |
| `frontend/src/components/narrative/PlayerInput.tsx` and `ChoiceList.tsx` | There is free text and direct choice submission, but no grounded inline completion or preparation/confirmation workflow. |
| `frontend/src/components/character/CharacterSheet.tsx` and `frontend/src/api/types.ts` | The sheet has generic stats and name/description abilities. Typed spell resources, usable actions and full rule-derived projections need new contracts. |
| `internal/rules/`, `internal/progression/leveling.go`, `internal/combat/` | Existing mechanics are a starting point, not evidence of SRD 5.2.1 conformance. Custom XP thresholds, skill-point presentation, condition handling and round resolution need reference-based review. |
| `internal/engine/runtime.go` and `turn_completer.go` | Web/TUI share the engine. Current tool effects precede history persistence; turns are not serialized per campaign. Address #29 and #30 before new resource spending. |
| `internal/journal/summarizer.go` and `frontend/src/components/journal/JournalPanel.tsx` | Short turn-range narrative summaries and notes already exist. Extend these capabilities for durable chapters rather than assuming a short summary is a complete story reader. |
| `frontend/src/hooks/useNarrative.ts`, `useWebSocket.ts` and `useRefreshAfterTurn.ts` | Preserve history/live hydration while fixing campaign isolation, reconnect reconciliation and event retention. |
| `docs/brand/edda/GUIDE.md`, `frontend/tailwind.config.ts`, `frontend/src/index.css` | The identity pack and application use different visual systems. The owner's choice here is to refresh the application system; keep identity assets and their usage rules intact. |

Retain the older HUD rollout's screen-fitting gameplay principle, semantic colors and shared control conventions. Adapt its widescreen arrangement for equal desktop/phone play: the frame follows the available viewport rather than requiring a landscape aspect ratio. Existing architecture vocabulary and the shared web/TUI engine boundary remain in use.

## Product structure and visual direction

Use four main destinations, with contextual shortcuts:

| Destination | Purpose |
| --- | --- |
| Play | Transcript, current scene, action composer, suggested actions and contextual combat controls. |
| Character | Abilities, saving throws, skills, attacks, spells, equipment, resources, conditions and progression. Inventory belongs here. |
| World | Quests, NPCs, places/world map, known facts, relationships and codex. Quest/location shortcuts can open context without leaving the story. |
| Journal | Story chapters, player notes and session/history index, with access to replay and export. A Read story shortcut opens the same chapter reader from Play. |

Desktop: keep a viewport-height gameplay frame with a compact context/status header, a central internally scrolling transcript, an anchored action area, navigation and an optional companion panel for the character or selected world context. Keep a comfortable reading width within that frame. Opening details preserves the reading position and action draft. Character also has a full-width view for deeper inspection and level-up choices, with its own internal scrolling.

Phone: keep the same viewport-fitting principle with a compact header, one internally scrolling reading column, a bottom action area, four labeled destinations and full-height detail sheets or views. Preserve the draft, scroll and selected context when moving between them. Account for safe-area insets, browser chrome and keyboard resizing. Adapt the arrangement without requiring rotation, zooming out or a fixed 16:9 aspect ratio.

### Main gameplay viewport contract

- At the supported default viewport sizes, the main gameplay page fits the screen without document-level vertical or horizontal scrolling. The player can read the current scene and reach action input, submission, navigation and essential status within that frame.
- Transcript history, long choices, character details and world information scroll in their own bounded regions. Keep the action composer and confirmation controls anchored while reading earlier transcript entries.
- Fit the gameplay workspace, not the entire transcript or full character sheet at once. Use selected panels, compact summaries, expandable choices and detail sheets to expose depth without extending the page below the screen.
- During combat, initiative/current actor, essential resources and action availability remain reachable in the frame. Detailed rolls, conditions and target information open within bounded panels or overlays.
- With the phone keyboard open, resize the visible workspace and keep the active composer and submit/confirm controls visible. Compact secondary chrome as needed; the keyboard must not cover the action controls or trap panel scrolling.
- Implement against dynamic viewport height with an appropriate fallback, flexible panel sizing and explicit scroll ownership. Keep text and touch targets readable; do not satisfy the requirement by shrinking the entire HUD.
- At increased text size or browser zoom, allow accessible reflow when the compact arrangement cannot fit. Preserve access and focus rather than clipping controls to enforce a rigid frame.
- This viewport contract applies to main gameplay. Chapter reading and other destinations can use layouts suited to their longer content.

Keep the dark canvas, warm reading text, Art Deco accents and semantic jade/sapphire/ruby/gold roles. Reduce repeated borders, badges, uppercase labels and timestamps around every paragraph. Use display typography for headings and readable body typography for long prose; compare the current body font with a reading alternative during prototypes before deciding to replace it. Consolidate any chosen sizes, colors, spacing and focus styles in shared tokens.

Create desktop and phone concepts from the same fixture before component replacement. Demonstrate the gameplay frame with a long transcript, long choices, a character action, combat and the phone keyboard open; also compare a story chapter. Visual review selects details within the agreed design language; it does not reopen the screen-fitting requirement or the product direction.

## Transcript and reading behavior

- Make player actions and game-master prose distinguishable without placing every entry in a heavy card. Use readable paragraph widths, spacing and dialogue formatting.
- Separate prose from mechanical receipts. Rolls, resource changes and effects can expand for detail without interrupting the story; critical costs and consequences remain visible.
- Keep the newest scene easy to reach. Auto-follow only while the player is already near the end; do not pull them away from earlier text they are reading.
- Scroll history within the transcript panel while keeping the gameplay frame, essential status and action controls in place. Long choices and growing action drafts must use bounded expansion or internal scrolling.
- Preserve position across section changes, reconnects and reloads. Provide search or a scene/history index for long campaigns, implemented after stable turn identifiers exist.
- Restore choices and action results from persisted history, not only from transient route state. Maintain access to replay and exports.
- Distinguish preparing, resolving, narrating, complete, retryable failure and reconnecting. A committed action with delayed narration must not look like an unsubmitted action.
- Incremental narrative text is a separate enhancement (#41). The reading redesign can ship with status updates and a final prose response; if token streaming is added, mark uncommitted text appropriately and reconcile it with the saved result.

## Character sheet and creation

Expose rule-derived values with a way to inspect their sources:

- Identity, species, background, class/subclass, level and advancement choices.
- Ability scores/modifiers, proficiency bonus, saving throws, skills, expertise and relevant passive values.
- HP/max/temporary HP, defenses, speed, initiative, conditions, concentration, exhaustion and death-save state when applicable.
- Attacks, weapon properties/masteries, damage and available actions; action, bonus action, reaction and movement availability.
- Known/prepared spells, casting ability, spell attack/save values, slots or other casting resources, upcasting, components, concentration, duration, range and targets.
- Equipment, currency, consumables, ammunition, charges, attunement and carrying limits as supported by the selected rules.
- Class features, feats, limited uses, recovery conditions, rests, Hit Dice and progression.

Values come from the server's character projection. The client must not independently maintain XP thresholds, infer available spell slots from prose, or treat arbitrary ability names as executable mechanics. Do not display internal character/location UUIDs as primary character information.

Selecting a usable entry opens its details and action preparation: what it does, legal targets/options, costs and relevant availability. Unavailable actions remain discoverable with a player-safe reason. Pure reference entries are not rendered as fake working controls.

Keep the choice-based creation flow's general sequence. Full-rules selections become catalogue-backed and validated together: origins, class, abilities, proficiencies, spells and starting equipment. Show derived values and outstanding required choices before world generation. The interview proposes a structured draft using the same catalogue; it cannot invent an unsupported build or bypass validation. Preserve inputs on errors and make world-building progress/recovery explicit (#33–#35).

Do not reuse the current skill-point or custom XP systems as if they were the chosen reference. Keep them under legacy rules while adding the appropriate derived proficiency and progression behavior for new rules-versioned characters.

## Action preparation, confirmation and resolution

Use one server-owned action path for sheet controls, selected choices and free text. Extend the existing API/engine rather than creating a separate gameplay engine in React.

1. **Prepare.** A structured selection or free-text intent is associated with actor, targets/options and the current campaign revision. Preparation is read-only: no resource consumption, dice, time advancement or hidden information disclosure. A language model can interpret text and propose intent; the rules module validates it.
2. **Preview.** Return the available action, required choices, costs, player-visible modifiers, validation reasons, rules version and state revision. Ambiguous or materially costly interpretations require a preview. Straightforward conversation/reference requests can remain a single submit operation when they have no mechanical commitment.
3. **Confirm.** Submit a unique client action identifier and the prepared action. Recheck ownership, state revision, targets and resources. If relevant state changed, return an updated preview instead of executing stale intent.
4. **Resolve and commit.** Serialize the campaign action; roll on the server; apply costs and outcomes; update combat/time as appropriate; persist a resolution receipt and campaign revision atomically. Retries return the same committed result. Reusing an identifier for a different request is rejected.
5. **Narrate.** Generate prose from the committed outcome. Failure or regeneration here never rerolls dice or spends resources again. Store narration against the stable resolution identifier; reconnect can retrieve the receipt while prose is pending.
6. **Present and reconcile.** Update the sheet and transcript from the authoritative result. Replay and history use the same persisted receipt; refresh relevant queries without optimistic spending.

Use short database transactions for validation/commit, not transactions held across external model requests. Persist pending work for recoverable narration using the existing database; a new queue service is not a requirement. Account for both deterministic mechanical effects and world changes: model-generated changes are validated proposals, and any later change becomes an explicit revisioned action rather than an unnoticed tool side effect.

Proposed contract concepts, to finalize in a focused contract PR:

| Concept | Required data and behavior |
| --- | --- |
| Rules reference | Ruleset ID/version, content revision and supported-rule metadata. |
| Character projection | Typed derived values, resources, conditions, capabilities and source references. |
| Prepared action | Stable action reference, actor, options, targets, costs, validation, expected revision; contains no rolled result. |
| Action resolution | Action/turn identifiers, resulting revision, committed effects, dice receipt, known outcomes and narration status. |
| Suggestion | Display text/completion, optional prepared-action reference, player-known reasons and source revision. |
| Chapter | Covered committed-turn range/IDs, knowledge boundary, text, generation state and revision history. |

Evolve `pkg/api`, frontend API types and WebSocket envelopes together. Keep current paths and add fields/messages where feasible; a breaking mechanical contract gets explicit version negotiation. Existing clients must not trigger legacy mutation behavior for an SRD campaign accidentally. The terminal client either uses the same preparation/confirmation path or reports the unsupported capability safely.

## Rules authority and narrative combat

Pin the [official SRD 5.2.1 PDF](https://media.dndbeyond.com/compendium-images/srd/5.2/SRD_CC_v5.2.1.pdf), its content hash, stable rule/content IDs and source-section references during rules ingestion. Preserve the attribution specified by the document in imported content and the appropriate credits/export surfaces. Treat the SRD catalogue as the content boundary; a familiar character-app interaction is a usability reference, not a request to reproduce its assets or a separate paid catalogue.

Build a conformance register covering checks/saves, attacks and damage, action economy, movement/range, conditions and timing, spellcasting, rests, recovery, character creation, advancement and supported content. Register every SRD entry or rule group with its implementation status, deterministic fixtures and known gaps. Include levels 1–20, SRD classes/subclasses, monster/summon stat blocks and exploration/environment rules in the target; scope each implementation batch honestly.

The rules module owns computations, dice and resource effects. The game master owns fictional interpretation, NPC intent and rulings where the reference leaves a decision open. It cannot change a committed roll, invent resources or replace a supported rule with a generic success check. A ruling affecting mechanics is recorded with the action and applied consistently.

Narrative combat still needs authoritative state: initiative and active combatant, action economy, movement, relevant distances/reach, sight/cover, target relationships, conditions, concentration and effect timing. Show the player enough context to select an action without requiring a map. Do not silently treat vague Near/Far labels as an exact range check.

When distance, cover or area membership is unresolved, establish the scene facts or request clarification before resolving an action that depends on them. Persist affected target sets and spatial rulings so the next turn does not reinterpret them. Geometry may be represented internally when needed without introducing a grid UI.

Reactions and meaningful resource choices remain player decisions at the correct trigger. Rolls on confirmation removes an extra dice prompt; it does not authorize the game master to spend a reaction or limited feature silently. Handle applicable NPC and ongoing-effect actions with the same resolver and durable receipts. Keep undiscovered enemy stats, secret difficulties and hidden entities out of previews, receipts, suggestions and accessibility text.

## Contextual action assistance

Offer a small set of relevant actions and reminders, not an exhaustive optimal strategy. Suggestions can draw on the current scene, current objective, character capabilities/resources, equipment, encountered people and player-known facts.

- Ground every recommendation in known state and a current revision. Recheck an accepted suggestion before execution.
- Show a brief useful reason when relevant: a known clue, a proficiency or an available spell. Never imply that a guessed story outcome is guaranteed.
- Include a gentle Where were we? reminder on return and an explicit help/reminder action. Assist with forgotten context without spoiling unknown information.
- Render a short suggested completion as faint text only when the composer can accept it. Tab accepts into the draft; Escape dismisses. Typing remains unrestricted.
- Intercept Tab only for an actual visible completion while editing. Shift+Tab, normal focus traversal, input-method composition and screen-reader navigation keep working. Provide an explicit accessible Accept suggestion control as well.
- On phone, tapping a suggestion inserts it into the draft or opens action preparation; it never sends by itself. Selecting a choice follows the same explicit submission/confirmation semantics.
- Generate or rank candidates once per relevant scene/revision and match prefixes locally. Do not call a paid model on each keystroke. Failed or slow assistance must not block typing or action submission.
- Replace obsolete suggestions after a turn, resource change, campaign switch or reconnect. Campaign state and suggestion caches never mix across campaigns.

## Faithful chapters and journal

Keep three concepts distinct: the played transcript, short recaps used to resume/recall play, and readable story chapters. Reuse journal storage/generation where suitable, but give chapters explicit source coverage and lifecycle rather than simply relabeling the existing short summaries.

Generate from committed player-visible events, available dialogue and narration. Compression and connective prose can improve flow; new events, motivations presented as fact, invented dialogue, changed outcomes and hidden lore cannot be added. Ignore instructions embedded in campaign text when generating the retelling.

Default to third-person past tense. Chapters follow plot developments and scene boundaries rather than requiring one paragraph per turn. The reader has chapter navigation, reading-position retention, comfortable typography, a source-history link and export. Display how far play is covered and whether newer events await a chapter.

On-demand generation covers the next completed range. Once persisted scene boundaries are available, schedule background generation at appropriate boundaries. Existing short ten-turn summaries can support recaps and memory but do not determine chapter structure. Generation failure does not block play, and retry does not create duplicate chapters. Regeneration produces a revision; it does not silently replace an already-read chapter or the source transcript.

Automated coverage checks establish source ranges, names/outcomes and knowledge boundaries. Curated human-reviewed campaign fixtures assess plot fidelity and readability; a passing structural check alone is not proof that generated prose is faithful.

## Compatibility, migration and recovery

- Add a rules version independently of the existing `narrative`, `light` and `crunch` mode identifiers. Preserve current campaigns as legacy until explicitly migrated.
- Add data and backfill only where the original value maps unambiguously. Never fabricate class, spell, proficiency or resource history from generic ability prose.
- New SRD characters use the new schema and derived projection. Existing campaigns retain current history and mode behavior; reading/export/replay remain available throughout rollout.
- If migration is offered, first produce a report of unmapped/ambiguous fields and a proposed character. Migration requires explicit owner choices and a recoverable snapshot. This plan does not authorize bulk conversion.
- Roll out additive database/API changes before client dependence. Keep code rollback able to read compatible data; do not roll back a release by dropping newly written gameplay records.
- Campaign ownership (#28), ordering (#29), consistent failure/retry (#30), creation recovery (#33–#35) and reset consistency (#36) are release dependencies where those operations are exposed.
- Preserve existing local dirty work, other checkouts and runtime services during implementation. No deploy or production data mutation is implied by this plan.

## Ordered delivery units

Each row is a bounded implementation slice or an explicitly repeated content batch. Large rule catalogues are split into repeated class/feature/spell-family PRs; they are not one giant implementation PR. Dependencies refer to this table, and linked review issues are reused rather than duplicated.

| ID | Unit and affected boundaries | Depends on | Smallest useful acceptance check |
| --- | --- | --- | --- |
| P01 | Repair tagged integration compilation and add the database CI gate (#31, then #32). Two PRs if needed. | None | Tagged query/handler tests run and the required gate fails when a database assertion breaks. |
| P02 | Enforce campaign ownership for all nested resources and WS (#28); `internal/server`, auth/handlers and store entry points. | None | Two users cannot read, mutate or connect to each other's campaign. |
| P03 | Make action/turn ordering explicit and prevent duplicate numbering (#29); engine/state/migrations. | P01 | Barrier-controlled same-campaign submissions cannot save duplicate turns; other campaigns remain independent. |
| P04 | Establish durable action identity, receipts and failure/retry recovery (#30); engine/game/query/API boundaries. | P02, P03 | Injected failures and repeat confirmations cannot lose committed effects or apply them twice. |
| P05 | Add deterministic player-journey fixtures and frontend behavior commands (#45); frontend test harness and CI. | None | Existing creation, play, history, reconnect and replay paths have a repeatable baseline. |
| P06 | Produce desktop/phone concepts and select reading/layout details; existing CSS/tokens/AppShell/play composition. | Interview decisions | Demonstrate a screen-fitting gameplay frame with long transcript/choices, sheet action, combat and phone keyboard open; also review the chapter reader. No live gameplay edits required. |
| P07 | Introduce rules/content version metadata, pinned SRD source/attribution and the conformance register; domain/rules/content/query boundaries. | P01 | Source version/hash is reproducible and every planned rule/content family has an explicit coverage status. |
| P08 | Add typed character state and server-derived sheet projection, keeping legacy data separate; domain/game/API/schema. | P07 | A fixed character fixture returns reproducible derived values and legacy reads remain unchanged. |
| P09 | Implement the common deterministic mechanics used by the first slice: dice, checks/saves, proficiency and supported effect/cost primitives; rules/combat/game. | P07, P08 | Known-roll fixtures assert the pinned rule outcome and no model chooses arithmetic or resources. |
| P10 | Add read-only preparation and explicit confirmation contracts to REST/WS and engine, with recoverable narration. | P04, P08, P09 | Preview spends nothing; stale/invalid confirmation spends nothing; successful retries return one receipt and one set of dice. |
| P11 | Refine the viewport-fitting play shell, internal scroll regions and anchored composer; frontend layout/narrative/hooks. | P05, P06 | Main gameplay has no document overflow at supported default sizes; complete a turn and read old history with essential controls reachable, including with the phone keyboard open. |
| P12 | Ship a first full-stack vertical slice: a valid low-level martial character, one supported attack, one skill action, one usable item and corresponding sheet controls. | P08–P11 | Select/preview/confirm each action; inspect roll/cost/outcome; reload and replay the same saved result. Mark the build as partial support. |
| P13 | Add structured initiative/action economy and narrative spatial/condition state, then reaction handling; combat/game/API/combat UI. | P09, P10, P12 | An encounter fixture checks range, initiative, movement, effect timing and reaction choices without a grid. Split encounter core and reactions into separate PRs. |
| P14 | Ship a casting slice with prepared/known spells, slots, targets, components, concentration and recovery; rules/game/API/sheet. | P09, P10, P12 | A caster chooses targets and slot, confirms once, spends once and retains correct concentration/resource state after reload. |
| P15 | Expand creation to catalogue-backed valid SRD builds using the existing choice/interview flow; world/startup/frontend workflow. | P07–P10, P14 | Both creation methods produce the same validated character contract and preserve drafts through errors. |
| P16 | Fix creation lifecycle and required setup errors (#33–#35), reset consistency (#36), with separate focused PRs. | P01, P02 | Fail each persistence stage; no incomplete world is reported playable and no reset can partially succeed. |
| P17 | Expand rules/content in bounded batches: remaining character features, classes/subclasses, spell families, conditions, items, rests, levels, monsters/summons and exploration rules; conformance register drives the order. | P13–P16 | Each batch has source-linked fixtures, UI availability and persistence tests; unsupported coverage is never hidden. Advancement and multiclass rules get their own slices. |
| P18 | Add grounded suggestions and reminders, then ghost-text/touch acceptance; engine/prompt/context/API/composer. | P10–P14 | A relevant known ability/clue can be suggested; a hidden fact or unavailable resource cannot; Tab/tap never submits. |
| P19 | Make summary/worldgen routing and recall intentional (#38–#40); server/LLM/journal/memory. | P04, P07 | Request-recording stubs prove the selected route and embedding combination; auto-summary failure cannot fail a committed turn. Keep each issue focused. |
| P20 | Add durable chapter generation/source coverage, followed by the journal reader/export; journal/query/handlers/frontend. | P04, P11, P19 | A completed fixture becomes a faithful saved chapter; retries/revisions and source-history links work without blocking play. Backend and reader can be separate PRs. |
| P21 | Fix campaign-switch isolation (#42), stream cancellation (#37) and bounded live-event processing (#43), in focused PRs after reproducing each concern. | P04, P05 | Delayed history, disconnect and backpressure fixtures cannot cross campaign state or leave blocked workers. |
| P22 | Complete navigation/accessibility and narrow-screen qualification (#44, #46); all account/library/setup/play/journal/world/replay surfaces. | P11–P15, P18, P20 | Validate gameplay viewport fit, internal scrolling and anchored controls; keyboard, touch, screen-reader, zoom and phone-keyboard checks complete all player workflows. |
| P23 | Evaluate incremental narrative streaming separately (#41). | P04, P10, P21 | Progressive text, committed result and narration retries reconcile without duplicate history or mechanical effects. This does not block the readable transcript. |
| P24 | Qualify the declared rules/content coverage and end-to-end campaign lifecycle before recommending full-rules play. | Applicable earlier units, P17 register | The acceptance scenarios below pass; gaps are explicit; migration, failure and existing-mode checks pass. |

P05/P06 and read-only design work can proceed while backend foundations are repaired. The first interactive sheet slice depends on reliable action commitment; it does not wait for every spell or chapter feature. P17 is a series of bounded releases, and its coverage register determines when the full detailed-rules target is actually complete. No delivery duration is asserted without estimating those batches.

## Acceptance scenarios and verification

| Scenario | Boundary and falsifiable check |
| --- | --- |
| Create and resume | Create a martial and a caster through choice-based setup; validate an interview-produced build; fail and recover generation; reload with the same valid character, world and opening choices. |
| Fit main gameplay on screen | At desktop 1366×768, 1440×900 and 1920×1080, tablet 768×1024, and phone 360×800 and 390×844 default viewports, assert no document-level overflow. Use long transcript/choices, combat and open details. Scroll only the intended panel and keep essential controls reachable. Repeat phone checks with browser chrome changes and the software keyboard open; verify accessible reflow at increased text size/zoom. |
| Read and act | Read earlier prose, inspect character/world context, return to the same position/draft, accept a completion, then submit explicitly. Repeat with phone keyboard and keyboard-only desktop navigation. |
| Use an ability | Prepare a skill action from the sheet and the equivalent free text. With the same recorded circumstances and fixed dice, assert equivalent costs and outcomes. |
| Spend resources | Use a consumable and cast an upcast concentration spell where supported. Recheck availability, target/range and costs; retry confirmation and reconnect; resources change once. |
| Resolve combat | Run multiple combatants through initiative, movement, bonus action/reaction choices, conditions, concentration and death/recovery cases. Assert actual state and receipts, not prose that claims a rule ran. |
| Reject stale intent | Open two clients; spend or change the relevant resource in one. The other must refresh its preview rather than execute an invalid prepared action. |
| Fail safely | Fail before commit, after commit but before narration, during client disconnect, and during narration retry. Assert one durable outcome and recoverable history, with no rerolls. |
| Read the story | Generate chapters for a curated campaign with a failed plan, a corrected misunderstanding and a secret not yet learned. Verify plot/outcome fidelity, source coverage, known-information boundary and recoverable generation failure. |
| Recover a long session | Switch campaigns with delayed history, play a long event fixture, reconnect and replay. Assert campaign isolation, bounded transport bookkeeping and stable history. |
| Preserve legacy play | Existing rules modes, old campaigns, TUI paths, saves, export and replay keep their documented behavior. No generic old character is silently labeled SRD-compliant. |

Use deterministic model doubles and injectable dice for reproducible mechanical tests. Separate those results from live-provider prose/latency qualification. Generated chapter quality and suggestion usefulness need curated review in addition to structural tests.

Existing commands to retain and extend:

```sh
rtk proxy env GOMAXPROCS=4 go test -short -race -p 2 ./...
rtk proxy env GOMAXPROCS=4 go vet -p 2 ./...
rtk proxy env GOMAXPROCS=4 go test -tags integration -run Integration -p 1 -timeout 5m ./internal/handlers ./internal/state/sqlc
rtk proxy pnpm --dir frontend build
rtk proxy pnpm --dir frontend lint
```

P05 must introduce and document the actual frontend unit/interaction and end-to-end commands; no such command is claimed to exist today. Add shared contract-shape tests for API/WS/history and migration round-trip checks as those contracts change. During implementation, use browser review at the gameplay viewport sizes above, keyboard/screen-reader review, and controlled interruption/reconnect tests. Check document bounds, panel scroll ownership and anchored-control visibility with automated layout assertions; separately check actual phone keyboard/browser-chrome behavior because a resized desktop viewport alone does not prove it. Local checks, hosted CI, browser/device qualification and any eventual deployment are separate results.

## Remaining assumptions and next checkpoint

The interview has resolved the major product choices. Remaining defaults are explicit above; implementation should not silently expand them into multiplayer, a tactical map or a new brand.

Before implementing a new rules family, its coverage entry must settle exact supported content and any SRD rulings left open. Ambiguous legacy conversions remain unmigrated until resolved. These are bounded content decisions, not reasons to keep the initial product interview open indefinitely.

Next checkpoint: review the desktop/phone concepts and the P07 conformance register, then start the first vertical gameplay slice after the ownership and commitment foundations. This document is the plan; its acceptance scenarios have not yet been implemented or run.
