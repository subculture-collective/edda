#!/usr/bin/env python3
"""Compare issue #240 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (239, 240, 241)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = ['Rod of Lordly Might Rod,']

# Each family's printed category line.
HEADERS = {
    'ring-of-three-wishes': 'Ring, Legendary',
    'ring-of-warmth': 'Ring, Uncommon (Requires Attunement)',
    'ring-of-water-walking': 'Ring, Uncommon',
    'ring-of-x-ray-vision': 'Ring, Rare (Requires Attunement)',
    'ring-of-the-ram': 'Ring, Rare (Requires Attunement)',
    'robe-of-eyes': 'Wondrous Item, Rare (Requires Attunement)',
    'robe-of-scintillating-colors': 'Wondrous Item, Very Rare (Requires Attunement)',
    'robe-of-stars': 'Wondrous Item, Very Rare (Requires Attunement)',
    'robe-of-useful-items': 'Wondrous Item, Uncommon',
    'robe-of-the-archmagi': 'Wondrous Item, Legendary (Requires Attunement by a Sorcerer, Warlock, or Wizard)',
    'rod-of-absorption': 'Rod, Very Rare (Requires Attunement)',
    'rod-of-alertness': 'Rod, Very Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'ring-of-three-wishes': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'none'}, 'you can expend 1 of its 3 charges'),
    'ring-of-the-ram': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This ring has 3 charges and regains 1d3 expended charges daily at dawn'),
    'robe-of-scintillating-colors': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This robe has 3 charges, and it regains 1d3 expended charges daily at dawn'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'ring-of-three-wishes': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('spell', 'wish', 'to cast Wish from it'),
            ('max_charges', 3, 'expend 1 of its 3 charges'),
            ('charge_cost', 1, 'you can expend 1 of its 3 charges'),
            ('after_last_charge', 'becomes-nonmagical', 'The ring becomes nonmagical when you use the last charge'),
        ],
    },
    'ring-of-warmth': {
        'standard': [
            ('trigger', 'take-cold-damage-while-wearing', 'If you take Cold damage while wearing this ring'),
            ('damage_type', 'cold', 'If you take Cold damage'),
            ('reduction', {'dice': 2, 'sides': 8, 'type': 'damage-reduction'}, 'reduces the damage you take by 2d8'),
            ('protected[0]', 'you', 'you and everything you wear and carry'),
            ('protected[1]', 'everything-you-wear-and-carry', 'everything you wear and carry'),
            ('protected_at_or_below_fahrenheit', 0, 'temperatures of 0 degrees Fahrenheit or lower'),
        ],
    },
    'ring-of-water-walking': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('spell', 'water-walk', 'you cast Water Walk from it'),
            ('target', 'yourself-only', 'targeting only yourself'),
        ],
    },
    'ring-of-x-ray-vision': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('action', 'magic', 'you can take a Magic action to gain X-ray vision'),
            ('range_feet', 30, 'X-ray vision with a range of 30 feet'),
            ('duration_minutes', 1, 'for 1 minute'),
            ('solid_objects', 'appear-transparent-and-pass-light', 'solid objects within that radius appear transparent and don’t prevent light from passing through them'),
            ('penetration.stone_feet', 1, 'penetrate 1 foot of stone'),
            ('penetration.common_metal_inches', 1, '1 inch of common metal'),
            ('penetration.wood_or_dirt_feet', 3, 'up to 3 feet of wood or dirt'),
            ('penetration.blocked_by[0]', 'thicker-substances', 'Thicker substances'),
            ('penetration.blocked_by[1]', 'thin-sheet-of-lead', 'a thin sheet of lead block the vision'),
            ('reuse.trigger', 'use-again-before-long-rest', 'Whenever you use the ring again before taking a Long Rest'),
            ('reuse.save_ability', 'constitution', 'DC 15 Constitution saving throw'),
            ('reuse.save_dc', 15, 'DC 15 Constitution saving throw'),
            ('reuse.failure_condition', 'exhaustion', 'or gain 1 Exhaustion level'),
            ('reuse.exhaustion_levels', 1, 'gain 1 Exhaustion level'),
        ],
    },
    'ring-of-the-ram': {
        'standard': [
            ('max_charges', 3, 'This ring has 3 charges'),
            ('recharge', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'regains 1d3 expended charges daily at dawn'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('action', 'magic', 'you can take a Magic action to expend 1 to 3 charges'),
            ('minimum_charges', 1, 'expend 1 to 3 charges to make a ranged spell attack'),
            ('maximum_charges', 3, 'expend 1 to 3 charges to make a ranged spell attack'),
            ('attack.kind', 'ranged-spell-attack', 'to make a ranged spell attack against one creature'),
            ('attack.targets', 1, 'against one creature'),
            ('attack.range_feet', 60, 'one creature you can see within 60 feet of yourself'),
            ('attack.form', 'spectral-rams-head', 'The ring produces a spectral ram’s head'),
            ('attack.attack_bonus', 7, 'makes its attack roll with a +7 bonus'),
            ('attack.damage_per_charge', {'dice': 2, 'sides': 10, 'type': 'force'}, 'for each charge you spend, the target takes 2d10 Force damage'),
            ('attack.push_feet_per_charge', 5, 'is pushed 5 feet away from you'),
            ('break_object.action', 'magic', 'charges as a Magic action to try to break a nonmagical object'),
            ('break_object.target[0]', 'nonmagical-object', 'to try to break a nonmagical object'),
            ('break_object.target[1]', 'not-worn-or-carried', 'that isn’t being worn or carried'),
            ('break_object.range_feet', 60, 'object you can see within 60 feet'),
            ('break_object.check_ability', 'strength', 'The ring makes a Strength check'),
            ('break_object.check_bonus_per_charge', 5, 'with a +5 bonus for each charge you spend'),
        ],
    },
    'robe-of-eyes': {
        'standard': [
            ('appearance', 'eyelike-patterns', 'adorned with eyelike patterns'),
            ('all_around_vision.check_ability', 'wisdom', 'Advantage on Wisdom (Perception) checks'),
            ('all_around_vision.check_skill', 'perception', 'Wisdom (Perception) checks'),
            ('all_around_vision.check_mode', 'advantage', 'gives you Advantage'),
            ('all_around_vision.check_requirement', 'relies-on-sight', 'checks that rely on sight'),
            ('senses[0]', 'darkvision', 'You have Darkvision and Truesight'),
            ('senses[1]', 'truesight', 'You have Darkvision and Truesight'),
            ('sense_range_feet', 120, 'both with a range of 120 feet'),
            ('drawbacks.condition', 'blinded', 'gives you the Blinded condition'),
            ('drawbacks.duration_minutes', 1, 'the Blinded condition for 1 minute'),
            ('drawbacks.triggers[0].spell', 'light', 'A Light spell cast on the robe'),
            ('drawbacks.triggers[0].placement', 'cast-on-the-robe', 'cast on the robe'),
            ('drawbacks.triggers[0].save_dc', 11, 'DC 11 for Light'),
            ('drawbacks.triggers[1].spell', 'daylight', 'Daylight spell cast within 5 feet of the robe'),
            ('drawbacks.triggers[1].placement', 'cast-within-5-feet-of-the-robe', 'within 5 feet of the robe'),
            ('drawbacks.triggers[1].save_dc', 15, 'DC 15 for Daylight'),
            ('drawbacks.save_ability', 'constitution', 'you make a Constitution saving throw'),
            ('drawbacks.save_when', 'end-of-each-of-your-turns', 'At the end of each of your turns'),
            ('drawbacks.success', 'ends-condition-on-yourself', 'ending the condition on yourself on a success'),
        ],
    },
    'robe-of-scintillating-colors': {
        'standard': [
            ('max_charges', 3, 'This robe has 3 charges'),
            ('recharge', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'it regains 1d3 expended charges daily at dawn'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('condition', 'while-wearing', 'While you wear it'),
            ('action', 'magic', 'you can take a Magic action and expend 1 charge'),
            ('charge_cost', 1, 'expend 1 charge to cause the garment'),
            ('effect', 'shifting-pattern-of-dazzling-hues', 'display a shifting pattern of dazzling hues'),
            ('duration', 'until-end-of-your-next-turn', 'until the end of your next turn'),
            ('bright_light_radius_feet', 30, 'Bright Light in a 30-foot radius'),
            ('dim_light_additional_feet', 30, 'Dim Light for an additional 30 feet'),
            ('attacks_against_you', 'disadvantage', 'creatures that can see you have Disadvantage on attack rolls against you'),
            ('stun.affects', 'creature-in-bright-light-that-can-see-you-on-activation', 'Any creature in the Bright Light that can see you when the robe’s power is activated'),
            ('stun.save_ability', 'wisdom', 'DC 15 Wisdom saving throw'),
            ('stun.save_dc', 15, 'DC 15 Wisdom saving throw'),
            ('stun.condition', 'stunned', 'or have the Stunned condition'),
            ('stun.condition_ends', 'when-effect-ends', 'until the effect ends'),
        ],
    },
    'robe-of-stars': {
        'standard': [
            ('appearance', 'black-or-dark-blue-robe-with-small-white-or-silver-stars', 'black or dark-blue robe is embroidered with small white or silver stars'),
            ('saving_throw_bonus', 1, '+1 bonus to saving throws while you wear it'),
            ('missiles.large_stars', 6, 'Six stars, located on the robe’s upper-front portion'),
            ('missiles.action', 'magic', 'you can take a Magic action to remove one of the stars'),
            ('missiles.spell', 'magic-missile', 'version of Magic Missile'),
            ('missiles.spell_level', 5, 'cast the level 5 version'),
            ('missiles.stars_return', {'dice': 1, 'sides': 6, 'type': 'stars'}, '1d6 removed stars reappear'),
            ('missiles.return_when', 'daily-at-dusk', 'Daily at dusk'),
            ('astral_travel.action', 'magic', 'you can take a Magic action to enter the Astral Plane'),
            ('astral_travel.destination', 'astral-plane', 'enter the Astral Plane'),
            ('astral_travel.takes', 'everything-worn-and-carried', 'along with everything you are wearing and carrying'),
            ('astral_travel.remain_until', 'magic-action-to-return', 'until you take a Magic action to return to the plane you were on'),
            ('astral_travel.reappear_at', 'last-space-or-nearest-unoccupied-space', 'the last space you occupied or, if that space is occupied, the nearest unoccupied space'),
        ],
    },
    'robe-of-useful-items': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing the robe'),
            ('action', 'magic', 'you can take a Magic action to detach one of the patches'),
            ('patch_becomes', 'object-or-creature-it-represents', 'causing it to become the object or creature it represents'),
            ('after_last_patch', 'ordinary-garment', 'the robe becomes an ordinary garment'),
            ('copies_of_each_fixed_patch', 2, 'The robe has two of each of the following patches'),
            ('fixed_patches[0]', 'bullseye-lantern-filled-and-lit', '• Bullseye Lantern (filled and lit)'),
            ('fixed_patches[1]', 'dagger', '• Dagger'),
            ('fixed_patches[2]', 'mirror', '• Mirror'),
            ('fixed_patches[3]', 'pole', '• Pole'),
            ('fixed_patches[4]', 'rope-coiled', '• Rope (coiled)'),
            ('fixed_patches[5]', 'sack', '• Sack'),
            ('other_patches', {'dice': 4, 'sides': 4, 'type': 'patches'}, 'the robe has 4d4 other patches'),
            ('other_patches_chosen_by', 'gm-chooses-or-rolls-randomly', 'The GM chooses the patches or determines them randomly'),
            ('table_die_sides', 100, '1d100 Patch'),
            ('patch_table[0].min', 1, '01–08 Bag of 100 GP'),
            ('patch_table[0].max', 8, '01–08 Bag of 100 GP'),
            ('patch_table[0].patch', 'bag-of-100-gp', '01–08 Bag of 100 GP'),
            ('patch_table[1].min', 9, '09–15 Silver coffer'),
            ('patch_table[1].max', 15, '09–15 Silver coffer'),
            ('patch_table[1].patch', 'silver-coffer', '09–15 Silver coffer'),
            ('patch_table[2].min', 16, '16–22 Iron door'),
            ('patch_table[2].max', 22, '16–22 Iron door'),
            ('patch_table[2].patch', 'iron-door', '16–22 Iron door'),
            ('patch_table[3].min', 23, '23–30 10 gems worth 100 GP each'),
            ('patch_table[3].max', 30, '23–30 10 gems worth 100 GP each'),
            ('patch_table[3].patch', 'gems-worth-100-gp-each', '23–30 10 gems worth 100 GP each'),
            ('patch_table[4].min', 31, '31–44 Wooden ladder'),
            ('patch_table[4].max', 44, '31–44 Wooden ladder'),
            ('patch_table[4].patch', 'wooden-ladder', '31–44 Wooden ladder'),
            ('patch_table[5].min', 45, '45–51 Riding Horse with a Riding Saddle'),
            ('patch_table[5].max', 51, '45–51 Riding Horse with a Riding Saddle'),
            ('patch_table[5].patch', 'riding-horse-with-riding-saddle', '45–51 Riding Horse with a Riding Saddle'),
            ('patch_table[6].min', 52, '52–59 Open pit'),
            ('patch_table[6].max', 59, '52–59 Open pit'),
            ('patch_table[6].patch', 'open-pit', '52–59 Open pit'),
            ('patch_table[7].min', 60, '60–68 4 Potions of Healing'),
            ('patch_table[7].max', 68, '60–68 4 Potions of Healing'),
            ('patch_table[7].patch', 'potions-of-healing', '60–68 4 Potions of Healing'),
            ('patch_table[8].min', 69, '69–75 Rowboat'),
            ('patch_table[8].max', 75, '69–75 Rowboat'),
            ('patch_table[8].patch', 'rowboat', '69–75 Rowboat'),
            ('patch_table[9].min', 76, '76–83 Spell Scroll containing one spell of level 1, 2, or'),
            ('patch_table[9].max', 83, '76–83 Spell Scroll containing one spell of level 1, 2, or'),
            ('patch_table[9].patch', 'spell-scroll-of-level-1-2-or-3', '76–83 Spell Scroll containing one spell of level 1, 2, or'),
            ('patch_table[10].min', 84, '84–90 2 Mastiffs'),
            ('patch_table[10].max', 90, '84–90 2 Mastiffs'),
            ('patch_table[10].patch', 'mastiffs', '84–90 2 Mastiffs'),
            ('patch_table[11].min', 91, '91–96 Window'),
            ('patch_table[11].max', 96, '91–96 Window'),
            ('patch_table[11].patch', 'window', '91–96 Window'),
            ('patch_table[12].min', 97, '97–00 Portable Ram'),
            ('patch_table[12].max', 100, '97–00 Portable Ram'),
            ('patch_table[12].patch', 'portable-ram', '97–00 Portable Ram'),
            ('patch_facts[0].patch', 'gems-worth-100-gp-each', '10 gems worth 100 GP each'),
            ('patch_facts[0].key', 'count', '10 gems worth 100 GP each'),
            ('patch_facts[0].amount', 10, '10 gems worth 100 GP each'),
            ('patch_facts[0].unit', 'gems', '10 gems worth 100 GP each'),
            ('patch_facts[1].patch', 'potions-of-healing', '4 Potions of Healing'),
            ('patch_facts[1].key', 'count', '4 Potions of Healing'),
            ('patch_facts[1].amount', 4, '4 Potions of Healing'),
            ('patch_facts[1].unit', 'potions', '4 Potions of Healing'),
            ('patch_facts[2].patch', 'mastiffs', '2 Mastiffs'),
            ('patch_facts[2].key', 'count', '2 Mastiffs'),
            ('patch_facts[2].amount', 2, '2 Mastiffs'),
            ('patch_facts[2].unit', 'mastiffs', '2 Mastiffs'),
            ('patch_facts[3].patch', 'spell-scroll-of-level-1-2-or-3', 'one spell of level 1, 2, or'),
            ('patch_facts[3].key', 'minimum-spell-level', 'one spell of level 1, 2, or'),
            ('patch_facts[3].amount', 1, 'one spell of level 1, 2, or'),
            ('patch_facts[3].unit', 'spell-level', 'one spell of level 1, 2, or'),
            ('patch_facts[4].patch', 'spell-scroll-of-level-1-2-or-3', '3 (your choice)'),
            ('patch_facts[4].key', 'maximum-spell-level', '3 (your choice)'),
            ('patch_facts[4].amount', 3, '3 (your choice)'),
            ('patch_facts[4].unit', 'spell-level', '3 (your choice)'),
            ('patch_facts[5].patch', 'bag-of-100-gp', 'Bag of 100 GP'),
            ('patch_facts[5].key', 'contents-value', 'Bag of 100 GP'),
            ('patch_facts[5].amount', 100, 'Bag of 100 GP'),
            ('patch_facts[5].unit', 'gp', 'Bag of 100 GP'),
            ('patch_facts[6].patch', 'silver-coffer', 'Silver coffer (1 foot long'),
            ('patch_facts[6].key', 'length', 'Silver coffer (1 foot long'),
            ('patch_facts[6].amount', 1, 'Silver coffer (1 foot long'),
            ('patch_facts[6].unit', 'feet', 'Silver coffer (1 foot long'),
            ('patch_facts[7].patch', 'silver-coffer', '6 inches wide and deep'),
            ('patch_facts[7].key', 'width-and-depth', '6 inches wide and deep'),
            ('patch_facts[7].amount', 6, '6 inches wide and deep'),
            ('patch_facts[7].unit', 'inches', '6 inches wide and deep'),
            ('patch_facts[8].patch', 'silver-coffer', 'deep) worth 500 GP'),
            ('patch_facts[8].key', 'value', 'deep) worth 500 GP'),
            ('patch_facts[8].amount', 500, 'deep) worth 500 GP'),
            ('patch_facts[8].unit', 'gp', 'deep) worth 500 GP'),
            ('patch_facts[9].patch', 'iron-door', 'up to 10 feet wide'),
            ('patch_facts[9].key', 'maximum-width', 'up to 10 feet wide'),
            ('patch_facts[9].amount', 10, 'up to 10 feet wide'),
            ('patch_facts[9].unit', 'feet', 'up to 10 feet wide'),
            ('patch_facts[10].patch', 'iron-door', 'and 10 feet high'),
            ('patch_facts[10].key', 'maximum-height', 'and 10 feet high'),
            ('patch_facts[10].amount', 10, 'and 10 feet high'),
            ('patch_facts[10].unit', 'feet', 'and 10 feet high'),
            ('patch_facts[11].patch', 'gems-worth-100-gp-each', 'gems worth 100 GP each'),
            ('patch_facts[11].key', 'value-each', 'gems worth 100 GP each'),
            ('patch_facts[11].amount', 100, 'gems worth 100 GP each'),
            ('patch_facts[11].unit', 'gp', 'gems worth 100 GP each'),
            ('patch_facts[12].patch', 'wooden-ladder', 'ladder (24 feet long)'),
            ('patch_facts[12].key', 'length', 'ladder (24 feet long)'),
            ('patch_facts[12].amount', 24, 'ladder (24 feet long)'),
            ('patch_facts[12].unit', 'feet', 'ladder (24 feet long)'),
            ('patch_facts[13].patch', 'open-pit', 'Open pit (a 10-foot Cube)'),
            ('patch_facts[13].key', 'cube-size', 'Open pit (a 10-foot Cube)'),
            ('patch_facts[13].amount', 10, 'Open pit (a 10-foot Cube)'),
            ('patch_facts[13].unit', 'feet', 'Open pit (a 10-foot Cube)'),
            ('patch_facts[14].patch', 'open-pit', 'place on the ground within 10 feet of yourself'),
            ('patch_facts[14].key', 'placement-range', 'place on the ground within 10 feet of yourself'),
            ('patch_facts[14].amount', 10, 'place on the ground within 10 feet of yourself'),
            ('patch_facts[14].unit', 'feet', 'place on the ground within 10 feet of yourself'),
            ('patch_facts[15].patch', 'rowboat', 'Rowboat (12 feet long)'),
            ('patch_facts[15].key', 'length', 'Rowboat (12 feet long)'),
            ('patch_facts[15].amount', 12, 'Rowboat (12 feet long)'),
            ('patch_facts[15].unit', 'feet', 'Rowboat (12 feet long)'),
            ('patch_facts[16].patch', 'window', 'Window (2 feet by 4 feet'),
            ('patch_facts[16].key', 'width', 'Window (2 feet by 4 feet'),
            ('patch_facts[16].amount', 2, 'Window (2 feet by 4 feet'),
            ('patch_facts[16].unit', 'feet', 'Window (2 feet by 4 feet'),
            ('patch_facts[17].patch', 'window', '2 feet by 4 feet'),
            ('patch_facts[17].key', 'height', '2 feet by 4 feet'),
            ('patch_facts[17].amount', 4, '2 feet by 4 feet'),
            ('patch_facts[17].unit', 'feet', '2 feet by 4 feet'),
            ('patch_facts[18].patch', 'window', 'up to 2 feet deep'),
            ('patch_facts[18].key', 'maximum-depth', 'up to 2 feet deep'),
            ('patch_facts[18].amount', 2, 'up to 2 feet deep'),
            ('patch_facts[18].unit', 'feet', 'up to 2 feet deep'),
            ('patch_rules[0].patch', 'iron-door', 'barred on one side of your choice'),
            ('patch_rules[0].rule', 'barred-on-one-side-of-your-choice', 'barred on one side of your choice'),
            ('patch_rules[1].patch', 'iron-door', 'which you can place in an opening you can reach'),
            ('patch_rules[1].rule', 'placed-in-an-opening-you-can-reach', 'which you can place in an opening you can reach'),
            ('patch_rules[2].patch', 'iron-door', 'it conforms to fit the opening, attaching and hinging itself'),
            ('patch_rules[2].rule', 'conforms-attaches-and-hinges-itself', 'it conforms to fit the opening, attaching and hinging itself'),
            ('patch_rules[3].patch', 'window', 'which you can place on a vertical surface you can reach'),
            ('patch_rules[3].rule', 'placed-on-a-vertical-surface-you-can-reach', 'which you can place on a vertical surface you can reach'),
        ],
    },
    'robe-of-the-archmagi': {
        'standard': [
            ('attunement_by_any_of[0]', 'sorcerer', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[1]', 'warlock', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[2]', 'wizard', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('armor.requirement', 'not-wearing-armor', 'If you aren’t wearing armor'),
            ('armor.base_armor_class', 15, 'your base Armor Class is 15 plus your Dexterity modifier'),
            ('armor.adds_ability_modifier', 'dexterity', 'plus your Dexterity modifier'),
            ('magic_resistance.save_mode', 'advantage', 'You have Advantage on saving throws'),
            ('magic_resistance.against[0]', 'spells', 'against spells and other magical effects'),
            ('magic_resistance.against[1]', 'other-magical-effects', 'against spells and other magical effects'),
            ('war_mage.increases[0]', 'spell-save-dc', 'Your spell save DC and spell attack bonus each increase'),
            ('war_mage.increases[1]', 'spell-attack-bonus', 'Your spell save DC and spell attack bonus each increase'),
            ('war_mage.bonus', 2, 'each increase by 2'),
        ],
    },
    'rod-of-absorption': {
        'standard': [
            ('condition', 'while-holding', 'While holding this rod'),
            ('absorb.action', 'reaction', 'you can take a Reaction to absorb a spell'),
            ('absorb.requires[0]', 'spell-targets-only-you', 'a spell that is targeting only you'),
            ('absorb.requires[1]', 'spell-creates-no-area-of-effect', 'and doesn’t create an area of effect'),
            ('absorb.result', 'spell-canceled-and-energy-stored', 'The absorbed spell’s effect is canceled, and the spell’s energy'),
            ('absorb.energy_level', 'level-of-spell-as-cast', 'The energy has the same level as the spell when it was cast'),
            ('absorb.casting_resources', 'wasted', 'any resources used to cast it are wasted'),
            ('absorb.if_it_cannot_store', 'rod-has-no-effect', 'the rod has no effect on that spell'),
            ('lifetime_capacity_levels', 50, 'up to 50 levels of energy over the course of its existence'),
            ('known_on_attunement[0]', 'lifetime-levels-absorbed', 'you know how many levels of energy the rod has absorbed over the course of its existence and how many levels of spell energy it currently has stored'),
            ('known_on_attunement[1]', 'levels-currently-stored', 'you know how many levels of energy the rod has absorbed over the course of its existence and how many levels of spell energy it currently has stored'),
            ('convert.requirement', 'spellcaster-holding-rod', 'If you are a spellcaster holding the rod'),
            ('convert.result', 'spell-slots-for-prepared-or-known-spells', 'convert energy stored in it into spell slots to cast spells you have prepared or know'),
            ('convert.level_limit', 'equal-to-or-lower-than-your-own-slots', 'only of a level equal to or lower than your own spell slots'),
            ('convert.maximum_slot_level', 5, 'up to a maximum of level 5'),
            ('convert.casting', 'stored-levels-replace-slots-otherwise-cast-as-normal', 'You use the stored levels in place of your slots but otherwise cast the spell as normal'),
            ('found_stored_levels', {'dice': 1, 'sides': 10, 'type': 'spell-energy-levels'}, 'A newly found rod typically has 1d10 levels of spell energy'),
            ('nonmagical_when[0]', 'cannot-absorb-more', 'A rod that can no longer absorb spell energy and has no energy remaining becomes nonmagical'),
            ('nonmagical_when[1]', 'no-energy-remaining', 'A rod that can no longer absorb spell energy and has no energy remaining becomes nonmagical'),
        ],
    },
    'rod-of-alertness': {
        'standard': [
            ('condition', 'while-holding', 'While holding the rod'),
            ('advantage_on[0]', 'wisdom-perception-checks', 'Advantage on Wisdom (Perception) checks'),
            ('advantage_on[1]', 'initiative-rolls', 'and on Initiative rolls'),
            ('spells[0]', 'detect-evil-and-good', '• Detect Evil and Good'),
            ('spells[1]', 'detect-magic', '• Detect Magic'),
            ('spells[2]', 'detect-poison-and-disease', '• Detect Poison and Disease'),
            ('spells[3]', 'see-invisibility', '• See Invisibility'),
            ('protective_aura.action', 'magic', 'As a Magic action, you can plant the haft end of the rod'),
            ('protective_aura.bright_light_radius_feet', 60, 'Bright Light in a 60-foot radius'),
            ('protective_aura.dim_light_additional_feet', 60, 'Dim Light for an additional 60 feet'),
            ('protective_aura.beneficiaries[0]', 'you', 'you and your allies gain'),
            ('protective_aura.beneficiaries[1]', 'your-allies', 'you and your allies gain'),
            ('protective_aura.benefit_requirement', 'while-in-the-bright-light', 'While in that Bright Light'),
            ('protective_aura.armor_class_bonus', 1, '+1 bonus to Armor Class and saving throws'),
            ('protective_aura.saving_throw_bonus', 1, '+1 bonus to Armor Class and saving throws'),
            ('protective_aura.senses', 'location-of-invisible-creatures-in-the-bright-light', 'sense the location of any Invisible creature that is also in the Bright Light'),
            ('protective_aura.duration_minutes', 10, 'the effect ends after 10 minutes'),
            ('protective_aura.ends_early_when', 'creature-takes-magic-action-to-pull-rod', 'or when a creature takes a Magic action to pull the rod from the ground'),
            ('protective_aura.reuse_after', 'next-dawn', 'this property can’t be used again until the next dawn'),
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
        raise SystemExit("usage: verify_srd_items_240.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_240.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 240]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 240 and catalogue["batch"] == "data-magic-items-16"
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
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 239-241")


main()
