# Rogue and Thief decision checkpoint (#329)

PENDING parent reservation, complete item/effect scope and POST-RANGER release; NOT publication or sealed proof.

This decision proposes bounded execution; it implements no Rogue runtime. Catalogue completion is source evidence only. No new native issues or dependencies have been published. #329 remains open. The item/scroll proposal requires the parent’s scope reservation; it is not yet the complete publication acceptance proof.

Source: English SRD 5.2.1, Rogue PDF61–63, Thief64, advancement23, multiclass24–26; SHA256 `8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87`. Source-derived records retain the full CC-BY-4.0 attribution and adaptation notice in the canonical catalogues and docs/srd-5.2.1-conformance.md. This plan and ownership matrix are Edda additions.

## All twenty rows

Class PB is preserved source metadata. Live PB and XP use the sum of authoritative admitted post-transition class levels. Fighter4/Rogue1→Rogue2 is total6,14000XP/PB3; Fighter4→Rogue1 is total5,6500XP/PB3. HP/features/Sneak Attack use Rogue class level.

| Rogue level | Class-table PB | Sneak Attack | Required new behavior owners |
|---|---|---|---|
| 1 | 2 | 1d6 | expertise, sneak-attack, languages, mastery |
| 2 | 2 | 1d6 | cunning-action |
| 3 | 2 | 2d6 | subclass, steady-aim, fast-hands-check, fast-hands-object, climber, jumper |
| 4 | 2 | 2d6 | improvement |
| 5 | 3 | 3d6 | strike, poison, trip, withdraw, uncanny-dodge |
| 6 | 3 | 3d6 | expertise |
| 7 | 3 | 4d6 | evasion, reliable-talent |
| 8 | 3 | 4d6 | improvement |
| 9 | 4 | 5d6 | stealth-attack |
| 10 | 4 | 5d6 | improvement |
| 11 | 4 | 6d6 | improved-strike |
| 12 | 4 | 6d6 | improvement |
| 13 | 5 | 7d6 | attunement, charges, scroll-check, scroll-light, scroll-detect |
| 14 | 5 | 7d6 | daze, knock-out, obscure |
| 15 | 5 | 8d6 | slippery-mind |
| 16 | 5 | 8d6 | improvement |
| 17 | 6 | 9d6 | reflexes |
| 18 | 6 | 9d6 | elusive |
| 19 | 6 | 10d6 | boon |
| 20 | 6 | 10d6 | luck |

## Runtime evidence

- `internal/rules/srd521/creation.go`: NewSupportedBuild fixed Fighter/Cleric admission only; no Rogue.
- `internal/rules/srd521/character.go`: Project class/level guard excludes Rogue; class label is insufficient.
- `internal/rules/srd521/catalogue.go`: ValidateReadiness requires implemented mechanic keys; metadata does not grant runtime.
- `internal/game/advancement_preview.go`: Fixed Fighter1->2 preview; no Rogue HP/eligibility adapter.
- `internal/game/spell_preview.go`: Fixed Cleric cast seams; no Thief scroll admission.
- `internal/engine/action_confirmation.go`: Durable receipt/transaction seam; no inferred Rogue feature execution.
- `internal/combat/resolver.go`: Shared combat arithmetic seam; no source-key Rogue entitlements.
- `internal/handlers/character_rules.go`: Shared sheet read route; projection guard still controls readiness.
- `frontend/src/components/start/CharacterGuidedForm.tsx`: Class selector labels do not qualify Rogue creation.
- `frontend/src/components/character/ActionPreviewPanel.tsx`: Shared confirmation UI does not qualify Rogue action options.
- `tui/app.go`: Trusted command path consumes shared engine; no separate Rogue entitlement.

## Proposed one-PR terminals

Preview inert: no writes or dice. Explicit confirmation revalidates source,actor,campaign,revision,target,range,timing,cost and qualification. Injectable dice; atomic choice/effect/resource/history persistence. Retry/replay/reload/reconnect returns same receipt,conflicting payload rejects,concurrent confirmations commit ONCE. Inject failures for complete isolated PostgreSQL rollback. Web AND TUI expose supported choices/unavailable reasons; accessible cancel/confirm/reconnect; changed web requires desktop/phone fixtures,build and focused lint. Automatic effects belong to triggering confirmed event,not optional new activation. Unknown/private facts unavailable without leaking hidden reasons. Short tests are not integration-tag PostgreSQL/browser/device proof. No production/inference/runtime promotion. Class-total levels integer1-20,not booleans/fractions; each admitted class positive,post-transition total<=20; later-class variants require own multiclass eligibility,first-character variants do not acquire a multiclass prerequisite at runtime. Spell Recall needs actual Spellcasting; ordinary feats have no blanket casting/encounter requirement. Slot-cast one/turn applies only to slot expenditure TO CAST; authoritative confirmed invalid target spends slot/no effect with source-private filtering; unauthorized/unverifiable facts unavailable.

### hp (class, Task Size M)

