import { TacticalReceipt, type TacticalMindReceipt } from './TacticalReceipt';
import { ActionSurgeReceipt } from './ActionSurgeReceipt';
import { EquipmentReceipt } from './EquipmentReceipt';
import { AdvancementReceipt } from './AdvancementReceipt';
import { useEffect, useRef, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { APIError } from '../../api/client';
import { confirmCharacterAction, readActionReceipt, publishConfirmedAction, type PreparedAction } from '../../api/actions';
import { EncounterSummary } from '../combat/EncounterSummary';
import { EffectReceipt } from './EffectReceipt';
import { SpellReceipt } from './SpellReceipt';
import { LightReceipt } from './LightReceipt';
import { FallReceipt } from './FallReceipt';
import { ZeroHPReceipt } from './ZeroHPReceipt';
import { HealingReceipt } from './HealingReceipt';
import { AttackReceipt } from './AttackReceipt';
import { CheckReceipt } from './CheckReceipt';
import type { CheckReceipt as Receipt, TurnResponse } from '../../api/types';
import { getCampaignCharacterInventory, previewCharacterAction } from '../../api/characters';
import type { ActionPreview, CharacterSheetProjection, Encounter } from '../../api/types';

// Matches srd521.LightColors; the server refuses any other color.
const lightColors = ['white', 'yellow', 'orange', 'red', 'green', 'blue', 'violet'];
const isLight = (intent: ActionPreview['intent']) => intent.kind.startsWith('light-') || !!intent.entry_id?.startsWith('light:');

export function ActionPreviewPanel({ campaignId, actorId, sheet, initialSelection = '', initialPreview, reaction }: {
  readonly campaignId: string;
  readonly actorId: string;
  readonly sheet: CharacterSheetProjection;
  readonly initialSelection?: string;
  readonly initialPreview?: ActionPreview;
  readonly reaction?: Encounter['reaction'];
}) {
  const client = useQueryClient();
  const inventory = useQuery({ queryKey: ['campaign', campaignId, 'inventory'], queryFn: () => getCampaignCharacterInventory(campaignId) });
  const [target, setTarget] = useState('');
  const [slotLevel, setSlotLevel] = useState('1');
  const [lightColor, setLightColor] = useState('white');
  const storageKey = `edda:prepared-action:${campaignId}:${actorId}`;
  const [pending, setPending] = useState<PreparedAction | null>(() => {
    try { return JSON.parse(sessionStorage.getItem(storageKey) ?? 'null') as PreparedAction | null; } catch { return null; }
  });
  const requestRef = useRef(pending);
  const restoredActionId = useRef(pending?.action_id);
  const [completed, setCompleted] = useState<TurnResponse | null>(null);
  const finish = (result: TurnResponse) => {
    setCompleted(result);
    setPending(null);

    try { sessionStorage.removeItem(storageKey); } catch { /* Storage is optional. */ }
    publishConfirmedAction(campaignId, result, requestRef.current?.intent.kind.startsWith('equipment-') ? 'Confirm equipment change' : requestRef.current?.intent.kind === 'action-surge' ? 'Confirm Action Surge' : requestRef.current?.intent.kind === 'advance-level' ? 'Confirm Fighter advancement' : requestRef.current && isLight(requestRef.current.intent) ? 'Confirm Light' : requestRef.current?.intent.kind === 'fall-resolve' ? 'Confirm Fall' : requestRef.current?.intent.kind === 'death-save' ? 'Confirm Death Saving Throw' : requestRef.current?.intent.kind === 'ally-aid' ? 'Confirm ally aid' : requestRef.current?.intent.kind === 'spell-cast' ? `Confirm ${requestRef.current.intent.entry_id === 'shield-of-faith' ? 'Shield of Faith' : requestRef.current.intent.entry_id === 'detect-magic' ? 'Detect Magic' : 'Cure Wounds'}` : requestRef.current?.intent.kind === 'detect-magic-study' ? 'Confirm Detect Magic study' : requestRef.current?.intent.kind === 'consumable-use' ? 'Confirm Potion of Healing' : requestRef.current?.intent.kind === 'combat-react' ? 'Confirm opportunity attack' : requestRef.current?.intent.kind === 'combat-decline-reaction' ? 'Decline opportunity attack' : requestRef.current?.intent.kind === 'martial-attack' ? 'Confirm Longsword attack' : `Confirm ${requestRef.current?.intent.skill ?? requestRef.current?.intent.ability ?? 'character'} check`);
    requestRef.current = null;
    void client.invalidateQueries({ queryKey: ['campaign', campaignId] });
  };
  useEffect(() => {
    if (!pending || pending.action_id !== restoredActionId.current) return;
    let active = true;
    void readActionReceipt(campaignId, pending.action_id).then(receipt => {
      if (active && receipt?.result) finish(receipt.result);
    }).catch(() => { /* The same action ID remains available for retry. */ });
    return () => { active = false; };
    // Reconcile requests restored on mount, never a live confirmation: its
    // receipt can exist before a lost response, while its retry ID is still needed.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [campaignId, pending?.action_id]);
  const confirmation = useMutation({
    mutationFn: (request: PreparedAction) => confirmCharacterAction(campaignId, request),
    onSuccess: finish,
    onError: error => {
      if (error instanceof APIError && error.status === 409 && error.message.includes('refresh')) {
        setPending(null); requestRef.current = null;
        try { sessionStorage.removeItem(storageKey); } catch { /* Storage is optional. */ }
        preview.reset();
        setProvidedPreview(undefined);
      }
    },
  });
  const [selection, setSelection] = useState(initialPreview?.intent.skill ? `skill:${initialPreview.intent.skill}` : initialSelection);
  const [providedPreview, setProvidedPreview] = useState(initialPreview);
  // Detect Magic's Magic Initiate free use has no slot level to choose.
  const slotLevels = (id?: string) => sheet.spellcasting?.spells.find(spell => spell.id === id)?.slot_levels ?? [];
  const preview = useMutation({ mutationFn: (value: string) => {
    const [kind, name] = value.split(':');
    // Light is a cantrip: its entry names a color and it spends no slot.
    if (value === 'spell-cast:light') return previewCharacterAction(campaignId, { actor_id: actorId, kind, entry_id: `light:${lightColor}` });
    // The server owns the declared fall: the client names it and, for water, the entry check.
    if (kind === 'fall-resolve') return previewCharacterAction(campaignId, { actor_id: actorId, kind, entry_id: sheet.fall?.id, ...(name ? { skill: name } : {}) });
    return previewCharacterAction(campaignId, {
      actor_id: actorId,
      ...(kind === 'spell-cast' && slotLevels(name).length > 0 ? { slot_level: Number(slotLevel) } : {}),
      kind: kind === 'skill' || kind === 'ability-check-fixture' ? 'ability-check' : kind,
      ...(kind === 'skill' ? { skill: name } : kind.startsWith('equipment-') || kind === 'entry' || kind === 'martial-attack' || kind === 'combat-start' || kind === 'combat-react' || kind === 'consumable-use' || kind === 'ally-aid' || kind === 'spell-cast' || kind === 'rest-short' || kind === 'advance-level' || kind === 'tactical-mind' || kind === 'ability-check-fixture' ? { entry_id: name } : { ability: name }),
    });
  } });
  const result = preview.data ?? (pending ? undefined : providedPreview);
  const trigger = (reaction ?? result?.encounter?.reaction)?.pending;
  return <section aria-label="Action preview" className="space-y-3 border-y border-gold/20 py-4">
    <h3 className="font-heading text-gold">Preview a character action</h3>
    <p className="text-sm text-pewter">Inspect an action without rolling dice or changing resources, time or history.</p>
    <label className="block text-sm">Character action
      <select value={selection} disabled={preview.isPending || !!pending} onChange={event => { setSelection(event.target.value); setProvidedPreview(undefined); setTarget(''); preview.reset(); confirmation.reset(); setCompleted(null); }} className="mt-2 block min-h-11 w-full max-w-full bg-midnight p-2 text-champagne">
        <option value="">Choose an action</option>
        <optgroup label="Ability checks">{Object.keys(sheet.abilities).map(name => <option key={name} value={`ability-check:${name}`}>{name} check</option>)}</optgroup>
        <optgroup label="Skills">{Object.keys(sheet.skills).map(name => <option key={name} value={`skill:${name}`}>{name.replaceAll('-', ' ')} check</option>)}</optgroup>
        <optgroup label="Saving throws">{Object.keys(sheet.saves).map(name => <option key={name} value={`saving-throw:${name}`}>{name} save</option>)}</optgroup>
        {!sheet.spellcasting && !sheet.scope.startsWith('fighter-2-') && <optgroup label="Advancement"><option value="advance-level:fixed">Fighter level 2 · fixed HP · retain Defense and mastery</option><option value="advance-level:roll">Fighter level 2 · roll HP · retain Defense and mastery</option></optgroup>}
        {sheet.scope.startsWith('fighter-2-') && <optgroup label="Tactical Mind"><option value="ability-check-fixture:easy-strength-athletics-dc10">Edda fixture · Easy Strength (Athletics) · DC 10</option><option value="tactical-mind:use">Tactical Mind · inspect saved failed check</option><option value="tactical-mind:decline">Decline Tactical Mind · close saved check window</option></optgroup>}
        <optgroup label="Equipped weapons">{inventory.data?.filter(item => item.equipped && item.quantity > 0 && item.item_type === 'weapon' && item.properties.srd521_id === 'longsword').map(item => <option key={item.id} value={`martial-attack:${item.id}`}>{item.name} (one hand)</option>)}</optgroup>
        {trigger && <optgroup label="Pending reaction"><option value={`combat-react:${trigger.weapon_id}`}>Longsword opportunity attack (reaction)</option><option value="combat-decline-reaction">Decline opportunity attack (no cost)</option></optgroup>}
        {sheet.capabilities.some(c => c.id === 'action-surge' && c.playable) && <optgroup label="Fighter level 2"><option value="action-surge">Action Surge (one extra Action except Magic)</option><option value="rest-short">Short Rest (one hour; recover Fighter resources)</option><option value="rest-long">Long Rest (eight hours; recover Fighter resources)</option></optgroup>}
        <optgroup label="Encounter turns"><option value="combat-end-turn">End turn</option><option value="combat-move">Spend 5 feet of movement</option><option value="combat-dash">Dash (action)</option><option value="combat-approach">Approach target 5 feet</option><option value="combat-retreat">Retreat from target 5 feet</option><option value="combat-stand">Stand up from Prone (half speed)</option>{inventory.data?.filter(item => item.equipped && item.quantity > 0 && item.properties.srd521_id === 'longsword').map(item => <option key={item.id} value={`combat-start:${item.id}`}>Start ordinary melee encounter ({item.name})</option>)}</optgroup>
        <optgroup label="Equipment">{inventory.data?.filter(item => item.quantity === 1 && item.item_type === 'armor' && ['chain-mail', 'shield'].includes(String(item.properties.srd521_id))).map(item => <option key={item.id} value={`equipment-${item.equipped ? 'remove' : 'equip'}:${item.id}`}>{item.equipped ? 'Remove' : 'Equip'} {item.name}</option>)}</optgroup>
        <optgroup label="Consumables">{inventory.data?.filter(item => item.quantity > 0 && item.item_type === 'consumable' && item.properties.srd521_id === 'potion-of-healing').map(item => <option key={item.id} value={`consumable-use:${item.id}`}>{item.name} (quantity {item.quantity}, select target)</option>)}</optgroup>
        {sheet.spellcasting && <optgroup label="Spells">{sheet.spellcasting.spells.filter(spell => spell.playable !== false).map(spell => <option key={spell.id} value={`spell-cast:${spell.id}`}>{spell.label} ({spell.prepared ? 'prepared' : 'not prepared'})</option>)}</optgroup>}
        {sheet.light && <optgroup label="Active Light"><option value={sheet.light.covered ? 'light-uncover' : 'light-cover'}>{sheet.light.covered ? 'Uncover the lit object' : 'Cover the lit object with something opaque'}</option><option value="light-dismiss">Dismiss Light (no action)</option></optgroup>}
        {sheet.spellcasting && <optgroup label="Concentration and recovery"><option value="concentration-end">End concentration (no action)</option>{sheet.concentration?.spell === 'detect-magic' && <option value="detect-magic-study">Study with Detect Magic (Magic action)</option>}<option value="wait-minute">Wait one quiet minute</option><option value="rest-short">Short Rest (one hour; no Hit Die)</option><option value="rest-short:hit-die">Short Rest (one hour; spend one d8 Hit Die)</option><option value="rest-long">Long Rest (eight hours; includes sleep)</option><option value="combat-start">Start ordinary encounter (Cleric)</option></optgroup>}
        {sheet.death_saves && !sheet.death_saves.dead && <optgroup label="At 0 HP"><option value="death-save">{sheet.death_saves.stable ? 'Wait to recover from Stable (1d4 hours)' : 'Roll a Death Saving Throw (outside encounters)'}</option>{!sheet.death_saves.stable && <option value="ally-aid:medicine">An ally tries to stabilize you (Medicine, DC 10)</option>}{inventory.data?.filter(item => item.quantity > 0 && item.item_type === 'consumable' && item.properties.srd521_id === 'potion-of-healing').map(item => <option key={item.id} value={`ally-aid:${item.id}`}>An ally gives you {item.name} (quantity {item.quantity})</option>)}</optgroup>}
        {sheet.fall && <optgroup label="Declared fall"><option value="fall-resolve">Resolve the declared fall{sheet.fall.landing === 'water' && ' (no Reaction)'}</option>{sheet.fall.landing === 'water' && <><option value="fall-resolve:athletics">Resolve it with a Reaction: Strength (Athletics)</option><option value="fall-resolve:acrobatics">Resolve it with a Reaction: Dexterity (Acrobatics)</option></>}</optgroup>}
        <optgroup label="Sheet entries">{sheet.capabilities.map(entry => <option key={entry.id} value={`entry:${entry.id}`}>{entry.label}</option>)}</optgroup>
      </select>
    </label>
    {selection === 'spell-cast:light' && <label className="block text-sm">Light color<select value={lightColor} disabled={!!pending || preview.isPending} onChange={event => { setLightColor(event.target.value); setProvidedPreview(undefined); setTarget(''); preview.reset(); }} className="block min-h-11 w-full bg-midnight p-2">{lightColors.map(color => <option key={color} value={color}>{color}</option>)}</select></label>}
    {selection.startsWith('spell-cast:') && selection !== 'spell-cast:light' && slotLevels(selection.slice('spell-cast:'.length)).length > 0 && <label className="block text-sm">Spell slot level<select value={slotLevel} disabled={!!pending || preview.isPending} onChange={event => { setSlotLevel(event.target.value); setTarget(''); preview.reset(); }} className="block min-h-11 w-full bg-midnight p-2">{sheet.spellcasting?.spells.find(spell => `spell-cast:${spell.id}` === selection)?.slot_levels.map(level => <option key={level} value={level}>Level {level} (one slot)</option>)}</select></label>}
    <button type="button" className="min-h-11 border border-gold/40 px-4 text-gold disabled:opacity-50" disabled={!selection || preview.isPending || !!pending} onClick={() => preview.mutate(selection)}>{preview.isPending ? 'Preparing preview…' : 'Preview action'}</button>
    {preview.isError && <p role="alert" className="text-ruby">Could not prepare this action. Try the preview again.</p>}
    {pending && !result && <button type="button" disabled={confirmation.isPending} className="min-h-11 text-gold" onClick={() => confirmation.mutate(pending)}>Retry same action</button>}
    {confirmation.isError && <p role="alert">{confirmation.error instanceof APIError && confirmation.error.status === 409 ? confirmation.error.message : 'Could not confirm delivery. Retry the same action; its ID is retained.'}</p>}
    {completed && <section aria-label="Saved action result"><p>Action saved.</p><p>{completed.narrative}</p>{completed.resolution_events?.filter(event => event.type === 'srd_tactical_mind').map((event,index) => <TacticalReceipt key={index} receipt={event.details.receipt as TacticalMindReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_action_surge').map((event,index) => <ActionSurgeReceipt key={index} receipt={event.details.receipt as import('../../api/types').ActionSurgeReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_equipment').map((event,index) => <EquipmentReceipt key={index} receipt={event.details.receipt as import('../../api/types').EquipmentReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_advancement').map((event, index) => <AdvancementReceipt key={index} receipt={event.details.receipt as import('../../api/types').AdvancementReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_effect').map((event,index) => <EffectReceipt key={index} receipt={event.details.receipt as import('../../api/types').EffectReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_spell').map((event,index) => <SpellReceipt key={index} receipt={event.details.receipt as import('../../api/types').SpellReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_light').map((event, index) => <LightReceipt key={index} receipt={event.details.receipt as import('../../api/types').LightReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_fall').map((event, index) => <FallReceipt key={index} receipt={event.details.receipt as import('../../api/types').FallReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_death_save').map((event, index) => <ZeroHPReceipt key={index} receipt={event.details.receipt as import('../../api/types').ZeroHPReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_consumable').map((event, index) => <HealingReceipt key={index} receipt={event.details.receipt as import('../../api/types').HealingReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_attack').map((event, index) => <AttackReceipt key={index} receipt={event.details.receipt as import('../../api/types').AttackReceipt} />)}{completed.resolution_events?.filter(event => event.type === 'srd_check').map((event, index) => <CheckReceipt key={index} receipt={event.details.receipt as Receipt} />)}</section>}
    {result?.encounter && <EncounterSummary encounter={result.encounter} />}
    {result && <div role="status" className="space-y-2 text-sm">
      {result.targets.some(option => option.reach_feet !== undefined) && <ul aria-label="Target ranges">{result.targets.map(option => <li key={option.id}>{option.label}: {option.distance_feet === undefined ? 'distance needs clarification' : `${option.distance_feet} ft`} · reach {option.reach_feet} ft. {option.unavailable_reason}</li>)}</ul>}
      {result.available ? <>
        <p>Preview ready. No action has been performed.</p>
        {result.modifier && <p>Modifier: {result.modifier.value}</p>}
        {result.intent.kind === 'death-save' || result.intent.kind === 'ally-aid' ? <p>{result.options.join(' ')}</p> : result.intent.kind === 'fall-resolve' ? <p>Declared fall: {result.options.join(' ')}</p> : <p>{result.intent.kind.startsWith('equipment-') ? 'Equipment change: ' : isLight(result.intent) ? 'Light: ' : 'Roll options: '}{result.options.join(', ')}{!result.intent.kind.startsWith('equipment-') && !isLight(result.intent) && '. Circumstances require adjudication.'}</p>}
        {(result.intent.kind === 'consumable-use' || result.intent.kind === 'ally-aid' || result.intent.kind === 'spell-cast' || result.intent.kind === 'martial-attack' || result.intent.kind === 'combat-start') && <label className="block">{result.intent.kind === 'ally-aid' ? 'Helping ally' : result.intent.kind === 'consumable-use' ? 'Potion target' : isLight(result.intent) ? 'Lit object' : result.intent.kind === 'spell-cast' ? 'Spell target' : 'Attack target'}<select className="block min-h-11 w-full bg-midnight p-2" value={target} disabled={!!pending || confirmation.isPending} onChange={event => setTarget(event.target.value)}><option value="">{result.intent.kind === 'ally-aid' ? 'Choose an ally' : 'Choose a known target'}</option>{result.targets.map(option => <option key={option.id} value={option.id} disabled={!!option.unavailable_reason}>{option.label}{option.distance_feet !== undefined ? ` · ${option.distance_feet} ft / reach ${option.reach_feet} ft` : ''}{option.unavailable_reason ? ` · ${option.unavailable_reason}` : ''}</option>)}</select></label>}
        {result.intent.kind !== 'fall-resolve' && result.intent.kind !== 'death-save' && result.intent.kind !== 'ally-aid' && <p>Known targets: {result.targets.length ? result.targets.map(target => target.label).join(', ') : 'No target selected. Scene targets require adjudication.'}</p>}
        <p>{result.costs.length ? result.costs.map(cost => `${cost.resource_id}: ${cost.amount}`).join(', ') : isLight(result.intent) || result.intent.kind === 'fall-resolve' || result.intent.kind === 'death-save' ? 'No direct resource cost.' : 'No direct resource cost for this check.'} {result.cost_status === 'encounter' ? 'Remaining turn resources are enforced on confirmation.' : result.cost_status === 'fixture-only' ? 'Prepared fixture outside encounter turns.' : 'Action and time costs require adjudication.'}</p>
        <details><summary className="min-h-11 cursor-pointer py-2">Inspect preview sources</summary>
          <ul>{result.modifier?.derivation.map((term, index) => <li key={index}>{term.label}: {term.value}<small className="block break-all text-pewter">{term.source_id}</small></li>)}</ul>
          {result.source_ids.map(source => <p key={source} className="break-all text-pewter">{source}</p>)}
          <p className="break-all text-pewter">{result.rules_version} · {result.expected_revision}</p>
        </details>
      </> : <p>{result.reason}</p>}
      {result.available && !completed && <button type="button" disabled={confirmation.isPending || ((result.intent.kind === 'consumable-use' || result.intent.kind === 'ally-aid' || result.intent.kind === 'spell-cast' || result.intent.kind === 'martial-attack' || result.intent.kind === 'combat-start') && !target && !pending)} className="min-h-11 border border-gold px-4 text-gold disabled:opacity-50" onClick={() => {
        const request = requestRef.current ?? { action_id: crypto.randomUUID(), expected_revision: result.expected_revision, intent: result.intent, ...(result.intent.kind === 'saving-throw' ? { target_id: actorId } : result.intent.kind === 'combat-react' ? { target_id: trigger?.target_id } : result.intent.kind === 'consumable-use' || result.intent.kind === 'ally-aid' || result.intent.kind === 'spell-cast' || result.intent.kind === 'martial-attack' || result.intent.kind === 'combat-start' ? { target_id: target } : {}) };
        requestRef.current = request; setPending(request);
        try { sessionStorage.setItem(storageKey, JSON.stringify(request)); } catch { /* In-memory retry still keeps the ID. */ }
        confirmation.mutate(request);
      }}>{confirmation.isPending ? 'Confirming…' : pending ? 'Retry same action' : 'Confirm action'}</button>}
      <p className="text-pewter">Action Surge supports only the level-two Human Soldier Fighter on its active ordinary encounter turn: one additional Action except Magic, one use per Short or Long Rest. Cure Wounds supports self or one prepared willing, conscious adjacent ally with a level 1 slot outside encounters. Shield of Faith supports self-casting with a level 1 slot. Shield of Faith takes a Bonus Action and concentration; it replaces the current concentration effect. Light supports one object you carry outside encounters: a Magic action and no slot, one hour, until cast again or dismissed; covering it is your declaration. A declared fall resolves outside encounters onto solid ground or water; the server holds its height and rolls the damage. Detect Magic supports only the recorded Magic Initiate free use outside encounters: no slot, ten-minute concentration, and only adjudicated, player-known magic within 30 feet. Potion of Healing supports conscious self-use or one prepared willing, conscious ally within 5 feet outside encounters; consumes one owned potion and one Bonus Action. Checks require scene adjudication. Martial attacks support only the prepared one-handed longsword fixture; encounter turns are limited to the ordinary two-actor fixture; mastery and other reaction features remain unavailable. Opportunity attacks require a saved explicit trigger.</p>
    </div>}
  </section>;
}
