#!/usr/bin/env python3
"""Compare issue #243 item identities, headers, table rows and values with the pinned PDF."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = "8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87"
PAGES = (246, 247, 248)
RARITY = {"common": "Common", "uncommon": "Uncommon", "rare": "Rare", "very-rare": "Very Rare", "legendary": "Legendary"}
CATEGORY = {"armor": "Armor (", "potion": "Potion, ", "ring": "Ring, ", "rod": "Rod, ", "scroll": "Scroll, ", "staff": "Staff, ", "weapon": "Weapon (", "wondrous-item": "Wondrous Item, "}
FOOTER = re.compile(r"\d+ System Reference Document 5\.2\.1")
# A category line under a heading. Each family block must hold only its own.
CATEGORY_LINE = re.compile(r"(?:Wondrous Item|Potion|Ring|Rod|Staff|Wand|Scroll|(?:Weapon|Armor) \([^)]*\)), (?:Common|Uncommon|Rare|Very Rare|Legendary|Artifact|Rarity Varies)")
NUMBER_WORDS = {"cantrip": 0, "a": 1, "an": 1, "once": 1, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10}
# The headings that follow the batch's last families in print.
FOLLOWING = ['Talisman of Pure Good Wondrous']

# Each family's printed category line.
HEADERS = {
    'staff-of-swarming-insects': 'Staff, Rare (Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard)',
    'staff-of-thunder-and-lightning': 'Staff, Very Rare (Requires Attunement)',
    'staff-of-withering': 'Staff, Rare (Requires Attunement)',
    'staff-of-the-magi': 'Staff, Legendary (Requires Attunement by a Sorcerer, Warlock, or Wizard)',
    'staff-of-the-python': 'Staff, Uncommon (Requires Attunement)',
    'staff-of-the-woodlands': 'Staff, Rare (Requires Attunement by a Druid)',
    'stone-of-controlling-earth-elementals': 'Wondrous Item, Rare',
    'stone-of-good-luck-luckstone': 'Wondrous Item, Uncommon (Requires Attunement)',
    'sun-blade': 'Weapon (Longsword), Rare (Requires Attunement)',
    'sword-of-life-stealing': 'Weapon (Glaive, Greatsword, Longsword, Rapier, Scimitar, or Shortsword), Rare (Requires Attunement)',
    'sword-of-sharpness': 'Weapon (Glaive, Greatsword, Longsword, or Scimitar), Very Rare (Requires Attunement)',
    'sword-of-wounding': 'Weapon (Glaive, Greatsword, Longsword, Rapier, Scimitar, or Shortsword), Rare (Requires Attunement)',
}

# Printed charge sentence of each family with a charge pool.
CHARGES = {
    'staff-of-swarming-insects': ({'mode': 'limited', 'maximum': 10, 'per_use': 1, 'recharge': 'regain-1d6-plus-4-at-dawn'}, 'This staff has 10 charges'),
    'staff-of-withering': ({'mode': 'limited', 'maximum': 3, 'per_use': 1, 'recharge': 'regain-1d3-at-dawn'}, 'This staff has 3 charges and regains 1d3 expended charges daily at dawn'),
    'staff-of-the-magi': ({'mode': 'limited', 'maximum': 50, 'per_use': 1, 'recharge': 'regain-4d6-plus-2-at-dawn'}, 'This staff has 50 charges'),
    'staff-of-the-woodlands': ({'mode': 'limited', 'maximum': 6, 'per_use': 1, 'recharge': 'regain-1d6-at-dawn'}, 'This staff has 6 charges'),
}

# Every payload leaf of every variant: (path, value, printed phrase). The
# catalogued payload must have exactly these leaves, and each phrase must be
# printed inside the family's own text. A "none" key or an empty list records
# that nothing is printed. source_pages is checked from the page layout.
VALUES = {
    'staff-of-swarming-insects': {
        'standard': [
            ('attunement_by_any_of[0]', 'bard', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[1]', 'cleric', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[2]', 'druid', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[3]', 'sorcerer', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[4]', 'warlock', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[5]', 'wizard', 'Requires Attunement by a Bard, Cleric, Druid, Sorcerer, Warlock, or Wizard'),
            ('max_charges', 10, 'This staff has 10 charges'),
            ('insect_cloud.action', 'magic', 'you can take a Magic action and expend 1 charge'),
            ('insect_cloud.charge_cost', 1, 'expend 1 charge to cause a swarm of harmless flying insects'),
            ('insect_cloud.swarm', 'harmless-flying-insects', 'a swarm of harmless flying insects'),
            ('insect_cloud.emanation_feet', 30, 'fill a 30-foot Emanation'),
            ('insect_cloud.duration_minutes', 10, 'The insects remain for 10 minutes'),
            ('insect_cloud.area', 'heavily-obscured-for-creatures-other-than-you', 'Heavily Obscured for creatures other than you'),
            ('insect_cloud.dispersed_by', 'strong-wind', 'A strong wind (like that created by Gust of Wind) disperses the swarm'),
            ('save_and_attack_source', 'your-spell-save-dc-and-spell-attack-modifier', 'using your spell save DC and spell attack modifier'),
            ('spells[0].spell', 'giant-insect', 'Giant Insect 4'),
            ('spells[0].charge_cost', 4, 'Giant Insect 4'),
            ('spells[1].spell', 'insect-plague', 'Insect Plague 5'),
            ('spells[1].charge_cost', 5, 'Insect Plague 5'),
            ('recharge.dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.bonus', 4, 'regains 1d6 + 4 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, a swarm of insects consumes and destroys the staff'),
            ('last_charge.result', 'insect-swarm-consumes-staff-then-disperses', 'a swarm of insects consumes and destroys the staff, then disperses'),
        ],
    },
    'staff-of-thunder-and-lightning': {
        'standard': [
            ('wielded_as', 'quarterstaff', 'wielded as a magic Quarterstaff'),
            ('attack_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('each_property_reuse_after', 'next-dawn', 'Once one of these properties is used, it can’t be used again until the next dawn'),
            ('lightning.trigger', 'hit-with-melee-attack', 'Lightning. When you hit with a melee attack using the staff'),
            ('lightning.extra_damage', {'dice': 2, 'sides': 6, 'type': 'lightning'}, 'an extra 2d6 Lightning damage'),
            ('lightning.action_required', False, 'an extra 2d6 Lightning damage (no action required)'),
            ('thunder.trigger', 'hit-with-melee-attack', 'Thunder. When you hit with a melee attack using the staff'),
            ('thunder.audible_feet', 300, 'crack of thunder audible out to 300 feet'),
            ('thunder.action_required', False, 'audible out to 300 feet (no action required)'),
            ('thunder.save_ability', 'constitution', 'The target you hit must succeed on a DC 17 Constitution'),
            ('thunder.save_dc', 17, 'The target you hit must succeed on a DC 17'),
            ('thunder.condition', 'stunned', 'or have the Stunned condition'),
            ('thunder.condition_ends', 'end-of-your-next-turn', 'until the end of your next turn'),
            ('thunder_and_lightning.trigger', 'immediately-after-melee-hit', 'Immediately after you hit with a melee attack using the staff'),
            ('thunder_and_lightning.action', 'bonus', 'you can take a Bonus Action to use the Lightning and Thunder properties'),
            ('thunder_and_lightning.uses[0]', 'lightning', 'use the Lightning and Thunder properties (see above) at the same time'),
            ('thunder_and_lightning.uses[1]', 'thunder', 'use the Lightning and Thunder properties (see above) at the same time'),
            ('thunder_and_lightning.spends_their_daily_uses', False, 'Doing so doesn’t expend the daily use of those properties, only the use of this one'),
            ('lightning_strike.action', 'magic', 'You can take a Magic action to cause a bolt of lightning'),
            ('lightning_strike.line_width_feet', 5, 'in a Line that is 5 feet wide'),
            ('lightning_strike.line_length_feet', 120, 'and 120 feet long'),
            ('lightning_strike.save_ability', 'dexterity', 'DC 17 Dexterity saving throw, taking 9d6'),
            ('lightning_strike.save_dc', 17, 'DC 17 Dexterity saving throw, taking 9d6'),
            ('lightning_strike.damage', {'dice': 9, 'sides': 6, 'type': 'lightning'}, 'taking 9d6 Lightning damage'),
            ('lightning_strike.success_damage', 'half', 'or half as much damage on a successful one'),
            ('thunderclap.action', 'magic', 'You can take a Magic action to cause the staff to produce a thunderclap'),
            ('thunderclap.audible_feet', 600, 'audible out to 600 feet'),
            ('thunderclap.emanation_feet', 60, 'within a 60-foot Emanation'),
            ('thunderclap.save_ability', 'constitution', 'DC 17 Constitution saving throw. On a failed save, a creature takes 2d6'),
            ('thunderclap.save_dc', 17, 'DC 17 Constitution saving throw. On a failed save, a creature takes 2d6'),
            ('thunderclap.damage', {'dice': 2, 'sides': 6, 'type': 'thunder'}, 'takes 2d6 Thunder damage'),
            ('thunderclap.failure_condition', 'deafened', 'has the Deafened condition'),
            ('thunderclap.condition_minutes', 1, 'the Deafened condition for 1 minute'),
            ('thunderclap.success', 'half-damage-only', 'On a successful save, a creature takes half as much damage only'),
        ],
    },
    'staff-of-withering': {
        'standard': [
            ('wielded_as', 'quarterstaff', 'wielded as a magic Quarterstaff'),
            ('hit_damage', 'as-normal-quarterstaff', 'On a hit, it deals damage as a normal Quarterstaff'),
            ('max_charges', 3, 'This staff has 3 charges'),
            ('recharge_dice', {'dice': 1, 'sides': 3, 'type': 'charges'}, 'regains 1d3 expended charges daily at dawn'),
            ('charge_cost', 1, 'you can expend 1 charge to deal an extra 2d10'),
            ('extra_damage', {'dice': 2, 'sides': 10, 'type': 'necrotic'}, 'an extra 2d10 Necrotic damage'),
            ('save_ability', 'constitution', 'DC 15 Constitution saving throw'),
            ('save_dc', 15, 'DC 15 Constitution saving throw'),
            ('failure.mode', 'disadvantage', 'the target has Disadvantage for 1 hour'),
            ('failure.duration_hours', 1, 'Disadvantage for 1 hour'),
            ('failure.applies_to[0]', 'strength-checks-and-saves', 'on any ability check or saving throw that uses Strength or Constitution'),
            ('failure.applies_to[1]', 'constitution-checks-and-saves', 'on any ability check or saving throw that uses Strength or Constitution'),
        ],
    },
    'staff-of-the-magi': {
        'standard': [
            ('attunement_by_any_of[0]', 'sorcerer', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[1]', 'warlock', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('attunement_by_any_of[2]', 'wizard', 'Requires Attunement by a Sorcerer, Warlock, or Wizard'),
            ('max_charges', 50, 'This staff has 50 charges'),
            ('wielded_as', 'quarterstaff', 'wielded as a magic Quarterstaff'),
            ('attack_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('spell_attack_bonus', 2, '+2 bonus to spell attack rolls'),
            ('spell_absorption.save_mode', 'advantage', 'you have Advantage on saving throws against spells'),
            ('spell_absorption.trigger', 'spell-from-another-creature-targets-only-you', 'when another creature casts a spell that targets only you'),
            ('spell_absorption.action', 'reaction', 'you can take a Reaction when another creature casts a spell'),
            ('spell_absorption.result', 'spell-canceled-and-absorbed', 'the staff absorbs the magic of the spell, canceling its effect'),
            ('spell_absorption.charges_gained', 'absorbed-spell-level', 'gaining a number of charges equal to the absorbed spell’s level'),
            ('spell_absorption.explodes_above_charges', 50, 'brings the staff’s total number of charges above 50'),
            ('spell_absorption.explosion', 'as-retributive-strike', 'the staff explodes as if you activated its Retributive Strike'),
            ('save_dc_source', 'your-spell-save-dc', 'using your spell save DC'),
            ('spells[0].spell', 'arcane-lock', 'Arcane Lock 0'),
            ('spells[0].printed_level', 'none', 'Arcane Lock 0'),
            ('spells[0].charge_cost', 0, 'Arcane Lock 0'),
            ('spells[1].spell', 'conjure-elemental', 'Conjure Elemental 7'),
            ('spells[1].printed_level', 'none', 'Conjure Elemental 7'),
            ('spells[1].charge_cost', 7, 'Conjure Elemental 7'),
            ('spells[2].spell', 'detect-magic', 'Detect Magic 0'),
            ('spells[2].printed_level', 'none', 'Detect Magic 0'),
            ('spells[2].charge_cost', 0, 'Detect Magic 0'),
            ('spells[3].spell', 'dispel-magic', 'Dispel Magic 3'),
            ('spells[3].printed_level', 'none', 'Dispel Magic 3'),
            ('spells[3].charge_cost', 3, 'Dispel Magic 3'),
            ('spells[4].spell', 'enlarge-reduce', 'Enlarge/Reduce 0'),
            ('spells[4].printed_level', 'none', 'Enlarge/Reduce 0'),
            ('spells[4].charge_cost', 0, 'Enlarge/Reduce 0'),
            ('spells[5].spell', 'fireball', 'Fireball (level 7 version) 7'),
            ('spells[5].printed_level', 'level-7-version', 'Fireball (level 7 version) 7'),
            ('spells[5].charge_cost', 7, 'Fireball (level 7 version) 7'),
            ('spells[6].spell', 'flaming-sphere', 'Flaming Sphere 2'),
            ('spells[6].printed_level', 'none', 'Flaming Sphere 2'),
            ('spells[6].charge_cost', 2, 'Flaming Sphere 2'),
            ('spells[7].spell', 'ice-storm', 'Ice Storm 4'),
            ('spells[7].printed_level', 'none', 'Ice Storm 4'),
            ('spells[7].charge_cost', 4, 'Ice Storm 4'),
            ('spells[8].spell', 'invisibility', 'Invisibility 2'),
            ('spells[8].printed_level', 'none', 'Invisibility 2'),
            ('spells[8].charge_cost', 2, 'Invisibility 2'),
            ('spells[9].spell', 'knock', 'Knock 2'),
            ('spells[9].printed_level', 'none', 'Knock 2'),
            ('spells[9].charge_cost', 2, 'Knock 2'),
            ('spells[10].spell', 'light', 'Light 0'),
            ('spells[10].printed_level', 'none', 'Light 0'),
            ('spells[10].charge_cost', 0, 'Light 0'),
            ('spells[11].spell', 'lightning-bolt', 'Lightning Bolt (level 7 version) 7'),
            ('spells[11].printed_level', 'level-7-version', 'Lightning Bolt (level 7 version) 7'),
            ('spells[11].charge_cost', 7, 'Lightning Bolt (level 7 version) 7'),
            ('spells[12].spell', 'mage-hand', 'Mage Hand 0'),
            ('spells[12].printed_level', 'none', 'Mage Hand 0'),
            ('spells[12].charge_cost', 0, 'Mage Hand 0'),
            ('spells[13].spell', 'passwall', 'Passwall 5'),
            ('spells[13].printed_level', 'none', 'Passwall 5'),
            ('spells[13].charge_cost', 5, 'Passwall 5'),
            ('spells[14].spell', 'plane-shift', 'Plane Shift 7'),
            ('spells[14].printed_level', 'none', 'Plane Shift 7'),
            ('spells[14].charge_cost', 7, 'Plane Shift 7'),
            ('spells[15].spell', 'protection-from-evil-and-good', 'Protection from Evil and Good 0'),
            ('spells[15].printed_level', 'none', 'Protection from Evil and Good 0'),
            ('spells[15].charge_cost', 0, 'Protection from Evil and Good 0'),
            ('spells[16].spell', 'telekinesis', 'Telekinesis 5'),
            ('spells[16].printed_level', 'none', 'Telekinesis 5'),
            ('spells[16].charge_cost', 5, 'Telekinesis 5'),
            ('spells[17].spell', 'wall-of-fire', 'Wall of Fire 4'),
            ('spells[17].printed_level', 'none', 'Wall of Fire 4'),
            ('spells[17].charge_cost', 4, 'Wall of Fire 4'),
            ('spells[18].spell', 'web', 'Web 2'),
            ('spells[18].printed_level', 'none', 'Web 2'),
            ('spells[18].charge_cost', 2, 'Web 2'),
            ('recharge.dice', {'dice': 4, 'sides': 6, 'type': 'charges'}, 'regains 4d6 + 2 expended charges daily at dawn'),
            ('recharge.bonus', 2, 'regains 4d6 + 2 expended charges daily at dawn'),
            ('recharge.when', 'daily-at-dawn', 'expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.regains_on', 20, 'On a 20, the staff regains'),
            ('last_charge.regains', {'dice': 1, 'sides': 12, 'type': 'charges'}, 'On a 20, the staff regains 1d12 + 1 charges'),
            ('last_charge.regains_bonus', 1, '1d12 + 1 charges'),
            ('retributive_strike.action', 'magic', 'You can take a Magic action to break the staff'),
            ('retributive_strike.area_emanation_feet', 30, 'fills a 30-foot Emanation'),
            ('retributive_strike.planar_travel_percent', 50, 'You have a 50 percent chance to instantly travel to a random plane of existence'),
            ('retributive_strike.your_damage_per_charge', 16, 'Force damage equal to 16 times the number of charges'),
            ('retributive_strike.other_damage_per_charge', 6, 'a creature takes Force damage equal to 6 times the number of charges'),
            ('retributive_strike.damage_type', 'force', 'you take Force damage'),
            ('retributive_strike.save_ability', 'dexterity', 'DC 17 Dexterity saving throw'),
            ('retributive_strike.save_dc', 17, 'DC 17 Dexterity saving throw'),
            ('retributive_strike.success_damage', 'half', 'a creature takes half as much damage'),
        ],
    },
    'staff-of-the-python': {
        'standard': [
            ('action', 'magic', 'As a Magic action, you can throw this staff'),
            ('throw_range_feet', 10, 'an unoccupied space within 10 feet of you'),
            ('creature', 'giant-constrictor-snake', 'causing the staff to become a Giant Constrictor Snake'),
            ('control', 'under-your-control', 'The snake is under your control'),
            ('initiative', 'shares-yours-and-acts-immediately-after', 'shares your Initiative count, taking its turn immediately after yours'),
            ('command.timing', 'your-turn', 'On your turn, you can mentally command the snake'),
            ('command.method', 'mentally', 'you can mentally command the snake'),
            ('command.action_required', False, '(no action required)'),
            ('command.range_feet', 60, 'if it is within 60 feet of you'),
            ('command.blocked_by_condition', 'incapacitated', 'you don’t have the Incapacitated condition'),
            ('command.options[0]', 'decide-action-and-movement', 'You decide what action the snake takes and where it moves during its turn, or you can issue it a general command'),
            ('command.options[1]', 'general-command', 'You decide what action the snake takes and where it moves during its turn, or you can issue it a general command'),
            ('command.without_commands', 'defends-itself', 'Absent commands from you, the snake defends itself'),
            ('revert.action', 'bonus', 'As a Bonus Action, you can command the snake to revert'),
            ('revert.location', 'snakes-current-space', 'revert to staff form in its current space'),
            ('revert.reuse_after_hours', 1, 'you can’t use the staff’s property again for 1 hour'),
            ('at_zero_hit_points', 'dies-reverts-then-staff-shatters', 'If the snake is reduced to 0 Hit Points, it dies and reverts to its staff form; the staff then shatters and is destroyed'),
            ('hit_points_on_early_revert', 'all-regained', 'If the snake reverts to staff form before losing all its Hit Points, it regains all of them'),
        ],
    },
    'staff-of-the-woodlands': {
        'standard': [
            ('attunement_by_any_of[0]', 'druid', '(Requires Attunement by a Druid)'),
            ('max_charges', 6, 'This staff has 6 charges'),
            ('wielded_as', 'quarterstaff', 'wielded as a magic Quarterstaff'),
            ('attack_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('spell_attack_bonus', 2, 'you have a +2 bonus to spell attack rolls'),
            ('save_dc_source', 'your-spell-save-dc', 'using your spell save DC'),
            ('spells[0].spell', 'animal-friendship', 'Animal Friendship 1'),
            ('spells[0].charge_cost', 1, 'Animal Friendship 1'),
            ('spells[1].spell', 'awaken', 'Awaken 5'),
            ('spells[1].charge_cost', 5, 'Awaken 5'),
            ('spells[2].spell', 'barkskin', 'Barkskin 2'),
            ('spells[2].charge_cost', 2, 'Barkskin 2'),
            ('spells[3].spell', 'locate-animals-or-plants', 'Locate Animals or Plants 2'),
            ('spells[3].charge_cost', 2, 'Locate Animals or Plants 2'),
            ('spells[4].spell', 'pass-without-trace', 'Pass without Trace 2'),
            ('spells[4].charge_cost', 2, 'Pass without Trace 2'),
            ('spells[5].spell', 'speak-with-animals', 'Speak with Animals 1'),
            ('spells[5].charge_cost', 1, 'Speak with Animals 1'),
            ('spells[6].spell', 'speak-with-plants', 'Speak with Plants 3'),
            ('spells[6].charge_cost', 3, 'Speak with Plants 3'),
            ('spells[7].spell', 'wall-of-thorns', 'Wall of Thorns 6'),
            ('spells[7].charge_cost', 6, 'Wall of Thorns 6'),
            ('tree_form.action', 'magic', 'You can take a Magic action to plant one end of the staff in earth'),
            ('tree_form.charge_cost', 1, 'and expend 1 charge to transform the staff'),
            ('tree_form.placement', 'one-end-in-earth-in-unoccupied-space', 'plant one end of the staff in earth in an unoccupied space'),
            ('tree_form.tree_height_feet', 60, 'The tree is 60 feet tall'),
            ('tree_form.trunk_diameter_feet', 5, 'has a 5-foot-diameter trunk'),
            ('tree_form.branch_radius_feet', 20, 'spread out in a 20-foot radius'),
            ('tree_form.aura', 'faint-transmutation-discerned-with-detect-magic', 'radiates a faint aura of Transmutation magic that can be discerned with the Detect Magic spell'),
            ('tree_form.revert', 'touch-tree-and-use-magic-action', 'While touching the tree and using a Magic action, you return the staff to its normal form'),
            ('tree_form.creatures_in_tree', 'fall', 'Any creature in the tree falls when the tree reverts'),
            ('recharge_dice', {'dice': 1, 'sides': 6, 'type': 'charges'}, 'regains 1d6 expended charges daily at dawn'),
            ('last_charge.die_sides', 20, 'If you expend the last charge, roll 1d20'),
            ('last_charge.destroyed_on', 1, 'On a 1, the staff loses its properties'),
            ('last_charge.result', 'becomes-nonmagical-quarterstaff', 'the staff loses its properties and becomes a nonmagical Quarterstaff'),
        ],
    },
    'stone-of-controlling-earth-elementals': {
        'standard': [
            ('weight_pounds', 5, 'this 5-pound stone'),
            ('condition', 'while-touching-stone-to-ground', 'While touching this 5-pound stone to the ground'),
            ('action', 'magic', 'you can take a Magic action to summon an Earth Elemental'),
            ('creature', 'earth-elemental', 'summon an Earth Elemental'),
            ('summon_range_feet', 30, 'an unoccupied space you choose within 30 feet of yourself'),
            ('obeys_commands', True, 'obeys your commands'),
            ('initiative', 'acts-immediately-after-you-on-your-count', 'takes its turn immediately after you on your Initiative count'),
            ('duration_hours', 1, 'The elemental disappears after 1 hour'),
            ('also_disappears_when[0]', 'it-dies', 'when it dies'),
            ('also_disappears_when[1]', 'dismissed-as-bonus-action', 'or when you dismiss it as a Bonus Action'),
            ('reuse_after', 'next-dawn', 'can’t be used this way again until the next dawn'),
        ],
    },
    'stone-of-good-luck-luckstone': {
        'standard': [
            ('appearance', 'polished-agate', 'this polished agate'),
            ('condition', 'while-on-your-person', 'While this polished agate is on your person'),
            ('ability_check_bonus', 1, '+1 bonus to ability checks and saving throws'),
            ('saving_throw_bonus', 1, '+1 bonus to ability checks and saving throws'),
        ],
    },
    'sun-blade': {
        'standard': [
            ('appearance', 'sword-hilt', 'appears to be a sword hilt'),
            ('blade.condition', 'while-grasping-hilt', 'While grasping the hilt'),
            ('blade.action', 'bonus', 'you can take a Bonus Action to cause a blade of pure radiance'),
            ('blade.options[0]', 'appear', 'to spring into existence or make the blade disappear'),
            ('blade.options[1]', 'disappear', 'to spring into existence or make the blade disappear'),
            ('blade.functions_as', 'longsword', 'this magic weapon functions as a Longsword'),
            ('blade.property', 'finesse', 'with the Finesse property'),
            ('blade.proficiency_if_proficient_with[0]', 'longswords', 'If you are proficient with Longswords or Shortswords, you are proficient with the Sun Blade'),
            ('blade.proficiency_if_proficient_with[1]', 'shortswords', 'If you are proficient with Longswords or Shortswords, you are proficient with the Sun Blade'),
            ('attack_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('damage_bonus', 2, '+2 bonus to attack rolls and damage rolls'),
            ('damage_type', 'radiant', 'deals Radiant damage instead of Slashing damage'),
            ('replaces_damage_type', 'slashing', 'instead of Slashing damage'),
            ('undead_extra_damage', {'dice': 1, 'sides': 8, 'type': 'radiant'}, 'an extra 1d8 Radiant damage'),
            ('sunlight.bright_radius_feet', 15, 'Bright Light in a 15-foot radius'),
            ('sunlight.dim_additional_feet', 15, 'Dim Light for an additional 15 feet'),
            ('sunlight.light_is_sunlight', True, 'The light is sunlight'),
            ('sunlight.adjust_action', 'magic', 'you can take a Magic action to expand or reduce its radius'),
            ('sunlight.adjust_step_feet', 5, 'by 5 feet each'),
            ('sunlight.maximum_each_feet', 30, 'to a maximum of 30 feet each'),
            ('sunlight.minimum_each_feet', 10, 'or a minimum of 10 feet each'),
        ],
    },
    'sword-of-life-stealing': {
        'standard': [
            ('trigger_natural_roll', 20, 'roll a 20 on the d20 for the attack roll'),
            ('target', 'creature', 'When you attack a creature with this magic weapon'),
            ('extra_damage', 15, 'an extra 15 Necrotic damage'),
            ('damage_type', 'necrotic', '15 Necrotic damage'),
            ('excluded_creature_types[0]', 'construct', 'if it isn’t a Construct or an Undead'),
            ('excluded_creature_types[1]', 'undead', 'if it isn’t a Construct or an Undead'),
            ('temporary_hit_points', 'equal-to-necrotic-damage-taken', 'you gain Temporary Hit Points equal to the amount of Necrotic damage taken'),
        ],
    },
    'sword-of-sharpness': {
        'standard': [
            ('object_hit', 'maximize-weapon-damage-dice', 'When you attack an object with this magic weapon and hit, maximize your weapon damage dice'),
            ('trigger_natural_roll', 20, 'roll a 20 on the d20 for the attack roll'),
            ('target', 'creature', 'When you attack a creature with this weapon'),
            ('extra_damage', 14, 'an extra 14 Slashing damage'),
            ('damage_type', 'slashing', '14 Slashing damage'),
            ('applies_condition', 'exhaustion', 'gains 1 Exhaustion level'),
            ('exhaustion_levels', 1, 'gains 1 Exhaustion level'),
        ],
    },
    'sword-of-wounding': {
        'standard': [
            ('trigger', 'hit-creature-with-attack', 'When you hit a creature with an attack using this magic weapon'),
            ('extra_damage', {'dice': 2, 'sides': 6, 'type': 'necrotic'}, 'an extra 2d6 Necrotic damage'),
            ('save_ability', 'constitution', 'DC 15 Constitution saving throw'),
            ('save_dc', 15, 'DC 15 Constitution saving throw'),
            ('failure', 'cannot-regain-hit-points', 'or be unable to regain Hit Points'),
            ('duration_hours', 1, 'unable to regain Hit Points for 1 hour'),
            ('repeat_save', 'end-of-each-of-its-turns', 'The target repeats the save at the end of each of its turns'),
            ('repeat_success', 'ends-effect-on-itself', 'ending the effect on itself on a success'),
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
        raise SystemExit("usage: verify_srd_items_243.py /path/to/SRD_CC_v5.2.1.pdf")
    pdf = Path(sys.argv[1])
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, "PDF pin mismatch"
    # Columns are read independently and chained in reading order across pages.
    stream, spans = "", {}
    for page in PAGES:
        start = len(stream)
        stream += " ".join(clean(text(pdf, page, x)).strip() for x in (0, 306))
        spans[page] = (start, len(stream))
        stream += " "

    catalogue = json.loads((ROOT / "internal/rules/srd521/testdata/items_243.json").read_text())
    inventory = json.loads((ROOT / "docs/investigations/issue-80-spell-item-inventory.json").read_text())
    expected = [(r["id"], r["name"], r["source_page"]) for r in inventory["records"] if r.get("catalogue_issue") == 243]
    assert len(expected) == 12
    families = catalogue["families"]
    assert catalogue["source"]["sha256"] == PIN and catalogue["catalogue_issue"] == 243 and catalogue["batch"] == "data-magic-items-19"
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
    print(f"ok: {len(families)} families, {sum(len(f['variants']) for f in families)} variants, {count} printed values, pp. 246-248")


main()
