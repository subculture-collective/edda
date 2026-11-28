import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { expect, it, vi } from 'vitest';
import { ActionPreviewPanel } from '../src/components/character/ActionPreviewPanel';
import { confirmCharacterAction, readActionReceipt } from '../src/api/actions';
import { getCampaignCharacterInventory, previewCharacterAction } from '../src/api/characters';
import type { ActionPreview, CharacterSheetProjection, Encounter } from '../src/api/types';
import fixture from './fixtures/srd-check-sheet.json';

vi.mock('../src/api/actions', () => ({ confirmCharacterAction: vi.fn(), readActionReceipt: vi.fn().mockResolvedValue(null), publishConfirmedAction: vi.fn() }));
vi.mock('../src/api/characters', () => ({ previewCharacterAction: vi.fn(), getCampaignCharacterInventory: vi.fn().mockResolvedValue([]) }));
const sheet = fixture.sheet as CharacterSheetProjection;
const ready: ActionPreview = {
  intent: { actor_id: 'hero', kind: 'ability-check', ability: 'strength', skill: 'athletics' },
  rules_version: 'srd-5.2.1', expected_revision: 'actor-v1:fixture', available: true,
  modifier: sheet.skills.athletics, targets: [], options: ['ordinary'], costs: [],
  cost_status: 'requires-adjudication', source_ids: ['srd-5.2.1:content/d20-tests/ability-checks'],
};
function mount(campaignId = 'alpha', actorId = 'hero', initialPreview?: ActionPreview) {
  return render(<QueryClientProvider client={new QueryClient()}><ActionPreviewPanel key={`${campaignId}:${actorId}`} campaignId={campaignId} actorId={actorId} sheet={sheet} initialPreview={initialPreview} /></QueryClientProvider>);
}

it('uses a supplied text preview until a manual selection deliberately discards it', async () => {
  mount('provided', 'hero', ready);
  expect(await screen.findByText('Preview ready. No action has been performed.')).toBeVisible();
  expect(previewCharacterAction).not.toHaveBeenCalled();
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'saving-throw:constitution');
  expect(screen.queryByText('Preview ready. No action has been performed.')).not.toBeInTheDocument();
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, intent: { actor_id: 'hero', kind: 'saving-throw', ability: 'constitution' } });
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(previewCharacterAction).toHaveBeenLastCalledWith('provided', { actor_id: 'hero', kind: 'saving-throw', ability: 'constitution' });
});

