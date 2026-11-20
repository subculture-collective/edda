#!/usr/bin/env python3
"""Compare issue #225 item identities, headers, tables and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (209, 210, 211)
CATEGORY = {"armor": "Armor", "weapon": "Weapon", "wondrous-item": "Wondrous Item"}
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
BASE = {
    ("any-medium-armor", "any-heavy-armor"): "Any Medium or Heavy",
    ("any-light-armor", "any-medium-armor", "any-heavy-armor"): "Any Light, Medium, or Heavy",
    ("any-ammunition",): "Any Ammunition",
    ("shield",): "Shield",
    ("plate-armor",): "Plate Armor",
}
EXCLUDED = {("hide-armor",): ", Except Hide Armor", (): ""}


def text(pdf, page, *crop):
    args = ["pdftotext", "-f", str(page), "-l", str(page)]
    if crop:
        args += ["-x", str(crop[0]), "-y", "0", "-W", "306", "-H", "792"]
    result = subprocess.run(args + ["-layout", str(pdf), "-"], check=True, capture_output=True, text=True)
    return result.stdout


def key(label):
    return re.sub(r"[^a-z0-9]+", "-", label.lower()).strip("-")


def header(family):
    """Rebuild the printed category and rarity line from the catalogued fields."""
    line = CATEGORY[family["category"]]
    base = family["base"]
    if base["options"]:
        line += " (" + BASE[tuple(base["options"])] + EXCLUDED[tuple(base["excluded"])] + ")"
    variants = family["variants"]
    if len(variants) == 1:
        rarity = RARITY[variants[0]["rarity"]]
    else:
        parts = [f"{RARITY[v['rarity']]} (+{v['id'][-1]})" for v in variants]
        rarity = ", ".join(parts[:-1]) + ", or " + parts[-1]
    attunement = {v["attunement"]["mode"] for v in variants}
    assert len(attunement) == 1
    return f"{line}, {rarity}" + (" (Requires Attunement)" if attunement == {"required"} else "")


def section(page, start, end):
    begin = page.index(start)
    return page[begin:page.index(end, begin + len(start))]


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_225.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave.
    pages = {p: " ".join(" ".join(text(pdf, p, x).split()) for x in (0, 306)) for p in PAGES}
    joined = {p: re.sub(r"(\w)- (\w)", r"\1\2", t) for p, t in pages.items()}

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_225.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 225]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 225
    assert len(families) == 12 and sum(len(f["variants"]) for f in families) == 16
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    for f in families:
        page = f["presentations"][0]["pages"]
        assert page[0] == page[1] and page[0] in PAGES, f["name"]
        # The heading is immediately followed by its category and rarity line.
        printed = f"{f['name']} {header(f)}"
        assert printed in pages[page[0]], printed
        rest = pages[page[0]].split(printed, 1)[1]
        assert not rest.startswith((" (Requires", ",", " or ")), printed
        for v in f["variants"]:
            assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}
            assert v["charges"] == {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
    # No GP price is printed for any of the twelve families.
    assert " GP" not in section(pages[209], "Adamantine Armor", "System Reference Document") + pages[210]
    assert " GP" not in section(pages[211], "Armor of Vulnerability", "Arrow-Catching Shield")
    details = {f["name"]: [v["details"] for v in f["variants"]] for f in families}

    # Ammunition of Slaying: every 1d100 row, then the save and damage.
    slaying = details["Ammunition of Slaying"][0]["slaying"]
    block = section(pages[209], "Ammunition of Slaying Weapon", "System Reference Document")
    rows = [(int(a), int(b) or 100, key(t)) for a, b, t in re.findall(r"(\d\d)–(\d\d) ([A-Z][a-z]+)", block)]
    assert sorted(rows) == [(r["min"], r["max"], r["creature_type"]) for r in slaying["creature_types"]] and len(rows) == 14
    damage = slaying["extra_damage"]
    assert f"DC {slaying['save_dc']} {slaying['save'].title()} saving throw" in joined[209]
    assert f"extra {damage['dice']}d{damage['sides']} {damage['type'].title()} damage" in joined[209]
    for snippet in ("GM chooses or determines randomly", "half as much extra damage on a successful one", "ammunition becomes nonmagical", "1d100 Creature Type"):
        assert snippet in joined[209], snippet

    # Ammunition and Armor, +1, +2, or +3: the bonus follows the rarity suffix.
    for name, field in (("Ammunition, +1, +2, or +3", "ammunition_bonus"), ("Armor, +1, +2, or +3", "armor_bonus")):
        for n, d in enumerate(details[name], 1):
            values = {v for k, v in d[field].items() if k.endswith("_bonus")}
            assert values == {n}, name
    ammo = details["Ammunition, +1, +2, or +3"][0]["ammunition_bonus"]
    assert ammo["typical_quantities"] == [10, 20] and ammo["value_pieces"] == 10
    for snippet in ("bonus to attack rolls and damage rolls", "The bonus is determined by the rarity of the ammunition",
                    "Once it hits a target, the ammunition is no longer magical", "quantities of ten or twenty pieces",
                    "Ten pieces of this ammunition are equivalent in value to a potion of the same rarity", "any Critical Hit against you becomes a normal hit"):
        assert snippet in joined[209], snippet
    for snippet in ("bonus to Armor Class while wearing this armor", "The bonus is determined by its rarity"):
        assert snippet in joined[210], snippet

    health = details["Amulet of Health"][0]["health"]
    assert f"Your {health['ability'].title()} is {health['score']} while you wear this amulet" in joined[209]
    assert f"{health['ability'].title()} is {health['no_effect_if_score_at_least']} or higher without it" in joined[209]
    assert "can’t be targeted by Divination spells or perceived through magical scrying sensors unless you allow it" in joined[209]

    # Amulet of the Planes: check, radius, every range and every plane in order.
    planes = details["Amulet of the Planes"][0]["planes"]
    block = section(joined[209], "Amulet of the Planes Wondrous Item", "Animated Shield")
    check = planes["check"]
    assert f"DC {check['dc']} {check['ability'].title()} ({check['skill'].title()}) check" in block
    assert f"within {planes['failure_radius_feet']} feet of you" in block and "you cast Plane Shift" in block and "take a Magic action" in block
    chunks = re.split(r"(\d\d)–(\d\d) ", block.split("1d100 Destination ")[1])[1:]
    assert len(chunks) == 15 and len(planes["destinations"]) == 5
    for row, a, b, body in zip(planes["destinations"], chunks[0::3], chunks[1::3], chunks[2::3]):
        assert (row["min"], row["max"]) == (int(a), int(b) or 100), body
        names = [key(n.removeprefix("the ")) for _, n in re.findall(r"on an? (\d), ([A-Za-z ]+?)(?:;|$)", body.strip())]
        assert names == [p["plane"] for p in row["planes"]] and [p["roll"] for p in row["planes"]] == list(range(1, len(names) + 1)), body
        die = re.search(r"rolling 1d(\d+)", body)
        assert row["sub_die_sides"] == (int(die.group(1)) if die else 0) == len(names)
        place = {"random-location-on-named-plane": "on the plane you named", "random-location-on-inner-plane": "on an Inner Plane",
                 "random-location-on-outer-plane": "on an Outer Plane", "random-location-on-astral-plane": "on the Astral Plane"}
        assert "Random location " + place[row["destination"]] in body

    shield = details["Animated Shield"][0]["animated"]
    assert shield["duration_seconds"] == 60 and shield["hands_free"]
    for snippet in ("While holding this Shield, you can take a Bonus Action to cause it to animate", "leaving your hands free",
                    "remains animate for 1 minute, until you take a Bonus Action to end this effect, or until you die or have the Incapacitated condition",
                    "falls to the ground or into your hand if you have one free"):
        assert snippet in joined[209], snippet

    # Apparatus of the Crab: statistics, then the two-column lever table.
    a = details["Apparatus of the Crab"][0]["apparatus"]
    left = joined[210].split("Armor, +1, +2, or +3")[0]
    stats = (f"AC {a['ac']}; HP {a['hp']}; Speed {a['speed_feet']} ft., Swim {a['swim_feet']} ft. "
             f"(or {a['speed_without_legs_feet']} ft. for both if the legs aren’t extended); Immunity to Poison and Psychic damage")
    pressure = a["pressure_damage"]
    for snippet in (stats, f"weighing {a['weight_pounds']} pounds", f"DC {a['catch_check']['dc']} Intelligence (Investigation) check",
                    "two Medium or smaller creatures", "Ten levers", "giant lobster", "is a Large object", "requires one pilot",
                    "airtight and watertight", f"enough air for {a['air_hours']} hours of breathing, divided by the number of breathing creatures inside",
                    "floats on water", f"depth of {a['max_depth_feet']} feet",
                    f"{pressure['dice']}d{pressure['sides']} {pressure['type'].title()} damage each minute", "take a Utilize action to move as many as two",
                    "goes back to its neutral position", "from left to right"):
        assert snippet in left, snippet
    assert (a["occupants"], a["pilots"], a["lever_count"], a["max_levers_per_action"], a["pressure_interval_seconds"]) == (2, 1, 10, 2, 60)
    table = text(pdf, 210).split("Apparatus of the Crab Levers\n")[1]
    assert re.match(r"\s*Lever Up\s+Down\n", table)
    numbers = [int(n) for n in re.findall(r"^ {1,12}(\d+) {2,}\S", table, re.M)]
    assert numbers == [lever["lever"] for lever in a["levers"]] == list(range(1, 11))
    flat = re.sub(r"(\w)- (\w)", r"\1\2", " ".join(table.split()))
    up_hit = a["claw_up_damage"]
    for snippet in ("Legs extend, allowing the apparatus to walk and swim", "Legs retract, reducing the apparatus’s Speed and Swim",
                    "Speed to 0 and making it unable to benefit from bonuses to speed", "Forward window shutter opens", "Forward window shutter closes",
                    "Side window shutters open (two per side)", "Side window shutters close (two per side)",
                    "Two claws extend from the front side of the", "The claws retract", "Each extended claw makes the following melee attack",
                    f"+{a['claw_attack_bonus']} to hit, reach {a['claw_reach_feet']} ft. Hit: {a['claw_up_average']} ({up_hit['dice']}d{up_hit['sides']}) {up_hit['type'].title()} damage",
                    "Hit: The target has the Grappled", f"condition (escape DC {a['claw_escape_dc']})", "walks or swims forward provided its", "walks or swims backward provided its",
                    f"turns {a['turn_degrees']} degrees counterclockwise", f"turns {a['turn_degrees']} degrees clockwise",
                    f"Bright Light in a {a['bright_light_feet']}-foot radius", f"Dim Light for an additional {a['dim_light_additional_feet']} feet", "The light turns off",
                    f"sinks up to {a['vertical_feet']} feet if it’s in liquid", f"rises up to {a['vertical_feet']} feet if it’s in liquid",
                    "The rear hatch unseals and opens", "The rear hatch closes and seals"):
        assert snippet in flat, snippet
    assert a["side_shutters_per_side"] == 2 and a["claws"] == 2

    inv = details["Armor of Invulnerability"][0]["invulnerability"]
    types = "Bludgeoning, Piercing, and Slashing damage"
    assert inv["resistances"] == inv["property_immunities"] == ["bludgeoning", "piercing", "slashing"]
    assert (inv["property_duration_seconds"], inv["property_uses"], inv["property_recharge"]) == (600, 1, "next-dawn")
    for snippet in (f"You have Resistance to {types} while you wear this armor", f"Metal Shell. You can take a Magic action to give yourself Immunity to {types} for 10 minutes or until you are no longer wearing the armor",
                    "Once this property is used, it can’t be used again until the next dawn"):
        assert snippet in joined[210], snippet

    # Armor of Resistance: every 1d10 row. The Down lever column follows the table.
    res = details["Armor of Resistance"][0]["resistance"]
    block = section(pages[210], "1d10 Damage Type 1d10 Damage Type", "Down").removeprefix("1d10 Damage Type 1d10 Damage Type")
    rows = sorted((int(n), key(t)) for n, t in re.findall(r"(\d+) ([A-Z][a-z]+)", block))
    assert rows == [(r["roll"], r["damage_type"]) for r in res["damage_types"]] and len(rows) == res["table_die_sides"] == 10
    assert "Resistance to one type of damage while you wear this armor. The GM chooses the type or determines it randomly by rolling" in joined[210]

    vul = details["Armor of Vulnerability"][0]["vulnerability"]
    assert vul["damage_types"] == ["bludgeoning", "piercing", "slashing"] and (vul["resistance_count"], vul["vulnerability_count"]) == (1, 2)
    assert vul["cursed"] and not vul["removing_armor_ends_curse"]
    for snippet in ("Resistance to one of the following damage types: Bludgeoning, Piercing, or Slashing. The GM chooses the type or determines it randomly",
                    "Curse. This armor is cursed", "revealed only when the Identify spell is cast on the armor or you attune to it",
                    "Attuning to the armor curses you until you are targeted by a Remove Curse spell or similar magic",
                    "removing the armor fails to end the curse", "Vulnerability to two of the three damage types associated with the armor (not the one to which it grants Resistance)"):
        assert snippet in joined[211], snippet
    print("#225 pinned PDF, twelve identities/pages/headers, sixteen variants, tables and source values pass")


if __name__ == "__main__":
    main()
