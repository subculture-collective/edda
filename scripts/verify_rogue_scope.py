#!/usr/bin/env python3
"""Verify Rogue audit checkpoint. Publication/sealed proof deliberately fails while pending."""
import argparse
from collections import Counter
import hashlib
import json
from pathlib import Path
import subprocess
from verify_fighter_scope import graph_read, acyclic
from rogue_capacity import project

ROOT=Path(__file__).resolve().parents[1]
PATH=ROOT/'docs/investigations/issue-329-rogue-scope.json'
RECEIPT=ROOT/'docs/investigations/issue-329-native-observation.json'
PIN='8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87'
AUDIT_SEAL='c25cf037e54639956bf0cb1bd3617c9bdf817049c054f75d6417e7abeb260d76'
OBSERVATION_SEAL='a09589a35cac1eb94d6936a3c2c035b15851d8571a286e19e0b30bb065962dc1'

def digest(x):return hashlib.sha256(json.dumps(x,sort_keys=True).encode()).hexdigest()

def pages(path,api,limit=3):
    assert '?' not in path
    seen=set();result=[]
    for page in range(1,limit+1):
        batch=api(f'{path}?per_page=100&page={page}')
        assert isinstance(batch,list) and len(batch)<=100
        ids=[x['id'] for x in batch]
        assert len(set(ids))==len(ids) and not seen.intersection(ids),'repeated page'
        seen.update(ids);result.extend(batch)
        if len(batch)<100:return result
    raise AssertionError('pagination call cap')

def parent_append(prefix,section,actual):
    assert actual.startswith(prefix+section),'original-only rollback or changed prefix/section'
    assert actual.count('<!-- edda-329:start -->')==actual.count('<!-- edda-329:end -->')==1,'altered duplicate marker'
    assert len(actual)<=65536

def caps(graph):
    g={x['number']:x for x in graph};assert len(g)==len(graph)
    for x in graph:
        assert len(x['blocked_by'])==len(set(x['blocked_by']))<=100
        assert set(x['blocked_by'])<=set(g),'missing dependency'
        assert x['parent'] is None or x['parent'] in g
    outgoing=Counter(n for x in graph for n in x['blocked_by']);children=Counter(x['parent'] for x in graph if x['parent'] is not None)
    assert max(outgoing.values(),default=0)<=100
    assert max(children.values(),default=0)<=100
    acyclic(g);acyclic({n:dict(blocked_by=[x['parent']] if x['parent'] is not None else []) for n,x in g.items()})

def validate(m):
    assert digest(m)==AUDIT_SEAL,'independent audited checkpoint seal'
    assert m['source_sha256']==PIN and m['base']=='952fab917329e00cceebf46a8bacd59cb1c30e99'
    assert m['accepted_post_ranger'] is None and m['publication'].startswith('PENDING')
    records=[]
    for name,c in m['catalogues'].items():
        assert c==json.loads((ROOT/f'internal/rules/srd521/testdata/character-{name}.json').read_text())
        assert c['source']['sha256']==PIN
        if name in ['rogue','thief']:records+=c['records']
    assert len(records)==107 and records==[x['record'] for x in m['records']]
    rows={x['key']:x for x in m['issues']};assert len(rows)==len(m['issues'])==66
    for x in m['records']:
        r=x['record'];assert set(x['key_owners'])==set(r['required_mechanics']+r['option_keys'])
        assert x['state']=='metadata only; Rogue/Thief unavailable'
        assert set(x['execution_owners'])<=set(rows)
        for owners in x['key_owners'].values():assert owners and set(owners)<=set(rows)
        if r['kind']=='feature' and not r['id'].endswith('/ability-score-improvement'):
            assert 'improvement' not in x['execution_owners'],'unrelated ASI owner'
    for n in range(1,21):
        row=rows[f'level-{n}'];assert 'hp' in row['depends'] and 'multiclass' in row['depends']
        assert set(m['grants'][str(n)])<=set(row['depends'])
        if n>1:assert f'level-{n-1}' in row['depends']
    assert {'evasion','reliable-talent'}<=set(rows['level-7']['depends'])
    assert {'hp','multiclass'}<=set(rows['later-1']['depends'])
    assert rows['hp']['prerequisites']==[968,58]
    assert not any(429 in x['prerequisites'] for x in rows.values())
    for x in rows.values():
        assert x['issue'] is None and x['publication']=='UNPUBLISHED'
        assert set(x['depends'])<=set(rows)
    adv={x['level']:x['advancement'] for x in m['catalogues']['advancement']['records']}
    assert adv[5]==dict(experience_points=6500,proficiency_bonus=3)
    assert adv[6]==dict(experience_points=14000,proficiency_bonus=3)
    for p,x in m['evidence'].items():assert hashlib.sha256((ROOT/p).read_bytes()).hexdigest()==x['sha256']
    for n in [47,54,283,337]:
        section=(ROOT/f'docs/investigations/issue-329-parent-{n}-PENDING.md').read_text()
        assert len(section)<=(90 if n==47 else 750)
        parent_append('',section,section)
    return m

def validate_observation(m,r):
    assert digest(r)==OBSERVATION_SEAL,'independent native observation seal'
    assert r['status']=='OBSERVATION ONLY; NO OWN NATIVE PUBLICATION; PENDING'
    assert r['owned_issues']==[] and r['observed_at'] and r['limits']
    caps(r['graph']);g={x['number']:x for x in r['graph']}
    for n in [285,294,316]:assert g[n]['state']=='closed'
    assert g[329]['parent']==283 and g[329]['blocked_by']==[285,294,316]
    assert g[968]['parent']==336 and g[968]['blocked_by']==[429]
    for n,c in r['consumed_contracts'].items():assert c['number']==int(n)
    c=r['consumed_contracts']['779'];assert hashlib.sha256(c['body'].encode()).hexdigest()=='556d1618d1cd2bd9c956fe38861d89e94ee14ff61a505b5d2b6a1fa46d19d150'
    assert c['parent']==50 and c['blocked_by']==[61,64,65,66,67,488,778]
    assert project(m,r['graph'])==json.loads((ROOT/'docs/investigations/issue-329-capacity-proposal.json').read_text())

def main():
    p=argparse.ArgumentParser(description=__doc__);p.add_argument('--pdf',required=True);p.add_argument('--live',action='store_true');p.add_argument('--publication',action='store_true');p.add_argument('--sealed',action='store_true');a=p.parse_args()
    assert hashlib.sha256(Path(a.pdf).read_bytes()).hexdigest()==PIN
    m=validate(json.loads(PATH.read_text()));r=json.loads(RECEIPT.read_text());validate_observation(m,r)
    for name in ['rogue','thief','advancement','multiclass','origin_feats','general_feats','fighting_styles','epic_boon_feats']:
        subprocess.run(['rtk','proxy','python3',str(ROOT/f'scripts/verify_srd_{name}.py'),a.pdf],check=True,stdout=subprocess.DEVNULL)
    if a.live:
        graph=graph_read();caps(graph);g={x['number']:x for x in graph}
        for n,c in r['consumed_contracts'].items():
            actual=g[int(n)]
            assert all(actual[k]==c[k] for k in ['title','body','parent','blocked_by']),f'consumed owner drift#{n}'
        project(m,graph)
    assert not(a.publication or a.sealed),'PENDING: native publication, item scope and accepted POST-RANGER prefix unavailable; checkpoint cannot qualify closure'
    print('Rogue checkpoint passes:107 records,20 rows,66 proposed terminals; PUBLICATION/SEALED proof PENDING')

if __name__=='__main__':main()
