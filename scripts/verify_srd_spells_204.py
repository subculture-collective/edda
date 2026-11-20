#!/usr/bin/env python3
"""Check issue #204 identities and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_204.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (126, 127, 128):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_204.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 204}
    records = catalogue["records"]
    assert len(records) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        assert r["name"] in pages[r["source_page"]], r["name"]
    # Each record's printed header must match the structured record, so a wrong
    # level, school, list, casting time, range, component or duration fails.
    schools = {"abjuration", "conjuration", "divination", "enchantment", "evocation", "illusion", "necromancy", "transmutation"}
    for r in records:
        page_text = pages[r["source_page"]]
        classes = ", ".join(c.capitalize() for c in r["lists"])
        if r["level"] == 0:
            head = f"{r['school'].capitalize()} Cantrip ({classes})"
        else:
            head = f"Level {r['level']} {r['school'].capitalize()} ({classes})"
        assert r["school"] in schools
        casting = {"action": "Action", "1-minute": "1 minute",
                   "bonus-action": "Bonus Action, which you take immediately after hitting a creature with a weapon" if r["casting_time"] == "bonus-action" and r["range"]["kind"] == "self" else "Bonus Action"}[r["casting_time"]]
        r_kind = r["range"]
        rng = {"feet": f"{r_kind['feet']} feet", "touch": "Touch", "self": "Self", "special": "Special"}[r_kind["kind"]]
        comps = (["V"] if r["components"]["verbal"] else []) + (["S"] if r["components"]["somatic"] else [])
        mat = r["components"]["material"]
        comp_text = ", ".join(comps + ([f"M ({mat['description']})"] if mat else []))
        seconds = r["duration_seconds"]
        durations = {0: "Instantaneous", 60: "1 minute", 3600: "1 hour", 28800: "8 hours"}
        dur = durations[seconds]
        if r["concentration"]:
            dur = "Concentration, up to " + dur
        header = f"{r['name']} {head} Casting Time: {casting} Range: {rng} Components: {comp_text} Duration: {dur}"
        # A header wrapped across two lines keeps its words, so compare on whole words.
        assert header in page_text or header in page_text.replace("- ", ""), (r["name"], header)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_204_test.go; this is not a prose parser.
    anchors = {
        126: ["One Humanoid you can see within range", "Advantage on the save if you or your allies are fighting it",
              "repeats the save, ending the spell on itself on a success", "no action required",
              "spell slot of level 6 (up to 10 minutes), 7 (up to 1 hour), or 8+ (up to 8 hours)",
              "Acid, Cold, Fire, Lightning, or Poison", "15-foot Cone", "3d6 damage of the chosen type",
              "increases by 1d6 for each spell slot level above 2", "a trance state to act as a dream messenger",
              "no more than ten words", "3d6 Psychic damage when it wakes up",
              "Weather Sensor", "for the next 24 hours", "This effect persists for 1 round", "fit in a 5-foot Cube",
              "a candle, a torch, or a campfire"],
        127: ["100-foot-radius circle", "Difficult Terrain", "1d6 fissures", "1d10 × 10 feet deep and 10 feet wide",
              "50 Bludgeoning damage", "12d6 Bludgeoning damage", "DC 20 Strength (Athletics) check",
              "1d10 Force damage", "two beams at level 5, three beams at level 11, and four beams at level 17",
              "Beckon Air", "Beckon Earth", "Beckon Fire", "Beckon Water", "Sculpt Element", "1-foot Cube", "for 1 hour",
              "Advantage on ability checks using the chosen ability", "one additional creature for each spell slot level above 2",
              "enlarges or reduces a creature or an object you can see within range"],
        128: ["Constitution saving throw", "extra 1d4 damage", "1d4 less damage on a hit", "can’t reduce the damage below 1",
              "A Large or larger creature has Advantage on this save", "1d6 Piercing damage at the start of each of its turns",
              "increases by 1d6 for each spell slot level above 1", "20-foot square within range", "Difficult Terrain",
              "Strength (Athletics) check against your spell save DC", "−10 penalty to Wisdom (Perception) checks and Passive Perception"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Dream"]["range"]["kind"] == "special"
    assert by_name["Dream"]["source_details"]["casting_seconds"] == 60
    assert by_name["Earthquake"]["source_details"]["earthquake"]["collapse_escape_dc"] == 20
    assert by_name["Enthrall"]["source_details"]["enthrall"]["failure_penalty"] == -10
    assert by_name["Enlarge/Reduce"]["source_details"]["enlarge"]["source_pages"] == [127, 128]
    assert [s["beams"] for s in by_name["Eldritch Blast"]["source_details"]["blast"]["beam_steps"]] == [2, 3, 4]
    print("#204 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
