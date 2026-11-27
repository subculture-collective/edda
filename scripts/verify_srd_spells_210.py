#!/usr/bin/env python3
"""Check issue #210 identities, printed headers and source anchors against the pinned PDF."""
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
        raise SystemExit("usage: verify_srd_spells_210.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, and
    # drop the page footer that sits between a header split across columns.
    pages = {}
    for page in (141, 142, 143):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(FOOTER.sub("", " ".join(result.stdout.split())))
        pages[page] = " ".join(" ".join(columns).split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_210.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: (r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 210}
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
    # offline value fixtures in spells_210_test.go; this is not a prose parser.
    anchors = {
        141: ["Choose a Humanoid that you can see within range", "one additional Humanoid for each spell slot level above 2", "30-foot Em", "Advantage on all saving throws", "Fiend or an Undead hits", "Blinded condition until the end of its next turn", "1d6 Force damage", "Wisdom (Perception or Survival)", "3–4 (up to 8 hours) or 5+ (up to 24 hours)", "30-foot Cube within range", "Charmed, the creature has the Inca", "Speed of 0", "ranged spell attack", "1d10 Piercing damage", "within 5 feet of it", "take 2d6 Cold damage", "increases by 1d6 for each spell slot level above 1", "20-foot-radius, 40-foot-high Cylinder", "2d10 Bludgeoning damage and 4d6 Cold damage", "until the end of your next turn"],
        142: ["creases by 1d10 for each spell slot level above 4", "Attunement, and how many charges", "ongoing spells are affecting the item", "If you instead touch a creature", "unknown or magical script", "language of the text", "A creature that has Truesight can read the hidden message", "immune to this", "doesn’t need to breathe, eat, or drink", "Burial. The target is entombed", "Chaining. Chains firmly rooted", "Hedged Prison.", "Minimus Containment. The target becomes 1 inch", "Slumber. The target has the Unconscious condition", "within the next decade", "cast with a level 9 spell slot", "20-foot-radius Sphere centered on a point within range", "Heavily Obscured", "10d8 Fire damage", "only once per turn"],
        143: ["cloud moves 10 feet away from you", "2d10 Necrotic damage", "creases by 1d10 for each spell slot level above 1", "Swarming locusts fill a 20-foot-radius Sphere", "Lightly Obscured and Difficult Terrain", "4d10 Piercing damage", "creases by 1d10 for each spell slot level above 5"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    print("#210 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
