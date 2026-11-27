#!/usr/bin/env python3
"""Compare issue #242 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (243, 244, 245, 246)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = ['Staff of Swarming Insects Staff,', 'Shield of Missile Attraction Armor']

# Each family's printed category line.
HEADERS = {
    'shield-1-2-or-3': 'Armor (Shield), Uncommon (+1), Rare (+2), or Very Rare (+3)',
    'slippers-of-spider-climbing': 'Wondrous Item, Uncommon (Requires Attunement)',
    'sovereign-glue': 'Wondrous Item, Legendary',
    'spell-scroll': 'Scroll, Rarity Varies',
    'spellguard-shield': 'Armor (Shield), Very Rare (Requires Attunement)',
    'sphere-of-annihilation': 'Wondrous Item, Legendary',
    'staff-of-charming': 'Staff, Rare (Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard)',
    'staff-of-fire': 'Staff, Very Rare (Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard)',
    'staff-of-frost': 'Staff, Very Rare (Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard)',
    'staff-of-healing': 'Staff, Rare (Requires Attunement by a Bard, Cleric, or Druid)',
    'staff-of-power': 'Staff, Very Rare (Requires Attunement by a Sorcerer, Warlock, or Wizard)',
    'staff-of-striking': 'Staff, Very Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'staff-of-charming': ({'mode': 'limited', 'maximum': 10, 'per_use': 1, 'recharge': 'regain-1d8-plus-2-at-dawn'}, 'This staff has 10 charges'),
    'staff-of-fire': ({'mode': 'limited', 'maximum': 10, 'per_use': 1, 'recharge': 'regain-1d6-plus-4-at-dawn'}, 'The staff has 10 charges'),
    'staff-of-frost': ({'mode': 'limited', 'maximum': 10, 'per_use': 1, 'recharge': 'regain-1d6-plus-4-at-dawn'}, 'The staff has 10 charges'),
    'staff-of-healing': ({'mode': 'limited', 'maximum': 10, 'per_use': 1, 'recharge': 'regain-1d6-plus-4-at-dawn'}, 'This staff has 10 charges'),
    'staff-of-power': ({'mode': 'limited', 'maximum': 20, 'per_use': 1, 'recharge': 'regain-2d8-plus-4-at-dawn'}, 'This staff has 20 charges'),
    'staff-of-striking': ({'mode': 'limited', 'maximum': 10, 'per_use': 1, 'recharge': 'regain-1d6-plus-4-at-dawn'}, 'The staff has 10 charges'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'shield-1-2-or-3': {
        'plus-1': [
            ('condition', 'while-holding', 'While holding this Shield'),
            ('armor_class_bonus', 1, 'Uncommon (+1)'),
            ('stacks_with_shield_bonus', True, 'in addition to the Shield’s normal bonus to AC'),
        ],
        'plus-2': [
            ('condition', 'while-holding', 'While holding this Shield'),
            ('armor_class_bonus', 2, 'Rare (+2)'),
            ('stacks_with_shield_bonus', True, 'in addition to the Shield’s normal bonus to AC'),
        ],
        'plus-3': [
            ('condition', 'while-holding', 'While holding this Shield'),
            ('armor_class_bonus', 3, 'Very Rare (+3)'),
            ('stacks_with_shield_bonus', True, 'in addition to the Shield’s normal bonus to AC'),
        ],
    },
    'slippers-of-spider-climbing': {
        'standard': [
            ('condition', 'while-wearing', 'While you wear these light shoes'),
            ('movement[0]', 'up-down-and-across-vertical-surfaces', 'you can move up, down, and across vertical surfaces and along ceilings'),
            ('movement[1]', 'along-ceilings', 'you can move up, down, and across vertical surfaces and along ceilings'),
            ('hands_free', True, 'while leaving your hands free'),
            ('climb_speed', 'equal-to-speed', 'You have a Climb Speed equal to your Speed'),
            ('excluded_surfaces', 'slippery-such-as-ice-or-oil', 'on a slippery surface, such as one covered by ice or oil'),
        ],
    },
    'sovereign-glue': {
        'standard': [
            ('appearance', 'viscous-milky-white-substance', 'This viscous, milky-white substance'),
            ('bond', 'permanent-adhesive-between-two-objects', 'form a permanent adhesive bond between any two objects'),
            ('storage', 'jar-or-flask-coated-with-oil-of-slipperiness', 'stored in a jar or flask that has been coated inside with Oil of Slipperiness'),
            ('found_ounces', {'dice': 1, 'sides': 6, 'type': 'ounces'}, 'a container contains 1d6 + 1 ounces'),
            ('found_ounces_bonus', 1, '1d6 + 1 ounces'),
            ('coverage_ounces', 1, 'One ounce of the glue'),
            ('coverage_square_feet', 1, 'cover a 1-foot square surface'),
            ('apply_action', 'utilize', 'takes a Utilize action'),
            ('set_minutes', 1, 'the applied glue takes 1 minute to set'),
            ('broken_only_by[0]', 'universal-solvent', 'broken only by the application of Universal Solvent or Oil of Etherealness, or with a Wish spell'),
            ('broken_only_by[1]', 'oil-of-etherealness', 'broken only by the application of Universal Solvent or Oil of Etherealness, or with a Wish spell'),
            ('broken_only_by[2]', 'wish', 'broken only by the application of Universal Solvent or Oil of Etherealness, or with a Wish spell'),
        ],
    },
    'spell-scroll': {
        'cantrip': [
            ('spell_level', 0, 'Cantrip Common 13 +5'),
            ('rarity', 'common', 'Cantrip Common 13 +5'),
            ('save_dc', 13, 'Cantrip Common 13 +5'),
            ('attack_bonus', 5, 'Cantrip Common 13 +5'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-1': [
            ('spell_level', 1, '1 Common 13 +5'),
            ('rarity', 'common', '1 Common 13 +5'),
            ('save_dc', 13, '1 Common 13 +5'),
            ('attack_bonus', 5, '1 Common 13 +5'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-2': [
            ('spell_level', 2, '2 Uncommon 13 +5'),
            ('rarity', 'uncommon', '2 Uncommon 13 +5'),
            ('save_dc', 13, '2 Uncommon 13 +5'),
            ('attack_bonus', 5, '2 Uncommon 13 +5'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-3': [
            ('spell_level', 3, '3 Uncommon 15 +7'),
            ('rarity', 'uncommon', '3 Uncommon 15 +7'),
            ('save_dc', 15, '3 Uncommon 15 +7'),
            ('attack_bonus', 7, '3 Uncommon 15 +7'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-4': [
            ('spell_level', 4, '4 Rare 15 +7'),
            ('rarity', 'rare', '4 Rare 15 +7'),
            ('save_dc', 15, '4 Rare 15 +7'),
            ('attack_bonus', 7, '4 Rare 15 +7'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-5': [
            ('spell_level', 5, '5 Rare 17 +9'),
            ('rarity', 'rare', '5 Rare 17 +9'),
            ('save_dc', 17, '5 Rare 17 +9'),
            ('attack_bonus', 9, '5 Rare 17 +9'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-6': [
            ('spell_level', 6, '6 Very Rare 17 +9'),
            ('rarity', 'very-rare', '6 Very Rare 17 +9'),
            ('save_dc', 17, '6 Very Rare 17 +9'),
            ('attack_bonus', 9, '6 Very Rare 17 +9'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-7': [
            ('spell_level', 7, '7 Very Rare 18 +10'),
            ('rarity', 'very-rare', '7 Very Rare 18 +10'),
            ('save_dc', 18, '7 Very Rare 18 +10'),
            ('attack_bonus', 10, '7 Very Rare 18 +10'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-8': [
            ('spell_level', 8, '8 Very Rare 18 +10'),
            ('rarity', 'very-rare', '8 Very Rare 18 +10'),
            ('save_dc', 18, '8 Very Rare 18 +10'),
            ('attack_bonus', 10, '8 Very Rare 18 +10'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
        'level-9': [
            ('spell_level', 9, '9 Legendary 19 +11'),
            ('rarity', 'legendary', '9 Legendary 19 +11'),
            ('save_dc', 19, '9 Legendary 19 +11'),
            ('attack_bonus', 11, '9 Legendary 19 +11'),
            ('bears', 'one-spell', 'bears the words of a single spell'),
            ('readable_if', 'spell-on-your-spell-list', 'If the spell is on your spell list, you can read the scroll and cast its spell'),
            ('material_components_required', False, 'without Material components'),
            ('otherwise', 'unintelligible', 'Otherwise, the scroll is unintelligible'),
            ('casting_time', 'spells-normal-casting-time', 'requires the spell’s normal casting time'),
            ('after_casting', 'crumbles-to-dust', 'Once the spell is cast, the scroll crumbles to dust'),
            ('lost_if_casting_interrupted', False, 'If the casting is interrupted, the scroll isn’t lost'),
            ('higher_level_check.trigger', 'higher-level-than-you-can-normally-cast', 'of a higher level than you can normally cast'),
            ('higher_level_check.ability', 'your-spellcasting-ability', 'you make an ability check using your spellcasting ability'),
            ('higher_level_check.dc_base', 10, 'The DC equals 10 plus the spell’s level. On a failed check, the spell disappears'),
            ('higher_level_check.dc_adds', 'spell-level', '10 plus the spell’s level'),
            ('higher_level_check.failure', 'spell-disappears-with-no-other-effect', 'the spell disappears from the scroll with no other effect'),
            ('copying.eligible', 'wizard-spell', 'A Wizard spell on a Spell Scroll can be copied into a spellbook'),
            ('copying.check_ability', 'intelligence', 'Intelligence (Arcana) check'),
            ('copying.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('copying.dc_base', 10, 'with a DC equal to 10 plus the spell’s level'),
            ('copying.dc_adds', 'spell-level', 'DC equal to 10 plus the spell’s level'),
            ('copying.success', 'spell-copied', 'On a successful check, the spell is copied'),
            ('copying.scroll_destroyed_either_way', True, 'Whether the check succeeds or fails, the Spell Scroll is destroyed'),
        ],
    },
    'spellguard-shield': {
        'standard': [
            ('condition', 'while-holding', 'While holding this Shield'),
            ('save_mode', 'advantage', 'you have Advantage on saving throws'),
            ('against[0]', 'spells', 'against spells and other magical effects'),
            ('against[1]', 'other-magical-effects', 'against spells and other magical effects'),
            ('spell_attacks_against_you', 'disadvantage', 'spell attack rolls have Disadvantage against you'),
        ],
    },
    'sphere-of-annihilation': {
        'standard': [
            ('diameter_feet', 2, 'This 2-foot-diameter black sphere'),
            ('nature', 'hole-in-the-multiverse', 'a hole in the multiverse'),
            ('obliterates', 'all-matter-passing-through-or-passed-through', 'obliterates all matter it passes through and all matter that passes through it'),
            ('artifact_exception', 'artifacts-unscathed-unless-susceptible', 'Unless an Artifact is susceptible to damage from a Sphere of Annihilation, it passes through the sphere unscathed'),
            ('touch_damage', {'dice': 8, 'sides': 10, 'type': 'force'}, 'isn’t wholly engulfed and obliterated by it takes 8d10 Force damage'),
            ('control.stationary_until_controlled', True, 'is stationary until someone takes control of it'),
            ('control.range_feet', 60, 'If you are within 60 feet of a sphere'),
            ('control.action', 'magic', 'you can take a Magic action to make a DC 25'),
            ('control.check_ability', 'intelligence', 'DC 25 Intelligence (Arcana) check'),
            ('control.check_skill', 'arcana', 'Intelligence (Arcana) check'),
            ('control.dc', 25, 'DC 25 Intelligence (Arcana) check'),
            ('control.lasts_until', 'start-of-your-next-turn', 'until the start of your next turn'),
            ('control.other_controller', 'loses-control', 'that creature loses control of the sphere'),
            ('control.failure_move_feet', 10, 'the sphere moves 10 feet toward you in a straight line'),
            ('move.action', 'bonus', 'you can take a Bonus Action to cause it to move'),
            ('move.feet_per_intelligence_modifier', 5, 'up to a number of feet equal to 5 times your Intelligence modifier'),
            ('move.minimum_feet', 5, '(minimum 5 feet)'),
            ('move.save_ability', 'dexterity', 'DC 19 Dexterity saving throw'),
            ('move.save_dc', 19, 'DC 19 Dexterity saving throw'),
            ('move.damage', {'dice': 8, 'sides': 10, 'type': 'force'}, 'or be touched by it, taking 8d10 Force damage'),
            ('move.at_zero_hit_points', 'obliterated-leaving-only-possessions', 'A creature reduced to 0 Hit Points by this damage is obliterated, leaving its possessions behind'),
            ('interaction_triggers[0]', 'planar-portal', 'comes into contact with a planar portal (such as that created by the Gate spell) or an extradimensional space'),
            ('interaction_triggers[1]', 'extradimensional-space', 'comes into contact with a planar portal (such as that created by the Gate spell) or an extradimensional space'),
            ('table_die_sides', 100, '1d100 Result'),
            ('interaction_table[0].min', 1, '01–50 The sphere is destroyed'),
            ('interaction_table[0].max', 50, '01–50 The sphere is destroyed'),
            ('interaction_table[0].result', 'sphere-destroyed', '01–50 The sphere is destroyed'),
            ('interaction_table[1].min', 51, '51–85 The sphere moves through the portal'),
            ('interaction_table[1].max', 85, '51–85 The sphere moves through the portal'),
            ('interaction_table[1].result', 'sphere-moves-through', '51–85 The sphere moves through the portal or into the extradimensional space'),
            ('interaction_table[2].min', 86, '86–00 A spatial rift'),
            ('interaction_table[2].max', 100, '86–00 A spatial rift'),
            ('interaction_table[2].result', 'spatial-rift-sends-all-to-random-plane', 'A spatial rift sends the sphere and each creature and object within 180 feet of the sphere to a random plane of existence'),
            ('rift_radius_feet', 180, 'within 180 feet of the sphere'),
        ],
    },
    'staff-of-charming': {
        'standard': [
            ('attunement_by_any_of[0]', 'bard', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[1]', 'cleric', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[2]', 'druid', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[3]', 'sorcerer', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[4]', 'warlock', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[5]', 'wizard', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('condition', 'while-holding', 'While holding the staff'),
            ('max_charges', 10, 'This staff has 10 charges'),
            ('cast_spell.charge_cost', 1, 'You can expend 1 of the staff’s charges to cast Charm Person'),
            ('cast_spell.spells[0]', 'charm-person', 'Charm Person'),
            ('cast_spell.spells[1]', 'command', 'Command, or'),
            ('cast_spell.spells[2]', 'comprehend-languages', 'Comprehend Languages'),
            ('cast_spell.save_dc_source', 'your-spell-save-dc', 'using your spell save DC'),
            ('reflect_enchantment.trigger', 'succeed-on-save-against-enchantment-targeting-only-you', 'If you succeed on a saving throw against an Enchantment spell that targets only you'),
            ('reflect_enchantment.action', 'reaction', 'you can take a Reaction to expend 1 charge'),
            ('reflect_enchantment.charge_cost', 1, 'to expend 1 charge from the staff'),
            ('reflect_enchantment.result', 'turn-spell-back-on-caster', 'turn the spell back on its caster as if you had cast the spell'),
            ('resist_enchantment.trigger', 'fail-save-against-enchantment-targeting-only-you', 'If you fail a saving throw against an Enchantment spell that targets only you'),
            ('resist_enchantment.result', 'failed-save-becomes-success', 'turn your failed save into a successful one'),
            ('resist_enchantment.reuse_after', 'next-dawn', 'again until the next dawn'),
            ('recharge.dice', {'dice': 1, 'sides': 8, 'type': 'charges'}, 'regains 1d8 + 2 expended charges daily at dawn'),
            ('recharge.bonus', 2, 'regains 1d8 + 2 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the staff crumbles to dust'),
            ('last_charge.result', 'crumbles-to-dust', 'the staff crumbles to dust and is destroyed'),
        ],
    },
    'staff-of-fire': {
        'standard': [
            ('attunement_by_any_of[0]', 'druid', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[1]', 'sorcerer', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[2]', 'warlock', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[3]', 'wizard', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('condition', 'while-holding', 'while you hold this staff'),
            ('resistance', 'fire', 'You have Resistance to Fire damage'),
            ('max_charges', 10, 'The staff has 10 charges'),
            ('save_dc_source', 'your-spell-save-dc', 'using your spell save DC'),
            ('spells[0].spell', 'burning-hands', 'Burning Hands 1'),
            ('spells[0].charge_cost', 1, 'Burning Hands 1'),
            ('spells[1].spell', 'fireball', 'Fireball 3'),
            ('spells[1].charge_cost', 3, 'Fireball 3'),
            ('spells[2].spell', 'wall-of-fire', 'Wall of Fire 4'),
            ('spells[2].charge_cost', 4, 'Wall of Fire 4'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.bonus', 4, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the staff crumbles into cinders'),
            ('last_charge.result', 'crumbles-into-cinders', 'the staff crumbles into cinders and is destroyed'),
        ],
    },
    'staff-of-frost': {
        'standard': [
            ('attunement_by_any_of[0]', 'druid', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[1]', 'sorcerer', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[2]', 'warlock', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[3]', 'wizard', 'Requires Attunement by a Druid, Sorcerer, Warlock, or Wizard'),
            ('condition', 'while-holding', 'while you hold this staff'),
            ('resistance', 'cold', 'You have Resistance to Cold damage'),
            ('max_charges', 10, 'The staff has 10 charges'),
            ('save_dc_source', 'your-spell-save-dc', 'using your spell save DC'),
            ('spells[0].spell', 'cone-of-cold', 'Cone of Cold 5'),
            ('spells[0].charge_cost', 5, 'Cone of Cold 5'),
            ('spells[1].spell', 'fog-cloud', 'Fog Cloud 1'),
            ('spells[1].charge_cost', 1, 'Fog Cloud 1'),
            ('spells[2].spell', 'ice-storm', 'Ice Storm 4'),
            ('spells[2].charge_cost', 4, 'Ice Storm 4'),
            ('spells[3].spell', 'wall-of-ice', 'Wall of Ice 4'),
            ('spells[3].charge_cost', 4, 'Wall of Ice 4'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.bonus', 4, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the staff turns to water'),
            ('last_charge.result', 'turns-to-water', 'the staff turns to water and is destroyed'),
        ],
    },
    'staff-of-healing': {
        'standard': [
            ('attunement_by_any_of[0]', 'bard', 'Requires Attunement by a Bard, Cleric, or Druid'),
            ('attunement_by_any_of[1]', 'cleric', 'Requires Attunement by a Bard, Cleric, or Druid'),
            ('attunement_by_any_of[2]', 'druid', 'Requires Attunement by a Bard, Cleric, or Druid'),
            ('condition', 'while-holding', 'While holding the staff'),
            ('max_charges', 10, 'This staff has 10 charges'),
            ('spell_modifier_source', 'your-spellcasting-ability-modifier', 'using your spellcasting ability modifier'),
            ('cure_wounds.charges_per_spell_level', 1, 'Cure Wounds 1 charge per spell level'),
            ('cure_wounds.maximum_spell_level', 4, '(maximum 4 for a level 4 spell)'),
            ('spells[0].spell', 'lesser-restoration', 'Lesser Restoration 2'),
            ('spells[0].charge_cost', 2, 'Lesser Restoration 2'),
            ('spells[1].spell', 'mass-cure-wounds', 'Mass Cure Wounds 5'),
            ('spells[1].charge_cost', 5, 'Mass Cure Wounds 5'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.bonus', 4, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the staff vanishes'),
            ('last_charge.result', 'vanishes-in-flash-of-light', 'the staff vanishes in a flash of light, lost forever'),
        ],
    },
    'staff-of-power': {
        'standard': [
            ('attunement_by_any_of[0]', 'sorcerer', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[1]', 'warlock', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[2]', 'wizard', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('max_charges', 20, 'This staff has 20 charges'),
            ('wielded_as', 'quarterstaff', 'wielded as a magic Quarterstaff'),
            ('attack_bonus', 2, 'magic Quarterstaff that grants a +2 bonus to attack rolls'),
            ('damage_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('holding_bonus', 2, 'While holding it, you gain a +2 bonus'),
            ('holding_bonus_applies_to[0]', 'armor-class', '+2 bonus to Armor Class, saving throws, and spell attack rolls'),
            ('holding_bonus_applies_to[1]', 'saving-throws', '+2 bonus to Armor Class, saving throws, and spell attack rolls'),
            ('holding_bonus_applies_to[2]', 'spell-attack-rolls', '+2 bonus to Armor Class, saving throws, and spell attack rolls'),
            ('save_dc_source', 'your-spell-save-dc', 'using your spell save DC'),
            ('spells[0].spell', 'cone-of-cold', 'Cone of Cold 5'),
            ('spells[0].printed_level', 'none', 'Cone of Cold 5'),
            ('spells[0].charge_cost', 5, 'Cone of Cold 5'),
            ('spells[1].spell', 'fireball', 'Fireball (level 5 version) 5'),
            ('spells[1].printed_level', 'level-5-version', 'Fireball (level 5 version) 5'),
            ('spells[1].charge_cost', 5, 'Fireball (level 5 version) 5'),
            ('spells[2].spell', 'globe-of-invulnerability', 'Globe of Invulnerability 6'),
            ('spells[2].printed_level', 'none', 'Globe of Invulnerability 6'),
            ('spells[2].charge_cost', 6, 'Globe of Invulnerability 6'),
            ('spells[3].spell', 'hold-monster', 'Hold Monster 5'),
            ('spells[3].printed_level', 'none', 'Hold Monster 5'),
            ('spells[3].charge_cost', 5, 'Hold Monster 5'),
            ('spells[4].spell', 'levitate', 'Levitate 2'),
            ('spells[4].printed_level', 'none', 'Levitate 2'),
            ('spells[4].charge_cost', 2, 'Levitate 2'),
            ('spells[5].spell', 'lightning-bolt', 'Lightning Bolt (level 5 version) 5'),
            ('spells[5].printed_level', 'level-5-version', 'Lightning Bolt (level 5 version) 5'),
            ('spells[5].charge_cost', 5, 'Lightning Bolt (level 5 version) 5'),
            ('spells[6].spell', 'magic-missile', 'Magic Missile 1'),
            ('spells[6].printed_level', 'none', 'Magic Missile 1'),
            ('spells[6].charge_cost', 1, 'Magic Missile 1'),
            ('spells[7].spell', 'ray-of-enfeeblement', 'Ray of Enfeeblement 1'),
            ('spells[7].printed_level', 'none', 'Ray of Enfeeblement 1'),
            ('spells[7].charge_cost', 1, 'Ray of Enfeeblement 1'),
            ('spells[8].spell', 'wall-of-force', 'Wall of Force 5'),
            ('spells[8].printed_level', 'none', 'Wall of Force 5'),
            ('spells[8].charge_cost', 5, 'Wall of Force 5'),
            ('recharge.dice', {'dice': 2, 'sides': 8, 'type': 'charges'}, 'regains 2d8 + 4 expended charges daily at dawn'),
            ('recharge.bonus', 4, 'regains 2d8 + 4 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.loses_properties_on', 1, 'On a 1, the staff retains its +2 bonus'),
            ('last_charge.retains', 'attack-and-damage-bonus-only', 'retains its +2 bonus to attack rolls and damage rolls but loses all other properties'),
            ('last_charge.regains_on', 20, 'On a 20, the staff regains'),
            ('last_charge.regains', {'dice': 1, 'sides': 8, 'type': 'charges'}, 'On a 20, the staff regains 1d8 + 2 charges'),
            ('last_charge.regains_bonus', 2, '1d8 + 2 charges'),
            ('retributive_strike.action', 'magic', 'You can take a Magic action to break the staff'),
            ('retributive_strike.area_emanation_feet', 30, 'fills a 30-foot Emanation'),
            ('retributive_strike.planar_travel_percent', 50, 'You have a 50 percent chance to instantly travel to a random plane of existence'),
            ('retributive_strike.your_damage_per_charge', 16, 'Force damage equal to 16 times the number of charges'),
            ('retributive_strike.other_damage_per_charge', 4, 'On a failed save, a creature takes Force damage equal to 4 times the number of charges'),
            ('retributive_strike.damage_type', 'force', 'you take Force damage'),
            ('retributive_strike.save_ability', 'dexterity', 'DC 17 Dexterity saving throw'),
            ('retributive_strike.save_dc', 17, 'DC 17 Dexterity saving throw'),
            ('retributive_strike.success_damage', 'half', 'a creature takes half as much damage'),
        ],
    },
    'staff-of-striking': {
        'standard': [
            ('wielded_as', 'quarterstaff', 'wielded as a magic Quarterstaff'),
            ('attack_bonus', 3, '+3 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 3, '+3 bonus to attack rolls and damage rolls'),
            ('max_charges', 10, 'The staff has 10 charges'),
            ('trigger', 'hit-with-melee-attack', 'When you hit with a melee attack using it'),
            ('maximum_charges_per_hit', 3, 'you can expend up to 3 charges'),
            ('extra_damage_per_charge', {'dice': 1, 'sides': 6, 'type': 'force'}, 'the target takes an extra 1d6 Force damage'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.bonus', 4, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the staff becomes'),
            ('last_charge.result', 'becomes-nonmagical-quarterstaff', 'the staff becomes a nonmagical Quarterstaff'),
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
        raise SystemExit("usage: verify_srd_items_242.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_242.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 242]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 242 and catalogue["batch"] == "data-magic-items-18"
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
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 243-246")


main()