it('keeps an existing stored prepared action ahead of a supplied text preview', async () => {
  const pending = { action_id: 'old-action', expected_revision: 'old-v1', intent: ready.intent };
  sessionStorage.setItem('edda:prepared-action:pending:hero', JSON.stringify(pending));
  mount('pending', 'hero', ready);
  expect(await screen.findByRole('button', { name: 'Retry same action' })).toBeVisible();
  expect(screen.queryByRole('button', { name: 'Confirm action' })).not.toBeInTheDocument();
  expect(sessionStorage.getItem('edda:prepared-action:pending:hero')).toContain('old-action');
});
it('shows exact reach and clarification before an unavailable attack can be confirmed', async () => {
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, available: false, reason: 'No legal target.', targets: [{id:'foe', label:'Practice foe', distance_feet:15, reach_feet:5, unavailable_reason:'Out of reach.'}] });
  mount('spatial');
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'skill:athletics');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(await screen.findByRole('list', { name: 'Target ranges' })).toHaveTextContent('Practice foe: 15 ft · reach 5 ft. Out of reach.');
  expect(screen.queryByRole('button', {name:'Confirm action'})).not.toBeInTheDocument();
  expect(confirmCharacterAction).not.toHaveBeenCalled();
});
it('previews a selected sheet skill without confirming or performing an action', async () => {
  vi.mocked(previewCharacterAction).mockResolvedValue(ready);
  mount();
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'skill:athletics');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(await screen.findByText('Preview ready. No action has been performed.')).toBeVisible();
  expect(previewCharacterAction).toHaveBeenCalledWith('alpha', { actor_id: 'hero', kind: 'ability-check', skill: 'athletics' });
  expect(screen.getByText('Modifier: 5')).toBeVisible();
  expect(screen.getByText(/No direct resource cost/)).toHaveTextContent('Action and time costs require adjudication.');
  expect(confirmCharacterAction).not.toHaveBeenCalled();
  expect(screen.getByRole('button', { name: 'Confirm action' })).toBeEnabled();
  await userEvent.click(screen.getByText('Inspect preview sources'));
  expect(screen.getByText(/actor-v1:fixture/)).toBeVisible();
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'saving-throw:constitution');
  expect(screen.queryByText('Preview ready. No action has been performed.')).not.toBeInTheDocument();
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(previewCharacterAction).toHaveBeenLastCalledWith('alpha', { actor_id: 'hero', kind: 'saving-throw', ability: 'constitution' });
});
it('shows safe unavailable reasons for sheet entries and allows retry after network failure', async () => {
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, available: false, modifier: undefined, reason: 'This sheet entry does not support action previews yet.' }).mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce(ready);
  mount();
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'entry:attack');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(await screen.findByText('This sheet entry does not support action previews yet.')).toBeVisible();
  expect(screen.queryByText('Modifier: 5')).not.toBeInTheDocument();
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'skill:athletics');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(await screen.findByRole('alert')).toHaveTextContent('Try the preview again');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(await screen.findByText('Preview ready. No action has been performed.')).toBeVisible();
});
it('blocks overlapping selections while a preview is pending', async () => {
  let finish!: (value: ActionPreview) => void;
  vi.mocked(previewCharacterAction).mockReturnValue(new Promise(resolve => { finish = resolve; }));
  mount();
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'skill:athletics');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(screen.getByLabelText('Character action')).toBeDisabled();
  expect(screen.getByRole('button', { name: 'Preparing preview…' })).toBeDisabled();
  finish(ready);
  expect(await screen.findByText('Preview ready. No action has been performed.')).toBeVisible();
});

it('confirms once and retries uncertain delivery with the same action identity', async () => {
  vi.mocked(previewCharacterAction).mockResolvedValue(ready);
  vi.mocked(confirmCharacterAction).mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce({ action_id: 'saved', turn_number: 1, narrative: 'Recorded check.', state_changes: [], resolution_events: [], combat_active: false });
  mount('retry-campaign');
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'skill:athletics');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  await userEvent.click(await screen.findByRole('button', { name: 'Confirm action' }));
  expect(await screen.findByRole('alert')).toHaveTextContent('ID is retained');
  const request = vi.mocked(confirmCharacterAction).mock.calls[0][1];
  expect(request.expected_revision).toBe(ready.expected_revision);
  expect(request.intent).toEqual(ready.intent);
  await userEvent.click(screen.getByRole('button', { name: 'Retry same action' }));
  expect(await screen.findByText('Action saved.')).toBeVisible();
  expect(vi.mocked(confirmCharacterAction).mock.calls[1][1]).toEqual(request);
});
it('retains a live confirmation identity when its receipt exists before a lost response', async () => {
  const campaign = 'receipt-before-response';
  const storageKey = `edda:prepared-action:${campaign}:hero`;
  let rejectConfirmation!: (error: Error) => void;
  vi.mocked(confirmCharacterAction).mockReturnValueOnce(new Promise((_, reject) => { rejectConfirmation = reject; }));
  vi.mocked(readActionReceipt).mockImplementation(async (_campaign, actionId) => ({ result: {
    action_id: actionId, turn_number: 2, narrative: 'Committed before response was lost.',
    state_changes: [], resolution_events: [], combat_active: false,
  } }));
  const view = mount(campaign, 'hero', ready);
  await userEvent.click(screen.getByRole('button', { name: 'Confirm action' }));
  const request = vi.mocked(confirmCharacterAction).mock.calls[0][1];
  expect(sessionStorage.getItem(storageKey)).toBe(JSON.stringify(request));
  expect(readActionReceipt).not.toHaveBeenCalled();
  await act(async () => { rejectConfirmation(new Error('response lost after commit')); });
  expect(await screen.findByRole('alert')).toHaveTextContent('ID is retained');
  expect(sessionStorage.getItem(storageKey)).toBe(JSON.stringify(request));
  view.unmount();
  mount(campaign);
  expect(await screen.findByText('Committed before response was lost.')).toBeVisible();
  expect(readActionReceipt).toHaveBeenCalledExactlyOnceWith(campaign, request.action_id);
  expect(confirmCharacterAction).toHaveBeenCalledTimes(1);
  expect(sessionStorage.getItem(storageKey)).toBeNull();
  vi.mocked(readActionReceipt).mockResolvedValue(null);
});

