#!/usr/bin/env python3
"""Check issue #196 identities and edition-specific anchors against pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
NAMES = ["Acid Arrow", "Acid Splash", "Aid", "Alarm", "Alter Self", "Animal Friendship",
         "Animal Messenger", "Animal Shapes", "Animate Dead", "Animate Objects",
         "Antilife Shell", "Antimagic Field"]


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_196.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (107, 108, 109):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_196.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 196]
    records = catalogue["records"]
    assert len(records) == 12
    assert [(r["id"], r["name"], r["source_page"]) for r in records] == expected
    assert [r["name"] for r in records] == NAMES
    assert catalogue["source"]["sha256"] == PIN
    assert catalogue["source"]["license"] == "CC-BY-4.0"
    for r in records:
        assert r["name"] in pages[r["source_page"]], r["name"]
        assert r["scope"] == "complete-source-record" and r["source_id"] == r["id"]
    # The batch ends where the next family starts; nothing after it is imported.
    assert "Antipathy/Sympathy" in pages[109]
    assert not any(r["name"].startswith("Antipathy") for r in records)
    # Each header is compared literally with the record it was transcribed into.
    headers = {
        107: ["Level 2 Evocation (Wizard) Casting Time: Action Range: 90 feet Components: V, S, M (powdered rhubarb leaf) Duration: Instantaneous",
              "Evocation Cantrip (Sorcerer, Wizard) Casting Time: Action Range: 60 feet Components: V, S Duration: Instantaneous",
              "Level 2 Abjuration (Bard, Cleric, Druid, Paladin, Ranger) Casting Time: Action Range: 30 feet Components: V, S, M (a strip of white cloth) Duration: 8 hours",
              "Level 1 Abjuration (Ranger, Wizard) Casting Time: 1 minute or Ritual Range: 30 feet Components: V, S, M (a bell and silver wire) Duration: 8 hours",
              "Level 2 Transmutation (Sorcerer, Wizard) Casting Time: Action Range: Self Components: V, S Duration: Concentration, up to 1 hour",
              "Level 1 Enchantment (Bard, Druid, Ranger) Casting Time: Action Range: 30 feet Components: V, S, M (a morsel of food) Duration: 24 hours"],
        108: ["Level 2 Enchantment (Bard, Druid, Ranger) Casting Time: Action or Ritual Range: 30 feet Components: V, S, M (a morsel of food) Duration: 24 hours",
              "Level 8 Transmutation (Druid) Casting Time: Action Range: 30 feet Components: V, S Duration: 24 hours",
              "Level 3 Necromancy (Cleric, Wizard) Casting Time: 1 minute Range: 10 feet Components: V, S, M (a drop of blood, a piece of flesh, and a pinch of bone dust) Duration: Instantaneous",
              "Level 5 Transmutation (Bard, Sorcerer, Wizard) Casting Time: Action Range: 120 feet Components: V, S Duration: Concentration, up to 1 minute"],
        109: ["Level 5 Abjuration (Druid) Casting Time: Action Range: Self Components: V, S Duration: Concentration, up to 1 hour",
              "Level 8 Abjuration (Cleric, Wizard) Casting Time: Action Range: Self Components: V, S, M (iron filings) Duration: Concentration, up to 1 hour"],
    }
    ordered = [(page, header) for page in (107, 108, 109) for header in headers[page]]
    assert len(ordered) == 12
    seconds = {"Instantaneous": 0, "8 hours": 28800, "24 hours": 86400, "up to 1 hour": 3600, "up to 1 minute": 60}
    for r, (page, header) in zip(records, ordered):
        assert r["source_page"] == page and (r["name"] + " " + header) in pages[page], r["name"]
        kind, rest = header.split(" (", 1)
        classes, rest = rest.split(") Casting Time: ", 1)
        casting, rest = rest.split(" Range: ", 1)
        distance, rest = rest.split(" Components: ", 1)
        components, duration = rest.split(" Duration: ", 1)
        if kind.endswith("Cantrip"):
            level, school = 0, kind.split()[0]
        else:
            _, level, school = kind.split()
        assert (r["level"], r["school"]) == (int(level), school.lower()), r["name"]
        assert r["lists"] == [c.lower() for c in classes.split(", ")], r["name"]
        assert r["ritual"] == casting.endswith(" or Ritual"), r["name"]
        assert r["casting_time"] == {"Action": "action", "1 minute": "1-minute"}[casting.removesuffix(" or Ritual")], r["name"]
        want_range = {"kind": "self", "feet": 0} if distance == "Self" else {"kind": "feet", "feet": int(distance.split()[0])}
        assert r["range"] == want_range, r["name"]
        material = components.split("M (", 1)[1].rstrip(")") if "M (" in components else None
        got = r["components"]
        assert got["verbal"] and got["somatic"] and components.startswith("V, S"), r["name"]
        assert (got["material"] or {}).get("description") == material, r["name"]
        if material:
            assert got["material"]["cost_cp"] == 0 and not got["material"]["consumed"], r["name"]
        concentration = duration.startswith("Concentration, ")
        assert r["concentration"] == concentration, r["name"]
        assert r["duration_seconds"] == seconds[duration.removeprefix("Concentration, ")], r["name"]
        assert r["slot_levels"] == ([] if int(level) == 0 else list(range(int(level), 10))), r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_196_test.go; this is not a prose parser.
    anchors = {
        107: ["4d4 Acid damage and 2d4 Acid dam- age at the end of its next turn", "half as much of the initial damage only",
              "(both initial and later) increases by 1d4 for each spell slot level above 2",
              "5-foot-radius Sphere", "Dexterity saving throw or take 1d6 Acid damage", "levels 5 (2d6), 11 (3d6), and 17 (4d6)",
              "Choose up to three creatures within range", "increase by 5 for the duration",
              "Points increase by 5 for each spell slot level above 2",
              "a door, a window, or an area within range that is no larger than a 20-foot Cube",
              "touches or enters the warded area", "designate creatures that won’t set off the alarm",
              "Audible Alarm.", "handbell for 10 seconds within 60 feet", "Mental Alarm.", "within 1 mile of the warded area",
              "awakens you if you’re asleep", "Choose one of the following options", "take a Magic action to replace",
              "Aquatic Adaptation.", "gain a Swim Speed equal to your Speed", "Change Appearance.",
              "none of your statistics change", "can’t appear as a creature of a different size", "Natural Weapons.",
              "claws (Slashing), fangs (Piercing), horns (Piercing), or hooves (Blud- geoning)", "it deals 1d6 damage",
              "spellcasting ability modifier for the attack and damage rolls rather than using Strength",
              "Target a Beast that you can see within range", "Wisdom saving throw or have the Charmed condition"],
        108: ["or one of your allies deals damage to the target, the spells ends",
              "one additional Beast for each spell slot level above 1",
              "A Tiny Beast of your choice", "Charisma saving throw", "Challenge Rating isn’t 0, it automatically succeeds",
              "a location you have visited", "up to twenty-five words", "about 25 miles per 24 hours or 50 miles if the Beast can fly",
              "the message is lost, and the Beast returns to where you cast the spell",
              "tion increases by 48 hours for each spell slot level above 2",
              "any number of willing creatures", "Large or smaller Beast", "Challenge Rating of 4 or lower",
              "dif- ferent form for each target", "take a Magic action to transform the targets again",
              "creature type; Hit Points; Hit Point Dice; alignment; ability to communicate; and Intelligence, Wisdom, and Charisma scores",
              "it can’t cast spells", "equipment melds into the new form",
              "Temporary Hit Points equal to the Hit Points of the first form", "ends it as a Bonus Action",
              "a pile of bones or a corpse of a Medium or Small Humanoid", "a Skeleton if you chose bones or a Zombie if you chose a corpse",
              "within 60 feet of you", "takes the Dodge action and moves only to avoid harm",
              "under your control for 24 hours", "up to four creatures", "two additional Undead crea- tures for each spell slot level above 3",
              "different corpse or pile of bones", "aren’t being worn or carried, aren’t fixed to a surface, and aren’t Gargantuan",
              "equal to your spellcasting ability modifier", "a Medium or smaller target counts as"],
        109: ["one object, a Large target counts as two, and a Huge target counts as three",
              "sprouts legs, and becomes a Construct that uses the Animated Object stat block",
              "until it is reduced to 0 Hit Points", "ally to you and your allies",
              "shares your Initiative count and takes its turn immediately after yours", "within 500 feet of you",
              "reverts to its object form, and any remain- ing damage carries over",
              "1d4 (Medium or smaller), 1d6 (Large), or 1d12 (Huge) for each spell slot level above 5",
              "Huge or Smaller Construct, Unaligned", "AC 15", "HP 10 (Medium or smaller), 20 (Large), 40 (Huge)", "Speed 30 ft.",
              "Str 16 +3 +3 Dex 10 +0 +0 Con 10 +0 +0", "Int 3 −4 −4 Wis 3 −4 −4 Cha 1 −5 −5",
              "Immunities Poison, Psychic; Charmed, Exhaustion, Frightened, Paralyzed, Poisoned",
              "Blindsight 30 ft.; Passive Perception 6", "Understands the languages you know",
              "CR None (XP 0; PB equals your Proficiency Bonus)",
              "Melee Attack Roll: Bonus equals your spell attack modifier, reach 5 ft.",
              "Force damage equal to 1d4 + 3 (Medium or smaller), 2d6 + 3 + your spellcasting ability modifier (Large), or 2d12 + 3 + your spellcasting ability modifier (Huge)",
              "10-foot Emanation for the duration", "other than Constructs and Undead from passing or reach- ing through it",
              "An affected creature can cast spells", "or make attacks with Ranged or Reach weapons through the barrier",
              "forced to pass through the barrier, the spell ends", "10-foot Em- anation",
              "cast spells, take Magic actions, or create other magical effects inside the aura",
              "Magical properties of magic items don’t work inside the aura or on anything inside it",
              "can’t extend into the aura, and no one can teleport into or out of it or use planar travel there",
              "Portals close temporarily", "except those cast by an Artifact or a deity, are suppressed",
              "counts against its duration", "Dispel Magic has no effect on the aura", "don’t nullify each other"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            assert snippet in pages[page], (page, snippet)
    by_name = {r["name"]: r for r in records}

    def details(name, key):
        return by_name[name]["source_details"][key]

    arrow = by_name["Acid Arrow"]
    assert [(e["dice_count"], e["die_sides"], e["damage_type"]) for e in arrow["effects"]] == [(4, 4, "acid"), (2, 4, "acid")]
    assert arrow["scaling"] == {"kind": "slot", "above_level": 2, "dice_count": 1, "die_sides": 4, "effect_ids": ["initial", "later"]}
    assert details("Acid Splash", "splash")["upgrades"] == [
        {"character_level": 5, "dice_count": 2, "die_sides": 6}, {"character_level": 11, "dice_count": 3, "die_sides": 6},
        {"character_level": 17, "dice_count": 4, "die_sides": 6}]
    assert details("Aid", "aid")["additional_hit_points_per_slot"] == 5
    assert [m["choice"] for m in details("Alarm", "alarm")["modes"]] == ["audible", "mental"]
    assert by_name["Alarm"]["source_details"]["casting_seconds"] == 60
    assert by_name["Animate Dead"]["source_details"]["casting_seconds"] == 60
    assert [g["growth"] for g in details("Alter Self", "alter")["weapons"]["growths"]] == ["claws", "fangs", "horns", "hooves"]
    assert details("Animal Friendship", "friendship")["source_pages"] == [107, 108]
    assert details("Animal Messenger", "messenger")["additional_duration_seconds_per_slot"] == 48 * 3600
    assert details("Animal Shapes", "shapes")["form_max_cr"] == 4
    assert details("Animate Dead", "dead")["additional_creatures_per_slot"] == 2
    objects = details("Animate Objects", "objects")
    assert objects["source_pages"] == [108, 109]
    block = objects["animated_object"]
    assert [(a["score"], a["modifier"], a["save"]) for a in block["abilities"]] == [
        (16, 3, 3), (10, 0, 0), (10, 0, 0), (3, -4, -4), (3, -4, -4), (1, -5, -5)]
    assert [h["count"] for h in block["hit_points"]] == [10, 20, 40]
    assert details("Antilife Shell", "antilife")["exempt_types"] == ["construct", "undead"]
    assert details("Antimagic Field", "antimagic")["emanation_feet"] == 10
    print("#196 pinned PDF, twelve identities/pages, headers and edition-specific anchors pass")


if __name__ == "__main__":
    main()
