#!/usr/bin/env python3
"""Build a source audit and unreserved proposal; performs no tracker writes."""
import hashlib
import json
from pathlib import Path

ROOT=Path(__file__).resolve().parents[1]
OUT=ROOT/"docs/investigations"
PIN="8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
BASE="952fab917329e00cceebf46a8bacd59cb1c30e99"

def build():
    cats={name:json.loads((ROOT/f"internal/rules/srd521/testdata/character-{name}.json").read_text()) for name in ["rogue","thief","advancement","multiclass"]}
    tasks=[]
    def add(key,fixture,prerequisites=(),depends=(),category="class"):
        tasks.append(dict(key=key,title="Rogue: "+key.replace("-"," "),fixture=fixture,prerequisites=list(prerequisites),depends=list(depends),category=category,project_size="M",issue=None,publication="UNPUBLISHED"))
    add("hp","PDF23/61: own D8/fixed5 adapter ONLY, qualification968 -> transaction429. First character8+CON; later level max(1,d8+CON) or max(1,5+CON). +2 CON fixed gives7; die1 CON-3 gives1; CON17->18 at total8 adds8 retroactively without healing. Reject D10/fixed6, reroll, first8 on later entry.",[968,58])
    add("traits","PDF61: DEX primary; DEX/INT saves; Simple and Martial with Finesse OR Light; Thieves’ Tools and Light armor. Reject Heavy armor, martial neither property, later-class saves.",[430,333,58])
    add("skills","PDF61: choose exactly4 distinct Acrobatics/Athletics/Deception/Insight/Intimidation/Investigation/Perception/Persuasion/Sleight of Hand/Stealth. Positive first4; negative3/5/duplicate/Survival.",[384,430])
    add("equipment","PDF61: A Leather Armor,2 Daggers,Shortsword,Shortbow,20 Arrows,Quiver,Thieves’ Tools,Burglar’s Pack,8GP OR B100GP. Positive exact A or B; reject19arrows,9GP,mixed/repeated grants. Inventory grant only; no inferred item effects.",[195,430])
    add("multiclass","PDF24-25/61: own DEX13 and EVERY admitted current-class primary minimum13. DEX12 fails independently; existing Fighter STR12/DEX12 fails. Gain D8,one Rogue skill,Thieves’ Tools,Light armor and level1 features; no saves,weapon/equipment/four-skill regrant. Positive Fighter4 DEX13 -> Rogue1; reject unqualified vector.",[335,430], ["hp"])
    add("expertise","PDF61: exactly2 EXISTING skill proficiencies at1 and2 MORE at6; PB2 expertise adds4 not6, PB3 adds6. Reject tools,unproficient/duplicate skills,third selection,grant6 at5.",[384,430,61])
    add("sneak-attack","PDF61-62: once PER TURN creature hit,Finesse OR Ranged weapon; Advantage OR ally within5ft/not Incapacitated AND no Disadvantage. Same weapon damage type; ceil(Rogue level/2)d6. Positive Rogue5 3d6,own turn and different off-turn; reject same turn repeat,ally6ft/incapacitated,disadvantage,unarmed/nonqualifying weapon,noncreature. Inject3/4/5=12 extra; no off-hand conflation.",[487])
    add("languages","PDF62/18: Thieves’ Cant plus1 from actual standard/rare language tables. Positive Common Sign Language; reject unknown/duplicate/two extra; preserve independent provenance.",[384,430])
    add("mastery","PDF62: two proficient weapon kinds; completed Long Rest may replace. Dagger/Shortbow positive; reject3,duplicate,unproficient,rest start/interruption. Consume selected property owner489-496,not label.",[384,428,430]+list(range(489,497)))
    add("cunning-action","PDF62: on OWN turn Bonus Action Dash/Disengage/Hide; one choice. Reject off-turn,spent Bonus Action,automatic extra action. Speed30 Dash adds30 movement; consumes780/781 and proposed hide output.",[780,781], ["hide"])
    add("subclass","PDF62/64: choose exactly Thief at Rogue3; no subclass at2,second choice or foreign subclass; feature grant class3/9/13/17 only.",[384])
    add("steady-aim","PDF62: Bonus Action only if NO movement this turn; Advantage NEXT attack on CURRENT turn; Speed0 remainder turn. Positive no move -> attack; reject moved1ft,second attack benefit,next turn carryover.")
    featdeps=[384,389]+list(range(385,402))+[487,427,428,778,779,962]
    add("improvement","PDF63/86-88: class4/8/10/12/16 ASI OR qualifying feat. +2one/+1two,max20. Positive18->20; reject21,duplicate,unqualified,wrong class level. NO blanket Spellcasting,encounter or Concentration prerequisite for ordinary selection.399 Spell Recall needs ACTUAL Spellcasting;386 Magic Initiate and scroll access do not supply it. Selected feat requires actual relevant487/427/428/778/779/962 execution,not269/272/270/268 planning closure.",featdeps)
    add("strike","PDF63: Rogue5 when Sneak Attack damage occurs choose1 effect; remove cost dice BEFORE rolling,resolve AFTER attack damage; DC8+DEX+livePB. DEX3/PB3 DC14,3d6 minus1=2d6. Reject miss,no Sneak Attack,negative dice,unsupported effect.",[],["sneak-attack"])
    add("poison","PDF63: cost1d6,Poisoner’s Kit ON PERSON; CON save failure Poisoned1minute; repeat end each target turn success ends own instance. DC14 fail13/succeed14; reject kit absent,client cure,clearing other poison,expiry61seconds.",[504,727], ["strike"])
    add("trip","PDF63: cost1d6 target Large OR SMALLER,DEX save fail Prone. DC14 roll13 vs14; Huge rejects,no cost on invalid target. Consume488 source-owned Prone.",[488],["strike"])
    add("withdraw","PDF63: cost1d6 immediately AFTER attack move up to HALF Speed without Opportunity Attacks; Speed30 allows15,not16; no teleport or next-turn movement credit.",[],["strike"])
    add("uncanny-dodge","PDF63: Reaction when attacker YOU SEE hits with attack roll; halve attack damage round DOWN. 9->4; reject unseen,save-only damage,no reaction,miss,whole round resistance.")
    add("evasion","PDF63: DEX save for HALF damage only; success0/fail HALF; cannot use Incapacitated. Damage11 ->0 or5. Reject CON save,no-half effect,incapacitated,ASI ownership.",[487])
    add("reliable-talent","PDF63: ability check using skill OR tool proficiency treats d20<=9 as10. Die9->10 plus own modifiers;10 stays10; reject unproficient,save,attack,change die11.",[430,61])
    add("improved-strike","PDF63: Rogue11 up to2 strike effects,pay EACH before rolling;6d6 minus1minus1=4d6. Reject3 effects,grant at10,free second effect.",[],["strike","poison","trip","withdraw"])
    add("daze","PDF63: Rogue14 cost2d6 CON save; failure NEXT turn only ONE of move/action/Bonus Action.7d6->5d6; reject two categories,next-next turn,wrong save. No invented Stunned condition.",[],["strike"])
    add("knock-out","PDF63: Rogue14 cost6d6 CON fail Unconscious1minute OR ANY damage; repeat save end each turn success ends own instance.7d6->1d6; reject five dice affordability,damage not ending,foreign instance removal.",[487],["strike"])
    add("obscure","PDF63: Rogue14 cost3d6 DEX fail Blinded until END next target turn.7d6->4d6; reject CON save,early start-turn end,foreign condition removal.",[960,963],["strike"])
    add("slippery-mind","PDF63: Rogue15 adds WIS and CHA save proficiency using livePB. PB5 bonus5 not10; reject level14,duplicate stacking,INT replacement.",[430])
    add("elusive","PDF63: Rogue18 no attack roll has Advantage AGAINST self unless Incapacitated. Positive unseen attacker loses advantage; negative Incapacitated allows normal calculation,save treated as attack.",[487])
    add("boon","PDF63/87-88: Rogue19 Epic Boon OR other qualified feat,Night Spirit recommendation only. Positive chosen qualified boon; reject auto-Night-Spirit,level18,unqualified feat;399 requires real Spellcasting.",featdeps)
    add("luck","PDF63: Rogue20 FAILED D20 Test -> roll20; once until completed Short OR Long Rest. Positive failed die2 ->20 (then ordinary modifiers); reject success,unknown outcome,second use,rest start/interrupted rest,automatic natural20 critical assumption.",[427,428])
    add("fast-hands-check","PDF64: Thief3 Bonus Action DEX(Sleight of Hand) lock/trap with Thieves’ Tools OR pocket. Check uses supplied trusted object DC15,14fails/15succeeds; reject no tools for lock/trap,unknown object facts,foreign pocket,spell substitution. Requires separately proposed bounded lock/trap/pocket outputs.",[],["subclass","object-check"])
    add("fast-hands-object","PDF64: Thief3 Bonus Action Utilize OR Magic to USE ITEM requiring that action; not arbitrary spell Magic action. Positive named qualified item route only; reject wrong activation,unsupported item,spent Bonus Action; exact mundane/magic outputs prerequisite.",[],["subclass","utilize","magic-item"])
    add("climber","PDF64: Thief3 Climb Speed=current Speed; Speed30 ->30,Speed0->0. Reject flat30,Swim grant; dynamic modifiers retained.",[],["subclass"])
    add("jumper","PDF64/183/185: Thief3 optional DEX rather than STR jump distance. STR10/DEX16 running long jump16 vs10; high jump3+DEXmod=6 vs3; reject adding both,flight,waiving distance/movement/run-up rules.",[],["subclass"])
    add("stealth-attack","PDF64: Thief9 Cunning Strike1d6; Hide’s Invisible survives THIS attack only if end turn behind Three-Quarters OR Total Cover.5d6->4d6; reject Half Cover,wrong source Invisible,no Hide,no end-turn check; save original Hide source and independently end on other triggers.",[962], ["strike","hide","subclass"])
    add("attunement","PDF64/102: Thief13 up to FOUR simultaneous,not five. Positive fourth eligible item; reject fifth,foreign/duplicate item,ignoring item-specific restrictions. No all-item activation.",[],["subclass","attune-item"])
    add("charges","PDF64: Thief13 on item PROPERTY expending charges roll1d6;6 skips charge spend,1-5 ordinary spend; preserve action cost and property effect. Positive6 keeps2charges,5 spends1of2; reject reroll/recharge/waived eligibility; finite proposed Wand fixture only.",[],["subclass","magic-item"])
    add("scroll-check","PDF64: Thief13 any Spell Scroll,INT casting ability. Cantrip/level1 reliable; higher INT(Arcana) DC10+SPELL LEVEL; level2 DC12 (11fail disintegrates,12casts),level9 DC19. No Spellcasting grant; no invented class-list restriction. Real spell effect must separately qualify; broad339spells remain unavailable.",[],["subclass","scroll-item"])
    add("scroll-light","PDF64/144/104-106: compose Thief scroll adapter with actual967 Light,INT; no learned/prepared/Spellcasting grant. Reject other-holder object,Huge,foreign scroll,slot spending or ritual/upcast. Own ongoing Light dismissal NO ACTION when not Incapacitated; Incapacitated alone does not end nonconcentration Light.",[967], ["scroll-check"])
    add("scroll-detect","PDF64/123/104-106: compose scroll adapter with906 Detect Magic,INT and spell concentration/visibility/barriers only. Reject invisible-object aura disclosure; no slot spend or Ritual from consuming scroll. Prepared Ritual elsewhere normal time+10min,no slot,no Ritual Adept,NOritual upcast; >=1min cast Magic each turn/casting Concentration; interrupted no slot spend.",[906],["scroll-check"])
    add("reflexes","PDF64: Thief17 TWO turns FIRST round every combat,normal Initiative and minus10. Initiative15 ->15and5; reject later round,reload granting third,legacy surprise gate; equal initiative ordering deterministic; each is actual distinct turn with action economy and per-turn features.",[],["subclass"])
    # Proposed shared fixtures are deliberately finite and require parent reservation.
    add("hide","PDF183: DC15 DEX(Stealth) while Heavily Obscured or Three-Quarters/Total Cover,out of enemy sight; save check total as discovery DC,Invisible source. Ends on louder-than-whisper sound,enemy finds,attack roll,Verbal spell.14fails/15succeeds; Half Cover alone rejects. Supreme Sneak owns its one attack exception.",[962],category="shared-action")
    add("object-check","PDF98/94/64: ONLY owned Lock/Thieves’ Tools lock operation,one trusted mechanical trap disarm,and one trusted pocket fixture; retain DEX(Sleight of Hand) Fast Hands exception. DC15 14/15 controls; no arbitrary traps/catalogue inference; unknown DC/ownership unavailable.",[195,430],category="shared-item")
    add("utilize","PDF99/191: ONLY owned Signal Whistle via Utilize, audible up to600ft. Positive trusted listener600ft; negative601ft,unknown acoustic facts,foreign whistle,spell action,invented consumption. Tinderbox exposed-fuel ignition is already Bonus Action onPDF100 and is NOT this Utilize fixture. Fast Hands supplies Bonus Action exception,not this owner.",[195],category="shared-item")
    add("magic-item","PDF magic item Wand of Magic Detection: ONLY its Detect Magic property,3charges,Magic action,1charge; actual906 effect,source-specific dawn recharge separate and UNAVAILABLE here. Reject0charges,foreign/unattuned as source requires,other wand property. Thief conserves charge only,never restores.",[195,906],category="shared-item")
    add("attune-item","PDF102: ONLY four distinct owned source identities: Cloak of Protection PDF216,Gauntlets of Ogre Power PDF223,Ring of Jumping PDF238,Ring of Protection PDF238. Attunement protocol ONLY,not their magical effects; Short Rest focus one item,normal cap3,Thief adapter4. Reject unknown restrictions,5items,foreign possession,rest interruption; reject duplicate Ring of Protection type and using the identification Short Rest as the attunement rest; item magical effects remain separately unavailable without actual owners.",[195,427],category="shared-item")
    add("scroll-item","PDF Spell Scroll: ONLY Light/Detect Magic canonical scroll identities; readable ownership and consume/disintegration atomic with effect/check receipt; reject foreign/unknown/duplicate spend. Higher-level check uses isolated supplied scroll level,not granted unsupported effect. All other339effect coverage excluded pending parent decision.",[195],category="shared-item")
    grants={1:["expertise","sneak-attack","languages","mastery"],2:["cunning-action"],3:["subclass","steady-aim","fast-hands-check","fast-hands-object","climber","jumper"],4:["improvement"],5:["strike","poison","trip","withdraw","uncanny-dodge"],6:["expertise"],7:["evasion","reliable-talent"],8:["improvement"],9:["stealth-attack"],10:["improvement"],11:["improved-strike"],12:["improvement"],13:["attunement","charges","scroll-check","scroll-light","scroll-detect"],14:["daze","knock-out","obscure"],15:["slippery-mind"],16:["improvement"],17:["reflexes"],18:["elusive"],19:["boon"],20:["luck"]}
    records=cats["rogue"]["records"]+cats["thief"]["records"]; by={r["id"]:r for r in records}
    for n in range(1,21):
        row=next(r for r in records if r["kind"]=="level" and r["level"]==n)
        add(f"level-{n}",f"PDF62 row{n}: first-character1 or adjacent Rogue{n-1}->{n} ONLY. Exact progression="+json.dumps(row["progression"],sort_keys=True)+". XP/livePB from complete canonical advancement by post-transition TOTAL admitted class levels,NOT class-table PB. Fighter4/Rogue1->Rogue2 total6=14000XP/PB3; Fighter4->Rogue1 total5=6500/PB3. Reject class-indexed300XP/PB2,noninteger/invalid levels,total21,skip,unqualified combination,duplicate grant; class dice/resources remain Rogue-level. Later-class variant MUST consume own multiclass eligibility; first1 consumes traits/skills/equipment. Independent feature ownership,not ASI fallback.",[],["hp","multiclass"]+(["traits","skills","equipment"] if n==1 else [f"level-{n-1}"])+grants[n])
    add("later-1","PDF24-25/61: Fighter4->Rogue1 TOTAL5 6500XP/PB3,own reduced traits and all level1 features,HPmax(1,5+CON) not8+CON; reject0XP/PB2,saves/weapons/equipment regrant,missing own eligibility,total21.",[],["hp","multiclass"]+grants[1])
    return cats,records,by,tasks,grants

