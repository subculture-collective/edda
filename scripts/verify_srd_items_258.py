#!/usr/bin/env python3
"""Compare issue #258 adventuring gear records with the pinned PDF.

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
    'clothes-fine': {
        'traits': [
            ('expensive-fabrics-with-expertly-crafted-details', 'made of expensive fabrics and adorned with expertly crafted details'),
            ('required-for-admission-to-some-events-and-locations', 'Some events and locations admit only people wearing these clothes'),
        ],
    },
    'clothes-travelers': {
        'traits': [
            ('resilient-garments-for-travel-in-various-environments', 'resilient garments designed for travel in various environments'),
        ],
    },
    'component-pouch': {
        'traits': [
            ('watertight', 'is watertight'),
            ('compartments-hold-all-free-material-components-of-your-spells', 'compartments that hold all the free Material components of your spells'),
        ],
    },
    'costume': {
        'traits': [
            ('advantage-on-checks-to-impersonate-represented-person-or-type-while-wearing', 'While wearing a Costume, you have Advantage on any ability check you make to impersonate the person or type of person it represents'),
        ],
    },
    'crowbar': {
        'traits': [
            ('advantage-on-strength-checks-where-leverage-applies', 'gives you Advantage on Strength checks where the Crowbar’s leverage can be applied'),
        ],
    },
    'diplomats-pack': {
    },
    'druidic-focus': {
        'traits': [
            ('carved-ribboned-or-painted-to-channel-primal-magic', 'is carved, tied with ribbon, or painted to channel primal magic'),
        ],
    },
    'dungeoneers-pack': {
    },
    'entertainers-pack': {
    },
    'explorers-pack': {
    },
    'flask': {
        'measures': [
            ('capacity', '1', 'pint', 'holds up to 1 pint'),
        ],
    },
    'grappling-hook': {
        'checks': ['catch-hook'],
        'measures': [
            ('throw-range', '50', 'feet', 'within 50 feet of yourself'),
        ],
        'traits': [
            ('throw-at-railing-ledge-or-other-catch', 'at a railing, a ledge, or another catch'),
            ('tied-rope-can-then-be-climbed', 'If you tied a Rope to the hook, you can then climb it'),
        ],
    },
}


def main():
    context = srd_gear.verify(sys.argv, 258, "data-adventuring-gear-03", 12, 14, EVIDENCE, ())
    print(f"ok: 12 gear families and 14 variants match their printed entries and tables")


main()
