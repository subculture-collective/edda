#!/usr/bin/env python3
"""Read-only Ranger decision328 source/ownership/graph audit; plans never prove publication."""
import argparse
from copy import deepcopy
import hashlib
import json
from pathlib import Path
import re
import subprocess
from verify_druid_scope import Reader, pages
from verify_monk_scope import next_cursor
from verify_barbarian_scope import acyclic
ROOT=Path(__file__).resolve().parents[1]
MANIFEST=ROOT/'docs/investigations/issue-328-ranger-scope.json'
RECEIPT=ROOT/'docs/investigations/issue-328-link-verification.json'
PIN='8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87'
PROJECT='PVT_kwDODjtmGM4BrFJ2'
ADVANCEMENT='internal/rules/srd521/testdata/character-advancement.json'
def sha(v):return hashlib.sha256(v.encode() if isinstance(v,str) else v).hexdigest()
def rows(m):return {x['key']:x for x in m['tasks']+m['gates']}
def path_to(rr,start,target):
 seen=set()
 def walk(k):
  if k==target:return True
  if k in seen:return False
  seen.add(k)
  return any(walk(d) for d in rr[k]['depends'])
 return walk(start)
def validate_plan(m):
 assert m['decision']==328 and m['source']['sha256']==PIN
 actual=[]
 for p,c in m['catalogues'].items():
  assert json.loads((ROOT/p).read_text())==c,'source operands'
  if p.endswith(('character-ranger.json','character-hunter.json')):actual+=c['records']
 assert len(actual)==93 and actual==m['records'],'source93'
 assert len(m['record_matrix'])==93
 rr=rows(m);assert len(rr)==len(m['tasks'])+len(m['gates'])
 acyclic({k:dict(blocked_by=x['depends']) for k,x in rr.items()})
 for x in m['tasks']:
  assert x['type']=='Task' and x['size']=='M'
  assert all(d in rr for d in x['depends'])
 for c,r in zip(m['record_matrix'],actual):
  assert c['record']==r and c['id']==r['id']
  assert set(c['key_owners'])==set(r['required_mechanics']+r['option_keys'])
  assert all(v in rr for v in c['key_owners'].values())
  assert all(c[k] for k in ['rules','game','engine','combat','database','handlers','web','tui'])
 # Independent semantic oracle: every named feature has its actual execution owner;
 # coherent rehashing of the manifest cannot move one into ASI or a selector.
 for r,c in zip(actual,m['record_matrix']):
  if r['kind']=='feature':
   slug=r['id'].rsplit('/',1)[-1]
   expected='tireless-hp' if slug=='tireless' else slug
   for k in r['required_mechanics']:
    owner=expected
    if slug=='tireless' and 'exhaustion' in k:owner='tireless-rest'
    if slug=='fighting-style' and ('/druidic-warrior/' in k or k.endswith('/cantrip-replacements')):owner='druidic-cantrips'
    assert c['key_owners'][k]==owner,'real feature execution leaf'
   assert r in rr[expected]['source_records'],'feature source context'
 for n in range(1,21):
  x=rr['level-'+str(n)];assert 'hp' in x['depends'] and 'later-entry' in x['depends'],'own transition adapters'
  assert x['source_records']==[r for r in actual if r['kind']=='level' and r['level']==n]
  source=x['source_records'][0]
  for f in source['progression']['feature_refs']:
   slug=f.rsplit('/',1)[-1];slug='tireless-hp' if slug=='tireless' else slug
   assert path_to(rr,x['key'],slug),'granting row execution path'
  for f in {3:['hunters-lore','hunters-prey'],7:['defensive-tactics'],11:['superior-hunters-prey'],15:['superior-hunters-defense']}.get(n,[]):assert path_to(rr,x['key'],f),'Hunter timing path'
  assert 'COMPLETE canonical advancement table' in x['fixture'] and 'sum of authoritative admitted class levels AFTER transition' in x['fixture'],'total/class operands'
 assert rr['hp']['prerequisites']==[968] and 429 not in sum([x['prerequisites'] for x in m['tasks']],[]),'shared seam'
 assert path_to(rr,'class-gate','hp') and path_to(rr,'class-gate','shared-gate')
 assert len(m['spell_ownership'])==59
 for spell,owners in m['spell_ownership'].items():
  assert owners
  for owner in owners:
   if isinstance(owner,str):assert rr[owner]['spell']==spell and path_to(rr,'shared-gate',owner)
 for x in m['tasks']:
  expected=[dict(id=c['id'],keys=sorted(k for k,v in c['key_owners'].items() if v==x['key'])) for c in m['record_matrix'] if x['key'] in c['key_owners'].values()]
  assert x['owned_contexts']==expected,'exact owned contexts'
  if not x['key'].startswith('level-'):
   assert x['source_records']==[c['record'] for c in m['record_matrix'] if x['key'] in c['key_owners'].values()],'real contextual source operands'
 for witness,expected in zip(m['operand_witnesses'],[(1,5,6500),(2,6,14000)]):
  ranger,total,xp=expected
  assert witness['after']=={'fighter':4,'ranger':ranger},'authoritative class vector'
  assert witness['ranger_row']==ranger and witness['total_row']==sum(witness['after'].values())==total,'class and total operands'
  canonical=next(r for r in m['catalogues'][ADVANCEMENT]['records'] if r['level']==total)['advancement']
  assert witness['xp']==canonical['experience_points']==xp and witness['pb']==canonical['proficiency_bonus']==3,'complete XP/PB table'
  assert witness['hp_mode']=='later-class-d10-or-fixed6','own later HP'
 for p,h in m['evidence_hashes'].items():assert sha((ROOT/p).read_bytes())==h,'runtime evidence changed'
 return rr