PDF23/61: own D8/fixed5 adapter ONLY, qualification968 -> transaction429. First character8+CON; later level max(1,d8+CON) or max(1,5+CON). +2 CON fixed gives7; die1 CON-3 gives1; CON17->18 at total8 adds8 retroactively without healing. Reject D10/fixed6, reroll, first8 on later entry.

Existing prerequisites: [968, 58]; proposed owners: []. UNPUBLISHED.

### traits (class, Task Size M)

PDF61: DEX primary; DEX/INT saves; Simple and Martial with Finesse OR Light; Thieves’ Tools and Light armor. Reject Heavy armor, martial neither property, later-class saves.

Existing prerequisites: [430, 333, 58]; proposed owners: []. UNPUBLISHED.

### skills (class, Task Size M)

PDF61: choose exactly4 distinct Acrobatics/Athletics/Deception/Insight/Intimidation/Investigation/Perception/Persuasion/Sleight of Hand/Stealth. Positive first4; negative3/5/duplicate/Survival.

Existing prerequisites: [384, 430]; proposed owners: []. UNPUBLISHED.

### equipment (class, Task Size M)

PDF61: A Leather Armor,2 Daggers,Shortsword,Shortbow,20 Arrows,Quiver,Thieves’ Tools,Burglar’s Pack,8GP OR B100GP. Positive exact A or B; reject19arrows,9GP,mixed/repeated grants. Inventory grant only; no inferred item effects.

Existing prerequisites: [195, 430]; proposed owners: []. UNPUBLISHED.

### multiclass (class, Task Size M)

PDF24-25/61: own DEX13 and EVERY admitted current-class primary minimum13. DEX12 fails independently; existing Fighter STR12/DEX12 fails. Gain D8,one Rogue skill,Thieves’ Tools,Light armor and level1 features; no saves,weapon/equipment/four-skill regrant. Positive Fighter4 DEX13 -> Rogue1; reject unqualified vector.

Existing prerequisites: [335, 430]; proposed owners: ['hp']. UNPUBLISHED.

### expertise (class, Task Size M)

PDF61: exactly2 EXISTING skill proficiencies at1 and2 MORE at6; PB2 expertise adds4 not6, PB3 adds6. Reject tools,unproficient/duplicate skills,third selection,grant6 at5.

Existing prerequisites: [384, 430, 61]; proposed owners: []. UNPUBLISHED.

### sneak-attack (class, Task Size M)

PDF61-62: once PER TURN creature hit,Finesse OR Ranged weapon; Advantage OR ally within5ft/not Incapacitated AND no Disadvantage. Same weapon damage type; ceil(Rogue level/2)d6. Positive Rogue5 3d6,own turn and different off-turn; reject same turn repeat,ally6ft/incapacitated,disadvantage,unarmed/nonqualifying weapon,noncreature. Inject3/4/5=12 extra; no off-hand conflation.

Existing prerequisites: [487]; proposed owners: []. UNPUBLISHED.

### languages (class, Task Size M)

PDF62/18: Thieves’ Cant plus1 from actual standard/rare language tables. Positive Common Sign Language; reject unknown/duplicate/two extra; preserve independent provenance.

Existing prerequisites: [384, 430]; proposed owners: []. UNPUBLISHED.

### mastery (class, Task Size M)

PDF62: two proficient weapon kinds; completed Long Rest may replace. Dagger/Shortbow positive; reject3,duplicate,unproficient,rest start/interruption. Consume selected property owner489-496,not label.

Existing prerequisites: [384, 428, 430, 489, 490, 491, 492, 493, 494, 495, 496]; proposed owners: []. UNPUBLISHED.

### cunning-action (class, Task Size M)

PDF62: on OWN turn Bonus Action Dash/Disengage/Hide; one choice. Reject off-turn,spent Bonus Action,automatic extra action. Speed30 Dash adds30 movement; consumes780/781 and proposed hide output.

Existing prerequisites: [780, 781]; proposed owners: ['hide']. UNPUBLISHED.

### subclass (class, Task Size M)

PDF62/64: choose exactly Thief at Rogue3; no subclass at2,second choice or foreign subclass; feature grant class3/9/13/17 only.

Existing prerequisites: [384]; proposed owners: []. UNPUBLISHED.

### steady-aim (class, Task Size M)

PDF62: Bonus Action only if NO movement this turn; Advantage NEXT attack on CURRENT turn; Speed0 remainder turn. Positive no move -> attack; reject moved1ft,second attack benefit,next turn carryover.

Existing prerequisites: []; proposed owners: []. UNPUBLISHED.

### improvement (class, Task Size M)

PDF63/86-88: class4/8/10/12/16 ASI OR qualifying feat. +2one/+1two,max20. Positive18->20; reject21,duplicate,unqualified,wrong class level. NO blanket Spellcasting,encounter or Concentration prerequisite for ordinary selection.399 Spell Recall needs ACTUAL Spellcasting;386 Magic Initiate and scroll access do not supply it. Selected feat requires actual relevant487/427/428/778/779/962 execution,not269/272/270/268 planning closure.

