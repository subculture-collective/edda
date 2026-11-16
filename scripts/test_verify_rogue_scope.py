"""Bounded controls for checkpoint claims; publication is intentionally unavailable."""
from copy import deepcopy
import json
from pathlib import Path
import unittest
from urllib.parse import parse_qs,urlparse
import verify_rogue_scope as v

class CheckpointTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.m=json.loads(v.PATH.read_text());cls.r=json.loads(v.RECEIPT.read_text())
    def test_positive(self):v.validate(self.m);v.validate_observation(self.m,self.r)
    def test_matrix_corruptions(self):
        for mutate in [lambda m:m['records'][0]['key_owners'].clear(),lambda m:m['issues'][0]['prerequisites'].clear(),lambda m:m['grants']['7'].remove('evasion'),lambda m:m['catalogues']['advancement']['records'][5]['advancement'].update(experience_points=300),lambda m:m.update(accepted_post_ranger='invented'),lambda m:m['issues'][0].update(issue=9999),lambda m:m['records'][0]['record']['parameters'].clear()]:
            m=deepcopy(self.m);mutate(m)
            with self.assertRaises(AssertionError):v.validate(m)
    def test_coherent_weakening(self):
        m=deepcopy(self.m);r=deepcopy(self.r)
        m['issues'][0]['prerequisites']=[]
        r['consumed_contracts']['968']['blocked_by']=[]
        with self.assertRaises(AssertionError):v.validate(m)
        with self.assertRaises(AssertionError):v.validate_observation(self.m,r)
    def test_source_numeric_and_grant_controls(self):
        for change in ['source','numeric','grant','total','feat','later']:
            m=deepcopy(self.m)
            if change=='source':m['source_sha256']='0'*64
            elif change=='numeric':m['catalogues']['advancement']['records'][5]['advancement']['experience_points']=300
            elif change=='grant':m['issues'][next(i for i,x in enumerate(m['issues']) if x['key']=='level-7')]['depends'].remove('evasion')
            elif change=='total':m['issues'][0]['fixture']=m['issues'][0]['fixture'].replace('D8','D10')
            elif change=='feat':m['records'][next(i for i,x in enumerate(m['records']) if x['record']['name']=='Evasion')]['execution_owners']=['improvement']
            else:m['issues'][-1]['depends'].remove('multiclass')
            with self.assertRaises(AssertionError):v.validate(m)
    def test_marker_counts_and_append(self):
        prefix='accepted prefix';section=Path('docs/investigations/issue-329-parent-47-PENDING.md').read_text()
        v.parent_append(prefix,section,prefix+section+'\nlegitimate later parent append')
        for actual in [prefix,prefix+section+section,prefix+section+'<!-- edda-329:start --> altered duplicate <!-- edda-329:end -->','changed'+section]:
            with self.assertRaises(AssertionError):v.parent_append(prefix,section,actual)
    def test_pagination_100_101(self):
        for count,expected_calls in [(100,2),(101,2),(99,1)]:
            calls=[]
            def fake(url):
                calls.append(url);assert len(calls)<=3
                q=parse_qs(urlparse(url).query,strict_parsing=True)
                assert set(q)=={'page','per_page'} and q['per_page']==['100']
                page=int(q['page'][0]);assert 1<=page<=3
                return [dict(id=n) for n in range((page-1)*100,min(count,page*100))]
            self.assertEqual(len(v.pages('exact/path',fake)),count);self.assertEqual(len(calls),expected_calls)
    def test_repeat_and_bound(self):
        calls=[]
        def fake(url):
            calls.append(url);assert len(calls)<=3
            return [dict(id=n) for n in range(100)]
        with self.assertRaisesRegex(AssertionError,'repeated'):v.pages('path',fake)
        calls=[]
        def endless(url):
            calls.append(url);assert len(calls)<=3
            page=int(parse_qs(urlparse(url).query)['page'][0]);return [dict(id=page*100+n) for n in range(100)]
        with self.assertRaisesRegex(AssertionError,'call cap'):v.pages('path',endless)
        self.assertEqual(len(calls),3)
    def test_caps_each_direction(self):
        for count in [100,101]:
            base=[dict(number=i,parent=None,blocked_by=[],state='open') for i in range(count+1)]
            for mode in ['incoming','outgoing','children']:
                g=deepcopy(base)
                if mode=='incoming':g[0]['blocked_by']=list(range(1,count+1))
                elif mode=='outgoing':
                    for x in g[1:]:x['blocked_by']=[0]
                else:
                    for x in g[1:]:x['parent']=0
                if count==100:v.caps(g)
                else:
                    with self.assertRaises(AssertionError):v.caps(g)
    def test_both_cycles(self):
        g=[dict(number=1,parent=None,blocked_by=[2]),dict(number=2,parent=None,blocked_by=[1])]
        with self.assertRaises(AssertionError):v.caps(g)
        for x in g:x['blocked_by']=[]
        g[0]['parent']=2;g[1]['parent']=1
        with self.assertRaises(AssertionError):v.caps(g)

if __name__=='__main__':unittest.main()
