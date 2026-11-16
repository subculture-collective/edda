#!/usr/bin/env python3
"""Read-only Warlock candidate proof; strict mode requires actual publication."""
import argparse
from collections import Counter
import hashlib
import importlib
import json
from pathlib import Path
import re
import subprocess

from warlock_contract_oracle import (CONTRACTS, SPLITS, TOME_CANTRIPS, TOME_RITUALS,
                                     FIEND_SPELLS, NAMED_SPELLS, body, level_fixture)

ROOT = Path(__file__).resolve().parents[1]
MANIFEST = ROOT / 'docs/investigations/issue-331-warlock-scope.json'
PIN = '8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87'
BASE = 'e7cd1ba8318ba296edd4fc5b6b84715ca11016f2'
NAMES = ['warlock', 'fiend', 'invocations_01', 'invocations_02', 'invocations_03']


def exact(a, b):
    if type(a) is not type(b):
        return False
    if isinstance(a, dict):
        return a.keys() == b.keys() and all(exact(a[k], b[k]) for k in a)
    if isinstance(a, list):
        return len(a) == len(b) and all(exact(x, y) for x, y in zip(a, b))
    return a == b


def source_records(pdf):
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, 'PDF pin'
    result = {}
    for name in NAMES:
        mod = importlib.import_module('verify_srd_' + name)
        records = mod.source_records(pdf) if name in ('warlock', 'invocations_02') else mod.source_records()
        result[name] = records
    assert [len(x) for x in result.values()] == [71, 17, 41, 49, 13]
    return result


def spell_rows(pdf, pages):
    result, level = [], None
    for page, x in pages:
        text = subprocess.check_output(['pdftotext', '-f', str(page), '-l', str(page),
            '-x', str(x), '-y', '0', '-W', '310', '-H', '792', '-layout', str(pdf), '-'], text=True)
        for line in text.splitlines():
            h = re.search(r'(?:Level (\d) \w+ Spells|Cantrips \(Level (0))', line)
            if h:
                level = int(h[1] or h[2])
            m = re.match(r'\s*(.+?)\s{2,}(Abjuration|Conjuration|Divination|Enchantment|Evocation|Illusion|Necromancy|Transmutation)\s+([CMR, —]+)\s*$', line)
            if m:
                result.append(dict(name=m[1].strip(), level=level, page=page, school=m[2], special=m[3].strip()))
    return result


def source_owner(record, key, by):
    if record['kind'] == 'level':
        return 'level-' + str(record['level'])
    if key.endswith('/hit-die-sides'):
        return 'hp'
    if '/multiclass/' in key:
        return 'later-entry'
    ident = record['id']
    if '/arcanum-' in ident:
        return 'mystic-arcanum'
    if '/pact-slots' in ident:
        return 'slots'
    if any('/' + s in ident for s in ('cantrips', 'prepared-spells', 'cantrip-eligibility', 'prepared-spell-eligibility')):
        return 'preparation'
    if record['kind'] == 'subclass':
        return 'warlock-subclass'
    if record['kind'] == 'feature':
        return ident.rsplit('/', 1)[-1]
    if record['parent_id'] in by:
        return source_owner(by[record['parent_id']], key, by)
    return 'creation'


def acyclic(graph):
    done, active = set(), set()
    def visit(n):
        assert n not in active, 'cycle'
        if n in done:
            return
        active.add(n)
        for d in graph.get(n, []):
            visit(d)
        active.remove(n)
        done.add(n)
    for n in graph:
        visit(n)


def caps(graph, parents):
    acyclic(graph)
    acyclic({k: [v] for k, v in parents.items()})
    assert all(len(v) == len(set(v)) and len(v) <= 100 for v in graph.values()), 'incoming cap'
    assert max(Counter(d for v in graph.values() for d in v).values(), default=0) <= 100, 'outgoing cap'
    assert max(Counter(parents.values()).values(), default=0) <= 100, 'children cap'


def bounded_pages(fetch, endpoint, limit=16):
    assert endpoint.endswith('per_page=100') and 'page=' not in endpoint.replace('per_page=', ''), 'exact query'
    out, seen, payloads = [], set(), set()
    for page in range(1, limit + 1):
        data = fetch(endpoint + '&page=' + str(page))
        assert isinstance(data, list) and len(data) <= 100, 'page shape'
        if not data:
            return out
        payload = json.dumps(data, sort_keys=True)
        assert payload not in payloads, 'repeated page'
        payloads.add(payload)
        for row in data:
            assert type(row['number']) is int and row['number'] not in seen, 'duplicate/invalid identity'
            seen.add(row['number'])
            out.append(row)
        if len(data) < 100:
            return out
    raise AssertionError('call cap')


def prefix_check(actual, prefix, append, marker):
    assert actual.startswith(prefix + append), 'immutable prefix and exact append'
    assert actual.count(marker) == 1, 'append occurrence'


