#!/usr/bin/env python3
"""Check issue #214 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_214.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (150, 151, 152):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_214.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 214}
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
        "action": "Action", "bonus-action": "Bonus Action", "1-minute": "1 minute", "10-minutes": "10 minutes",
        "reaction": "Reaction, which you take in response to taking damage from a creature that you can see within 60 feet of yourself",
    }
    timed = {6: "1 round", 60: "1 minute", 600: "10 minutes", 3600: "1 hour", 7200: "2 hours",
             28800: "8 hours", 86400: "24 hours", 864000: "10 days"}
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
            rng = {"touch": "Touch", "self": "Self", "sight": "Sight"}[kind]
        c = r["components"]
        mat = c["material"]
        parts = (["V"] if c["verbal"] else []) + (["S"] if c["somatic"] else []) + ([f"M ({mat['description']})"] if mat else [])
        # Jump prints the singular "Component:" label.
        label = "Component" if r["name"] == "Jump" else "Components"
        if r["duration"] == "instantaneous":
            duration = "Instantaneous"
        elif r["duration"] == "until-dispelled":
            duration = "Until dispelled"
        else:
            duration = timed[r["duration_seconds"]]
        if r["concentration"]:
            duration = "Concentration, up to " + duration
        header = f"{r['name']} {head} Casting Time: {casting} Range: {rng} {label}: {', '.join(parts)} Duration: {duration}"
        assert header in text, (r["name"], header)
        assert f"{r['name']} {school if r['level'] == 0 else 'Level'}" in pages[page], (r["name"], page)
        # Material costs are printed in GP; the record stores copper pieces.
        if mat:
            costs = [int(v.replace(",", "")) for v in re.findall(r"worth ([\d,]+)\+ GP", mat["description"])]
            total = sum(costs) * 100
            if r["name"] == "Legend Lore":
                # 250+ GP of incense and four ivory strips worth 50+ GP each.
                assert costs == [250, 50] and "four ivory strips" in mat["description"]
                total = (250 + 4 * 50) * 100
            assert mat["cost_cp"] == total, r["name"]
            assert mat["consumed"] == ("which the spell consumes" in mat["description"]), r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_214_test.go; this is not a prose parser.
    anchors = {
        150: ["Three illusory duplicates", "roll a d6 for each of your remaining duplicates", "3 or higher", "Blinded condition, Blindsight, or Truesight", "up to twice your Speed", "intangible and invulnerable", "surrounded by silvery mist", "up to 30 feet to an unoccupied space you can see", "last 24 hours and that lasted no more than 10 min", "Remove Curse or Greater Restoration", "7 days ago (level 6 spell slot), 30 days ago", "365 days ago (level 8 spell slot)", "(level 9 spell slot)"],
        151: ["5-foot-radius, 40-foot-high Cylinder", "up to 60 feet", "2d10 Radiant damage", "reverts to its true form", "creases by 1d10 for each spell slot level above 2", "no larger than 40 feet on a side", "half the area’s largest dimension", "pillar up to 20 feet high", "natural stone or stone construction", "no larger than 10 feet in any dimension", "magical scrying sensors", "up to 5 feet wide, 8 feet tall, and 20 feet deep", "30-foot Emana", "+10 bonus to Dex", "leave no tracks", "Intelligence saving throw", "no larger than a 10-foot Cube"],
        152: ["rationalizes any il", "phantasmal bridge", "deal 2d8 Psychic damage", "within 5 feet of the phantasm", "4d10 Psychic damage", "takes half as much damage, and the spell ends", "creases by 1d10 for each spell slot level above 4", "A Large, quasi-real, horselike creature", "saddle, bit, and bridle", "more than 10 feet away from the steed", "Riding Horse stat block", "100 feet and can travel 13 miles in an hour", "rider 1 minute to dismount", "if the steed takes any damage"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#214 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
