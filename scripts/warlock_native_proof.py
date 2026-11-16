"""Read-only actual-candidate/native proof; ROOT release is deliberately unset."""
import hashlib
import json
from pathlib import Path
import subprocess

from verify_warlock_scope import ROOT, exact, caps, prefix_check
from warlock_native_plan import build, projection, PROJECT

# Set only from ROOT's explicitly accepted POST-SORCERER release. The observed
# graph baseline below is NOT that release and can never stand in for it.
ACCEPTED_POST_SORCERER_PREFIXES = None
APPROVED_FINITE_EXTENSIONS = ()
BASELINE = ROOT / 'docs/investigations/issue-331-native-baseline.json'
FOREIGN = ROOT / 'docs/investigations/issue-331-sorcerer-boundary.json'
PROPOSAL = ROOT / 'docs/investigations/issue-331-native-proposal.json'


def validate_proposal(m, actual, baseline, foreign):
    expected = build(m, foreign, baseline['graph'])
    assert exact(actual, expected), 'complete actual candidate native bodies/parents/prerequisites/TaskSize'
    graph = {x['number']: x for x in baseline['graph']}
    assert len(graph) == 1017, 'observed baseline extent, not a final prefix'
    assert hashlib.sha256(graph[779]['body'].encode()).hexdigest() == '556d1618d1cd2bd9c956fe38861d89e94ee14ff61a505b5d2b6a1fa46d19d150'
    assert graph[779]['parent'] == 50 and graph[779]['blocked_by'] == [61, 64, 65, 66, 67, 488, 778]
    assert graph[968]['parent'] == 336 and graph[968]['blocked_by'] == [429]
    assert [x['number'] for x in m['source_prerequisites']] == [285, 296, 302, 303, 304, 318]
    for source in m['source_prerequisites']:
        old = graph[source['number']]
        assert source['state'] == old['state'] == 'CLOSED'
        assert source['body'] == old['body'] and source['title'] == old['title']
    for entry in m['spell_overlap_audit']:
        for owner in entry['candidate_owners']:
            old = graph[owner['issue']]
            assert owner['body'] == old['body'] and owner['title'] == old['title'], 'whole existing contract reuse'
    wizard = json.loads((ROOT / 'docs/investigations/issue-331-wizard-boundary.json').read_text())
    r = projection(actual, baseline['graph'], foreign, wizard)
    caps(r['graph'], r['parents'])
    assert not r['unresolved_foreign'], 'no empty/dangling foreign heads'
    assert all(parent in r['graph'] for parent in r['parents'].values()), 'complete containment graph'
    assert len(actual['tasks']) == 254 and len(actual['gates']) == 27, 'finite281 actual owned nodes'
    for old in baseline['graph']:
        if old['number'] not in (279, 280, 337):
            assert r['graph'][str(old['number'])] == [str(x) for x in old['blocked_by']], 'no unapproved old prerequisite mutation'
    assert r['children283'] == 95 and r['outgoing429'] == 100
    assert r['root279'] == 100 and r['root280'] == 95
    assert all(429 not in t['prerequisites'] for t in actual['tasks'] + actual['gates'])
    assert next(t for t in actual['tasks'] if t['key'] == 'hp')['prerequisites'] == [774, 968]
    return r


class Reader:
    """Bounded real GitHub queries with no repeated cursors or silent API hangs."""
    def __init__(self):
        self.calls = 0

    def __call__(self, query, variables=None):
        self.calls += 1
        assert self.calls <= 80, 'native read call cap'
        assert 'mutation' not in query, 'read-only proof'
        p = subprocess.run(['gh', 'api', 'graphql', '--input', '-'],
            input=json.dumps({'query': query, 'variables': variables or {}}),
            capture_output=True, text=True, check=True, timeout=45)
        result = json.loads(p.stdout)
        assert not result.get('errors'), 'GraphQL error'
        return result['data']['repository']


def read_live(issue_map):
    api, cursor, cursors, seen, graph = Reader(), None, set(), set(), []
    query = 'query($cursor:String){repository(owner:"subculture-collective",name:"edda"){issues(first:100,after:$cursor){pageInfo{hasNextPage endCursor}nodes{number title body state parent{number} blockedBy(first:100){pageInfo{hasNextPage}nodes{number}}}}}}'
    for _ in range(20):
        page = api(query, {'cursor': cursor})['issues']
        for row in page['nodes']:
            assert row['number'] not in seen, 'duplicate issue/page'
            seen.add(row['number'])
            assert not row['blockedBy']['pageInfo']['hasNextPage'], 'incoming cap'
            graph.append(dict(number=row['number'], title=row['title'], body=row['body'] or '',
                state=row['state'], parent=(row['parent'] or {}).get('number'),
                blocked_by=sorted(n['number'] for n in row['blockedBy']['nodes'])))
        if not page['pageInfo']['hasNextPage']:
            break
        cursor = page['pageInfo']['endCursor']
        assert cursor and cursor not in cursors, 'missing/repeated cursor'
        cursors.add(cursor)
    else:
        raise AssertionError('native graph page cap')
    project = {}
    nums = sorted(set(issue_map.values()) | {779, 968})
    for start in range(0, len(nums), 25):
        batch = nums[start:start+25]
        fields = ' '.join(f'i{n}:issue(number:{n})' + '{issueType{name} projectItems(first:100){pageInfo{hasNextPage}nodes{project{id} fieldValueByName(name:"Size"){... on ProjectV2ItemFieldSingleSelectValue{name}}}}}' for n in batch)
        result = api('{repository(owner:"subculture-collective",name:"edda"){' + fields + '}}')
        for n in batch:
            row = result['i' + str(n)]
            assert not row['projectItems']['pageInfo']['hasNextPage'], 'project page bound'
            entries = [x for x in row['projectItems']['nodes'] if x['project']['id'] == PROJECT]
            assert len(entries) == 1, 'exact Project membership'
            project[str(n)] = dict(type=row['issueType']['name'], size=entries[0]['fieldValueByName']['name'])
    return dict(graph=graph, project=project, api_calls=api.calls)


