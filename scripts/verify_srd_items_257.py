#!/usr/bin/env python3
"""Compare issue #257 adventuring gear records with the pinned PDF.

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
    'block-and-tackle': {
        'measures': [
            ('hoist-multiplier', '4', 'times-normal-lift', 'hoist up to four times the weight you can normally lift'),
        ],
    },
    'book': {
        'measures': [
            ('knowledge-check-bonus', '5', 'bonus', '+5 bonus'),
        ],
        'traits': [
            ('contains-fiction-or-nonfiction', 'contains fiction or nonfiction'),
            ('bonus-requires-consulting-accurate-nonfiction-about-topic', 'If you consult an accurate nonfiction Book about its topic'),
            ('bonus-to-intelligence-arcana-history-nature-or-religion-checks', 'to Intelligence (Arcana, History, Nature, or Religion) checks you make about that topic'),
        ],
    },
    'bottle-glass': {
        'measures': [
            ('capacity', '3/2', 'pints', 'holds up to 11/2 pints'),
        ],
    },
    'bucket': {
        'measures': [
            ('capacity-volume', '1/2', 'cubic-foot', 'up to half a cubic foot of contents'),
        ],
    },
    'burglars-pack': {
    },
    'caltrops': {
        'checks': ['avoid-damage-and-halt'],
        'damage': ['failed-save'],
        'measures': [
            ('covered-square-side', '5', 'feet', '5-foot-square area'),
            ('placement-range', '5', 'feet', 'within 5 feet of yourself'),
            ('failed-save-speed', '0', 'speed', 'Speed reduced to 0'),
            ('recovery-time', '10', 'minutes', 'It takes 10 minutes to recover the Caltrops'),
        ],
        'traits': [
            ('spread-from-bag', 'you can spread Caltrops from their bag'),
            ('triggers-on-first-entry-on-a-turn', 'A creature that enters this area for the first time on a turn'),
            ('speed-penalty-lasts-until-start-of-next-turn', 'until the start of its next turn'),
        ],
    },
    'candle': {
        'measures': [
            ('burn-duration', '1', 'hour', 'For 1 hour'),
            ('bright-light-radius', '5', 'feet', 'Bright Light in a 5-foot radius'),
            ('additional-dim-light', '5', 'feet', 'Dim Light for an additional 5 feet'),
        ],
        'traits': [
            ('sheds-light-while-lit', 'a lit Candle sheds'),
        ],
    },
    'case-crossbow-bolt': {
        'measures': [
            ('capacity', '20', 'bolts', 'holds up to 20 Bolts'),
        ],
    },
    'case-map-or-scroll': {
        'measures': [
            ('paper-capacity', '10', 'sheets', 'up to 10 sheets of paper'),
            ('parchment-capacity', '5', 'sheets', '5 sheets of parchment'),
        ],
        'traits': [
            ('holds-paper-or-parchment', 'sheets of paper or'),
        ],
    },
    'chain': {
        'checks': ['wrap-creature', 'escape', 'burst'],
        'measures': [
            ('wrap-range', '5', 'feet', 'within 5 feet of yourself'),
        ],
        'traits': [
            ('target-unwilling-grappled-incapacitated-or-restrained-creature', 'unwilling creature within 5 feet of yourself that has the Grappled, Incapacitated, or Restrained condition'),
            ('bound-legs-restrain-until-escape', 'If the creature’s legs are bound, the creature has the Restrained condition until it escapes'),
        ],
    },
    'chest': {
        'measures': [
            ('capacity-volume', '12', 'cubic-feet', 'up to 12 cubic feet of contents'),
        ],
    },
    'climbers-kit': {
        'measures': [
            ('maximum-fall-from-anchor', '25', 'feet', 'can’t fall more than 25 feet from the anchor point'),
            ('maximum-movement-from-anchor', '25', 'feet', 'can’t move more than 25 feet from there'),
        ],
        'traits': [
            ('includes-boot-tips-gloves-pitons-and-harness', 'includes boot tips, gloves, pitons, and a harness'),
            ('anchor-yourself', 'use the Climber’s Kit to anchor yourself'),
            ('moving-farther-requires-undoing-anchor', 'without undoing the anchor'),
        ],
    },
}


def main():
    context = srd_gear.verify(sys.argv, 257, "data-adventuring-gear-02", 12, 12, EVIDENCE, ())
    print(f"ok: 12 gear families and 12 variants match their printed entries and tables")


main()
