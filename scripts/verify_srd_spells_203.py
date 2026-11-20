#!/usr/bin/env python3
"""Check issue #203 identities and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_203.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (123, 124, 125, 126):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        # A page footer can fall between a header split across the two columns.
        pages[page] = " ".join(" ".join(columns).replace(f"{page} System Reference Document 5.2.1", "").split())
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_203.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 203}
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
        casting = "Action or Ritual" if r["ritual"] else {"action": "Action", "bonus-action": "Bonus Action"}[r["casting_time"]]
        if r["name"] == "Divine Smite":
            casting += ", which you take immediately after hitting a target with a Melee weapon or an Unarmed Strike"
        r_kind = r["range"]
        rng = {"feet": f"{r_kind['feet']} feet", "touch": "Touch", "self": "Self"}[r_kind["kind"]]
        comps = (["V"] if r["components"]["verbal"] else []) + (["S"] if r["components"]["somatic"] else [])
        mat = r["components"]["material"]
        comp_text = ", ".join(comps + ([f"M ({mat['description']})"] if mat else []))
        seconds = r["duration_seconds"]
        durations = {0: "Instantaneous", 60: "1 minute", 3600: "1 hour"}
        dur = durations[seconds]
        if r["concentration"]:
            dur = "Concentration, up to " + dur
        label = "Component" if r["name"] == "Divine Smite" else "Components"  # the printed page says Component here
        header = f"{r['name']} {head} Casting Time: {casting} Range: {rng} {label}: {comp_text} Duration: {dur}"
        # A header wrapped across two lines keeps its words, so compare on whole words.
        assert header in page_text or header in page_text.replace("- ", ""), (r["name"], header)
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_203_test.go; this is not a prose parser.
    anchors = {
        123: ["You teleport to a location within range", "200 feet straight downward", "300 feet upward to the northwest at a 45-degree angle",
              "teleport one willing creature", "within 5 feet of you when you teleport", "within 5 feet of your destination space"],
        124: ["each take 4d6 Force damage", "teleportation fails", "1 foot shorter or taller", "heavier or lighter", "same basic arrangement of limbs",
              "Study action", "Intelligence (Investigation) check against your spell save DC", "green ray", "wall created by Wall of Force",
              "10d6 + 40 Force damage", "gray dust", "True Resurrection or a Wish spell", "Large or smaller nonmagical object", "10-foot-Cube portion",
              "increases by 3d6 for each spell slot level above 6", "Disadvantage on attack rolls against you", "Break Enchantment",
              "Charmed or Frightened condition", "Charisma saving throw", "sent to the Shadowfell", "sent to the Feywild",
              "spell of level 3 or lower", "DC 10 plus that spell", "automatically end a spell on the target",
              "equal to or less than the level of the spell slot you use", "discordant melody", "3d6 Psychic damage",
              "use its Reaction, if available", "safest route", "half as much damage only",
              "increases by 1d6 for each spell slot level above 1"],
        125: ["one question about a specific goal", "within 7 days", "cryptic rhyme", "cumulative 25 percent", "no answer",
              "extra 1d4 Radiant damage", "extra 2d8 Radiant damage", "increases by 1d8 if the target is a Fiend or an Undead",
              "increases by 1d8 for each spell slot level above 1", "50 Hit Points or fewer", "Divine Word Effects",
              "0–20", "21–30", "31–40", "41–50", "The target dies", "Blinded, Deafened, and Stunned conditions for 1 hour",
              "Blinded and Deafened conditions for 10 minutes", "Deafened condition for 1 minute", "24 hours by any means short of a Wish spell",
              "One Beast you can see within range", "slot of level 5 (up to 10 minutes), 6 (up to 1 hour), or 7+ (up to 8 hours)",
              "Advantage on the save if you or your allies are fighting it", "repeats the save", "telepathic link", "Attack that creature",
              "but must take your own Reaction to do so", "One creature you can see within range must"],
        126: ["telepathic link with the Charmed target", "Fetch that object", "level 9 spell slot (up to 8 hours)"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            # Words hyphenated at a line end rejoin in the de-hyphenated text.
            assert snippet in pages[page] or snippet in pages[page].replace("- ", ""), (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Divination"]["components"]["material"]["cost_cp"] == 2500
    assert by_name["Dimension Door"]["source_details"]["dimension_door"]["source_pages"] == [123, 124]
    assert by_name["Dominate Monster"]["source_details"]["dominate_monster"]["source_pages"] == [125, 126]
    assert by_name["Disintegrate"]["effects"][0]["bonus"] == 40
    print("#203 pinned PDF, twelve identities/pages, headers and source anchors pass")


if __name__ == "__main__":
    main()
