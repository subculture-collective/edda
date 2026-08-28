# Dummy turn prompts for rules/combat/state testing

Use these as disposable playtest prompts after creating a campaign. They are written to exercise different runtime states, tool groups, and rules modes. Prefer local dev only.

These are only player turn inputs. For full hypothetical LLM messages with serialized game state and tool sets, see [`assembled-prompt-fixtures.md`](./assembled-prompt-fixtures.md).

```bash
export EDDA_BASE_URL=http://localhost:18080
export EDDA_TOKEN=<local-test-token>
export CAMPAIGN_ID=<campaign-id>
python3 ~/.config/opencode/blacktower/managed-skills/play-edda/scripts/edda_api.py action "$CAMPAIGN_ID" "<prompt>"
```

## 1. Opening exploration / choices fallback

**Best state:** fresh campaign, no active quest yet.  
**Rules mode:** `narrative` or `light`.  
**Exercises:** scene description, `skill_check`, location/fact extraction, post-turn choice pass.

```text
I stop at the edge of the first strange landmark, study the ground, smell the air, and look for one safe path forward and one sign that someone or something passed this way recently.
```

Expected signals:
- Possible Perception/Survival/Investigation check.
- A useful clue or hazard may become a world fact.
- Response should return 3-5 `choices` from the post-turn choice pass.

## 2. D&D-style uncertain skill check

**Best state:** `rules_mode=light`, exploration or dungeon-like location.  
**Exercises:** DC selection, failure/success narration, `resolution_events` visibility.

```text
I use my rope and pitons to climb the crumbling black wall, moving slowly and testing each handhold before putting my weight on it.
```

Expected signals:
- Athletics or Acrobatics `skill_check`.
- `resolution_events` should expose roll, total, DC, success/failure.
- Narrative should not mention the DC or roll directly.

## 3. Narrative-rules low-mechanics resolution

**Best state:** `rules_mode=narrative`, social/mystery scene.  
**Exercises:** resolving without over-rolling; lore persistence.

```text
I sit across from the frightened witness, set my weapon aside, and ask them to tell me what they saw in whatever order the memory comes back.
```

Expected signals:
- May use Persuasion/Insight if stakes are uncertain, but should not force mechanics if trust is obvious.
- Durable revelations should become `world_fact` entries.
- If the witness creates a concrete goal, quest extraction may create/update a quest.

## 4. Quest creation from a clear objective

**Best state:** no active quest or an obvious untracked objective.  
**Exercises:** post-turn quest extraction, quest objectives, related entity links.

```text
I accept the elder's request: I will find the three missing ward-stones, learn which one is cracked, and bring them back before nightfall so the village barrier can be restored.
```

Expected signals:
- `create_quest` should create a short-term quest.
- Quest should include ordered objectives, not just a title/description.
- If `GET /quests` shows `objectives: []` while the state-change payload had objectives, diagnose quest read-model/persistence.

## 5. Quest progress / objective completion

**Best state:** active quest with at least one objective.  
**Exercises:** `update_quest` or `complete_objective`, avoiding duplicate quests.

```text
I place the first ward-stone into my pack, mark its location on my map, and compare its crack pattern against the elder's description so I can tell whether this is the damaged one.
```

Expected signals:
- Existing quest should update rather than duplicate.
- Objective may complete or a new sub-objective may be added.
- Inventory may update if the ward-stone becomes a carried item.

## 6. Inventory gain and item use

**Best state:** exploration near a discovered object.  
**Exercises:** `create_item`/`add_item`, item properties, durable state-change projection.

```text
I wrap the silver thorn in cloth, add it to my pack, and note exactly where I found it and what it was touching.
```

Expected signals:
- Item should be created/added to inventory.
- A fact may record where/why it matters.
- Later inventory fetch should show the item.

## 7. NPC creation / dialogue / relationship