it('reconciles a saved receipt after remount instead of rolling again', async () => {
  const request = { action_id: 'retained', expected_revision: ready.expected_revision, intent: ready.intent };
  sessionStorage.setItem('edda:prepared-action:reload:hero', JSON.stringify(request));
  vi.mocked(readActionReceipt).mockResolvedValueOnce({ result: { action_id: 'retained', turn_number: 2, narrative: 'Saved before disconnect.', state_changes: [], resolution_events: [], combat_active: false } });
  mount('reload');
  expect(await screen.findByText('Action saved.')).toBeVisible();
  expect(readActionReceipt).toHaveBeenCalledWith('reload', 'retained');
  expect(confirmCharacterAction).not.toHaveBeenCalled();
});

it('selects an equipped longsword and known target, then confirms only on request', async () => {
 vi.mocked(getCampaignCharacterInventory).mockResolvedValueOnce([{ id:'sword',campaign_id:'martial',name:'Longsword',description:'',rarity:'common',item_type:'weapon',properties:{srd521_id:'longsword'},equipped:true,quantity:1 }]);
 vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready,intent:{actor_id:'hero',kind:'martial-attack',entry_id:'sword'},targets:[{id:'foe',label:'Practice foe'}],costs:[{resource_id:'attack-action',amount:1}],cost_status:'fixture-only' });
 vi.mocked(confirmCharacterAction).mockResolvedValueOnce({action_id:'saved',turn_number:2,narrative:'Longsword hit, 8 slashing.',state_changes:[],resolution_events:[],combat_active:false});
 mount('martial');
 await screen.findByRole('option',{name:'Longsword (one hand)'});
 await userEvent.selectOptions(screen.getByLabelText('Character action'),'martial-attack:sword');
 await userEvent.click(screen.getByRole('button',{name:'Preview action'}));
 expect(await screen.findByRole('button',{name:'Confirm action'})).toBeDisabled();
 await userEvent.selectOptions(screen.getByLabelText('Attack target'),'foe');
 expect(confirmCharacterAction).not.toHaveBeenCalled();
 await userEvent.click(screen.getByRole('button',{name:'Confirm action'}));
 expect(await screen.findByText('Action saved.')).toBeVisible();
 expect(vi.mocked(confirmCharacterAction).mock.calls[0][1]).toMatchObject({target_id:'foe',intent:{kind:'martial-attack',entry_id:'sword'}});
});

