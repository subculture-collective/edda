#!/usr/bin/env python3
"""Compare issue #236 item identities, headers and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
# pp. 233-235 hold the twelve families: the Oil of Slipperiness continues from p. 233 to
# p. 234 and the Portable Hole from p. 234 to p. 235.
PAGES = (233, 234, 235)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor", "wondrous-item": "Wondrous Item", "potion": "Potion"}
OPTION = {"half-plate-armor": "Half Plate Armor", "plate-armor": "Plate Armor"}
# The source's own attunement qualifier, printed after "Requires Attunement".
ATTUNEMENT_BY = {"Pearl of Power": " by a Spellcaster"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# The next printed heading that is not one of this batch's families ends the last block.
FOLLOWING = "Potion of Clairvoyance Potion, Rare"


def text(pdf, page, x):
    args = ["pdftotext", "-f", str(page), "-l", str(page), "-x", str(x), "-y", "0", "-W", "306", "-H", "792", "-layout", str(pdf), "-"]
    return subprocess.run(args, check=True, capture_output=True, text=True).stdout


def clean(raw):
    return re.sub(r"(\w)- (\w)", r"\1\2", " ".join(FOOTER.sub(" ", " ".join(raw.split())).split()))


def header(family):
    """Rebuild the printed category and rarity line from the catalogued fields."""
    variants = family["variants"]
    assert len(variants) == 1, family["name"]
    options = family["base"]["options"]
    assert not family["base"]["excluded"]
    label = CATEGORY[family["category"]]
    if options:
        label += " (" + " or ".join(OPTION[o] for o in options) + ")"
    line = f"{label}, {RARITY[variants[0]['rarity']]}"
    if variants[0]["attunement"]["mode"] == "required":
        line += " (Requires Attunement" + ATTUNEMENT_BY.get(family["name"], "") + ")"
    return line


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert phrase in source, f"{label} {value}: {phrase}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_236.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream = " ".join(clean(text(pdf, p, x)) for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_236.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 236]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 236 and catalogue["batch"] == "data-magic-items-12"
    assert len(families) == 12 and [len(f["variants"]) for f in families] == [1] * 12
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["presentations"][0]["pages"][0] for f in families] == [233, 233, 234, 234, 234, 234, 234, 234, 234, 234, 234, 235]
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
        v = f["variants"][0]
        assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}
        assert " GP" not in block[f["name"]]
    detail = {f["name"]: f["variants"][0]["details"] for f in families}
    charges = {f["name"]: f["variants"][0]["charges"] for f in families}
    none = {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}
    three = {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "regain-1d3-at-dawn"}
    for name in names:
        assert charges[name] == (three if name.startswith("Pipes of ") else none), name
        assert ("charge" in block[name].lower()) == name.startswith("Pipes of "), name

    c = detail["Oil of Sharpness"]["oil_of_sharpness"]
    b = block["Oil of Sharpness"]
    assert (c["coats_melee_weapons"], c["coats_ammunition"]) == (1, 20) and "One vial of this oil can coat one Melee weapon or twenty pieces of ammunition" in b
    assert c["requires_nonmagical"] and "ammunition and Melee weapons that are nonmagical and deal Slashing or Piercing damage are affected" in b
    tied(c["damage_types"], {"slashing": "deal Slashing or Piercing damage", "piercing": "Slashing or Piercing damage"}, b, "sharpness damage")
    assert c["application_minutes"] == 1 and "Applying the oil takes 1 minute" in b
    assert (c["weapon_bonus"], c["ammunition_bonus"]) == (3, 3) and "turning the coated weapon into a +3 Weapon or the coated ammunition into +3 Ammunition" in b
    tied([c["appearance"]], {"clear-gelatinous-oil-with-ultrathin-silver-shards": "This clear, gelatinous oil sparkles with tiny, ultrathin silver shards"}, b, "sharpness appearance")

    c = detail["Oil of Slipperiness"]["oil_of_slipperiness"]
    b = block["Oil of Slipperiness"]
    assert c["source_pages"] == [233, 234] and "Grease spell in that area for 8 hours" in b
    assert c["coverage"] == {"max_size": "medium", "includes_worn_and_carried_equipment": True, "extra_vial_per_size_above": "medium"}
    assert "One vial of this oil can cover one Medium or smaller creature, along with the equipment it’s wearing and carrying (one additional vial is required for each size category above Medium)" in b
    assert c["application_minutes"] == 10 and "Applying the oil takes 10 minutes" in b
    tied([c["body_spell"]], {"freedom-of-movement": "gains the effect of the Freedom of Movement spell for 8 hours"}, b, "slipperiness body")
    assert c["body_duration_hours"] == 8
    tied([c["pour_action"]], {"magic": "poured on the ground as a Magic action"}, b, "slipperiness pour")
    assert c["pour_area_square_feet"] == 10 and "where it covers a 10-foot square" in b
    tied([c["pour_spell"]], {"grease": "duplicating the effect of the Grease spell"}, b, "slipperiness grease")
    assert c["pour_duration_hours"] == 8
    tied([c["appearance"]], {"sticky-black-unguent-thick-and-heavy-flows-quickly-when-poured": "This sticky, black unguent is thick and heavy, but it flows quickly when poured"}, b, "slipperiness appearance")

    c = detail["Pearl of Power"]["pearl_of_power"]
    b = block["Pearl of Power"]
    tied(c["attunement_by_any_of"], {"spellcaster": "(Requires Attunement by a Spellcaster)"}, headers["Pearl of Power"], "pearl attunement")
    tied([c["condition"]], {"on-your-person": "While this pearl is on your person"}, b, "pearl condition")
    tied([c["action"]], {"magic": "you can take a Magic action to regain"}, b, "pearl action")
    assert c["spell_slots_regained"] == 1 and c["max_slot_level"] == 3 and "regain one expended spell slot of level 3 or lower" in b
    assert (c["uses"], c["recharge"]) == (1, "next-dawn") and "Once you use the pearl, it can’t be used again until the next dawn" in b

    c = detail["Periapt of Health"]["periapt_of_health"]
    b = block["Periapt of Health"]
    tied([c["condition"]], {"while-wearing": "While wearing this pendant"}, b, "health condition")
    tied([c["action"]], {"magic": "you can take a Magic action to regain 2d4 + 2 Hit Points"}, b, "health action")
    assert c["healing_dice"] == {"dice": 2, "sides": 4, "type": "healing"} and c["healing_bonus"] == 2
    assert (c["uses"], c["recharge"]) == (1, "next-dawn") and "Once used, this property can’t be used again until the next dawn" in b
    tied([c["poisoned_save_mode"]], {"advantage": "you have Advantage on saving throws"}, b, "health save mode")
    tied(c["poisoned_save_applies_to"], {"avoid-poisoned": "to avoid or end the Poisoned condition", "end-poisoned": "or end the Poisoned condition"}, b, "health save scope")

    c = detail["Periapt of Proof against Poison"]["periapt_of_proof_against_poison"]
    b = block["Periapt of Proof against Poison"]
    tied([c["appearance"]], {"delicate-silver-chain-with-brilliant-cut-black-gem-pendant": "This delicate silver chain has a brilliant-cut black gem pendant"}, b, "proof appearance")
    tied([c["condition"]], {"while-wearing": "While you wear it"}, b, "proof condition")
    tied(c["condition_immunities"], {"poisoned": "Immunity to the Poisoned condition"}, b, "proof conditions")
    tied(c["damage_immunities"], {"poison": "and Poison damage"}, b, "proof damage")

    c = detail["Periapt of Wound Closure"]["periapt_of_wound_closure"]
    b = block["Periapt of Wound Closure"]
    tied([c["condition"]], {"while-wearing": "While wearing this pendant, you gain the following benefits"}, b, "closure condition")
    p = c["life_preservation"]
    tied([p["trigger"]], {"make-death-saving-throw": "Whenever you make a Death Saving Throw"}, b, "closure trigger")
    assert (p["changes_roll_at_or_below"], p["changed_to"]) == (9, 10) and "you can change a roll of 9 or lower to a 10" in b
    tied([p["result"]], {"failed-save-becomes-successful": "turning a failed save into a successful one"}, b, "closure result")
    h = c["natural_healing_boost"]
    tied([h["trigger"]], {"roll-hit-point-die-to-regain-hit-points": "Whenever you roll a Hit Point Die to regain Hit Points"}, b, "boost trigger")
    assert h["multiplier"] == 2 and "double the number of Hit Points it restores" in b

    c = detail["Philter of Love"]["philter_of_love"]
    b = block["Philter of Love"]
    tied([c["trigger"]], {"next-creature-you-see-after-drinking": "The next time you see a creature within 10 minutes after drinking this philter"}, b, "philter trigger")
    assert c["window_minutes"] == 10
    assert (c["condition"], c["charmed_by"], c["duration_hours"]) == ("charmed", "that-creature", 1)
    assert "you are charmed by that creature and have the Charmed condition for 1 hour" in b
    tied([c["appearance"]], {"rose-hued-effervescent-liquid-with-one-heart-shaped-bubble": "This rose-hued, effervescent liquid contains one easy-to-miss bubble shaped like a heart"}, b, "philter appearance")

    c = detail["Pipes of Haunting"]["pipes_of_haunting"]
    b = block["Pipes of Haunting"]
    assert "These pipes have 3 charges and regain 1d3 expended charges daily at dawn" in b
    assert c["action"] == "magic" and c["charge_cost"] == 1 and "You can take a Magic action to play them and expend 1 charge" in b
    tied([c["effect"]], {"eerie-spellbinding-tune": "to create an eerie, spellbinding tune"}, b, "haunting effect")
    tied([c["targets"]], {"each-creature-of-your-choice": "Each creature of your choice within 30 feet of you"}, b, "haunting targets")
    assert c["range_feet"] == 30 and c["save"] == {"ability": "wisdom", "dc": 15}
    assert "must succeed on a DC 15 Wisdom saving throw" in b
    assert (c["condition"], c["duration_seconds"]) == ("frightened", 60) and "have the Frightened condition for 1 minute" in b
    tied([c["repeat_save"]], {"end-of-each-of-its-turns": "A creature that fails the save repeats it at the end of each of its turns"}, b, "haunting repeat")
    tied([c["success_ends"]], {"effect-on-itself": "ending the effect on itself on a success"}, b, "haunting success")
    assert c["immunity_hours_after_success"] == 24 and "A creature that succeeds on its save is immune to the effect of these pipes for 24 hours" in b

    c = detail["Pipes of the Sewers"]["pipes_of_the_sewers"]
    b = block["Pipes of the Sewers"]
    p = c["passive"]
    tied([p["condition"]], {"on-your-person": "While these pipes are on your person"}, b, "sewers condition")
    tied(p["creatures"], {"ordinary-rat": "ordinary rats", "giant-rat": "and giant rats"}, b, "sewers creatures")
    tied([p["attitude"]], {"indifferent": "are Indifferent toward you"}, b, "sewers attitude")
    tied([p["unless"]], {"you-threaten-or-harm-them": "won’t attack you unless you threaten or harm them"}, b, "sewers unless")
    assert "The pipes have 3 charges and regain 1d3 expended charges daily at dawn" in b
    k = c["call_swarms"]
    assert (k["play_action"], k["call_action"]) == ("magic", "bonus")
    assert "If you play the pipes as a Magic action, you can take a Bonus Action to expend 1 to 3 charges" in b and (k["min_charges"], k["max_charges"]) == (1, 3)
    assert k["swarm"] == "swarm-of-rats" and k["swarms_per_charge"] == 1 and "calling forth one Swarm of Rats with each expended charge" in b
    tied([k["range"]], {"half-a-mile": "if enough rats are within half a mile of you"}, b, "sewers range")
    tied([k["insufficient_rats"]], {"charge-is-wasted": "If there aren’t enough rats to form a swarm, the charge is wasted"}, b, "sewers shortfall")
    tied([k["availability_determined_by"]], {"gm": "(as determined by the GM)"}, b, "sewers gm")
    tied([k["swarm_movement"]], {"toward-music-by-shortest-available-route-not-otherwise-controlled": "Called swarms move toward the music by the shortest available route but aren’t under your control otherwise"}, b, "sewers movement")
    s = c["sway_swarms"]
    tied([s["trigger"]], {"uncontrolled-swarm-of-rats-within-30-feet-while-you-play": "Whenever a Swarm of Rats that isn’t under another creature’s control comes within 30 feet of you while you are playing the pipes"}, b, "sway trigger")
    assert s["trigger_range_feet"] == 30 and s["save"] == {"ability": "wisdom", "dc": 15} and "the swarm makes a DC 15 Wisdom saving throw" in b
    tied([s["success_result"]], {"behaves-normally-and-cannot-be-swayed": "the swarm behaves as it normally would and can’t be swayed by the pipes’ music for the next 24 hours"}, b, "sway success")
    assert s["success_immunity_hours"] == 24
    tied([s["failure_result"]], {"friendly-to-you-and-your-allies": "becomes Friendly to you and your allies"}, b, "sway failure")
    tied([s["failure_lasts_while"]], {"you-play-the-pipes-each-round-as-magic-action": "for as long as you continue to play the pipes each round as a Magic action"}, b, "sway while")
    assert s["friendly_swarm_obeys_commands"] and "A Friendly swarm obeys your commands" in b
    tied([s["uncommanded_behavior"]], {"defends-itself-and-takes-no-other-actions": "If you issue no commands to a Friendly swarm, it defends itself but otherwise takes no actions"}, b, "sway uncommanded")
    tied([s["control_ends_when"]], {"friendly-swarm-starts-turn-more-than-30-feet-away": "If a Friendly swarm starts its turn more than 30 feet away from you, your control over that swarm ends"}, b, "sway control")
    assert s["control_end_distance_feet"] == 30
    tied([s["after_control_ends"]], {"behaves-normally-and-cannot-be-swayed-for-24-hours": "and the swarm behaves as it normally would and can’t be swayed by the pipes’ music for the next 24 hours"}, b, "sway after")

    c = detail["Plate Armor of Etherealness"]["plate_armor_of_etherealness"]
    b = block["Plate Armor of Etherealness"]
    tied([c["condition"]], {"while-wearing": "While you’re wearing this armor"}, b, "plate condition")
    tied([c["action"]], {"magic": "you can take a Magic action and use a command word"}, b, "plate action")
    tied([c["activation"]], {"command-word": "use a command word"}, b, "plate activation")
    tied([c["spell"]], {"etherealness": "gain the effect of the Etherealness spell"}, b, "plate spell")
    tied(c["ends_on"], {"armor-removed": "The spell ends immediately if you remove the armor", "magic-action-repeat-command-word": "or take a Magic action to repeat the command word"}, b, "plate end")
    assert (c["uses"], c["recharge"]) == (1, "next-dawn") and "can’t be used again until the next dawn" in b

    c = detail["Portable Hole"]["portable_hole"]
    b = block["Portable Hole"]
    assert c["source_pages"] == [234, 235] and "creates an extradimensional hole 10 feet deep" in b
    tied([c["cloth"]], {"fine-black-silk-soft-folds-to-handkerchief-size": "This fine black cloth, soft as silk, is folded up to the dimensions of a handkerchief"}, b, "hole cloth")
    assert c["diameter_feet"] == 6 and "It unfolds into a circular sheet 6 feet in diameter" in b
    u = c["unfold"]
    tied([u["action"]], {"magic": "You can take a Magic action to unfold a Portable Hole"}, b, "hole unfold")
    tied([u["surface"]], {"on-or-against-solid-surface": "place it on or against a solid surface"}, b, "hole surface")
    assert u["depth_feet"] == 10
    tied([u["shape"]], {"cylindrical": "The cylindrical space within the hole"}, b, "hole shape")
    tied([u["plane"]], {"different-plane-of-existence": "exists on a different plane of existence"}, b, "hole plane")
    assert u["cannot_create_open_passages"] and "so it can’t be used to create open passages" in b
    tied([u["exit"]], {"creature-inside-open-hole-can-climb-out": "Any creature inside an open Portable Hole can exit the hole by climbing out of it"}, b, "hole exit")
    k = c["close"]
    tied([k["action"]], {"magic": "You can take a Magic action to close a Portable Hole"}, b, "hole close")
    tied([k["contents"]], {"creatures-and-objects-remain-in-extradimensional-space": "any creatures or objects within remain in the extradimensional space"}, b, "hole contents")
    tied([k["weight"]], {"next-to-nothing": "the hole weighs next to nothing"}, b, "hole weight")
    e = c["escape_when_folded"]
    tied([e["action"]], {"action": "can take an action to make a DC 10 Strength (Athletics) check"}, b, "hole escape action")
    assert e["check"] == {"ability": "strength", "skill": "athletics", "dc": 10}
    tied([e["success"]], {"forces-its-way-out": "the creature forces its way out"}, b, "hole escape success")
    assert e["appears_within_feet"] == 5 and "appears within 5 feet of the Portable Hole" in b
    assert e["closed_air_hours"] == 1 and "holds enough air for 1 hour of breathing" in b
    tied([e["air_divided_by"]], {"number-of-breathing-creatures-inside": "divided by the number of breathing creatures inside"}, b, "hole air")
    n = c["inside_extradimensional_space"]
    tied(n["containers"], {"bag-of-holding": "a Bag of Holding", "handy-haversack": "Handy Haversack", "similar-item": "or similar item"}, b, "hole containers")
    tied([n["result"]], {"both-items-destroyed-and-gate-to-astral-plane-opens": "instantly destroys both items and opens a gate to the Astral Plane"}, b, "hole result")
    tied([n["gate_origin"]], {"where-one-item-was-placed-inside-the-other": "The gate originates where the one item was placed inside the other"}, b, "hole origin")
    assert n["pulls_in_within_feet"] == 10 and "Any creature within 10 feet of the gate" in b
    tied([n["exempt"]], {"creature-behind-total-cover": "and not behind Total Cover"}, b, "hole exempt")
    tied([n["destination"]], {"random-location-on-astral-plane": "deposited in a random location on the Astral Plane"}, b, "hole destination")
    assert n["gate_closes_afterward"] and "The gate then closes" in b
    assert n["one_way"] and not n["can_be_reopened"] and "The gate is one-way only and can’t be reopened" in b

    c = detail["Potion of Animal Friendship"]["potion_of_animal_friendship"]
    b = block["Potion of Animal Friendship"]
    tied([c["trigger"]], {"drink-this-potion": "When you drink this potion"}, b, "animal trigger")
    tied([c["spell"]], {"animal-friendship": "cast the level 3 version of the Animal Friendship spell"}, b, "animal spell")
    assert c["spell_level"] == 3 and c["save_dc"] == 13 and "(save DC 13)" in b
    tied([c["appearance"]], {"muddy-liquid": "Agitating this potion’s muddy liquid"}, b, "animal appearance")
    tied(c["appearance_items"], {"fish-scale": "a fish scale", "hummingbird-feather": "a hummingbird feather", "cat-claw": "a cat claw", "squirrel-hair": "or a squirrel hair"}, b, "animal items")
    print("#236 pinned PDF, twelve identities/pages/headers, 12 source-listed variants and source values pass")


if __name__ == "__main__":
    main()