def audit():
    cats,records,by,tasks,grants=build()
    features={"expertise":["expertise"],"sneak-attack":["sneak-attack"],"thieves-cant":["languages"],"weapon-mastery":["mastery"],"cunning-action":["cunning-action"],"rogue-subclass":["subclass"],"steady-aim":["steady-aim"],"ability-score-improvement":["improvement"],"cunning-strike":["strike","poison","trip","withdraw"],"uncanny-dodge":["uncanny-dodge"],"evasion":["evasion"],"reliable-talent":["reliable-talent"],"improved-cunning-strike":["improved-strike"],"devious-strikes":["daze","knock-out","obscure"],"slippery-mind":["slippery-mind"],"elusive":["elusive"],"epic-boon":["boon"],"stroke-of-luck":["luck"],"fast-hands":["fast-hands-check","fast-hands-object"],"second-story-work":["climber","jumper"],"supreme-sneak":["stealth-attack"],"use-magic-device":["attunement","charges","scroll-check","scroll-light","scroll-detect"],"thiefs-reflexes":["reflexes"]}
    def owners(r):
        ident=r["id"];slug=ident.rsplit("/",1)[-1]
        if r["kind"]=="level":return ["level-"+str(r["level"])]
        if r["kind"]=="class":return ["hp","traits","skills","multiclass"]
        if r["kind"]=="subclass":return ["subclass"]
        if slug in features:return features[slug]
        for text,target in [("starting-equipment","equipment"),("multiclass-skill","multiclass"),("skills","skills"),("strike-poison","poison"),("strike-trip","trip"),("strike-withdraw","withdraw"),("strike-daze","daze"),("daze-restricted","daze"),("strike-knock-out","knock-out"),("strike-obscure","obscure"),("scroll","scroll-check"),("stealth-","stealth-attack"),("sleight-of-hand","fast-hands-check"),("use-an-object","fast-hands-object")]:
            if text in slug:return [target]
        return owners(by[r["parent_id"]])
    evidence={
      "internal/rules/srd521/creation.go":"NewSupportedBuild fixed Fighter/Cleric admission only; no Rogue.",
      "internal/rules/srd521/character.go":"Project class/level guard excludes Rogue; class label is insufficient.",
      "internal/rules/srd521/catalogue.go":"ValidateReadiness requires implemented mechanic keys; metadata does not grant runtime.",
      "internal/game/advancement_preview.go":"Fixed Fighter1->2 preview; no Rogue HP/eligibility adapter.",
      "internal/game/spell_preview.go":"Fixed Cleric cast seams; no Thief scroll admission.",
      "internal/engine/action_confirmation.go":"Durable receipt/transaction seam; no inferred Rogue feature execution.",
      "internal/combat/resolver.go":"Shared combat arithmetic seam; no source-key Rogue entitlements.",
      "internal/handlers/character_rules.go":"Shared sheet read route; projection guard still controls readiness.",
      "frontend/src/components/start/CharacterGuidedForm.tsx":"Class selector labels do not qualify Rogue creation.",
      "frontend/src/components/character/ActionPreviewPanel.tsx":"Shared confirmation UI does not qualify Rogue action options.",
      "tui/app.go":"Trusted command path consumes shared engine; no separate Rogue entitlement."}
    evidence={p:dict(finding=f,sha256=hashlib.sha256((ROOT/p).read_bytes()).hexdigest()) for p,f in evidence.items()}
    matrix=[]
    for r in records:
        mapping={}
        for key in r["required_mechanics"]+r["option_keys"]:
            own=owners(r)
            if r["kind"]=="class":own=["hp" if "hit-die" in key else "multiclass" if "/multiclass/" in key else "traits"]
            if r["id"].endswith("/use-magic-device"):
                own=["attunement"] if "attun" in key else ["charges"] if "charge" in key else ["scroll-check"]
            if r["id"].endswith("/second-story-work"):own=["climber"] if "climb" in key else ["jumper"]
            mapping[key]=own
        matrix.append(dict(record=r,execution_owners=owners(r),key_owners=mapping,state="metadata only; Rogue/Thief unavailable",evidence={"engine":list(evidence)[:8],"web":list(evidence)[8:10],"tui":["tui/app.go"]}))
    for t in tasks:
        t["source_records"]=[r["record"]["id"] for r in matrix if t["key"] in r["execution_owners"]]
        t["source_keys"]=sorted({k for r in matrix for k,v in r["key_owners"].items() if t["key"] in v})
    common="Preview inert: no writes or dice. Explicit confirmation revalidates source,actor,campaign,revision,target,range,timing,cost and qualification. Injectable dice; atomic choice/effect/resource/history persistence. Retry/replay/reload/reconnect returns same receipt,conflicting payload rejects,concurrent confirmations commit ONCE. Inject failures for complete isolated PostgreSQL rollback. Web AND TUI expose supported choices/unavailable reasons; accessible cancel/confirm/reconnect; changed web requires desktop/phone fixtures,build and focused lint. Automatic effects belong to triggering confirmed event,not optional new activation. Unknown/private facts unavailable without leaking hidden reasons. Short tests are not integration-tag PostgreSQL/browser/device proof. No production/inference/runtime promotion. Class-total levels integer1-20,not booleans/fractions; each admitted class positive,post-transition total<=20; later-class variants require own multiclass eligibility,first-character variants do not acquire a multiclass prerequisite at runtime. Spell Recall needs actual Spellcasting; ordinary feats have no blanket casting/encounter requirement. Slot-cast one/turn applies only to slot expenditure TO CAST; authoritative confirmed invalid target spends slot/no effect with source-private filtering; unauthorized/unverifiable facts unavailable."
    manifest=dict(schema_version=1,decision=329,base=BASE,source_sha256=PIN,publication="PENDING parent reservation, complete item/effect scope and POST-RANGER release; NOT publication or sealed proof",catalogues=cats,records=matrix,issues=tasks,grants=grants,evidence=evidence,common_acceptance=common,accepted_post_ranger=None,supplements={str(n):"PENDING ROOT serial publication; no accepted prefix captured" for n in [47,54,283,337]})
    (OUT/"issue-329-rogue-scope.json").write_text(json.dumps(manifest,indent=2)+"\n")
    lines=["# Rogue and Thief decision checkpoint (#329)","",manifest["publication"]+".","", "This decision proposes bounded execution; it implements no Rogue runtime. Catalogue completion is source evidence only. No new native issues or dependencies have been published. #329 remains open. The item/scroll proposal requires the parent’s scope reservation; it is not yet the complete publication acceptance proof.","", "Source: English SRD 5.2.1, Rogue PDF61–63, Thief64, advancement23, multiclass24–26; SHA256 `"+PIN+"`. Source-derived records retain the full CC-BY-4.0 attribution and adaptation notice in the canonical catalogues and docs/srd-5.2.1-conformance.md. This plan and ownership matrix are Edda additions.","", "## All twenty rows", "", "Class PB is preserved source metadata. Live PB and XP use the sum of authoritative admitted post-transition class levels. Fighter4/Rogue1→Rogue2 is total6,14000XP/PB3; Fighter4→Rogue1 is total5,6500XP/PB3. HP/features/Sneak Attack use Rogue class level.","", "| Rogue level | Class-table PB | Sneak Attack | Required new behavior owners |", "|---|---|---|---|"]
    for n in range(1,21):
        r=next(r for r in records if r["kind"]=="level" and r["level"]==n)
        lines.append('| {} | {} | {}d6 | {} |'.format(n,r['progression']['proficiency_bonus'],(n+1)//2,', '.join(grants[n])))
    lines += ['', '## Runtime evidence', '']
    lines += [f"- `{p}`: {e['finding']}" for p,e in evidence.items()]
    lines += ['', '## Proposed one-PR terminals', '',common,'']
    lines += [f"### {t['key']} ({t['category']}, Task Size M)\n\n{t['fixture']}\n\nExisting prerequisites: {t['prerequisites']}; proposed owners: {t['depends']}. UNPUBLISHED.\n" for t in tasks]
    lines += ['', '## Source-key matrix', '', 'The JSON companion retains every complete source record, all required and option keys, execution owners and hashed engine/web/TUI evidence. All states are unavailable.', '', '| Source record | Execution owners |', '|---|---|']
    lines += ['| `'+r['record']['id']+'` | '+', '.join(r['execution_owners'])+' |' for r in matrix]
    lines += ['', '## Publication boundary', '', 'Proposed qualification gate belongs under329, with class/level/shared cohorts below it and mandatory terminal paths to337 after329 closes. No additional child under283. No direct429 edge; own HP consumes968. Existing Monk/Paladin bodies stay exact. Reserve ONE shared qualification gateway under192 before any shared write. Dependency and containment DAGs checked separately. Final markers must append to ROOT-released accepted POST-RANGER prefixes; none is accepted here.', '', 'UMD source permits any Spell Scroll, but qualification is effect-specific. Recommendation: Light and Detect Magic composition only, all other variants explicitly unavailable. This partial slice cannot complete all-scroll qualification. Parent must choose finite remaining-effect scope. Four exact attunable identities bound protocol tests only; their magical effects do not become available from attunement. Ordinary trap/pocket trusted facts cannot be inferred. Broader item/effect coverage remains unresolved; do not publish generic all-item terminals.']
    (OUT/'issue-329-rogue-decision.md').write_text('\n'.join(lines)+'\n')
    for n in [47,54,283,337]:
        content='Rogue #329: publication pending.'
        if n!=47:content='Rogue #329 source audit only; bounded execution roster and native paths PENDING parent reservation/publication. No runtime promotion. Final seal waits for accepted POST-RANGER release.'
        section=f'\n<!-- edda-329:start -->\n{content}\n<!-- edda-329:end -->\n'
        assert len(section)<=(90 if n==47 else 750)
        (OUT/f'issue-329-parent-{n}-PENDING.md').write_text(section)
    print(len(records),'records',len(tasks),'proposed terminals; zero tracker writes')
    return manifest

if __name__=='__main__':audit()