it('previews a potion quantity and effect and confirms self-use only', async () => {
 vi.mocked(getCampaignCharacterInventory).mockResolvedValueOnce([{ id:'potion',campaign_id:'healing',name:'Potion of Healing',description:'',rarity:'common',item_type:'consumable',properties:{srd521_id:'potion-of-healing'},equipped:false,quantity:2 }]);
 vi.mocked(previewCharacterAction).mockResolvedValueOnce({...ready,intent:{actor_id:'hero',kind:'consumable-use',entry_id:'potion'},modifier:undefined,targets:[{id:'hero',label:'Mira (self)'}],options:['Quantity: 2; regain 2d4 + 2 HP, capped at 12 (current 1); self-use only'],costs:[{resource_id:'item:potion',amount:1},{resource_id:'bonus-action',amount:1}],cost_status:'fixture-only'});
 vi.mocked(confirmCharacterAction).mockResolvedValueOnce({action_id:'saved',turn_number:2,narrative:'Mira regains 6 HP.',state_changes:[],resolution_events:[],combat_active:false});
 mount('healing');
 await screen.findByRole('option',{name:'Potion of Healing (quantity 2, select target)'});
 await userEvent.selectOptions(screen.getByLabelText('Character action'),'consumable-use:potion');
 await userEvent.click(screen.getByRole('button',{name:'Preview action'}));
 expect(await screen.findByText(/Roll options: Quantity: 2; regain 2d4/)).toBeVisible();
 expect(screen.getByRole('button',{name:'Confirm action'})).toBeDisabled();
 await userEvent.selectOptions(screen.getByLabelText('Potion target'),'hero');
 expect(confirmCharacterAction).not.toHaveBeenCalled();
 await userEvent.click(screen.getByRole('button',{name:'Confirm action'}));
 expect(await screen.findByText('Action saved.')).toBeVisible();
 expect(vi.mocked(confirmCharacterAction).mock.calls[0][1]).toMatchObject({target_id:'hero',intent:{kind:'consumable-use',entry_id:'potion'}});
});


it('requires an explicit encounter target and shows the durable budget', async () => {
 vi.mocked(getCampaignCharacterInventory).mockResolvedValueOnce([{ id:'sword',campaign_id:'turns',name:'Longsword',description:'',rarity:'common',item_type:'weapon',properties:{srd521_id:'longsword'},equipped:true,quantity:1 }]);
 vi.mocked(previewCharacterAction).mockResolvedValueOnce({...ready,intent:{actor_id:'hero',kind:'combat-start',entry_id:'sword'},targets:[{id:'foe',label:'Practice foe'}],costs:[],cost_status:'encounter'}).mockResolvedValueOnce({...ready,intent:{actor_id:'hero',kind:'combat-end-turn'},costs:[],cost_status:'encounter',encounter:{id:'encounter',status:'active',round:2,turn:0,actors:[{id:'hero',name:'Mira',player:true,hp:11,max_hp:12,initiative:20,initiative_die:19},{id:'foe',name:'Practice foe',player:false,hp:18,max_hp:20,initiative:19,initiative_die:19}],budget:{action:false,bonus_action:true,movement:15}}});
 vi.mocked(confirmCharacterAction).mockResolvedValueOnce({action_id:'saved',turn_number:2,narrative:'Initiative rolled.',state_changes:[],resolution_events:[],combat_active:false});
 mount('turns');
 await screen.findByRole('option',{name:'Start ordinary melee encounter (Longsword)'});
 await userEvent.selectOptions(screen.getByLabelText('Character action'),'combat-start:sword');
 await userEvent.click(screen.getByRole('button',{name:'Preview action'}));
 expect(await screen.findByRole('button',{name:'Confirm action'})).toBeDisabled();
 await userEvent.selectOptions(screen.getByLabelText('Attack target'),'foe');
 await userEvent.click(screen.getByRole('button',{name:'Confirm action'}));
 expect(await screen.findByText('Action saved.')).toBeVisible();
 expect(vi.mocked(confirmCharacterAction).mock.calls[0][1]).toMatchObject({target_id:'foe',intent:{kind:'combat-start',entry_id:'sword'}});
 await userEvent.selectOptions(screen.getByLabelText('Character action'),'combat-end-turn');
 await userEvent.click(screen.getByRole('button',{name:'Preview action'}));
 expect(await screen.findByText("Round 2 · Mira's turn")).toBeVisible();
 expect(screen.getByRole('list',{name:'Initiative order'})).toHaveTextContent('Practice foe: initiative 19');
 expect(screen.getByText('Action: spent · Bonus action: available · Movement: 15 ft')).toBeVisible();
});