def read_live(m):
 api=Reader();graph=[];cursor=None;seen=set();cursors=set()
 for _ in range(20):
  q='query($cursor:String){repository(owner:"subculture-collective",name:"edda"){issues(first:100,after:$cursor){pageInfo{hasNextPage endCursor}nodes{number title body state parent{number} blockedBy(first:100){pageInfo{hasNextPage}nodes{number}}}}}}'
  b=api('graphql',{'query':q,'variables':{'cursor':cursor}})['data']['repository']['issues']
  for x in b['nodes']:
   assert x['number'] not in seen,'repeated issue page';seen.add(x['number'])
   assert not x['blockedBy']['pageInfo']['hasNextPage'],'incoming cap'
   graph.append(dict(number=x['number'],title=x['title'],body=x['body'] or '',state=x['state'].lower(),parent=(x['parent'] or {}).get('number'),blocked_by=sorted(v['number'] for v in x['blockedBy']['nodes'])))
  if not b['pageInfo']['hasNextPage']:break
  cursor=next_cursor(b['pageInfo'],cursors)
 else:raise AssertionError('graph page cap')
 project={}
 nums=[x['issue'] for x in m['tasks']+m['gates'] if x['issue']]+[968]
 for start in range(0,len(nums),25):
  batch=nums[start:start+25]
  fields=' '.join('i%d:issue(number:%d){issueType{name} projectItems(first:100){pageInfo{hasNextPage}nodes{project{id} fieldValueByName(name:"Size"){... on ProjectV2ItemFieldSingleSelectValue{name}}}}}'%(n,n) for n in batch)
  data=api('graphql',{'query':'{repository(owner:"subculture-collective",name:"edda"){'+fields+'}}'})['data']['repository']
  for n in batch:
   x=data['i'+str(n)];assert not x['projectItems']['pageInfo']['hasNextPage']
   p=[p for p in x['projectItems']['nodes'] if p['project']['id']==PROJECT];assert len(p)==1
   project[str(n)]=dict(type=x['issueType']['name'],size=p[0]['fieldValueByName']['name'])
 return dict(graph=graph,project=project,api_calls=api.calls)

def validate_live(m,r):
 rr=validate_plan(m);assert all(x['issue'] for x in rr.values()),'unpublished plan cannot qualify'
 g={x['number']:x for x in r['graph']};assert len(g)==len(r['graph'])
 assert g[328]['state']=='open','decision must remain open until merge'
 for n in [285,293,315]:assert g[n]['state']=='closed'
 acyclic(g);acyclic({n:dict(blocked_by=[x['parent']] if x['parent'] else []) for n,x in g.items()})
 children={};out={}
 for x in g.values():
  assert len(x['blocked_by'])==len(set(x['blocked_by']))<=100,'incoming cap'
  children[x['parent']]=children.get(x['parent'],0)+1
  for n in x['blocked_by']:out[n]=out.get(n,0)+1
 assert all(v<=100 for k,v in children.items() if k is not None),'child cap'
 assert all(v<=100 for v in out.values()),'outgoing cap'
 assert g[968]['parent']==336 and g[968]['blocked_by']==[429]
 for x in rr.values():
  a=g[x['issue']];expected=set(x['prerequisites'])|{rr[k]['issue'] for k in x['depends']}
  assert a['title']==x['title'] and a['body']==x['body'],'owned title/body'
  assert a['parent']==x.get('parent',rr[x['parent_key']]['issue'] if x.get('parent_key') else None),'owned parent'
  assert set(a['blocked_by'])==expected,'owned prerequisite'
  assert r['project'][str(x['issue'])]==dict(type='Task',size='M'),'metadata'
 for key,roots in [('shared-gate',[279,280]),('class-gate',[337])]:
  assert all(rr[key]['issue'] in g[n]['blocked_by'] for n in roots),'gate root edge'
 for n,old in m['reused'].items():
  a=g[int(n)];assert a['body']==old['body'] and a['parent']==old['parent'],'reused owner body/parent'
  assert set(old['blocked_by'])<=set(a['blocked_by']),'historical prerequisites'
 return dict(incoming={str(n):len(g[n]['blocked_by']) for n in [279,280,337]},outgoing429=out[429],children283=children[283],nodes=len(g))

