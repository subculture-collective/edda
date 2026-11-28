#!/usr/bin/env python3
"""Compare issue #261 adventuring gear records with the pinned PDF. It also checks the p. 99 Potion of Healing entry against the Healing contract fixture, which owns it.

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
    'pouch': {
        'measures': [
            ('capacity-weight', '6', 'pounds', 'holds up to 6 pounds'),
            ('capacity-volume', '1/5', 'cubic-foot', 'within one-fifth of a cubic foot'),
        ],
    },
    'priests-pack': {
    },
    'quiver': {
        'measures': [
            ('capacity', '20', 'arrows', 'holds up to 20 Arrows'),
        ],
    },
    'ram-portable': {
        'measures': [
            ('strength-check-bonus', '4', 'bonus', '+4 bonus to the Strength check'),
            ('helpers', '1', 'character', 'One other character can help you'),
        ],
        'traits': [
            ('used-to-break-down-doors', 'to break down doors'),
            ('helper-grants-advantage-on-check', 'giving you Advantage on this check'),
        ],
    },
    'rations': {
        'traits': [
            ('travel-ready-food-including-jerky-dried-fruit-hardtack-and-nuts', 'travel-ready food, including jerky, dried fruit, hardtack, and nuts'),
            ('not-eating-risks-malnutrition', 'See “Malnutrition” in “Rules Glossary” for the risks of not eating'),
        ],
    },
    'robe': {
        'traits': [
            ('vocational-or-ceremonial-significance', 'has vocational or ceremonial significance'),
            ('required-colors-or-symbols-for-admission-to-some-events-and-locations', 'Some events and locations admit only people wearing a Robe bearing certain colors or symbols'),
        ],
    },
    'rope': {
        'checks': ['tie-knot', 'burst', 'escape'],
        'traits': [
            ('bind-only-unwilling-grappled-incapacitated-or-restrained-creature', 'You can bind an unwilling creature with the Rope only if the creature has the Grappled, Incapacitated, or Restrained condition'),
            ('bound-legs-restrain-until-escape', 'If the creature’s legs are bound, the creature has the Restrained condition until it escapes'),
        ],
    },
    'sack': {
        'measures': [
            ('capacity-weight', '30', 'pounds', 'holds up to 30 pounds'),
            ('capacity-volume', '1', 'cubic-foot', 'within 1 cubic foot'),
        ],
    },
    'scholars-pack': {
    },
    'shovel': {
        'measures': [
            ('digging-time', '1', 'hour', 'Working for 1 hour'),
            ('hole-side', '5', 'feet', '5 feet on each side'),
        ],
        'traits': [
            ('digs-hole-in-soil-or-similar-material', 'dig a hole that is 5 feet on each side in soil or similar material'),
        ],
    },
    'signal-whistle': {
        'measures': [
            ('audible-distance', '600', 'feet', 'heard up to 600 feet away'),
        ],
        'traits': [
            ('blown-to-produce-sound', 'When blown as a Utilize action, a Signal Whistle produces a sound'),
        ],
    },
}


def main():
    context = srd_gear.verify(sys.argv, 261, "data-adventuring-gear-06", 11, 11, EVIDENCE, ("srd-5.2.1:content/adventuring-gear/potion-of-healing",))
    srd_gear.check_potion_of_healing(context)
    print(f"ok: 11 gear families and 11 variants match their printed entries and tables")


main()
