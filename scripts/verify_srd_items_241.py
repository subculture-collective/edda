#!/usr/bin/env python3
"""Compare issue #241 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (241, 242, 243, 244)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = ['Slippers of Spider Climbing Wondrous', 'Shield, +1, +2, or +3 Armor']

# Each family's printed category line.
HEADERS = {
    'rod-of-lordly-might': 'Rod, Legendary (Requires Attunement)',
    'rod-of-resurrection': 'Rod, Legendary (Requires Attunement)',
    'rod-of-rulership': 'Rod, Rare (Requires Attunement)',
    'rod-of-security': 'Rod, Very Rare',
    'rope-of-climbing': 'Wondrous Item, Uncommon',
    'rope-of-entanglement': 'Wondrous Item, Rare',
    'scarab-of-protection': 'Wondrous Item, Legendary (Requires Attunement)',
    'scimitar-of-speed': 'Weapon (Scimitar), Very Rare (Requires Attunement)',
    'sending-stones': 'Wondrous Item, Uncommon',
    'sentinel-shield': 'Armor (Shield), Uncommon',
    'shield-of-missile-attraction': 'Armor (Shield), Rare (Requires Attunement)',
    'shield-of-the-cavalier': 'Armor (Shield), Very Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'rod-of-resurrection': ({'mode': 'limited', 'maximum': 5, 'per_use': 1, 'recharge': 'regain-1-at-dawn'}, 'The rod has 5 charges'),
    'scarab-of-protection': ({'mode': 'limited', 'maximum': 12, 'per_use': 1, 'recharge': 'none'}, 'The scarab has 12 charges'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'rod-of-lordly-might': {
        'standard': [
            ('functions_as', 'mace', 'functions as a magic Mace'),
            ('attack_bonus', 3, 'a magic Mace that grants a +3 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 3, 'a +3 bonus to attack rolls and damage rolls made with it'),
            ('button_count', 6, 'six different buttons'),
            ('buttons.press_action', 'bonus', 'You can press one of the following buttons as a Bonus Action'),
            ('buttons.effect_lasts_until[0]', 'different-button-pushed', 'a button’s effect lasts until you push a different button or until you push the same button again'),
            ('buttons.effect_lasts_until[1]', 'same-button-pushed-again', 'a button’s effect lasts until you push a different button or until you push the same button again'),
            ('buttons.fiery_blade.button', 1, 'Button 1. A fiery blade'),
            ('buttons.fiery_blade.bright_light_radius_feet', 40, 'Bright Light in a 40-foot radius'),
            ('buttons.fiery_blade.dim_light_additional_feet', 40, 'Dim Light for an additional 40 feet'),
            ('buttons.fiery_blade.functions_as[0]', 'longsword', 'a magic Longsword or Shortsword (your choice)'),
            ('buttons.fiery_blade.functions_as[1]', 'shortsword', 'a magic Longsword or Shortsword (your choice)'),
            ('buttons.fiery_blade.extra_damage', {'dice': 2, 'sides': 6, 'type': 'fire'}, 'an extra 2d6 Fire damage on a hit'),
            ('buttons.battleaxe.button', 2, 'Button 2. The rod’s flanged head folds down'),
            ('buttons.battleaxe.functions_as', 'battleaxe', 'into a magic Battleaxe'),
            ('buttons.battleaxe.attack_and_damage_bonus', 3, 'magic Battleaxe that grants a +3 bonus'),
            ('buttons.spear.button', 3, 'Button 3.'),
            ('buttons.spear.functions_as', 'spear', 'into a magic Spear'),
            ('buttons.spear.haft_feet', 6, 'lengthens into a 6-foot haft'),
            ('buttons.spear.attack_and_damage_bonus', 3, 'magic Spear that grants a +3 bonus'),
            ('buttons.climbing_pole.button', 4, 'Button 4. The rod transforms into a climbing pole'),
            ('buttons.climbing_pole.maximum_length_feet', 50, 'up to 50 feet long (you specify the length)'),
            ('buttons.climbing_pole.buttons_stay_in_reach', True, 'the rod’s buttons remain within your reach'),
            ('buttons.climbing_pole.anchors[0]', 'bottom-spike', 'a spike at the bottom and three hooks at the top anchor the pole'),
            ('buttons.climbing_pole.anchors[1]', 'three-top-hooks', 'a spike at the bottom and three hooks at the top anchor the pole'),
            ('buttons.climbing_pole.anchor_surface', 'as-hard-as-granite', 'In surfaces as hard as granite'),
            ('buttons.climbing_pole.bar_length_inches', 3, 'Horizontal bars 3 inches long'),
            ('buttons.climbing_pole.bar_spacing_feet', 1, '1 foot apart, forming a ladder'),
            ('buttons.climbing_pole.capacity_pounds', 4000, 'bear up to 4,000 pounds'),
            ('buttons.climbing_pole.reverts_on[0]', 'more-weight', 'More weight or lack of solid anchoring causes the rod to revert'),
            ('buttons.climbing_pole.reverts_on[1]', 'lack-of-solid-anchoring', 'More weight or lack of solid anchoring causes the rod to revert'),
            ('buttons.battering_ram.button', 5, 'Button 5. The rod transforms into a handheld battering ram'),
            ('buttons.battering_ram.check_bonus', 10, 'grants its user a +10 bonus'),
            ('buttons.battering_ram.check_ability', 'strength', 'to Strength (Athletics) checks'),
            ('buttons.battering_ram.check_skill', 'athletics', 'Strength (Athletics) checks'),
            ('buttons.battering_ram.check_subjects[0]', 'doors', 'to break through doors, barricades, and other barriers'),
            ('buttons.battering_ram.check_subjects[1]', 'barricades', 'to break through doors, barricades, and other barriers'),
            ('buttons.battering_ram.check_subjects[2]', 'other-barriers', 'to break through doors, barricades, and other barriers'),
            ('buttons.normal_form.button', 6, 'Button 6. The rod assumes or remains in its normal form'),
            ('buttons.normal_form.indicates', 'magnetic-north', 'indicates magnetic north'),
            ('buttons.normal_form.without_magnetic_north', 'nothing-happens', 'Nothing happens if this function of the rod is used in a location that has no magnetic north'),
            ('buttons.normal_form.knowledge[0]', 'approximate-depth-below-ground', 'knowledge of your approximate depth beneath the ground or your height above it'),
            ('buttons.normal_form.knowledge[1]', 'height-above-ground', 'knowledge of your approximate depth beneath the ground or your height above it'),
            ('drain_life.trigger', 'hit-creature-with-melee-attack', 'When you hit a creature with a melee attack using the rod'),
            ('drain_life.save_ability', 'constitution', 'DC 17 Constitution saving throw. On a failed save, the target takes an extra 4d6'),
            ('drain_life.save_dc', 17, 'DC 17 Constitution saving throw. On a failed save, the target takes an extra 4d6'),
            ('drain_life.extra_damage', {'dice': 4, 'sides': 6, 'type': 'necrotic'}, 'an extra 4d6 Necrotic damage'),
            ('drain_life.you_regain', 'half-of-that-necrotic-damage', 'you regain a number of Hit Points equal to half that Necrotic damage'),
            ('drain_life.reuse_after', 'next-dawn', 'can’t be used again until the next dawn'),
            ('paralyze.trigger', 'hit-creature-with-melee-attack', 'When you hit a creature with a melee attack using the rod'),
            ('paralyze.save_ability', 'constitution', 'DC 17 Constitution saving throw. On a failed save, the target has the Paralyzed'),
            ('paralyze.save_dc', 17, 'DC 17 Constitution saving throw. On a failed save, the target has the Paralyzed'),
            ('paralyze.condition', 'paralyzed', 'the target has the Paralyzed condition'),
            ('paralyze.duration_minutes', 1, 'the Paralyzed condition for 1 minute'),
            ('paralyze.repeat_save', 'end-of-each-of-its-turns', 'The target repeats the save at the end of each of its turns, ending the effect on a success'),
            ('paralyze.reuse_after', 'next-dawn', 'can’t be used again until the next dawn'),
            ('terrify.condition', 'frightened', 'a target has the Frightened condition'),
            ('terrify.action', 'magic', 'you can take a Magic action to force each creature'),
            ('terrify.range_feet', 30, 'each creature you can see within 30 feet of yourself'),
            ('terrify.save_ability', 'wisdom', 'DC 17 Wisdom saving throw'),
            ('terrify.save_dc', 17, 'DC 17 Wisdom saving throw'),
            ('terrify.duration_minutes', 1, 'the Frightened condition for 1 minute'),
            ('terrify.repeat_save', 'end-of-each-of-its-turns', 'A Frightened target repeats the save at the end of each of its turns'),
            ('terrify.repeat_success', 'ends-effect-on-itself', 'ending the effect on itself on a success'),
            ('terrify.reuse_after', 'next-dawn', 'can’t be used again until the next dawn'),
        ],
    },
    'rod-of-resurrection': {
        'standard': [
            ('max_charges', 5, 'The rod has 5 charges'),
            ('condition', 'while-holding', 'While you hold it'),
            ('spells[0].spell', 'heal', 'Heal (expends 1 charge)'),
            ('spells[0].charge_cost', 1, 'Heal (expends 1 charge)'),
            ('spells[1].spell', 'resurrection', 'or Resurrection (expends 5 charges)'),
            ('spells[1].charge_cost', 5, 'Resurrection (expends 5 charges)'),
            ('recharge_charges', 1, 'The rod regains 1 expended charge daily at dawn'),
            ('recharge_when', 'daily-at-dawn', 'expended charge daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the rod disappears'),
            ('last_charge.result', 'disappears-in-harmless-burst-of-radiance', 'disappears in a harmless burst of radiance'),
        ],
    },
    'rod-of-rulership': {
        'standard': [
            ('action', 'magic', 'You can take a Magic action to present the rod'),
            ('targets', 'each-creature-of-your-choice-you-can-see', 'command obedience from each creature of your choice that you can see'),
            ('range_feet', 120, 'within 120 feet of yourself'),
            ('save_ability', 'wisdom', 'DC 15 Wisdom saving throw'),
            ('save_dc', 15, 'DC 15 Wisdom saving throw'),
            ('condition', 'charmed', 'or have the Charmed condition'),
            ('duration_hours', 8, 'the Charmed condition for 8 hours'),
            ('regards_you_as', 'trusted-leader', 'the creature regards you as its trusted leader'),
            ('ends_early_when[0]', 'harmed-by-you-or-allies', 'If harmed by you or your allies or commanded to do something contrary to its nature'),
            ('ends_early_when[1]', 'commanded-against-its-nature', 'If harmed by you or your allies or commanded to do something contrary to its nature'),
            ('reuse_after', 'next-dawn', 'can’t be used again until the next dawn'),
        ],
    },
    'rod-of-security': {
        'standard': [
            ('condition', 'while-holding', 'While holding this rod'),
            ('action', 'magic', 'you can take a Magic action to activate it'),
            ('travelers.includes_you', True, 'instantly transports you and up to 199 other willing creatures'),
            ('travelers.other_willing_creatures', 199, 'up to 199 other willing creatures you can see'),
            ('destination', 'demiplane', 'to a demiplane'),
            ('form_chosen_by', 'you', 'You choose the form the demiplane takes'),
            ('demiplane.provides', 'enough-water-and-food', 'enough water and food to sustain its visitors'),
            ('demiplane.environment_can_harm', False, 'environment can’t harm its occupants'),
            ('demiplane.objects_exist_only_there', True, 'Everything else that can be interacted with there can exist only there'),
            ('demiplane.aging', False, 'creatures don’t age while there'),
            ('demiplane.time_passes_normally', True, 'although time passes normally'),
            ('healing.per', 'each-hour-spent-in-demiplane', 'For each hour spent in the demiplane'),
            ('healing.as_if_hit_point_dice_spent', 1, 'as if it had spent 1 Hit Point Die'),
            ('stay.base_days', 200, 'up to 200 days divided by the number of creatures present'),
            ('stay.divided_by', 'creatures-present', 'divided by the number of creatures present'),
            ('stay.rounding', 'down', '(round down)'),
            ('ends_when[0]', 'time-runs-out', 'When the time runs out or you take a Magic action to end the effect'),
            ('ends_when[1]', 'magic-action-to-end', 'When the time runs out or you take a Magic action to end the effect'),
            ('return_to', 'original-or-nearest-unoccupied-space', 'the location they occupied when you activated the rod or an unoccupied space nearest that location'),
            ('reuse_after_days', 10, 'can’t be used again until 10 days have passed'),
        ],
    },
    'rope-of-climbing': {
        'standard': [
            ('length_feet', 60, 'This 60-foot length of rope'),
            ('capacity_pounds', 3000, 'can hold up to 3,000 pounds'),
            ('condition', 'while-holding-one-end', 'While holding one end of the rope'),
            ('action', 'magic', 'you can take a Magic action to command the other end of the rope to animate'),
            ('move_limit', 'rope-length-away-from-you', 'up to the rope’s length away from you'),
            ('move_feet_per_turn', 10, 'That end moves 10 feet on your turn when you first command it and 10 feet at the start of each of your subsequent turns'),
            ('moves_until[0]', 'reaching-destination', 'until reaching its destination or until you tell it to stop'),
            ('moves_until[1]', 'told-to-stop', 'until reaching its destination or until you tell it to stop'),
            ('other_commands[0]', 'fasten', 'You can also tell the rope to fasten itself securely to an object or to unfasten itself, to knot or unknot itself, or to coil itself for carrying'),
            ('other_commands[1]', 'unfasten', 'You can also tell the rope to fasten itself securely to an object or to unfasten itself, to knot or unknot itself, or to coil itself for carrying'),
            ('other_commands[2]', 'knot', 'You can also tell the rope to fasten itself securely to an object or to unfasten itself, to knot or unknot itself, or to coil itself for carrying'),
            ('other_commands[3]', 'unknot', 'You can also tell the rope to fasten itself securely to an object or to unfasten itself, to knot or unknot itself, or to coil itself for carrying'),
            ('other_commands[4]', 'coil', 'You can also tell the rope to fasten itself securely to an object or to unfasten itself, to knot or unknot itself, or to coil itself for carrying'),
            ('knotted.knot_spacing_feet', 1, 'large knots appear at 1-foot intervals'),
            ('knotted.length_feet', 50, 'the rope shortens to a 50-foot length'),
            ('knotted.check_mode', 'advantage', 'grants Advantage on ability checks made to climb using the rope'),
            ('knotted.check_subject', 'climb-using-the-rope', 'ability checks made to climb using the rope'),
            ('armor_class', 20, 'The rope has AC 20'),
            ('hit_points', 20, 'AC 20, HP 20'),
            ('immunities[0]', 'poison', 'Immunity to Poison and Psychic damage'),
            ('immunities[1]', 'psychic', 'Immunity to Poison and Psychic damage'),
            ('regains_hit_points', 1, 'It regains 1 Hit Point every 5 minutes'),
            ('regain_interval_minutes', 5, '1 Hit Point every 5 minutes'),
            ('regain_needs_hit_points', 1, 'as long as it has at least 1 Hit Point'),
            ('destroyed_at_hit_points', 0, 'If the rope drops to 0 Hit Points, it is destroyed'),
        ],
    },
    'rope-of-entanglement': {
        'standard': [
            ('length_feet', 30, 'This rope is 30 feet long'),
            ('condition', 'while-holding-one-end', 'While holding one end of the rope'),
            ('action', 'magic', 'you can take a Magic action to command the other end to dart forward'),
            ('targets', 1, 'entangle one creature you can see'),
            ('range_feet', 20, 'within 20 feet of yourself'),
            ('save_ability', 'dexterity', 'DC 15 Dexterity saving throw'),
            ('save_dc', 15, 'DC 15 Dexterity saving throw'),
            ('applies_condition', 'restrained', 'or have the Restrained condition'),
            ('release[0].method', 'let-go-of-your-end', 'by letting go of your end of the rope'),
            ('release[0].rope_coils_in', 'targets-space', 'causing the rope to coil up in the target’s space'),
            ('release[1].method', 'bonus-action-to-repeat-command', 'or by using a Bonus Action to repeat the command'),
            ('release[1].rope_coils_in', 'your-hand', '(causing the rope to coil up in your hand)'),
            ('escape.action', 'action', 'A target Restrained by the rope can take an action'),
            ('escape.choices[0].ability', 'strength', 'DC 15 Strength (Athletics)'),
            ('escape.choices[0].skill', 'athletics', 'Strength (Athletics)'),
            ('escape.choices[1].ability', 'dexterity', 'Dexterity (Acrobatics) check'),
            ('escape.choices[1].skill', 'acrobatics', 'Dexterity (Acrobatics)'),
            ('escape.dc', 15, 'DC 15 Strength (Athletics) or Dexterity (Acrobatics) check'),
            ('escape.success', 'no-longer-restrained', 'On a successful check, the target is no longer Restrained by the rope'),
            ('recoil.trigger', 'target-escapes-while-you-hold-rope', 'If you’re still holding onto the rope when a target escapes from it'),
            ('recoil.action', 'reaction', 'you can take a Reaction to command the rope to coil up in your hand'),
            ('recoil.otherwise', 'rope-coils-in-targets-space', 'otherwise, the rope coils up in the target’s space'),
            ('armor_class', 20, 'The rope has AC 20'),
            ('hit_points', 20, 'AC 20, HP 20'),
            ('immunities[0]', 'poison', 'Immunity to Poison and Psychic damage'),
            ('immunities[1]', 'psychic', 'Immunity to Poison and Psychic damage'),
            ('regains_hit_points', 1, 'It regains 1 Hit Point every 5 minutes'),
            ('regain_interval_minutes', 5, '1 Hit Point every 5 minutes'),
            ('regain_needs_hit_points', 1, 'as long as it has at least 1 Hit Point'),
            ('destroyed_at_hit_points', 0, 'If the rope drops to 0 Hit Points, it is destroyed'),
        ],
    },
    'scarab-of-protection': {
        'standard': [
            ('appearance', 'beetle-shaped-medallion', 'This beetle-shaped medallion'),
            ('condition', 'while-on-your-person', 'while it is on your person'),
            ('benefit_count', 3, 'provides three benefits'),
            ('armor_class_bonus', 1, '+1 bonus to Armor Class'),
            ('preservation.max_charges', 12, 'The scarab has 12 charges'),
            ('preservation.triggers[0]', 'fail-save-against-necromancy-spell', 'If you fail a saving throw against a Necromancy spell'),
            ('preservation.triggers[1]', 'fail-save-against-harmful-effect-from-undead', 'or a harmful effect originating from an Undead'),
            ('preservation.action', 'reaction', 'you can take a Reaction to expend 1 charge'),
            ('preservation.charge_cost', 1, 'to expend 1 charge and turn the failed save into a successful one'),
            ('preservation.result', 'failed-save-becomes-success', 'turn the failed save into a successful one'),
            ('preservation.after_last_charge', 'crumbles-to-powder-and-is-destroyed', 'The scarab crumbles into powder and is destroyed when its last charge is expended'),
            ('spell_resistance.save_mode', 'advantage', 'You have Advantage on saving throws against spells'),
            ('spell_resistance.against', 'spells', 'saving throws against spells'),
        ],
    },
    'scimitar-of-speed': {
        'standard': [
            ('attack_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('extra_attack_action', 'bonus', 'you can make one attack with it as a Bonus Action'),
            ('extra_attacks', 1, 'make one attack with it'),
            ('frequency', 'each-of-your-turns', 'on each of your turns'),
        ],
    },
    'sending-stones': {
        'standard': [
            ('stones_per_pair', 'pair', 'Sending Stones come in pairs'),
            ('matching', 'carved-to-match', 'each stone carved to match the other'),
            ('condition', 'while-touching-one-stone', 'While you touch one stone'),
            ('spell', 'sending', 'you can cast Sending from it'),
            ('target', 'bearer-of-other-stone', 'The target is the bearer of the other stone'),
            ('if_no_bearer', 'you-learn-this-and-do-not-cast', 'you know that fact as soon as you use the stone, and you don’t cast the spell'),
            ('shared_reuse_after', 'next-dawn', 'the stones can’t be used again until the next dawn'),
            ('other_stone_when_one_destroyed', 'nonmagical', 'If one of the stones in a pair is destroyed, the other one becomes nonmagical'),
        ],
    },
    'sentinel-shield': {
        'standard': [
            ('condition', 'while-holding', 'While holding this Shield'),
            ('advantage_on[0]', 'initiative-rolls', 'Advantage on Initiative rolls'),
            ('advantage_on[1]', 'wisdom-perception-checks', 'and Wisdom (Perception) checks'),
            ('emblem', 'eye', 'emblazoned with a symbol of an eye'),
        ],
    },
    'shield-of-missile-attraction': {
        'standard': [
            ('condition', 'while-holding', 'While holding this Shield'),
            ('defense', 'resistance', 'you have Resistance to damage'),
            ('defense_against', 'damage-from-ranged-weapon-attacks', 'damage from attacks made with Ranged weapons'),
            ('curse.begins_on', 'attunement', 'Attuning to it curses you'),
            ('curse.ends_on', 'remove-curse-or-similar-magic', 'until you are targeted by a Remove Curse spell or similar magic'),
            ('curse.removing_shield_ends_curse', False, 'Removing the Shield fails to end the curse on you'),
            ('curse.trigger', 'ranged-weapon-attack-targets-nearby-creature', 'Whenever an attack with a Ranged weapon targets a creature'),
            ('curse.radius_feet', 10, 'a creature within 10 feet of you'),
            ('curse.result', 'you-become-the-target', 'the curse causes you to become the target instead'),
        ],
    },
    'shield-of-the-cavalier': {
        'standard': [
            ('condition', 'while-holding', 'While holding this Shield'),
            ('armor_class_bonus', 2, '+2 bonus to Armor Class'),
            ('stacks_with_shield_bonus', True, 'This bonus is in addition to the Shield’s normal bonus to AC'),
            ('forceful_bash.trigger', 'attack-action', 'When you take the Attack action'),
            ('forceful_bash.attacks_replaced', 1, 'you can make one of the attack rolls using the Shield'),
            ('forceful_bash.reach_feet', 5, 'against a target within 5 feet of yourself'),
            ('forceful_bash.attack_roll_adds[0]', 'proficiency-bonus', 'Apply your Proficiency Bonus and Strength modifier to the attack roll'),
            ('forceful_bash.attack_roll_adds[1]', 'strength-modifier', 'Apply your Proficiency Bonus and Strength modifier to the attack roll'),
            ('forceful_bash.damage', {'dice': 2, 'sides': 6, 'type': 'force'}, 'the Shield deals Force damage to the target equal to 2d6 + 2 plus your Strength modifier'),
            ('forceful_bash.damage_bonus', 2, '2d6 + 2'),
            ('forceful_bash.damage_adds_ability_modifier', 'strength', '2d6 + 2 plus your Strength modifier'),
            ('forceful_bash.push_feet', 10, 'push it up to 10 feet directly away from yourself'),
            ('forceful_bash.prone_if', 'creature-your-size-or-smaller', 'If the creature is your size or smaller'),
            ('forceful_bash.applies_condition', 'prone', 'giving it the Prone condition'),
            ('protective_field.action', 'reaction', 'As a Reaction, when you or an ally'),
            ('protective_field.triggers[0]', 'you-or-nearby-ally-targeted-by-attack', 'is targeted by an attack'),
            ('protective_field.triggers[1]', 'you-or-nearby-ally-saves-against-area-effect', 'or makes a saving throw against an area of effect'),
            ('protective_field.ally_range_feet', 5, 'an ally you can see within 5 feet of you'),
            ('protective_field.emanation_feet', 5, 'immobile 5-foot Emanation originating from you'),
            ('protective_field.immobile', True, 'create an immobile'),
            ('protective_field.not_fully_contained', 'pushed-to-nearest-unoccupied-space-outside', 'not fully contained within it are pushed into the nearest unoccupied spaces outside it'),
            ('protective_field.triggering_effect', 'no-effect-inside', 'has no effect on creatures and objects inside the Emanation'),
            ('protective_field.concentration_required', True, 'lasts as long as you maintain Concentration'),
            ('protective_field.max_duration_minutes', 1, 'up to 1 minute'),
            ('protective_field.passage', 'nothing-in-or-out', 'Nothing can pass into or out of the Emanation'),
            ('protective_field.isolation[0]', 'inside-cannot-be-damaged-from-outside', 'can’t be damaged by attacks or effects originating from outside, nor can a creature inside the Emanation damage anything outside it'),
            ('protective_field.isolation[1]', 'inside-cannot-damage-outside', 'can’t be damaged by attacks or effects originating from outside, nor can a creature inside the Emanation damage anything outside it'),
            ('protective_field.reuse_after', 'next-dawn', 'can’t be used again until the next dawn'),
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
        raise SystemExit("usage: verify_srd_items_241.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_241.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 241]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 241 and catalogue["batch"] == "data-magic-items-17"
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
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 241-243")


main()