def expected_tasks(records):
    by = {r['id']: r for r in records}
    tasks = []
    def add(key, fixture, depends=(), role='implementation', parent='features-cohort'):
        tasks.append(dict(key=key, title='Implement Warlock ' + key.replace('-', ' ') if role == 'implementation' else 'Qualify Warlock ' + key.replace('-', ' '),
            fixture=fixture, body=body(key, fixture), type='Task', size='M', parent=parent,
            issue=None, prerequisites=[968] if key == 'hp' else [430] if key == 'creation' else [],
            depends=sorted(set(depends)), role=role))
    for key, fixture in CONTRACTS.items():
        add(key, fixture, SPLITS.get(key, {}), 'qualification' if key in SPLITS else 'implementation')
    for parent, leaves in SPLITS.items():
        for key, fixture in leaves.items():
            add(key, fixture)
    for level in range(1, 21):
        row = by['srd-5.2.1:character/level/warlock/' + str(level)]
        deps = ['hp', 'preparation', 'eldritch-invocations', 'slots'] + [r.rsplit('/', 1)[-1] for r in row['progression']['feature_refs']]
        if level in (3, 6, 10, 14):
            deps += ['warlock-subclass']
        if level == 1:
            deps += ['creation', 'later-entry']
        if level >= 3:
            deps += ['dark-ones-blessing', 'fiend-spells']
        if level >= 6:
            deps += ['dark-ones-own-luck']
        if level >= 10:
            deps += ['fiendish-resilience']
        if level >= 14:
            deps += ['hurl-through-hell']
        add('level-' + str(level), level_fixture(level, row), deps, parent='levels-cohort')
    return tasks


def validate_candidate(m, expected):
    assert m['decision'] == 331 and m['base'] == BASE
    records = [r for rs in expected.values() for r in rs]
    assert exact(m['records'], records), 'actual candidate source records'
    for name, rs in expected.items():
        path = ROOT / ('internal/rules/srd521/testdata/character-' + name.replace('_', '-') + '.json')
        assert exact(json.loads(path.read_text())['records'], rs), 'actual catalogue source'
    by = {r['id']: r for r in records}
    assert len(by) == 191 and len(m['record_matrix']) == 191
    matrix = []
    for r in records:
        keys = list(dict.fromkeys(r['required_mechanics'] + r['option_keys']))
        matrix.append(dict(record=r, key_owners={k: source_owner(r, k, by) for k in keys},
                           status='unavailable', evidence=m['consumer_evidence']))
    assert exact(m['record_matrix'], matrix), 'complete contextual owner matrix'
    tasks = {t['key']: t for t in m['tasks']}
    expected_rows = {t['key']: t for t in expected_tasks(records)}
    expected_keys = expected_rows.keys()
    assert tasks.keys() == expected_keys, 'complete finite class task roster'
    for key, t in tasks.items():
        assert exact(t, expected_rows[key]), ('independent complete Task contract', key)
        fixture = expected_rows[key]['fixture']
        if key.startswith('level-'):
            level = int(key[6:])
            row = by['srd-5.2.1:character/level/warlock/' + str(level)]
            fixture = level_fixture(level, row)
            assert 'hp' in t['depends'] and 'preparation' in t['depends'] and 'eldritch-invocations' in t['depends'], 'every granting path'
            features = {x.rsplit('/', 1)[-1] for x in row['progression']['feature_refs']}
            assert features <= set(t['depends']), 'all repeated feature grants'
            if level in (3, 6, 10, 14):
                assert 'warlock-subclass' in t['depends'], 'subclass path'
        assert t['fixture'] == fixture and t['body'] == body(key, fixture), ('complete candidate body', key)
        assert t['type'] == 'Task' and t['size'] == 'M' and t['parent'] in ('features-cohort', 'levels-cohort')
        assert t['prerequisites'] == ([968] if key == 'hp' else [430] if key == 'creation' else []), 'owner prerequisites'
    caps({k: v['depends'] for k, v in tasks.items()}, {k: v['parent'] for k, v in tasks.items()})
    assert len(m['spells']['base']) == 72 and len(m['spells']['consumed']) == 117
    assert m['publication']['state'] == 'PENDING', 'checkpoint must remain honest'


def validate_publication(m, receipt, live=False):
    assert live, 'strict proof requires live native readback, never forged published flags'
    assert receipt and receipt['state'] == 'PUBLISHED', 'PENDING native publication'
    assert receipt.get('root_accepted_post_sorcerer_release'), 'ROOT final prefix release'
    assert receipt.get('root_serial_appends') == [47, 54, 283, 337], 'ROOT appends'
    # Until complete split/semantic ownership approval exists no receipt can certify this draft.
    assert m['publication'].get('complete_semantic_approval') is True, 'complete finite split approval'
    raise AssertionError('Native publication validator not sealed: exact native Task/Project/graph receipt required')


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--pdf', type=Path, required=True)
    p.add_argument('--source-class-checkpoint', action='store_true')
    args = p.parse_args()
    m = json.loads(MANIFEST.read_text())
    expected = source_records(args.pdf)
    validate_candidate(m, expected)
    base = spell_rows(args.pdf, [(75, 0), (75, 310), (76, 0)])
    assert exact(m['spells']['base'], base), 'direct base list'
    allrows = spell_rows(args.pdf, [(p, x) for p in range(30, 84) for x in (0, 310)])
    assert sorted({r['name'] for r in allrows if r['level'] == 0}) == sorted(TOME_CANTRIPS)
    assert sorted({r['name'] for r in allrows if r['level'] == 1 and 'R' in r['special']}) == sorted(TOME_RITUALS)
    consumed = sorted({r['name'] for r in base} | set(TOME_CANTRIPS + TOME_RITUALS + FIEND_SPELLS + NAMED_SPELLS))
    assert m['spells']['consumed'] == consumed, 'full consumed union'
    if not args.source_class_checkpoint:
        validate_publication(m, m.get('native_receipt'))
    print('PASS source/class checkpoint only: 191 records,28 Invocation identities,20 rows,72+Tome/Fiend/named=117 spell identities; native proof PENDING')


if __name__ == '__main__':
    main()