Existing prerequisites: [384, 389, 385, 386, 387, 388, 389, 390, 391, 392, 393, 394, 395, 396, 397, 398, 399, 400, 401, 487, 427, 428, 778, 779, 962]; proposed owners: []. UNPUBLISHED.

### strike (class, Task Size M)

PDF63: Rogue5 when Sneak Attack damage occurs choose1 effect; remove cost dice BEFORE rolling,resolve AFTER attack damage; DC8+DEX+livePB. DEX3/PB3 DC14,3d6 minus1=2d6. Reject miss,no Sneak Attack,negative dice,unsupported effect.

Existing prerequisites: []; proposed owners: ['sneak-attack']. UNPUBLISHED.

### poison (class, Task Size M)

PDF63: cost1d6,Poisoner’s Kit ON PERSON; CON save failure Poisoned1minute; repeat end each target turn success ends own instance. DC14 fail13/succeed14; reject kit absent,client cure,clearing other poison,expiry61seconds.

Existing prerequisites: [504, 727]; proposed owners: ['strike']. UNPUBLISHED.

### trip (class, Task Size M)

PDF63: cost1d6 target Large OR SMALLER,DEX save fail Prone. DC14 roll13 vs14; Huge rejects,no cost on invalid target. Consume488 source-owned Prone.

Existing prerequisites: [488]; proposed owners: ['strike']. UNPUBLISHED.

### withdraw (class, Task Size M)

PDF63: cost1d6 immediately AFTER attack move up to HALF Speed without Opportunity Attacks; Speed30 allows15,not16; no teleport or next-turn movement credit.

Existing prerequisites: []; proposed owners: ['strike']. UNPUBLISHED.

### uncanny-dodge (class, Task Size M)

PDF63: Reaction when attacker YOU SEE hits with attack roll; halve attack damage round DOWN. 9->4; reject unseen,save-only damage,no reaction,miss,whole round resistance.

Existing prerequisites: []; proposed owners: []. UNPUBLISHED.

### evasion (class, Task Size M)

PDF63: DEX save for HALF damage only; success0/fail HALF; cannot use Incapacitated. Damage11 ->0 or5. Reject CON save,no-half effect,incapacitated,ASI ownership.

Existing prerequisites: [487]; proposed owners: []. UNPUBLISHED.

### reliable-talent (class, Task Size M)

PDF63: ability check using skill OR tool proficiency treats d20<=9 as10. Die9->10 plus own modifiers;10 stays10; reject unproficient,save,attack,change die11.

Existing prerequisites: [430, 61]; proposed owners: []. UNPUBLISHED.

### improved-strike (class, Task Size M)

PDF63: Rogue11 up to2 strike effects,pay EACH before rolling;6d6 minus1minus1=4d6. Reject3 effects,grant at10,free second effect.

Existing prerequisites: []; proposed owners: ['strike', 'poison', 'trip', 'withdraw']. UNPUBLISHED.

### daze (class, Task Size M)

PDF63: Rogue14 cost2d6 CON save; failure NEXT turn only ONE of move/action/Bonus Action.7d6->5d6; reject two categories,next-next turn,wrong save. No invented Stunned condition.

Existing prerequisites: []; proposed owners: ['strike']. UNPUBLISHED.

### knock-out (class, Task Size M)

PDF63: Rogue14 cost6d6 CON fail Unconscious1minute OR ANY damage; repeat save end each turn success ends own instance.7d6->1d6; reject five dice affordability,damage not ending,foreign instance removal.

Existing prerequisites: [487]; proposed owners: ['strike']. UNPUBLISHED.

### obscure (class, Task Size M)

PDF63: Rogue14 cost3d6 DEX fail Blinded until END next target turn.7d6->4d6; reject CON save,early start-turn end,foreign condition removal.

Existing prerequisites: [960, 963]; proposed owners: ['strike']. UNPUBLISHED.

### slippery-mind (class, Task Size M)

PDF63: Rogue15 adds WIS and CHA save proficiency using livePB. PB5 bonus5 not10; reject level14,duplicate stacking,INT replacement.

Existing prerequisites: [430]; proposed owners: []. UNPUBLISHED.

### elusive (class, Task Size M)

PDF63: Rogue18 no attack roll has Advantage AGAINST self unless Incapacitated. Positive unseen attacker loses advantage; negative Incapacitated allows normal calculation,save treated as attack.

Existing prerequisites: [487]; proposed owners: []. UNPUBLISHED.

### boon (class, Task Size M)

PDF63/87-88: Rogue19 Epic Boon OR other qualified feat,Night Spirit recommendation only. Positive chosen qualified boon; reject auto-Night-Spirit,level18,unqualified feat;399 requires real Spellcasting.

Existing prerequisites: [384, 389, 385, 386, 387, 388, 389, 390, 391, 392, 393, 394, 395, 396, 397, 398, 399, 400, 401, 487, 427, 428, 778, 779, 962]; proposed owners: []. UNPUBLISHED.

