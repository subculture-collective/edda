#!/usr/bin/env python3
"""Compare issue #227 item identities, headers, tables and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (214, 215)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare"}
# Printed quantities the catalogue stores in seconds.
SECONDS = {600: "10 minutes", 3600: "1 hour", 14400: "4 hours", 60: "1 minute"}


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
    assert family["category"] == "wondrous-item" and family["base"] == {"options": [], "excluded": []}, family["name"]
    (variant,) = family["variants"]
    assert variant["id"] == family["id"] + "/variant/standard" and variant["name"] == family["name"]
    line = f"Wondrous Item, {RARITY[variant['rarity']]}"
    return line + (" (Requires Attunement)" if variant["attunement"]["mode"] == "required" else "")


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert phrase in source, f"{label} {value}: {phrase}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_227.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave, then
    # join both pages: Candle of Invocation continues from p. 214 onto p. 215.
    pages = {p: " ".join(" ".join(text(pdf, p, x).split()) for x in (0, 306)) for p in PAGES}
    for p in PAGES:
        assert pages[p].count(f" {p} System Reference Document 5.2.1") == 1, p
    stream = " ".join(pages[p].replace(f" {p} System Reference Document 5.2.1", "") for p in PAGES)
    joined = re.sub(r"(\w)- (\w)", r"\1\2", stream)

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_227.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 227]
    (planned,) = [b for b in inventory["batches"] if b["key"] == "data-magic-items-03"]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 227 and catalogue["batch"] == planned["key"]
    assert len(families) == 12 and sum(len(f["variants"]) for f in families) == 12
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [(m["name"], m["source_page"]) for m in planned["members"]] == [(n, p) for _, n, p in expected]
    assert planned["pages"] == ", ".join(str(p) for p in PAGES)

    # Each heading is immediately followed by its category and rarity line, on
    # its inventory page. Blocks run from one heading to the next printed one.
    blocks, starts = {}, []
    for f in families:
        page = f["presentations"][0]["pages"]
        assert page[0] == page[1] and page[0] in PAGES, f["name"]
        printed = f"{f['name']} {header(f)}"
        assert pages[page[0]].count(printed) == 1 and joined.count(printed) == 1, printed
        rest = joined.split(printed, 1)[1]
        assert not rest.startswith((" (Requires", ",", " or ")), printed
        starts.append((joined.index(printed), printed, f["name"]))
        (v,) = f["variants"]
        assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}
        assert v["charges"] == {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
    # The twelve families are consecutive in the source; the next heading is not ours.
    end = joined.index("Censer of Controlling Air Elementals Wondrous Item")
    assert [s[2] for s in sorted(starts)] == sorted(f["name"] for f in families) and max(starts)[0] < end
    assert joined.index("Boots of Speed Wondrous Item") == 0, "p. 214 must open with the first family"
    ordered = sorted(starts) + [(end, "", "")]
    for (begin, printed, name), nxt in zip(ordered, ordered[1:]):
        blocks[name] = joined[begin + len(printed):nxt[0]]
    # No price or charge pool is printed for any of the twelve families.
    ours = joined[:end]
    assert " GP" not in ours and "charge" not in ours.lower()
    details = {f["name"]: next(iter(f["variants"][0]["details"].values())) for f in families}
    assert [next(iter(f["variants"][0]["details"])) for f in families] == [
        "boots_speed", "boots_striding", "boots_winterlands", "bowl", "bracers_archery", "bracers_defense",
        "brazier", "brooch", "broom", "candle", "cape", "carpet"]

    d, block = details["Boots of Speed"], blocks["Boots of Speed"]
    tied([d["condition"]], {"while-wearing": "While you wear these boots"}, block, "speed condition")
    tied([d["activate_action"]], {"bonus-action": "you can take a Bonus Action"}, block, "speed action")
    tied([d["activation"]], {"click-heels-together": "to click the boots’ heels together"}, block, "speed activation")
    tied([d["speed_multiplier"]], {2: "the boots double your Speed"}, block, "speed multiplier")
    tied([d["opportunity_attacks_against_wearer"]], {"disadvantage": "any creature that makes an Opportunity Attack against you has Disadvantage on the attack roll"}, block, "speed attacks")
    tied([d["end_trigger"]], {"click-heels-together-again": "If you click your heels together again, you end the effect"}, block, "speed end")
    assert f"used the boots’ property for a total of {SECONDS[d['total_use_seconds']]}" in block
    tied([d["after_total_use"]], {"ceases-to-function-for-wearer": "the magic ceases to function for you"}, block, "speed limit")
    tied([d["recharge"]], {"long-rest": "until you finish a Long Rest"}, block, "speed recharge")

    d, block = details["Boots of Striding and Springing"], blocks["Boots of Striding and Springing"]
    tied([d["condition"]], {"while-wearing": "While you wear these boots"}, block, "striding condition")
    assert f"your Speed becomes {d['speed_becomes_feet']} feet unless your Speed is higher" in block and d["speed_unless"] == "speed-is-higher"
    tied(d["speed_not_reduced_by"], {"carrying-weight-over-carrying-capacity": "your Speed isn’t reduced by you carrying weight in excess of your carrying capacity",
                                     "wearing-heavy-armor": "or wearing Heavy Armor"}, block, "striding reductions")
    tied([d["jump_frequency"]], {"once-on-each-of-your-turns": "Once on each of your turns"}, block, "striding frequency")
    assert f"you can jump up to {d['jump_distance_feet']} feet by spending only {d['jump_movement_cost_feet']} feet of movement" in block

    d, block = details["Boots of the Winterlands"], blocks["Boots of the Winterlands"]
    tied([d["appearance"]], {"furred-boots": "These furred boots are snug and feel warm"}, block, "winterlands appearance")
    tied([d["condition"]], {"while-wearing": "While wearing them, you gain the following benefits"}, block, "winterlands condition")
    cold, strider = d["cold_resistance"], d["winter_strider"]
    assert [key(n) for n in re.findall(r"([A-Z][a-z]+ [A-Z][a-z]+)\. You", block)] == [cold["benefit"], strider["benefit"]]
    assert f"You have Resistance to {cold['resistance'].title()} damage" in block
    assert f"can tolerate temperatures of {cold['tolerated_temperature_fahrenheit']} degrees Fahrenheit or lower without any additional protection" in block
    assert (cold["temperature_bound"], cold["without"]) == ("or-lower", "additional-protection")
    assert strider["ignores"] == "difficult-terrain" and "You ignore Difficult Terrain created by " + " or ".join(strider["created_by"]) + "." in block

    # Bowl and Brazier print the same summoning property around their own vessel.
    for name, vessel, payload, opening in (
            ("Bowl of Commanding Water Elementals", "bowl", "water-elemental", "While this bowl is filled with water and you are within {} feet of it"),
            ("Brazier of Commanding Fire Elementals", "brazier", "fire-elemental", "While you are within {} feet of this brazier")):
        d, block = details[name], blocks[name]
        s = d["summon"]
        assert opening.format(s["within_feet"]) in block, name
        tied([s["action"]], {"magic": "you can take a Magic action"}, block, name + " action")
        assert s["summons"] == payload and f"to summon a {payload.replace('-', ' ').title()}." in block
        tied([s["appears"]], {f"unoccupied-space-as-close-to-{vessel}-as-possible": f"appears in an unoccupied space as close to the {vessel} as possible"}, block, name + " appears")
        tied(s["elemental_traits"], {"understands-your-languages": "understands your languages", "obeys-your-commands": "obeys your commands",
                                     "acts-immediately-after-you-on-your-initiative-count": "takes its turn immediately after you on your Initiative count"}, block, name + " traits")
        tied(s["ends_when"], {"duration-elapsed": f"disappears after {SECONDS[s['duration_seconds']]}", "elemental-dies": "when it dies",
                              "dismissed-as-bonus-action": "or when you dismiss it as a Bonus Action"}, block, name + " ends")
        assert s["dismiss_action"] == "bonus-action" and s["uses"] == 1
        tied([s["recharge"]], {"next-dawn": f"The {vessel} can’t be used this way again until the next dawn"}, block, name + " recharge")
    d, block = details["Bowl of Commanding Water Elementals"], blocks["Bowl of Commanding Water Elementals"]
    assert set(d) == {"vessel_condition", "summon", "dimensions_approximate", "diameter_feet", "depth", "capacity_gallons"}
    assert d["vessel_condition"] == "filled-with-water" and d["dimensions_approximate"] is True and d["depth"] == "half-of-diameter"
    assert f"The bowl is about {d['diameter_feet']} foot in diameter and half as deep. It holds about {d['capacity_gallons']} gallons." in block
    assert set(details["Brazier of Commanding Fire Elementals"]) == {"summon"}

    d, block = details["Bracers of Archery"], blocks["Bracers of Archery"]
    tied([d["condition"]], {"while-wearing": "While wearing these bracers"}, block, "archery condition")
    weapons = " and ".join(w.title() for w in d["weapon_proficiencies"])
    assert f"you have proficiency with the {weapons}" in block and d["damage_bonus_weapons"] == d["weapon_proficiencies"]
    assert f"you gain a +{d['damage_bonus']} bonus to damage rolls made with such weapons" in block

    d, block = details["Bracers of Defense"], blocks["Bracers of Defense"]
    tied([d["condition"]], {"while-wearing": "While wearing these bracers"}, block, "defense condition")
    assert f"you gain a +{d['ac_bonus']} bonus to Armor Class" in block
    tied(d["requires"], {"wearing-no-armor": "if you are wearing no armor", "using-no-shield": "and using no Shield"}, block, "defense requirements")

    d, block = details["Brooch of Shielding"], blocks["Brooch of Shielding"]
    tied([d["condition"]], {"while-wearing": "While wearing this brooch"}, block, "brooch condition")
    assert [r.title() for r in d["resistances"]] == re.findall(r"you have Resistance to (\w+) damage", block)
    tied([d["immunity_to_damage_from"]], {"magic-missile-spell": "you have Immunity to damage from the Magic Missile spell"}, block, "brooch immunity")

    d, block = details["Broom of Flying"], blocks["Broom of Flying"]
    tied([d["appearance"]], {"wooden-broom": "This wooden broom"}, block, "broom appearance")
    tied([d["inactive"]], {"functions-like-mundane-broom": "functions like a mundane broom until"}, block, "broom inactive")
    tied([d["hover_requirement"]], {"stand-astride": "you stand astride it"}, block, "broom requirement")
    tied([d["hover_action"]], {"magic": "take a Magic action to make it hover beneath you"}, block, "broom hover action")
    assert f"It has a Fly Speed of {d['fly_speed_feet']} feet. It can carry up to {d['capacity_pounds']} pounds" in block
    assert f"its Fly Speed becomes {d['reduced_fly_speed_feet']} feet while carrying over {d['reduced_when_over_pounds']} pounds" in block
    tied(d["hover_ends"], {"rider-lands": "The broom stops hovering when you land", "no-longer-riding": "or when you’re no longer riding it"}, block, "broom hover ends")
    tied([d["send_action"]], {"magic": "As a Magic action, you can send the broom"}, block, "broom send action")
    tied([d["send_mode"]], {"travels-alone": "to travel alone to a destination"}, block, "broom send mode")
    assert f"a destination within {d['send_range_miles']} mile of you" in block
    tied(d["send_requirements"], {"name-the-location": "if you name the location", "familiar-with-location": "and are familiar with it"}, block, "broom send requirements")
    tied([d["recall_action"]], {"magic": "The broom comes back to you when you take a Magic action"}, block, "broom recall action")
    tied(d["recall_requirements"], {"command-word": "and use a command word",
                                    "broom-still-within-range": f"if the broom is still within {d['recall_range_miles']} mile of you"}, block, "broom recall requirements")

    # Candle of Invocation: the text starts on p. 214; all sixteen rows of the
    # two-column 1d100 table are on p. 215. A printed 00 is 100.
    d, block = details["Candle of Invocation"], blocks["Candle of Invocation"]
    tied([d["activation"]], {"lit": "This candle’s magic is activated when the candle is lit"}, block, "candle activation")
    tied([d["activate_action"]], {"magic": "which requires a Magic action"}, block, "candle action")
    assert f"After burning for {SECONDS[d['burn_seconds']]}, the candle is destroyed" in block and d["after_burn"] == "destroyed"
    assert d["snuff_early"] is True and "You can snuff it out early for use at a later time" in block
    assert f"Deduct the time it burned in increments of {SECONDS[d['burn_deduction_increment_seconds']]} from its total burn time" in block
    assert f"While lit, the candle sheds {d['light'].title()} Light in a {d['light_radius_feet']}-foot radius" in block
    tied([d["you_in_light"]], {"advantage-on-d20-tests": "While you are within that light, you have Advantage on D20 Tests"}, block, "candle advantage")
    casters = " or ".join(c.title() for c in d["free_casters"])
    assert f"a {casters} in the light can cast level {d['free_cast_spell_level']} spells they have prepared without expending spell slots" in block
    assert (d["free_cast_requirement"], d["free_cast_cost"]) == ("prepared", "no-spell-slot")
    tied([d["gate_condition"]], {"lit-for-first-time": "when you light the candle for the first time"}, block, "candle gate condition")
    assert f"you can cast {d['gate_casts'].title()} with it. Doing so destroys the candle" in block and d["gate_destroys_candle"] is True
    tied([d["gate_links_to"]], {"outer-plane": "links to a particular Outer Plane"}, block, "candle link")
    tied(d["plane_selection"], {"gm-chooses": "chosen by the GM", "random-table": "or determined by rolling on the following table"}, block, "candle selection")
    table = block.split("rolling on the following table.", 1)[1]
    assert table.count(f"1d{d['table_die_sides']} Outer Plane") == 2
    rows = [(int(a), int(b) or 100, key(n)) for a, b, n in re.findall(r"(\d\d)–(\d\d) ((?:[A-Z][a-z]+ ?)+)", table)]
    assert sorted(rows) == [(r["min"], r["max"], r["plane"]) for r in d["outer_planes"]] and len(rows) == 16
    assert d["table_page"] == 215 and "1d100 Outer Plane" in pages[215] and "1d100 Outer Plane" not in pages[214]
    assert "This candle’s magic is activated" in pages[214] and "Ysgard" in pages[215]

    d, block = details["Cape of the Mountebank"], blocks["Cape of the Mountebank"]
    tied([d["scent"]], {"brimstone": "This cape smells faintly of brimstone"}, block, "cape scent")
    tied([d["condition"]], {"while-wearing": "While wearing it"}, block, "cape condition")
    assert f"you can use it to cast {d['casts'].replace('-', ' ').title()} as a Magic action" in block and d["action"] == "magic" and d["uses"] == 1
    tied([d["recharge"]], {"next-dawn": "This property can’t be used again until the next dawn"}, block, "cape recharge")
    tied([d["on_teleport"]], {"leaves-cloud-of-smoke": "When you teleport with that spell, you leave behind a cloud of smoke"}, block, "cape smoke")
    tied([d["vacated_space"]], {"lightly-obscured": "The space you left is Lightly Obscured by that smoke"}, block, "cape space")
    tied([d["obscured_until"]], {"end-of-your-next-turn": "until the end of your next turn"}, block, "cape until")

    # Carpet of Flying: the four sizes are rows of the 1d100 table, in order.
    d, block = details["Carpet of Flying"], blocks["Carpet of Flying"]
    tied([d["effect"]], {"hovers-and-flies": "You can make this carpet hover and fly"}, block, "carpet effect")
    tied([d["activate_action"]], {"magic": "by taking a Magic action"}, block, "carpet action")
    tied([d["activation_requirement"]], {"command-word": "and using the carpet’s command word"}, block, "carpet command")
    tied([d["control"]], {"moves-according-to-your-directions": "It moves according to your directions"}, block, "carpet control")
    assert f"if you are within {d['control_range_feet']} feet of it" in block
    assert d["size_count"] == 4 and "Four sizes of Carpet of Flying exist" in block
    tied(d["size_selection"], {"gm-chooses": "The GM chooses the size of a given carpet", "random-table": "or determines it randomly by rolling on the following table"}, block, "carpet selection")
    assert d["max_load_multiplier"] == 2 and "A carpet can carry up to twice the weight shown on the table" in block
    tied([d["overloaded_fly_speed"]], {"halved": "but its Fly Speed is halved"}, block, "carpet overload")
    tied([d["overloaded_when"]], {"carrying-more-than-normal-capacity": "if it carries more than its normal capacity"}, block, "carpet overload condition")
    assert f"1d{d['table_die_sides']} Size Capacity Fly Speed" in block
    rows = re.findall(r"(\d\d)–(\d\d) (\d) ft\. × (\d) ft\. (\d+) lb\. (\d+) feet", block)
    assert [tuple(int(x) for x in r[:1]) + (int(r[1]) or 100,) + tuple(int(x) for x in r[2:]) for r in rows] == [
        (r["min"], r["max"], r["width_feet"], r["length_feet"], r["capacity_pounds"], r["fly_speed_feet"]) for r in d["sizes"]] and len(rows) == 4
    print("#227 pinned PDF, twelve identities/pages/headers, twelve variants, tables and source values pass")


if __name__ == "__main__":
    main()
