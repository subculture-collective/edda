#!/usr/bin/env python3
"""Check issue #200 identities and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_200.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (116, 117, 118):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_200.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 200}
    records = catalogue["records"]
    assert len(records) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        assert r["name"] in pages[r["source_page"]], r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_200_test.go; this is not a prose parser.
    anchors = {
        116: ["Command Level 1 Enchantment (Bard, Cleric, Paladin)", "You speak a one-word command",
              "shortest and most direct route", "within 5 feet of you", "drops whatever it is holding",
              "moving away from you by the fastest available means", "has the Prone condition",
              "doesn’t move and takes no action or Bonus Action",
              "one additional creature for each spell slot level above 1",
              "Commune Level 5 Divination (Cleric)", "Casting Time: 1 minute or Ritual", "M (incense)",
              "up to three questions", "cumulative 25 percent chance", "Long Rest",
              "Commune with Nature Level 5 Divination (Druid, Ranger)", "within 3 miles",
              "limited to 300 feet"],
        117: ["castles and settlements", "doesn’t function where nature has been replaced by construction", "Choose three of the following facts",
              "Challenge Rating 10+ creature", "Celestial, an Elemental, a Fey, a Fiend, or an Undead",
              "Locations of bodies of water", "Casting Time: Action or Ritual",
              "a pinch of soot and salt", "It takes about 1 minute to read one page of text",
              "Compulsion Level 4 Enchantment (Bard)", "taking the safest route",
              "a small crystal or glass cone", "60-foot Cone", "8d8 Cold damage", "frozen statue until it thaws",
              "1d8 for each spell slot level above 5", "three nut shells", "10-foot-radius Sphere",
              "roll 1d10 at the start of each of its turns", "1, north; 2, east; 3, south; or 4, west",
              "The target chooses its behavior", "radius increases by 5 feet for each spell slot level above 4",
              "Conjure Animals Level 3 Conjuration (Druid, Ranger)", "Large pack of spectral, intangible animals",
              "wolves, serpents, or birds", "Advantage on Strength saving throws", "up to 30 feet"],
        118: ["3d10 Slashing damage", "1d10 for each spell slot level above 3",
              "Conjure Celestial Level 7 Conjuration (Cleric)", "10-foot-radius, 40-foot-high Cylinder",
              "4d12 plus your spellcasting ability modifier", "6d12 Radiant damage",
              "1d12 for each spell slot level above 7",
              "Conjure Elemental Level 5 Conjuration (Druid, Wizard)", "air (Lightning), earth (Thunder), fire (Fire), or water (Cold)",
              "8d8 damage of the spirit’s type", "4d8 damage of the spirit’s type",
              "Conjure Fey Level 6 Conjuration (Druid)", "3d12 plus your", "Frightened condition until the start of your next turn",
              "within 30 feet of the space it left",
              "Conjure Minor Elementals Level 4 Conjuration (Druid, Wizard)", "15-foot Emanation",
              "extra 2d8 damage", "Acid, Cold, Fire, or Light-", "Terrain for your enemies",
              "1d8 for each spell slot level above 4"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            assert snippet in pages[page], (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Commune"]["ritual"] and by_name["Commune with Nature"]["ritual"]
    assert by_name["Comprehend Languages"]["ritual"]
    assert by_name["Commune with Nature"]["source_details"]["nature"]["source_pages"] == [116, 117]
    assert by_name["Conjure Animals"]["source_details"]["animals"]["source_pages"] == [117, 118]
    assert by_name["Confusion"]["source_details"]["confusion"]["additional_radius_feet_per_slot"] == 5
    print("#200 pinned PDF, twelve identities/pages and source anchors pass")


if __name__ == "__main__":
    main()
