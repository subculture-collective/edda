#!/usr/bin/env python3
"""Compare issue #252 armor records with the Armor table in the pinned PDF.

The expected records are built from the printed table rows and their category
headings, not from a list kept in this script.
"""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGE = 92
ISSUE = 252
BATCH = "data-armor-02"
COUNT = 1
HEADING = re.compile(r"(Light|Medium|Heavy) Armor \((\d+) Minutes? to (?:(Don or Doff)|Don and (\d+) Minutes? to Doff)\)")
SHIELD_HEADING = "Shield (Utilize Action to Don or Doff)"
ROW = re.compile(
    r"(?P<name>[A-Z][a-z]+(?: [A-Z][a-z]+)*)\s{2,}"
    r"(?P<ac>\d+ \+ Dex modifier(?: \(max \d+\))?|\+\d+|\d+)\s{2,}"
    r"(?P<strength>—|Str \d+)\s{2,}(?P<stealth>—|Disadvantage)\s{2,}"
    r"(?P<weight>\d+) lb\.\s+(?P<amount>\d+(?:,\d{3})*) (?P<currency>CP|SP|GP)"
)


def table_rows(page):
    """Rows in print order, each with its category heading's don and doff time."""
    rows, section, started = [], None, False
    for raw in page.splitlines():
        line = raw.strip()
        if not started:
            started = line.startswith("Armor ") and "Armor Class (AC)" in line and line.endswith("Cost")
            continue
        if not line or line.endswith("System Reference Document 5.2.1"):
            continue
        heading = HEADING.fullmatch(line)
        if heading:
            don = int(heading[2])
            doff = don if heading[3] else int(heading[4])
            section = (heading[1].lower(), {"action": "none", "don_minutes": don, "doff_minutes": doff})
            continue
        if line == SHIELD_HEADING:
            section = ("shield", {"action": "utilize", "don_minutes": 0, "doff_minutes": 0})
            continue
        match = ROW.fullmatch(line)
        assert match and section, ("unparsed table line", line)
        rows.append(dict(match.groupdict(), section=section))
    return rows


def record(row):
    """The catalogue payload, mechanic keys, equip mode and cost of a row."""
    category, times = row["section"]
    shield = category == "shield"
    armor_class = {"base": 0, "bonus": 0, "dexterity": "none", "dexterity_maximum": 0}
    if row["ac"].startswith("+"):
        assert shield, row
        armor_class["bonus"] = int(row["ac"][1:])
        class_key = f"armor-class.bonus-{armor_class['bonus']}"
    else:
        found = re.fullmatch(r"(\d+)( \+ Dex modifier(?: \(max (\d+)\))?)?", row["ac"])
        assert found and not shield, row
        armor_class["base"] = int(found[1])
        class_key = f"armor-class.{found[1]}"
        if found[3]:
            armor_class.update(dexterity="capped", dexterity_maximum=int(found[3]))
            class_key += f"-plus-dexterity-max-{found[3]}"
        elif found[2]:
            armor_class["dexterity"] = "full"
            class_key += "-plus-dexterity"
    strength = 0 if row["strength"] == "—" else int(row["strength"].removeprefix("Str "))
    stealth = row["stealth"] == "Disadvantage"
    payload = {"armor_category": category, "armor_class": armor_class, "strength_requirement": strength,
               "stealth_disadvantage": stealth, "don_doff": times, "weight_pounds": int(row["weight"])}
    equip = "wield" if shield else "wear"
    keys = [f"item.{equip}", f"armor.{category}", "armor.training", class_key]
    if strength:
        keys.append(f"strength.{strength}")
    if stealth:
        keys.append("stealth.disadvantage")
    if shield:
        keys += ["don.utilize-action", "doff.utilize-action"]
    else:
        keys += [f"don.minutes-{times['don_minutes']}", f"doff.minutes-{times['doff_minutes']}"]
    cost = {"kind": "listed", "amount": int(row["amount"].replace(",", "")), "currency": row["currency"].lower(),
            "source_pages": [PAGE, PAGE]}
    return payload, keys, equip, cost


def same(got, want):
    """Equality that does not treat False, 0 and 0.0 as the same value."""
    return json.dumps(got, sort_keys=True) == json.dumps(want, sort_keys=True)


def main():
    if len(sys.argv) != 2:
        raise SystemExit(f"usage: verify_srd_items_{ISSUE}.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    args = ["pdftotext", "-f", str(PAGE), "-l", str(PAGE), "-layout", str(pdf), "-"]
    rows = table_rows(subprocess.run(args, check=True, capture_output=True, text=True).stdout)
    assert len(rows) == 13, len(rows)
    printed = {r["name"]: r for r in rows}
    assert len(printed) == 13

    catalogue = json.loads((ROOT / f"internal/rules/srd521/testdata/items_{ISSUE}.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    armor = [r for r in inventory["records"] if r["family"] == "armor"]
    # The planning inventory lists exactly the printed armor rows.
    assert sorted(printed) == sorted(r["name"] for r in armor)
    expected = [(r["id"], r["name"], r["source_page"]) for r in armor if r.get("catalogue_issue") == ISSUE]
    assert len(expected) == COUNT
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == ISSUE and catalogue["batch"] == BATCH
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected

    for f in families:
        name = f["name"]
        slug = name.lower().replace(" ", "-")
        payload, keys, equip, cost = record(printed[name])
        assert f["id"] == f"srd-5.2.1:content/items/{slug}" and f["kind"] == "mundane" and f["category"] == "armor", name
        assert f["base"] == {"options": [], "excluded": []} and len(f["presentations"]) == 1, name
        assert f["presentations"][0] == {"id": f"srd-5.2.1:content/armor/{slug}", "pages": [PAGE, PAGE]}, name
        assert len(f["variants"]) == 1, name
        v = f["variants"][0]
        assert v["id"] == f["id"] + "/variant/standard" and v["name"] == name and v["rarity"] == "none", name
        assert v["attunement"] == {"mode": "none", "restrictions": []}, name
        assert same(v["charges"], {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}), name
        assert v["equip"] == {"mode": equip, "restrictions": []} and v["use"] == {"mode": "none", "restrictions": []}, name
        assert same(v["cost"], cost), (name, v["cost"], cost)
        assert v["required_mechanic_keys"] == keys, (name, keys)
        assert same(v["details"], {"armor": payload}), (name, v["details"], payload)
    print(f"ok: {len(families)} armor records match their printed table rows, p. {PAGE}")


main()
