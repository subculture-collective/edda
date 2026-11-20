#!/usr/bin/env python3
"""Compare issue #226 item identities, headers, tables and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (211, 212, 213)
CATEGORY = {"armor": "Armor", "weapon": "Weapon", "wondrous-item": "Wondrous Item"}
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
BASE = {("shield",): "Shield", ("battleaxe", "greataxe", "halberd"): "Battleaxe, Greataxe, or Halberd"}
UNITS = {"hours": "hour", "minutes": "minute", "days": "day"}
# Every string key of every payload, tied to the phrase it adapts. A list field
# maps its keys in catalogued order. Numbers, dice and table rows are compared
# separately in main(); "color" and "giant_types" follow their variant.
KEYS = {
    "arrow_catching": {
        "condition": {"while-wielding": "while you wield this Shield"},
        "ac_bonus_against": {"ranged-attack-rolls": "against ranged attack rolls"},
        "redirect_action": {"reaction": "you can take a Reaction"},
        "redirect_trigger": {"attacker-makes-ranged-attack-roll-against-target-within-range": "Whenever an attacker makes a ranged attack roll against a target within"},
        "redirect_effect": {"become-target-instead": "to become the target of the attack instead"},
    },
    "beans": {
        "material": {"heavy-cloth": "This heavy cloth bag"},
        "when_empty": {"becomes-nonmagical": "becomes a nonmagical item when it no longer contains any beans"},
        "dump_trigger": {"dump-one-or-more-beans": "If you dump one or more beans out of the bag"},
        "area_shape": {"sphere": "-foot-radius Sphere"},
        "area_center": {"dumped-beans": "centered on them"},
        "save": {"dexterity": "Dexterity saving throw"},
        "save_success": {"half-damage": "half as much damage on a successful one"},
        "plant_steps": {"remove-bean-from-bag": "If you remove a bean from the bag", "plant-in-dirt-or-sand": "plant it in dirt or sand", "water-it": "and then water it"},
        "effect_origin": {"ground-where-planted": "from the ground where it was planted"},
        "effect_selection": {"gm-chooses": "The GM can choose an effect from the following table", "random-table": "or determine it randomly"},
    },
    "devouring": {
        "resembles": {"bag-of-holding": "resembles a Bag of Holding"},
        "nature": {"feeding-orifice-of-gigantic-extradimensional-creature": "is a feeding orifice for a gigantic extradimensional creature"},
        "closes_orifice": {"turning-bag-inside-out": "Turning the bag inside out closes the orifice"},
        "creature_senses": {"whatever-is-placed-inside": "can sense whatever is placed inside the bag"},
        "devours": {"animal-matter-placed-wholly-inside": "Animal or vegetable matter placed wholly in the bag",
                    "vegetable-matter-placed-wholly-inside": "Animal or vegetable matter placed wholly in the bag"},
        "devoured_outcome": {"devoured-and-lost-forever": "is devoured and lost forever"},
        "pull_trigger": {"part-of-living-creature-placed-inside": "When part of a living creature is placed in the bag"},
        "escape_action": {"action": "A creature inside the bag can take an action to try to escape"},
        "rescue_action": {"action": "Another creature can take an action to reach into the bag"},
        "rescue_condition": {"puller-not-pulled-inside-first": "provided the puller isn’t pulled inside the bag first"},
        "turn_start_inside": {"creature-devoured-body-destroyed": "Any creature that starts its turn inside the bag is devoured, its body destroyed"},
        "stores": {"inanimate-objects": "Inanimate objects can be stored in the bag"},
        "swallow_frequency": {"once-each-day": "once each day"},
        "swallow_effect": {"objects-spat-out-into-another-plane-of-existence": "spits them out into another plane of existence"},
        "swallow_determined_by": {"gm-determines-time-and-plane": "The GM determines the time and plane"},
        "destroyed_by": {"pierced": "If the bag is pierced or torn", "torn": "If the bag is pierced or torn"},
        "destroyed_effect": {"contents-transported-to-random-astral-plane-location": "transported to a random location on the Astral Plane"},
    },
    "holding": {
        "retrieve_action": {"utilize": "requires a Utilize action"},
        "destroyed_by": {"overloaded": "If the bag is overloaded, pierced, or torn", "pierced": "If the bag is overloaded, pierced, or torn", "torn": "If the bag is overloaded, pierced, or torn"},
        "destroyed_effect": {"contents-scattered-in-astral-plane": "its contents are scattered in the Astral Plane"},
        "inside_out_effect": {"contents-spill-unharmed": "its contents spill forth unharmed"},
        "inside_out_reuse": {"must-be-put-right-before-reuse": "the bag must be put right before it can be used again"},
        "air_divided_by": {"breathing-creatures-inside": "divided by the number of breathing creatures inside"},
        "nesting_spaces": {"handy-haversack": "Handy Haversack", "portable-hole": "Portable Hole", "similar-item": "or similar item"},
        "nesting_effects": {"both-items-destroyed": "instantly destroys both items", "gate-to-astral-plane-opens": "opens a gate to the Astral Plane"},
        "gate_origin": {"where-one-item-was-placed-inside-the-other": "The gate originates where the one item was placed inside the other"},
        "gate_shape": {"sphere": "-foot-radius Sphere centered on the gate"},
        "gate_effect": {"creatures-sucked-to-random-astral-plane-location": "is sucked through it to a random location on the Astral Plane"},
    },
    "tricks": {
        "cloth_colors": {"gray": "made from gray, rust, or tan cloth", "rust": "made from gray, rust, or tan cloth", "tan": "made from gray, rust, or tan cloth"},
        "appearance": {"appears-empty": "appears empty"},
        "contents": {"small-fuzzy-object": "a small, fuzzy object"},
        "action": {"magic": "You can take a Magic action to pull the fuzzy object"},
        "transforms": {"on-landing-into-creature-from-color-table": "When the object lands, it transforms into a creature you determine by rolling on the table that corresponds to the bag’s color"},
        "stat_block_source": {"monsters": "See “Monsters” for the creature’s stat block"},
        "vanishes": {"next-dawn": "vanishes at the next dawn", "reduced-to-0-hit-points": "or when it is reduced to 0 Hit Points"},
        "attitude": {"friendly-to-user-and-allies": "Friendly to you and your allies"},
        "initiative": {"immediately-after-user": "acts immediately after you on your Initiative count"},
        "command_action": {"bonus-action": "You can take a Bonus Action to command"},
        "command_scope": {"movement-and-action-on-next-turn": "how the creature moves and what action it takes on its next turn"},
        "without_orders": {"acts-appropriate-to-nature": "acts in a fashion appropriate to its nature"},
        "uses_reset": {"next-dawn": "can’t be used again until the next dawn"},
    },
    "force": {
        "appearance": {"small-black-sphere": "This small black sphere"},
        "action": {"magic": "You can take a Magic action to throw the bead"},
        "area_shape": {"sphere": "-foot-radius Sphere on impact"},
        "save": {"dexterity": "Dexterity saving throw"},
        "save_success": {"no-damage": "must succeed on a DC 15 Dexterity saving throw or take"},
        "trapped_if_all": {"failed-save": "Any creature that failed the save", "completely-within-area": "and is completely within the area is trapped inside this sphere"},
        "pushed_if_any": {"succeeded-on-save": "Creatures that succeeded on the save", "partially-within-area": "or are partially within the area are pushed away"},
        "push_effect": {"pushed-from-center-until-outside": "pushed away from the center of the sphere until they are no longer inside it"},
        "passes_through": {"breathable-air": "Only breathable air can pass through the sphere’s wall"},
        "blocked": {"attacks": "No attack or other effect can pass through", "other-effects": "No attack or other effect can pass through"},
        "move_action": {"utilize": "can take a Utilize action to push against the sphere’s wall"},
        "move_actor": {"enclosed-creature": "An enclosed creature"},
        "move_distance": {"half-creature-speed": "up to half the creature’s Speed"},
    },
    "nourishment": {
        "appearance": {"flavorless-gelatinous-bead": "This flavorless, gelatinous bead"},
        "consumed_by": {"dissolves-on-tongue": "dissolves on your tongue"},
        "equivalent": {"rations": "of Rations"},
    },
    "dwarvenkind": {
        "condition": {"while-wearing": "While wearing this belt"},
        "benefits": {"dwarvish": "Dwarvish. You know Dwarvish", "friend-of-dwarvenkind": "Friend of Dwarvenkind.", "toughness": "Toughness."},
        "language": {"dwarvish": "You know Dwarvish"},
        "check_ability": {"charisma": "Charisma (Persuasion) checks"},
        "check_skill": {"persuasion": "Charisma (Persuasion) checks"},
        "check_audience": {"dwarves": "with dwarves and duergar", "duergar": "with dwarves and duergar"},
        "beard_condition": {"while-attuned": "while attuned to the belt"},
        "beard_timing": {"each-day-at-dawn": "each day at dawn"},
        "beard_effects": {"full-beard-if-able-to-grow-one": "growing a full beard if you can grow one", "thicker-beard-if-already-bearded": "or a thicker beard if you already have one"},
        "additional_condition": {"not-dwarf-or-duergar": "If you aren’t a dwarf or duergar"},
        "additional_benefits": {"darkvision": "Darkvision. You have Darkvision", "resilience": "Resilience."},
        "resistance": {"poison": "Resistance to Poison damage"},
        "save_advantage": {"avoid-or-end-poisoned-condition": "Advantage on saving throws you make to avoid or end the Poisoned condition"},
    },
    "giant_strength": {
        "condition": {"while-wearing": "While wearing this belt"},
        "ability": {"strength": "your Strength changes to a score granted by the belt"},
        "score_determined_by": {"giant-type": "The type of giant determines the score"},
    },
    "berserker": {
        "hit_point_condition": {"while-attuned": "while you are attuned to this weapon"},
        "curse_trigger": {"becoming-attuned": "becoming attuned to it extends the curse to you"},
        "curse_effects": {"unwilling-to-part-with-weapon": "you are unwilling to part with the weapon", "keeps-weapon-within-reach": "keeping it within reach at all times",
                          "disadvantage-on-attack-rolls-with-other-weapons": "Disadvantage on attack rolls with weapons other than this one"},
        "berserk_trigger": {"another-creature-damages-user-while-weapon-in-possession": "Whenever another creature damages you while the weapon is in your possession"},
        "berserk_save": {"wisdom": "Wisdom saving throw"},
        "berserk_save_failure": {"go-berserk": "or go berserk"},
        "berserk_ends": {"start-of-turn-with-no-creatures-in-range-seen-or-heard": "This berserk state ends when you start your turn and there are no creatures within"},
        "berserk_target": {"nearest-creature-seen-or-heard": "you regard the creature nearest to you that you can see or hear as your enemy"},
        "berserk_tie_break": {"choose-one-at-random": "choose one at random"},
        "berserk_compulsion": {"move-as-close-to-target-as-possible": "you must move as close to the creature as possible", "take-attack-action-targeting-it": "and take the Attack action, targeting the creature"},
        "berserk_unreachable": {"turn-ends-after-all-available-movement-used": "your turn ends after you’ve used up all your available movement"},
        "berserk_retarget_triggers": {"target-dies": "If the creature dies", "target-no-longer-seen-or-heard": "or can no longer be seen or heard by you"},
        "berserk_retarget": {"next-nearest-creature-seen-or-heard": "the next nearest creature that you can see or hear becomes your new target"},
    },
    "elvenkind": {
        "condition": {"while-wearing": "While you wear these boots"},
        "advantage_ability": {"dexterity": "Dexterity (Stealth) checks"},
        "advantage_skill": {"stealth": "Dexterity (Stealth) checks"},
    },
    "levitation": {
        "condition": {"while-wearing": "While you wear these boots"},
        "spell": {"levitate": "you can cast Levitate"},
        "target": {"self": "on yourself"},
    },
}
SPECIAL = {"tricks": {"color"}, "giant_strength": {"giant_types"}}
# Bag of Beans rows: the effect key and each fact key, tied to the row's text.
BEAN_ROWS = [
    ("toadstools-sprout", "toadstools sprout", {"eater-rolls-any-die": "If a creature eats a toadstool, roll any die",
                                                "odd-roll-save-or-poison-damage-and-poisoned": "On an odd roll, the eater must succeed on",
                                                "even-roll-gains-temporary-hit-points": "On an even roll, the eater gains"}),
    ("geyser-erupts", "A geyser erupts and spouts", {"gm-chooses-liquid": "(GM’s choice)", "liquid-water": "water,", "liquid-beer": "beer,", "liquid-mayonnaise": "mayonnaise,",
                                                    "liquid-tea": "tea,", "liquid-vinegar": "vinegar,", "liquid-wine": "wine,", "liquid-oil": "or oil"}),
    ("treant-sprouts", "A Treant sprouts", {"roll-any-die": "Roll any die", "odd-roll-chaotic-evil": "On an odd roll, the treant is Chaotic Evil",
                                            "even-roll-chaotic-good": "On an even roll, the treant is Chaotic Good"}),
    ("threatening-statue-rises", "stone statue in your likeness rises", {
        "animate-but-immobile-stone-statue-in-user-likeness": "An animate but immobile stone statue in your likeness rises",
        "makes-verbal-threats-against-user": "makes verbal threats against you",
        "if-user-leaves-directs-newcomers-to-find-and-attack-user": "directs the newcomers to find and attack you",
        "knows-user-location-on-same-plane": "If you are on the same plane of existence as the statue, it knows where you are"}),
    ("green-campfire-springs-forth", "A campfire with green flames springs forth", {"green-flames": "A campfire with green flames", "ends-early-if-extinguished": "or until it is extinguished"}),
    ("shrieker-fungi-sprout", "Three Shrieker Fungi sprout", {}),
    ("bright-pink-toads-crawl-forth", "bright-pink toads crawl forth", {
        "touched-toad-transforms-into-large-or-smaller-monster-of-gm-choice": "Whenever a toad is touched, it transforms into a Large or smaller monster of the GM’s choice",
        "monster-acts-by-alignment-and-nature": "acts in accordance with its alignment and nature",
        "monster-disappears-in-bright-pink-smoke": "then disappears in a puff of bright-pink smoke"}),
    ("hungry-bulette-attacks", "A hungry Bulette burrows up and attacks", {"burrows-up-and-attacks": "burrows up and attacks"}),
    ("fruit-tree-grows", "A fruit tree grows", {"potion-fruit-act-as-randomly-determined-potions": "of which act as randomly determined potions", "picked-fruit-remains": "Picked fruit remains"}),
    ("rainbow-egg-nest-springs-up", "rainbow-colored eggs springs up", {
        "any-creature-that-eats-an-egg-saves": "Any creature that eats an egg makes a", "increase-is-permanent": "a creature permanently increases its lowest ability score",
        "randomly-choose-among-equally-low-scores": "randomly choosing among equally low scores", "failed-save-internal-explosion": "from an internal explosion"}),
    ("pyramid-bursts-upward", "60-foot-square base bursts upward", {
        "burial-chamber-inside": "Inside is a burial chamber",
        "contains-mummy-mummy-lord-or-other-undead-of-gm-choice": "containing a Mummy, a Mummy Lord, or some other Undead of the GM’s choice",
        "sarcophagus-contains-treasure-of-gm-choice": "Its sarcophagus contains treasure of the GM’s choice"}),
    ("giant-beanstalk-sprouts", "A giant beanstalk sprouts", {
        "height-of-gm-choice": "growing to a height of the GM’s choice", "top-leads-where-gm-chooses": "The top leads where the GM chooses",
        "example-great-view": "such as to a great view", "example-cloud-giant-castle": "a cloud giant’s castle", "example-another-plane-of-existence": "or another plane of existence"}),
]


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
    assert base["excluded"] == [], family["name"]
    if base["options"]:
        line += " (" + BASE[tuple(base["options"])] + ")"
    rarities = {v["rarity"] for v in family["variants"]}
    # Only Belt of Giant Strength prints "Rarity Varies"; its table gives each row's rarity.
    rarity = RARITY[rarities.pop()] if len(rarities) == 1 else "Rarity Varies"
    attunement = {v["attunement"]["mode"] for v in family["variants"]}
    assert len(attunement) == 1
    return f"{line}, {rarity}" + (" (Requires Attunement)" if attunement == {"required"} else "")


def section(page, start, end):
    begin = page.index(start)
    return page[begin:page.index(end, begin + len(start))]


def tied(values, phrases, source, label):
    """Require the catalogued keys to be exactly the mapped ones, each with its printed phrase."""
    assert list(values) == list(phrases), f"{label}: {values}"
    for value, phrase in phrases.items():
        assert phrase in source, f"{label} {value}: {phrase}"


def tie_payload(name, payload, source):
    """Tie every string key of a payload; an unmapped string field fails."""
    fields = KEYS[name]
    strings = [k for k, v in payload.items() if isinstance(v, str) or isinstance(v, list) and v and isinstance(v[0], str)]
    assert strings == [k for k in payload if k in fields or k in SPECIAL.get(name, ())], f"{name}: unmapped or missing string field in {strings}"
    for field, phrases in fields.items():
        value = payload[field]
        tied([value] if isinstance(value, str) else value, phrases, source, f"{name} {field}")


def roll(q):
    """Print a quantity the way the source does: 5d4, 1d4 + 4 or a bare number."""
    if not q["dice"]:
        return str(q["bonus"])
    return f"{q['dice']}d{q['sides']}" + (f" + {q['bonus']}" if q["bonus"] else "")


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_226.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Read columns independently so continuation lines do not interleave. The
    # page footer closes the left column; drop it so each page reads left, then right.
    pages = {p: " ".join(" ".join(text(pdf, p, x).split()) for x in (0, 306)) for p in PAGES}
    for p in PAGES:
        footer = f" {p} System Reference Document 5.2.1"
        assert pages[p].count(footer) == 1, p
        pages[p] = pages[p].replace(footer, "")
    joined = {p: re.sub(r"(\w)- (\w)", r"\1\2", t) for p, t in pages.items()}

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_226.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 226]
    planned = [b for b in inventory["batches"] if b["key"] == "data-magic-items-02"]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 226 and catalogue["batch"] == "data-magic-items-02"
    assert len(families) == 12 and sum(len(f["variants"]) for f in families) == 18
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert len(planned) == 1 and planned[0]["issue"] == 226 and planned[0]["pages"] == "211, 212, 213"
    assert [(m["name"], m["source_page"]) for m in planned[0]["members"]] == [(name, page) for _, name, page in expected]
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
    # The batch owns p. 211 from Arrow-Catching Shield, all of p. 212 and p. 213;
    # no GP price or charge is printed for any of the twelve families.
    own = " ".join((joined[211][joined[211].index("Arrow-Catching Shield Armor"):], joined[212], joined[213]))
    assert " GP" not in own and "charge" not in own.lower()
    # Exactly twelve category lines are printed inside the batch's pages, in order.
    lines = re.findall(r"(?:Armor \(|Weapon \(|Wondrous Item, )[^.]*?(?:Common|Uncommon|Rare|Legendary|Varies)(?: \(Requires Attunement\))?", own)
    assert lines == [header(f) for f in families], lines
    details = {f["name"]: [v["details"] for v in f["variants"]] for f in families}
    names = {f["name"]: [v["name"] for v in f["variants"]] for f in families}

    shield = details["Arrow-Catching Shield"][0]["arrow_catching"]
    for snippet in (f"You gain a +{shield['ac_bonus']} bonus to Armor Class against ranged attack rolls while you wield this Shield",
                    "This bonus is in addition to the Shield’s normal bonus to AC",
                    f"Whenever an attacker makes a ranged attack roll against a target within {shield['redirect_range_feet']} feet of you, "
                    "you can take a Reaction to become the target of the attack instead"):
        assert snippet in joined[211], snippet
    assert shield["adds_to_shield_bonus"]
    tie_payload("arrow_catching", shield, section(joined[211], "Arrow-Catching Shield Armor", "Bag of Beans Wondrous Item"))

    # Bag of Beans: the bag, the explosion, then every 1d100 row across both columns.
    beans = details["Bag of Beans"][0]["beans"]
    block = joined[211][joined[211].index("Bag of Beans Wondrous Item"):]
    damage = beans["damage"]
    for snippet in (f"heavy cloth bag contains {beans['beans_dice']}d{beans['beans_die_sides']} dry beans when found",
                    "weighs half a pound regardless of how many beans it contains", "becomes a nonmagical item when it no longer contains any beans",
                    f"dump one or more beans out of the bag, they explode in a {beans['area_radius_feet']}-foot-radius Sphere centered on them",
                    "All the dumped beans are destroyed in the explosion", "each creature in the Sphere, including you",
                    f"makes a DC {beans['save_dc']} {beans['save'].title()} saving throw, taking {damage['dice']}d{damage['sides']} {damage['type'].title()} damage "
                    "on a failed save or half as much damage on a successful one",
                    "remove a bean from the bag, plant it in dirt or sand, and then water it", "the bean disappears",
                    "produces an effect 1 minute later from the ground where it was planted",
                    "The GM can choose an effect from the following table or determine it randomly", "1d100 Effect"):
        assert snippet in block, snippet
    assert (beans["weight_pounds_numerator"], beans["weight_pounds_denominator"], beans["plant_delay_seconds"], beans["table_die_sides"]) == (1, 2, 60, 100)
    assert beans["weight_ignores_contents"] and beans["dumped_beans_destroyed"] and beans["affects_user"] and beans["bean_disappears"]
    tie_payload("beans", beans, block.split("1d100 Effect ", 1)[0])
    # The table starts in the left column (01, 02–10) and continues in the right (11–20 … 96–00).
    table = block.split("1d100 Effect ", 1)[1].replace("1d100 Effect ", "")
    # A row label is 01 or a range; "DC 15 Constitution" inside a row is not one.
    chunks = re.split(r"(?:^|(?<!DC) )(\d\d)(?:–(\d\d))? (?=[A-Z0-9])", table)[1:]
    rows = [(int(a), int(b) if b else int(a), body) for a, b, body in zip(chunks[0::3], chunks[1::3], chunks[2::3])]
    rows = sorted((a, b or 100, body) for a, b, body in rows)
    assert [(a, b) for a, b, _ in rows] == [(r["min"], r["max"]) for r in beans["effects"]] and len(rows) == 12, [(a, b) for a, b, _ in rows]
    for row, (_, _, body) in zip(beans["effects"], rows):
        printed = sorted(re.findall(r"\d+d\d+(?: \+ \d+)?|\d+", body))
        for q in row["quantities"]:
            if q["label"] in ("treants", "bulettes", "beanstalks"):
                assert q == {"label": q["label"], "dice": 0, "sides": 0, "bonus": 1, "unit": "count"} and re.match(r"A (hungry |giant )?(Treant|Bulette|beanstalk)", body), body
                continue
            if q["label"] == "shrieker-fungi":
                assert q["bonus"] == 3 and body.startswith("Three Shrieker Fungi"), body
                continue
            want = roll(q)
            assert want in printed, (row["effect"], want, body)
            printed.remove(want)
            if q["unit"] in UNITS:
                assert f"{want} {UNITS[q['unit']]}" in body, (want, body)
            elif q["unit"].endswith("-damage"):
                assert f"{want} {q['unit'].removesuffix('-damage').title()} damage" in body, (want, body)
            elif q["unit"].endswith("-save-dc"):
                assert f"DC {want} {q['unit'].removesuffix('-save-dc').title()} saving throw" in body, (want, body)
        # Every number printed in the row is catalogued.
        assert printed == [], (row["effect"], printed)
        for creature in row["creatures"]:
            # The row prints the plural "Shrieker Fungi"; the key is the singular creature.
            printed_name = {"shrieker-fungus": "shrieker fungi"}.get(creature, creature.replace("-", " "))
            assert printed_name in body.lower(), (creature, body)
    # Each row's effect and fact keys are tied to that row's own text.
    assert [r["effect"] for r in beans["effects"]] == [effect for effect, _, _ in BEAN_ROWS]
    for row, (_, _, body), (effect, phrase, facts) in zip(beans["effects"], rows, BEAN_ROWS):
        assert phrase in body, (effect, phrase)
        tied(row["facts"], facts, body, f"bean row {effect} facts")

    dev = details["Bag of Devouring"][0]["devouring"]
    block = section(joined[212], "Bag of Devouring Wondrous Item", "Bag of Holding Wondrous Item")
    for snippet in ("resembles a Bag of Holding but is a feeding orifice for a gigantic extradimensional creature", "Turning the bag inside out closes the orifice",
                    "can sense whatever is placed inside the bag", "Animal or vegetable matter placed wholly in the bag is devoured and lost forever",
                    f"When part of a living creature is placed in the bag, as happens when someone reaches inside it, there is a {dev['pull_chance_percent']} percent chance that the creature is pulled inside",
                    f"can take an action to try to escape, doing so with a successful DC {dev['escape_check']['dc']} Strength (Athletics) check",
                    f"can take an action to reach into the bag to pull a creature out, doing so with a successful DC {dev['rescue_check']['dc']} Strength (Athletics) check, provided the puller isn’t pulled inside the bag first",
                    "Any creature that starts its turn inside the bag is devoured, its body destroyed", "Inanimate objects can be stored in the bag, which can hold a cubic foot of such material",
                    "once each day, the bag swallows any objects inside it and spits them out into another plane of existence", "The GM determines the time and plane",
                    "If the bag is pierced or torn, it is destroyed, and anything contained within it is transported to a random location on the Astral Plane"):
        assert snippet in block, snippet
    assert dev["capacity_cubic_feet"] == 1
    tie_payload("devouring", dev, block)
    assert {dev[k]["ability"] + "-" + dev[k]["skill"] for k in ("escape_check", "rescue_check")} == {"strength-athletics"}

    hold = details["Bag of Holding"][0]["holding"]
    block = section(joined[212], "Bag of Holding Wondrous Item", "Bag of Tricks Wondrous Item")
    for snippet in (f"roughly {hold['interior_square_side_feet']} feet square and {hold['interior_depth_feet']} feet deep on the inside",
                    f"can hold up to {hold['capacity_pounds']} pounds, not exceeding a volume of {hold['capacity_cubic_feet']} cubic feet",
                    f"The bag weighs {hold['weight_pounds']} pounds, regardless of its contents", "Retrieving an item from the bag requires a Utilize action",
                    "If the bag is overloaded, pierced, or torn, it is destroyed, and its contents are scattered in the Astral Plane",
                    "If the bag is turned inside out, its contents spill forth unharmed, but the bag must be put right before it can be used again",
                    f"enough air for {hold['air_minutes']} minutes of breathing, divided by the number of breathing creatures inside",
                    "inside an extradimensional space created by a Handy Haversack, Portable Hole, or similar item instantly destroys both items and opens a gate to the Astral Plane",
                    "The gate originates where the one item was placed inside the other",
                    f"Any creature within a {hold['gate_radius_feet']}-foot-radius Sphere centered on the gate is sucked through it to a random location on the Astral Plane",
                    "The gate then closes. The gate is one-way and can’t be reopened"):
        assert snippet in block, snippet
    assert hold["gate_closes"] and hold["gate_one_way"] and not hold["gate_reopenable"] and hold["interior_approximate"] and hold["weight_ignores_contents"]
    tie_payload("holding", hold, block)

    # Bag of Tricks: the shared text, then one two-column 1d8 table per cloth colour.
    tricks = [d["tricks"] for d in details["Bag of Tricks"]]
    block = section(joined[212], "Bag of Tricks Wondrous Item", "Bead of Force Wondrous Item")
    first = tricks[0]
    for snippet in ("This bag made from gray, rust, or tan cloth appears empty", "reveals the presence of a small, fuzzy object",
                    f"You can take a Magic action to pull the fuzzy object from the bag and throw it up to {first['throw_feet']} feet",
                    "When the object lands, it transforms into a creature you determine by rolling on the table that corresponds to the bag’s color",
                    "See “Monsters” for the creature’s stat block", "The creature vanishes at the next dawn or when it is reduced to 0 Hit Points",
                    "The creature is Friendly to you and your allies, and it acts immediately after you on your Initiative count",
                    "You can take a Bonus Action to command how the creature moves and what action it takes on its next turn",
                    "In the absence of such orders, the creature acts in a fashion appropriate to its nature",
                    "Once three fuzzy objects have been pulled from the bag, the bag can’t be used again until the next dawn"):
        assert snippet in block, snippet
    assert [t["color"] for t in tricks] == first["cloth_colors"] == ["gray", "rust", "tan"]
    tables = re.split(r"(Gray|Rust|Tan) Bag of Tricks ", block)[1:]
    assert tables[0::2] == [n.removesuffix(" Bag of Tricks") for n in names["Bag of Tricks"]]
    for trick, body in zip(tricks, tables[1::2]):
        assert body.startswith("1d8 Creature 1d8 Creature ") and (trick["uses"], trick["uses_reset"], trick["table_die_sides"]) == (3, "next-dawn", 8)
        printed = sorted((int(n), key(c)) for n, c in re.findall(r"(\d) ([A-Z][a-z]+(?: [A-Z][a-z]+)?)", body.removeprefix("1d8 Creature 1d8 Creature ")))
        assert printed == [(r["roll"], r["creature"]) for r in trick["creatures"]] and len(printed) == 8, printed
        tie_payload("tricks", trick, block)

    # Bead of Force starts on p. 212 and continues at the top of p. 213.
    force = details["Bead of Force"][0]["force"]
    assert force["source_pages"] == [212, 213]
    start = joined[212][joined[212].index("Bead of Force Wondrous Item"):]
    block = re.sub(r"(\w)- (\w)", r"\1\2", start + " " + joined[213].split(" Bead of Nourishment Wondrous Item")[0])
    damage = force["damage"]
    found = f"{force['found_dice']}d{force['found_die_sides']} + {force['found_bonus']}"
    for snippet in (f"This small black sphere measures {force['diameter_inches_numerator']}/{force['diameter_inches_denominator']} of an inch in diameter and weighs an ounce",
                    f"Typically, {found} Beads of Force are found together", f"You can take a Magic action to throw the bead up to {force['throw_feet']} feet",
                    f"The bead explodes in a {force['area_radius_feet']}-foot-radius Sphere on impact and is destroyed",
                    f"Each creature in the Sphere must succeed on a DC {force['save_dc']} {force['save'].title()} saving throw or take {damage['dice']}d{damage['sides']} {damage['type'].title()} damage",
                    "A sphere of transparent force then encloses the area for 1 minute", "Any creature that failed the save and is completely within the area is trapped inside this sphere",
                    "Creatures that succeeded on the save or are partially within the area are pushed away from the center of the sphere until they are no longer inside it",
                    "Only breathable air can pass through the sphere’s wall. No attack or other effect can pass through",
                    "An enclosed creature can take a Utilize action to push against the sphere’s wall, moving the sphere up to half the creature’s Speed",
                    f"The sphere can be picked up, and its magic causes it to weigh only {force['sphere_weight_pounds']} pound, regardless of the weight of creatures inside"):
        assert snippet in block, snippet
    assert (force["weight_ounces"], force["sphere_duration_seconds"]) == (1, 60) and "half as much" not in block
    assert force["destroyed_on_impact"] and force["can_be_picked_up"] and force["weight_ignores_creatures"]
    tie_payload("force", force, block)

    food = details["Bead of Nourishment"][0]["nourishment"]
    block = section(joined[213], "Bead of Nourishment Wondrous Item", "Belt of Dwarvenkind Wondrous Item")
    assert f"This flavorless, gelatinous bead dissolves on your tongue and provides as much nourishment as {food['nourishment_days']} day of Rations" in block
    tie_payload("nourishment", food, block)

    dwarf = details["Belt of Dwarvenkind"][0]["dwarvenkind"]
    block = section(joined[213], "Belt of Dwarvenkind Wondrous Item", "Belt of Giant Strength Wondrous Item")
    for snippet in ("While wearing this belt, you gain the following benefits", "Dwarvish. You know Dwarvish",
                    "Friend of Dwarvenkind. You have Advantage on Charisma (Persuasion) checks made to interact with dwarves and duergar",
                    f"Toughness. Your Constitution increases by {dwarf['constitution_increase']}, to a maximum of {dwarf['constitution_maximum']}",
                    f"while attuned to the belt, you have a {dwarf['beard_chance_percent']} percent chance each day at dawn of growing a full beard if you can grow one, or a thicker beard if you already have one",
                    "If you aren’t a dwarf or duergar, you gain the following additional benefits while wearing the belt",
                    f"Darkvision. You have Darkvision with a range of {dwarf['darkvision_feet']} feet", "Resilience. You have Resistance to Poison damage",
                    "You also have Advantage on saving throws you make to avoid or end the Poisoned condition"):
        assert snippet in block, snippet
    tie_payload("dwarvenkind", dwarf, block)
    # Benefit names are the printed run-in headings, in order.
    assert [key(h) for h in re.findall(r"(Dwarvish|Friend of Dwarvenkind|Toughness|Darkvision|Resilience)\. You", block)] == dwarf["benefits"] + dwarf["additional_benefits"]

    # Belt of Giant Strength: every row of the Belt / Str. / Rarity table.
    belts = [d["giant_strength"] for d in details["Belt of Giant Strength"]]
    for snippet in ("While wearing this belt, your Strength changes to a score granted by the belt", "The type of giant determines the score (see the table below)",
                    "The item has no effect on you if your Strength without the belt is equal to or greater than the belt’s score"):
        assert snippet in joined[213], snippet
    table = section(joined[213], "Belt Str. Rarity ", " Berserker Axe Weapon").removeprefix("Belt Str. Rarity ")
    # The "frost or stone" row wraps, so its closing word follows the rarity.
    table = table.replace("(frost or 23 Very Rare stone)", "(frost or stone) 23 Very Rare")
    printed = re.findall(r"(Belt of Giant Strength \(([a-z ]+)\)) (\d+) (Rare|Very Rare|Legendary)", table)
    assert len(printed) == 5 and " ".join(f"{n} {s} {r}" for n, _, s, r in printed) == table, table
    variants = next(f["variants"] for f in families if f["name"] == "Belt of Giant Strength")
    for (name, giants, score, rarity), belt, v in zip(printed, belts, variants):
        assert (name, giants.split(" or "), int(score), rarity) == (v["name"], belt["giant_types"], belt["score"], RARITY[v["rarity"]]), name
        assert belt["no_effect_if_score_at_least"] == belt["score"] and v["id"].endswith("/variant/" + key(giants))
        tie_payload("giant_strength", belt, section(joined[213], "Belt of Giant Strength Wondrous Item", " Berserker Axe Weapon"))

    axe = details["Berserker Axe"][0]["berserker"]
    block = section(joined[213], "Berserker Axe Weapon", " Boots of Elvenkind Wondrous Item")
    for snippet in (f"You gain a +{axe['attack_bonus']} bonus to attack rolls and damage rolls made with this magic weapon",
                    f"while you are attuned to this weapon, your Hit Point maximum increases by {axe['hit_point_maximum_per_level']} for each level you have attained",
                    "Curse. This weapon is cursed, and becoming attuned to it extends the curse to you",
                    "you are unwilling to part with the weapon, keeping it within reach at all times", "You also have Disadvantage on attack rolls with weapons other than this one",
                    f"Whenever another creature damages you while the weapon is in your possession, you must succeed on a DC {axe['berserk_save_dc']} {axe['berserk_save'].title()} saving throw or go berserk",
                    f"This berserk state ends when you start your turn and there are no creatures within {axe['berserk_end_range_feet']} feet of you that you can see or hear",
                    "While berserk, you regard the creature nearest to you that you can see or hear as your enemy", "If there are multiple possible creatures, choose one at random",
                    "On each of your turns, you must move as close to the creature as possible and take the Attack action, targeting the creature",
                    "your turn ends after you’ve used up all your available movement",
                    "If the creature dies or can no longer be seen or heard by you, the next nearest creature that you can see or hear becomes your new target"):
        assert snippet in block, snippet
    assert axe["cursed"] and axe["attack_bonus"] == axe["damage_bonus"] == 1
    tie_payload("berserker", axe, block)

    elf = details["Boots of Elvenkind"][0]["elvenkind"]
    assert elf["silent_steps"] and elf["any_surface"]
    assert "While you wear these boots, your steps make no sound, regardless of the surface you are moving across" in joined[213]
    assert f"You also have Advantage on {elf['advantage_ability'].title()} ({elf['advantage_skill'].title()}) checks" in joined[213]
    tie_payload("elvenkind", elf, section(joined[213], "Boots of Elvenkind Wondrous Item", "Boots of Levitation Wondrous Item"))

    lev = details["Boots of Levitation"][0]["levitation"]
    block = joined[213].split("Boots of Levitation Wondrous Item, Rare (Requires Attunement) ", 1)[1]
    assert block == (f"While you wear these boots, you can cast {lev['spell'].title()} on your{lev['target']}."), block
    tie_payload("levitation", lev, block)
    assert set(KEYS) == {name for variants in details.values() for d in variants for name in d}
    print("#226 pinned PDF, twelve identities/pages/headers, eighteen variants, tables and source values pass")


if __name__ == "__main__":
    main()