it('requires an explicit casting slot and self target before confirmation', async () => {
 const caster = {...sheet, spellcasting: { ability:'wisdom', modifier:sheet.skills.athletics, spells:[{id:'cure-wounds',label:'Cure Wounds',known:true,prepared:true,components:['verbal','somatic'],components_available:true,slot_levels:[1],source_id:'spell'}] }};
 vi.mocked(previewCharacterAction).mockResolvedValueOnce({...ready,intent:{actor_id:'hero',kind:'spell-cast',entry_id:'cure-wounds',slot_level:1},targets:[{id:'hero',label:'Mira (self; touch)'}],costs:[{resource_id:'spell-slot-1',amount:1},{resource_id:'magic-action',amount:1}]});
 vi.mocked(confirmCharacterAction).mockResolvedValueOnce({action_id:'cast',turn_number:1,narrative:'Cure Wounds saved.',state_changes:[],combat_active:false,resolution_events:[]});
 render(<QueryClientProvider client={new QueryClient()}><ActionPreviewPanel campaignId="casting" actorId="hero" sheet={caster} /></QueryClientProvider>);
 await userEvent.selectOptions(screen.getByLabelText('Character action'),'spell-cast:cure-wounds');
 expect(screen.getByLabelText('Spell slot level')).toHaveValue('1');
 await userEvent.click(screen.getByRole('button',{name:'Preview action'}));
 expect(await screen.findByRole('button',{name:'Confirm action'})).toBeDisabled();
 expect(previewCharacterAction).toHaveBeenCalledWith('casting',{actor_id:'hero',kind:'spell-cast',entry_id:'cure-wounds',slot_level:1});
 await userEvent.selectOptions(screen.getByLabelText('Spell target'),'hero');
 expect(confirmCharacterAction).not.toHaveBeenCalled();
 await userEvent.click(screen.getByRole('button',{name:'Confirm action'}));
 expect(await screen.findByText('Action saved.')).toBeVisible();
 expect(vi.mocked(confirmCharacterAction).mock.calls[0][1]).toMatchObject({target_id:'hero',intent:{kind:'spell-cast',entry_id:'cure-wounds',slot_level:1}});
});

it('routes Tactical Mind decline through preview and requires explicit confirmation', async () => {
  const tactical = { ...ready, intent: { actor_id: 'hero', kind: 'tactical-mind', entry_id: 'decline' }, modifier: undefined, options: ['Close the saved failed-check window without rolling or spending.'] };
  vi.mocked(previewCharacterAction).mockResolvedValueOnce(tactical);
  render(<QueryClientProvider client={new QueryClient()}><ActionPreviewPanel campaignId="tactical" actorId="hero" sheet={{ ...sheet, scope: 'fighter-2-chain-mail-defense-fixture' }} /></QueryClientProvider>);
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'tactical-mind:decline');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(previewCharacterAction).toHaveBeenLastCalledWith('tactical', { actor_id: 'hero', kind: 'tactical-mind', entry_id: 'decline' });
  expect(await screen.findByText(/Close the saved failed-check window/)).toBeVisible();
  expect(confirmCharacterAction).not.toHaveBeenCalled();
  expect(screen.getByRole('button', { name: 'Confirm action' })).toBeEnabled();
});