def source_check(m,pdf):
 assert sha(pdf.read_bytes())==PIN,'wrong PDF'
 for s in ['ranger','hunter','advancement','multiclass','origin_feats','fighting_styles','general_feats','epic_boon_feats']:
  subprocess.run(['python3',str(ROOT/'scripts'/('verify_srd_'+s+'.py')),str(pdf)],check=True)
 raw=subprocess.check_output(['pdftotext','-f','107','-l','175','-raw',str(pdf),'-'],text=True)
 headings=list(re.finditer(r'^([^\n]+)\n(?:Level [0-9] |[A-Za-z]+ Cantrip)',raw,re.M))
 for x in m['tasks']:
  if not x.get('source_excerpt'):continue
  found=[i for i,h in enumerate(headings) if h.group(1)==x['spell']];assert len(found)==1
  i=found[0];text=raw[headings[i].start():headings[i+1].start() if i+1<len(headings) else len(raw)]
  text=re.sub(r'\f?System Reference Document 5.2.1\n\d+\n','',text).strip()
  assert text==x['source_excerpt'],'spell operands'

def project_capacity(m,r):
 """Temporary preflight only: proposed edges are never recorded as live evidence."""
 rr=validate_plan(m);g={x['number']:deepcopy(x) for x in r['graph']}
 for key,x in rr.items():
  node=x['issue'] or 'planned:'+key
  assert node not in g or x['issue'], 'proposal collision'
  g[node]=dict(parent=x.get('parent',rr[x['parent_key']]['issue'] or 'planned:'+x['parent_key'] if x.get('parent_key') else None),blocked_by=x['prerequisites']+[rr[d]['issue'] or 'planned:'+d for d in x['depends']])
  for root in x.get('roots',[]):
   if node not in g[root]['blocked_by']:g[root]['blocked_by'].append(node)
 outgoing={};children={}
 for n,x in g.items():
  assert len(x['blocked_by'])==len(set(x['blocked_by']))<=100,'projected incoming cap'
  for d in x['blocked_by']:outgoing[d]=outgoing.get(d,0)+1
  children[x['parent']]=children.get(x['parent'],0)+1
 assert all(v<=100 for v in outgoing.values()),'projected outgoing cap'
 assert all(v<=100 for k,v in children.items() if k is not None),'projected child cap'
 acyclic(g);acyclic({n:dict(blocked_by=[x['parent']] if x['parent'] else []) for n,x in g.items()})
 return dict(incoming={str(n):len(g[n]['blocked_by']) for n in [279,280,337]},outgoing={str(n):outgoing.get(n,0) for n in [61,64,65,66,67,273,353,429]},children={str(n):children.get(n,0) for n in [277,283,328,336]})

def supplement_plan(current,supplement,prefix,pinned_hash):
 assert sha(prefix)==pinned_hash and supplement['previous_sha256']==pinned_hash,'accepted immutable prefix'
 assert current.startswith(prefix),'prefix rollback or alteration'
 start='<!-- edda-328:start -->';end='<!-- edda-328:end -->';section=supplement['append']
 assert current.count(start)==current.count(end)==current.count(section)<=1,'partial altered or duplicate section'
 desired=current if section in current else current+section
 assert desired.startswith(current) and desired.count(section)==1 and len(desired)<=65536
 return desired

def main():
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--plan',action='store_true');p.add_argument('--capacity-live',action='store_true');p.add_argument('--live',action='store_true');p.add_argument('--pdf',type=Path);a=p.parse_args()
 m=json.loads(MANIFEST.read_text());validate_plan(m)
 if a.pdf:source_check(m,a.pdf)
 if a.capacity_live:
  print('PROJECTION ONLY, not publication:',project_capacity(m,read_live(m)));return
 if a.plan:print('PASS planned source/ownership DAG only; NO published graph or final prefix proof');return
 assert m['accepted_parent_prefix'] is not None,'post-Paladin immutable prefix release pending'
 r=read_live(m) if a.live else json.loads(RECEIPT.read_text());counts=validate_live(m,r)
 if a.live:RECEIPT.write_text(json.dumps(r,indent=2)+'\n')
 print('PASS live native receipt',counts)
if __name__=='__main__':main()
