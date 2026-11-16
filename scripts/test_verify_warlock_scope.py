"""Contrary probes against the actual Warlock candidate, not a synthetic receipt."""
import copy
import json
from pathlib import Path
import unittest

import verify_warlock_scope as v


class WarlockCandidateProof(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.candidate = json.loads(v.MANIFEST.read_text())
        cls.expected = v.source_records(Path('/tmp/edda-wave1/SRD_CC_v5.2.1.pdf'))

    def reject(self, mutation):
        candidate = copy.deepcopy(self.candidate)
        mutation(candidate)
        with self.assertRaises(AssertionError):
            v.validate_candidate(candidate, self.expected)

    def test_actual_candidate(self):
        v.validate_candidate(self.candidate, self.expected)

    def test_every_record_and_context_required(self):
        for n in range(191):
            self.reject(lambda m, n=n: m['records'].pop(n))
            self.reject(lambda m, n=n: m['record_matrix'].pop(n))

    def test_every_context_owner_independently_protected(self):
        for n, row in enumerate(self.candidate['record_matrix']):
            for key in row['key_owners']:
                self.reject(lambda m, n=n, key=key: m['record_matrix'][n]['key_owners'].__setitem__(key, 'contrary-owner'))

    def test_numeric_json_is_typed(self):
        for n, record in enumerate(self.candidate['records']):
            for key, value in record.get('parameters', {}).items():
                if type(value) is int:
                    for bad in [True, float(value)]:
                        self.reject(lambda m, n=n, key=key, bad=bad: m['records'][n]['parameters'].__setitem__(key, bad))

    def test_all_independent_grants(self):
        for n, task in enumerate(self.candidate['tasks']):
            if task['key'].startswith('level-'):
                for dep in task['depends']:
                    if dep in ('hp', 'preparation', 'eldritch-invocations', 'ability-score-improvement', 'mystic-arcanum', 'warlock-subclass'):
                        self.reject(lambda m, n=n, dep=dep: m['tasks'][n]['depends'].remove(dep))

    def test_coherent_body_fixture_and_fake_receipt_weakening(self):
        for n in range(len(self.candidate['tasks'])):
            def weaken(m, n=n):
                m['tasks'][n]['fixture'] = 'Accept any known spell with no target or resource validation.'
                m['tasks'][n]['body'] = v.body(m['tasks'][n]['key'], m['tasks'][n]['fixture'])
                m['native_receipt'] = {'body': m['tasks'][n]['body'], 'state': 'PUBLISHED'}
            self.reject(weaken)

    def test_actual_native_flags_cannot_forge_publication(self):
        for live in [False, True]:
            receipt = {'state': 'PUBLISHED', 'root_accepted_post_sorcerer_release': True,
                       'root_serial_appends': [47, 54, 283, 337]}
            with self.assertRaises(AssertionError):
                v.validate_publication(self.candidate, receipt, live=live)

    def test_hp_adapter_and_metadata(self):
        n = next(n for n, t in enumerate(self.candidate['tasks']) if t['key'] == 'hp')
        self.reject(lambda m: m['tasks'][n].__setitem__('prerequisites', [429]))
        for field, bad in [('type', 'Feature'), ('size', 'L'), ('parent', 283)]:
            self.reject(lambda m, field=field, bad=bad: m['tasks'][n].__setitem__(field, bad))

    def test_dependency_and_containment_separate_caps(self):
        v.caps({'a': ['b']}, {'b': 'a'})
        for graph, parents in [({'a': ['b'], 'b': ['a']}, {}), ({}, {'a': 'b', 'b': 'a'}),
                               ({'a': list(range(101))}, {}),
                               ({n: ['root'] for n in range(101)}, {}),
                               ({}, {n: 'parent' for n in range(101)})]:
            with self.assertRaises(AssertionError):
                v.caps(graph, parents)

    def test_bounded_actual_queries(self):
        endpoint = '/issues?state=all&per_page=100'
        calls = []
        def fetch(url):
            calls.append(url)
            return [{'number': n} for n in range(100)]
        with self.assertRaises(AssertionError):
            v.bounded_pages(fetch, endpoint)
        self.assertEqual(calls, [endpoint + '&page=1', endpoint + '&page=2'])
        calls.clear()
        def endless(url):
            calls.append(url)
            return [{'number': len(calls) * 100 + n} for n in range(100)]
        with self.assertRaises(AssertionError):
            v.bounded_pages(endless, endpoint, limit=3)
        self.assertEqual(len(calls), 3)

    def test_exact_prefix_and_permitted_later_append(self):
        prefix, append = 'old owner body\n', '<!-- edda331 --> exact root supplement\n'
        v.prefix_check(prefix + append + 'later legitimate append\n', prefix, append, '<!-- edda331 -->')
        for actual in [append, prefix, prefix + append.replace('exact', 'weak'), prefix + append * 2]:
            with self.assertRaises(AssertionError):
                v.prefix_check(actual, prefix, append, '<!-- edda331 -->')


if __name__ == '__main__':
    unittest.main()
