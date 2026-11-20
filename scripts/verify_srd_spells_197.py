#!/usr/bin/env python3
"""Check issue #197 identities and edition-specific anchors against pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
NAMES = ["Antipathy/Sympathy", "Arcane Eye", "Arcane Hand", "Arcane Lock", "Arcane Sword",
         "Arcanist’s Magic Aura", "Astral Projection", "Augury", "Aura of Life", "Awaken",
         "Bane", "Banishment"]
# Source header lines, compared verbatim against the page that starts each spell
# (or, for Arcanist's Magic Aura, the page its header continues onto).
HEADERS = {
    "Antipathy/Sympathy": [(109, "Level 8 Enchantment (Bard, Druid, Wizard) Casting Time: 1 hour Range: 60 feet "
                                 "Components: V, S, M (a mix of vinegar and honey) Duration: 10 days")],
    "Arcane Eye": [(110, "Level 4 Divination (Wizard) Casting Time: Action Range: 30 feet "
                         "Components: V, S, M (a bit of bat fur) Duration: Concentration, up to 1 hour")],
    "Arcane Hand": [(110, "Level 5 Evocation (Sorcerer, Wizard) Casting Time: Action Range: 120 feet "
                          "Components: V, S, M (an eggshell and a glove) Duration: Concentration, up to 1 minute")],
    "Arcane Lock": [(110, "Level 2 Abjuration (Wizard) Casting Time: Action Range: Touch "
                          "Components: V, S, M (gold dust worth 25+ GP, which the spell consumes) "
                          "Duration: Until dispelled")],
    "Arcane Sword": [(110, "Level 7 Evocation (Bard, Wizard) Casting Time: Action Range: 90 feet "
                           "Components: V, S, M (a miniature sword worth 250+ GP) "
                           "Duration: Concentration, up to 1 minute")],
    "Arcanist’s Magic Aura": [(110, "Level 2 Illusion (Wizard) Casting Time: Action Range: Touch"),
                                   (111, "Components: V, S, M (a small square of silk) Duration: 24 hours")],
    "Astral Projection": [(111, "Level 9 Necromancy (Cleric, Warlock, Wizard) Casting Time: 1 hour Range: 10 feet "
                                "Components: V, S, M (for each of the spell’s targets, one jacinth worth "
                                "1,000+ GP and one silver bar worth 100+ GP, all of which the spell consumes) "
                                "Duration: Until dispelled")],
    "Augury": [(111, "Level 2 Divination (Cleric, Druid, Wizard) Casting Time: 1 minute or Ritual Range: Self "
                     "Components: V, S, M (specially marked sticks, bones, cards, or other divinatory tokens "
                     "worth 25+ GP) Duration: Instantaneous")],
    "Aura of Life": [(111, "Level 4 Abjuration (Cleric, Paladin) Casting Time: Action Range: Self "
                           "Components: V Duration: Concentration, up to 10 minutes")],
    "Awaken": [(111, "Level 5 Transmutation (Bard, Druid) Casting Time: 8 hours Range: Touch "
                     "Components: V, S, M (an agate worth 1,000+ GP, which the spell consumes) "
                     "Duration: Instantaneous")],
    "Bane": [(112, "Level 1 Enchantment (Bard, Cleric, Warlock) Casting Time: Action Range: 30 feet "
                   "Components: V, S, M (a drop of blood) Duration: Concentration, up to 1 minute")],
    "Banishment": [(112, "Level 4 Abjuration (Cleric, Paladin, Sorcerer, Warlock, Wizard) Casting Time: Action "
                         "Range: 30 feet Components: V, S, M (a pentacle) Duration: Concentration, up to 1 minute")],
}
CASTING = {"action": "Action", "1-minute": "1 minute or Ritual", "1-hour": "1 hour", "8-hours": "8 hours"}
DURATION = {60: "1 minute", 600: "10 minutes", 3600: "1 hour", 86400: "24 hours", 864000: "10 days"}


def header_from_record(r):
    """Rebuild the source header from the structured record alone."""
    lists = ", ".join(name.capitalize() for name in r["lists"])
    rng = {"touch": "Touch", "self": "Self"}.get(r["range"]["kind"], f"{r['range']['feet']} feet")
    c = r["components"]
    parts = [letter for letter, on in (("V", c["verbal"]), ("S", c["somatic"])) if on]
    if c["material"]:
        parts.append(f"M ({c['material']['description']})")
    if r["duration"] == "timed":
        duration = DURATION[r["duration_seconds"]]
        if r["concentration"]:
            duration = "Concentration, up to " + duration
    else:
        duration = {"instantaneous": "Instantaneous", "until-dispelled": "Until dispelled"}[r["duration"]]
    return (f"Level {r['level']} {r['school'].capitalize()} ({lists}) Casting Time: {CASTING[r['casting_time']]} "
            f"Range: {rng} Components: {', '.join(parts)} Duration: {duration}")


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_197.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (109, 110, 111, 112):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_197.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    owned = [r for r in inventory["records"] if r.get("catalogue_issue") == 197]
    records = catalogue["records"]
    assert len(records) == 12
    assert [r["name"] for r in records] == NAMES
    assert [(r["id"], r["name"], r["source_page"]) for r in records] == \
        [(r["id"], r["name"], r["source_page"]) for r in owned]
    assert all(r["family"] == "spells" for r in owned)
    assert catalogue["source"]["sha256"] == PIN
    assert catalogue["source"]["license"] == "CC-BY-4.0"
    for r in records:
        assert r["name"] in pages[r["source_page"]], r["name"]
        assert r["scope"] == "complete-source-record"
        assert r["slot_levels"] == list(range(r["level"], 10)), r["name"]
        # The header rebuilt from structured fields must equal the source header,
        # and each source header fragment must be on its page.
        rebuilt = header_from_record(r)
        assert rebuilt == " ".join(text for _, text in HEADERS[r["name"]]), (r["name"], rebuilt)
        for page, text in HEADERS[r["name"]]:
            assert text in pages[page], (r["name"], page, text)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_197_test.go; this is not a prose parser.
    anchors = {
        109: ["antipathy or sympathy", "one creature or object that is Huge or smaller",
              "red dragons, goblins, or vampires", "Wisdom saving throw when it comes within 120 feet",
              "Antipathy. The creature has the Frightened", "as far away as possible",
              "Sympathy. The creature has the Charmed", "as close as possible", "is within 5 feet of the target",
              "If the target damages the Charmed creature"],
        110: ["Ending the Effect.", "more than 120 feet away", "immune to it for 1 minute",
              "Invisible, invulnerable eye", "see in every direction", "Darkvision with a range of 30 feet",
              "As a Bonus Action, you can move the eye up to 30 feet", "1 inch in diameter",
              "Large hand of shimmering magical", "AC 20", "equal to your Hit Point maximum",
              "0 Hit Points, the spell ends", "move the hand up to 60 feet", "Clenched Fist.",
              "melee spell attack", "5d8 Force damage", "Forceful Hand.", "Strength saving throw",
              "up to 5 feet plus a number of feet equal to five times your spellcasting ability modifier",
              "Grasping Hand.", "Dexterity saving throw", "Grappled condition",
              "escape DC equal to your spell save DC", "4d6 plus your spellcasting ability modifier",
              "Interposing Hand.", "Half Cover", "Difficult Terrain for your enemies",
              "Clenched Fist increases by 2d8", "Grasping Hand increases by 2d6",
              "slot level above 5", "closed door, window, gate, container, or hatch",
              "any nonmagical means", "any creatures you designate", "within 5 feet of the object, unlocks it for 1 minute",
              "within 5 feet of the sword", "4d12 plus your spellcasting ability modifier",
              "move the sword up to 30 feet", "same target or a different one"],
        111: ["willing crea- ture or an object that isn’t being worn or carried", "every day for 30 days",
              "Mask (Creature).", "False Aura (Object).", "such as Detect Magic",
              "nonmagical object appear magical", "magic item appear nonmagical", "school of magic",
              "up to eight willing creatures", "already on that plane", "Unconscious",
              "doesn’t need food or air, and doesn’t age", "after 1 foot", "both die",
              "0 Hit Points, the spell ends for that target", "Magic action to dismiss",
              "within the next 30 minutes", "Weal Good Woe Bad Weal and woe Good and bad Indifference Neither good nor bad",
              "cumulative 25 percent", "30-foot Emanation", "Resistance to Necrotic damage",
              "Hit Point maximums can’t be reduced", "gains 1 Hit Point",
              "Beast or Plant creature with an Intelligence of 3 or less", "natural plant that isn’t a creature"],
        112: ["Intelligence of 10", "speak one language you know", "Awakened Shrub or Awakened Tree",
              "Charmed condition for 30 days", "chooses its attitude toward you",
              "Up to three creatures of your choice that you can see", "Charisma saving throw",
              "subtract 1d4", "one additional creature for each spell slot level above 1",
              "One creature that you can see within range", "harmless demiplane", "Incapacitated",
              "nearest unoccupied space", "an Aberration, a Celestial, an Ele- mental, a Fey, or a Fiend",
              "if the spell lasts for 1 minute", "one additional creature for each spell slot level above 4"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            assert snippet in pages[page], (page, snippet)
    by_name = {r["name"]: r for r in records}
    details = {name: next(v for v in r["source_details"].values() if isinstance(v, dict))
               for name, r in by_name.items()}
    assert details["Antipathy/Sympathy"]["source_pages"] == [109, 110]
    assert details["Arcanist’s Magic Aura"]["source_pages"] == [110, 111]
    assert details["Awaken"]["source_pages"] == [111, 112]
    assert [m["choice"] for m in details["Antipathy/Sympathy"]["modes"]] == ["antipathy", "sympathy"]
    assert [o["choice"] for o in details["Arcane Hand"]["options"]] == \
        ["clenched-fist", "forceful-hand", "grasping-hand", "interposing-hand"]
    assert details["Arcane Hand"]["slot_dice"] == [
        {"effect_id": "clenched-fist", "dice_count": 2, "die_sides": 8},
        {"effect_id": "grasping-crush", "dice_count": 2, "die_sides": 6}]
    assert [(o["omen"], o["results"]) for o in details["Augury"]["omens"]] == [
        ("Weal", "Good"), ("Woe", "Bad"), ("Weal and woe", "Good and bad"),
        ("Indifference", "Neither good nor bad")]
    costs = {name: (r["components"]["material"] or {}).get("cost_cp") for name, r in by_name.items()}
    assert costs["Arcane Lock"] == 2500 and by_name["Arcane Lock"]["components"]["material"]["consumed"]
    assert costs["Arcane Sword"] == 25000 and not by_name["Arcane Sword"]["components"]["material"]["consumed"]
    assert costs["Astral Projection"] == 110000 and by_name["Astral Projection"]["components"]["material"]["consumed"]
    assert details["Astral Projection"]["jacinth_minimum_cost_cp"] == 100000
    assert details["Astral Projection"]["silver_bar_minimum_cost_cp"] == 10000
    assert costs["Augury"] == 2500 and not by_name["Augury"]["components"]["material"]["consumed"]
    assert by_name["Augury"]["ritual"] and sum(r["ritual"] for r in records) == 1
    assert costs["Awaken"] == 100000 and by_name["Awaken"]["components"]["material"]["consumed"]
    assert by_name["Aura of Life"]["components"]["material"] is None
    assert details["Bane"]["base_targets"] == 3 and details["Banishment"]["base_targets"] == 1
    print("#197 pinned PDF, twelve identities/pages, source headers and edition-specific anchors pass")


if __name__ == "__main__":
    main()
