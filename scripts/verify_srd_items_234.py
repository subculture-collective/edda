#!/usr/bin/env python3
"""Compare issue #234 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
# pp. 228 and 229 name the families; the Manual of Golems table is at the top of p. 230.
PAGES = (228, 229, 230)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# Printed category line of each weapon family, before the rarity.
CATEGORY = {
    "javelin-of-lightning": "Weapon (Javelin)",
    "luck-blade": "Weapon (Glaive, Greatsword, Longsword, Rapier, Scimitar, Sickle, or Shortsword)",
    "mace-of-disruption": "Weapon (Mace)",
    "mace-of-smiting": "Weapon (Mace)",
    "mace-of-terror": "Weapon (Mace)",
}
# The heading after the last family, printed below the Manual of Golems table.
FOLLOWING = "Manual of Quickness of Action Wondrous Item"


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return FOOTER.sub(" ", " ".join(raw.split()))


def flat(value):
    """Ignore spaces and hyphens so a word split across a line end still matches."""
    return re.sub(r"[\s\-]", "", value)


def has(source, phrase):
    return flat(phrase) in flat(source)


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert has(source, phrase), f"{label} {value}: {phrase}"


def header(family):
    variants = family["variants"]
    attune = variants[0]["attunement"]["mode"] == "required"
    suffix = " (Requires Attunement)" if attune else ""
    category = CATEGORY.get(family["id"].rsplit("/", 1)[1], "Wondrous Item")
    assert (family["category"] == "weapon") == category.startswith("Weapon")
    return f"{category}, {RARITY[variants[0]['rarity']]}{suffix}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_234.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages. Only the
    # left column of p. 230 is read: the Manual of Golems table is its first block.
    stream = " ".join(clean(text(pdf, p, x)) for p, x in ((228, 0), (228, 306), (229, 0), (229, 306), (230, 0)))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_234.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 234]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 234 and catalogue["batch"] == "data-magic-items-10"
    assert len(families) == 12 and all(len(f["variants"]) == 1 for f in families)
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["presentations"][0]["pages"][0] for f in families] == [228, 228, 228, 228] + [229] * 8
    names = [f["name"] for f in families]
    starts = []
    for f in families:
        printed = f"{f['name']} {header(f)}"
        assert stream.count(printed) == 1, printed
        starts.append(stream.index(printed))
    assert starts == sorted(starts), "printed order differs from the catalogue order"
    ends = starts[1:] + [stream.index(FOLLOWING)]
    block = {f["name"]: stream[s:e] for f, s, e in zip(families, starts, ends)}
    detail = {f["name"]: f["variants"][0]["details"] for f in families}
    charges = {f["name"]: f["variants"][0]["charges"] for f in families}
    none = {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
    for f in families:
        page = f["presentations"][0]["pages"]
        assert f["variants"][0]["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}, f["name"]
    # The golem table prints costs, but they are creation supplies, never a price.
    assert " GP" not in " ".join(b for n, b in block.items() if n != "Manual of Golems")
    assert all(c == none for n, c in charges.items() if n not in ("Luck Blade", "Mace of Terror"))

    d = detail["Iron Bands"]["iron_bands"]
    b = block["Iron Bands"]
    tied([d["appearance"]], {"rusty-iron-sphere": "This rusty iron sphere"}, b, "bands appearance")
    assert (d["diameter_inches"], d["weight_pounds"]) == (3, 1) and has(b, "measures 3 inches in diameter and weighs 1 pound")
    assert d["action"] == "magic" and d["throw_range_feet"] == 60 and has(b, "Magic action to throw the sphere at a Huge or smaller creature you can see within 60 feet of yourself")
    assert (d["target_max_size"], d["target_visibility"]) == ("huge", "can-see")
    tied([d["opens_into"]], {"tangle-of-metal-bands": "it opens into a tangle of metal bands"}, b, "bands opens")
    assert (d["attack_kind"], d["attack_ability"], d["attack_adds_proficiency_bonus"]) == ("ranged", "dexterity", True)
    assert has(b, "Make a ranged attack roll with an attack bonus equal to your Dexterity modifier plus your Proficiency Bonus")
    assert d["hit_condition"] == "restrained" and has(b, "the target has the Restrained condition")
    assert d["release_action"] == "bonus" and has(b, "until you take a Bonus Action to issue a command that releases it")
    tied(d["contracts_on"], {"release-command": "Doing so", "missed-attack": "missing with the attack"}, b, "bands contracts")
    assert has(b, "causes the bands to contract and become a sphere once more")
    assert d["break_action"] == "action" and has(b, "can take an action to make a DC 20 Strength (Athletics) check to break the iron bands")
    assert d["break_check"] == {"ability": "strength", "skill": "athletics", "dc": 20}
    tied([d["break_by"]], {"creature-that-can-touch-bands": "A creature that can touch the bands, including the one Restrained"}, b, "bands breaker")
    tied(d["break_success"], {"item-destroyed": "the item is destroyed", "restrained-creature-freed": "the Restrained creature is freed"}, b, "bands success")
    tied([d["break_failure"]], {"further-attempts-by-that-creature-fail": "any further attempts made by that creature automatically fail"}, b, "bands failure")
    assert d["break_failure_lockout_hours"] == 24 and has(b, "until 24 hours have elapsed")
    tied([d["reuse"]], {"next-dawn": "they can’t be used again until the next dawn"}, b, "bands reuse")

    d = detail["Iron Flask"]["iron_flask"]
    b = block["Iron Flask"]
    tied([d["condition"]], {"while-holding": "While holding this brass-stoppered iron flask"}, b, "flask condition")
    assert d["action"] == "magic" and d["range_feet"] == 60 and has(b, "Magic action to target a creature that you can see within 60 feet of yourself")
    assert d["target_visibility"] == "can-see" and d["requires_empty_flask"] is True and has(b, "If the flask is empty")
    tied([d["target_native_to"]], {"other-plane-than-yours": "native to a plane of existence other than the one you’re on"}, b, "flask native")
    assert (d["save_ability"], d["save_dc"]) == ("wisdom", 17) and has(b, "DC 17 Wisdom saving throw")
    tied([d["save_failure"]], {"trapped-in-flask": "be trapped in the flask"}, b, "flask failure")
    tied([d["previously_trapped_save"]], {"advantage": "it has Advantage on the save"}, b, "flask repeat")
    assert has(b, "If the target has been trapped by the flask before")
    tied([d["trapped_until"]], {"released": "remains in the flask until released"}, b, "flask until")
    assert d["capacity_creatures"] == 1 and has(b, "can hold only one creature at a time")
    tied(d["trapped_effects"], {"does-not-age": "doesn’t age", "does-not-need-to-breathe": "doesn’t need to breathe", "does-not-need-to-eat": "breathe, eat", "does-not-need-to-drink": "eat, or drink"}, b, "flask trapped")
    assert d["release_action"] == "magic" and has(b, "take a Magic action to remove the flask’s stopper and release the creature")
    tied([d["release_by"]], {"remove-stopper": "remove the flask’s stopper"}, b, "flask release")
    assert d["obeys_for_hours"] == 1 and has(b, "obeys your commands for 1 hour")
    assert d["understands_commands_in_any_language"] is True and has(b, "understanding those commands even if it doesn’t know the language in which the commands are given")
    tied(d["self_defense_when"], {"no-commands-issued": "If you issue no commands", "command-likely-to-cause-death-or-imprisonment": "a command that is likely to result in its death or imprisonment"}, b, "flask defence")
    tied([d["self_defense_behavior"]], {"defends-itself-takes-no-other-actions": "it defends itself but otherwise takes no actions"}, b, "flask behaviour")
    tied([d["after_duration"]], {"acts-per-normal-disposition-and-alignment": "acts in accordance with its normal disposition and alignment"}, b, "flask after")
    tied([d["identify_reveals"]], {"whether-flask-contains-creature": "An Identify spell reveals if the flask contains a creature"}, b, "flask identify")
    tied([d["type_determined_by"]], {"opening-flask": "the only way to determine the type of creature is to open the flask"}, b, "flask type")
    tied([d["newly_found"]], {"might-contain-creature-chosen-by-gm": "A newly discovered Iron Flask might already contain a creature chosen by the GM"}, b, "flask found")

    d = detail["Javelin of Lightning"]["javelin_of_lightning"]
    b = block["Javelin of Lightning"]
    assert (d["damage_swap_trigger"], d["damage_swap_from"], d["damage_swap_to"]) == ("attack-roll-hit", "piercing", "lightning")
    assert has(b, "Each time you make an attack roll with this magic weapon and hit, you can have it deal Lightning damage instead of Piercing damage")
    k = d["lightning_bolt"]
    tied([k["trigger"]], {"throw-at-target-within-120-feet": "When you throw this weapon at a target no farther than 120 feet from you"}, b, "bolt trigger")
    assert k["range_feet"] == 120 and k["forgo_attack_roll"] is True and has(b, "you can forgo making a ranged attack roll and instead turn the weapon into a bolt of lightning")
    assert (k["shape"], k["width_feet"]) == ("line", 5) and has(b, "forms a 5-foot-wide Line between you and the target")
    tied([k["line_between"]], {"you-and-target": "between you and the target"}, b, "bolt line")
    tied(k["affects"], {"target": "The target", "each-other-creature-in-line-excluding-you": "each other creature in the Line (excluding you)"}, b, "bolt affects")
    assert (k["save_ability"], k["save_dc"]) == ("dexterity", 13) and has(b, "makes a DC 13 Dexterity saving throw")
    assert k["damage"] == {"dice": 4, "sides": 6, "type": "lightning"} and has(b, "taking 4d6 Lightning damage on a failed save or half as much damage on a successful one")
    tied([k["on_success"]], {"half-damage": "half as much damage on a successful one"}, b, "bolt success")
    tied([k["returns"]], {"weapon-reappears-in-hand-after-damage": "Immediately after dealing this damage, the weapon reappears in your hand"}, b, "bolt returns")
    tied([k["reuse"]], {"next-dawn": "This property can’t be used again until the next dawn"}, b, "bolt reuse")

    d = detail["Lantern of Revealing"]["lantern_of_revealing"]
    b = block["Lantern of Revealing"]
    tied([d["condition"]], {"while-lit": "While lit, this hooded lantern"}, b, "lantern condition")
    assert (d["burn_hours"], d["oil_pints"]) == (6, 1) and has(b, "burns for 6 hours on 1 pint of oil")
    assert (d["bright_light_feet"], d["dim_light_additional_feet"]) == (30, 30) and has(b, "Bright Light in a 30-foot radius and Dim Light for an additional 30 feet")
    tied(d["reveals"], {"invisible-creatures": "Invisible creatures", "invisible-objects": "and objects are visible"}, b, "lantern reveals")
    tied([d["reveals_when"]], {"in-bright-light": "as long as they are in the lantern’s Bright Light"}, b, "lantern when")
    assert d["hood_action"] == "utilize" and has(b, "You can take a Utilize action to lower the hood")
    tied([d["hooded_light"]], {"dim-light": "reducing the lantern’s light to Dim Light in a 5-foot radius"}, b, "lantern hood")
    assert d["hooded_radius_feet"] == 5

    d = detail["Luck Blade"]["luck_blade"]
    b = block["Luck Blade"]
    assert (d["attack_bonus"], d["damage_bonus"]) == (1, 1) and has(b, "You gain a +1 bonus to attack rolls and damage rolls made with this magic weapon")
    assert d["save_bonus"] == 1 and has(b, "you also gain a +1 bonus to saving throws")
    tied([d["save_bonus_when"]], {"weapon-on-your-person": "While the weapon is on your person"}, b, "luck save when")
    k = d["luck"]
    tied([k["condition"]], {"weapon-on-your-person": "If the weapon is on your person"}, b, "luck condition")
    assert k["action"] == "none" and has(b, "call on its luck (no action required)")
    tied([k["effect"]], {"reroll-one-failed-d20-test-must-use-second-roll": "reroll one failed D20 Test"}, b, "luck effect")
    tied([k["forbidden_while"]], {"incapacitated": "if you don’t have the Incapacitated condition"}, b, "luck forbidden")
    tied([k["reroll_use"]], {"must-use-second-roll": "You must use the second roll"}, b, "luck reroll")
    tied([k["reuse"]], {"next-dawn": "Once used, this property can’t be used again until the next dawn"}, b, "luck reuse")
    k = d["wish"]
    assert k["initial_charges"] == {"dice": 1, "sides": 3, "type": "charges"} and has(b, "The weapon has 1d3 charges")
    tied([k["condition"]], {"while-holding": "While holding it"}, b, "wish condition")
    assert (k["charge_cost"], k["spell"]) == (1, "wish") and has(b, "you can expend 1 charge and cast Wish from it")
    tied([k["reuse"]], {"next-dawn": "Once used, this property can’t be used again until the next dawn"}, b, "wish reuse")
    tied([k["property_lost_when"]], {"no-charges": "The weapon loses this property if it has no charges"}, b, "wish lost")
    # The source prints no maximum or recharge: 3 is the top of the 1d3 starting roll.
    assert charges["Luck Blade"] == {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "none"} and not has(b, "regains")

    d = detail["Mace of Disruption"]["mace_of_disruption"]
    b = block["Mace of Disruption"]
    assert d["extra_damage"] == {"dice": 2, "sides": 6, "type": "radiant"} and has(b, "that creature takes an extra 2d6 Radiant damage")
    tied(d["extra_damage_targets"], {"fiend": "When you hit a Fiend", "undead": "or an Undead with this magic weapon"}, b, "disruption targets")
    assert d["destroy_at_hit_points_or_fewer"] == 25 and has(b, "If the target has 25 Hit Points or fewer after taking this damage")
    assert (d["save_ability"], d["save_dc"]) == ("wisdom", 15) and has(b, "must succeed on a DC 15 Wisdom saving throw or be destroyed")
    tied([d["save_failure"]], {"destroyed": "or be destroyed"}, b, "disruption failure")
    tied([d["save_success_condition"]], {"frightened": "the creature has the Frightened condition"}, b, "disruption success")
    tied([d["success_condition_ends"]], {"end-of-your-next-turn": "until the end of your next turn"}, b, "disruption ends")
    tied([d["light_condition"]], {"while-holding": "While you hold this weapon"}, b, "disruption light")
    assert (d["bright_light_feet"], d["dim_light_additional_feet"]) == (20, 20) and has(b, "Bright Light in a 20-foot radius and Dim Light for an additional 20 feet")

    d = detail["Mace of Smiting"]["mace_of_smiting"]
    b = block["Mace of Smiting"]
    assert (d["attack_bonus"], d["damage_bonus"]) == (1, 1) and has(b, "You gain a +1 bonus to attack rolls and damage rolls made with this magic weapon")
    assert d["construct_bonus"] == 3 and has(b, "The bonus increases to +3 when you use the weapon to attack a Construct")
    tied([d["critical_trigger"]], {"natural-20-on-attack-roll": "When you roll a 20 on an attack roll made with this weapon"}, b, "smiting trigger")
    assert (d["critical_extra_damage"], d["critical_damage_type"]) == (7, "bludgeoning") and has(b, "the target takes an extra 7 Bludgeoning damage")
    assert d["construct_critical_extra_damage"] == 14 and has(b, "or 14 Bludgeoning damage if it’s a Construct")
    assert d["construct_destroyed_at_hit_points_or_fewer"] == 25 and has(b, "If a Construct has 25 Hit Points or fewer after taking this damage, it is destroyed")

    d = detail["Mace of Terror"]["mace_of_terror"]
    b = block["Mace of Terror"]
    assert charges["Mace of Terror"] == {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "regain-1d3-at-dawn"} and d["max_charges"] == 3
    assert has(b, "This magic weapon has 3 charges and regains 1d3 expended charges daily at dawn")
    assert d["regain_expended_charges"] == {"dice": 1, "sides": 3, "type": "charges"}
    tied([d["regain_when"]], {"daily-at-dawn": "daily at dawn"}, b, "terror regain")
    tied([d["condition"]], {"while-holding": "While holding the weapon"}, b, "terror condition")
    assert d["action"] == "magic" and d["charge_cost"] == 1 and has(b, "you can take a Magic action and expend 1 charge to release a wave of terror from it")
    assert d["range_feet"] == 30 and has(b, "within 30 feet of you")
    tied([d["targets"]], {"each-creature-of-your-choice": "Each creature of your choice"}, b, "terror targets")
    assert (d["save_ability"], d["save_dc"]) == ("wisdom", 15) and has(b, "must succeed on a DC 15 Wisdom saving throw or have the Frightened condition for 1 minute")
    assert d["fail_condition"] == "frightened" and d["duration_seconds"] == 60
    tied(d["frightened_behavior"], {"must-move-as-far-from-you-as-possible": "must spend its turns trying to move as far away from you as it can", "cannot-make-opportunity-attacks": "it can’t make Opportunity Attacks", "action-only-dash-or-escape-movement-preventing-effect": "it can use only the Dash action or try to escape from an effect that prevents it from moving"}, b, "terror behaviour")
    tied([d["no_movement_action"]], {"dodge": "the creature can take the Dodge action"}, b, "terror dodge")
    assert has(b, "If it has nowhere it can move")
    tied([d["repeat_save"]], {"end-of-each-of-its-turns": "At the end of each of its turns, a creature repeats the save"}, b, "terror repeat")
    tied([d["success_ends"]], {"effect-on-itself": "ending the effect on itself on a success"}, b, "terror ends")

    d = detail["Mantle of Spell Resistance"]["mantle_of_spell_resistance"]
    b = block["Mantle of Spell Resistance"]
    tied([d["condition"]], {"while-wearing": "while you wear this cloak"}, b, "mantle condition")
    assert (d["roll"], d["mode"]) == ("saving-throw", "advantage") and has(b, "You have Advantage on saving throws against spells")
    tied([d["against"]], {"spells": "against spells"}, b, "mantle against")

    for name, ability, intro in (("Manual of Bodily Health", "constitution", "This book contains health and nutrition tips, and its words are charged with magic"), ("Manual of Gainful Exercise", "strength", "This book describes fitness exercises, and its words are charged with magic")):
        d = detail[name]["manual_of_" + name.split(" of ")[1].lower().replace(" ", "_")]
        b = block[name]
        assert has(b, intro) and has(b, "If you spend 48 hours over a period of 6 days or fewer studying the book’s contents and practicing its guidelines")
        assert (d["study_hours"], d["within_days"]) == (48, 6)
        assert d["ability"] == ability and has(b, f"your {ability.capitalize()} increases by 2, to a maximum of 30")
        assert (d["increase"], d["maximum"]) == (2, 30)
        assert d["then_loses_magic"] is True and d["regains_magic_after_years"] == 100 and has(b, "The manual then loses its magic but regains it in a century")

    d = detail["Manual of Golems"]["manual_of_golems"]
    b = block["Manual of Golems"]
    assert d["source_pages"] == [229, 230]
    tied([d["type_chosen_by"]], {"gm-or-1d20-table": "The GM chooses the type or determines it randomly by rolling on the accompanying table"}, b, "golems chosen")
    assert d["table_die_sides"] == 20 and "1d20 Golem Time Cost" in b
    table = re.findall(r"(\d+)(?:–(\d+))? (Clay|Flesh|Iron|Stone) Golem (\d+) days ([\d,]+) GP", b)
    assert len(table) == 4 == len(d["golems"]), table
    for (lo, hi, kind, days, cost), row in zip(table, d["golems"]):
        assert row == {"golem": kind.lower() + "-golem", "table_min": int(lo), "table_max": int(hi or lo), "time_days": int(days), "cost_gp": int(cost.replace(",", ""))}, row
    assert d["decipher_requires_spellcaster"] is True and (d["decipher_spell_slots"], d["decipher_slot_level"]) == (2, 5)
    assert has(b, "you must be a spellcaster with at least two level 5 spell slots")
    assert d["unable_user_damage"] == {"dice": 6, "sides": 6, "type": "psychic"} and has(b, "takes 6d6 Psychic damage")
    assert d["work_uninterrupted"] is True and d["max_rest_hours_per_day"] == 8 and has(b, "working without interruption with the manual at hand and resting no more than 8 hours per day")
    tied([d["supplies_cost"]], {"pay-table-cost": "You must also pay the specified cost to purchase supplies"}, b, "golems supplies")
    tied([d["on_completion"]], {"book-consumed-in-eldritch-flames": "the book is consumed in eldritch flames"}, b, "golems completion")
    tied([d["animation"]], {"ashes-sprinkled-on-golem": "when the ashes of the manual are sprinkled on it"}, b, "golems animation")
    tied([d["stat_block_in"]], {"monsters": "See “Monsters” for the golem’s stat block"}, b, "golems stat block")
    tied([d["golem_control"]], {"under-your-control-understands-and-obeys": "The golem is under your control, and it understands and obeys your commands"}, b, "golems control")

    # Every payload is covered above, and every family was matched by name.
    assert list(detail) == names
    print(f"issue #234: {len(families)} families match pinned SRD 5.2.1 pp. 228-230")


if __name__ == "__main__":
    main()
