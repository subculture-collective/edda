#!/usr/bin/env python3
"""Independent source envelope and operand proof for all 50 consumed Paladin spells."""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT=Path(__file__).resolve().parents[1]
PIN='8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87'
# Source transcription, independent of the decision manifest and its clause table.
ENVELOPES={
'Aid':(2,'Action','30 feet','8 hours'),
'Aura of Life':(4,'Action','Self','Concentration, up to 10 minutes'),
'Banishment':(4,'Action','30 feet','Concentration, up to 1 minute'),
'Beacon of Hope':(3,'Action','30 feet','Concentration, up to 1 minute'),
'Bless':(1,'Action','30 feet','Concentration, up to 1 minute'),
'Command':(1,'Action','60 feet','Instantaneous'),
'Commune':(5,'1 minute or Ritual','Self','1 minute'),
'Create Food and Water':(3,'Action','30 feet','Instantaneous'),
'Cure Wounds':(1,'Action','Touch','Instantaneous'),
'Daylight':(3,'Action','60 feet','1 hour'),
'Death Ward':(4,'Action','Touch','8 hours'),
'Detect Evil and Good':(1,'Action','Self','Concentration, up to 10 minutes'),
'Detect Magic':(1,'Action or Ritual','Self','Concentration, up to 10 minutes'),
'Detect Poison and Disease':(1,'Action or Ritual','Self','Concentration, up to 10 minutes'),
'Dispel Evil and Good':(5,'Action','Self','Concentration, up to 1 minute'),
'Dispel Magic':(3,'Action','120 feet','Instantaneous'),
'Divine Favor':(1,'Bonus Action','Self','1 minute'),
'Divine Smite':(1,'Bonus Action, which you take immediately after hitting a target with a Melee weapon or an Unarmed Strike','Self','Instantaneous'),
'Find Steed':(2,'Action','30 feet','Instantaneous'),
'Flame Strike':(5,'Action','60 feet','Instantaneous'),
'Freedom of Movement':(4,'Action','Touch','1 hour'),
'Geas':(5,'1 minute','60 feet','30 days'),
'Gentle Repose':(2,'Action or Ritual','Touch','10 days'),
'Greater Restoration':(5,'Action','Touch','Instantaneous'),
'Guardian of Faith':(4,'Action','30 feet','8 hours'),
'Guidance':(0,'Action','Touch','Concentration, up to 1 minute'),
'Heroism':(1,'Action','Touch','Concentration, up to 1 minute'),
'Lesser Restoration':(2,'Bonus Action','Touch','Instantaneous'),
'Light':(0,'Action','Touch','1 hour'),
'Locate Creature':(4,'Action','Self','Concentration, up to 1 hour'),
'Locate Object':(2,'Action','Self','Concentration, up to 10 minutes'),
'Magic Circle':(3,'1 minute','10 feet','1 hour'),
'Magic Weapon':(2,'Bonus Action','Touch','1 hour'),
'Mending':(0,'1 minute','Touch','Instantaneous'),
'Prayer of Healing':(2,'10 minutes','30 feet','Instantaneous'),
'Protection from Evil and Good':(1,'Action','Touch','Concentration up to 10 minutes'),
'Protection from Poison':(2,'Action','Touch','1 hour'),
'Purify Food and Drink':(1,'Action or Ritual','10 feet','Instantaneous'),
'Raise Dead':(5,'1 hour','Touch','Instantaneous'),
'Remove Curse':(3,'Action','Touch','Instantaneous'),
'Resistance':(0,'Action','Touch','Concentration, up to 1 minute'),
'Revivify':(3,'Action','Touch','Instantaneous'),
'Sacred Flame':(0,'Action','60 feet','Instantaneous'),
'Searing Smite':(1,'Bonus Action, which you take immediately after hitting a target with a Melee weapon or an Unarmed Strike','Self','1 minute'),
'Shield of Faith':(1,'Bonus Action','60 feet','Concentration, up to 10 minutes'),
'Shining Smite':(2,'Bonus Action, which you take immediately after hitting a creature with a Melee weapon or an Unarmed Strike','Self','Concentration, up to 1 minute'),
'Spare the Dying':(0,'Action','15 feet','Instantaneous'),
'Thaumaturgy':(0,'Action','30 feet','Up to 1 minute'),
'Warding Bond':(2,'Action','Touch','1 hour'),
'Zone of Truth':(2,'Action','60 feet','10 minutes')}

