#!/usr/bin/env python3
"""Compare issue #249 weapon records with the Weapons table in the pinned PDF.

The expected records are built from the printed table rows, not from a list
kept in this script.
"""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGE = 91
SECTION = re.compile(r"(Simple|Martial) (Melee|Ranged) Weapons")
# One space may separate a long name from its damage, as for Heavy Crossbow.
ROW = re.compile(
    r"(?P<name>[A-Z][a-z]+(?: [A-Z][a-z]+)?)\s+(?P<damage>\d+d\d+|\d+) (?P<type>Bludgeoning|Piercing|Slashing)\s+"
    r"(?P<properties>.*?)\s+(?P<mastery>Cleave|Graze|Nick|Push|Sap|Slow|Topple|Vex)\s+"
    r"(?P<weight>—|\d+(?:/\d+)? lb\.)\s+(?P<amount>\d+) (?P<currency>CP|SP|GP)"
)
PROPERTY = re.compile(r"(Ammunition|Finesse|Heavy|Light|Loading|Reach|Thrown|Two-Handed|Versatile)(?: \(([^()]*)\))?")


def table_rows(page):
    """Rows in print order, with a wrapped Properties cell joined to its row."""
    rows, section, started = [], None, False
    for raw in page.splitlines():
        line = raw.strip()
        if not line:
            continue
        if line == "Weapons":
            started = True
            continue
        if not started or line.startswith("Name ") or line.endswith("System Reference Document 5.2.1"):
            continue
        heading = SECTION.fullmatch(line)
        if heading:
            section = (heading[1].lower(), heading[2].lower())
            continue
        match = ROW.fullmatch(line)
        if match:
            assert section, line
            rows.append(dict(match.groupdict(), section=section))
            continue
        # Only a continuation of the previous Properties cell may remain.
        assert rows and PROPERTY.fullmatch(line.rstrip(",")), ("unparsed table line", line)
        rows[-1]["properties"] += " " + line
    return rows


def record(row):
    """The catalogue payload and mechanic keys that a printed row requires."""
    category, attack = row["section"]
    kind = row["type"].lower()
    if "d" in row["damage"]:
        dice, sides = map(int, row["damage"].split("d"))
        damage = {"dice": dice, "sides": sides, "flat": 0, "type": kind}
    else:
        damage = {"dice": 0, "sides": 0, "flat": int(row["damage"]), "type": kind}
    properties, ranges, ammunition = [], {"normal_feet": 0, "long_feet": 0}, "none"
    versatile, exception = {"dice": 0, "sides": 0}, "none"
    if row["properties"] != "—":
        parts = re.split(r",\s*(?![^()]*\))", row["properties"])
        for part in parts:
            match = PROPERTY.fullmatch(part)
            assert match, (row["name"], part)
            name, argument = match[1].lower(), match[2]
            properties.append(name)
            if name in ("ammunition", "thrown"):
                found = re.fullmatch(r"Range (\d+)/(\d+)" + (r"; (Arrow|Bolt|Bullet|Needle)" if name == "ammunition" else ""), argument)
                assert found, (row["name"], argument)
                ranges = {"normal_feet": int(found[1]), "long_feet": int(found[2])}
                if name == "ammunition":
                    ammunition = found[3].lower()
            elif name == "versatile":
                dice, sides = map(int, argument.split("d"))
                versatile = {"dice": dice, "sides": sides}
            elif argument is not None:
                assert (name, argument) == ("two-handed", "unless mounted"), (row["name"], part)
                exception = "unless-mounted"
    weight = {"listed": False, "pounds_numerator": 0, "pounds_denominator": 1}
    if row["weight"] != "—":
        top, _, bottom = row["weight"].removesuffix(" lb.").partition("/")
        weight = {"listed": True, "pounds_numerator": int(top), "pounds_denominator": int(bottom or 1)}
    mastery = row["mastery"].lower()
    payload = {"proficiency_category": category, "attack_kind": attack, "damage": damage, "properties": properties,
               "range": ranges, "ammunition": ammunition, "versatile_damage": versatile,
               "two_handed_exception": exception, "mastery": mastery, "weight": weight}
    keys = ["item.wield", "weapon.attack", f"weapon.{category}", f"weapon.{attack}", f"damage.{row['damage']}-{kind}"]
    keys += [f"property.{p}" for p in properties]
    if ranges["normal_feet"]:
        keys.append(f"range.{ranges['normal_feet']}-{ranges['long_feet']}")
    if ammunition != "none":
        keys.append(f"ammunition.{ammunition}")
    if versatile["dice"]:
        keys.append(f"damage.two-handed-{versatile['dice']}d{versatile['sides']}")
    if exception != "none":
        keys.append(f"property.two-handed-{exception}")
    keys.append(f"mastery.{mastery}")
    cost = {"kind": "listed", "amount": int(row["amount"]), "currency": row["currency"].lower(), "source_pages": [PAGE, PAGE]}
    return payload, keys, cost


def same(got, want):
    """Equality that does not treat False, 0 and 0.0 as the same value."""
    return json.dumps(got, sort_keys=True) == json.dumps(want, sort_keys=True)


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_249.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    args = ["pdftotext", "-f", str(PAGE), "-l", str(PAGE), "-layout", str(pdf), "-"]
    rows = table_rows(subprocess.run(args, check=True, capture_output=True, text=True).stdout)
    assert len(rows) == 38, len(rows)
    printed = {r["name"]: r for r in rows}
    assert len(printed) == 38

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_249.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    weapons = [r for r in inventory["records"] if r["family"] == "weapons"]
    # The planning inventory lists every printed weapon except Heavy Crossbow.
    assert set(printed) - {r["name"] for r in weapons} == {"Heavy Crossbow"}
    assert not {r["name"] for r in weapons} - set(printed)
    expected = [(r["id"], r["name"], r["source_page"]) for r in weapons if r.get("catalogue_issue") == 249]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 249 and catalogue["batch"] == "data-weapons-03"
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected

    for f in families:
        name = f["name"]
        slug = name.lower().replace(" ", "-")
        payload, keys, cost = record(printed[name])
        assert f["id"] == f"srd-5.2.1:content/items/{slug}" and f["kind"] == "mundane" and f["category"] == "weapon", name
        assert f["base"] == {"options": [], "excluded": []} and f["presentations"][0]["pages"] == [PAGE, PAGE], name
        assert len(f["variants"]) == 1, name
        v = f["variants"][0]
        assert v["id"] == f["id"] + "/variant/standard" and v["name"] == name and v["rarity"] == "none", name
        assert v["attunement"] == {"mode": "none", "restrictions": []}, name
        assert same(v["charges"], {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}), name
        assert v["equip"] == {"mode": "wield", "restrictions": []} and v["use"] == {"mode": "weapon-attack", "restrictions": []}, name
        assert same(v["cost"], cost), (name, v["cost"], cost)
        assert v["required_mechanic_keys"] == keys, (name, keys)
        assert same(v["details"], {"weapon": payload}), (name, v["details"], payload)
    print(f"ok: {len(families)} weapons match their printed table rows, p. {PAGE}")


main()
