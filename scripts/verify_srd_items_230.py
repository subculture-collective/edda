#!/usr/bin/env python3
"""Compare issue #230 item identities, headers, tables and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (218, 219, 220)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor", "weapon": "Weapon", "wondrous-item": "Wondrous Item", "potion": "Potion"}
OPTION = {"scale-mail": "Scale Mail", "any-simple-weapon": "Any Simple", "any-martial-weapon": "Martial", "half-plate-armor": "Half Plate Armor",
          "plate-armor": "Plate Armor", "warhammer": "Warhammer", "chain-mail": "Chain Mail", "chain-shirt": "Chain Shirt"}
# The source's own attunement qualifier, printed after "Requires Attunement".
ATTUNEMENT_BY = {"Dwarven Thrower": " by a Dwarf or a Creature Attuned to a Belt of Dwarvenkind"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# The next printed heading that is not one of this batch's families ends the last block.
FOLLOWING = "Energy Bow Weapon (Longbow or Shortbow), Very Rare"


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return re.sub(r"(\w)- (\w)", r"\1\2", " ".join(FOOTER.sub(" ", " ".join(raw.split())).split()))


def header(family):
    """Rebuild the printed category and rarity line from the catalogued fields."""
    rarities = {v["rarity"] for v in family["variants"]}
    attunements = {v["attunement"]["mode"] for v in family["variants"]}
    assert len(rarities) == 1 and len(attunements) == 1, family["name"]
    options = family["base"]["options"]
    assert not family["base"]["excluded"]
    label = CATEGORY[family["category"]]
    if options:
        if options == ["any-simple-weapon", "any-martial-weapon"]:
            label += " (Any Simple or Martial)"
        else:
            label += " (" + " or ".join(OPTION[o] for o in options) + ")"
    line = f"{label}, {RARITY[rarities.pop()]}"
    if attunements.pop() == "required":
        line += " (Requires Attunement" + ATTUNEMENT_BY.get(family["name"], "") + ")"
    return line


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert phrase in source, f"{label} {value}: {phrase}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_230.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order, so Dragon Scale
    # Mail continues from p. 218 to p. 219 and Dust of Dryness across its column.
    stream = " ".join(clean(text(pdf, p, x)) for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_230.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 230]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 230
    assert len(families) == 12 and [len(f["variants"]) for f in families] == [10, 1, 1, 1, 1, 1, 1, 1, 1, 4, 1, 1]
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    names = [f["name"] for f in families]
    headers = {f["name"]: f"{f['name']} {header(f)} " for f in families}
    for printed in headers.values():
        assert stream.count(printed) == 1, printed
    assert stream.count(FOLLOWING) == 1
    starts = sorted((stream.index(p), n) for n, p in headers.items())
    assert [n for _, n in starts] == names, "printed order differs from catalogue order"
    block = {}
    for i, (at, name) in enumerate(starts):
        end = starts[i + 1][0] if i + 1 < len(starts) else stream.index(FOLLOWING)
        block[name] = stream[at + len(headers[name]): end]
    for f in families:
        page = f["presentations"][0]["pages"]
        assert page[0] == page[1] and page[0] in PAGES, f["name"]
        # The page is the one on which the heading is printed.
        pages = [p for p in PAGES for x in (0, 306) if headers[f["name"]] in clean(text(pdf, p, x))]
        assert pages == [page[0]], (f["name"], pages)
        for v in f["variants"]:
            assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}
            assert v["charges"] == {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
        assert " GP" not in block[f["name"]] and "charge" not in block[f["name"]].lower()
    variants = {f["name"]: [(v["id"].rsplit("/", 1)[1], v["name"], v["details"]) for v in f["variants"]] for f in families}

    # Dragon Scale Mail: ten variants from the dragon/resistance table.
    b = block["Dragon Scale Mail"]
    assert " Dragon Resistance Dragon Resistance " in b
    rows = re.findall(r"(Black|Blue|Brass|Bronze|Copper) (Acid|Lightning|Fire|Poison|Cold) (Gold|Green|Red|Silver|White) (Acid|Lightning|Fire|Poison|Cold)", b)
    assert len(rows) == 5
    pairs = sorted([(d.lower(), r.lower()) for row in rows for d, r in (row[0:2], row[2:4])])
    got = sorted((v[2]["scale_mail"]["dragon"], v[2]["scale_mail"]["resistance"]) for v in variants["Dragon Scale Mail"])
    assert pairs == got and len(set(p[0] for p in pairs)) == 10, pairs
    assert [v[0] for v in variants["Dragon Scale Mail"]] == sorted(v[0] for v in variants["Dragon Scale Mail"])
    for slug, name, detail in variants["Dragon Scale Mail"]:
        c = detail["scale_mail"]
        assert name == f"Dragon Scale Mail ({slug})" and c["dragon"] == slug
        tied([c["condition"]], {"while-wearing": "While wearing this armor"}, b, "scale mail condition")
        assert c["ac_bonus"] == 1 and "you gain a +1 bonus to Armor Class" in b
        tied([c["breath_weapon_save_mode"]], {"advantage": "you have Advantage on saving throws against the breath weapons of Dragons"}, b, "scale mail breath")
        assert c["breath_weapon_source"] == "dragons"
        assert "Resistance to one damage type determined by the kind of dragon that provided the scales" in b
        tied([c["sense_action"]], {"magic": "you can focus your senses as a Magic action"}, b, "scale mail sense action")
        tied([c["sense_target"]], {"closest-dragon-of-same-type": "the closest dragon within 30 miles of yourself that is of the same type as the armor"}, b, "scale mail target")
        assert c["sense_range_miles"] == 30
        tied(c["sense_reveals"], {"distance": "discern the distance", "direction": "and direction"}, b, "scale mail reveals")
        assert (c["sense_uses"], c["sense_recharge"]) == (1, "next-dawn") and "can’t be used again until the next dawn" in b

    c = variants["Dragon Slayer"][0][2]["slayer"]
    b = block["Dragon Slayer"]
    assert (c["attack_bonus"], c["damage_bonus"]) == (1, 1) and "+1 bonus to attack rolls and damage rolls made with this magic weapon" in b
    assert c["extra_damage"] == {"dice": 3, "sides": 6, "type": "weapon-type"} and "an extra 3d6 damage of the weapon’s type" in b
    tied([c["extra_damage_target"]], {"dragon": "if the target is a Dragon"}, b, "slayer target")

    c = variants["Dust of Disappearance"][0][2]["disappearance"]
    b = block["Dust of Disappearance"]
    tied([c["appearance"]], {"fine-sand": "This powder resembles fine sand"}, b, "disappearance appearance")
    assert c["uses"] == 1 and "enough of it for one use" in b
    tied([c["action"]], {"utilize": "When you take a Utilize action to throw the dust into the air"}, b, "disappearance action")
    assert c["area_emanation_feet"] == 10 and "within a 10-foot Emanation originating from you" in b
    tied(c["affects"], {"you": "you and", "each-creature": "each creature", "each-object": "and object"}, b, "disappearance affects")
    tied([c["grants_condition"]], {"invisible": "have the Invisible condition"}, b, "disappearance condition")
    assert c["duration"] == {"dice": 2, "sides": 4, "type": "minutes"} and "for 2d4 minutes" in b
    assert c["duration_same_for_all"] and "The duration is the same for all subjects" in b
    tied([c["consumed_when"]], {"magic-takes-effect": "the dust is consumed when its magic takes effect"}, b, "disappearance consumed")
    tied(c["ends_for_creature_after"], {"makes-attack-roll": "makes an attack roll", "deals-damage": "deals damage", "casts-spell": "or casts a spell"}, b, "disappearance end")

    c = variants["Dust of Dryness"][0][2]["dryness"]
    b = block["Dust of Dryness"]
    tied([c["container"]], {"small-packet": "This small packet contains"}, b, "dryness container")
    assert c["pinches"] == {"dice": 1, "sides": 6, "type": "pinches"} and c["pinches_bonus"] == 4 and "1d6 + 4 pinches of dust" in b
    tied([c["water_action"]], {"utilize": "As a Utilize action, you can sprinkle a pinch of the dust over water"}, b, "dryness water action")
    assert c["water_cube_feet"] == 15 and "up to a 15-foot Cube of water" in b
    tied([c["water_result"]], {"marble-sized-pellet": "into one marble-sized pellet"}, b, "dryness pellet")
    tied([c["pellet_location"]], {"floats-or-rests-near-where-sprinkled": "which floats or rests near where the dust was sprinkled"}, b, "dryness location")
    tied([c["pellet_weight"]], {"negligible": "The pellet’s weight is negligible"}, b, "dryness weight")
    tied([c["smash_action"]], {"utilize": "A creature can take a Utilize action to smash the pellet"}, b, "dryness smash action")
    tied([c["smash_surface"]], {"hard-surface": "against a hard surface"}, b, "dryness surface")
    tied([c["smash_result"]], {"pellet-shatters-and-releases-absorbed-water": "shatter and release the water the dust absorbed"}, b, "dryness result")
    tied([c["smash_ends"]], {"destroys-pellet-and-ends-its-magic": "Doing so destroys the pellet and ends its magic"}, b, "dryness end")
    tied([c["elemental_action"]], {"utilize": "As a Utilize action, you can sprinkle a pinch of the dust on an Elemental"}, b, "dryness elemental action")
    assert c["elemental_range_feet"] == 5 and "within 5 feet of yourself" in b, "continuation across the column break"
    tied([c["elemental_target"]], {"elemental-composed-mostly-of-water": "that is composed mostly of water (such as a Water Elemental)"}, b, "dryness target")
    assert (c["save_ability"], c["save_dc"], c["save_success"]) == ("constitution", 13, "half-damage") and "makes a DC 13 Constitution saving throw" in b
    assert c["damage"] == {"dice": 10, "sides": 6, "type": "necrotic"} and "taking 10d6 Necrotic damage on a failed save or half as much damage on a successful one" in b

    c = variants["Dust of Sneezing and Choking"][0][2]["sneezing"]
    b = block["Dust of Sneezing and Choking"]
    tied([c["container"]], {"small-container": "Found in a small container"}, b, "sneezing container")
    tied([c["resembles"]], {"dust-of-disappearance": "this powder resembles Dust of Disappearance"}, b, "sneezing resembles")
    assert c["identify_reveals"] and "Identify reveals it to be such" in b
    assert c["uses"] == 1 and "enough of it for one use" in b
    tied([c["action"]], {"utilize": "As a Utilize action, you can throw the dust into the air"}, b, "sneezing action")
    assert c["area_emanation_feet"] == 30 and c["affects_user"] and "yourself and every creature in a 30-foot Emanation originating from you" in b
    assert (c["save_ability"], c["save_dc"]) == ("constitution", 15) and "make a DC 15 Constitution saving throw" in b
    tied(c["auto_succeed_creature_types"], {"constructs": "Constructs", "elementals": "Elementals", "oozes": "Oozes", "plants": "Plants", "undead": "and Undead succeed on the save automatically"}, b, "sneezing types")
    tied([c["failure_effect"]], {"sneezes-uncontrollably": "begins sneezing uncontrollably"}, b, "sneezing failure")
    tied(c["failure_conditions"], {"incapacitated": "it has the Incapacitated condition", "suffocating": "and is suffocating"}, b, "sneezing conditions")
    tied([c["repeat_save"]], {"end-of-each-of-its-turns": "repeats the save at the end of each of its turns"}, b, "sneezing repeat")
    tied(c["ends_on"], {"successful-save": "ending the effect on itself on a success", "lesser-restoration": "targeted by a Lesser Restoration spell"}, b, "sneezing end")

    c = variants["Dwarven Plate"][0][2]["plate"]
    b = block["Dwarven Plate"]
    tied([c["condition"]], {"while-wearing": "While wearing this armor"}, b, "plate condition")
    assert c["ac_bonus"] == 2 and "you gain a +2 bonus to Armor Class" in b
    tied([c["reaction_trigger"]], {"effect-moves-you-against-your-will-along-ground": "if an effect moves you against your will along the ground"}, b, "plate trigger")
    tied([c["reaction_effect"]], {"reduce-distance-moved": "you can take a Reaction to reduce the distance you are moved"}, b, "plate effect")
    assert c["reduce_distance_up_to_feet"] == 10 and "by up to 10 feet" in b

    c = variants["Dwarven Thrower"][0][2]["thrower"]
    b = block["Dwarven Thrower"]
    tied(c["attunement_by_any_of"], {"dwarf": "", "creature-attuned-to-belt-of-dwarvenkind": ""}, "", "thrower attunement")
    assert (c["attack_bonus"], c["damage_bonus"]) == (3, 3) and "+3 bonus to attack rolls and damage rolls made with this magic weapon" in b
    assert c["thrown_property"] and (c["normal_range_feet"], c["long_range_feet"]) == (20, 60)
    assert "It has the Thrown property with a normal range of 20 feet and a long range of 60 feet" in b
    assert c["ranged_hit_extra_damage"] == {"dice": 1, "sides": 8, "type": "force"} and "When you hit with a ranged attack using this weapon, it deals an extra 1d8 Force damage" in b
    assert c["giant_target_extra_damage"] == {"dice": 2, "sides": 8, "type": "force"} and "or an extra 2d8 Force damage if the target is a Giant" in b
    tied([c["giant_target"]], {"giant": "if the target is a Giant"}, b, "thrower giant")
    tied([c["returns_to"], c["returns_when"]], {"your-hand": "flies back to your hand", "immediately-after-hitting-or-missing": "Immediately after hitting or missing"}, b, "thrower return")

    c = variants["Efficient Quiver"][0][2]["quiver"]
    b = block["Efficient Quiver"]
    assert c["extradimensional"] and "connects to an extradimensional space" in b and "three compartments" in b
    assert c["max_weight_pounds"] == 2 and "never weighing more than 2 pounds" in b
    comp = c["compartments"]
    assert [(x["length"], x["capacity"]) for x in comp] == [("shortest", 60), ("midsize", 18), ("longest", 6)]
    assert "The shortest compartment can hold up to 60 Arrows, Bolts, or similar objects" in b
    assert "The midsize compartment holds up to 18 Javelins or similar objects" in b
    assert "The longest compartment holds up to 6 long objects, such as bows, Quarterstaffs, or Spears" in b
    assert [x["holds"] for x in comp] == [["arrows", "bolts", "similar-objects"], ["javelins", "similar-objects"], ["long-objects", "bows", "quarterstaffs", "spears"]]
    tied([c["draw_as"]], {"regular-quiver-or-scabbard": "as if doing so from a regular quiver or scabbard"}, b, "quiver draw")

    c = variants["Efreeti Bottle"][0][2]["bottle"]
    b = block["Efreeti Bottle"]
    tied([c["appearance"]], {"painted-brass-bottle": "this painted brass bottle"}, b, "bottle appearance")
    tied([c["action"]], {"magic": "When you take a Magic action to remove the stopper"}, b, "bottle action")
    tied([c["smoke_ends"]], {"end-of-your-turn": "At the end of your turn"}, b, "bottle smoke end")
    tied([c["smoke_end_effect"]], {"flash-of-harmless-fire": "the smoke disappears with a flash of harmless fire"}, b, "bottle smoke effect")
    assert c["summons"] == "efreeti" and c["appears_within_feet"] == 30 and "an Efreeti appears in an unoccupied space within 30 feet of you" in b
    tied([c["table_applies_to"]], {"first-opening": "The first time the bottle is opened, the GM rolls on the following table"}, b, "bottle table")
    assert c["table_die_sides"] == 10 and " 1d10 Effect 1 The efreeti attacks you" in b
    r1, r2, r3 = c["effects"]
    assert [(r["min"], r["max"]) for r in c["effects"]] == [(1, 1), (2, 9), (10, 10)] and "2–9 The efreeti understands" in b and " 10 The efreeti understands" in b
    tied([r1["effect"]], {"efreeti-attacks-you": "The efreeti attacks you"}, b, "bottle row 1")
    assert r1["fight_rounds"] == 5 and "After fighting for 5 rounds, the efreeti disappears, and the bottle loses its magic" in b
    assert r1["ends"] == ["after-fighting-rounds"] and r1["bottle_after"] == "bottle-loses-magic" and not (r1["understands_owner_languages"] or r1["obeys_owner_commands"])
    assert r2["understands_owner_languages"] and r2["obeys_owner_commands"] and "understands your languages and obeys your commands for 1 hour" in b
    assert r2["duration_seconds"] == 3600 and r2["ends"] == ["duration-elapsed"]
    tied([r2["bottle_after"]], {"returns-to-bottle-with-new-stopper": "after which it returns to the bottle, and a new stopper contains it"}, b, "bottle row 2 end")
    assert r2["stopper_locked_hours"] == 24 and "The stopper can’t be removed for 24 hours" in b
    assert r2["same_effect_next_openings"] == 2 and "The next two times the bottle is opened, the same effect occurs" in b
    assert r2["escapes_on_opening"] == 4 and "If the bottle is opened a fourth time, the efreeti escapes and disappears, and the bottle loses its magic" in b
    tied([r2["effect"]], {"efreeti-serves-you": "obeys your commands"}, b, "bottle row 2")
    assert r3["understands_owner_languages"] and not r3["obeys_owner_commands"] and "understands your languages and can cast Wish once for you" in b
    assert r3["wishes"] == 1 and r3["duration_seconds"] == 3600 and r3["bottle_after"] == "bottle-loses-magic"
    tied(r3["ends"], {"grants-wish": "It disappears when it grants the wish", "duration-elapsed": "or after 1 hour, and the bottle loses its magic"}, b, "bottle row 10 end")
    tied([r3["effect"]], {"efreeti-casts-wish-once": "can cast Wish once for you"}, b, "bottle row 10")
    for r in c["effects"]:
        assert r["fight_rounds"] == (5 if r is r1 else 0) and r["wishes"] == (1 if r is r3 else 0)

    # Elemental Gem: four variants from the gem/elemental table.
    b = block["Elemental Gem"]
    assert " Gem Summoned Elemental Blue sapphire Air Elemental Emerald Water Elemental Red corundum Fire Elemental Yellow diamond Earth Elemental" in b
    rows = re.findall(r"(Blue sapphire|Emerald|Red corundum|Yellow diamond) (Air|Water|Fire|Earth) Elemental", b)
    assert [(g.lower().replace(" ", "-"), e.lower() + "-elemental") for g, e in rows] == [(v[2]["gem"]["gem"], v[2]["gem"]["elemental"]) for v in variants["Elemental Gem"]]
    for slug, name, detail in variants["Elemental Gem"]:
        c = detail["gem"]
        assert name == "Elemental Gem (" + slug.replace("-", " ") + ")" and c["gem"] == slug
        tied([c["action"]], {"utilize": "When you take a Utilize action to break the gem"}, b, "gem action")
        tied([c["gem_after"]], {"ceases-to-be-magical": "the gem ceases to be magical"}, b, "gem after")
        tied([c["location"]], {"unoccupied-space-nearest-broken-gem": "appears in an unoccupied space as close to the broken gem as possible"}, b, "gem location")
        assert c["understands_owner_languages"] and "understands your languages" in b and c["obeys_owner_commands"] and "obeys your commands" in b
        tied([c["initiative"]], {"immediately-after-summoner": "takes its turn immediately after you on your Initiative count"}, b, "gem initiative")
        assert c["duration_seconds"] == 3600
        tied(c["end_triggers"], {"duration-elapsed": "disappears after 1 hour", "dies": "when it dies", "dismissed-as-bonus-action": "when you dismiss it as a Bonus Action"}, b, "gem end")

    c = variants["Elixir of Health"][0][2]["elixir"]
    b = block["Elixir of Health"]
    assert c["cures_magical_contagions"] and "you are cured of all magical contagions" in b
    tied(c["ends_conditions"], {"blinded": "Blinded", "deafened": "Deafened", "paralyzed": "Paralyzed", "poisoned": "and Poisoned"}, b, "elixir conditions")
    assert "the following conditions end on you" in b
    tied([c["appearance"]], {"clear-red-liquid-with-tiny-bubbles-of-light": "The clear, red liquid has tiny bubbles of light in it"}, b, "elixir appearance")

    c = variants["Elven Chain"][0][2]["chain"]
    b = block["Elven Chain"]
    tied([c["condition"]], {"while-wearing": "while you wear this armor"}, b, "chain condition")
    assert c["ac_bonus"] == 1 and "You gain a +1 bonus to Armor Class" in b
    assert c["trained_with_this_armor"] and "You are considered trained with this armor even if you lack training with Medium or Heavy armor" in b
    tied(c["trained_even_without"], {"medium-armor": "Medium", "heavy-armor": "or Heavy armor"}, b, "chain training")
    print("#230 pinned PDF, twelve identities/pages/headers, 24 source-listed variants, tables and source values pass")


if __name__ == "__main__":
    main()
