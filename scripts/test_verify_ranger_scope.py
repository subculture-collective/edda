"""Independent contrary ownership, source, operand and pagination controls."""
from copy import deepcopy
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch
import verify_ranger_scope as v
from verify_druid_scope import pages, Reader

class Controls(unittest.TestCase):
 def setUp(self):self.m=json.loads(v.MANIFEST.read_text())
 def reject(self,change):
  m=deepcopy(self.m);change(m)
  with self.assertRaises((AssertionError,KeyError)):v.validate_plan(m)
 def test_canonical(self):v.validate_plan(self.m)
 def test_rehashed_feature_owner_and_context(self):
  # Coherent reassignment also changes source context, not just a checksum.
  for slug in ['extra-attack','roving','feral-senses','precise-hunter','relentless-hunter','foe-slayer','natures-veil','hunters-lore','superior-hunters-defense']:
   def mutate(m):
    c=next(c for c in m['record_matrix'] if c['record']['kind']=='feature' and c['id'].endswith('/'+slug))
    c['key_owners']={k:'ability-score-improvement' for k in c['key_owners']}
    old=next(x for x in m['tasks'] if x['key']==slug)
    asi=next(x for x in m['tasks'] if x['key']=='ability-score-improvement')
    asi['source_records']+=old['source_records'];old['source_records']=[]
    m['record_owners'][c['id']]='ability-score-improvement'
    m['ownership_sha256']=v.sha(json.dumps(m['record_owners'],sort_keys=True))
   self.reject(mutate)
 def test_every_feature_native_path(self):
  for n in range(1,21):
   def mutate(m):
    row=next(x for x in m['tasks'] if x['key']=='level-'+str(n));row['depends']=[]
   self.reject(mutate)
 def test_later_admission_and_hp_mandatory(self):
  for dep in ['later-entry','hp']:
   self.reject(lambda m:next(x for x in m['tasks'] if x['key']=='level-2')['depends'].remove(dep))
 def test_total_operand(self):
  self.reject(lambda m:next(x for x in m['tasks'] if x['key']=='level-2').update(fixture='Class2 XP300 PB2 from sourceclassrow'))
 def test_unequal_class_total_controls(self):
  for field,value in [('ranger_row',6),('total_row',2),('xp',300),('pb',2),('hp_mode','first10plusCON')]:
   self.reject(lambda m:m['operand_witnesses'][1].update({field:value}))
 def test_source_numeric(self):
  def mutate(m):
   p=next(p for p in m['catalogues'] if p.endswith('character-ranger.json'))
   m['catalogues'][p]['records'][0]['parameters']['srd-5.2.1:mechanic/ranger/hit-die-sides']=12
  self.reject(mutate)
 def test_dependency_cycle(self):
  self.reject(lambda m:next(x for x in m['gates'] if x['key']=='execution-dependencies')['depends'].append('level-2'))
 def test_unknown_owner(self):
  self.reject(lambda m:m['record_matrix'][0]['key_owners'].update({'invented':'invented'}))
 def test_live_refuses_plan(self):
  with self.assertRaises(AssertionError):v.validate_live(self.m,{'graph':[]})
 def test_wrong_pdf(self):
  with tempfile.TemporaryDirectory() as tmp:
   p=Path(tmp)/'wrong.pdf';p.write_bytes(b'wrong source')
   with self.assertRaisesRegex(AssertionError,'wrong PDF'):v.source_check(self.m,p)
 def test_exact_rest_query_pages(self):
  calls=[]
  def api(path):
   calls.append(path)
   return [{'id':i} for i in range(100)] if len(calls)==1 else [{'id':100}]
  self.assertEqual(len(list(pages(api,'repos/o/r/issues?state=all'))),101)
  self.assertEqual(calls,['repos/o/r/issues?state=all&per_page=100&page=1','repos/o/r/issues?state=all&per_page=100&page=2'])
 def test_repeated_rest_page(self):
  with self.assertRaisesRegex(AssertionError,'repeated API page'):list(pages(lambda _: [{'id':i} for i in range(100)],'repos/o/r/issues'))
 def test_cursor_missing_repeat(self):
  for cursor,seen in [(None,set()),('same',{'same'})]:
   with self.assertRaises(AssertionError):v.next_cursor({'endCursor':cursor},seen)
 def test_parent_prefix_and_later_appends(self):
  prefix='accepted Paladin prefix';section='\n<!-- edda-328:start -->\nsection\n<!-- edda-328:end -->\n'
  item=dict(append=section,previous_sha256=v.sha(prefix));h=v.sha(prefix)
  self.assertEqual(v.supplement_plan(prefix,item,prefix,h),prefix+section)
  self.assertEqual(v.supplement_plan(prefix+section+'later append',item,prefix,h),prefix+section+'later append')
  for current in [prefix[:-1],prefix+section+section,prefix+section.replace('section','altered'),prefix+'<!-- edda-328:start -->']:
   with self.assertRaises(AssertionError):v.supplement_plan(current,item,prefix,h)
 def test_api_readonly_and_callcap(self):
  r=Reader();r.calls=350
  with self.assertRaisesRegex(AssertionError,'call cap'):r('graphql',{'query':'{}'})
  with self.assertRaisesRegex(AssertionError,'read-only'):Reader()('graphql',{'query':'mutation{bad}'})
if __name__=='__main__':unittest.main()