it('routes the explicitly selected Easy Athletics DC 10 fixture without accepting a DC', async () => {
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, intent: { actor_id: 'hero', kind: 'ability-check', entry_id: 'easy-strength-athletics-dc10' }, options: ['Edda bounded adjudication fixture: Easy Strength (Athletics), DC 10.'] });
  render(<QueryClientProvider client={new QueryClient()}><ActionPreviewPanel campaignId="easy" actorId="hero" sheet={{ ...sheet, scope: 'fighter-2-chain-mail-defense-fixture' }} /></QueryClientProvider>);
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'ability-check-fixture:easy-strength-athletics-dc10');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(await screen.findByText(/Edda bounded adjudication fixture: Easy Strength/)).toBeVisible();
  expect(previewCharacterAction).toHaveBeenLastCalledWith('easy', { actor_id: 'hero', kind: 'ability-check', entry_id: 'easy-strength-athletics-dc10' });
  expect(confirmCharacterAction).not.toHaveBeenCalled();
});
it('previews the Magic Initiate Detect Magic free use without a slot level and offers study while concentrating', async () => {
  const detect = { id: 'detect-magic', label: 'Detect Magic (Magic Initiate free use)', known: true, prepared: true, playable: true, components: ['verbal', 'somatic'], components_available: true, slot_levels: [], source_id: 'srd-5.2.1:content/spells/detect-magic' };
  const cleric = { ...sheet, spellcasting: { ability: 'wisdom', modifier: sheet.skills.insight, spells: [detect] }, concentration: { id: 'cast', spell: 'detect-magic', started_at: 0, expires_at: 600 } } as CharacterSheetProjection;
  vi.mocked(previewCharacterAction).mockResolvedValue({ ...ready, intent: { actor_id: 'hero', kind: 'spell-cast', entry_id: 'detect-magic' } });
  render(<QueryClientProvider client={new QueryClient()}><ActionPreviewPanel campaignId="detect" actorId="hero" sheet={cleric} /></QueryClientProvider>);
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'spell-cast:detect-magic');
  expect(screen.queryByLabelText('Spell slot level')).not.toBeInTheDocument();
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(previewCharacterAction).toHaveBeenLastCalledWith('detect', { actor_id: 'hero', kind: 'spell-cast', entry_id: 'detect-magic' });
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'detect-magic-study');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(previewCharacterAction).toHaveBeenLastCalledWith('detect', { actor_id: 'hero', kind: 'detect-magic-study', ability: undefined });
});

it('keeps a reaction chosen from a preview that offered it', async () => {
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, intent: { actor_id: 'hero', kind: 'combat-end-turn' }, encounter: { id: 'encounter', status: 'active', round: 1, turn: 1, actors: [{ id: 'hero', name: 'Mira', player: true, hp: 11, max_hp: 12, initiative: 20, initiative_die: 19 }, { id: 'foe', name: 'Practice foe', player: false, hp: 20, max_hp: 20, initiative: 19, initiative_die: 19 }], budget: { action: false, bonus_action: true, movement: 30 }, reaction: { available: true, pending: { id: 'trigger', target_id: 'foe', weapon_id: 'sword', space_revision: 'scene' } } } });
  mount('reaction-from-preview');
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'combat-end-turn');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  await screen.findByRole('option', { name: 'Decline opportunity attack (no cost)' });
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'combat-decline-reaction');
  expect(screen.getByLabelText('Character action')).toHaveValue('combat-decline-reaction');
  expect(screen.getByRole('button', { name: 'Preview action' })).toBeEnabled();
});

