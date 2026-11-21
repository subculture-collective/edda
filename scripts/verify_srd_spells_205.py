#!/usr/bin/env python3
"""Check issue #205 identities and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
CASTING = {
    "etherealness": "Action", "expeditious-retreat": "Bonus Action", "eyebite": "Action", "fabricate": "10 minutes",
    "faerie-fire": "Action", "faithful-hound": "Action", "false-life": "Action", "fear": "Action",
    "feather-fall": "Reaction, which you take when you or a creature you can see within 60 feet of you falls",
    "find-familiar": "1 hour or Ritual", "find-steed": "Action", "find-traps": "Action",
}
DURATIONS = {
    "etherealness": "Up to 8 hours", "expeditious-retreat": "Concentration, up to 10 minutes", "eyebite": "Concentration, up to 1 minute",
    "fabricate": "Instantaneous", "faerie-fire": "Concentration, up to 1 minute", "faithful-hound": "8 hours",
    "false-life": "Instantaneous", "fear": "Concentration, up to 1 minute", "feather-fall": "1 minute",
    "find-familiar": "Instantaneous", "find-steed": "Instantaneous", "find-traps": "Instantaneous",
}


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_205.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (128, 129, 130, 131):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_205.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 205}
    records = catalogue["records"]
    assert len(records) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    # Each record's printed header must match the structured record, so a wrong
    # level, school, list, casting time, range, component or duration fails.
    for r in records:
        slug = r["id"].rsplit("/", 1)[1]
        page_text = pages[r["source_page"]]
        classes = ", ".join(c.capitalize() for c in r["lists"])
        head = f"Level {r['level']} {r['school'].capitalize()} ({classes})"
        rng = {"feet": f"{r['range']['feet']} feet", "touch": "Touch", "self": "Self"}[r["range"]["kind"]]
        comps = (["V"] if r["components"]["verbal"] else []) + (["S"] if r["components"]["somatic"] else [])
        mat = r["components"]["material"]
        comp_text = ", ".join(comps + ([f"M ({mat['description']})"] if mat else []))
        # The source prints the singular label for Find Steed.
        label = "Component:" if slug == "find-steed" else "Components:"
        header = f"{r['name']} {head} Casting Time: {CASTING[slug]} Range: {rng} {label} {comp_text} Duration: {DURATIONS[slug]}"
        assert header in page_text, (r["name"], header)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_205_test.go; this is not a prose parser.
    anchors = {
        128: ["every foot of movement costs an extra foot", "more than 60 feet away", "Force damage equal to twice the number of feet you are moved",
              "a plane that doesn’t border it", "up to three willing creatures (including yourself) for each spell slot level above 7",
              "within 10 feet of you when you cast the spell"],
        129: ["take the Dash action, and until the spell ends", "again as a Bonus Action", "within 60 feet of you that you can see",
              "Wisdom saving throw", "Asleep. The target has the Unconscious condition", "Panicked. The target has the Frightened con",
              "at least 60 feet away from you", "Sickened. The target has the Poisoned condition",
              "Large or smaller object", "10-foot Cube or eight connected 5-foot Cubes", "no larger than Medium (contained within a 5-foot Cube)",
              "Creatures and magic items can’t be created", "Artisan’s Tools", "20-foot Cube within range", "blue, green, or violet light",
              "Dim Light in a 10-", "more than 300 feet apart", "within 30 feet of it", "Truesight with a range of 30 feet",
              "within 5 feet of it", "take 4d8 Force damage", "move the hound up to 30 feet",
              "2d4 + 4 Temporary Hit Points", "5 additional Temporary Hit Points for each spell slot level above 1"],
        130: ["30-foot Cone", "drop whatever it is holding", "Wisdom saving throw", "up to five falling creatures", "60 feet per round",
              "Bat, Cat, Frog, Hawk, Lizard, Octopus, Owl, Rat, Raven, Spider, Weasel", "Challenge Rating of 0", "Celestial, Fey, or Fiend",
              "within 100 feet of you", "Bonus Action, you can see through", "unoccupied space within 30 feet of you", "One Familiar Only",
              "Otherworldly Steed stat block", "a horse, a camel, a dire wolf, or an elk", "Use the spell slot’s level for the spell’s level in the stat block"],
        131: ["sense any trap within range that is within line of sight", "Alarm or Glyph of Warding spell or a mechanical pit trap",
              "an unstable ceiling", "a hidden sinkhole", "reveals that a trap is present but not its location",
              "general nature of the danger posed"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Find Familiar"]["components"]["material"]["cost_cp"] == 1000
    assert by_name["Etherealness"]["scaling"]["above_level"] == 7
    assert by_name["Faithful Hound"]["effects"][2]["dice_count"] == 4
    print("#205 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
