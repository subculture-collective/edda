#!/usr/bin/env python3
"""Check issue #201 identities, headers and source anchors against the pinned PDF."""
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
    if r["casting_time"] == "reaction":
        casting = "Reaction, which you take when you see a creature within 60 feet of yourself casting a spell with Verbal, Somatic, or Material components"
    else:
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
    # Contagion's source heading prints the singular "Component".
    label = "Component" if r["name"] == "Contagion" and len(parts) == 2 else "Components"
    return (f"{r['name']} Level {r['level']} {school} ({lists}) Casting Time: {casting} "
            f"Range: {rng} {label}: {', '.join(parts)} Duration: {duration}")


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_201.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (118, 119, 120, 121):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_201.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 201}
    records = catalogue["records"]
    assert len(records) == 12 and len(expected) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        assert header(r) in pages[r["source_page"]], header(r)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_201_test.go; this is not a prose parser.
    anchors = {
        118: ["Conjure Woodland Beings", "10-foot Emanation for the duration"],
        119: ["Wisdom saving throw", "5d8 Force damage", "only once per turn",
              "Disengage action as a Bonus Action", "creases by 1d8 for each spell slot level above 4",
              "DC 15 Intelligence saving throw", "up to five questions", "6d6 Psychic damage",
              "Incapacitated condition until you finish a Long Rest", "Greater Restoration",
              "11d8 Necrotic damage", "three successes or failures", "lasts for 7 days",
              "level 5 or lower", "casting time of an action", "only one Contingency spell",
              "Bright Light in a 20-foot radius and Dim Light for an additional 20 feet",
              "Cube up to 100 feet on a side", "rise by as much as 20 feet"],
        120: ["25 percent chance of capsizing", "Whirlpool", "at least 50 feet square and 25 feet deep", "5 feet wide at the base",
              "up to 50 feet wide at the top", "25 feet tall", "pulled 10 feet toward it",
              "2d8 Bludgeoning", "Strength (Athletics)", "within 5 miles of you",
              "1d4 × 10 minutes", "Torrential rain, driving hail, or blizzard", "Heat wave",
              "Freezing", "Moderate wind", "Gale", "the slot isn’t expended",
              "45 pounds of food and 30 gallons of fresh water", "spoils after 24 hours"],
        121: ["up to 10 gallons", "10 additional gallons of water", "increases by 5 feet",
              "can cast this spell only at night", "up to three corpses of Medium or Small Humanoids",
              "within 120 feet of you", "24 hours", "level 7 spell slot", "four Ghouls",
              "five Ghouls", "two Ghasts or Wights", "six Ghouls, three Ghasts or Wights, or two Mummies",
              "no larger than a 5-foot Cube", "Vegetable matter 24 hours", "Stone or crystal 12 hours",
              "Precious metals 1 hour", "Gems 10 minutes", "Adamantine or mithral 1 minute",
              "Cube in- creases by 5 feet for each spell slot level above 5"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            assert snippet in pages[page], (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Contingency"]["components"]["material"]["cost_cp"] == 150000
    assert by_name["Continual Flame"]["components"]["material"]["consumed"]
    assert by_name["Create Undead"]["components"]["material"]["cost_cp"] == 15000
    assert by_name["Contact Other Plane"]["ritual"]
    assert by_name["Creation"]["duration"] == "special"
    assert by_name["Counterspell"]["source_details"]["casting_triggers"]
    assert by_name["Control Water"]["source_details"]["water_control"]["source_pages"] == [119, 120]
    assert by_name["Conjure Woodland Beings"]["source_details"]["woodland"]["source_pages"] == [118, 119]
    print("#201 pinned PDF, twelve identities/pages, rebuilt headers and source anchors pass")


if __name__ == "__main__":
    main()
