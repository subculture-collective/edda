#!/usr/bin/env python3
"""Independent source comparator for stat-block batch 02 (#152).

Re-reads the pinned PDF columns, rebuilds every record from the printed block
and compares the whole catalogue. No Go, network or runtime is involved.
"""
import hashlib
import json
from pathlib import Path
import re
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]
PIN = '8974902d109d6e63672d7c490bde9ccf052410503d9cfa768237154fbc5e3d87'
CATALOGUE = ROOT/'internal/rules/srd521/testdata/monsters-statblocks-02.json'
NOTICE = ('Adapted the eleven stat blocks Allosaurus and Ancient Black, Blue, Brass, Bronze, Copper, Gold, Green, Red, '
          'Silver and White Dragon on PDF pages 265–344 into typed mechanical fields and implementation keys, and '
          'recorded the Animals section heading on PDF page 344 as an index entry without a stat block. Retained every '
          'skill, immunity, sense, trait, action, spell list and legendary action, with continuation text across '
          'columns and pages. Mechanic keys and typed details are Edda additions. No playable NPC behavior is granted.')

# Printed trait and action names per record, in source order. Pages are inclusive.
BLOCKS = [
    ('animals/allosaurus', 'Allosaurus', (344, 344), [], ['Bite', 'Claws'], []),
    ('monsters/ancient-black-dragon', 'Ancient Black Dragon', (265, 265), ['Amphibious', 'Legendary Resistance'],
     ['Multiattack', 'Rend', 'Acid Breath', 'Spellcasting'], ['Cloud of Insects', 'Frightful Presence', 'Pounce']),
    ('monsters/ancient-blue-dragon', 'Ancient Blue Dragon', (267, 267), ['Legendary Resistance'],
     ['Multiattack', 'Rend', 'Lightning Breath', 'Spellcasting'], ['Cloaked Flight', 'Sonic Boom', 'Tail Swipe']),
    ('monsters/ancient-brass-dragon', 'Ancient Brass Dragon', (269, 269), ['Legendary Resistance'],
     ['Multiattack', 'Rend', 'Fire Breath', 'Sleep Breath', 'Spellcasting'], ['Blazing Light', 'Pounce', 'Scorching Sands']),
    ('monsters/ancient-bronze-dragon', 'Ancient Bronze Dragon', (271, 271), ['Amphibious', 'Legendary Resistance'],
     ['Multiattack', 'Rend', 'Lightning Breath', 'Repulsion Breath', 'Spellcasting'], ['Guiding Light', 'Pounce', 'Thunderclap']),
    ('monsters/ancient-copper-dragon', 'Ancient Copper Dragon', (277, 277), ['Legendary Resistance'],
     ['Multiattack', 'Rend', 'Acid Breath', 'Slowing Breath', 'Spellcasting'], ['Giggling Magic', 'Mind Jolt', 'Pounce']),
    ('monsters/ancient-gold-dragon', 'Ancient Gold Dragon', (292, 292), ['Amphibious', 'Legendary Resistance'],
     ['Multiattack', 'Rend', 'Fire Breath', 'Spellcasting', 'Weakening Breath'], ['Banish', 'Guiding Light', 'Pounce']),
    ('monsters/ancient-green-dragon', 'Ancient Green Dragon', (294, 295), ['Amphibious', 'Legendary Resistance'],
     ['Multiattack', 'Rend', 'Poison Breath', 'Spellcasting'], ['Mind Invasion', 'Noxious Miasma', 'Pounce']),
    ('monsters/ancient-red-dragon', 'Ancient Red Dragon', (319, 319), ['Legendary Resistance'],
     ['Multiattack', 'Rend', 'Fire Breath', 'Spellcasting'], ['Commanding Presence', 'Fiery Rays', 'Pounce']),
    ('monsters/ancient-silver-dragon', 'Ancient Silver Dragon', (325, 325), ['Legendary Resistance'],
     ['Multiattack', 'Rend', 'Cold Breath', 'Paralyzing Breath', 'Spellcasting'], ['Chill', 'Cold Gale', 'Pounce']),
    ('monsters/ancient-white-dragon', 'Ancient White Dragon', (341, 341), ['Ice Walk', 'Legendary Resistance'],
     ['Multiattack', 'Rend', 'Cold Breath'], ['Freezing Burst', 'Frightful Presence', 'Pounce']),
]
HEADINGS = [('animals/animals', 'Animals', (344, 344))]
SIZES = 'Tiny|Small|Medium|Large|Huge|Gargantuan'
HEAD = re.compile(r'^(' + SIZES + r') ([A-Z][a-z]+)(?: \(([A-Za-z]+)\))?, ([A-Z][A-Za-z ]+)$')
ABILITY = {'Strength': 'str', 'Dexterity': 'dex', 'Constitution': 'con', 'Intelligence': 'int', 'Wisdom': 'wis', 'Charisma': 'cha'}
WORDS = {'one': 1, 'two': 2, 'three': 3}
DICE = r'(\d+) \((\d+)d(\d+)(?: \+ (\d+))?\)'
# Effect key, display name and source clause. Every match is retained verbatim.
EFFECTS = [
    ('recharge', 'Recharge', r'Recharge \d–\d'),
    ('saving-throw', 'Saving throw', r'(?:Strength|Dexterity|Constitution|Intelligence|Wisdom|Charisma) Saving Throw: DC \d+'),
    ('line-area', 'Line area', r'each creature in an? \d+-foot-long, \d+-foot-wide Line'),
    ('cone-area', 'Cone area', r'each creature (?:that isn’t currently affected by this breath )?in a \d+-foot Cone'),
    ('unaffected-targets-only', 'Unaffected targets only', r'that isn’t currently affected by this breath'),
    ('sphere-area', 'Sphere area', r'each creature in a \d+-foot-radius Sphere centered on a point the dragon can see within \d+ feet'),
    ('single-target-in-range', 'Single target in range', r'one creature the dragon can see within \d+ feet'),
    ('half-damage-on-success', 'Half damage on success', r'Success: Half damage(?: only)?'),
    ('repeated-save', 'Repeated save', r'repeats the save'),
    ('second-failure-escalation', 'Second failure escalation', r'Second Failure'),
    ('timed-automatic-success', 'Timed automatic success', r'After 1 minute, it succeeds automatically'),
    ('ends-on-success', 'Ends on success', r'ending the effect on itself on a success'),
    ('ends-on-damage-or-wake-action', 'Ends on damage or wake action', r'This effect ends for the target if it takes damage or a creature within 5 feet of it takes an action to wake it'),
    ('forced-movement', 'Forced movement', r'pushed up to \d+ feet straight away from the dragon'),
    ('no-reactions', 'No Reactions', r'can’t take Reactions'),
    ('speed-halved', 'Speed halved', r'(?:its|the target’s) Speed is halved'),
    ('speed-zero', 'Speed zero', r'the target’s Speed is 0'),
    ('action-or-bonus-action-only', 'Action or Bonus Action only', r'it can take either an action or a Bonus Action on its turn, not both'),
    ('strength-d20-test-disadvantage', 'Strength D20 Test Disadvantage', r'Disadvantage on Strength-based D20 Tests'),
    ('damage-roll-penalty', 'Damage roll penalty', r'subtracts \d+ \(\d+d\d+\) from its damage rolls'),
    ('concentration-save-disadvantage', 'Concentration save Disadvantage', r'Disadvantage on saving throws to maintain Concentration'),
    ('d20-test-penalty', 'D20 Test penalty', r'rolls \d+d\d+ whenever it makes an ability check or attack roll and subtracts the number rolled from the D20 Test'),
    ('armor-class-penalty', 'Armor Class penalty', r'takes a −\d+ penalty to AC'),
    ('demiplane-banishment', 'Demiplane banishment', r'is transported to a harmless demiplane'),
    ('reappearance-placement', 'Reappearance placement', r'it reappears in an unoccupied space of the dragon’s choice within \d+ feet of the dragon'),
    ('once-per-turn-cycle', 'Once per turn cycle', r'can’t take this action again until the start of its next turn'),
    ('half-speed-movement', 'Half Speed movement', r'moves up to half its Speed'),
    ('half-fly-speed-movement', 'Half Fly Speed movement', r'it can fly up to half its Fly Speed'),
    ('attack-sequence', 'Attack sequence', r'makes (?:one|two|three) [A-Z][a-z]+ attacks?'),
    ('multiattack-replacement', 'Multiattack replacement', r'It can replace (?:one|two) attacks? with a use of [^.]+'),
    ('spellcasting-use', 'Spellcasting use', r'uses Spellcasting to cast'),
    ('self-target', 'Self target', r'on itself, and'),
    ('no-material-components', 'No Material components', r'requiring no Material components'),
    ('charisma-spellcasting-ability', 'Charisma spellcasting ability', r'using Charisma as the spellcasting ability'),
    ('spell-save-dc', 'Spell save DC', r'spell save DC \d+'),
    ('spell-attack-bonus', 'Spell attack bonus', r'\+\d+ to hit with spell attacks'),
    ('at-will-casting', 'At Will casting', r'At Will:'),
    ('daily-casting', 'Daily casting', r'1/Day Each:'),
    ('target-size-limit', 'Target size limit', r'If the target is a Large or smaller creature'),
    ('straight-line-movement-trigger', 'Straight-line movement trigger', r'the allosaurus moved 30\+ feet straight toward it immediately before the hit'),
    ('bonus-attack', 'Bonus attack', r'the allosaurus can make one Bite attack against it'),
    ('shapechange-form-limits', 'Shapechange form limits', r'Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell'),
    ('duration', 'Duration', r'[Uu]ntil the (?:end|start) of (?:its|the target’s|the dragon’s) next turn|for \d+ minutes?'),
]


