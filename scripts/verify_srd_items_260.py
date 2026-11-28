#!/usr/bin/env python3
"""Compare issue #260 adventuring gear records with the pinned PDF.

srd_gear.py builds weights, prices, variant rows, pack contents, actions,
checks, damage and spellcasting classes from the printed text. The purposes,
measure names and trait keys are hand-written; EVIDENCE ties each measure and
trait to a phrase the entry must print.
"""
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))
import srd_gear  # noqa: E402

EVIDENCE = {
    'magnifying-glass': {
        'measures': [
            ('fire-ignition-time', '5', 'minutes', 'about 5 minutes for the fire to ignite'),
        ],
        'traits': [
            ('advantage-on-checks-to-appraise-or-inspect-highly-detailed-item', 'grants Advantage on any ability check made to appraise or inspect a highly detailed item'),
            ('lighting-fire-requires-sunlight-bright-light-and-tinder', 'Lighting a fire with a Magnifying Glass requires light as bright as sunlight to focus, tinder to ignite'),
        ],
    },
    'manacles': {
        'checks': ['bind-creature', 'escape', 'burst', 'pick-lock-without-key'],
        'measures': [
            ('bind-range', '5', 'feet', 'within 5 feet of yourself'),
        ],
        'traits': [
            ('target-unwilling-small-or-medium-grappled-incapacitated-or-restrained-creature', 'bind an unwilling Small or Medium creature within 5 feet of yourself that has the Grappled, Incapacitated, or Restrained condition'),
            ('bound-creature-has-disadvantage-on-attack-rolls', 'a creature has Disadvantage on attack rolls'),
            ('restrained-if-attached-to-fixed-chain-or-hook', 'the creature is Restrained if the Manacles are attached to a chain or hook that is fixed in place'),
            ('comes-with-key', 'comes with a key'),
            ('thieves-tools-can-pick-without-key', 'Without the key, a creature can use Thieves’ Tools to pick the Manacles’ lock'),
        ],
    },
    'map': {
        'measures': [
            ('wayfinding-check-bonus', '5', 'bonus', '+5 bonus'),
        ],
        'traits': [
            ('bonus-requires-consulting-accurate-map', 'If you consult an accurate Map'),
            ('bonus-to-wisdom-survival-checks-to-find-way-in-mapped-place', 'to Wisdom (Survival) checks you make to find your way in the place represented on it'),
        ],
    },
    'mirror': {
        'traits': [
            ('handheld-steel', 'A handheld steel Mirror'),
            ('useful-for-cosmetics-peeking-around-corners-and-signals', 'personal cosmetics but also for peeking around corners and reflecting light as a signal'),
        ],
    },
    'net': {
        'checks': ['avoid-restrained', 'escape'],
        'measures': [
            ('throw-range', '15', 'feet', 'within 15 feet of yourself'),
            ('escape-helper-range', '5', 'feet', 'a creature within 5 feet of it'),
            ('armor-class', '10', 'armor-class', 'AC 10'),
            ('hit-points', '5', 'hit-points', '5 HP'),
        ],
        'traits': [
            ('throw-net-in-place-of-one-attack', 'replace one of your attacks with throwing a Net'),
            ('target-one-creature-you-can-see', 'Target a creature you can see'),
            ('failed-save-restrains-until-escape', 'have the Restrained condition until it escapes'),
            ('huge-or-larger-target-succeeds-automatically', 'The target succeeds automatically if it is Huge or larger'),
            ('successful-check-frees-restrained-creature', 'freeing the Restrained creature on a success'),
            ('immune-to-bludgeoning-poison-and-psychic-damage', 'Immunity to Bludgeoning, Poison, and Psychic damage'),
            ('destroying-net-frees-target-and-ends-effect', 'also frees the target, ending the effect'),
        ],
    },
    'oil': {
        'checks': ['avoid-being-covered-in-oil'],
        'damage': ['extra-from-burning-oil', 'burning-area'],
        'measures': [
            ('throw-range', '20', 'feet', 'within 20 feet of yourself'),
            ('drying-time', '1', 'minute', 'after 1 minute'),
            ('covered-square-side', '5', 'feet', '5-foot-square area'),
            ('pour-range', '5', 'feet', 'within 5 feet of yourself'),
            ('area-burn-rounds', '2', 'rounds', '2 rounds from when the oil was lit'),
            ('area-burn-seconds', '12', 'seconds', 'or 12 seconds'),
            ('fuel-burn-duration', '6', 'hours', 'burns for 6 hours in a Lamp or Lantern'),
        ],
        'traits': [
            ('douse-creature-object-or-space-or-use-as-fuel', 'You can douse a creature, object, or space with Oil or use it as fuel'),
            ('throw-flask-in-place-of-one-attack', 'replace one of your attacks with throwing an Oil flask'),
            ('target-one-creature-or-object', 'Target one creature or object'),
            ('extra-damage-if-fire-damage-taken-before-oil-dries', 'If the target takes Fire damage before the oil dries'),
            ('pour-flask-on-level-ground', 'pour an Oil flask on level ground'),
            ('area-damages-on-entry-or-turn-end-once-per-turn', 'A creature can take this damage only once per turn'),
            ('fuel-for-lamps-and-lanterns', 'Oil serves as fuel for Lamps and Lanterns'),
            ('fuel-duration-need-not-be-consecutive', 'That duration doesn’t need to be consecutive'),
            ('extinguish-and-rekindle-until-total-duration-burned', 'rekindle it again until it has burned for a total of 6 hours'),
        ],
    },
    'paper': {
        'measures': [
            ('handwritten-words-per-sheet', '250', 'words', 'about 250 handwritten words'),
        ],
    },
    'parchment': {
        'measures': [
            ('handwritten-words-per-sheet', '250', 'words', 'about 250 handwritten words'),
        ],
    },
    'perfume': {
        'measures': [
            ('vial-volume', '4', 'ounces', '4-ounce vial'),
            ('benefit-duration', '1', 'hour', 'For 1 hour after applying Perfume'),
            ('influence-range', '5', 'feet', 'within 5 feet of you'),
        ],
        'traits': [
            ('advantage-on-charisma-persuasion-to-influence-indifferent-humanoid', 'Advantage on Charisma (Persuasion) checks made to influence an Indifferent Humanoid'),
        ],
    },
    'poison-basic': {
        'damage': ['extra-from-poisoned-weapon-or-ammunition'],
        'measures': [
            ('coated-weapons', '1', 'weapon', 'coat one weapon'),
            ('coated-ammunition', '3', 'pieces', 'up to three pieces of ammunition'),
            ('potency-duration', '1', 'minute', 'retains potency for 1 minute'),
        ],
        'traits': [
            ('coat-one-weapon-or-ammunition-from-vial', 'use a vial of Basic Poison to coat one weapon or'),
            ('applies-on-piercing-or-slashing-damage', 'A creature that takes Piercing or Slashing damage from the poisoned weapon or ammunition'),
            ('potency-ends-early-when-damage-is-dealt', 'or until its damage is dealt, whichever comes first'),
        ],
    },
    'pole': {
        'measures': [
            ('length', '10', 'feet', 'is 10 feet long'),
            ('touch-range', '10', 'feet', 'up to 10 feet away'),
        ],
        'traits': [
            ('advantage-when-vaulting-in-high-or-long-jump-athletics-check', 'If you must make a Strength (Athletics) check as part of a High or Long Jump, you can use the Pole to vault, giving yourself Advantage on the check'),
        ],
    },
    'pot-iron': {
        'measures': [
            ('capacity', '1', 'gallon', 'holds up to 1 gallon'),
        ],
    },
}


def main():
    context = srd_gear.verify(sys.argv, 260, "data-adventuring-gear-05", 12, 12, EVIDENCE, ())
    print(f"ok: 12 gear families and 12 variants match their printed entries and tables")


main()