it('keeps an unacknowledged request when its menu entry disappears, then clears the selection', async () => {
  const caster = (covered: boolean) => ({ ...sheet, light: { id: 'light', item_id: 'mace', color: 'blue', started_at: 0, expires_at: 3600, bright_feet: 20, dim_feet: 20, covered } }) as CharacterSheetProjection;
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, intent: { actor_id: 'hero', kind: 'light-cover' }, modifier: undefined });
  vi.mocked(confirmCharacterAction).mockRejectedValueOnce(new Error('offline')).mockResolvedValueOnce({ action_id: 'saved', turn_number: 1, narrative: 'Covered.', state_changes: [], resolution_events: [], combat_active: false });
  const client = new QueryClient();
  const panel = (covered: boolean) => <QueryClientProvider client={client}><ActionPreviewPanel campaignId="cover-retry" actorId="hero" sheet={caster(covered)} /></QueryClientProvider>;
  const view = render(panel(false));
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'light-cover');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  await userEvent.click(await screen.findByRole('button', { name: 'Confirm action' }));
  expect(await screen.findByRole('alert')).toHaveTextContent('ID is retained');
  // The refreshed sheet shows the commit before the client has its response.
  view.rerender(panel(true));
  expect(screen.getByText('Preview ready. No action has been performed.')).toBeVisible();
  await userEvent.click(screen.getByRole('button', { name: 'Retry same action' }));
  const calls = vi.mocked(confirmCharacterAction).mock.calls.slice(-2);
  expect(calls[1][1]).toEqual(calls[0][1]);
  expect(await screen.findByText('Covered.')).toBeVisible();
  expect(screen.getByLabelText('Character action')).toHaveValue('');
  expect(screen.getByRole('button', { name: 'Preview action' })).toBeDisabled();
});

it('lists saved receipts in the order the server sent them', async () => {
  const zero = { rules_version: 'srd-5.2.1', event: 'death-save', reason: 'Death Saving Throw after the attack.', hp_before: 0, hp_after: 0, max_hp: 12, before: {}, after: { successes: 1 }, source_ids: [] };
  const attack = { rules_version: 'srd-5.2.1', weapon: 'Unarmed Strike', target: 'Mira', outcome: 'hit', damage: 1, damage_type: 'bludgeoning', die: 12, total: 14, damage_dice: [], damage_modifier: 1, modifier: { value: 2, derivation: [] }, source_ids: [] };
  const effect = { rules_version: 'srd-5.2.1', reason: 'Concentration expired after the round.', time_before: 0, time_after: 6, hp_before: 0, hp_after: 0, slots_before: 0, slots_after: 0, hit_dice_before: 1, hit_dice_after: 1, source_ids: [] };
  vi.mocked(confirmCharacterAction).mockResolvedValueOnce({ action_id: 'saved', turn_number: 2, narrative: 'Player turn ended.', state_changes: [], combat_active: false, resolution_events: [
    { type: 'srd_attack', label: 'NPC Unarmed Strike', outcome: 'hit', details: { receipt: attack } },
    { type: 'srd_death_save', label: '0 HP', outcome: zero.reason, details: { receipt: zero } },
    { type: 'srd_effect', label: 'Caster effects and recovery', outcome: effect.reason, details: { receipt: effect } },
  ] });
  mount('ordered', 'hero', ready);
  await userEvent.click(screen.getByRole('button', { name: 'Confirm action' }));
  const saved = await screen.findByRole('region', { name: 'Saved action result' });
  expect(saved.textContent).toMatch(/Unarmed Strike against Mira.*Death Saving Throw after the attack\..*Concentration expired after the round\./);
});

it('shows a saved Detect Magic receipt in the result panel', async () => {
  const receipt = { rules_version: 'srd-5.2.1', kind: 'detect-magic-study', spell: 'detect-magic', use: '', magic_initiate_uses_before: 0, magic_initiate_uses_after: 0, slots_before: 1, slots_after: 1, time_before: 0, time_after: 0, presence: true, auras: [], reason: 'Detect Magic study; Magic action; no slot or feat use. No sensed magic is on a visible creature or object.', source_ids: [] };
  vi.mocked(confirmCharacterAction).mockResolvedValueOnce({ action_id: 'saved', turn_number: 4, narrative: 'Recorded study.', state_changes: [], combat_active: false, resolution_events: [{ type: 'srd_detect_magic', label: 'Detect Magic', outcome: receipt.reason, details: { receipt } }] });
  mount('detect-magic', 'hero', ready);
  await userEvent.click(screen.getByRole('button', { name: 'Confirm action' }));
  expect(await screen.findByRole('region', { name: 'Saved Detect Magic receipt' })).toHaveTextContent(receipt.reason);
});

