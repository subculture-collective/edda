#!/usr/bin/env python3
"""Compare issue #244 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (248, 249, 250)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wand": "Wand, ", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = ['Wand of Enemy Detection Wand']

# Each family's printed category line.
HEADERS = {
    'talisman-of-pure-good': 'Wondrous Item, Legendary (Requires Attunement by a Cleric or Paladin)',
    'talisman-of-ultimate-evil': 'Wondrous Item, Legendary (Requires Attunement)',
    'talisman-of-the-sphere': 'Wondrous Item, Legendary (Requires Attunement)',
    'thunderous-greatclub': 'Weapon (Greatclub), Very Rare (Requires Attunement)',
    'tome-of-clear-thought': 'Wondrous Item, Very Rare',
    'tome-of-leadership-and-influence': 'Wondrous Item, Very Rare',
    'tome-of-understanding': 'Wondrous Item, Very Rare',
    'trident-of-fish-command': 'Weapon (Trident), Uncommon (Requires Attunement)',
    'universal-solvent': 'Wondrous Item, Legendary',
    'vicious-weapon': 'Weapon (Any Simple or Martial), Rare',
    'vorpal-sword': 'Weapon (Glaive, Greatsword, Longsword, or Scimitar), Legendary (Requires Attunement)',
    'wand-of-binding': 'Wand, Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'talisman-of-pure-good': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'none'}, 'The talisman has 7 charges'),
    'talisman-of-ultimate-evil': ({'mode': 'limited', 'maximum': 6, 'per_use': 1, 'recharge': 'none'}, 'The talisman has 6 charges'),
    'trident-of-fish-command': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This magic weapon has 3 charges, and it regains 1d3 expended charges daily at dawn'),
    'wand-of-binding': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'talisman-of-pure-good': {
        'standard': [
            ('attunement_by_any_of[0]', 'cleric', 'Requires Attunement by a Cleric or Paladin'),
            ('attunement_by_any_of[1]', 'paladin', 'Requires Attunement by a Cleric or Paladin'),
            ('touch_damage.creature_types[0]', 'fiend', 'A Fiend or an Undead that touches the talisman'),
            ('touch_damage.creature_types[1]', 'undead', 'A Fiend or an Undead that touches the talisman'),
            ('touch_damage.damage', {'dice': 8, 'sides': 6, 'type': 'radiant'}, 'takes 8d6 Radiant damage'),
            ('touch_damage.repeats', 'end-of-each-turn-holding-or-carrying', 'takes the damage again each time it ends its turn holding or carrying the talisman'),
            ('holy_symbol', True, 'You can use the talisman as a Holy Symbol'),
            ('spell_attack_bonus', 2, 'You gain a +2 bonus to spell attack rolls while you wear or hold it'),
            ('max_charges', 7, 'The talisman has 7 charges'),
            ('rebuke.action', 'magic', 'you can take a Magic action to expend 1 charge'),
            ('rebuke.charge_cost', 1, 'a Magic action to expend 1 charge'),
            ('rebuke.target', 'one-visible-creature-on-ground', 'target one creature you can see on the ground'),
            ('rebuke.range_feet', 120, 'on the ground within 120 feet of yourself'),
            ('rebuke.save_ability', 'dexterity', 'DC 20 Dexterity saving throw'),
            ('rebuke.save_dc', 20, 'DC 20 Dexterity saving throw'),
            ('rebuke.disadvantage_types[0]', 'fiend', 'If the target is a Fiend or an Undead, it has Disadvantage on the save'),
            ('rebuke.disadvantage_types[1]', 'undead', 'If the target is a Fiend or an Undead, it has Disadvantage on the save'),
            ('rebuke.failure', 'destroyed-leaving-no-remains', 'the target falls into the fissure and is destroyed, leaving no remains'),
            ('rebuke.success_damage', {'dice': 4, 'sides': 6, 'type': 'psychic'}, 'takes 4d6 Psychic damage from the ordeal'),
            ('last_charge', 'disperses-into-golden-light-and-is-destroyed', 'the talisman disperses into motes of golden light and is destroyed'),
        ],
    },
    'talisman-of-ultimate-evil': {
        'standard': [
            ('touch_damage.excluded_creature_types[0]', 'fiend', 'A creature that isn’t a Fiend or an Undead that touches the talisman'),
            ('touch_damage.excluded_creature_types[1]', 'undead', 'A creature that isn’t a Fiend or an Undead that touches the talisman'),
            ('touch_damage.damage', {'dice': 8, 'sides': 6, 'type': 'necrotic'}, 'takes 8d6 Necrotic damage'),
            ('touch_damage.repeats', 'end-of-each-turn-holding-or-carrying', 'takes the damage again each time it ends its turn holding or carrying the talisman'),
            ('holy_symbol', True, 'You can use the talisman as a Holy Symbol'),
            ('spell_attack_bonus', 2, 'You gain a +2 bonus to spell attack rolls while you wear or hold it'),
            ('max_charges', 6, 'The talisman has 6 charges'),
            ('ultimate_end.action', 'magic', 'you can take a Magic action to expend 1 charge'),
            ('ultimate_end.charge_cost', 1, 'a Magic action to expend 1 charge'),
            ('ultimate_end.target', 'one-visible-creature-on-ground', 'target one creature you can see on the ground'),
            ('ultimate_end.range_feet', 120, 'on the ground within 120 feet of yourself'),
            ('ultimate_end.save_ability', 'dexterity', 'DC 20 Dexterity saving throw'),
            ('ultimate_end.save_dc', 20, 'DC 20 Dexterity saving throw'),
            ('ultimate_end.disadvantage_types[0]', 'celestial', 'If the target is a Celestial, it has Disadvantage on the save'),
            ('ultimate_end.failure', 'destroyed-leaving-no-remains', 'the target falls into the fissure and is destroyed, leaving no remains'),
            ('ultimate_end.success_damage', {'dice': 4, 'sides': 6, 'type': 'psychic'}, 'takes 4d6 Psychic damage from the ordeal'),
            ('last_charge', 'dissolves-into-slime-and-is-destroyed', 'the talisman dissolves into foul-smelling slime and is destroyed'),
        ],
    },
    'talisman-of-the-sphere': {
        'standard': [
            ('condition', 'while-holding-or-wearing', 'While holding or wearing this talisman'),
            ('check.ability', 'intelligence', 'any Intelligence (Arcana) check'),
            ('check.skill', 'arcana', 'any Intelligence (Arcana) check'),
            ('check.mode', 'advantage', 'you have Advantage on any Intelligence (Arcana) check'),
            ('check.purpose', 'control-sphere-of-annihilation', 'you make to control a Sphere of Annihilation'),
            ('move.timing', 'start-of-turn-while-in-control', 'when you start your turn in control of a Sphere of Annihilation'),
            ('move.action', 'magic', 'you can take a Magic action to move it'),
            ('move.base_feet', 10, 'move it 10 feet'),
            ('move.additional_feet_per_intelligence_modifier', 10, 'additional feet equal to 10 times your Intelligence modifier'),
            ('move.straight_line_required', False, 'This movement doesn’t have to be in a straight line'),
        ],
    },
    'thunderous-greatclub': {
        'standard': [
            ('strength.score', 20, 'your Strength is 20'),
            ('strength.condition', 'while-attuned-unless-already-equal-or-higher', 'unless your Strength is already equal to or greater than that score'),
            ('creature_extra_damage', {'dice': 1, 'sides': 8, 'type': 'thunder'}, 'an extra 1d8 Thunder damage to any creature it hits'),
            ('object_extra_damage', {'dice': 3, 'sides': 8, 'type': 'thunder'}, 'an extra 3d8 Thunder damage to objects it hits'),
            ('object_filter', 'not-worn-or-carried', 'objects it hits that aren’t being worn or carried'),
            ('clap.action', 'magic', 'Clap of Thunder. As a Magic action'),
            ('clap.strike', 'hard-surface', 'strike the weapon against a hard surface'),
            ('clap.audible_feet', 300, 'audible out to 300 feet'),
            ('clap.cone_feet', 30, 'create a 30-foot Cone of thunderous energy'),
            ('clap.save_ability', 'strength', 'DC 15 Strength saving throw'),
            ('clap.save_dc', 15, 'DC 15 Strength saving throw'),
            ('clap.condition', 'prone', 'or have the Prone condition'),
            ('clap.object_damage', {'dice': 3, 'sides': 8, 'type': 'thunder'}, 'carried take 3d8 Thunder damage'),
            ('earthquake.action', 'magic', 'Earthquake. As a Magic action'),
            ('earthquake.strike', 'ground', 'strike the weapon against the ground'),
            ('earthquake.circle_radius_feet', 50, 'in a 50-foot-radius circle centered on the point of impact'),
            ('earthquake.structure_damage', 50, 'take 50 Bludgeoning damage'),
            ('earthquake.structure_damage_type', 'bludgeoning', 'take 50 Bludgeoning damage'),
            ('earthquake.prone_save_ability', 'dexterity', 'must succeed on a DC 20 Dexterity saving throw or have the Prone condition'),
            ('earthquake.prone_save_dc', 20, 'must succeed on a DC 20 Dexterity saving throw or have the Prone condition'),
            ('earthquake.condition', 'prone', 'or have the Prone condition. If that creature is also concentrating'),
            ('earthquake.concentration_save_ability', 'constitution', 'DC 20 Constitution saving throw, or its Concentration is broken'),
            ('earthquake.concentration_save_dc', 20, 'DC 20 Constitution saving throw, or its Concentration is broken'),
            ('earthquake.fissure.depth_feet', 30, 'cause a 30-foot-deep'),
            ('earthquake.fissure.width_feet', 10, '10-foot-wide fissure'),
            ('earthquake.fissure.save_ability', 'dexterity', 'must make a DC 20 Dexterity saving throw, falling into the fissure'),
            ('earthquake.fissure.save_dc', 20, 'must make a DC 20 Dexterity saving throw, falling into the fissure'),
            ('earthquake.fissure.failure', 'falls-into-fissure', 'falling into the fissure on a failed save'),
            ('earthquake.fissure.success', 'moves-with-fissure-edge', 'or moving with the fissure’s edge on a successful one'),
            ('earthquake.fissure.structures', 'collapse-into-fissure', 'Any structure on a spot where the fissure opens collapses into the fissure'),
            ('earthquake.reuse_after', 'next-dawn', 'Once you use this property, it can’t be used again until the next dawn'),
        ],
    },
    'tome-of-clear-thought': {
        'standard': [
            ('contents', 'memory-and-logic-exercises', 'This book contains memory and logic exercises'),
            ('study.hours', 48, 'If you spend 48 hours'),
            ('study.within_days', 6, 'over a period of 6 days or fewer'),
            ('study.activity', 'study-contents-and-practice-guidelines', 'studying the book’s contents and practicing its guidelines'),
            ('ability', 'intelligence', 'your Intelligence increases by 2'),
            ('increase', 2, 'your Intelligence increases by 2'),
            ('maximum', 30, 'to a maximum of 30'),
            ('magic_after_use', 'lost', 'The manual then loses its magic but regains it in a century'),
            ('magic_returns_after_centuries', 1, 'regains it in a century'),
        ],
    },
    'tome-of-leadership-and-influence': {
        'standard': [
            ('contents', 'guidelines-for-influencing-and-charming', 'This book contains guidelines for influencing and charming others'),
            ('study.hours', 48, 'If you spend 48 hours'),
            ('study.within_days', 6, 'over a period of 6 days or fewer'),
            ('study.activity', 'study-contents-and-practice-guidelines', 'studying the book’s contents and practicing its guidelines'),
            ('ability', 'charisma', 'your Charisma increases by 2'),
            ('increase', 2, 'your Charisma increases by 2'),
            ('maximum', 30, 'to a maximum of 30'),
            ('magic_after_use', 'lost', 'The manual then loses its magic but regains it in a century'),
            ('magic_returns_after_centuries', 1, 'regains it in a century'),
        ],
    },
    'tome-of-understanding': {
        'standard': [
            ('contents', 'intuition-and-insight-exercises', 'This book contains intuition and insight exercises'),
            ('study.hours', 48, 'If you spend 48 hours'),
            ('study.within_days', 6, 'over a period of 6 days or fewer'),
            ('study.activity', 'study-contents-and-practice-guidelines', 'studying the book’s contents and practicing its guidelines'),
            ('ability', 'wisdom', 'your Wisdom increases by 2'),
            ('increase', 2, 'your Wisdom increases by 2'),
            ('maximum', 30, 'to a maximum of 30'),
            ('magic_after_use', 'lost', 'The manual then loses its magic, but regains it in a century'),
            ('magic_returns_after_centuries', 1, 'regains it in a century'),
        ],
    },
    'trident-of-fish-command': {
        'standard': [
            ('max_charges', 3, 'This magic weapon has 3 charges'),
            ('recharge_dice', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'regains 1d3 expended charges daily at dawn'),
            ('condition', 'while-carrying', 'While you carry it'),
            ('charge_cost', 1, 'you can expend 1 charge'),
            ('spell', 'dominate-beast', 'to cast Dominate Beast'),
            ('save_dc', 15, '(save DC 15)'),
            ('target', 'beast-with-swim-speed', 'on a Beast that has a Swim Speed'),
        ],
    },
    'universal-solvent': {
        'standard': [
            ('appearance', 'milky-liquid-with-strong-alcohol-smell', 'This tube holds milky liquid with a strong alcohol smell'),
            ('ounces_when_found', {'dice': 1, 'sides': 6, 'type': 'ounces'}, 'When found, a tube contains 1d6 + 1 ounces'),
            ('ounces_when_found_bonus', 1, '1d6 + 1 ounces'),
            ('action', 'utilize', 'You can take a Utilize action'),
            ('minimum_ounces_per_pour', 1, 'pour 1 or more ounces of solvent'),
            ('target', 'surface-within-reach', 'onto a surface within reach'),
            ('dissolves.maximum_square_feet', 1, 'up to 1 square foot of adhesive it touches'),
            ('dissolves.timing', 'instantly', 'Each ounce instantly dissolves'),
            ('dissolves.includes', 'sovereign-glue', 'including Sovereign Glue'),
        ],
    },
    'vicious-weapon': {
        'standard': [
            ('trigger', 'hit-creature', 'to any creature it hits'),
            ('extra_damage', {'dice': 2, 'sides': 6, 'type': 'same-as-weapon'}, 'an extra 2d6 damage'),
            ('damage_type_rule', 'same-as-weapon-normal-damage', 'This extra damage is of the same type as the weapon’s normal damage'),
        ],
    },
    'vorpal-sword': {
        'standard': [
            ('attack_bonus', 3, 'You gain a +3 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 3, 'You gain a +3 bonus to attack rolls and damage rolls'),
            ('ignores_resistance_to', 'slashing', 'the weapon ignores Resistance to Slashing damage'),
            ('head.target', 'creature-with-at-least-one-head', 'attack a creature that has at least one head'),
            ('head.trigger_natural_roll', 20, 'roll a 20 on the d20 for the attack roll'),
            ('head.effect', 'cut-off-one-head', 'you cut off one of the creature’s heads'),
            ('head.death', 'dies-if-it-cannot-survive-without-head', 'The creature dies if it can’t survive without the lost head'),
            ('head.immune_when[0]', 'immunity-to-slashing-damage', 'if it has Immunity to Slashing damage'),
            ('head.immune_when[1]', 'no-head-or-no-need-for-one', 'if it doesn’t have or need a head'),
            ('head.immune_when[2]', 'gm-decides-too-big', 'if the GM decides that the creature is too big for its head to be cut off'),
            ('head.immune_extra_damage', 30, 'takes an extra 30 Slashing damage from the hit'),
            ('head.immune_extra_damage_type', 'slashing', 'extra 30 Slashing damage'),
            ('head.legendary_resistance.cost', 'one-daily-use', 'it can expend one daily use of that trait'),
            ('head.legendary_resistance.result', 'keeps-head-and-takes-extra-damage', 'to avoid losing its head, taking the extra damage instead'),
        ],
    },
    'wand-of-binding': {
        'standard': [
            ('max_charges', 7, 'This wand has 7 charges'),
            ('save_dc', 17, '(save DC 17)'),
            ('spells[0].spell', 'hold-monster', 'Hold Monster 5'),
            ('spells[0].charge_cost', 5, 'Hold Monster 5'),
            ('spells[1].spell', 'hold-person', 'Hold Person 2'),
            ('spells[1].charge_cost', 2, 'Hold Person 2'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
}


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


def numbers(phrase):
    """The numbers a phrase prints, as digits or words. A printed 00 means 100."""
    digits = re.findall(r"\d+", re.sub(r"(?<=\d),(?=\d)", "", phrase))
    found = {int(d) for d in digits} | ({100} if "00" in digits else set())
    return found | {NUMBER_WORDS[w] for w in re.findall(r"[a-z]+", phrase.lower()) if w in NUMBER_WORDS}


def leaves(value, path=""):
    if isinstance(value, dict) and set(value) != {"dice", "sides", "type"}:
        for key, child in value.items():
            if key != "source_pages":
                yield from leaves(child, f"{path}.{key}" if path else key)
    elif isinstance(value, list) and value:
        for i, child in enumerate(value):
            yield from leaves(child, f"{path}[{i}]")
    else:
        yield path, value


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_244.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_244.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 244]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 244 and catalogue["batch"] == "data-magic-items-20"
    assert [(f["presentations"][0]["id"], f["name"], f["presentations"][0]["pages"][0]) for f in families] == expected
    assert [f["id"].rsplit("/", 1)[1] for f in families] == list(HEADERS) == list(VALUES)

    starts = {}
    for f in families:
        slug = f["id"].rsplit("/", 1)[1]
        printed = f"{f['name']} {HEADERS[slug]}"
        assert stream.count(printed) == 1, printed
        starts[slug] = stream.index(printed)
    bounds = sorted(starts.values())
    for heading in FOLLOWING:
        bounds.append(stream.index(heading, bounds[0]))
    count = 0
    for f in families:
        slug = f["id"].rsplit("/", 1)[1]
        start = starts[slug]
        end = min([b for b in bounds if b > start] + [len(stream)])
        block = stream[start:end].strip()
        assert len(CATEGORY_LINE.findall(block)) == 1, (slug, CATEGORY_LINE.findall(block))
        header = HEADERS[slug]
        assert header.startswith(CATEGORY[f["category"]]), (slug, header)
        pages = [p for p, (a, b) in spans.items() if a < start + len(block) and start < b]
        assert f["presentations"][0]["pages"] == [pages[0], pages[0]], (slug, pages)
        pool = CHARGES.get(slug)
        if pool is None:
            assert not re.search(r"(?:has|have) \d+ charges", block), slug
        else:
            assert has(block, pool[1]), (slug, pool[1])
            assert pool[0]["maximum"] in numbers(pool[1]), slug
        assert [v["id"].rsplit("/", 1)[1] for v in f["variants"]] == list(VALUES[slug]), slug
        for v in f["variants"]:
            vslug = v["id"].rsplit("/", 1)[1]
            label = f"{slug}/{vslug}"
            assert ("Requires Attunement" in header) == (v["attunement"]["mode"] == "required"), label
            assert "Rarity Varies" in header or RARITY[v["rarity"]] in header, label
            assert v["attunement"]["restrictions"] == [], label
            assert v["charges"] == (pool[0] if pool else {"mode": "none", "maximum": 0, "per_use": 0, "recharge": "none"}), label
            assert v["cost"] == {"kind": "unspecified", "amount": 0, "currency": "none", "source_pages": f["presentations"][0]["pages"]}, label
            assert len(v["details"]) == 1, label
            payload = next(iter(v["details"].values()))
            assert payload.get("rarity", v["rarity"]) == v["rarity"], label
            # A family whose text crosses a page records both pages.
            assert payload.get("source_pages") == (pages if len(pages) > 1 else None), (label, pages)
            want = VALUES[slug][vslug]
            got = list(leaves(payload))
            assert [p for p, _ in got] == [row[0] for row in want], (label, [p for p, _ in got])
            for (path, value), (_, printed_value, phrase) in zip(got, want):
                where = (label, path, value, phrase)
                assert type(value) is type(printed_value) and value == printed_value, where
                if phrase is None:
                    assert value == [], where
                    continue
                assert has(block, phrase), where
                if isinstance(value, dict):
                    assert f"{value['dice']}d{value['sides']}" in phrase, where
                elif type(value) is int:
                    assert value in numbers(phrase), where
                count += 1
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 248-250")


main()
