#!/usr/bin/env python3
"""Check issue #221 identities, printed headers and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
FOOTER = re.compile(r"\b\d+ System Reference Document 5\.2\.1\b")


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_221.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (166, 167, 168, 169):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_221.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 221}
    records = catalogue["records"]
    assert len(records) == 12
    assert {r["id"]: (r["name"], r["source_page"]) for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        details = next(v for v in r["source_details"].values() if isinstance(v, dict))
        # A continuation page list must start at the printed page and be consecutive.
        if "source_pages" in details:
            span = details["source_pages"]
            assert span[0] == r["source_page"] and len(span) > 1
            assert span == list(range(span[0], span[-1] + 1))
    # Each record's printed header must match the structured record, so a wrong
    # level, school, list, casting time, range, component or duration fails.
    castings = {
        "action": "Action",
        "1-minute": "1 minute",
    }
    timed = {
        6: "1 round",
        60: "1 minute",
        600: "10 minutes",
        3600: "1 hour",
        28800: "8 hours",
    }
    # Headers with a printed duration form of their own.
    printed_duration = {
        "Symbol": "Until dispelled or triggered",
        "Thaumaturgy": "Up to 1 minute",
    }
    for r in records:
        page = r["source_page"]
        # A header may continue onto the next page.
        text = pages[page] + " " + pages.get(page + 1, "")
        classes = ", ".join(c.capitalize() for c in r["lists"])
        school = r["school"].capitalize()
        head = f"{school} Cantrip ({classes})" if r["level"] == 0 else f"Level {r['level']} {school} ({classes})"
        casting = castings[r["casting_time"]] + (" or Ritual" if r["ritual"] else "")
        kind, feet = r["range"]["kind"], r["range"]["feet"]
        if kind == "feet":
            rng = f"{feet} feet"
        else:
            rng = {"touch": "Touch", "self": "Self"}[kind]
        c = r["components"]
        mat = c["material"]
        parts = (["V"] if c["verbal"] else []) + (["S"] if c["somatic"] else []) + ([f"M ({mat['description']})"] if mat else [])
        label = "Components"
        if r["duration"] == "instantaneous":
            duration = "Instantaneous"
        elif r["duration"] == "until-dispelled":
            duration = "Until dispelled"
        else:
            duration = timed.get(r["duration_seconds"], "")
        if r["concentration"]:
            duration = "Concentration, up to " + duration
        duration = printed_duration.get(r["name"], duration)
        header = f"{r['name']} {head} Casting Time: {casting} Range: {rng} {label}: {', '.join(parts)} Duration: {duration}"
        # A header line hyphenated at its end rejoins in the de-hyphenated text.
        assert header in text or header in text.replace("- ", ""), (r["name"], header)
        assert f"{r['name']} {school if r['level'] == 0 else 'Level'}" in pages[page], (r["name"], page)
        # Material costs are printed in GP; the record stores copper pieces.
        if mat:
            costs = [int(v.replace(",", "")) for v in re.findall(r"worth ([\d,]+)\+ GP", mat["description"])]
            total = sum(costs) * 100
            assert mat["cost_cp"] == total, r["name"]
            assert mat["consumed"] == ("which the spell consumes" in mat["description"]), r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_221_test.go; this is not a prose parser.
    anchors = {
        166: ["uses the Draconic Spirit stat block", "Use the spell slot’s level for the spell’s level in the stat block", "Large Dragon, Neutral", "AC 14 + the spell’s level", "HP 50 + 10 for each spell level above 5", "Speed 30 ft., Fly 60 ft., Swim 30 ft.", "Str 19 +4 +4 Dex 14 +2 +2 Con 17 +3 +3", "Int 10 +0 +0 Wis 14 +2 +2 Cha 14 +2 +2", "Resistances Acid, Cold, Fire, Lightning, Poison", "Blindsight 30 ft., Darkvision 60 ft.", "Passive Perception 12", "half the spell’s level (round down)", "reach 10 feet. Hit: 1d6 + 4 + the spell’s level", "30-foot Cone. Failure: 2d6 damage", "5-foot-wide, 60-foot-long Line", "takes 6d8 Radiant damage", "Bright Light in a 30-foot radius"],
        167: ["60-foot-radius Sphere", "takes 12d6 Radiant damage", "larger than 10 feet in diameter", "moved more than 10 feet from where you cast this spell", "taking 10d10 Necrotic damage", "move at least 30 feet away from the glyph", "Unconscious condition for 10 minutes", "move a Huge or smaller creature", "move it up to 30 feet in any direction"],
        168: ["telekinetic grip, such as manipulating a simple tool,", "up to eight willing creatures of your choice within range", "rolls 1d100 and consults", "Very familiar 01–05 06–13 14–24 25–00", "Seen casually 01–33 34–43 44–53 54–00", "01–43 44–53 54–73 74–00", "False destination 01–50 51–00", "within the last six months", "takes 3d10 Force damage", "2d12 miles away", "5-foot-radius circle", "within 5 feet of the destination circle"],
        169: ["sigil sequences for two destinations", "studying it for 1 minute", "every day for 365 days", "up to three of its 1-minute effects active", "up to three times as loud as normal", "15-foot Cube originating from you", "takes 2d8 Thunder damage and is pushed 10 feet", "audible within 300 feet", "creases by 1d8 for each spell slot level above 1", "take 1d4 + 1 turns in a row", "more than 1,000 feet", "A 10-foot Emanation", "Spells of level 3 or lower"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#221 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