**Best state:** entering a populated or haunted location.  
**Exercises:** `create_npc`, `npc_dialogue`, relationship or awareness state.

```text
I call out softly to whoever is hiding behind the shrine: "I am not here to take anything. Tell me your name and what frightened you."
```

Expected signals:
- If a new NPC appears, `create_npc` should persist them.
- If they reveal durable lore, facts should persist.
- Encountered NPCs panel/API should include the NPC if the player meets them.

## 8. Combat initiation

**Best state:** location with hostile creature or imminent attack.  
**Rules mode:** `crunch`.  
**Exercises:** `initiate_combat`, combat state, tactical narration.

```text
When the bone-wolf lunges from the mist and blocks the only bridge, I draw my blade, brace behind my shield, and try to hold it back before it reaches the wounded scout.
```

Expected signals:
- Should call `initiate_combat` if threat is immediate.
- `combat_active` should become true.
- Narrative should establish positions and immediate stakes.

## 9. Combat round / attack or defense

**Best state:** `combat_active=true`.  
**Rules mode:** `crunch`.  
**Exercises:** combat round tools, HP/status changes, enemy tactics.

```text
I feint left, strike at the bone-wolf's exposed ribs, then step between it and the scout so it has to go through me.
```

Expected signals:
- Attack/action should resolve through combat or skill tools.
- HP/status may change via durable tools.
- Combat should remain active unless clearly resolved.

## 10. Retreat / combat resolution

**Best state:** `combat_active=true`, dangerous fight.  
**Exercises:** ending combat, movement, consequences.

```text
This fight is turning bad. I grab the scout by the shoulder, throw down my last ration bundle to distract the beast, and retreat toward the marked standing stone.
```

Expected signals:
- May remove/use an item.
- May move player or resolve combat if escape succeeds.
- Failure should produce consequences without silently killing the character.

## 11. Rest and recovery

**Best state:** wounded or after intense scene.  
**Exercises:** `rest`, time advancement, HP/status recovery.

```text
I make a hidden camp under the fallen cedar, bind my wounds, eat a ration, and rest until the forest sounds normal again.
```

Expected signals:
- Should use `rest` or `advance_time` if meaningful time passes.
- Inventory may decrement rations if tracked.
- HP/status may recover depending on rules mode.

## 12. Time pressure / countdown

**Best state:** active quest with a deadline.  
**Exercises:** campaign time, quest update, consequences.

```text
I ignore the safer road and take the flooded shortcut because the ward will fail at dusk if I do not return in time.
```

Expected signals:
- May use Survival/Athletics check.
- Should advance time or mark time pressure in quest state/facts.
- Should present a meaningful consequence for speed vs safety.

## 13. Fact revision / contradiction pressure

**Best state:** established fact later challenged by evidence.  
**Exercises:** `revise_fact` vs duplicate facts; consistency.

```text
I compare the elder's story with the inscription on the ward-stone. If they contradict each other, I trust the inscription and revise what we know about who broke the barrier.
```

Expected signals:
- Should not create duplicate contradictory facts.
- Should revise or qualify existing lore if evidence supports it.

## 14. Save/resume continuity check

**Best state:** after several turns with facts, quests, and inventory.  
**Exercises:** memory/state recall after reload. Use after refreshing UI or fetching history.

```text
Before acting, I review what I already know: current quest, known hazards, carried items, and the last promise I made. Then I choose the most urgent next step.
```

Expected signals:
- Narrative should reference actual active quest/facts/items.
- Should not hallucinate old choices or invent completed objectives.

## 15. Failure-with-progress check

**Best state:** any investigation scene.  
**Exercises:** failed checks should still move fiction forward.

```text
I search the ruined shrine for the hidden compartment the map hinted at, but I do it quickly because the mist is getting thicker.
```

Expected signals:
- On failure, the turn should still reveal a cost, complication, or partial clue.
- Avoid dead-end `state_changes: []` unless no durable change truly occurred.
