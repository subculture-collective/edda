#!/usr/bin/env python3
"""Check issue #202 identities and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_202.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (121, 122, 123):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_202.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 202}
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
        casting = "Action or Ritual" if r["ritual"] else "Action"
        r_kind = r["range"]
        rng = {"feet": f"{r_kind['feet']} feet", "touch": "Touch", "self": "Self"}[r_kind["kind"]]
        comps = (["V"] if r["components"]["verbal"] else []) + (["S"] if r["components"]["somatic"] else [])
        mat = r["components"]["material"]
        comp_text = ", ".join(comps + ([f"M ({mat['description']})"] if mat else []))
        seconds = r["duration_seconds"]
        durations = {0: "Instantaneous", 60: "1 minute", 600: "10 minutes", 3600: "1 hour", 28800: "8 hours"}
        dur = durations[seconds]
        if r["concentration"]:
            dur = "Concentration, up to " + dur
        header = f"{r['name']} {head} Casting Time: {casting} Range: {rng} Components: {comp_text} Duration: {dur}"
        # A header wrapped across two lines keeps its words, so compare on whole words.
        assert header in page_text, (r["name"], header)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_202_test.go; this is not a prose parser.
    anchors = {
        121: ["equal to 2d8 plus your spellcasting ability modifier",
              "increases by 2d8 for each spell slot level above 1"],
        122: ["up to four torch-size lights", "Dim Light in a 10-", "60 feet to a space within range", "within 20 feet of another light",
              "15-foot-radius Sphere", "15-foot Emanation", "level 2 or lower",
              "Darkvision with a range of 150 feet", "60-foot-radius Sphere",
              "additional 60 feet", "60-foot Emanation", "level 3 or lower",
              "drops to 1 Hit Point", "kill it instantly",
              "20-foot-radius Sphere centered on that point",
              "base damage is 12d6", "increases by 1d6 whenever your turn ends",
              "up to 40 feet", "increases by 1d6 for each spell slot level above 7"],
        123: ["shadowy Medium door", "30 feet in each dimension",
              "wood or stone (your choice)", "Prone condition","within 30 feet of yourself", "Aberration, Celestial, Elemental, Fey, Fiend, or Undead",
              "Hallow spell is active", "1 foot of stone, dirt, or wood; 1 inch of metal; or a thin sheet of lead",
              "faint aura", "school of magic", "poisonous or venomous creatures",
              "contagions within 30 feet of yourself", "Sense Thoughts. You sense", "Read Thoughts. Target",
              "Wisdom saving throw", "Intelligence (Arcana) check against your spell save DC"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Detect Thoughts"]["components"]["material"]["cost_cp"] == 1
    assert by_name["Dancing Lights"]["slot_levels"] == []
    assert by_name["Delayed Blast Fireball"]["source_details"]["fireball"]["base_dice_count"] == 12
    print("#202 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
