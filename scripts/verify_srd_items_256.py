#!/usr/bin/env python3
"""Compare issue #256 adventuring gear records with the pinned PDF.

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
    'acid': {
        'checks': ['avoid-damage'],
        'damage': ['thrown-vial'],
        'measures': [
            ('throw-range', '20', 'feet', 'within 20 feet of yourself'),
        ],
        'traits': [
            ('throw-vial-in-place-of-one-attack', 'replace one of your attacks with throwing a vial of Acid'),
            ('target-one-creature-or-object-you-can-see', 'Target one creature or object you can see'),
        ],
    },
    'alchemists-fire': {
        'checks': ['avoid-damage-and-burning'],
        'damage': ['thrown-flask'],
        'measures': [
            ('throw-range', '20', 'feet', 'within 20 feet of yourself'),
        ],
        'traits': [
            ('throw-flask-in-place-of-one-attack', 'replace one of your attacks with throwing a flask of Alchemist’s Fire'),
            ('target-one-creature-or-object-you-can-see', 'Target one creature or object you can see'),
            ('failed-save-starts-burning', 'start burning'),
        ],
    },
    'ammunition': {
        'traits': [
            ('required-by-weapon-with-ammunition-property', 'Ammunition is required by a weapon that has the Ammunition property'),
            ('weapon-description-specifies-ammunition-type', 'A weapon’s description specifies the type of ammunition used by the weapon'),
            ('bought-in-listed-amounts', 'the amount you get when you buy them'),
            ('storage-bought-separately', 'storage must be bought separately'),
        ],
    },
    'antitoxin': {
        'measures': [
            ('benefit-duration', '1', 'hour', 'for 1 hour'),
        ],
        'traits': [
            ('drink-one-vial', 'you can drink a vial of Antitoxin'),
            ('advantage-on-saves-to-avoid-or-end-poisoned', 'Advantage on saving throws to avoid or end the Poisoned condition'),
        ],
    },
    'arcane-focus': {
        'traits': [
            ('bejeweled-or-carved-to-channel-arcane-magic', 'is bejeweled or carved to channel arcane magic'),
        ],
    },
    'backpack': {
        'measures': [
            ('capacity-weight', '30', 'pounds', 'holds up to 30 pounds'),
            ('capacity-volume', '1', 'cubic-foot', 'within 1 cubic foot'),
        ],
        'traits': [
            ('can-serve-as-saddlebag', 'It can also serve as a saddlebag'),
        ],
    },
    'ball-bearings': {
        'checks': ['avoid-prone'],
        'measures': [
            ('covered-square-side', '10', 'feet', '10-foot-square area'),
            ('placement-range', '10', 'feet', 'within 10 feet of yourself'),
            ('recovery-time', '10', 'minutes', 'It takes 10 minutes to recover the Ball Bearings'),
        ],
        'traits': [
            ('spill-from-pouch-over-level-area', 'you can spill Ball Bearings from their pouch. They spread to cover a level'),
            ('triggers-on-first-entry-on-a-turn', 'A creature that enters this area for the first time on a turn'),
            ('failed-save-causes-prone', 'or have the Prone condition'),
        ],
    },
    'barrel': {
        'measures': [
            ('liquid-capacity', '40', 'gallons', 'up to 40 gallons of liquid'),
            ('dry-goods-capacity', '4', 'cubic-feet', 'up to 4 cubic feet of dry goods'),
        ],
        'traits': [
            ('holds-liquid-or-dry-goods', 'of liquid or up to'),
        ],
    },
    'basket': {
        'measures': [
            ('capacity-weight', '40', 'pounds', 'holds up to 40 pounds'),
            ('capacity-volume', '2', 'cubic-feet', 'within 2 cubic feet'),
        ],
    },
    'bedroll': {
        'measures': [
            ('occupants', '1', 'creature', 'sleeps one Small or Medium creature'),
        ],
        'traits': [
            ('occupant-is-small-or-medium', 'Small or Medium creature'),
            ('automatic-success-on-saves-against-extreme-cold-while-inside', 'While in a Bedroll, you automatically succeed on saving throws against extreme cold'),
        ],
    },
    'bell': {
        'measures': [
            ('audible-distance', '60', 'feet', 'heard up to 60 feet away'),
        ],
        'traits': [
            ('rung-to-produce-sound', 'When rung as a Utilize action, a Bell produces a sound'),
        ],
    },
    'blanket': {
        'traits': [
            ('advantage-on-saves-against-extreme-cold-while-wrapped', 'While wrapped in a blanket, you have Advantage on saving throws against extreme cold'),
        ],
    },
}


def main():
    context = srd_gear.verify(sys.argv, 256, "data-adventuring-gear-01", 12, 20, EVIDENCE, ())
    print(f"ok: 12 gear families and 20 variants match their printed entries and tables")


main()
