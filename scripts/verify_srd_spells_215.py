#!/usr/bin/env python3
"""Check issue #215 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_215.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (152, 153, 154, 155):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_215.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 215}
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
        "10-minutes": "10 minutes",
        "1-hour": "1 hour",
    }
    timed = {
        3600: "1 hour",
        86400: "24 hours",
    }
    # Headers whose printed casting time the record splits into typed fields.
    printed_casting = {
        "Plant Growth": "Action (Overgrowth) or 8 hours (Enrichment)",
    }
    # Headers with a printed duration form of their own.
    printed_duration = {
        "Prestidigitation": "Up to 1 hour",
    }
    for r in records:
        page = r["source_page"]
        # A header may continue onto the next page.
        text = pages[page] + " " + pages.get(page + 1, "")
        classes = ", ".join(c.capitalize() for c in r["lists"])
        school = r["school"].capitalize()
        head = f"{school} Cantrip ({classes})" if r["level"] == 0 else f"Level {r['level']} {school} ({classes})"
        casting = printed_casting.get(r["name"]) or castings[r["casting_time"]] + (" or Ritual" if r["ritual"] else "")
        kind, feet = r["range"]["kind"], r["range"]["feet"]
        if kind == "feet":
            rng = f"{feet} feet"
        else:
            rng = {"touch": "Touch", "self": "Self"}[kind]
        c = r["components"]
        mat = c["material"]
        parts = (["V"] if c["verbal"] else []) + (["S"] if c["somatic"] else []) + ([f"M ({mat['description']})"] if mat else [])
        # These headers print the singular "Component:" label.
        label = "Component" if r["name"] in {"Power Word Heal", "Power Word Kill"} else "Components"
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
    # offline value fixtures in spells_215_test.go; this is not a prose parser.
    anchors = {
        152: ["an otherworldly entity for aid", "a Celestial, an Elemental, or a Fiend loyal to it", "payment worth 100 GP per minute", "1,000 GP per hour", "up to 10 days) requires 10,000 GP per day", "a jewel worth 1,000+ GP"],
        153: ["extended to match the duration of this spell", "7 (30 days), 8 (180 days), and 9 (366 days)", "You and up to eight willing creatures who link hands in a circle", "100-foot-radius Sphere", "spend 4 feet of movement for every 1 foot it moves", "half-mile radius", "enriched for 365 days", "twice the normal amount of food", "takes 1d12 Poison damage", "5 (2d12), 11 (3d12), and 17 (4d12)", "shape-shift into a Beast form", "retains its alignment, personality, creature type, Hit Points, and Hit Point Dice"],
        154: ["gear melds into the new form", "regains all its Hit Points", "Charmed, Frightened, Paralyzed, Poisoned, or Stunned condition", "100 Hit Points or fewer, it dies", "takes 12d12 Psychic damage", "150 Hit Points or fewer, it has the Stunned condition", "Up to five creatures of your choice", "regain 2d8 Hit Points", "creases by 1d8 for each spell slot level above 2", "up to three of its non-instantaneous effects active", "1 cubic foot", "Eight rays of light flash from you in a 60-foot Cone", "Red. Failed Save: 12d6 Fire damage", "Green. Failed Save: 12d6 Poison damage"],
        155: ["Blue. Failed Save: 12d6 Cold damage", "saves three times, the condition ends", "teleports to another plane of existence", "The target is struck by two rays. Roll twice, rerolling any 8"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#215 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