OPERANDS={
'Bless': ['up to three creatures','adds 1d4','one additional creature for each spell slot level above 1','a Holy Symbol worth 5+ GP'],
'Command':['Approach.','within 5 feet','Drop.','Flee.','Grovel.','Halt.','no action or Bonus Action','one additional creature'],
'Detect Magic':['within 30 feet','take the Magic action','1 foot of stone, dirt, or wood; 1 inch of metal; or a thin sheet of lead'],
'Detect Poison and Disease':['magical contagions within 30 feet','kind of poison, creature, or contagion','1 foot of stone, dirt, or wood; 1 inch of metal; or a thin sheet of lead'],
'Divine Smite':['extra 2d8 Radiant','increases by 1d8 if the target is a Fiend or an Undead','increases by 1d8 for each spell slot level above 1'],
'Find Steed':['Large, rideable','Celestial, Fey, or Fiend','controlled mount','Incapacitated','immediately after yours','drops to 0 Hit Points or if you die'],
'Prayer of Healing':['Up to five creatures','entire casting','benefits of a Short Rest','2d8 Hit Points','until that creature finishes a Long Rest','increases by 1d8'],
'Magic Circle':['10-foot-radius, 20-foot-tall Cylinder','100+ GP','reverse direction','increases by 1 hour'],
'Magic Weapon':['nonmagical weapon','+1 bonus to attack rolls and damage rolls','+2 with a level 3–5 spell slot','+3 with a level 6+ spell slot'],
'Geas':['5d10 Psychic','no more than once each day','automatically succeeds','suicidal command','7 or 8','365 days','level 9'],
'Raise Dead':['no longer than 10 days','wasn’t Undead','1 Hit Point','500+ GP','−4 penalty to D20 Tests','reduced by 1'],
'Resistance':['Acid, Bludgeoning, Cold, Fire, Lightning, Necrotic, Piercing, Poison, Radiant, Slashing, or Thunder','reduces the total damage taken by 1d4','only once per turn'],
'Sacred Flame':['Dexterity saving throw','1d8 Radiant','no benefit from Half Cover or Three-Quarters Cover','5 (2d8), 11 (3d8), and 17 (4d8)'],
'Spare the Dying':['0 Hit Points and isn’t dead','becomes Stable','5 (30 feet), 11 (60 feet), and 17 (120 feet)'],
'Warding Bond':['another creature that is willing','50+ GP each','within 60 feet','+1 bonus to AC and saving throws','same amount of damage','you drop to 0 Hit Points','cast again on either'],
'Guardian of Faith':['20 Radiant damage','total of 60 damage','first time on a turn or starts its turn'],
'Flame Strike':['10-foot-radius, 40-foot-high Cylinder','5d6 Fire damage and 5d6 Radiant','Fire damage and the Radiant damage increase by 1d6'],
'Locate Creature':['1,000 feet','within 30 feet','different form','any thickness of lead'],
'Locate Object':['1,000 feet','within 30 feet','any thickness of lead']}


def source(pdf):
    assert hashlib.sha256(pdf.read_bytes()).hexdigest()==PIN
    text=subprocess.run(['rtk','proxy','pdftotext','-f','104','-l','175','-raw',str(pdf),'-'],capture_output=True,text=True,check=True).stdout
    pages=text.split('\f');lines=[]
    for p in pages:
        page=p.splitlines()
        if page and page[0].startswith('System Reference Document'):page=page[2:]
        lines+=page
    inventory=json.loads((ROOT/'docs/investigations/issue-80-spell-item-inventory.json').read_text())
    names={x['name'] for x in inventory['records'] if x['family']=='spells'}
    starts=[i for i,l in enumerate(lines[:-1]) if l in names and (lines[i+1].startswith('Level ') or 'Cantrip' in lines[i+1])]
    blocks={}
    for pos,i in enumerate(starts):
        end=starts[pos+1] if pos+1<len(starts) else len(lines)
        value='\n'.join(lines[i:end])
        # Preserve the source compound unit in wrapped "40-foot-high".
        value=re.sub(r'(?<=foot)-\s*\n','-',value)
        value=re.sub(r'-\s*\n','',value)
        blocks[lines[i]]=' '.join(value.split())
    return blocks


def validate(blocks):
    for name,(level,casting,range_,duration) in ENVELOPES.items():
        text=blocks[name]
        found=re.search(r'Level (\d+)',text.split('Casting Time:')[0])
        assert (int(found[1]) if found else 0)==level,(name,'level')
        assert text.split('Casting Time: ',1)[1].split(' Range:',1)[0]==casting,(name,'casting')
        assert text.split('Range: ',1)[1].split(' Component',1)[0]==range_,(name,'range')
        assert text.split('Duration: ',1)[1].startswith(duration),(name,'duration')
    for name,clauses in OPERANDS.items():
        for clause in clauses:assert clause in blocks[name],(name,clause)


def main():
    blocks=source(Path(sys.argv[1]));validate(blocks)
    # Each corrupted envelope and existing numeric clause must be independently rejected.
    count=0
    for name in ENVELOPES:
        bad=dict(blocks);bad[name]=bad[name].replace('Range:','Wrong range:',1)
        try:validate(bad)
        except (AssertionError,IndexError):count+=1
        else:raise AssertionError('accepted source envelope corruption')
    operand_controls=0
    for name,clauses in OPERANDS.items():
        for clause in clauses:
            bad=dict(blocks);bad[name]=bad[name].replace(clause,'CORRUPTED SOURCE OPERAND')
            try:validate(bad)
            except (AssertionError,IndexError):operand_controls+=1
            else:raise AssertionError('accepted numeric/alternative source corruption')
    print(f'PASS:50 exact consumed spell envelopes;{sum(map(len,OPERANDS.values()))} source-bound numeric/alternative clauses;{count} contrary envelopes/{operand_controls} contrary operands rejected')


if __name__=='__main__':main()