### luck (class, Task Size M)

PDF63: Rogue20 FAILED D20 Test -> roll20; once until completed Short OR Long Rest. Positive failed die2 ->20 (then ordinary modifiers); reject success,unknown outcome,second use,rest start/interrupted rest,automatic natural20 critical assumption.

Existing prerequisites: [427, 428]; proposed owners: []. UNPUBLISHED.

### fast-hands-check (class, Task Size M)

PDF64: Thief3 Bonus Action DEX(Sleight of Hand) lock/trap with Thieves’ Tools OR pocket. Check uses supplied trusted object DC15,14fails/15succeeds; reject no tools for lock/trap,unknown object facts,foreign pocket,spell substitution. Requires separately proposed bounded lock/trap/pocket outputs.

Existing prerequisites: []; proposed owners: ['subclass', 'object-check']. UNPUBLISHED.

### fast-hands-object (class, Task Size M)

PDF64: Thief3 Bonus Action Utilize OR Magic to USE ITEM requiring that action; not arbitrary spell Magic action. Positive named qualified item route only; reject wrong activation,unsupported item,spent Bonus Action; exact mundane/magic outputs prerequisite.

Existing prerequisites: []; proposed owners: ['subclass', 'utilize', 'magic-item']. UNPUBLISHED.

### climber (class, Task Size M)

PDF64: Thief3 Climb Speed=current Speed; Speed30 ->30,Speed0->0. Reject flat30,Swim grant; dynamic modifiers retained.

Existing prerequisites: []; proposed owners: ['subclass']. UNPUBLISHED.

### jumper (class, Task Size M)

PDF64/183/185: Thief3 optional DEX rather than STR jump distance. STR10/DEX16 running long jump16 vs10; high jump3+DEXmod=6 vs3; reject adding both,flight,waiving distance/movement/run-up rules.

Existing prerequisites: []; proposed owners: ['subclass']. UNPUBLISHED.

### stealth-attack (class, Task Size M)

PDF64: Thief9 Cunning Strike1d6; Hide’s Invisible survives THIS attack only if end turn behind Three-Quarters OR Total Cover.5d6->4d6; reject Half Cover,wrong source Invisible,no Hide,no end-turn check; save original Hide source and independently end on other triggers.

Existing prerequisites: [962]; proposed owners: ['strike', 'hide', 'subclass']. UNPUBLISHED.

### attunement (class, Task Size M)

PDF64/102: Thief13 up to FOUR simultaneous,not five. Positive fourth eligible item; reject fifth,foreign/duplicate item,ignoring item-specific restrictions. No all-item activation.

Existing prerequisites: []; proposed owners: ['subclass', 'attune-item']. UNPUBLISHED.

### charges (class, Task Size M)

PDF64: Thief13 on item PROPERTY expending charges roll1d6;6 skips charge spend,1-5 ordinary spend; preserve action cost and property effect. Positive6 keeps2charges,5 spends1of2; reject reroll/recharge/waived eligibility; finite proposed Wand fixture only.

Existing prerequisites: []; proposed owners: ['subclass', 'magic-item']. UNPUBLISHED.

### scroll-check (class, Task Size M)

PDF64: Thief13 any Spell Scroll,INT casting ability. Cantrip/level1 reliable; higher INT(Arcana) DC10+SPELL LEVEL; level2 DC12 (11fail disintegrates,12casts),level9 DC19. No Spellcasting grant; no invented class-list restriction. Real spell effect must separately qualify; broad339spells remain unavailable.

Existing prerequisites: []; proposed owners: ['subclass', 'scroll-item']. UNPUBLISHED.

### scroll-light (class, Task Size M)

PDF64/144/104-106: compose Thief scroll adapter with actual967 Light,INT; no learned/prepared/Spellcasting grant. Reject other-holder object,Huge,foreign scroll,slot spending or ritual/upcast. Own ongoing Light dismissal NO ACTION when not Incapacitated; Incapacitated alone does not end nonconcentration Light.

Existing prerequisites: [967]; proposed owners: ['scroll-check']. UNPUBLISHED.

### scroll-detect (class, Task Size M)

PDF64/123/104-106: compose scroll adapter with906 Detect Magic,INT and spell concentration/visibility/barriers only. Reject invisible-object aura disclosure; no slot spend or Ritual from consuming scroll. Prepared Ritual elsewhere normal time+10min,no slot,no Ritual Adept,NOritual upcast; >=1min cast Magic each turn/casting Concentration; interrupted no slot spend.

Existing prerequisites: [906]; proposed owners: ['scroll-check']. UNPUBLISHED.

### reflexes (class, Task Size M)

PDF64: Thief17 TWO turns FIRST round every combat,normal Initiative and minus10. Initiative15 ->15and5; reject later round,reload granting third,legacy surprise gate; equal initiative ordering deterministic; each is actual distinct turn with action economy and per-turn features.

