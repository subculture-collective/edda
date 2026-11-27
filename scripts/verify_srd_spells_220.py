#!/usr/bin/env python3
"""Check issue #220 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_220.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (164, 165, 166):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_220.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 220}
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
        "bonus-action": "Bonus Action",
    }
    timed = {
        60: "1 minute",
        600: "10 minutes",
        3600: "1 hour",
        28800: "8 hours",
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
            rng = {5280: "1 mile"}.get(feet, f"{feet} feet")
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
    # offline value fixtures in spells_220_test.go; this is not a prose parser.
    anchors = {
        164: ["The corpse must have a mouth", "target of this spell within the past 10 days", "ask the corpse up to five questions", "immobile 30-foot Emanation", "within the past day", "a Climb Speed equal to its Speed", "one additional creature for each spell slot level above 2", "20-foot-radius Sphere", "takes 2d4 Piercing damage for every 5 feet", "15-foot Emanation"],
        165: ["takes 3d8 Radiant damage", "3d8 Necrotic damage", "creases by 1d8 for each spell slot level above 3", "Force damage equal to 1d8 plus your spellcasting ability modifier", "move the force up to 20 feet", "creases by 1d8 for every slot level above 2", "emits Dim Light in a 10-foot radius", "20-foot-radius Sphere of yellow, nauseating gas", "no more than 5 feet in any dimension", "up to two hinges and a latch", "Bludgeoning, Piercing, and Slashing damage", "a radius of 300 feet", "take 2d6 Thunder damage", "takes 4d6 Acid damage", "six bolts of lightning", "10d6 Lightning damage"],
        166: ["takes 2d6 Bludgeoning damage", "Turns 5–10", "1d6 Cold damage", "no more than 25 words"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#220 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