def slug(name):
    return re.sub(r'[^a-z0-9]+', '-', name.lower().replace('’', '')).strip('-')


def signed(text):
    return int(text.replace('−', '-').replace('+', ''))


def dice(m, kind):
    return dict(average=int(m[1]), dice=int(m[2]), sides=int(m[3]), modifier=int(m[4] or 0), type=kind)


def column_lines(pdf, first, last):
    out = []
    for page in range(first, last + 1):
        for x in (0, 306):
            raw = subprocess.check_output(['pdftotext', '-f', str(page), '-l', str(page), '-x', str(x), '-y', '0',
                                           '-W', '306', '-H', '792', '-layout', str(pdf), '-'], text=True)
            for line in raw.replace('­', '').replace('\f', '\n').split('\n'):
                line = ' '.join(line.split())
                if line and not re.fullmatch(r'\d+ System Reference Document 5\.2\.1', line):
                    out.append(line)
    return out


def join(lines):
    text = ''
    for line in lines:
        if text.endswith('-'):
            text = (text if text.endswith('foot-') else text[:-1]) + line
        else:
            text = (text + ' ' + line).strip()
    return text


def block_lines(pdf, name, pages):
    # One extra page lets the reader see where a block that fills its last page ends.
    lines = column_lines(pdf, pages[0], pages[1] + 1)
    starts = [i for i in range(len(lines) - 1) if lines[i] == name and HEAD.match(lines[i + 1])]
    assert len(starts) == 1, name
    end = next(i for i in range(starts[0] + 2, len(lines) - 1) if HEAD.match(lines[i + 1]))
    body = lines[starts[0]:end]
    while not body[-1].endswith('.'):
        body.pop()  # family heading printed above the next block
    return body


