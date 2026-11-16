"""Checkpoint proofs stay honest until native publication is complete."""
from copy import deepcopy
import json
import unittest
import verify_paladin_scope as proof


class PaladinCheckpointTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.manifest=json.loads(proof.MANIFEST.read_text())

    def test_complete_checkpoint_and_corruption_controls(self):
        self.assertEqual(proof.validate(self.manifest),87)
        proof.probes(self.manifest)

    def test_every_class_application_consumes_own_hp_adapter(self):
        rows={x['key']:x for x in self.manifest['issues']}
        self.assertEqual([x['key'] for x in rows.values() if 429 in x['prerequisites']],['hp'])
        for key in [f'level-{n}' for n in range(1,21)]+['later-1']:
            self.assertIn('hp',rows[key]['depends'])

    def test_spell_effects_are_not_membership(self):
        rows={x['key']:x for x in self.manifest['shared_issues']}
        self.assertIn(352,rows['bless-extension']['prerequisites'])
        self.assertIn(265,rows['detect-extension']['prerequisites'])
        self.assertIn(625,rows['curse-item']['prerequisites'])
        self.assertFalse(any(x['spell']=='Dispel Magic' for x in rows.values()))

    def test_unpublished_issue_id_cannot_be_advertised_as_receipt(self):
        m=deepcopy(self.manifest);m['issues'][0]['issue']=99999
        with self.assertRaises(AssertionError):proof.validate(m)


if __name__=='__main__':unittest.main()
