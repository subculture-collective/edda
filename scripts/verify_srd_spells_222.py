#!/usr/bin/env python3
"""Check issue #222 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_222.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (169, 170, 171, 172):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_222.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 222}
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
        "1-hour": "1 hour",
    }
    timed = {
        36: "6 rounds",
        60: "1 minute",
        3600: "1 hour",
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
            rng = "1 mile" if feet == 5280 else f"{feet} feet"
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
            # True Strike's weapon is priced in copper pieces.
            total += sum(int(v) for v in re.findall(r"worth ([\d,]+)\+ CP", mat["description"]))
            assert mat["cost_cp"] == total, r["name"]
            assert mat["consumed"] == ("which the spell consumes" in mat["description"]), r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_222_test.go; this is not a prose parser.
    anchors = {
        169: ["understand any spoken or signed language that it hears or sees"],
        170: ["knows at least one language can understand it", "using 5 feet of movement", "within 500 feet", "within 5 feet of the destination tree", "only once on each of your turns", "equal to or less than the target’s Challenge Rating", "retains its Hit Points, Hit Point Dice, alignment, and personality", "Challenge Rating of 9 or lower", "If the spell lasts more than an hour", "no longer than 200 years", "except old age"],
        171: ["within 10 feet of you", "Truesight with a range of 120 feet", "(1d6), 11 (2d6), and 17 (3d6)", "300 feet long, 300 feet high, and 50 feet thick", "taking 6d10 Bludgeoning damage", "moves 50 feet away from you", "take 5d10 Bludgeoning damage", "reduced by 1d10", "It has AC 10, 1 Hit Point, and a Strength of 2", "move up to 15 feet", "more than 60 feet away from you", "takes 3d6 Necrotic damage", "half the amount of Necrotic damage dealt", "creases by 1d6 for each spell slot level above 3"],
        172: ["take 1d6 Psychic damage", "(2d6), 11 (3d6), and 17 (4d6)", "1-foot-diameter ball of acid", "20-foot-radius Sphere", "takes 10d4 Acid damage", "another 5d4 Acid damage", "increases by 2d4 for each spell slot level above 4"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#222 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
