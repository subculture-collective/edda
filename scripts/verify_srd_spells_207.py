#!/usr/bin/env python3
"""Check issue #207 identities, headers and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
CASTING = {"action": "Action", "1-minute": "1 minute", "10-minutes": "10 minutes"}
UNITS = ((86400, "day"), (3600, "hour"), (60, "minute"))


def span(seconds):
    for size, unit in UNITS:
        if seconds % size == 0:
            n = seconds // size
            return f"{n} {unit}" + ("" if n == 1 else "s")
    raise AssertionError(seconds)


def header(r):
    """Rebuild the printed source header from the structured record."""
    school = r["school"].capitalize()
    lists = ", ".join(c.capitalize() for c in r["lists"])
    casting = CASTING[r["casting_time"]] + (" or Ritual" if r["ritual"] else "")
    rng = r["range"]
    rng = f"{rng['feet']} feet" if rng["kind"] == "feet" else rng["kind"].capitalize()
    c = r["components"]
    parts = [p for p, on in (("V", c["verbal"]), ("S", c["somatic"])) if on]
    if c["material"]:
        parts.append(f"M ({c['material']['description']})")
    if r["duration"] == "instantaneous":
        duration = "Instantaneous"
    elif r["duration"] == "until-dispelled":
        duration = "Until dispelled"
    elif r["duration"] == "special":
        duration = "Special"
    else:
        duration = ("Concentration, up to " if r["concentration"] else "") + span(r["duration_seconds"])
    label = "Components"
    return (f"{r['name']} Level {r['level']} {school} ({lists}) Casting Time: {casting} "
            f"Range: {rng} {label}: {', '.join(parts)} Duration: {duration}")


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_207.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (133, 134, 135, 136):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_207.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 207}
    records = catalogue["records"]
    assert len(records) == 12 and len(expected) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        assert header(r) in pages[r["source_page"]], header(r)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_207_test.go; this is not a prose parser.
    anchors = {
        133: ["20-foot-radius Sphere", "Heavily Ob- scured", "Gust of Wind", "increases by 20 feet for each spell slot level above 1",
              "40,000 square feet of floor space", "height of 30 feet above the floor", "can’t teleport into the area or use portals",
              "Astral Plane, the Ethereal Plane, the Feywild, the Shadowfell, or the Plane Shift", "Aberrations, Celestials, Ele- mentals, Fey, Fiends, and Undead",
              "5d10 Radiant or Necrotic damage", "designate a password", "can’t overlap with the area of another Forbiddance", "every day for 30 days in the same location",
              "consumed on the last casting"],
        134: ["up to 20 feet on a side", "1/2-inch diameter bars spaced 1/2 inch apart", "up to 10 feet on a side", "Charisma saving throw",
              "extends into the Ethereal Plane", "can’t be dispelled by Dispel Magic", "Advantage on D20 Tests", "Disadvantage on attack rolls against it",
              "unaffected by Difficult Terrain", "Swim Speed equal to its Speed", "spend 5 feet of movement", "one additional creature for each spell slot level above 4",
              "60-foot-radius Sphere", "10d6 Cold damage", "6 inches over an area 30 feet square", "lasts for 1 minute", "Strength (Athletics) check against your spell save DC",
              "(to a range of 40 feet)", "After 1 minute", "increases by 1d6 for each spell slot level above 6", "misty cloud", "drops to 0 Hit Points"],
        135: ["Fly Speed of 10 feet", "Resistance to Bludgeoning, Piercing, and Slashing damage", "Immunity to the Prone condition", "Advantage on Strength, Dexterity, and Constitution saving",
              "can’t attack or cast spells", "one additional creature for each spell slot level above 3", "5 to 20 feet in diam- eter", "pseudonym, title, or nick- name doesn’t work",
              "free to act as the GM deems appropriate", "5d10 Psychic damage", "no more than once each day", "level 7 or 8 spell slot, the duration is 365 days", "level 9 spell slot, the spell lasts until",
              "Remove Curse, Greater Restoration, or Wish", "can’t become Undead", "Raise Dead", "giant centipede, spider, or wasp", "Dodge action"],
        136: ["Use the spell slot’s level for the spell’s level in the stat block", "Large Beast, Unaligned", "AC 11 + the spell’s level", "HP 30 + 10 for each spell level above 4",
              "Speed 40 ft., Climb 40 ft., Fly 40 ft. (Wasp only)", "Str 17 +3 +3", "Dex 13 +1 +1", "Con 15 +2 +2", "Int 4 −3 −3", "Wis 14 +2 +2", "Cha 3 −4 −4",
              "Darkvision 60 ft.; Passive Perception 12", "CR None (XP 0; PB equals your Proficiency Bonus)", "half this spell’s level (round down)", "reach 10 ft. Hit: 1d6 + 3 plus the spell’s level Piercing damage plus 1d4 Poison damage",
              "range 60 ft. Hit: 1d10 + 3 plus the spell’s level Bludgeoning damage", "Speed is reduced to 0 until the start of the in- sect’s next turn", "Venomous Spew (Centipede Only)",
              "Poisoned condition until the start of the insect’s next turn", "replace the number you roll with a 15", "indicates that you are being truthful"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Forbiddance"]["ritual"] and by_name["Forbiddance"]["casting_time"] == "10-minutes"
    assert by_name["Gentle Repose"]["ritual"] and by_name["Gentle Repose"]["casting_time"] == "action"
    assert by_name["Forbiddance"]["components"]["material"]["cost_cp"] == 100000
    assert by_name["Forcecage"]["components"]["material"]["consumed"]
    assert by_name["Gate"]["components"]["material"]["cost_cp"] == 500000
    assert by_name["Geas"]["duration_seconds"] == 30 * 86400
    assert by_name["Gaseous Form"]["source_details"]["gaseous"]["source_pages"] == [134, 135]
    assert by_name["Giant Insect"]["source_details"]["insect"]["source_pages"] == [135, 136]
    assert len(by_name["Giant Insect"]["source_details"]["insect"]["giant_insect"]["abilities"]) == 6
    print("#207 pinned PDF, twelve identities/pages, rebuilt headers and source anchors pass")


if __name__ == "__main__":
    main()
