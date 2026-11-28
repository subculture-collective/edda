#!/usr/bin/env python3
"""Check issue #224 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_224.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (175,):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_224.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 224}
    records = catalogue["records"]
    assert len(records) == 3
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
    }
    timed = {
        600: "10 minutes",
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
            assert mat["cost_cp"] == total, r["name"]
            assert mat["consumed"] == ("which the spell consumes" in mat["description"]), r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_224_test.go; this is not a prose parser.
    anchors = {
        175: ["any other spell of level 8 or lower", "up to 25,000 GP in value", "no more than 300 feet in any dimension", "up to twenty creatures that you can see", "up to ten creatures that you can see Resistance", "This Resistance is permanent", "for 8 hours", "replace one of your feats", "made within the last round", "you take 1d10 Necrotic damage per level of that spell", "your Strength score becomes 3 for 2d4 days", "decreases by 2 days", "33 percent chance", "your wish fails", "up to five willing creatures within 5 feet", "15-foot-radius Sphere", "makes a Charisma sav"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#224 pinned PDF, three identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