def split_entries(text, names):
    marks = []
    at = 0
    for name in names:
        m = re.compile(re.escape(name) + r'(?: \([^)]*\))?\. ').search(text, at)
        assert m, name
        marks.append(m.start())
        at = m.end()
    assert not marks or marks[0] == 0, text[:40]
    return [text[a:b].strip() for a, b in zip(marks, marks[1:] + [len(text)])]


def spell_list(text):
    spells = []
    for label, frequency in (('At Will:', 'at-will'), ('1/Day Each:', '1/day')):
        part = text.split(label)[1].split('1/Day Each:')[0]
        for item in re.split(r', (?![^(]*\))', part.strip()):
            m = re.fullmatch(r'([A-Z][A-Za-z ]+?)(?: \(level (\d) version\))?(?: \((Beast[^)]*)\))?', item.strip())
            assert m, item
            spells.append(dict(name=m[1], level=int(m[2] or 0), frequency=frequency, restriction=m[3] or ''))
    return spells


def action(name, text, section, spells):
    body = text[text.index('. ') + 2:]
    effects, spans = [], []
    for key, label, pattern in EFFECTS:
        for m in re.finditer(pattern, text):
            spans.append((m.start(), m.end(), 'effect/' + key, label, m[0]))
    for m in re.finditer(DICE + r' ([A-Z][a-z]+) damage', body):
        at = text.index(body) + m.start()
        spans.append((at, at + len(m[0]), 'effect/%s-damage' % m[5].lower(), m[5] + ' damage', m[0]))
    for m in re.finditer(r'has the ([A-Z][a-z]+) condition', text):
        spans.append((m.start(), m.end(), 'effect/%s-condition' % m[1].lower(), m[1] + ' condition', m[0]))
    for spell in sorted({s['name'] for s in spells}, key=len, reverse=True):
        for m in re.finditer(r'\b' + re.escape(spell) + r'(?: \(level \d version\))?', body):
            at = text.index(body) + m.start()
            spans.append((at, at + len(m[0]), 'effect/spell-' + slug(spell), spell + ' spell', m[0]))
    cycle = [s for s in spans if s[2] == 'effect/once-per-turn-cycle']
    spans = [s for s in spans if not (s[2] == 'effect/duration' and any(c[0] <= s[0] and s[1] <= c[1] for c in cycle))]
    for _, _, key, label, clause in sorted(spans):
        effects.append(dict(key=key, name=label, text=clause))
    detail = dict(action_id=slug(name), section=section, recharge=None, save=None, area=None, target_range_feet=0,
                  damage=[], conditions=[], push_feet=0, duration_minutes=0, penalties=[], once_per_turn_cycle=False,
                  movement='', attacks=[], replacement=None, casts=[], spellcasting=None)
    record = dict(id=slug(name), key=('legendary-action/' if section == 'legendary-action' else 'action/') + slug(name),
                  name=name, text=text, effects=effects)
    m = re.search(r'Melee Attack Roll: \+(\d+)(?: to hit)?, reach (\d+) ft\. Hit: ' + DICE + r' ([A-Z][a-z]+) damage', body)
    if m:
        record['key'] = 'action/melee-attack'
        record['attack'] = dict(kind='melee', bonus=int(m[1]), reach=int(m[2]), range=0, long_range=0,
                                damage=dict(average=int(m[3]), dice=int(m[4]), sides=int(m[5]), modifier=int(m[6] or 0)),
                                damage_type=m[7].lower())
    if name == 'Multiattack':
        record['key'] = 'action/multiattack'
    trigger = 'hit' if m else 'failure'
    for d in re.finditer(DICE + r' ([A-Z][a-z]+) damage', body):
        detail['damage'].append(dict(dice(d, d[5].lower()), trigger=trigger))
    m = re.search(r'Recharge (\d)–(\d)', text)
    if m:
        detail['recharge'] = dict(min=int(m[1]), max=int(m[2]))
    m = re.search(r'(\w+) Saving Throw: DC (\d+)', body)
    if m:
        detail['save'] = dict(ability=ABILITY[m[1]], dc=int(m[2]), half_on_success='Success: Half damage' in body)
    m = re.search(r'in an? (\d+)-foot-long, (\d+)-foot-wide Line', body)
    if m:
        detail['area'] = dict(shape='line', feet=int(m[1]), width_feet=int(m[2]), range_feet=0)
    m = re.search(r'in a (\d+)-foot Cone', body)
    if m:
        detail['area'] = dict(shape='cone', feet=int(m[1]), width_feet=0, range_feet=0)
    m = re.search(r'in a (\d+)-foot-radius Sphere centered on a point the dragon can see within (\d+) feet', body)
    if m:
        detail['area'] = dict(shape='sphere', feet=int(m[1]), width_feet=0, range_feet=int(m[2]))
    m = re.search(r'one creature the dragon can see within (\d+) feet', body)
    if m:
        detail['target_range_feet'] = int(m[1])
    detail['conditions'] = [c.lower() for c in re.findall(r'has the ([A-Z][a-z]+) condition', body)]
    m = re.search(r'pushed up to (\d+) feet', body)
    if m:
        detail['push_feet'] = int(m[1])
    m = re.search(r'(?:for|After) (\d+) minutes?', body)
    if m:
        detail['duration_minutes'] = int(m[1])
    m = re.search(r'subtracts (\d+) \((\d+)d(\d+)\) from its damage rolls', body)
    if m:
        detail['penalties'].append(dict(applies_to='damage-roll', flat=0, average=int(m[1]), dice=int(m[2]), sides=int(m[3])))
    m = re.search(r'rolls (\d+)d(\d+) whenever it makes an ability check or attack roll', body)
    if m:
        detail['penalties'].append(dict(applies_to='d20-test', flat=0, average=0, dice=int(m[1]), sides=int(m[2])))
    m = re.search(r'takes a −(\d+) penalty to AC', body)
    if m:
        detail['penalties'].append(dict(applies_to='armor-class', flat=int(m[1]), average=0, dice=0, sides=0))
    detail['once_per_turn_cycle'] = 'can’t take this action again until the start of its next turn' in body
    if 'moves up to half its Speed' in body:
        detail['movement'] = 'half-speed'
    if 'fly up to half its Fly Speed' in body:
        detail['movement'] = 'half-fly-speed'
    for m in re.finditer(r'(?:makes|make) (one|two|three) ([A-Z][a-z]+) attacks?', body):
        detail['attacks'].append(dict(action_id=slug(m[2]), count=WORDS[m[1]]))
    m = re.search(r'It can replace (one|two) attacks? with a use of ([^.]+)\.', body)
    if m:
        options = []
        for option in re.split(r' or ', re.sub(r'\([AB]\) ', '', m[2])):
            cast = re.fullmatch(r'Spellcasting to cast ([A-Z][A-Za-z ]+?)(?: \(level (\d) version\))?', option)
            options.append(dict(action_id='spellcasting', spell=cast[1], level=int(cast[2] or 0)) if cast
                           else dict(action_id=slug(option), spell='', level=0))
        detail['replacement'] = dict(attacks=WORDS[m[1]], options=options)
    if name != 'Multiattack':
        for m in re.finditer(r'(uses Spellcasting to cast|The dragon casts) ([A-Z][A-Za-z ]+?)(?: \(level (\d) version\))?(?= on itself|[.,])', body):
            detail['casts'].append(dict(spell=m[2], level=int(m[3] or 0), via_spellcasting=m[1].startswith('uses')))
    if name == 'Spellcasting' or 'The dragon casts Fear' in body:
        m = re.search(r'using Charisma as the spellcasting ability \(spell save DC (\d+)(?:, \+(\d+) to hit with spell attacks)?\)', body)
        detail['spellcasting'] = dict(ability='cha', save_dc=int(m[1]), attack_bonus=int(m[2]) if m[2] else None,
                                      spells=spell_list(body) if name == 'Spellcasting' else [])
    return record, detail


