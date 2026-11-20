#!/usr/bin/env python3
"""Compare issue #232 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
# pp. 223 and 224 name eleven families; the Hat of Many Spells' table continues on
# p. 225, where the Headband of Intellect is printed.
PAGES = (223, 224, 225)
RARITY = {"uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# Printed category line of each family, before the rarity.
CATEGORY = {
    "giant-slayer": "Weapon (Any Simple or Martial)",
    "glamoured-studded-leather": "Armor (Studded Leather Armor)",
    "hammer-of-thunderbolts": "Weapon (Maul or Warhammer)",
}
# The next heading after the last family on p. 225.
FOLLOWING = "Helm of Brilliance Wondrous Item"


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
    variant = family["variants"][0]
    slug = family["id"].rsplit("/", 1)[1]
    category = CATEGORY.get(slug, "Wondrous Item")
    assert (family["category"] == "weapon") == category.startswith("Weapon")
    assert (family["category"] == "armor") == category.startswith("Armor")
    suffix = ""
    if variant["attunement"]["mode"] == "required":
        suffix = " (Requires Attunement by a Wizard)" if slug == "hat-of-many-spells" else " (Requires Attunement)"
    return f"{category}, {RARITY[variant['rarity']]}{suffix}"


def roll(outcome_min, outcome_max, printed):
    """The printed lead-in of one sub-roll outcome: 'on a 1, ', 'on an 8, ' or 'on a 4–5, '."""
    label = str(outcome_min) if outcome_min == outcome_max else f"{outcome_min}–{outcome_max}"
    return f"on {'an' if label.startswith('8') else 'a'} {label}, {printed}"


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_232.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream = " ".join(clean(text(pdf, p, x)) for p in PAGES for x in (0, 306))

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_232.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 232]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 232 and catalogue["batch"] == "data-magic-items-08"
    assert len(families) == 12 and all(len(f["variants"]) == 1 for f in families)
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["presentations"][0]["pages"][0] for f in families] == [223, 223, 223, 224, 224, 224, 224, 224, 224, 224, 224, 225]
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
        v = f["variants"][0]
        assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": page}, f["name"]
        assert v["attunement"]["restrictions"] == [] and v["name"] == f["name"] and v["id"].endswith("/variant/standard")
    for name in charges:
        if name not in ("Gem of Seeing", "Hammer of Thunderbolts"):
            assert charges[name] == none, name

    d = detail["Gem of Seeing"]["gem_of_seeing"]
    b = block["Gem of Seeing"]
    assert charges["Gem of Seeing"] == {"mode": "limited", "maximum": 3, "per_use": 1, "recharge": "regain-1d3-at-dawn"} and d["max_charges"] == 3
    assert has(b, "This gem has 3 charges")
    assert d["action"] == "magic" and d["charge_cost"] == 1 and has(b, "As a Magic action, you can expend 1 charge")
    assert d["duration_seconds"] == 600 and has(b, "For the next 10 minutes")
    assert d["sense"] == "truesight" and d["sense_range_feet"] == 120 and has(b, "you have Truesight out to 120 feet")
    tied([d["condition"]], {"peering-through-gem": "when you peer through the gem"}, b, "seeing condition")
    assert d["recharge"] == {"dice": 1, "sides": 3, "type": "charges"} and has(b, "The gem regains 1d3 expended charges daily at dawn")
    tied([d["recharge_when"]], {"daily-at-dawn": "expended charges daily at dawn"}, b, "seeing recharge")

    d = detail["Giant Slayer"]["giant_slayer"]
    b = block["Giant Slayer"]
    assert (d["attack_bonus"], d["damage_bonus"]) == (1, 1) and has(b, "You gain a +1 bonus to attack rolls and damage rolls made with this magic weapon")
    tied([d["trigger"]], {"hit-giant-with-weapon": "When you hit a Giant with this weapon"}, b, "slayer trigger")
    tied([d["target_creature_type"]], {"giant": "the Giant takes an extra"}, b, "slayer type")
    assert d["extra_damage"] == {"dice": 2, "sides": 6, "type": "weapon-type"} and has(b, "takes an extra 2d6 damage of the weapon’s type")
    assert (d["save_ability"], d["save_dc"], d["condition"]) == ("strength", 15, "prone") and has(b, "succeed on a DC 15 Strength saving throw or have the Prone condition")

    d = detail["Glamoured Studded Leather"]["glamoured_studded_leather"]
    b = block["Glamoured Studded Leather"]
    assert d["armor_class_bonus"] == 1 and has(b, "you gain a +1 bonus to Armor Class")
    assert d["action"] == "bonus" and has(b, "You can also take a Bonus Action to cause the armor")
    tied([d["appearance"]], {"normal-clothing-or-other-armor": "assume the appearance of a normal set of clothing or some other kind of armor"}, b, "glamour appearance")
    tied(d["chosen_details"], {"color": "including color, style, and accessories", "style": "color, style, and accessories", "accessories": "style, and accessories"}, b, "glamour details")
    assert d["retains_bulk_and_weight"] is True and has(b, "the armor retains its normal bulk and weight")
    tied(d["ends_on"], {"use-property-again": "lasts until you use this property again", "doff-armor": "or doff the armor"}, b, "glamour ends")

    d = detail["Gloves of Missile Snaring"]["gloves_of_missile_snaring"]
    b = block["Gloves of Missile Snaring"]
    tied([d["trigger"]], {"hit-by-ranged-or-thrown-weapon-attack-roll": "If you’re hit by an attack roll made with a Ranged or Thrown weapon while wearing these gloves"}, b, "snaring trigger")
    assert d["action"] == "reaction" and has(b, "you can take a Reaction to reduce the damage")
    assert d["reduction"] == {"dice": 1, "sides": 10, "type": "damage-reduction"} and d["reduction_bonus_ability"] == "dexterity"
    assert has(b, "reduce the damage by 1d10 plus your Dexterity modifier if you have a free hand") and d["requires_free_hand"] is True
    tied([d["catch_condition"]], {"damage-reduced-to-0": "If you reduce the damage to 0"}, b, "snaring catch condition")
    tied([d["catchable"]], {"ammunition-or-weapon-small-enough-to-hold-in-that-hand": "you can catch the ammunition or weapon if it is small enough for you to hold in that hand"}, b, "snaring catchable")

    d = detail["Gloves of Swimming and Climbing"]["gloves_of_swimming_and_climbing"]
    b = block["Gloves of Swimming and Climbing"]
    tied([d["condition"]], {"while-wearing": "While wearing these gloves"}, b, "swim condition")
    tied([d["climb_speed"]], {"equal-to-speed": "you have a Climb Speed and a Swim Speed equal to your Speed"}, b, "climb speed")
    tied([d["swim_speed"]], {"equal-to-speed": "a Climb Speed and a Swim Speed equal to your Speed"}, b, "swim speed")
    assert d["check_bonus"] == 5 and (d["check_ability"], d["check_skill"]) == ("strength", "athletics")
    assert has(b, "you gain a +5 bonus to Strength (Athletics) checks made to climb or swim")
    tied(d["check_subjects"], {"climb": "checks made to climb or swim", "swim": "made to climb or swim"}, b, "swim subjects")

    d = detail["Gloves of Thievery"]["gloves_of_thievery"]
    b = block["Gloves of Thievery"]
    tied([d["appearance"]], {"imperceptible-while-worn": "These gloves are imperceptible while worn"}, b, "thievery appearance")
    tied([d["condition"]], {"while-wearing": "While wearing them"}, b, "thievery condition")
    assert d["check_bonus"] == 5 and (d["check_ability"], d["check_skill"]) == ("dexterity", "sleight-of-hand")
    assert has(b, "you gain a +5 bonus to Dexterity (Sleight of Hand) checks")

    d = detail["Goggles of Night"]["goggles_of_night"]
    b = block["Goggles of Night"]
    tied([d["appearance"]], {"dark-lenses": "these dark lenses"}, b, "goggles appearance")
    tied([d["condition"]], {"while-wearing": "While wearing these dark lenses"}, b, "goggles condition")
    assert d["darkvision_feet"] == 60 and has(b, "you have Darkvision out to 60 feet")
    assert d["existing_darkvision_increase_feet"] == 60 and has(b, "If you already have Darkvision, wearing the goggles increases its range by 60 feet")

    d = detail["Hammer of Thunderbolts"]["hammer_of_thunderbolts"]
    b = block["Hammer of Thunderbolts"]
    assert (d["attack_bonus"], d["damage_bonus"]) == (1, 1) and has(b, "You gain a +1 bonus to attack rolls and damage rolls made with this magic weapon")
    assert charges["Hammer of Thunderbolts"] == {"mode": "limited", "maximum": 5, "per_use": 1, "recharge": "regain-1d4-plus-1-at-dawn"} and d["max_charges"] == 5
    assert has(b, "The weapon has 5 charges")
    h = d["hurl"]
    assert h["charge_cost"] == 1 and has(b, "You can expend 1 charge and make a ranged attack with the weapon")
    assert (h["thrown_normal_range_feet"], h["thrown_long_range_feet"]) == (20, 60)
    assert has(b, "hurling it as if it had the Thrown property with a normal range of 20 feet and a long range of 60 feet")
    assert h["thunderclap_audible_feet"] == 300 and h["thunderclap_requires_hit"] is True and has(b, "If the attack hits, the weapon unleashes a thunderclap audible out to 300 feet")
    assert h["burst_radius_feet"] == 30 and has(b, "The target and every creature within 30 feet of it other than you")
    tied([h["burst_excludes"]], {"you": "other than you"}, b, "hurl excludes")
    assert (h["save_ability"], h["save_dc"], h["condition"]) == ("constitution", 17, "stunned")
    assert has(b, "must succeed on a DC 17 Constitution saving throw or have the Stunned condition")
    tied([h["condition_ends"]], {"end-of-your-next-turn": "until the end of your next turn"}, b, "hurl condition end")
    tied([h["returns_to_hand"]], {"immediately-after-hitting-or-missing": "Immediately after hitting or missing, the weapon flies back to your hand"}, b, "hurl return")
    assert d["recharge"] == {"dice": 1, "sides": 4, "type": "charges"} and d["recharge_bonus"] == 1 and has(b, "regains 1d4 + 1 expended charges daily at dawn")
    tied([d["recharge_when"]], {"daily-at-dawn": "expended charges daily at dawn"}, b, "hammer recharge")
    g = d["giants_bane"]
    assert g["requires_attunement_to_weapon"] is True and has(b, "While you are attuned to the weapon and wearing either a Belt of Giant Strength or Gauntlets of Ogre Power")
    tied(g["requires_wearing_either"], {"belt-of-giant-strength": "wearing either a Belt of Giant Strength or", "gauntlets-of-ogre-power": "or Gauntlets of Ogre Power to which you are also attuned"}, b, "bane worn")
    assert g["requires_attunement_to_worn_item"] is True and has(b, "to which you are also attuned")
    assert g["natural_roll"] == 20 and has(b, "When you roll a 20 on the d20 for an attack roll made with this weapon against a Giant")
    tied([g["target_creature_type"]], {"giant": "against a Giant"}, b, "bane type")
    assert (g["save_ability"], g["save_dc"]) == ("constitution", 17) and has(b, "must succeed on a DC 17 Constitution saving throw or die")
    tied([g["failure"]], {"die": "saving throw or die"}, b, "bane failure")
    m = d["might_of_giants"]
    tied(m["affects"], {"belt-of-giant-strength": "bestowed by your Belt of Giant Strength or", "gauntlets-of-ogre-power": "or Gauntlets of Ogre Power increases by 4"}, b, "might affects")
    assert (m["strength_increase"], m["strength_maximum"]) == (4, 30) and has(b, "increases by 4, to a maximum of 30")

    d = detail["Handy Haversack"]["handy_haversack"]
    b = block["Handy Haversack"]
    assert d["weight_pounds"] == 5 and has(b, "The haversack always weighs 5 pounds, regardless of its contents")
    assert d["side_pouches"] == {"count": 2, "capacity_pounds": 200, "capacity_cubic_feet": 25}
    assert has(b, "central pouch and two side pouches") and has(b, "Each side pouch can hold up to 200 pounds of material, not exceeding a volume of 25 cubic feet")
    assert d["central_pouch"] == {"count": 1, "capacity_pounds": 500, "capacity_cubic_feet": 64}
    assert has(b, "The central pouch can hold up to 500 pounds of material, not exceeding a volume of 64 cubic feet")
    tied(d["retrieval_actions"], {"utilize": "requires a Utilize action or", "bonus-action": "or a Bonus Action (your choice)"}, b, "haversack retrieval")
    assert d["retrieved_item_on_top"] is True and has(b, "the item is always magically on top")
    tied(d["rupture_triggers"], {"overloaded": "If any of its pouches is overloaded,", "pierced": "overloaded, pierced, or torn", "torn": "pierced, or torn, the haversack ruptures and is destroyed"}, b, "haversack rupture")
    assert d["contents_lost_when_destroyed"] is True and has(b, "its contents are lost forever")
    assert d["artifact_turns_up_again"] is True and has(b, "although an Artifact always turns up again somewhere")
    tied([d["inside_out_result"]], {"contents-spill-unharmed-and-haversack-must-be-put-right": "its contents spill forth unharmed, and the haversack must be put right before it can be used again"}, b, "haversack inside out")
    assert d["air_minutes_per_pouch"] == 10 and has(b, "Each pouch of the haversack holds enough air for 10 minutes of breathing")
    tied([d["air_divided_by"]], {"breathing-creatures-inside": "divided by the number of breathing creatures inside"}, b, "haversack air")
    n = d["nested_extradimensional_space"]
    tied(n["items"], {"bag-of-holding": "created by a Bag of Holding,", "portable-hole": "Bag of Holding, Portable Hole, or", "similar-item": "Portable Hole, or similar item"}, b, "nested items")
    tied([n["result"]], {"both-items-destroyed": "instantly destroys both items"}, b, "nested result")
    tied([n["gate_to"]], {"astral-plane": "opens a gate to the Astral Plane"}, b, "nested gate")
    assert n["gate_radius_feet"] == 10 and has(b, "Any creature within 10 feet of the gate")
    tied([n["spared_by"]], {"total-cover": "not behind Total Cover is sucked through it"}, b, "nested cover")
    tied([n["destination"]], {"random-location-on-astral-plane": "deposited in a random location on the Astral Plane"}, b, "nested destination")
    assert n["gate_closes_after"] is True and has(b, "The gate then closes")
    assert n["gate_one_way"] is True and n["gate_reopens"] is False and has(b, "The gate is one-way only and can’t be reopened")

    d = detail["Hat of Disguise"]["hat_of_disguise"]
    b = block["Hat of Disguise"]
    tied([d["condition"]], {"while-wearing": "While wearing this hat"}, b, "disguise condition")
    tied([d["spell"]], {"disguise-self": "you can cast the Disguise Self spell"}, b, "disguise spell")
    tied([d["ends_when"]], {"hat-removed": "The spell ends if the hat is removed"}, b, "disguise end")

    d = detail["Hat of Many Spells"]["hat_of_many_spells"]
    b = block["Hat of Many Spells"]
    assert d["source_pages"] == [224, 225] and "86–90" in b and "96–00" in b and "Headband of Intellect Wondrous" not in b
    tied([d["attunement_class"]], {"wizard": "(Requires Attunement by a Wizard)"}, b, "hat class")
    tied([d["condition"]], {"while-holding": "While holding the hat, you can use it as a Spellcasting Focus"}, b, "hat condition")
    tied([d["focus_spell_list"]], {"wizard": "as a Spellcasting Focus for your Wizard spells"}, b, "hat focus")
    tied([d["added_somatic_component"]], {"reach-into-hat-and-pull-spell-out": "gains a special Somatic component: you must reach into the hat and “pull” the spell out of it"}, b, "hat somatic")
    u = d["unknown_spell"]
    assert u["min_spell_level"] == 1 and has(b, "try to cast a level 1+ spell you don’t know")
    tied([u["spell_list"]], {"wizard": "The spell must be on the Wizard spell list"}, b, "hat list")
    tied([u["level_limit"]], {"level-you-can-cast": "it must be of a level you can cast"}, b, "hat level")
    assert u["material_cost_limit_gp"] == 1000 and has(b, "it can’t have Material components costing more than 1,000 GP")
    tied([u["slot_expended"]], {"spell-level": "you must expend a spell slot of the spell’s level"}, b, "hat slot")
    assert (u["check_ability"], u["check_skill"], u["check_dc_base"]) == ("intelligence", "arcana", 10)
    assert has(b, "make an Intelligence (Arcana) check (DC 10 plus the spell’s level)")
    tied([u["check_dc_adds"]], {"spell-level": "(DC 10 plus the spell’s level)"}, b, "hat dc adds")
    tied([u["success_casting_time"]], {"normal-casting-time": "you cast the spell using its normal casting time"}, b, "hat success")
    tied([u["success_lockout_until"]], {"short-or-long-rest": "you can’t use this property again until you finish a Short or Long Rest"}, b, "hat lockout")
    tied([u["failure_result"]], {"fail-to-cast-and-random-effect": "you fail to cast the spell and a random effect occurs instead"}, b, "hat failure")
    assert u["table_die_sides"] == 100 and has(b, "determined by rolling on the following table") and has(b, "1d100 Effect")
    tied([d["spell_stats"]], {"your-spell-save-dc-and-spell-attack-bonus": "Any spell you cast from the hat uses your spell save DC and spell attack bonus"}, b, "hat stats")
    # Each printed row: its range and opening words, then its further roll. A printed 00 is stored as 100.
    printed_rows = [
        ("01–50", "You cast a random spell determined by rolling 1d10", {
            "enlarge-reduce-enlarge-effect": "Enlarge/Reduce (enlarge effect)", "enlarge-reduce-reduce-effect": "Enlarge/Reduce (reduce effect)", "faerie-fire": "Faerie Fire",
            "fireball": "Fireball", "gust-of-wind": "Gust of Wind", "invisibility-cast-on-yourself": "Invisibility (cast on yourself)", "lightning-bolt": "Lightning Bolt",
            "phantasmal-force": "Phantasmal Force", "polymorph": "Polymorph", "stinking-cloud": "Stinking Cloud"}, 10),
        ("51–55", "You have the Stunned condition until the end of your next turn, believing something awesome just happened.", {}, 0),
        ("56–60", "A harmless swarm of butterflies fills a 10-foot Cube within 30 feet of yourself. The swarm disperses after 1 minute.", {}, 0),
        ("61–65", "You pull a nonmagical object out of the hat. Roll 1d4 to determine the object", {
            "vial-of-acid": "a vial of Acid", "flask-of-alchemists-fire": "a flask of Alchemist’s Fire", "crowbar": "a Crowbar", "lit-torch": "a lit Torch"}, 4),
        ("66–70", "You suffer a bout of “magic sickness” and have the Poisoned condition for 1 hour.", {}, 0),
        ("71–75", "You have the Petrified condition until the end of your next turn.", {}, 0),
        ("76–80", "You pull a nonmagical object out of the hat. Roll 1d4 to determine the object", {
            "dagger": "a Dagger", "rope-with-grappling-hook-tied-to-one-end": "a Rope with a Grappling Hook tied to one end", "bag-of-caltrops": "a bag of Caltrops", "gem-worth-50-gp": "a gem worth 50 GP"}, 4),
        ("81–85", "A creature appears in an unoccupied space as close to you as possible. The creature isn’t under your control and acts as it normally would, and it disappears after 1 hour or when it drops to 0 Hit Points. Roll 1d4 to determine the creature", {
            "camel": "a Camel", "constrictor-snake": "a Constrictor Snake", "elephant": "an Elephant", "mule": "a Mule"}, 4),
        ("86–90", "A Hostile Swarm of Bats flies out of the hat, occupies your space, and attacks you.", {}, 0),
        ("91–95", "A vertical, 10-foot-diameter, two-way portal to another plane of existence opens in an unoccupied space within 30 feet of you and remains open until the end of your next turn. The GM determines where it leads.", {}, 0),
        ("96–00", "You pull a magic item out of the hat. Roll 1d6 to determine the item’s rarity", {"common": "Common", "uncommon": "Uncommon", "rare": "Rare"}, 6),
    ]
    table = d["failure_table"]
    assert len(table) == len(printed_rows) == 11
    for r, (span, opening, outcomes, sub) in zip(table, printed_rows):
        lo, hi = span.split("–")
        hi = 100 if hi == "00" else int(hi)
        assert (r["min"], r["max"]) == (int(lo), hi), span
        assert has(b, f"{span} {opening}"), span
        assert r["sub_die_sides"] == sub and len(r["outcomes"]) == len(outcomes), span
        assert [o["result"] for o in r["outcomes"]] == list(outcomes), span
        for o, printed in zip(r["outcomes"], outcomes.values()):
            assert has(b, roll(o["min"], o["max"], printed)), (span, o)
    assert table[-1]["max"] == 100 and has(b, "96–00 You pull a magic item")
    facts = {2: [("cube-size", 10, "feet"), ("range", 30, "feet"), ("disperses-after", 1, "minutes")], 4: [("duration", 1, "hours")], 6: [("gem-value", 50, "gp")],
             7: [("duration", 1, "hours")], 9: [("portal-diameter", 10, "feet"), ("range", 30, "feet")], 10: [("disappears-after", 1, "hours")]}
    for i, r in enumerate(table):
        assert [(x["key"], x["amount"], x["unit"]) for x in r["facts"]] == facts.get(i, []), i
    assert has(b, "A harmless swarm of butterflies fills a 10-foot Cube within 30 feet") and has(b, "disperses after 1 minute")
    assert has(b, "have the Poisoned condition for 1 hour") and has(b, "a gem worth 50 GP") and has(b, "it disappears after 1 hour")
    assert has(b, "portal to another plane of existence opens in an unoccupied space within 30 feet") and has(b, "which disappears after 1 hour if it’s not consumed or destroyed before then")
    traits = {1: ["until-end-of-your-next-turn", "believing-something-awesome-happened"], 2: ["swarm-fills-cube-within-range-of-yourself"], 5: ["until-end-of-your-next-turn"],
              7: ["unoccupied-space-as-close-to-you-as-possible", "not-under-your-control", "acts-as-it-normally-would", "disappears-when-duration-elapses", "disappears-at-0-hit-points"],
              8: ["flies-out-of-hat", "occupies-your-space", "attacks-you"],
              9: ["vertical", "unoccupied-space-within-range", "open-until-end-of-your-next-turn", "gm-determines-destination"],
              10: ["gm-chooses-item", "disappears-if-not-consumed-or-destroyed"]}
    for i, r in enumerate(table):
        assert r["traits"] == traits.get(i, []), i
    assert has(b, "The GM chooses the item") and has(b, "The creature isn’t under your control and acts as it normally would")

    d = detail["Headband of Intellect"]["headband_of_intellect"]
    b = block["Headband of Intellect"]
    tied([d["condition"]], {"while-wearing": "while you wear this headband"}, b, "headband condition")
    assert d["ability"] == "intelligence" and d["score"] == 19 and has(b, "Your Intelligence is 19 while you wear this headband")
    assert d["no_effect_at_or_above_score"] == 19 and has(b, "It has no effect on you if your Intelligence is 19 or higher without it")
    print("ok: 12 families, 12 variants, pp. 223-225")


main()
