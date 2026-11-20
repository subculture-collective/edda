#!/usr/bin/env python3
"""Check issue #199 identities and source anchors against the pinned PDF."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_spells_199.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {}
    for page in (114, 115, 116):
        columns = []
        for x in (0, 306):
            result = subprocess.run(
                ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x),
                 "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"],
                check=True, capture_output=True, text=True,
            )
            columns.append(" ".join(result.stdout.split()))
        pages[page] = " ".join(columns)
    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/spells_199.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = {r["id"]: r["source_page"] for r in inventory["records"] if r.get("catalogue_issue") == 199}
    records = catalogue["records"]
    assert len(records) == 12
    assert {r["id"]: r["source_page"] for r in records} == expected
    assert catalogue["source"]["sha256"] == PIN
    for r in records:
        assert r["name"] in pages[r["source_page"]], r["name"]
    # Source-linked numeric and alternative anchors supplement the full literal
    # offline value fixtures in spells_199_test.go; this is not a prose parser.
    anchors = {
        114: ["3d10 Lightning", "10 feet tall with a 60-foot radius", "within 5 feet of",
              "increases by 1d10", "20-foot-radius Sphere", "Immunity to the Charmed and",
              "Indifferent about", "three silver pins", "Three bolts then leap",
              "within 30 feet of the first target", "10d8 Lightning",
              "One additional bolt leaps", "Duration: 1 hour"],
        115: ["Advantage if you or your allies are fighting it", "Friendly to you",
              "one additional creature for each spell slot level above 4",
              "above 1", "1d10 Necrotic", "5 (2d10), 11 (3d10), and 17",
              "Acid, Cold, Fire, Lightning, Poison, or Thun-", "3d8 damage",
              "a diamond worth 50+ GP", "within 30 feet of the target",
              "level 2+ spell slot", "level of the slot expended",
              "60-foot-radius", "8d8 Necrotic", "2d8 for each spell slot level above 6",
              "100+ GP", "jeweled horn for hearing", "glass eye for seeing", "1 mile",
              "10 minutes", "luminous orb", "Casting Time: 1 hour",
              "1,000+ GP", "2,000+ GP"],
        116: ["at least 1 cubic inch", "120 days", "none of the original’s equipment",
              "can’t be revived", "20-foot-radius Sphere of yellow-green",
              "5d8 Poison", "only once per turn", "10 feet away from you",
              "1d8 for each spell slot level above 5", "15-foot Cone",
              "a pinch of colorful sand", "Blinded condition until the end of"],
    }
    for page, snippets in anchors.items():
        for snippet in snippets:
            assert snippet in pages[page], (page, snippet)
    by_name = {r["name"]: r for r in records}
    assert by_name["Clone"]["components"]["material"]["cost_cp"] == 300000
    assert by_name["Chill Touch"]["slot_levels"] == []
    assert by_name["Chain Lightning"]["source_details"]["chain"]["secondary_bolts"] == 3
    print("#199 pinned PDF, twelve identities/pages and source anchors pass")


if __name__ == "__main__":
    main()