def parse(pdf, ident, name, pages, trait_names, action_names, legendary_names):
    lines = block_lines(pdf, name, pages)
    head = HEAD.match(lines[1])
    m = re.fullmatch(r'AC (\d+) Initiative ([+−]\d+) \((\d+)\)', lines[2])
    hp = re.fullmatch(DICE.replace(r'(\d+) \(', r'HP (\d+) \(', 1), lines[3])
    speeds = {}
    for part in lines[4].removeprefix('Speed ').split(', '):
        s = re.fullmatch(r'(?:(Burrow|Climb|Fly|Swim) )?(\d+) ft\.', part)
        speeds[(s[1] or 'walk').lower()] = int(s[2])
    assert lines[5] == 'MOD SAVE MOD SAVE MOD SAVE', lines[5]
    abilities = {a.lower(): dict(score=int(b), modifier=signed(c), save=signed(d))
                 for a, b, c, d in re.findall(r'(Str|Dex|Con|Int|Wis|Cha) (\d+) ([+−]\d+) ([+−]\d+)', lines[6] + ' ' + lines[7])}
    assert len(abilities) == 6 and re.sub(r'(Str|Dex|Con|Int|Wis|Cha) \d+ [+−]\d+ [+−]\d+', '', lines[6] + lines[7]).strip() == ''
    section = next(i for i, line in enumerate(lines) if line in ('Traits', 'Actions'))
    fields, label = {}, None
    for line in lines[8:section]:
        m2 = re.match(r'(Skills|Resistances|Vulnerabilities|Immunities|Gear|Senses|Languages|CR) (.*)', line)
        if m2:
            label = m2[1]
            fields[label] = [m2[2]]
        else:
            fields[label].append(line)
    fields = {k: join(v) for k, v in fields.items()}
    assert set(fields) <= {'Skills', 'Immunities', 'Senses', 'Languages', 'CR'}, fields
    skills = {slug(a): signed(b) for a, b in re.findall(r'([A-Z][a-z]+) ([+−]\d+)', fields.get('Skills', ''))}
    damage, _, conditions = fields.get('Immunities', '').partition('; ')
    senses = fields['Senses'].split('; ')
    special = {slug(a): int(b) for a, b in re.findall(r'([A-Z][a-z]+) (\d+) ft\.', senses[0])} if len(senses) == 2 else {}
    cr = re.fullmatch(r'(\S+) \(XP ([\d,]+)(?:, or ([\d,]+) in lair)?; PB \+(\d)\)', fields['CR'])
    parts, current = {}, None
    for line in lines[section:]:
        if line in ('Traits', 'Actions', 'Bonus Actions', 'Reactions', 'Legendary Actions'):
            current = parts.setdefault(line, [])
        else:
            current.append(line)
    assert set(parts) <= {'Traits', 'Actions', 'Legendary Actions'}, set(parts)
    traits = [dict(key='trait/' + slug(n), name=n, text=t)
              for n, t in zip(trait_names, split_entries(join(parts.get('Traits', [])), trait_names))]
    spells = spell_list(next(t for n, t in zip(action_names, split_entries(join(parts['Actions']), action_names))
                             if n == 'Spellcasting')) if 'Spellcasting' in action_names else []
    if any(n == 'Frightful Presence' for n in legendary_names) and not spells:
        spells = [dict(name='Fear')]
    actions, details = [], []
    for n, t in zip(action_names, split_entries(join(parts['Actions']), action_names)):
        a, d = action(n, t, 'action', spells)
        actions.append(a)
        details.append(d)
    legendary = None
    if legendary_names:
        text = join(parts['Legendary Actions'])
        intro = text[:text.index(legendary_names[0] + '. ')].strip()
        uses = re.match(r'Legendary Action Uses: (\d+) \((\d+) in Lair\)\.', intro)
        resist = re.search(r'Legendary Resistance \((\d)/Day, or (\d)/Day in Lair\)', traits[-1]['text'])
        legendary = dict(uses=int(uses[1]), lair_uses=int(uses[2]), resistance_uses=int(resist[1]),
                         resistance_lair_uses=int(resist[2]), text=intro)
        for n, t in zip(legendary_names, split_entries(text[len(intro):].strip(), legendary_names)):
            a, d = action(n, t, 'legendary-action', spells)
            actions.append(a)
            details.append(d)
    required = ['monster/stat-block']
    required += ['speed/' + k for k in speeds if k != 'walk'] + ['skill/' + k for k in skills]
    required += ['immunity/%s-damage' % d.lower() for d in damage.split(', ') if d]
    required += ['immunity/%s-condition' % c.lower() for c in conditions.split(', ') if c]
    required += ['sense/' + k for k in special]
    if cr[3]:
        required.append('monster/lair-experience')
    if legendary:
        required += ['monster/legendary-actions', 'monster/lair-legendary-uses']
    for item in traits + actions:
        required += [item['key']] + [e['key'] for e in item.get('effects', [])]
    monster = dict(
        id='srd-5.2.1:content/' + ident, name=lines[0], source_pages=list(pages), sizes=[head[1].lower()], type=head[2].lower(),
        alignment=head[4], ac=int(m[1]), hp=dict(average=int(hp[1]), dice=int(hp[2]), sides=int(hp[3]), modifier=int(hp[4] or 0)),
        initiative=dict(modifier=signed(m[2]), score=int(m[3])), speeds=speeds, abilities=abilities,
        senses=dict(passive_perception=int(senses[-1].removeprefix('Passive Perception ')), special=special),
        languages=[] if fields['Languages'] == 'None' else fields['Languages'].split(', '), gear=[], cr=cr[1],
        xp=int(cr[2].replace(',', '')), proficiency_bonus=int(cr[4]), traits=traits, actions=actions,
        required_mechanics=list(dict.fromkeys(required)))
    return dict(monster=monster, tags=[head[3].lower()] if head[3] else [], skills=skills, damage_vulnerabilities=[],
                damage_resistances=[], damage_immunities=[d.lower() for d in damage.split(', ') if d],
                condition_immunities=[c.lower() for c in conditions.split(', ') if c],
                lair_xp=int(cr[3].replace(',', '')) if cr[3] else 0, legendary=legendary, action_details=details)


