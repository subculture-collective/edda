#!/usr/bin/env python3
"""Compare issue #229 item identities, headers, tables and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (216, 217, 218)
CATEGORY = {"armor": "Armor", "weapon": "Weapon", "wondrous-item": "Wondrous Item"}
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary", "artifact": "Artifact"}
BASE = {
    ("dagger",): "Dagger",
    ("greatsword", "longsword", "rapier", "scimitar", "shortsword"): "Greatsword, Longsword, Rapier, Scimitar, or Shortsword",
    ("any-melee-weapon",): "Any Melee Weapon",
    ("any-light-armor", "any-medium-armor", "any-heavy-armor"): "Any Light, Medium, or Heavy",
}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
NEXT_HEADING = "Dragon Scale Mail Armor (Scale Mail), Very Rare (Requires Attunement)"


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return " ".join(re.sub(r"(\w)- (\w)", r"\1\2", FOOTER.sub(" ", " ".join(raw.split()))).split())


def key(label):
    return re.sub(r"[^a-z0-9]+", "-", label.lower()).strip("-")


def header(family):
    """Rebuild the printed category and rarity line from the catalogued fields."""
    (variant,) = family["variants"]
    assert variant["id"] == family["id"] + "/variant/standard" and variant["name"] == family["name"]
    options, line = tuple(family["base"]["options"]), CATEGORY[family["category"]]
    assert not family["base"]["excluded"] and bool(options) == (family["category"] != "wondrous-item"), family["name"]
    if options:
        line += f" ({BASE[options]})"
    line += f", {RARITY[variant['rarity']]}"
    return line + (" (Requires Attunement)" if variant["attunement"]["mode"] == "required" else "")


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert phrase in source, f"{label} {value}: {phrase}"


def dice_text(d):
    return f"{d['dice']}d{d['sides']}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_229.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order, so the Deck of
    # Illusions table continues from the end of p. 217 to the start of p. 218.
    columns = {(p, x): clean(text(pdf, p, x)) for p in PAGES for x in (0, 306)}
    stream = " ".join(columns[(p, x)] for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_229.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 229]
    planned = [b for b in inventory["batches"] if b["key"] == "data-magic-items-05"]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 229 and catalogue["batch"] == "data-magic-items-05"
    assert len(families) == 12 and all(len(f["variants"]) == 1 for f in families)
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert len(planned) == 1 and planned[0]["issue"] == 229 and planned[0]["pages"] == "216, 217, 218"
    assert [(m["name"], m["source_page"]) for m in planned[0]["members"]] == [(name, page) for _, name, page in expected]

    # Each family's section runs from its printed heading to the next family's heading.
    names = [f["name"] for f in families]
    starts = []
    block = {}
    for f in families:
        page = f["presentations"][0]["pages"]
        assert page[0] == page[1] and page[0] in PAGES, f["name"]
        printed = f"{f['name']} {header(f)} "
        assert stream.count(printed) == 1, printed
        starts.append(stream.index(printed) + len(printed))
        cost = f["variants"][0]["cost"]
        assert cost == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}
    assert starts == sorted(starts), "family headings must follow the catalogue order"
    assert stream.count(NEXT_HEADING) == 1
    ends = [stream.index(f"{n} {header(f)} ") for n, f in list(zip(names, families))[1:]] + [stream.index(NEXT_HEADING)]
    for f, begin, end in zip(families, starts, ends):
        block[f["name"]] = stream[begin:end]
        assert " GP" not in block[f["name"]], f["name"]
    # Printed page of each heading: the left or right column of its own page.
    for f in families:
        page = f["presentations"][0]["pages"][0]
        printed = f"{f['name']} {header(f)} "
        assert any(printed in columns[(page, x)] for x in (0, 306)), f"{f['name']} is not on p. {page}"
    detail = {f["name"]: f["variants"][0]["details"] for f in families}
    charges = {f["name"]: f["variants"][0]["charges"] for f in families}
    none = {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}

    c = detail["Crystal Ball of Telepathy"]["telepathy"]
    b = block["Crystal Ball of Telepathy"]
    tied([c["condition"]], {"while-touching": "While touching this crystal orb"}, b, "telepathy condition")
    tied([c["spell"]], {"scrying": "you can cast Scrying (save DC 17) with it"}, b, "telepathy spell")
    assert c["save_dc"] == 17
    tied([c["communication_mode"]], {"telepathy": "communicate telepathically"}, b, "telepathy mode")
    tied([c["communication_targets"]], {"creatures-you-can-see-within-range-of-sensor": "with creatures you can see within 30 feet of the spell’s sensor"}, b, "telepathy targets")
    assert c["communication_range_feet"] == 30
    tied([c["communication_range_from"]], {"spell-sensor": "of the spell’s sensor"}, b, "telepathy origin")
    tied([c["extra_spell"]], {"suggestion": "You can also cast Suggestion (save DC 17)"}, b, "telepathy extra spell")
    assert c["extra_save_dc"] == 17
    tied([c["extra_targets"]], {"one-of-those-creatures": "on one of those creatures"}, b, "telepathy extra target")
    tied([c["extra_delivery"]], {"through-the-sensor": "through the sensor"}, b, "telepathy delivery")
    assert not c["extra_concentration_required"] and "You don’t need to concentrate on this Suggestion to maintain it" in b
    tied([c["extra_ends_with"]], {"scrying-ends": "but it ends if Scrying ends"}, b, "telepathy end")
    assert (c["extra_uses"], c["extra_recharge"]) == (1, "next-dawn") and "You can’t cast Suggestion in this way again until the next dawn" in b
    assert charges["Crystal Ball of Telepathy"] == none and "charge" not in b.lower()

    c = detail["Crystal Ball of True Seeing"]["true_seeing"]
    b = block["Crystal Ball of True Seeing"]
    tied([c["condition"]], {"while-touching": "While touching this crystal orb"}, b, "true seeing condition")
    tied([c["spell"]], {"scrying": "you can cast Scrying (save DC 17) with it"}, b, "true seeing spell")
    assert c["save_dc"] == 17
    tied([c["sense"]], {"truesight": "you have Truesight"}, b, "true seeing sense")
    assert c["sense_range_feet"] == 120 and "with a range of 120 feet" in b
    tied([c["sense_centered_on"]], {"spell-sensor": "centered on the spell’s sensor"}, b, "true seeing centre")

    c = detail["Cube of Force"]["cube_of_force"]
    b = block["Cube of Force"]
    tied([c["appearance"]], {"cube": "This cube"}, b, "cube appearance")
    assert c["width_inches"] == 1 and "about an inch across" in b
    tied([c["markings"]], {"distinct-marking-on-each-face": "Each face has a distinct marking on it"}, b, "cube markings")
    tied([c["activation"]], {"press-face": "You can press one of those faces"}, b, "cube activation")
    tied([c["cost"]], {"expend-charges-required-for-face": "expend the number of charges required for it"}, b, "cube cost")
    assert c["save_dc"] == 17 and "(save DC 17), as shown in the Cube of Force Faces table" in b
    assert c["starting_charges"] == 10 and "The cube starts with 10 charges" in b
    assert (c["regain_dice"], c["regain_timing"]) == ({"dice": 1, "sides": 6, "type": "charges"}, "dawn") and "regains 1d6 expended charges daily at dawn" in b
    assert charges["Cube of Force"] == {"mode": "limited", "maximum": 10, "per_use": 0, "recharge": "regain-1d6-at-dawn"}
    table = b.split("Cube of Force Faces Spell Charge Cost ", 1)[1]
    printed = [(key(s), int(n)) for s, n in re.findall(r"([A-Z][A-Za-z ]+?) (\d+)(?= |$)", table)]
    assert [(r["spell"], r["charge_cost"]) for r in c["faces"]] == printed and len(printed) == 6, printed

    c = detail["Cubic Gate"]["cubic_gate"]
    b = block["Cubic Gate"]
    tied([c["appearance"]], {"cube": "This cube"}, b, "gate appearance")
    assert c["width_inches"] == 3 and "is 3 inches across" in b
    tied([c["radiates"]], {"palpable-magical-energy": "radiates palpable magical energy"}, b, "gate radiance")
    assert c["sides"] == 6 and "The six sides of the cube are each keyed to a different plane of existence" in b
    assert c["material_plane_sides"] == 1 and "one of which is the Material Plane" in b
    tied([c["other_planes_determined_by"]], {"gm": "linked to planes determined by the GM"}, b, "gate planes")
    assert c["max_charges"] == 3 and "The cube has 3 charges" in b
    assert (c["regain_dice"], c["regain_timing"]) == ({"dice": 1, "sides": 3, "type": "charges"}, "dawn") and "regains 1d3 expended charges daily at dawn" in b
    tied([c["action"]], {"magic": "As a Magic action"}, b, "gate action")
    assert c["charge_cost"] == 1 and "expend 1 of the cube’s charges" in b
    tied([c["gate_spell"]], {"gate": "Gate. Pressing one side of the cube, you cast Gate"}, b, "gate spell")
    tied([c["gate_activation"]], {"press-one-side": "Pressing one side of the cube, you cast Gate"}, b, "gate activation")
    tied([c["gate_effect"]], {"opens-portal-to-plane-keyed-to-side": "opening a portal to the plane of existence keyed to that side"}, b, "gate effect")
    tied([c["plane_shift_spell"]], {"plane-shift": "Plane Shift. Pressing one side of the cube twice, you cast Plane Shift"}, b, "plane shift spell")
    tied([c["plane_shift_activation"]], {"press-one-side-twice": "Pressing one side of the cube twice"}, b, "plane shift activation")
    tied([c["plane_shift_effect"]], {"transports-targets-to-plane-keyed-to-side": "transporting the targets to the plane of existence keyed to that side"}, b, "plane shift effect")
    assert charges["Cubic Gate"] == {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "regain-1d3-at-dawn"}

    c = detail["Dagger of Venom"]["venom"]
    b = block["Dagger of Venom"]
    assert (c["attack_bonus"], c["damage_bonus"]) == (1, 1) and "You gain a +1 bonus to attack rolls and damage rolls made with this magic weapon" in b
    tied([c["apply_action"]], {"bonus-action": "You can take a Bonus Action to magically coat"}, b, "venom action")
    tied([c["coats"]], {"blade": "coat the blade with poison"}, b, "venom coat")
    assert c["poison_duration_seconds"] == 60 and "The poison remains for 1 minute" in b
    tied(c["poison_ends_on"], {"duration-elapsed": "The poison remains for 1 minute", "attack-hits-creature": "or until an attack using this weapon hits a creature"}, b, "venom end")
    tied([c["save"]], {"constitution": "DC 15 Constitution saving throw"}, b, "venom save")
    assert c["save_dc"] == 15
    d = c["failure_damage"]
    assert (dice_text(d), d["type"]) == ("2d10", "poison") and "take 2d10 Poison damage" in b
    tied([c["failure_condition"]], {"poisoned": "have the Poisoned condition for 1 minute"}, b, "venom condition")
    assert c["condition_duration_seconds"] == 60
    assert (c["uses"], c["recharge"]) == (1, "next-dawn") and "The weapon can’t be used this way again until the next dawn" in b

    c = detail["Dancing Sword"]["dancing"]
    b = block["Dancing Sword"]
    tied([c["activate_action"]], {"bonus-action": "You can take a Bonus Action to toss this magic weapon into the air"}, b, "sword action")
    tied([c["activation"]], {"toss-into-air": "toss this magic weapon into the air"}, b, "sword activation")
    assert c["hover_fly_feet"] == 30 and "begins to hover, flies up to 30 feet" in b
    assert c["attack_range_feet"] == 5 and "attacks one creature of your choice within 5 feet of itself" in b
    tied([c["attack_roll"]], {"uses-your-attack-roll": "The weapon uses your attack roll"}, b, "sword attack roll")
    assert c["adds_ability_modifier_to_damage"] and "adds your ability modifier to damage rolls" in b
    tied([c["move_action"]], {"bonus-action": "While the weapon hovers, you can take a Bonus Action to cause it to fly"}, b, "sword move action")
    assert (c["move_feet"], c["move_within_feet_of_you"]) == (30, 30) and "fly up to 30 feet to another spot within 30 feet of you" in b
    assert c["attack_with_move"] and "As part of the same Bonus Action, you can cause the weapon to attack one creature within 5 feet of the weapon" in b
    assert c["attacks_before_return"] == 4 and "After the hovering weapon attacks for the fourth time" in b
    tied([c["returns_to"]], {"your-hand": "it flies back to you and tries to return to your hand"}, b, "sword return")
    tied([c["no_hand_free"]], {"falls-to-ground-in-your-space": "If you have no hand free, the weapon falls to the ground in your space"}, b, "sword no hand")
    tied([c["no_path"]], {"moves-as-close-as-possible-then-falls": "it moves as close to you as it can and then falls to the ground"}, b, "sword no path")
    tied(c["hover_ends_on"], {"you-grasp-it": "if you grasp it", "more-than-30-feet-away": "or are more than 30 feet away from it"}, b, "sword hover end")
    assert c["hover_end_distance_feet"] == 30

    c = detail["Decanter of Endless Water"]["decanter"]
    b = block["Decanter of Endless Water"]
    tied([c["appearance"]], {"stoppered-flask": "This stoppered flask"}, b, "decanter appearance")
    assert c["weight_pounds"] == 2 and "The decanter weighs 2 pounds" in b
    tied([c["action"]], {"magic": "You can take a Magic action to remove the stopper"}, b, "decanter action")
    tied([c["opens"]], {"remove-stopper": "remove the stopper"}, b, "decanter opens")
    tied(c["water_types"], {"fresh": "fresh water or salt water", "salt": "salt water (your choice)"}, b, "decanter water")
    tied([c["water_choice"]], {"your-choice": "(your choice)"}, b, "decanter choice")
    tied([c["stops_at"]], {"start-of-your-next-turn": "The water stops pouring out at the start of your next turn"}, b, "decanter stop")
    assert "issue one of three command words" in b and len(c["commands"]) == 3
    for r, word, gallons in zip(c["commands"], ("Splash", "Fountain", "Geyser"), (1, 5, 30)):
        assert r == {"word": word.lower(), "gallons": gallons}
        assert f"{word}. The decanter produces {gallons} gallon" in b, word
    tied([c["geyser_shape"]], {"line": "in a Line 30 feet long and 1 foot wide"}, b, "geyser shape")
    assert (c["geyser_length_feet"], c["geyser_width_feet"]) == (30, 1)
    tied([c["geyser_aim"]], {"one-direction-no-action-if-holding": "If you’re holding the decanter, you can aim the geyser in one direction (no action required)"}, b, "geyser aim")
    tied([c["geyser_targets"]], {"one-creature-in-line": "One creature of your choice in the Line"}, b, "geyser targets")
    tied([c["geyser_save"]], {"strength": "DC 13 Strength saving throw"}, b, "geyser save")
    d = c["geyser_damage"]
    assert c["geyser_save_dc"] == 13 and (dice_text(d), d["type"]) == ("1d4", "bludgeoning") and "take 1d4 Bludgeoning damage" in b
    tied([c["geyser_failure_condition"]], {"prone": "have the Prone condition"}, b, "geyser condition")
    tied([c["object_condition"]], {"not-worn-or-carried": "that isn’t being worn or carried"}, b, "geyser object")
    assert c["object_max_pounds"] == 200 and "weighs no more than 200 pounds" in b
    tied([c["object_effect"]], {"knocked-over": "The object is knocked over by the geyser"}, b, "geyser object effect")

    c = detail["Deck of Illusions"]["illusions"]
    b = block["Deck of Illusions"]
    assert (c["full_deck_cards"], c["creature_cards"], c["mirror_cards"]) == (34, 32, 2)
    assert "A full deck has 34 cards: 32 depicting specific creatures and two with a mirrored surface" in b
    d = c["missing_cards_dice"]
    assert (dice_text(d), d["type"], c["missing_cards_subtract"]) == ("1d20", "cards", 1) and "usually missing 1d20 − 1 cards" in b
    tied([c["functions_only_if"]], {"cards-drawn-at-random": "functions only if its cards are drawn at random"}, b, "deck random")
    tied([c["draw_action"]], {"magic": "You can take a Magic action to draw a card at random from the deck and throw"}, b, "deck draw")
    assert c["throw_range_feet"] == 30 and "to the ground at a point within 30 feet of yourself" in b
    tied([c["illusion"]], {"looks-and-behaves-like-real-creature": "looks and behaves like a real creature of its kind"}, b, "deck illusion")
    assert not c["can_do_harm"] and "except that it can do no harm" in b
    tied([c["move_action"]], {"magic": "you can take a Magic action to move it anywhere"}, b, "deck move")
    assert c["move_while_within_feet"] == 120 and "While you are within 120 feet of the illusory creature" in b
    tied([c["move_requires"]], {"can-see-illusion": "and can see it"}, b, "deck move sight")
    assert c["move_distance_feet"] == 30
    tied([c["move_from"]], {"its-card": "within 30 feet of its card"}, b, "deck move origin")
    tied([c["physical_interaction"]], {"reveals-false-objects-pass-through": "Any physical interaction with the illusory creature reveals it to be false, because objects pass through it"}, b, "deck touch")
    tied([c["inspect_action"]], {"study": "A creature that takes a Study action to visually inspect the illusory creature"}, b, "deck inspect")
    assert c["inspect_check"] == {"ability": "intelligence", "skill": "investigation", "dc": 15} and "successful DC 15 Intelligence (Investigation) check" in b
    tied([c["inspect_result"]], {"identifies-as-illusion": "identifies it as an illusion"}, b, "deck inspect result")
    tied(c["lasts_until"], {"card-moved": "The illusion lasts until its card is moved", "illusion-dispelled": "or the illusion is dispelled"}, b, "deck duration")
    tied(c["dispelled_by"], {"dispel-magic": "using a Dispel Magic spell", "similar-effect": "or a similar effect"}, b, "deck dispel")
    tied([c["after_end"]], {"card-image-disappears": "the image on its card disappears"}, b, "deck end")
    assert not c["card_reusable"] and "that card can’t be used again" in b
    tied([c["stat_block_source"]], {"monsters": "appear in “Monsters.”"}, b, "deck stat blocks")
    tied([c["stat_block_exception"]], {"card-drawer": "(except the card drawer)"}, b, "deck stat block exception")
    # The table starts on p. 217 and its header repeats at the top of p. 218.
    assert c["table_die_sides"] == 100 and c["table_page"] == 218
    assert "79–81 Oni" in columns[(217, 306)] and "82–84" not in columns[(217, 306)] and columns[(218, 0)].startswith("1d100 Illusion* 82–84 Priest")
    table = b.split("Deck of Illusions 1d100 Illusion* ", 1)[1].split(" *Stat blocks", 1)[0].replace("1d100 Illusion* ", "")
    parts = re.split(r"(\d\d)–(\d\d) ", table)
    assert parts[0] == "" and len(parts) == 1 + 3 * 33
    printed = [(int(lo), 100 if hi == "00" else int(hi), key(name.removeprefix("The "))) for lo, hi, name in zip(parts[1::3], parts[2::3], (p.strip() for p in parts[3::3]))]
    assert [(r["min"], r["max"], r["creature"]) for r in c["illusions"]] == printed
    assert printed[-1] == (97, 100, "card-drawer") and "97–00 The card drawer" in b and len(printed) == c["creature_cards"] + 1

    c = detail["Defender"]["defender"]
    b = block["Defender"]
    assert (c["attack_bonus"], c["damage_bonus"]) == (3, 3) and "You gain a +3 bonus to attack rolls and damage rolls made with this magic weapon" in b
    tied([c["transfer_timing"]], {"first-attack-each-turn": "The first time you attack with the weapon on each of your turns"}, b, "defender timing")
    tied(c["transfer_from"], {"attack-rolls": "reduce the bonus to your attack rolls", "damage-rolls": "and damage rolls"}, b, "defender source")
    tied([c["transfer_to"]], {"armor-class": "transfer some or all of the weapon’s bonus to your Armor Class"}, b, "defender target")
    tied([c["transfer_amount"]], {"some-or-all": "transfer some or all of"}, b, "defender amount")
    assert (c["example_attack_damage_bonus"], c["example_armor_class_bonus"]) == (1, 2) and "to +1 and gain a +2 bonus to Armor Class" in b
    tied([c["adjustment_lasts"]], {"until-start-of-your-next-turn": "remain in effect until the start of your next turn"}, b, "defender duration")
    assert c["must_hold_for_armor_class_bonus"] and "you must hold the weapon to gain a bonus to AC from it" in b

    c = detail["Demon Armor"]["demon_armor"]
    b = block["Demon Armor"]
    tied([c["condition"]], {"while-wearing": "While wearing this armor"}, b, "demon armor condition")
    assert c["armor_class_bonus"] == 1 and "you gain a +1 bonus to Armor Class" in b
    tied([c["language"]], {"abyssal": "you know Abyssal"}, b, "demon armor language")
    tied([c["gauntlets"]], {"clawed-gauntlets": "the armor’s clawed gauntlets"}, b, "demon armor gauntlets")
    d = c["unarmed_damage"]
    assert (dice_text(d), d["type"]) == ("1d8", "slashing") and "to deal 1d8 Slashing damage" in b
    tied([c["unarmed_replaces_damage_type"]], {"bludgeoning": "instead of the usual Bludgeoning damage"}, b, "demon armor replaced type")
    assert (c["unarmed_attack_bonus"], c["unarmed_damage_bonus"]) == (1, 1) and "you gain a +1 bonus to the attack and damage rolls of your Unarmed Strikes" in b
    assert c["cursed"] and "this cursed armor" in b
    tied([c["curse_trigger"]], {"don-armor": "Once you don this cursed armor"}, b, "demon armor curse")
    assert c["cannot_doff"] and "you can’t doff it" in b
    tied(c["doff_allowed_if"], {"targeted-by-remove-curse": "unless you are targeted by a Remove Curse spell", "targeted-by-similar-magic": "or similar magic"}, b, "demon armor doff")
    tied(c["curse_disadvantage"], {"attack-rolls-against-demons": "you have Disadvantage on attack rolls against demons",
                                   "saving-throws-against-demon-spells-and-special-abilities": "and on saving throws against their spells and special abilities"}, b, "demon armor disadvantage")

    c = detail["Dimensional Shackles"]["shackles"]
    b = block["Dimensional Shackles"]
    tied([c["place_action"]], {"utilize": "You can take a Utilize action to place these shackles on a creature"}, b, "shackles action")
    tied([c["place_target_condition"]], {"incapacitated": "that has the Incapacitated condition"}, b, "shackles target")
    tied([c["fit_size_min"], c["fit_size_max"]], {"small": "fit a creature of Small", "large": "to Large size"}, b, "shackles size")
    tied([c["prevents_movement"]], {"any-extradimensional-movement": "from using any method of extradimensional movement"}, b, "shackles prevent")
    tied(c["prevented_examples"], {"teleportation": "including teleportation", "travel-to-different-plane": "or travel to a different plane of existence"}, b, "shackles examples")
    tied([c["does_not_prevent"]], {"passing-through-interdimensional-portal": "They don’t prevent the creature from passing through an interdimensional portal"}, b, "shackles allow")
    tied([c["remove_action"]], {"utilize": "can take a Utilize action to remove them"}, b, "shackles remove")
    tied(c["removable_by"], {"you": "You and any creature", "creature-you-designate-when-you-use-the-shackles": "any creature you designate when you use the shackles"}, b, "shackles removers")
    assert c["escape_interval_days"] == 30 and "Once every 30 days, the bound creature can make a DC 30 Strength (Athletics) check" in b, "continuation across the column break"
    assert c["escape_check"] == {"ability": "strength", "skill": "athletics", "dc": 30}
    tied([c["escape_success"]], {"breaks-free-and-destroys-shackles": "the creature breaks free and destroys the shackles"}, b, "shackles escape")

    c = detail["Dragon Orb"]["dragon_orb"]
    b = block["Dragon Orb"]
    tied([c["appearance"]], {"etched-crystal-globe": "An orb is an etched crystal globe"}, b, "orb appearance")
    assert (c["diameter_inches"], c["used_diameter_inches"]) == (10, 20) and "about 10 inches in diameter" in b and "it grows to about 20 inches in diameter" in b
    tied([c["interior"]], {"swirling-mist": "mist swirls inside it"}, b, "orb interior")
    tied([c["peer_action"]], {"magic": "you can take a Magic action to peer into the orb’s depths"}, b, "orb action")
    tied([c["peer_requires"]], {"attuned": "While attuned to an orb"}, b, "orb requirement")
    tied([c["save"]], {"charisma": "DC 15 Charisma saving throw"}, b, "orb save")
    assert c["save_dc"] == 15
    tied([c["success_effect"]], {"control-orb-while-attuned": "you control the orb for as long as you remain attuned to it"}, b, "orb success")
    tied([c["failure_condition"]], {"charmed": "the orb imposes the Charmed condition on you"}, b, "orb failure")
    tied([c["failure_lasts"]], {"while-attuned": "for as long as you remain attuned to it. While you are Charmed"}, b, "orb failure duration")
    assert c["charmed_cannot_end_attunement"] and "you can’t voluntarily end your Attunement to it" in b
    tied([c["charmed_spell"]], {"suggestion": "the orb casts Suggestion on you"}, b, "orb charmed spell")
    tied([c["charmed_spell_frequency"]], {"at-will": "at will"}, b, "orb frequency")
    assert c["charmed_spell_save_dc"] == 18 and "Suggestion on you at will (save DC 18)" in b
    tied([c["charmed_urge"]], {"work-toward-evil-ends-orb-desires": "urging you to work toward the evil ends it desires"}, b, "orb urge")
    tied(c["desire_examples"], {"annihilation-of-a-particular-society-or-organization": "the annihilation of a particular society or organization",
                                "freedom-from-the-orb": "freedom from the orb", "spread-suffering-in-the-world": "to spread suffering in the world",
                                "advance-worship-of-tiamat": "to advance the worship of Tiamat", "something-else-the-gm-decides": "or something else the GM decides"}, b, "orb desires")
    assert c["max_charges"] == 7 and "The orb has 7 charges" in b
    d = c["regain_dice"]
    assert (dice_text(d), d["type"], c["regain_bonus"], c["regain_timing"]) == ("1d4", "charges", 3, "dawn") and "regains 1d4 + 3 expended charges daily at dawn" in b
    tied([c["spells_require"]], {"control-orb": "If you control the orb, you can cast one of the spells on the following table"}, b, "orb spells")
    assert charges["Dragon Orb"] == {"mode": "limited", "maximum": 7, "per_use": 0, "recharge": "regain-1d4-plus-3-at-dawn"}
    table = b.split("Spell Charge Cost ", 1)[1].split(" Call Dragons.", 1)[0]
    printed = []
    for label, cost in re.findall(r"([A-Z][A-Za-z ]*?(?: \([^)]*\))?) (\d+)(?= |$)", table):
        level = re.search(r"\(level (\d+) version\)", label)
        dc = re.search(r"\(save DC (\d+)\)", label)
        printed.append((key(re.sub(r"\(.*\)", "", label)), int(level.group(1)) if level else 0, int(dc.group(1)) if dc else 0, int(cost)))
    assert [(r["spell"], r["spell_level"], r["save_dc"], r["charge_cost"]) for r in c["spells"]] == printed and len(printed) == 5, printed
    tied([c["call_action"]], {"magic": "you can take a Magic action to cause the orb to issue a telepathic call"}, b, "orb call")
    tied([c["call_requires"]], {"control-orb": "Call Dragons. While you control the orb"}, b, "orb call requirement")
    assert c["call_range_miles"] == 40 and "extends in all directions for 40 miles" in b
    tied([c["call_directions"]], {"all-directions": "in all directions"}, b, "orb call directions")
    tied([c["call_targets"]], {"chromatic-dragons-in-range": "Chromatic dragons in range"}, b, "orb call targets")
    tied([c["call_effect"]], {"compelled-to-come-to-orb-by-most-direct-route": "feel compelled to come to the orb as soon as possible by the most direct route"}, b, "orb call effect")
    tied([c["call_unaffected"]], {"dragon-deities-such-as-tiamat": "Dragon deities such as Tiamat are unaffected by this call"}, b, "orb call exception")
    tied([c["call_attitude"]], {"might-be-hostile-toward-you": "might be Hostile toward you"}, b, "orb call attitude")
    assert c["call_cooldown_seconds"] == 3600 and "it can’t be used again for 1 hour" in b
    assert c["armor_class"] == 20 and "A Dragon Orb has AC 20" in b
    tied(c["destroyed_by"], {"damage-from-plus-3-weapon": "destroyed if it takes damage from a +3 Weapon", "disintegrate-spell": "or a Disintegrate spell"}, b, "orb destruction")
    assert c["nothing_else_harms"] and "Nothing else can harm it" in b

    for name in names:
        if name not in ("Cube of Force", "Cubic Gate", "Dragon Orb"):
            assert charges[name] == none, name
    print("#229 pinned PDF, twelve identities/pages/headers, one variant each, charges, tables and source values pass")


if __name__ == "__main__":
    main()
