#!/usr/bin/env python3
"""Compare issue #238 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (236, 237, 238)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = ['Ring of Free Action Ring']

# Each family's printed category line.
HEADERS = {
    'potion-of-poison': 'Potion, Uncommon',
    'potion-of-resistance': 'Potion, Uncommon',
    'potion-of-speed': 'Potion, Very Rare',
    'potion-of-vitality': 'Potion, Very Rare',
    'potion-of-water-breathing': 'Potion, Uncommon',
    'quarterstaff-of-the-acrobat': 'Weapon (Quarterstaff), Very Rare (Requires Attunement)',
    'ring-of-animal-influence': 'Ring, Rare',
    'ring-of-djinni-summoning': 'Ring, Legendary (Requires Attunement)',
    'ring-of-elemental-command': 'Ring, Legendary (Requires Attunement)',
    'ring-of-evasion': 'Ring, Rare (Requires Attunement)',
    'ring-of-feather-falling': 'Ring, Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'ring-of-animal-influence': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This ring has 3 charges, and it regains 1d3 expended charges daily at dawn'),
    'ring-of-elemental-command': ({'mode': 'limited', 'maximum': 5, 'per_use': 1, 'recharge': 'regain-1d4-plus-1-at-dawn'}, 'The ring has 5 charges and regains 1d4 + 1 expended charges daily at dawn'),
    'ring-of-evasion': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This ring has 3 charges, and it regains 1d3 expended charges daily at dawn'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'potion-of-poison': {
        'standard': [
            ('disguise', 'looks-smells-and-tastes-like-healing-or-other-beneficial-potion', 'looks, smells, and tastes like a Potion of Healing or another beneficial potion'),
            ('masked_by', 'illusion-magic', 'poison masked by illusion magic'),
            ('revealed_by', 'identify', 'Identify reveals its true nature'),
            ('damage', {'dice': 4, 'sides': 6, 'type': 'poison'}, 'you take 4d6 Poison damage'),
            ('save_ability', 'constitution', 'DC 13 Constitution saving throw'),
            ('save_dc', 13, 'DC 13 Constitution saving throw'),
            ('condition', 'poisoned', 'or have the Poisoned condition'),
            ('duration_hours', 1, 'the Poisoned condition for 1 hour'),
        ],
    },
    'potion-of-resistance': {
        'standard': [
            ('defense', 'resistance', 'you have Resistance to one type of damage'),
            ('duration_hours', 1, 'one type of damage for 1 hour'),
            ('type_chosen_by', 'gm-chooses-or-rolls-randomly', 'The GM chooses the type or determines it randomly'),
            ('table_die_sides', 10, '1d10 Damage Type'),
            ('damage_types[0].roll', 1, '1 Acid'),
            ('damage_types[0].damage_type', 'acid', '1 Acid'),
            ('damage_types[1].roll', 2, '2 Cold'),
            ('damage_types[1].damage_type', 'cold', '2 Cold'),
            ('damage_types[2].roll', 3, '3 Fire'),
            ('damage_types[2].damage_type', 'fire', '3 Fire'),
            ('damage_types[3].roll', 4, '4 Force'),
            ('damage_types[3].damage_type', 'force', '4 Force'),
            ('damage_types[4].roll', 5, '5 Lightning'),
            ('damage_types[4].damage_type', 'lightning', '5 Lightning'),
            ('damage_types[5].roll', 6, '6 Necrotic'),
            ('damage_types[5].damage_type', 'necrotic', '6 Necrotic'),
            ('damage_types[6].roll', 7, '7 Poison'),
            ('damage_types[6].damage_type', 'poison', '7 Poison'),
            ('damage_types[7].roll', 8, '8 Psychic'),
            ('damage_types[7].damage_type', 'psychic', '8 Psychic'),
            ('damage_types[8].roll', 9, '9 Radiant'),
            ('damage_types[8].damage_type', 'radiant', '9 Radiant'),
            ('damage_types[9].roll', 10, '10 Thunder'),
            ('damage_types[9].damage_type', 'thunder', '10 Thunder'),
        ],
    },
    'potion-of-speed': {
        'standard': [
            ('spell', 'haste', 'the effect of the Haste spell'),
            ('duration_minutes', 1, 'the Haste spell for 1 minute'),
            ('concentration_required', False, '(no Concentration required)'),
            ('lethargy_when_it_ends', False, 'without suffering the wave of lethargy that typically occurs when the effect ends'),
            ('appearance', 'yellow-fluid-streaked-with-black-that-swirls-on-its-own', 'yellow fluid is streaked with black and swirls on its own'),
        ],
    },
    'potion-of-vitality': {
        'standard': [
            ('removes', 'all-exhaustion-levels', 'removes any Exhaustion levels you have'),
            ('ends_condition', 'poisoned', 'ends the Poisoned condition on you'),
            ('hit_die_healing', 'maximum-for-any-hit-point-die-spent', 'you regain the maximum number of Hit Points for any Hit Point Die you spend'),
            ('duration_hours', 24, 'For the next 24 hours'),
            ('appearance', 'crimson-liquid-pulses-with-dull-light-like-a-heartbeat', 'crimson liquid regularly pulses with dull light, calling to mind a heartbeat'),
        ],
    },
    'potion-of-water-breathing': {
        'standard': [
            ('grants', 'breathe-underwater', 'You can breathe underwater'),
            ('duration_hours', 24, 'underwater for 24 hours after drinking this potion'),
            ('appearance', 'cloudy-green-fluid-smelling-of-the-sea-with-jellyfish-like-bubble', 'cloudy green fluid smells of the sea and has a jellyfish-like bubble floating in it'),
        ],
    },
    'quarterstaff-of-the-acrobat': {
        'standard': [
            ('attack_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('light.condition', 'while-holding', 'While holding this weapon, you can cause it to emit'),
            ('light.color', 'green', 'emit green Dim Light'),
            ('light.dim_radius_feet', 10, 'green Dim Light out to 10 feet'),
            ('light.activation[0]', 'bonus-action', 'either as a Bonus Action'),
            ('light.activation[1]', 'after-rolling-initiative', 'or after you roll Initiative'),
            ('light.extinguish_action', 'bonus', 'or you can extinguish the light as a Bonus Action'),
            ('alter_form.action', 'bonus', 'you can take a Bonus Action to alter its form'),
            ('alter_form.forms[0]', 'six-inch-rod', 'turning it into a 6-inch rod'),
            ('alter_form.forms[1]', 'ten-foot-pole', 'or a 10-foot pole'),
            ('alter_form.forms[2]', 'quarterstaff', 'or reverting it a Quarterstaff'),
            ('alter_form.rod_inches', 6, 'a 6-inch rod'),
            ('alter_form.pole_feet', 10, 'a 10-foot pole'),
            ('alter_form.elongation_limit', 'surrounding-space', 'the weapon will elongate only as far as the surrounding space allows'),
            ('acrobatic_assist.forms[0]', 'quarterstaff', 'Acrobatic Assist (Quarterstaff and 10-Foot Pole Forms Only)'),
            ('acrobatic_assist.forms[1]', 'ten-foot-pole', 'Acrobatic Assist (Quarterstaff and 10-Foot Pole Forms Only)'),
            ('acrobatic_assist.check_ability', 'dexterity', 'Advantage on Dexterity (Acrobatics) checks'),
            ('acrobatic_assist.check_skill', 'acrobatics', 'Dexterity (Acrobatics) checks'),
            ('acrobatic_assist.check_mode', 'advantage', 'you have Advantage on Dexterity'),
            ('attack_deflection.form', 'quarterstaff', 'Attack Deflection (Quarterstaff Form Only)'),
            ('attack_deflection.trigger', 'hit-by-attack-while-holding', 'When you are hit by an attack while holding the weapon'),
            ('attack_deflection.action', 'reaction', 'you can take a Reaction to twirl the weapon around you'),
            ('attack_deflection.armor_class_bonus', 5, '+5 bonus to your Armor Class against the triggering attack'),
            ('attack_deflection.possible_result', 'attack-may-miss-you', 'potentially causing the attack to miss you'),
            ('attack_deflection.reuse_after', 'short-or-long-rest', 'until you finish a Short or Long Rest'),
            ('ranged_weapon.form', 'quarterstaff', 'Ranged Weapon (Quarterstaff Form Only)'),
            ('ranged_weapon.property', 'thrown', 'has the Thrown property'),
            ('ranged_weapon.normal_range_feet', 30, 'a normal range of 30 feet'),
            ('ranged_weapon.long_range_feet', 120, 'a long range of 120 feet'),
            ('ranged_weapon.returns_to_hand', 'immediately-after-ranged-attack', 'Immediately after you make a ranged attack with the weapon, it flies back to your hand'),
        ],
    },
    'ring-of-animal-influence': {
        'standard': [
            ('max_charges', 3, 'This ring has 3 charges'),
            ('recharge', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'it regains 1d3 expended charges daily at dawn'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('condition', 'while-wearing', 'While wearing the ring'),
            ('charge_cost', 1, 'you can expend 1 charge to cast one of the following spells'),
            ('save_dc', 13, '(save DC 13)'),
            ('spells[0]', 'animal-friendship', '• Animal Friendship'),
            ('spells[1]', 'fear', '• Fear (affects Beasts only)'),
            ('spells[2]', 'speak-with-animals', '• Speak with Animals'),
            ('fear_affects_only', 'beasts', 'Fear (affects Beasts only)'),
        ],
    },
    'ring-of-djinni-summoning': {
        'standard': [
            ('condition', 'while-wearing', 'While wearing this ring'),
            ('action', 'magic', 'you can take a Magic action to summon a particular Djinni'),
            ('creature', 'djinni', 'summon a particular Djinni'),
            ('origin', 'elemental-plane-of-air', 'from the Elemental Plane of Air'),
            ('summon_range_feet', 120, 'an unoccupied space you choose within 120 feet of yourself'),
            ('concentration_required', True, 'It remains as long as you maintain Concentration'),
            ('max_duration_hours', 1, 'to a maximum of 1 hour'),
            ('also_ends_at', 'zero-hit-points', 'or until it drops to 0 Hit Points'),
            ('attitude', 'friendly-to-you-and-allies', 'the djinni is Friendly to you and your allies'),
            ('obeys_commands', True, 'it obeys your commands'),
            ('without_command', 'defends-itself-only', 'the djinni defends itself against attackers but takes no other actions'),
            ('cooldown_hours', 24, 'it can’t be summoned again for 24 hours'),
            ('nonmagical_when', 'djinni-dies', 'the ring becomes nonmagical if the djinni dies'),
            ('often_created_by', 'djinn-as-gifts', 'often created by the djinn they summon and given to mortals as gifts'),
        ],
    },
    'ring-of-elemental-command': {
        'air': [
            ('linked_plane', 'air', 'Air. You know'),
            ('plane_chosen_by', 'gm-chooses-or-rolls-randomly', 'The GM chooses or randomly determines the linked plane'),
            ('bane.target_creature_type', 'elemental', 'Advantage on attack rolls against Elementals'),
            ('bane.your_attacks', 'advantage', 'you have Advantage on attack rolls against Elementals'),
            ('bane.attacks_against_you', 'disadvantage', 'they have Disadvantage on attack rolls against you'),
            ('compulsion.action', 'magic', 'you can take a Magic action to try to compel an Elemental'),
            ('compulsion.range_feet', 60, 'an Elemental you see within 60 feet of yourself'),
            ('compulsion.save_ability', 'wisdom', 'DC 18 Wisdom saving throw'),
            ('compulsion.save_dc', 18, 'DC 18 Wisdom saving throw'),
            ('compulsion.condition', 'charmed', 'the Elemental has the Charmed condition'),
            ('compulsion.condition_ends', 'start-of-your-next-turn', 'until the start your next turn'),
            ('compulsion.control', 'you-determine-its-move-and-action-on-its-next-turn', 'you determine what it does with its move and action on its next turn'),
            ('focus.language', 'auran', 'You know Auran'),
            ('focus.defense', 'resistance', 'you have Resistance to Lightning damage'),
            ('focus.damage_type', 'lightning', 'you have Resistance to Lightning damage'),
            ('focus.speed_kind', 'fly', 'you have a Fly Speed equal to your Speed and can hover'),
            ('focus.speed', 'equal-to-speed', 'you have a Fly Speed equal to your Speed and can hover'),
            ('focus.other_benefits[0]', 'can-hover', 'and can hover'),
            ('spell_save_dc', 18, 'which has a save DC of 18'),
            ('spells[0].spell', 'chain-lightning', 'Chain Lightning (3 charges)'),
            ('spells[0].charge_cost', 3, 'Chain Lightning (3 charges)'),
            ('spells[1].spell', 'feather-fall', 'Feather Fall (0 charges)'),
            ('spells[1].charge_cost', 0, 'Feather Fall (0 charges)'),
            ('spells[2].spell', 'gust-of-wind', 'Gust of Wind (2 charges)'),
            ('spells[2].charge_cost', 2, 'Gust of Wind (2 charges)'),
            ('spells[3].spell', 'wind-wall', 'Wind Wall (1 charge)'),
            ('spells[3].charge_cost', 1, 'Wind Wall (1 charge)'),
            ('max_charges', 5, 'The ring has 5 charges'),
            ('recharge', {'dice': 1, 'sides': 4, 'type': 'charges'}, 'regains 1d4 + 1 expended charges daily at dawn'),
            ('recharge_bonus', 1, '1d4 + 1 expended charges'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
        ],
        'earth': [
            ('linked_plane', 'earth', 'Earth. You know'),
            ('plane_chosen_by', 'gm-chooses-or-rolls-randomly', 'The GM chooses or randomly determines the linked plane'),
            ('bane.target_creature_type', 'elemental', 'Advantage on attack rolls against Elementals'),
            ('bane.your_attacks', 'advantage', 'you have Advantage on attack rolls against Elementals'),
            ('bane.attacks_against_you', 'disadvantage', 'they have Disadvantage on attack rolls against you'),
            ('compulsion.action', 'magic', 'you can take a Magic action to try to compel an Elemental'),
            ('compulsion.range_feet', 60, 'an Elemental you see within 60 feet of yourself'),
            ('compulsion.save_ability', 'wisdom', 'DC 18 Wisdom saving throw'),
            ('compulsion.save_dc', 18, 'DC 18 Wisdom saving throw'),
            ('compulsion.condition', 'charmed', 'the Elemental has the Charmed condition'),
            ('compulsion.condition_ends', 'start-of-your-next-turn', 'until the start your next turn'),
            ('compulsion.control', 'you-determine-its-move-and-action-on-its-next-turn', 'you determine what it does with its move and action on its next turn'),
            ('focus.language', 'terran', 'You know Terran'),
            ('focus.defense', 'resistance', 'you have Resistance to Acid damage'),
            ('focus.damage_type', 'acid', 'you have Resistance to Acid damage'),
            ('focus.speed_kind', 'none', 'Earth. You know Terran'),
            ('focus.speed', 'none', 'Earth. You know Terran'),
            ('focus.other_benefits[0]', 'rubble-rocks-or-dirt-is-not-difficult-terrain', 'Terrain composed of rubble, rocks, or dirt isn’t Difficult Terrain for you'),
            ('focus.other_benefits[1]', 'move-through-solid-earth-or-rock-as-difficult-terrain', 'you can move through solid earth or rock as if those areas were Difficult Terrain'),
            ('focus.other_benefits[2]', 'passage-does-not-disturb-matter', 'without disturbing the matter through which you pass'),
            ('focus.other_benefits[3]', 'shunted-out-if-turn-ends-in-solid-earth-or-rock', 'If you end your turn in solid earth or rock, you are shunted out to the nearest unoccupied space you last occupied'),
            ('spell_save_dc', 18, 'which has a save DC of 18'),
            ('spells[0].spell', 'earthquake', 'Earthquake (5 charges)'),
            ('spells[0].charge_cost', 5, 'Earthquake (5 charges)'),
            ('spells[1].spell', 'stone-shape', 'Stone Shape (2 charges)'),
            ('spells[1].charge_cost', 2, 'Stone Shape (2 charges)'),
            ('spells[2].spell', 'stoneskin', 'Stoneskin (3 charges)'),
            ('spells[2].charge_cost', 3, 'Stoneskin (3 charges)'),
            ('spells[3].spell', 'wall-of-stone', 'Wall of Stone (3 charges)'),
            ('spells[3].charge_cost', 3, 'Wall of Stone (3 charges)'),
            ('max_charges', 5, 'The ring has 5 charges'),
            ('recharge', {'dice': 1, 'sides': 4, 'type': 'charges'}, 'regains 1d4 + 1 expended charges daily at dawn'),
            ('recharge_bonus', 1, '1d4 + 1 expended charges'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
        ],
        'fire': [
            ('linked_plane', 'fire', 'Fire. You know'),
            ('plane_chosen_by', 'gm-chooses-or-rolls-randomly', 'The GM chooses or randomly determines the linked plane'),
            ('bane.target_creature_type', 'elemental', 'Advantage on attack rolls against Elementals'),
            ('bane.your_attacks', 'advantage', 'you have Advantage on attack rolls against Elementals'),
            ('bane.attacks_against_you', 'disadvantage', 'they have Disadvantage on attack rolls against you'),
            ('compulsion.action', 'magic', 'you can take a Magic action to try to compel an Elemental'),
            ('compulsion.range_feet', 60, 'an Elemental you see within 60 feet of yourself'),
            ('compulsion.save_ability', 'wisdom', 'DC 18 Wisdom saving throw'),
            ('compulsion.save_dc', 18, 'DC 18 Wisdom saving throw'),
            ('compulsion.condition', 'charmed', 'the Elemental has the Charmed condition'),
            ('compulsion.condition_ends', 'start-of-your-next-turn', 'until the start your next turn'),
            ('compulsion.control', 'you-determine-its-move-and-action-on-its-next-turn', 'you determine what it does with its move and action on its next turn'),
            ('focus.language', 'ignan', 'You know Ignan'),
            ('focus.defense', 'immunity', 'you have Immunity to Fire damage'),
            ('focus.damage_type', 'fire', 'you have Immunity to Fire damage'),
            ('focus.speed_kind', 'none', 'Fire. You know Ignan'),
            ('focus.speed', 'none', 'Fire. You know Ignan'),
            ('focus.other_benefits', [], None),
            ('spell_save_dc', 18, 'which has a save DC of 18'),
            ('spells[0].spell', 'burning-hands', 'Burning Hands (1 charge)'),
            ('spells[0].charge_cost', 1, 'Burning Hands (1 charge)'),
            ('spells[1].spell', 'fireball', 'Fireball (2 charges)'),
            ('spells[1].charge_cost', 2, 'Fireball (2 charges)'),
            ('spells[2].spell', 'fire-storm', 'Fire Storm (4 charges)'),
            ('spells[2].charge_cost', 4, 'Fire Storm (4 charges)'),
            ('spells[3].spell', 'wall-of-fire', 'Wall of Fire (3 charges)'),
            ('spells[3].charge_cost', 3, 'Wall of Fire (3 charges)'),
            ('max_charges', 5, 'The ring has 5 charges'),
            ('recharge', {'dice': 1, 'sides': 4, 'type': 'charges'}, 'regains 1d4 + 1 expended charges daily at dawn'),
            ('recharge_bonus', 1, '1d4 + 1 expended charges'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
        ],
        'water': [
            ('linked_plane', 'water', 'Water. You know'),
            ('plane_chosen_by', 'gm-chooses-or-rolls-randomly', 'The GM chooses or randomly determines the linked plane'),
            ('bane.target_creature_type', 'elemental', 'Advantage on attack rolls against Elementals'),
            ('bane.your_attacks', 'advantage', 'you have Advantage on attack rolls against Elementals'),
            ('bane.attacks_against_you', 'disadvantage', 'they have Disadvantage on attack rolls against you'),
            ('compulsion.action', 'magic', 'you can take a Magic action to try to compel an Elemental'),
            ('compulsion.range_feet', 60, 'an Elemental you see within 60 feet of yourself'),
            ('compulsion.save_ability', 'wisdom', 'DC 18 Wisdom saving throw'),
            ('compulsion.save_dc', 18, 'DC 18 Wisdom saving throw'),
            ('compulsion.condition', 'charmed', 'the Elemental has the Charmed condition'),
            ('compulsion.condition_ends', 'start-of-your-next-turn', 'until the start your next turn'),
            ('compulsion.control', 'you-determine-its-move-and-action-on-its-next-turn', 'you determine what it does with its move and action on its next turn'),
            ('focus.language', 'aquan', 'You know Aquan'),
            ('focus.defense', 'none', 'Water. You know Aquan'),
            ('focus.damage_type', 'none', 'Water. You know Aquan'),
            ('focus.speed_kind', 'swim', 'you gain a Swim Speed of 60 feet'),
            ('focus.speed', '60-feet', 'you gain a Swim Speed of 60 feet'),
            ('focus.other_benefits[0]', 'can-breathe-underwater', 'and you can breathe underwater'),
            ('spell_save_dc', 18, 'which has a save DC of 18'),
            ('spells[0].spell', 'create-or-destroy-water', 'Create or Destroy Water (1 charge)'),
            ('spells[0].charge_cost', 1, 'Create or Destroy Water (1 charge)'),
            ('spells[1].spell', 'ice-storm', 'Ice Storm (2 charges)'),
            ('spells[1].charge_cost', 2, 'Ice Storm (2 charges)'),
            ('spells[2].spell', 'tsunami', 'Tsunami (5 charges)'),
            ('spells[2].charge_cost', 5, 'Tsunami (5 charges)'),
            ('spells[3].spell', 'wall-of-ice', 'Wall of Ice (3 charges)'),
            ('spells[3].charge_cost', 3, 'Wall of Ice (3 charges)'),
            ('spells[4].spell', 'water-walk', 'Water Walk (2 charges)'),
            ('spells[4].charge_cost', 2, 'Water Walk (2 charges)'),
            ('max_charges', 5, 'The ring has 5 charges'),
            ('recharge', {'dice': 1, 'sides': 4, 'type': 'charges'}, 'regains 1d4 + 1 expended charges daily at dawn'),
            ('recharge_bonus', 1, '1d4 + 1 expended charges'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
        ],
    },
    'ring-of-evasion': {
        'standard': [
            ('max_charges', 3, 'This ring has 3 charges'),
            ('recharge', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'it regains 1d3 expended charges daily at dawn'),
            ('recharge_when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('trigger', 'fail-dexterity-saving-throw-while-wearing', 'When you fail a Dexterity saving throw while wearing the ring'),
            ('action', 'reaction', 'you can take a Reaction to expend 1 charge'),
            ('charge_cost', 1, 'to expend 1 charge to succeed'),
            ('result', 'succeed-on-that-save', 'to succeed on that save instead'),
        ],
    },
    'ring-of-feather-falling': {
        'standard': [
            ('condition', 'while-falling-and-wearing', 'When you fall while wearing this ring'),
            ('descent_feet_per_round', 60, 'you descend 60 feet per round'),
            ('takes_falling_damage', False, 'take no damage from falling'),
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


def healing(stream):
    """Potions of Healing is held by the Healing contract fixture; compare its rows."""
    block = stream[stream.index("Potions of Healing Potion, Rarity Varies"):stream.index("Potion of Heroism Potion, Rare")]
    contract = json.loads((ROOT / "internal/rules/srd521/testdata/item-families.json").read_text())
    family = contract["families"][0]
    assert contract["source"]["sha256"] == PIN and family["id"] == "srd-5.2.1:content/items/potion-of-healing"
    assert {"id": "srd-5.2.1:content/magic-items/potions-of-healing", "pages": [236, 236]} in family["presentations"]
    rows = [("basic", "Potion of Healing 2d4 + 2 Common"), ("greater", "4d4 + 4 Uncommon (greater)"),
            ("superior", "8d4 + 8 Rare (superior)"), ("supreme", "10d4 + 20 Very Rare (supreme)")]
    assert [v["id"].rsplit("/", 1)[1] for v in family["variants"]] == [slug for slug, _ in rows]
    for v, (slug, printed) in zip(family["variants"], rows):
        h = v["healing"]
        dice = f"{h['dice']}d{h['sides']} + {h['bonus']} {RARITY[v['rarity']]}"
        assert dice in printed and has(block, printed), (slug, dice, printed)


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_srd_items_238.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_238.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 238]
    assert len(expected) == 12
    # The Healing contract fixture holds the twelfth family; see healing().
    expected = [e for e in expected if e[0] != "srd-5.2.1:content/magic-items/potions-of-healing"]
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 238 and catalogue["batch"] == "data-magic-items-14"
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
    healing(stream)
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 236-238")


main()