Existing prerequisites: []; proposed owners: ['subclass']. UNPUBLISHED.

### hide (shared-action, Task Size M)

PDF183: DC15 DEX(Stealth) while Heavily Obscured or Three-Quarters/Total Cover,out of enemy sight; save check total as discovery DC,Invisible source. Ends on louder-than-whisper sound,enemy finds,attack roll,Verbal spell.14fails/15succeeds; Half Cover alone rejects. Supreme Sneak owns its one attack exception.

Existing prerequisites: [962]; proposed owners: []. UNPUBLISHED.

### object-check (shared-item, Task Size M)

PDF98/94/64: ONLY owned Lock/Thieves’ Tools lock operation,one trusted mechanical trap disarm,and one trusted pocket fixture; retain DEX(Sleight of Hand) Fast Hands exception. DC15 14/15 controls; no arbitrary traps/catalogue inference; unknown DC/ownership unavailable.

Existing prerequisites: [195, 430]; proposed owners: []. UNPUBLISHED.

### utilize (shared-item, Task Size M)

PDF99/191: ONLY owned Signal Whistle via Utilize, audible up to600ft. Positive trusted listener600ft; negative601ft,unknown acoustic facts,foreign whistle,spell action,invented consumption. Tinderbox exposed-fuel ignition is already Bonus Action onPDF100 and is NOT this Utilize fixture. Fast Hands supplies Bonus Action exception,not this owner.

Existing prerequisites: [195]; proposed owners: []. UNPUBLISHED.

### magic-item (shared-item, Task Size M)

PDF magic item Wand of Magic Detection: ONLY its Detect Magic property,3charges,Magic action,1charge; actual906 effect,source-specific dawn recharge separate and UNAVAILABLE here. Reject0charges,foreign/unattuned as source requires,other wand property. Thief conserves charge only,never restores.

Existing prerequisites: [195, 906]; proposed owners: []. UNPUBLISHED.

### attune-item (shared-item, Task Size M)

PDF102: ONLY four distinct owned source identities: Cloak of Protection PDF216,Gauntlets of Ogre Power PDF223,Ring of Jumping PDF238,Ring of Protection PDF238. Attunement protocol ONLY,not their magical effects; Short Rest focus one item,normal cap3,Thief adapter4. Reject unknown restrictions,5items,foreign possession,rest interruption; reject duplicate Ring of Protection type and using the identification Short Rest as the attunement rest; item magical effects remain separately unavailable without actual owners.

Existing prerequisites: [195, 427]; proposed owners: []. UNPUBLISHED.

### scroll-item (shared-item, Task Size M)

PDF Spell Scroll: ONLY Light/Detect Magic canonical scroll identities; readable ownership and consume/disintegration atomic with effect/check receipt; reject foreign/unknown/duplicate spend. Higher-level check uses isolated supplied scroll level,not granted unsupported effect. All other339effect coverage excluded pending parent decision.

Existing prerequisites: [195]; proposed owners: []. UNPUBLISHED.

### level-1 (class, Task Size M)

PDF62 row1: first-character1 or adjacent Rogue0->1 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/expertise", "srd-5.2.1:character/feature/rogue/sneak-attack", "srd-5.2.1:character/feature/rogue/thieves-cant", "srd-5.2.1:character/feature/rogue/weapon-mastery"], "proficiency_bonus": 2}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'traits', 'skills', 'equipment', 'expertise', 'sneak-attack', 'languages', 'mastery']. UNPUBLISHED.

### level-2 (class, Task Size M)

PDF62 row2: first-character1 or adjacent Rogue1->2 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/cunning-action"], "proficiency_bonus": 2}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-1', 'cunning-action']. UNPUBLISHED.

### level-3 (class, Task Size M)

PDF62 row3: first-character1 or adjacent Rogue2->3 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/rogue-subclass", "srd-5.2.1:character/feature/rogue/steady-aim"], "proficiency_bonus": 2}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-2', 'subclass', 'steady-aim', 'fast-hands-check', 'fast-hands-object', 'climber', 'jumper']. UNPUBLISHED.

### level-4 (class, Task Size M)

PDF62 row4: first-character1 or adjacent Rogue3->4 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/ability-score-improvement"], "proficiency_bonus": 2}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-3', 'improvement']. UNPUBLISHED.

### level-5 (class, Task Size M)

PDF62 row5: first-character1 or adjacent Rogue4->5 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/cunning-strike", "srd-5.2.1:character/feature/rogue/uncanny-dodge"], "proficiency_bonus": 3}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-4', 'strike', 'poison', 'trip', 'withdraw', 'uncanny-dodge']. UNPUBLISHED.

### level-6 (class, Task Size M)

PDF62 row6: first-character1 or adjacent Rogue5->6 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/expertise"], "proficiency_bonus": 3}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-5', 'expertise']. UNPUBLISHED.

### level-7 (class, Task Size M)

PDF62 row7: first-character1 or adjacent Rogue6->7 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/evasion", "srd-5.2.1:character/feature/rogue/reliable-talent"], "proficiency_bonus": 3}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-6', 'evasion', 'reliable-talent']. UNPUBLISHED.

