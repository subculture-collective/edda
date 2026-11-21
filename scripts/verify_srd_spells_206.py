#!/usr/bin/env python3
"""Check issue #206 identities and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_206.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (131, 132, 133):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_206.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 206}
    records = catalogue["records"]
    assert len(records) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        assert r["name"] in pages[r["source_page"]], r["name"]
    # Each record's printed header must match the structured record, so a wrong
    # level, school, list, casting time, range, component or duration fails.
    schools = {"abjuration", "conjuration", "divination", "enchantment", "evocation", "illusion", "necromancy", "transmutation"}
    for r in records:
        page_text = pages[r["source_page"]]
        classes = ", ".join(c.capitalize() for c in r["lists"])
        if r["level"] == 0:
            head = f"{r['school'].capitalize()} Cantrip ({classes})"
        else:
            head = f"Level {r['level']} {r['school'].capitalize()} ({classes})"
        assert r["school"] in schools
        casting = {"action": "Action", "bonus-action": "Bonus Action", "1-minute": "1 minute"}[r["casting_time"]]
        if r["ritual"]:
            casting += " or Ritual"
        r_kind = r["range"]
        rng = {"feet": f"{r_kind['feet']} feet", "touch": "Touch", "self": "Self"}[r_kind["kind"]]
        comps = (["V"] if r["components"]["verbal"] else []) + (["S"] if r["components"]["somatic"] else [])
        mat = r["components"]["material"]
        comp_text = ", ".join(comps + ([f"M ({mat['description']})"] if mat else []))
        seconds = r["duration_seconds"]
        durations = {0: "Instantaneous", 60: "1 minute", 600: "10 minutes", 3600: "1 hour", 86400: "1 day"}
        dur = durations[seconds]
        if r["concentration"]:
            dur = "Concentration, up to " + dur
        header = f"{r['name']} {head} Casting Time: {casting} Range: {rng} Components: {comp_text} Duration: {dur}"
        # A header wrapped across two lines keeps its words, so compare on whole words.
        assert header in page_text, (r["name"], header)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_206_test.go; this is not a prose parser.
    anchors = {
        131: ["most direct physical route", "You must be familiar with the location", "on another plane of existence, a moving des",
              "you know how far it is and in what direction it lies", "which path is the most direct",
              "7d8 + 30 Necrotic damage", "Constitution saving throw", "rises at the start of your next turn as a Zombie",
              "follows your verbal orders", "20-foot-radius Sphere centered on that point", "taking 8d6 Fire damage",
              "start burning", "increases by 1d6 for each spell slot level above 3"],
        132: ["1d10 Fire damage", "increases by 1d10", "levels 5 (2d10), 11 (3d10), and 17 (4d10)",
              "Bright Light in a 10-foot radius and Dim Light for an additional 10 feet", "warm shield or a chill shield",
              "Resistance to Cold damage", "Resistance to Fire damage", "within 5 feet of you hits you with a melee attack roll",
              "2d8 Fire dam", "2d8 Cold damage from a chill shield", "up to ten 10-foot Cubes", "7d10 Fire damage",
              "scimitar", "evoke it again as a Bonus Action", "Magic action", "3d6 plus your spellcasting ability modifier",
              "above 2", "10-foot-radius, 40-foot", "5d6 Fire damage and 5d6 Radiant damage", "above 5",
              "5-foot-diameter sphere of fire", "ends its turn within 5 feet of the sphere", "2d6 Fire damage",
              "move the sphere up to 30 feet", "barriers up to 5 feet tall", "pits up to 10 feet wide",
              "Bright Light in a 20-foot radius and Dim Light for an additional 20 feet"],
        133: ["Constitution saving throw", "Restrained condition", "its Speed is 0 until the start of your next turn",
              "Constructs automatically succeed", "three times", "Petrified condition", "needn\u2019t be consecutive",
              "Greater Resto", "3 feet in diameter and 1 inch thick", "floats 3 feet above the ground", "up to 500 pounds",
              "within 20 feet of it", "elevation change of 10 feet or more", "more than 100 feet from the disk",
              "willing creature", "Fly Speed of 60 feet and can hover", "falls if it is still aloft",
              "one additional creature for each spell slot level above 3"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Find the Path"]["components"]["material"]["cost_cp"] == 10000
    assert by_name["Find the Path"]["duration_seconds"] == 86400
    assert by_name["Fire Bolt"]["slot_levels"] == []
    assert by_name["Floating Disk"]["ritual"] is True
    assert by_name["Finger of Death"]["effects"][0]["bonus"] == 30
    print("#206 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
