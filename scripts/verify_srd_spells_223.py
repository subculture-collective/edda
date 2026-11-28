#!/usr/bin/env python3
"""Check issue #223 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_223.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (172, 173, 174):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_223.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 223}
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
        60: "1 minute",
        600: "10 minutes",
        3600: "1 hour",
        28800: "8 hours",
        86400: "24 hours",
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
        header = f"{r['name']} {head} Casting Time: {casting} Range: {rng} {label}: {', '.join(parts)} Duration: {duration}"
        # A header line hyphenated at its end rejoins in the de-hyphenated text.
        assert header in text or header in text.replace("- ", ""), (r["name"], header)
        assert f"{r['name']} {school if r['level'] == 0 else 'Level'}" in pages[page], (r["name"], page)
        # Material costs are printed in GP; the record stores copper pieces.
        if mat:
            costs = [int(v.replace(",", "")) for v in re.findall(r"worth ([\d,]+)\+ GP", mat["description"])]
            total = sum(costs) * 100
            # Warding Bond prices each ring of a pair.
            if "a pair of" in mat["description"] and " each" in mat["description"]:
                total *= 2
            assert mat["cost_cp"] == total, r["name"]
            assert mat["consumed"] == ("which the spell consumes" in mat["description"]), r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_223_test.go; this is not a prose parser.
    anchors = {
        172: ["60 feet long, 20 feet high, and 1 foot thick", "20 feet in diameter, 20 feet high, and 1 foot thick", "taking 5d8 Fire damage", "within 10 feet of that side", "creases by 1d8 for each spell slot level above 4", "ten 10-foot-by-10-foot panels", "1/4 inch thick", "A Disintegrate spell destroys the wall", "ten 10-foot-square panels", "taking 10d6 Cold damage", "It has AC 12 and 30 Hit Points per 10-foot section", "taking 5d6 Cold damage"],
        173: ["increases by 2d6 and the damage from passing through the sheet of frigid air increases by 1d6 for each spell slot level above 6", "wall is 6 inches thick", "that are only 3 inches thick", "span greater than 20 feet", "halve the size of each panel", "has AC 15 and 30 Hit Points per inch of thickness", "60 feet long, 10 feet high, and 5 feet thick", "20-foot diameter and is up to 20 feet high and 5 feet thick", "taking 7d8 Pierc", "it must spend 4 feet of movement", "taking 7d8 Slashing damage", "damage increase by 1d8 for each spell slot level above 6", "+1 bonus to AC and saving throws", "separated by more than 60 feet", "up to ten willing creatures of your choice within range"],
        174: ["Up to ten willing creatures of your choice", "fill a 20-foot Cube", "have a depth of 5 feet", "burns away in 1 round, dealing 2d4 Fire damage", "30-foot-radius Sphere", "takes 10d10 Psychic damage", "takes 5d10 Psychic damage", "up to ten willing creatures of your choice within range assume gaseous forms", "Fly Speed of 300 feet", "Reverting takes 1 minute", "descends 60 feet per round for 1 minute", "50 feet long, 15 feet high, and 1 foot thick", "taking 4d8 Blud", "Small or smaller flying creatures"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#223 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
