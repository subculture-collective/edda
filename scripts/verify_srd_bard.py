#!/usr/bin/env python3
"""Offline #287 comparison with the pinned licensed SRD pp31–33; no network."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = '8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87'
ID = 'srd-5.2.1:character/'
KEY = 'srd-5.2.1:mechanic/'


def verify(pdf):
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, 'PDF pin mismatch'
    text = subprocess.check_output(['pdftotext', '-f', '31', '-l', '33', '-layout', str(pdf), '-'], text=True)
    c = json.loads((ROOT / 'internal/rules/srd521/testdata/character-bard.json').read_text())
    assert c['source']['sha256'] == PIN
    records = {r['id']: r for r in c['records']}
    assert len(records) == len(c['records']), 'duplicate canonical ID'
    inventory = json.loads((ROOT / 'docs/investigations/issue-74-character-inventory.json').read_text())
    expected = {r['id']: r['source_page'] for r in inventory['records'] if r['batch_key'] == 'class-bard'}
    canonical = [r for r in c['records'] if r['kind'] in ('class', 'feature')]
    assert len(canonical) == len(expected) == 13
    assert {s['inventory_id']: s['page'] for r in canonical for s in r['sources']} == expected
    rows = re.findall(r'^\s+(\d{1,2})\s+\+(\d)\s+(.+?)\s{2,}D(\d+)\s+(\d)\s+(\d+)\s+([\d —]+)\s*$', text, re.M)
    assert [int(r[0]) for r in rows] == list(range(1, 21)), '20 source rows required'

    def count(kind, slug, level):
        return next((v['value'] for v in records[ID + kind + '/bard/' + slug]['values_by_level'] if v['level'] == level), 0)

    for level, bonus, names, die, cantrips, prepared, slots in rows:
        n = int(level)
        row = records[f'{ID}level/bard/{n}']
        p = row['progression']
        assert row['sources'][0]['page'] == 31
        assert p['proficiency_bonus'] == int(bonus)
        assert row['parameters'] == {KEY + 'bard/bardic-die-sides': int(die)}
        actual = [records[ref]['name'] for ref in p['feature_refs']]
        actual += ['Subclass feature' for _ in p['feature_keys']]
        # Rows 1 and 2 wrap onto a second printed line.
        if n == 1:
            names += ' Spellcasting'
        if n == 2:
            names += ' Jack of All Trades'
        assert ', '.join(actual).lower().replace(',', '') == names.lower().replace(',', '').strip().replace('—', ''), (n, actual, names)
        assert count('choice', 'cantrips', n) == int(cantrips)
        assert count('choice', 'prepared-spells', n) == int(prepared)
        expected_slots = [0 if v == '—' else int(v) for v in slots.split()]
        assert len(expected_slots) == 9
        assert [count('resource', f'spell-slot-{j}', n) for j in range(1, 10)] == expected_slots
    # Extract prose by column to avoid interleaving opposite-column sentences.
    columns = [subprocess.check_output(['pdftotext', '-f', str(p), '-l', str(p), '-x', str(x), '-y', '0',
                '-W', '310', '-H', '792', '-layout', str(pdf), '-'], text=True) for p in (31, 32, 33) for x in (0, 310)]
    flat = ' '.join(re.sub(r'-\n\s*', '', ' '.join(columns).replace('\u00ad', '')).split())
    anchors = [
        'D8 per Bard level', 'Dexterity and Charisma', 'Choose any 3 skills', 'Choose 3 Musical Instruments',
        'Simple weapons', 'Light armor', '2 Daggers', '19 GP', '90 GP', 'one skill', 'one Musical Instrument',
        'As a Bonus Action', '60 feet', 'who can see or hear you', 'only one Bardic Inspiration die at a time',
        'within the next hour', 'fails a D20 Test', 'expended when it’s rolled', 'Charisma modifier (minimum of once)',
        'd8 at level 5', 'd10 at level 10', 'd12 at level 15',
        'Dancing Lights and Vicious', 'Mockery are recommended', 'replace one of your cantrips', 'levels 4 and 10',
        'four level 1', 'Charm Person, Color', 'Spray, Dissonant Whispers, and Healing Word',
        'must be of a level for which you have spell slots', 'always have prepared', 'don’t count',
        'replace one spell', 'Charisma is your spellcasting', 'Musical Instrument as a Spellcasting Focus',
        'Expertise', 'two of your skill proficiencies', 'Performance and Persuasion', 'level 9', 'two more',
        'half your Proficiency Bonus (round down)', 'skill proficiency you lack', 'doesn’t otherwise use your Proficiency Bonus',
        'College of Lore', 'Bard level or lower', 'levels 8, 12, and 16', 'for which you qualify',
        'Short or Long Rest', 'expend a spell slot (no action', 'regain one expended use',
        '30 feet', 'Charmed or Frightened', 'Reaction', 'rerolled', 'Advantage',
        'including this level', 'Prepared Spells number', 'increases', 'Bard, Cleric, Druid, and Wizard spell lists',
        'whenever you replace', 'count as Bard spells', 'roll Initiative', 'until you have two', 'fewer than that',
        'Boon of Spell Recall', 'Power Word Heal and Power Word Kill', 'always', 'second creature', 'within 10 feet of the first target',
    ]
    missing = [a for a in anchors if a not in flat]
    assert not missing, f'missing source anchors: {missing}'
    print(f'SRD {PIN}: 13 identities/pages, 20 Bard rows, all nine slot columns and prose anchors match pp31–33')


if __name__ == '__main__':
    if len(sys.argv) != 2:
        raise SystemExit('usage: verify_srd_bard.py /path/to/SRD_CC_v5.2.1.pdf')
    verify(Path(sys.argv[1]))
