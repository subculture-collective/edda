#!/usr/bin/env python3
"""Compare issue #245 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (250, 251, 252)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wand": "Wand, ", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = []

# Each family's printed category line.
HEADERS = {
    'wand-of-enemy-detection': 'Wand, Rare (Requires Attunement)',
    'wand-of-fear': 'Wand, Rare (Requires Attunement)',
    'wand-of-fireballs': 'Wand, Rare (Requires Attunement by a Spellcaster)',
    'wand-of-lightning-bolts': 'Wand, Rare (Requires Attunement by a Spellcaster)',
    'wand-of-magic-detection': 'Wand, Uncommon',
    'wand-of-magic-missiles': 'Wand, Uncommon',
    'wand-of-paralysis': 'Wand, Rare (Requires Attunement by a Spellcaster)',
    'wand-of-polymorph': 'Wand, Very Rare (Requires Attunement by a Spellcaster)',
    'wand-of-secrets': 'Wand, Uncommon',
    'wand-of-web': 'Wand, Uncommon (Requires Attunement by a Spellcaster)',
    'wand-of-wonder': 'Wand, Rare (Requires Attunement)',
    'wand-of-the-war-mage-1-2-or-3': 'Wand, Uncommon (+1), Rare (+2), or Very Rare (+3) (Requires Attunement by a Spellcaster)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'wand-of-enemy-detection': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-fear': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-fireballs': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-lightning-bolts': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-magic-detection': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This wand has 3 charges'),
    'wand-of-magic-missiles': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-paralysis': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-polymorph': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-secrets': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This wand has 3 charges and regains 1d3 expended charges daily at dawn'),
    'wand-of-web': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
    'wand-of-wonder': ({'mode': 'limited', 'maximum': 7, 'per_use': 1, 'recharge': 'regain-1d6-plus-1-at-dawn'}, 'This wand has 7 charges'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'wand-of-enemy-detection': {
        'standard': [
            ('max_charges', 7, 'This wand has 7 charges'),
            ('action', 'magic', 'you can take a Magic action to expend 1 charge'),
            ('charge_cost', 1, 'a Magic action to expend 1 charge'),
            ('duration_minutes', 1, 'For 1 minute'),
            ('reveals', 'direction-of-nearest-hostile-creature', 'you know the direction of the nearest creature Hostile to you'),
            ('range_feet', 60, 'within 60 feet'),
            ('reveals_distance', False, 'but not its distance from you'),
            ('senses[0]', 'invisible', 'creatures that are Invisible'),
            ('senses[1]', 'ethereal', 'Invisible, ethereal'),
            ('senses[2]', 'disguised', 'ethereal, disguised'),
            ('senses[3]', 'hidden', 'disguised, or hidden'),
            ('senses[4]', 'in-plain-sight', 'as well as those in plain sight'),
            ('ends_when', 'you-stop-holding-the-wand', 'The effect ends if you stop holding the wand'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-fear': {
        'standard': [
            ('max_charges', 7, 'This wand has 7 charges'),
            ('save_dc', 15, '(save DC 15)'),
            ('spells[0].spell', 'command', 'Command (“flee” or “grovel” only) 1'),
            ('spells[0].limit', 'flee-or-grovel-only', 'Command (“flee” or “grovel” only) 1'),
            ('spells[0].charge_cost', 1, 'Command (“flee” or “grovel” only) 1'),
            ('spells[1].spell', 'fear', 'Fear (60-foot Cone) 3'),
            ('spells[1].limit', '60-foot-cone', 'Fear (60-foot Cone) 3'),
            ('spells[1].charge_cost', 3, 'Fear (60-foot Cone) 3'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-fireballs': {
        'standard': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('max_charges', 7, 'This wand has 7 charges'),
            ('maximum_charges_per_cast', 3, 'expend no more than 3 charges'),
            ('spell', 'fireball', 'to cast Fireball'),
            ('save_dc', 15, '(save DC 15)'),
            ('base_charge_cost', 1, 'For 1 charge, you cast the level 3 version of the spell'),
            ('base_spell_level', 3, 'you cast the level 3 version of the spell'),
            ('levels_per_additional_charge', 1, 'You can increase the spell’s level by 1 for each additional charge you expend'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-lightning-bolts': {
        'standard': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('max_charges', 7, 'This wand has 7 charges'),
            ('maximum_charges_per_cast', 3, 'expend no more than 3 charges'),
            ('spell', 'lightning-bolt', 'to cast Lightning Bolt'),
            ('save_dc', 15, '(save DC 15)'),
            ('base_charge_cost', 1, 'For 1 charge, you cast the level 3 version of the spell'),
            ('base_spell_level', 3, 'you cast the level 3 version of the spell'),
            ('levels_per_additional_charge', 1, 'You can increase the spell’s level by 1 for each additional charge you expend'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-magic-detection': {
        'standard': [
            ('max_charges', 3, 'This wand has 3 charges'),
            ('charge_cost', 1, 'expend 1 charge to cast Detect Magic'),
            ('spell', 'detect-magic', 'cast Detect Magic'),
            ('recharge_dice', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'The wand regains 1d3 expended charges daily at dawn'),
        ],
    },
    'wand-of-magic-missiles': {
        'standard': [
            ('max_charges', 7, 'This wand has 7 charges'),
            ('maximum_charges_per_cast', 3, 'expend no more than 3 charges'),
            ('spell', 'magic-missile', 'to cast Magic Missile'),
            ('base_charge_cost', 1, 'For 1 charge, you cast the level 1 version of the spell'),
            ('base_spell_level', 1, 'you cast the level 1 version of the spell'),
            ('levels_per_additional_charge', 1, 'You can increase the spell’s level by 1 for each additional charge you expend'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-paralysis': {
        'standard': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('max_charges', 7, 'This wand has 7 charges'),
            ('action', 'magic', 'you can take a Magic action to expend 1 charge'),
            ('charge_cost', 1, 'a Magic action to expend 1 charge'),
            ('target', 'one-visible-creature', 'toward a creature you can see'),
            ('range_feet', 60, 'within 60 feet of yourself'),
            ('save_ability', 'constitution', 'DC 15 Constitution saving throw'),
            ('save_dc', 15, 'DC 15 Constitution saving throw'),
            ('condition', 'paralyzed', 'have the Paralyzed condition'),
            ('duration_minutes', 1, 'the Paralyzed condition for 1 minute'),
            ('repeat_save', 'end-of-each-of-its-turns', 'At the end of each of the target’s turns, it repeats the save'),
            ('repeat_success', 'ends-effect-on-itself', 'ending the effect on itself on a success'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-polymorph': {
        'standard': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('max_charges', 7, 'This wand has 7 charges'),
            ('charge_cost', 1, 'expend 1 charge to cast Polymorph'),
            ('spell', 'polymorph', 'to cast Polymorph'),
            ('save_dc', 15, '(save DC 15)'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-secrets': {
        'standard': [
            ('max_charges', 3, 'This wand has 3 charges'),
            ('recharge_dice', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'regains 1d3 expended charges daily at dawn'),
            ('action', 'magic', 'you can take a Magic action to expend 1 charge'),
            ('charge_cost', 1, 'a Magic action to expend 1 charge'),
            ('detects[0]', 'secret-door', 'if a secret door or trap'),
            ('detects[1]', 'trap', 'if a secret door or trap'),
            ('range_feet', 60, 'is within 60 feet of you'),
            ('result', 'pulses-and-points-at-nearest', 'the wand pulses and points at the one nearest to you'),
        ],
    },
    'wand-of-web': {
        'standard': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('max_charges', 7, 'This wand has 7 charges'),
            ('charge_cost', 1, 'expend 1 charge to cast Web'),
            ('spell', 'web', 'to cast Web'),
            ('save_dc', 13, '(save DC 13)'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into ashes'),
            ('last_charge.result', 'crumbles-into-ashes', 'the wand crumbles into ashes and is destroyed'),
        ],
    },
    'wand-of-wonder': {
        'standard': [
            ('max_charges', 7, 'This wand has 7 charges'),
            ('action', 'magic', 'you can take a Magic action to expend 1 charge'),
            ('charge_cost', 1, 'a Magic action to expend 1 charge'),
            ('point_range_feet', 120, 'choosing a point within 120 feet of yourself'),
            ('spell_save_dc', 15, 'Spells cast from the wand have a save DC of 15'),
            ('minimum_spell_range_feet', 120, 'it becomes 120 feet when cast from the wand'),
            ('multiple_subjects', 'gm-determines-randomly', 'the GM determines randomly which among them are affected'),
            ('table_die_sides', 100, '1d100 Effect'),
            ('effects[0].roll_from', 1, '01–20 You cast a spell originating from the chosen point'),
            ('effects[0].roll_to', 20, '01–20 You cast a spell originating from the chosen point'),
            ('effects[0].effect', 'cast-spell-from-1d10-table', '01–20 You cast a spell originating from the chosen point'),
            ('effects[1].roll_from', 21, '21–25 Nothing happens at the chosen point of origin'),
            ('effects[1].roll_to', 25, '21–25 Nothing happens at the chosen point of origin'),
            ('effects[1].effect', 'you-are-stunned', '21–25 Nothing happens at the chosen point of origin'),
            ('effects[2].roll_from', 26, '26–30 You cast Gust of Wind'),
            ('effects[2].roll_to', 30, '26–30 You cast Gust of Wind'),
            ('effects[2].effect', 'cast-gust-of-wind', '26–30 You cast Gust of Wind'),
            ('effects[3].roll_from', 31, '31–35 Nothing happens at the chosen point of origin. Instead, you take 1d6 Psychic damage'),
            ('effects[3].roll_to', 35, '31–35 Nothing happens at the chosen point of origin. Instead, you take 1d6 Psychic damage'),
            ('effects[3].effect', 'you-take-psychic-damage', '31–35 Nothing happens at the chosen point of origin. Instead, you take 1d6 Psychic damage'),
            ('effects[4].roll_from', 36, '36–40 Heavy rain falls'),
            ('effects[4].roll_to', 40, '36–40 Heavy rain falls'),
            ('effects[4].effect', 'heavy-rain', '36–40 Heavy rain falls'),
            ('effects[5].roll_from', 41, '41–45 A cloud of 600 oversized butterflies'),
            ('effects[5].roll_to', 45, '41–45 A cloud of 600 oversized butterflies'),
            ('effects[5].effect', 'butterfly-cloud', '41–45 A cloud of 600 oversized butterflies'),
            ('effects[6].roll_from', 46, '46–50 You cast Lightning Bolt'),
            ('effects[6].roll_to', 50, '46–50 You cast Lightning Bolt'),
            ('effects[6].effect', 'cast-lightning-bolt', '46–50 You cast Lightning Bolt'),
            ('effects[7].roll_from', 51, '51–55 The creature closest to the chosen point of origin is enlarged'),
            ('effects[7].roll_to', 55, '51–55 The creature closest to the chosen point of origin is enlarged'),
            ('effects[7].effect', 'enlarge-closest-creature', '51–55 The creature closest to the chosen point of origin is enlarged'),
            ('effects[8].roll_from', 56, '56–60 A magically formed creature appears'),
            ('effects[8].roll_to', 60, '56–60 A magically formed creature appears'),
            ('effects[8].effect', 'creature-appears-from-1d4-table', '56–60 A magically formed creature appears'),
            ('effects[9].roll_from', 61, '61–64 Grass covers a 60-foot-radius circle of ground'),
            ('effects[9].roll_to', 64, '61–64 Grass covers a 60-foot-radius circle of ground'),
            ('effects[9].effect', 'grass-grows', '61–64 Grass covers a 60-foot-radius circle of ground'),
            ('effects[10].roll_from', 65, '65–68 An object of the GM’s choice disappears into the Ethereal Plane'),
            ('effects[10].roll_to', 68, '65–68 An object of the GM’s choice disappears into the Ethereal Plane'),
            ('effects[10].effect', 'object-moves-to-ethereal-plane', '65–68 An object of the GM’s choice disappears into the Ethereal Plane'),
            ('effects[11].roll_from', 69, '69–72 Nothing happens at the chosen point of origin. Instead, you shrink'),
            ('effects[11].roll_to', 72, '69–72 Nothing happens at the chosen point of origin. Instead, you shrink'),
            ('effects[11].effect', 'you-shrink', '69–72 Nothing happens at the chosen point of origin. Instead, you shrink'),
            ('effects[12].roll_from', 73, '73–77 Leaves grow from the creature nearest to the chosen point of origin'),
            ('effects[12].roll_to', 77, '73–77 Leaves grow from the creature nearest to the chosen point of origin'),
            ('effects[12].effect', 'leaves-grow-on-nearest-creature', '73–77 Leaves grow from the creature nearest to the chosen point of origin'),
            ('effects[13].roll_from', 78, '78–82 Nothing happens at the chosen point of origin. Instead, a burst of colorful, shimmering light'),
            ('effects[13].roll_to', 82, '78–82 Nothing happens at the chosen point of origin. Instead, a burst of colorful, shimmering light'),
            ('effects[13].effect', 'blinding-light-burst', '78–82 Nothing happens at the chosen point of origin. Instead, a burst of colorful, shimmering light'),
            ('effects[14].roll_from', 83, '83–87 Nothing happens at the chosen point of origin. Instead, you cast Invisibility on yourself'),
            ('effects[14].roll_to', 87, '83–87 Nothing happens at the chosen point of origin. Instead, you cast Invisibility on yourself'),
            ('effects[14].effect', 'you-turn-invisible', '83–87 Nothing happens at the chosen point of origin. Instead, you cast Invisibility on yourself'),
            ('effects[15].roll_from', 88, '88–92 Nothing happens at the chosen point of origin'),
            ('effects[15].roll_to', 92, '88–92 Nothing happens at the chosen point of origin'),
            ('effects[15].effect', 'gem-stream', '88–92 Nothing happens at the chosen point of origin'),
            ('effects[16].roll_from', 93, '93–97 You cast Polymorph'),
            ('effects[16].roll_to', 97, '93–97 You cast Polymorph'),
            ('effects[16].effect', 'polymorph-closest-creature-from-1d4-table', '93–97 You cast Polymorph'),
            ('effects[17].roll_from', 98, '98–00 The creature closest to the chosen point of origin'),
            ('effects[17].roll_to', 100, '98–00 The creature closest to the chosen point of origin'),
            ('effects[17].effect', 'petrify-closest-creature', '98–00 The creature closest to the chosen point of origin'),
            ('spell_table.die_sides', 10, 'Roll 1d10 to determine the spell'),
            ('spell_table.rows[0].roll_from', 1, 'on a 1–2, Darkness'),
            ('spell_table.rows[0].roll_to', 2, 'on a 1–2, Darkness'),
            ('spell_table.rows[0].result', 'darkness', 'on a 1–2, Darkness'),
            ('spell_table.rows[1].roll_from', 3, 'on a 3–4, Faerie Fire'),
            ('spell_table.rows[1].roll_to', 4, 'on a 3–4, Faerie Fire'),
            ('spell_table.rows[1].result', 'faerie-fire', 'on a 3–4, Faerie Fire'),
            ('spell_table.rows[2].roll_from', 5, 'on a 5–6, Fireball'),
            ('spell_table.rows[2].roll_to', 6, 'on a 5–6, Fireball'),
            ('spell_table.rows[2].result', 'fireball', 'on a 5–6, Fireball'),
            ('spell_table.rows[3].roll_from', 7, 'on a 7–8, Slow'),
            ('spell_table.rows[3].roll_to', 8, 'on a 7–8, Slow'),
            ('spell_table.rows[3].result', 'slow', 'on a 7–8, Slow'),
            ('spell_table.rows[4].roll_from', 9, 'on a 9–10, Stinking Cloud'),
            ('spell_table.rows[4].roll_to', 10, 'on a 9–10, Stinking Cloud'),
            ('spell_table.rows[4].result', 'stinking-cloud', 'on a 9–10, Stinking Cloud'),
            ('stunned.condition', 'stunned', 'you have the Stunned condition'),
            ('stunned.until', 'start-of-your-next-turn', 'until the start of your next turn'),
            ('stunned.belief', 'something-awesome-just-happened', 'believing something awesome just happened'),
            ('gust_line', 'from-you-to-point-of-origin', 'You cast Gust of Wind. The Line created by the spell extends from you to the chosen point of origin'),
            ('psychic_damage', {'dice': 1, 'sides': 6, 'type': 'psychic'}, 'you take 1d6 Psychic damage'),
            ('rain.duration_minutes', 1, 'Heavy rain falls for 1 minute'),
            ('rain.height_feet', 120, 'in a 120-foot-high'),
            ('rain.radius_feet', 60, '60-foot-radius Cylinder centered on the chosen point of origin'),
            ('rain.area', 'lightly-obscured', 'the area of effect is Lightly Obscured'),
            ('butterflies.count', 600, 'A cloud of 600 oversized butterflies'),
            ('butterflies.height_feet', 60, 'fills a 60-foot-high'),
            ('butterflies.radius_feet', 30, '30-foot-radius Cylinder'),
            ('butterflies.duration_minutes', 10, 'The butterflies remain for 10 minutes'),
            ('butterflies.area', 'heavily-obscured', 'the area of effect is Heavily Obscured'),
            ('lightning_line', 'from-you-to-point-of-origin', 'You cast Lightning Bolt. The Line created by the spell extends from you to the chosen point of origin'),
            ('enlarge.effect', 'as-enlarge-reduce', 'enlarged as if you had cast Enlarge/ Reduce on it'),
            ('enlarge.fallback', 'you-if-target-cannot-be-affected', 'If the target isn’t you and can’t be affected by that spell, you become the target instead'),
            ('creature.placement', 'unoccupied-space-near-point-of-origin', 'in an unoccupied space as close to the chosen point of origin as possible'),
            ('creature.control', 'not-under-your-control', 'The creature isn’t under your control'),
            ('creature.behavior', 'acts-as-it-normally-would', 'acts as it normally would'),
            ('creature.duration_hours', 1, 'disappears after 1 hour'),
            ('creature.also_disappears_at_hit_points', 0, 'or when it drops to 0 Hit Points'),
            ('creature.table.die_sides', 4, 'Roll 1d4 to determine which creature appears'),
            ('creature.table.rows[0].roll_from', 1, 'On a 1, a Rhinoceros appears'),
            ('creature.table.rows[0].roll_to', 1, 'On a 1, a Rhinoceros appears'),
            ('creature.table.rows[0].result', 'rhinoceros', 'On a 1, a Rhinoceros appears'),
            ('creature.table.rows[1].roll_from', 2, 'on a 2, an Elephant appears'),
            ('creature.table.rows[1].roll_to', 2, 'on a 2, an Elephant appears'),
            ('creature.table.rows[1].result', 'elephant', 'on a 2, an Elephant appears'),
            ('creature.table.rows[2].roll_from', 3, 'on a 3–4, a Rat appears'),
            ('creature.table.rows[2].roll_to', 4, 'on a 3–4, a Rat appears'),
            ('creature.table.rows[2].result', 'rat', 'on a 3–4, a Rat appears'),
            ('grass.radius_feet', 60, 'Grass covers a 60-foot-radius circle of ground'),
            ('grass.center', 'as-close-to-point-of-origin-as-possible', 'with the center of that circle as close to the chosen point of origin as possible'),
            ('grass.size_multiplier', 10, 'grows to ten times its normal size'),
            ('grass.overgrown_minutes', 1, 'remains overgrown for 1 minute'),
            ('ethereal.object', 'gm-choice-not-worn-or-carried', 'The object must be neither worn nor carried'),
            ('ethereal.range_feet', 120, 'within 120 feet of the chosen point of origin'),
            ('ethereal.maximum_dimension_feet', 10, 'no larger than 10 feet in any dimension'),
            ('ethereal.without_object', 'nothing-happens', 'If there are no such objects in range, nothing happens'),
            ('shrink.effect', 'as-enlarge-reduce-on-yourself', 'you shrink as if you had cast Enlarge/ Reduce on yourself'),
            ('shrink.duration_minutes', 1, 'remain in that state for 1 minute'),
            ('leaves.fall_off_after_hours', 24, 'the leaves turn brown and fall off after 24 hours'),
            ('leaves.early_removal', 'picked-off', 'Unless they are picked off'),
            ('light_burst.emanation_feet', 30, 'extends from you in a 30-foot Emanation'),
            ('light_burst.save_ability', 'constitution', 'DC 15 Constitution saving throw or have the Blinded condition'),
            ('light_burst.save_dc', 15, 'DC 15 Constitution saving throw or have the Blinded condition'),
            ('light_burst.condition', 'blinded', 'have the Blinded condition'),
            ('light_burst.duration_minutes', 1, 'the Blinded condition for 1 minute'),
            ('light_burst.repeat_save', 'end-of-each-of-its-turns', 'A creature repeats the save at the end of each of its turns'),
            ('invisibility_target', 'yourself', 'you cast Invisibility on yourself'),
            ('gems.count_dice', {'dice': 1, 'sides': 4, 'type': 'gems'}, 'a stream of 1d4 × 10 gems'),
            ('gems.count_multiplier', 10, '1d4 × 10 gems'),
            ('gems.value_gp_each', 1, 'each worth 1 GP'),
            ('gems.line_length_feet', 30, 'in a Line 30 feet long'),
            ('gems.line_width_feet', 5, 'and 5 feet wide'),
            ('gems.direction', 'toward-point-of-origin', 'and 5 feet wide toward the chosen point of origin'),
            ('gems.damage_per_gem', 1, 'Each gem deals 1 Bludgeoning damage'),
            ('gems.damage_type', 'bludgeoning', '1 Bludgeoning damage'),
            ('gems.damage_division', 'total-divided-equally-in-line', 'the total damage of the gems is divided equally among all creatures in the Line'),
            ('polymorph_table.die_sides', 4, 'Roll 1d4 to determine the target’s new form'),
            ('polymorph_table.rows[0].roll_from', 1, 'On a 1, the new form is a Black Bear'),
            ('polymorph_table.rows[0].roll_to', 1, 'On a 1, the new form is a Black Bear'),
            ('polymorph_table.rows[0].result', 'black-bear', 'On a 1, the new form is a Black Bear'),
            ('polymorph_table.rows[1].roll_from', 2, 'on a 2, the new form is a Giant Wasp'),
            ('polymorph_table.rows[1].roll_to', 2, 'on a 2, the new form is a Giant Wasp'),
            ('polymorph_table.rows[1].result', 'giant-wasp', 'on a 2, the new form is a Giant Wasp'),
            ('polymorph_table.rows[2].roll_from', 3, 'on a 3–4, the new form is a Frog'),
            ('polymorph_table.rows[2].roll_to', 4, 'on a 3–4, the new form is a Frog'),
            ('polymorph_table.rows[2].result', 'frog', 'on a 3–4, the new form is a Frog'),
            ('petrify.save_ability', 'constitution', 'makes a DC 15 Constitution saving throw'),
            ('petrify.save_dc', 15, 'makes a DC 15 Constitution saving throw'),
            ('petrify.first_failure', 'restrained-and-turning-to-stone', 'the creature has the Restrained condition and begins to turn to stone'),
            ('petrify.repeat_save', 'end-of-its-next-turn', 'the creature repeats the save at the end of its next turn'),
            ('petrify.repeat_success', 'effect-ends', 'On a successful save, the effect ends'),
            ('petrify.repeat_failure', 'petrified-instead-of-restrained', 'the creature has the Petrified condition instead of the Restrained condition'),
            ('petrify.petrified_until', 'freed-by-greater-restoration-or-similar-magic', 'until the creature is freed by the Greater Restoration spell or similar magic'),
            ('last_charge.die_sides', 20, 'If you expend the wand’s last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the wand crumbles into dust'),
            ('last_charge.result', 'crumbles-into-dust', 'the wand crumbles into dust and is destroyed'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.bonus', 1, 'regains 1d6 + 1 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
        ],
    },
    'wand-of-the-war-mage-1-2-or-3': {
        'plus-1': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('condition', 'while-holding', 'While holding this wand'),
            ('spell_attack_bonus', 1, 'Uncommon (+1)'),
            ('ignores_half_cover', True, 'you ignore Half Cover when making a spell attack roll'),
        ],
        'plus-2': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('condition', 'while-holding', 'While holding this wand'),
            ('spell_attack_bonus', 2, 'Rare (+2)'),
            ('ignores_half_cover', True, 'you ignore Half Cover when making a spell attack roll'),
        ],
        'plus-3': [
            ('attunement_by_any_of[0]', 'spellcaster', 'Requires Attunement by a Spellcaster'),
            ('condition', 'while-holding', 'While holding this wand'),
            ('spell_attack_bonus', 3, 'Very Rare (+3)'),
            ('ignores_half_cover', True, 'you ignore Half Cover when making a spell attack roll'),
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
        raise SystemExit("usage: verify_srd_items_245.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_245.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 245]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 245 and catalogue["batch"] == "data-magic-items-21"
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
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 250-252")


main()
