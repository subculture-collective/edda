#!/usr/bin/env python3
"""Check issue #212 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_212.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (145, 146, 147):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_212.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 212}
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
    # offline value fixtures in spells_212_test.go; this is not a prose parser.
    anchors = {
        145: ["within 1,000 feet of you", "within 30 feet—at least", "apparel, jewelry, furniture, tool, or weapon", "thickness of lead blocks", "Speed increases by 10 feet", "base AC be", "13 plus its Dexterity modifier", "target dons armor", "spectral, floating hand", "more than 30 feet away from you", "pour the contents out of a vial", "carry more than 10 pounds", "10-foot-radius, 20-foot-tall Cylinder", "Celestials, Elementals, Fey, Fiends, or Undead", "first succeed on a Charisma saving throw", "Disadvantage on attack rolls"],
        146: ["possessed by or gain the Charmed or Frightened condition", "reverse direction", "creases by 1 hour for each spell slot level above 3", "catatonic state", "up to 100 feet out of the container", "Protection from Evil and Good or Magic Circle", "again for 24 hours", "Your Hit Points, Hit Point Dice, Strength, Dexterity, Constitution, Speed, and senses", "your own spellcasting DC", "When the spell ends, the container is destroyed", "three glowing darts", "1d4 + 1 Force damage", "one more dart for each spell slot level above 1", "25 words or fewer", "as long as 10 minutes", "occur within 30 feet of the object", "silver bell"],
        147: ["+1 bonus to attack rolls and damage rolls", "spell ends early if you cast it again", "+2 with a level 3–5 spell slot", "+3 with a level 6+ spell slot", "5 feet wide and 10 feet tall", "50 contiguous 10-foot Cubes", "nine- course banquet for up to 100 people", "A staff of 100 near-transparent servants", "no larger than a 20-foot Cube", "can’t deal damage or cause conditions", "lasts until dispelled, without requiring Concentration", "level 4+ spell slot", "up to six crea", "30-foot-radius Sphere", "equal to 5d8", "creases by 1d8 for each spell slot level above 5"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#212 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
