#!/usr/bin/env python3
"""Check issue #208 identities and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_208.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (136, 137, 138, 139):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_208.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 208}
    records = catalogue["records"]
    assert len(records) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        assert r["name"] in pages[r["source_page"]], r["name"]
        # A continuation page list must start at the printed page and be consecutive.
        details = next(iter(v for k, v in r["source_details"].items() if isinstance(v, dict)))
        if "source_pages" in details:
            assert details["source_pages"][0] == r["source_page"]
            assert details["source_pages"] == list(range(details["source_pages"][0], details["source_pages"][-1] + 1))
    # Each record's printed header must match the structured record, so a wrong
    # level, school, list, casting time, range, component or duration fails.
    schools = {"abjuration", "conjuration", "divination", "enchantment", "evocation", "illusion", "necromancy", "transmutation"}
    castings = {"action": "Action", "1-hour": "1 hour", "24-hours": "24 hours"}
    timed = {6: "1 round", 60: "1 minute", 28800: "8 hours", 86400: "24 hours"}
    for r in records:
        page_text = pages[r["source_page"]]
        classes = ", ".join(c.capitalize() for c in r["lists"])
        if r["level"] == 0:
            head = f"{r['school'].capitalize()} Cantrip ({classes})"
        else:
            head = f"Level {r['level']} {r['school'].capitalize()} ({classes})"
        assert r["school"] in schools
        assert r["ritual"] is False
        r_kind = r["range"]
        rng = {"feet": f"{r_kind['feet']} feet", "touch": "Touch", "self": "Self"}[r_kind["kind"]]
        comps = (["V"] if r["components"]["verbal"] else []) + (["S"] if r["components"]["somatic"] else [])
        mat = r["components"]["material"]
        comp_text = ", ".join(comps + ([f"M ({mat['description']})"] if mat else []))
        # Guidance prints the singular "Component:" label.
        label = "Component" if r["name"] == "Guidance" else "Components"
        if r["duration"] == "instantaneous":
            dur = "Instantaneous"
        elif r["duration"] == "until-dispelled":
            dur = "Until dispelled or triggered" if r["source_details"]["end_timing"] == "dispelled-or-triggered" else "Until dispelled"
        else:
            dur = timed[r["duration_seconds"]]
        if r["concentration"]:
            dur = "Concentration, up to " + dur
        header = f"{r['name']} {head} Casting Time: {castings[r['casting_time']]} Range: {rng} {label}: {comp_text} Duration: {dur}"
        # A header wrapped across two lines keeps its words, so compare on whole words.
        assert header in page_text, (r["name"], header)
        # Material costs are printed in GP; the record stores copper pieces.
        if mat and "GP" in mat["description"]:
            gp = int(re.search(r"worth ([\d,]+)\+ GP", mat["description"]).group(1).replace(",", ""))
            assert mat["cost_cp"] == gp * 100, r["name"]
            assert mat["consumed"] == ("which the spell consumes" in mat["description"]), r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_208_test.go; this is not a prose parser.
    anchors = {
        136: ["immobile, shimmering barrier", "Emanation around you", "level 5 or lower", "has no effect on them",
              "above 6", "10 feet in diameter", "more than 10 feet", "Wisdom (Perception) check against your spell save DC",
              "explosive rune or a spell glyph", "20-foot-radius Sphere", "5d8 Acid, Cold, Fire, Lightning, or Thunder damage"],
        137: ["increases by 1d8 for each spell slot level above 3", "level 3 or lower", "Hostile creatures",
              "Ten berries", "Bonus Action to eat one berry", "restores 1 Hit Point", "10- foot square", "Prone condition",
              "Invisible condition", "1 Exhaustion level", "Charmed or Petrified condition", "Attunement to a cursed magic item",
              "ability scores", "Hit Point maximum", "Large spectral guardian", "within 10 feet of the guardian",
              "20 Radiant damage", "total of 60 damage", "2,500 square feet", "20 feet tall", "one 50-foot square",
              "one hundred 5-foot squares", "twenty-five 10-foot squares", "within 5 feet of the warded area"],
        138: ["Dispel Magic has no effect on Guards and Wards itself", "If all four are dispelled", "365 days", "50 percent chance",
              "Arcane Lock spell", "up to ten doors", "regrow in 10 minutes", "Dancing Lights in four corridors",
              "Magic Mouth in two locations", "Stinking Cloud in two locations", "Gust of Wind in one corridor or room",
              "Suggestion in one 5-foot square", "adds 1d4 to any ability check", "ranged spell attack", "4d6 Radiant damage",
              "Advantage", "increases by 1d6 for each spell slot level above 1", "Line of strong wind 60 feet long and 10 feet wide",
              "Strength saving throw", "pushed 15 feet away", "2 feet of movement for every 1 foot", "50 percent chance to extinguish them",
              "Bonus Action on your later turns", "radius up to 60 feet", "Hallowed Ward", "Aberration, Celestial, Elemental, Fey, Fiend, or Undead",
              "Extra Effect", "Courage.", "Darkness. Darkness fills the area"],
        139: ["Daylight. Bright light fills the area", "Peaceful Rest", "Extradimensional Interference", "Fear. Creatures", "Resistance. Creatures",
              "Silence. No sound", "Tongues. Creatures", "Vulnerability. Creatures"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Glyph of Warding"]["components"]["material"]["cost_cp"] == 20000
    assert by_name["Hallow"]["casting_time"] == "24-hours"
    assert by_name["Guidance"]["slot_levels"] == []
    assert by_name["Hallow"]["source_details"]["hallow"]["source_pages"] == [138, 139]
    assert len(by_name["Hallow"]["source_details"]["hallow"]["extra_effects"]) == 10
    print("#208 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