def expected(pdf):
    assert hashlib.sha256(pdf.read_bytes()).hexdigest() == PIN, 'PDF pin'
    source = json.loads((ROOT/'internal/rules/srd521/testdata/monsters.json').read_text())['source']
    assert source['sha256'] == PIN
    headings = []
    for ident, name, pages in HEADINGS:
        lines = column_lines(pdf, *pages)
        at = lines.index(name)
        assert not HEAD.match(lines[at + 1]) and HEAD.match(lines[at + 2]), 'heading must precede a named block'
        headings.append(dict(id='srd-5.2.1:content/' + ident, name=name, kind='section-heading', source_pages=list(pages)))
    return dict(schema_version=1, source=dict(source, changes=NOTICE), headings=headings,
                stat_blocks=[parse(pdf, *block) for block in BLOCKS])


def verify(pdf, candidate=None):
    want = expected(pdf)
    actual = json.loads((candidate or CATALOGUE).read_text())
    assert list(actual) == list(want), 'top-level members'
    for key in ('schema_version', 'source', 'headings'):
        assert actual[key] == want[key], key
    assert len(actual['stat_blocks']) == len(want['stat_blocks']) == 11
    for got, exp in zip(actual['stat_blocks'], want['stat_blocks']):
        for key in exp:
            if key != 'monster':
                assert got.get(key) == exp[key], (exp['monster']['name'], key)
        for key in exp['monster']:
            assert got['monster'].get(key) == exp['monster'][key], (exp['monster']['name'], key)
        assert got == exp, exp['monster']['name']
    inventory = json.loads((ROOT/'docs/investigations/issue-82-world-inventory.json').read_text())
    owned = {r['id']: (r['name'], r['source_page']) for r in inventory['records'] if r['catalogue_issue'] == 152}
    batch = next(b for b in inventory['batches'] if b['issue'] == 152)
    found = {b['monster']['id']: (b['monster']['name'], b['monster']['source_pages'][0]) for b in actual['stat_blocks']}
    found.update({h['id']: (h['name'], h['source_pages'][0]) for h in actual['headings']})
    assert found == owned and len(found) == batch['record_count'] == 12, 'inventory identities and pages'
    print('pp. 265–344: eleven stat blocks and the Animals heading match the pinned PDF and the #82 inventory')


if __name__ == '__main__':
    if len(sys.argv) == 3 and sys.argv[1] == '--expected':
        print(json.dumps(expected(Path(sys.argv[2])), ensure_ascii=False, indent=2))
    elif len(sys.argv) in (2, 3):
        verify(Path(sys.argv[1]), Path(sys.argv[2]) if len(sys.argv) == 3 else None)
    else:
        raise SystemExit('usage: verify_srd_statblocks_02.py PDF [CATALOGUE] | --expected PDF')