def validate_native(m, plan, receipt, snapshot):
    assert ACCEPTED_POST_SORCERER_PREFIXES is not None, 'ROOT accepted POST-SORCERER prefix release pending'
    assert set(ACCEPTED_POST_SORCERER_PREFIXES) == {47, 54, 283, 337}, 'complete ROOT release'
    assert receipt and receipt.get('issue_map'), 'actual native publication missing'
    ids = receipt['issue_map']
    rows = {t['key']: t for t in plan['tasks'] + plan['gates']}
    assert ids.keys() == rows.keys() and len(set(ids.values())) == len(ids)
    assert all(type(v) is int and v > 0 for v in ids.values()), 'typed actual native identities'
    foreign = receipt['foreign_issue_map']
    assert all(type(v) is int and v > 0 for v in foreign.values()), 'typed actual foreign identities'
    def resolve(value):
        if isinstance(value, int) or str(value).isdigit():
            return int(value)
        return foreign[value] if ':' in value else ids[value]
    graph = {t['number']: t for t in snapshot['graph']}
    assert len(graph) == len(snapshot['graph'])
    caps({n: t['blocked_by'] for n, t in graph.items()}, {n: t['parent'] for n, t in graph.items() if t['parent']})
    assert snapshot['project']['779'] == dict(type='Task', size='M')
    assert snapshot['project']['968'] == dict(type='Task', size='M')
    for key, expected in rows.items():
        actual = graph[ids[key]]
        assert actual['title'] == expected['title'] and actual['body'] == expected.get('native_body', expected['body']), ('exact native body/title', key)
        assert actual['parent'] == resolve(expected['parent'])
        required = set(expected['prerequisites']) | {resolve(d) for d in expected['depends']}
        extensions = {resolve(d) for owner, d in APPROVED_FINITE_EXTENSIONS if owner == key}
        assert set(actual['blocked_by']) == required | extensions, ('exact native prerequisites', key)
        assert snapshot['project'][str(ids[key])] == dict(type='Task', size='M')
        assert actual['state'] == 'OPEN', 'future implementation/qualification remains mandatory'
    sorcerer = json.loads(FOREIGN.read_text())
    wizard = json.loads((ROOT / 'docs/investigations/issue-331-wizard-boundary.json').read_text())
    foreign_rows = {'330:' + t['key']: t for t in sorcerer['tasks'] + sorcerer['shared_split_proposal'] + sorcerer['gates']}
    foreign_rows.update({'332:' + t['key']: t for t in wizard['tasks'] + wizard['gates']})
    def external(value, namespace):
        value = str(value)
        if value.isdigit():
            return int(value)
        if value == 'warlock-shared-output':
            return ids['shared-output']
        if value == 'wizard-shared-output':
            return foreign['332:shared-output']
        key = value if ':' in value else namespace + ':' + value
        return ids[key[4:]] if key.startswith('331:') else foreign[key]
    for key, number in foreign.items():
        expected = foreign_rows[key]
        actual = graph[number]
        assert actual['title'] == expected['title'] and actual['body'] == expected.get('native_body', expected['body']), ('actual foreign whole contract', key)
        assert actual['parent'] == external(expected['parent'], key[:3])
        required = set(expected['prerequisites']) | {external(d, key[:3]) for d in expected['depends']}
        extensions = {external(d, key[:3]) for owner, d in APPROVED_FINITE_EXTENSIONS if owner == key}
        assert set(actual['blocked_by']) == required | extensions, ('actual foreign prerequisites', key)
        assert snapshot['project'][str(number)] == dict(type='Task', size='M')
    baseline = json.loads(BASELINE.read_text())['graph']
    for old in baseline:
        actual = graph[old['number']]
        assert actual['title'] == old['title'] and actual['parent'] == old['parent'], 'old title/parent'
        if old['number'] in (279, 280, 337):
            assert set(old['blocked_by']) <= set(actual['blocked_by']), 'historical root edge removed'
        else:
            assert actual['blocked_by'] == old['blocked_by'], 'old owner prerequisites changed'
        if old['number'] not in (47, 54, 283, 337):
            assert actual['body'] == old['body'], 'old owner whole body changed'
    for n, prefix in ACCEPTED_POST_SORCERER_PREFIXES.items():
        append = (ROOT / f'docs/investigations/issue-331-append-{n}.md').read_text()
        assert len(append) <= (70 if n == 47 else 750)
        prefix_check(graph[n]['body'], prefix, append, '<!-- edda331 -->')
    assert ids['class-gate'] in graph[337]['blocked_by']
    aggregate = foreign['330:shared-aggregate']
    assert graph[aggregate]['parent'] == 191
    assert set(graph[aggregate]['blocked_by']) == {foreign['330:shared-output'], ids['shared-output'], foreign['332:shared-output']}
    assert aggregate in graph[279]['blocked_by'] and aggregate in graph[280]['blocked_by']
    assert ids['class-gate'] not in graph[aggregate]['blocked_by']
    # Decision closure itself may occur; all actual implementation edges remain.
    assert graph[331]['state'] in ('OPEN', 'CLOSED')