### level-8 (class, Task Size M)

PDF62 row8: first-character1 or adjacent Rogue7->8 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/ability-score-improvement"], "proficiency_bonus": 3}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-7', 'improvement']. UNPUBLISHED.

### level-9 (class, Task Size M)

PDF62 row9: first-character1 or adjacent Rogue8->9 ONLY. Exact progression={"feature_keys": ["srd-5.2.1:mechanic/rogue/subclass-feature"], "feature_refs": [], "proficiency_bonus": 4}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-8', 'stealth-attack']. UNPUBLISHED.

### level-10 (class, Task Size M)

PDF62 row10: first-character1 or adjacent Rogue9->10 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/ability-score-improvement"], "proficiency_bonus": 4}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-9', 'improvement']. UNPUBLISHED.

### level-11 (class, Task Size M)

PDF62 row11: first-character1 or adjacent Rogue10->11 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/improved-cunning-strike"], "proficiency_bonus": 4}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-10', 'improved-strike']. UNPUBLISHED.

### level-12 (class, Task Size M)

PDF62 row12: first-character1 or adjacent Rogue11->12 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/ability-score-improvement"], "proficiency_bonus": 4}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-11', 'improvement']. UNPUBLISHED.

### level-13 (class, Task Size M)

PDF62 row13: first-character1 or adjacent Rogue12->13 ONLY. Exact progression={"feature_keys": ["srd-5.2.1:mechanic/rogue/subclass-feature"], "feature_refs": [], "proficiency_bonus": 5}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-12', 'attunement', 'charges', 'scroll-check', 'scroll-light', 'scroll-detect']. UNPUBLISHED.

### level-14 (class, Task Size M)

PDF62 row14: first-character1 or adjacent Rogue13->14 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/devious-strikes"], "proficiency_bonus": 5}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-13', 'daze', 'knock-out', 'obscure']. UNPUBLISHED.

### level-15 (class, Task Size M)

PDF62 row15: first-character1 or adjacent Rogue14->15 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/slippery-mind"], "proficiency_bonus": 5}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-14', 'slippery-mind']. UNPUBLISHED.

### level-16 (class, Task Size M)

PDF62 row16: first-character1 or adjacent Rogue15->16 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/ability-score-improvement"], "proficiency_bonus": 5}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-15', 'improvement']. UNPUBLISHED.

### level-17 (class, Task Size M)

PDF62 row17: first-character1 or adjacent Rogue16->17 ONLY. Exact progression={"feature_keys": ["srd-5.2.1:mechanic/rogue/subclass-feature"], "feature_refs": [], "proficiency_bonus": 6}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-16', 'reflexes']. UNPUBLISHED.

### level-18 (class, Task Size M)

PDF62 row18: first-character1 or adjacent Rogue17->18 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/elusive"], "proficiency_bonus": 6}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-17', 'elusive']. UNPUBLISHED.

### level-19 (class, Task Size M)

PDF62 row19: first-character1 or adjacent Rogue18->19 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/epic-boon"], "proficiency_bonus": 6}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-18', 'boon']. UNPUBLISHED.

### level-20 (class, Task Size M)

PDF62 row20: first-character1 or adjacent Rogue19->20 ONLY. Exact progression={"feature_keys": [], "feature_refs": ["srd-5.2.1:character/feature/rogue/stroke-of-luck"], "proficiency_bonus": 6}. XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'level-19', 'luck']. UNPUBLISHED.

### later-1 (class, Task Size M)

PDF24-25/61: Fighter4->Rogue1 TOTAL5 6500XP/PB3,own reduced traits and all level1 features,HPmax(1,5+CON) not8+CON; reject0XP/PB2,saves/weapons/equipment regrant,missing own eligibility,total21.

Existing prerequisites: []; proposed owners: ['hp', 'multiclass', 'expertise', 'sneak-attack', 'languages', 'mastery']. UNPUBLISHED.


## Source-key matrix

The JSON companion retains every complete source record, all required and option keys, execution owners and hashed engine/web/TUI evidence. All states are unavailable.

