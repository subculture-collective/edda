#!/usr/bin/env python3
"""Compare issue #262 adventuring gear records with the pinned PDF. It also checks the p. 99 Spell Scroll entry against the #242 records, which own it.

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
    'spikes-iron': {
        'measures': [
            ('bundle-amount', '10', 'spikes', 'bundles of ten'),
        ],
        'traits': [
            ('hammer-with-blunt-object-into-wood-earth-or-similar-material', 'use a blunt object, such as a Light Hammer, to hammer a spike into wood, earth, or a similar material'),
            ('jam-door-shut-or-anchor-rope-or-chain', 'jam a door shut or to then tie a Rope or Chain to the Spike'),
        ],
    },
    'spyglass': {
        'measures': [
            ('magnification', '2', 'times-size', 'magnified to twice their size'),
        ],
    },
    'string': {
        'measures': [
            ('length', '10', 'feet', 'is 10 feet long'),
        ],
        'traits': [
            ('can-be-knotted', 'You can tie a knot in it'),
        ],
    },
    'tent': {
        'measures': [
            ('occupants', '2', 'creatures', 'sleeps up to two'),
        ],
        'traits': [
            ('occupants-are-small-or-medium', 'Small or Medium creatures'),
        ],
    },
    'tinderbox': {
        'measures': [
            ('other-fire-lighting-time', '1', 'minute', 'Lighting any other fire takes 1 minute'),
        ],
        'traits': [
            ('small-container-of-flint-fire-steel-and-tinder', 'a small container holding flint, fire steel, and tinder (usually dry cloth soaked in light oil) used to kindle a fire'),
            ('lights-candle-lamp-lantern-torch-or-exposed-fuel', 'Using it to light a Candle, Lamp, Lantern, or Torch—or anything else with exposed fuel—takes a Bonus Action'),
        ],
    },
    'torch': {
        'damage': ['hit'],
        'measures': [
            ('burn-duration', '1', 'hour', 'burns for 1 hour'),
            ('bright-light-radius', '20', 'feet', 'Bright Light in a 20-foot radius'),
            ('additional-dim-light', '20', 'feet', 'Dim Light for an additional 20 feet'),
        ],
        'traits': [
            ('attack-with-torch-as-simple-melee-weapon', 'attack with the Torch, using it as a Simple Melee weapon'),
        ],
    },
    'vial': {
        'measures': [
            ('capacity', '4', 'ounces', 'holds up to 4 ounces'),
        ],
    },
    'waterskin': {
        'measures': [
            ('capacity', '4', 'pints', 'holds up to 4 pints'),
        ],
        'traits': [
            ('insufficient-water-risks-dehydration', 'If you don’t drink sufficient water, you risk dehydration'),
        ],
    },
}


def main():
    context = srd_gear.verify(sys.argv, 262, "data-adventuring-gear-07", 8, 8, EVIDENCE, ("srd-5.2.1:content/adventuring-gear/spell-scroll",))
    srd_gear.check_spell_scroll(context)
    print(f"ok: 8 gear families and 8 variants match their printed entries and tables")


main()