it('drops the confirmed preview while its menu entry stays offered', async () => {
  vi.mocked(previewCharacterAction).mockResolvedValueOnce(ready);
  vi.mocked(confirmCharacterAction).mockResolvedValueOnce({ action_id: 'saved', turn_number: 1, narrative: 'Recorded check.', state_changes: [], resolution_events: [], combat_active: false });
  mount('stale-preview');
  await userEvent.selectOptions(screen.getByLabelText('Character action'), 'skill:athletics');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  await userEvent.click(await screen.findByRole('button', { name: 'Confirm action' }));
  expect(await screen.findByText('Action saved.')).toBeVisible();
  expect(screen.queryByText('Preview ready. No action has been performed.')).not.toBeInTheDocument();
  expect(screen.queryByText('Modifier: 5')).not.toBeInTheDocument();
  // The entry stays chosen, ready for another preview.
  expect(screen.getByLabelText('Character action')).toHaveValue('skill:athletics');
  expect(screen.getByRole('button', { name: 'Preview action' })).toBeEnabled();
  // Previewing it again replaces the saved result with a confirmable preview.
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, expected_revision: 'actor-v2:fixture' });
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  expect(await screen.findByRole('button', { name: 'Confirm action' })).toBeEnabled();
  expect(screen.queryByText('Action saved.')).not.toBeInTheDocument();
});

// The Character tab has no saved reaction of its own: only a preview offers it.
it.each(['supplied', 'requested'])('stops offering a reaction once it is confirmed from a %s preview', async source => {
  const encounter: Encounter = { id: 'encounter', status: 'active', round: 1, turn: 0, actors: [{ id: 'foe', name: 'Practice foe', player: false, initiative: 19, initiative_die: 19 }], budget: { action: false, bonus_action: true, movement: 30 }, reaction: { available: true, pending: { id: 'trigger', target_id: 'foe', weapon_id: 'sword', space_revision: 'scene' } } };
  const offered: ActionPreview = { ...ready, intent: { actor_id: 'hero', kind: 'combat-end-turn' }, available: false, reason: 'Choose or decline the pending reaction before continuing.', modifier: undefined, cost_status: 'encounter', encounter };
  vi.mocked(previewCharacterAction).mockReset();
  if (source === 'requested') vi.mocked(previewCharacterAction).mockResolvedValueOnce(offered);
  vi.mocked(previewCharacterAction).mockResolvedValueOnce({ ...ready, intent: { actor_id: 'hero', kind: 'combat-decline-reaction' }, modifier: undefined, cost_status: 'encounter', encounter });
  vi.mocked(confirmCharacterAction).mockResolvedValueOnce({ action_id: 'saved', turn_number: 3, narrative: 'Opportunity attack declined.', state_changes: [], resolution_events: [], combat_active: false });
  mount(`reaction-${source}`, 'hero', source === 'supplied' ? offered : undefined);
  const menu = screen.getByLabelText('Character action');
  if (source === 'requested') {
    await userEvent.selectOptions(menu, 'combat-end-turn');
    await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  }
  await screen.findByRole('option', { name: 'Decline opportunity attack (no cost)' });
  await userEvent.selectOptions(menu, 'combat-decline-reaction');
  await userEvent.click(screen.getByRole('button', { name: 'Preview action' }));
  await userEvent.click(await screen.findByRole('button', { name: 'Confirm action' }));
  expect(await screen.findByText('Action saved.')).toBeVisible();
  expect(screen.queryByRole('option', { name: 'Decline opportunity attack (no cost)' })).not.toBeInTheDocument();
  expect(menu).toHaveValue('');
  expect(screen.getByRole('button', { name: 'Preview action' })).toBeDisabled();
});