| Source record | Execution owners |
|---|---|
| `srd-5.2.1:character/class/rogue` | hp, traits, skills, multiclass |
| `srd-5.2.1:character/choice/rogue/skills` | skills |
| `srd-5.2.1:character/choice/rogue/multiclass-skill` | multiclass |
| `srd-5.2.1:character/choice/rogue/starting-equipment` | equipment |
| `srd-5.2.1:character/level/rogue/1` | level-1 |
| `srd-5.2.1:character/level/rogue/2` | level-2 |
| `srd-5.2.1:character/level/rogue/3` | level-3 |
| `srd-5.2.1:character/level/rogue/4` | level-4 |
| `srd-5.2.1:character/level/rogue/5` | level-5 |
| `srd-5.2.1:character/level/rogue/6` | level-6 |
| `srd-5.2.1:character/level/rogue/7` | level-7 |
| `srd-5.2.1:character/level/rogue/8` | level-8 |
| `srd-5.2.1:character/level/rogue/9` | level-9 |
| `srd-5.2.1:character/level/rogue/10` | level-10 |
| `srd-5.2.1:character/level/rogue/11` | level-11 |
| `srd-5.2.1:character/level/rogue/12` | level-12 |
| `srd-5.2.1:character/level/rogue/13` | level-13 |
| `srd-5.2.1:character/level/rogue/14` | level-14 |
| `srd-5.2.1:character/level/rogue/15` | level-15 |
| `srd-5.2.1:character/level/rogue/16` | level-16 |
| `srd-5.2.1:character/level/rogue/17` | level-17 |
| `srd-5.2.1:character/level/rogue/18` | level-18 |
| `srd-5.2.1:character/level/rogue/19` | level-19 |
| `srd-5.2.1:character/level/rogue/20` | level-20 |
| `srd-5.2.1:character/feature/rogue/expertise` | expertise |
| `srd-5.2.1:character/feature/rogue/sneak-attack` | sneak-attack |
| `srd-5.2.1:character/feature/rogue/thieves-cant` | languages |
| `srd-5.2.1:character/feature/rogue/weapon-mastery` | mastery |
| `srd-5.2.1:character/feature/rogue/cunning-action` | cunning-action |
| `srd-5.2.1:character/feature/rogue/rogue-subclass` | subclass |
| `srd-5.2.1:character/feature/rogue/steady-aim` | steady-aim |
| `srd-5.2.1:character/feature/rogue/ability-score-improvement` | improvement |
| `srd-5.2.1:character/feature/rogue/cunning-strike` | strike, poison, trip, withdraw |
| `srd-5.2.1:character/feature/rogue/uncanny-dodge` | uncanny-dodge |
| `srd-5.2.1:character/feature/rogue/evasion` | evasion |
| `srd-5.2.1:character/feature/rogue/reliable-talent` | reliable-talent |
| `srd-5.2.1:character/feature/rogue/improved-cunning-strike` | improved-strike |
| `srd-5.2.1:character/feature/rogue/devious-strikes` | daze, knock-out, obscure |
| `srd-5.2.1:character/feature/rogue/slippery-mind` | slippery-mind |
| `srd-5.2.1:character/feature/rogue/elusive` | elusive |
| `srd-5.2.1:character/feature/rogue/epic-boon` | boon |
| `srd-5.2.1:character/feature/rogue/stroke-of-luck` | luck |
| `srd-5.2.1:character/choice/rogue/expertise-initial` | expertise |
| `srd-5.2.1:character/prerequisite/rogue/expertise-initial` | expertise |
| `srd-5.2.1:character/choice/rogue/expertise-level-6` | expertise |
| `srd-5.2.1:character/prerequisite/rogue/expertise-level-6` | expertise |
| `srd-5.2.1:character/choice/rogue/sneak-attack-eligibility` | sneak-attack |
| `srd-5.2.1:character/prerequisite/rogue/sneak-attack-eligibility` | sneak-attack |
| `srd-5.2.1:character/choice/rogue/additional-language` | languages |
| `srd-5.2.1:character/choice/rogue/weapon-mastery` | mastery |
| `srd-5.2.1:character/prerequisite/rogue/weapon-mastery-proficiency` | mastery |
| `srd-5.2.1:character/choice/rogue/cunning-action` | cunning-action |
| `srd-5.2.1:character/choice/rogue/rogue-subclass` | subclass |
| `srd-5.2.1:character/prerequisite/rogue/steady-aim-unmoved` | steady-aim |
| `srd-5.2.1:character/choice/rogue/ability-score-improvement` | improvement |
| `srd-5.2.1:character/prerequisite/rogue/ability-score-improvement-qualification` | improvement |
| `srd-5.2.1:character/choice/rogue/epic-boon` | boon |
| `srd-5.2.1:character/prerequisite/rogue/epic-boon-qualification` | boon |
| `srd-5.2.1:character/choice/rogue/cunning-strike` | strike, poison, trip, withdraw |
| `srd-5.2.1:character/choice/rogue/devious-strikes` | daze, knock-out, obscure |
| `srd-5.2.1:character/prerequisite/rogue/devious-strikes-options` | daze, knock-out, obscure |
| `srd-5.2.1:character/choice/rogue/strike-poison` | poison |
| `srd-5.2.1:character/prerequisite/rogue/strike-poison-selected` | poison |
| `srd-5.2.1:character/prerequisite/rogue/strike-poison-numbers` | poison |
| `srd-5.2.1:character/choice/rogue/strike-trip` | trip |
| `srd-5.2.1:character/prerequisite/rogue/strike-trip-selected` | trip |
| `srd-5.2.1:character/choice/rogue/strike-withdraw` | withdraw |
| `srd-5.2.1:character/prerequisite/rogue/strike-withdraw-selected` | withdraw |
| `srd-5.2.1:character/prerequisite/rogue/strike-withdraw-numbers` | withdraw |
| `srd-5.2.1:character/choice/rogue/strike-daze` | daze |
| `srd-5.2.1:character/prerequisite/rogue/strike-daze-selected` | daze |
| `srd-5.2.1:character/choice/rogue/daze-restricted-turn` | daze |
| `srd-5.2.1:character/choice/rogue/strike-knock-out` | knock-out |
| `srd-5.2.1:character/prerequisite/rogue/strike-knock-out-selected` | knock-out |
| `srd-5.2.1:character/prerequisite/rogue/strike-knock-out-numbers` | knock-out |
| `srd-5.2.1:character/choice/rogue/strike-obscure` | obscure |
| `srd-5.2.1:character/prerequisite/rogue/strike-obscure-selected` | obscure |
| `srd-5.2.1:character/prerequisite/rogue/improved-cunning-strike-operands` | improved-strike |
| `srd-5.2.1:character/resource/rogue/stroke-of-luck` | luck |
| `srd-5.2.1:character/subclass/rogue/thief/identity` | subclass |
| `srd-5.2.1:character/prerequisite/rogue/thief/identity-level` | subclass |
| `srd-5.2.1:character/prerequisite/rogue/thief/identity-selection` | subclass |
| `srd-5.2.1:character/feature/rogue/thief/fast-hands` | fast-hands-check, fast-hands-object |
| `srd-5.2.1:character/prerequisite/rogue/thief/fast-hands-level` | fast-hands-check, fast-hands-object |
| `srd-5.2.1:character/choice/rogue/thief/fast-hands` | fast-hands-check, fast-hands-object |
| `srd-5.2.1:character/choice/rogue/thief/sleight-of-hand` | fast-hands-check |
| `srd-5.2.1:character/prerequisite/rogue/thief/sleight-of-hand-check` | fast-hands-check |
| `srd-5.2.1:character/choice/rogue/thief/use-an-object` | fast-hands-object |
| `srd-5.2.1:character/prerequisite/rogue/thief/use-an-object-selection` | fast-hands-object |
| `srd-5.2.1:character/prerequisite/rogue/thief/sleight-of-hand-selection` | fast-hands-check |
| `srd-5.2.1:character/feature/rogue/thief/second-story-work` | climber, jumper |
| `srd-5.2.1:character/prerequisite/rogue/thief/second-story-work-level` | climber, jumper |
| `srd-5.2.1:character/feature/rogue/thief/supreme-sneak` | stealth-attack |
| `srd-5.2.1:character/prerequisite/rogue/thief/supreme-sneak-level` | stealth-attack |
| `srd-5.2.1:character/choice/rogue/thief/stealth-attack` | stealth-attack |
| `srd-5.2.1:character/prerequisite/rogue/thief/stealth-attack-strike` | stealth-attack |
| `srd-5.2.1:character/prerequisite/rogue/thief/stealth-attack-hide` | stealth-attack |
| `srd-5.2.1:character/choice/rogue/thief/stealth-cover` | stealth-attack |
| `srd-5.2.1:character/prerequisite/rogue/thief/stealth-cover-timing` | stealth-attack |
| `srd-5.2.1:character/feature/rogue/thief/use-magic-device` | attunement, charges, scroll-check, scroll-light, scroll-detect |
| `srd-5.2.1:character/prerequisite/rogue/thief/use-magic-device-level` | attunement, charges, scroll-check, scroll-light, scroll-detect |
| `srd-5.2.1:character/choice/rogue/thief/scrolls` | scroll-check |
| `srd-5.2.1:character/prerequisite/rogue/thief/scrolls-level` | scroll-check |
| `srd-5.2.1:character/choice/rogue/thief/scroll-check` | scroll-check |
| `srd-5.2.1:character/prerequisite/rogue/thief/scroll-check-condition` | scroll-check |
| `srd-5.2.1:character/feature/rogue/thief/thiefs-reflexes` | reflexes |
| `srd-5.2.1:character/prerequisite/rogue/thief/thiefs-reflexes-level` | reflexes |

## Publication boundary

Proposed qualification gate belongs under329, with class/level/shared cohorts below it and mandatory terminal paths to337 after329 closes. No additional child under283. No direct429 edge; own HP consumes968. Existing Monk/Paladin bodies stay exact. Reserve ONE shared qualification gateway under192 before any shared write. Dependency and containment DAGs checked separately. Final markers must append to ROOT-released accepted POST-RANGER prefixes; none is accepted here.

UMD source permits any Spell Scroll, but qualification is effect-specific. Recommendation: Light and Detect Magic composition only, all other variants explicitly unavailable. This partial slice cannot complete all-scroll qualification. Parent must choose finite remaining-effect scope. Four exact attunable identities bound protocol tests only; their magical effects do not become available from attunement. Ordinary trap/pocket trusted facts cannot be inferred. Broader item/effect coverage remains unresolved; do not publish generic all-item terminals.
