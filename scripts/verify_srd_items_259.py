#!/usr/bin/env python3
"""Compare issue #259 adventuring gear records with the pinned PDF.

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
    'healers-kit': {
        'measures': [
            ('uses', '10', 'uses', 'has ten uses'),
            ('stabilize-target-hit-points', '0', 'hit-points', 'has 0 Hit Points'),
        ],
        'traits': [
            ('expend-one-use-to-stabilize-unconscious-creature', 'expend one of its uses to stabilize an Unconscious creature'),
            ('no-wisdom-medicine-check-needed', 'without needing to make a Wisdom (Medicine) check'),
        ],
    },
    'holy-symbol': {
        'traits': [
            ('bejeweled-or-painted-to-channel-divine-magic', 'is bejeweled or painted to channel divine magic'),
            ('table-states-how-each-symbol-is-carried', 'The table indicates whether a Holy Symbol needs to be held, worn, or borne on fabric (such as a tabard or banner) or a Shield'),
        ],
    },
    'holy-water': {
        'checks': ['avoid-damage'],
        'damage': ['thrown-flask-against-fiend-or-undead'],
        'measures': [
            ('throw-range', '20', 'feet', 'within 20 feet of yourself'),
        ],
        'traits': [
            ('throw-flask-in-place-of-one-attack', 'replace one of your attacks with throwing a flask of Holy Water'),
            ('target-one-creature-you-can-see', 'Target one creature you can see'),
            ('damages-only-fiend-or-undead', 'if it is a Fiend or an Undead'),
        ],
    },
    'hunting-trap': {
        'checks': ['avoid-damage-and-halt', 'free-creature'],
        'damage': ['failed-save', 'each-failed-escape-check'],
        'measures': [
            ('failed-save-speed', '0', 'speed', 'Speed reduced to 0'),
            ('typical-chain-length', '3', 'feet', 'typically 3 feet'),
        ],
        'traits': [
            ('sawtooth-steel-ring-snaps-shut-on-pressure-plate', 'a sawtooth steel ring that snaps shut when a creature steps on a pressure plate in the center'),
            ('chain-affixed-to-immobile-object', 'affixed by a heavy chain to an immobile object'),
            ('speed-penalty-lasts-until-start-of-next-turn', 'until the start of its next turn'),
            ('movement-limited-by-chain-until-free', 'until the creature breaks free of the trap, its movement is limited by the length of the chain'),
            ('check-frees-self-or-creature-within-reach', 'freeing itself or another creature within its reach on a success'),
        ],
    },
    'ink': {
        'measures': [
            ('bottle-volume', '1', 'ounce', '1-ounce bottle'),
            ('writable-pages', '500', 'pages', 'write about 500 pages'),
        ],
    },
    'ink-pen': {
        'traits': [
            ('writes-or-draws-with-ink', 'Using Ink, an Ink Pen is used to write or draw'),
        ],
    },
    'jug': {
        'measures': [
            ('capacity', '1', 'gallon', 'holds up to 1 gallon'),
        ],
    },
    'ladder': {
        'measures': [
            ('height', '10', 'feet', 'is 10 feet tall'),
        ],
        'traits': [
            ('must-climb-to-move-up-or-down', 'You must climb to move up or down it'),
        ],
    },
    'lamp': {
        'measures': [
            ('bright-light-radius', '15', 'feet', 'Bright Light in a 15-foot radius'),
            ('additional-dim-light', '30', 'feet', 'Dim Light for an additional 30 feet'),
        ],
        'traits': [
            ('burns-oil-as-fuel', 'burns Oil as fuel'),
        ],
    },
    'lantern-bullseye': {
        'measures': [
            ('bright-light-cone', '60', 'feet', 'Bright Light in a 60-foot Cone'),
            ('additional-dim-light', '60', 'feet', 'Dim Light for an additional 60 feet'),
        ],
        'traits': [
            ('burns-oil-as-fuel', 'burns Oil as fuel'),
        ],
    },
    'lantern-hooded': {
        'measures': [
            ('bright-light-radius', '30', 'feet', 'Bright Light in a 30-foot radius'),
            ('additional-dim-light', '30', 'feet', 'Dim Light for an additional 30 feet'),
            ('lowered-hood-dim-light-radius', '5', 'feet', 'Dim Light in a 5-foot radius'),
        ],
        'traits': [
            ('burns-oil-as-fuel', 'burns Oil as fuel'),
            ('hood-can-be-lowered-or-raised', 'lower the hood, reducing the light to Dim Light in a 5-foot radius, or raise it again'),
        ],
    },
    'lock': {
        'checks': ['pick-lock-without-key'],
        'traits': [
            ('comes-with-key', 'comes with a key'),
            ('thieves-tools-can-pick-without-key', 'Without the key, a creature can use Thieves’ Tools to pick this Lock'),
        ],
    },
}


def main():
    context = srd_gear.verify(sys.argv, 259, "data-adventuring-gear-04", 12, 14, EVIDENCE, ())
    print(f"ok: 12 gear families and 14 variants match their printed entries and tables")


main()
