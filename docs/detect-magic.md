# Detect Magic through Magic Initiate (#265)

The fixed Human Acolyte Cleric records Magic Initiate (Cleric; Wisdom) with
Detect Magic and one free use (`stats.srd521.casting.choices.magic_initiate_uses`).
This slice makes only that free use playable, outside encounters. Slot casting,
Ritual casting and every other feat stay unavailable. Source: SRD 5.2.1 p. 87
(Magic Initiate) and p. 123 (Detect Magic).

## Actions

| Intent | Cost | Effect |
| --- | --- | --- |
| `spell-cast`, `entry_id: detect-magic`, no `slot_level` | Magic Initiate use 1→0, Magic action, no slot | Starts ten-minute concentration with the action ID; replaces current concentration and unlinks a replaced Shield of Faith recipient. The receipt says only whether adjudicated magic is sensed. |
| `detect-magic-study` | Magic action | Needs live Detect Magic concentration and sensed magic. Shows a faint aura on each visible sensed bearer, with the school only for spell-created effects. |

A `slot_level` on the cast is refused. Without a recorded choice, or with the
use spent, the preview gives the reason. A Long Rest restores the use; ending
concentration, damage saves, expiry and rests reuse the #77 lifecycle.

## TUI

The commands act for the campaign owner's first active character, through the
same preview and confirm flow as the web client.

| Command | Effect |
| --- | --- |
| `/detect-magic` | Shows any live concentration and its expiry, or why a study is unavailable. Writes nothing. |
| `/detect-magic cast` | Previews the free-use cast. |
| `/study-detect-magic` | Previews the study. |
| `/confirm-detect-magic` | Commits the preview. A retry reuses the action ID. |
| `/cancel-detect-magic` | Drops the preview. |
| `/end-concentration <caster UUID>` | Previews ending the caster's current concentration, whatever the spell. `/confirm-faith` commits it and `/cancel-faith` drops it. |

## Scene facts

Sensing reads only `npcs.properties.srd521_magic` on NPCs at the caster's
location that already appear in the session log. The `srd521_` prefix is
reserved for trusted fixtures, so model-created NPCs cannot set it. Every key
is required and non-null; anything else is skipped and never disclosed:

```json
{"profile":"detect-magic-v1","actor_ids":["<pc id>"],"location_id":"<location id>",
 "hidden":false,"distance_feet":20,"barrier":"none","visible":true,
 "origin":"spell","school":"abjuration"}
```

- `hidden: true` is rejected: hidden lore is never sensed.
- `barrier` is `none` or one of the source blockers `stone-dirt-or-wood-1-foot`,
  `metal-1-inch`, `lead-thin-sheet`; a blocker hides the bearer.
- Bearers beyond 30 feet are not sensed.
- `origin` is `spell` with a school, or `other` with an empty school.

When nothing is sensed, the receipt says that this does not prove that no magic
is present. Objects without NPC facts, encounter use and other senses remain
unsupported. Preview fingerprints the facts, so a changed ruling requires a
fresh preview.

## Checks

- `internal/rules/srd521/detect_magic_test.go`: source-record agreement, fact
  validation, cast, study, expiry, Long Rest and damage.
- `internal/handlers/detect_magic_integration_test.go`: PostgreSQL cast and
  study receipts, replacement, replay, rollback, stale facts, reload,
  missing or spent feat use, and the TUI preview.
- `tui/detect_magic_test.go`: command parsing, cast retry, status, study and
  cancel.

This work includes material from the System Reference Document 5.2.1 (“SRD
5.2.1”) by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd,
licensed under CC-BY-4.0 (https://creativecommons.org/licenses/by/4.0/legalcode).
Spell facts are adapted into Edda's structured rules.
