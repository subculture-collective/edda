import { test, expect, signIn } from './fixtures';
import { campaign, character, result, timestamp } from '../fixtures/data';
import fixture from '../fixtures/srd-check-sheet.json' with { type: 'json' };
import type { WebSocketRoute } from '@playwright/test';

// This joins the declared partial SRD boundary with the retained legacy
// play/reconnect path. It intentionally uses the fixture server, so it cannot
// prove database state, a physical phone session, or SRD free-text execution.
for (const viewport of [{ width: 1440, height: 900 }, { width: 390, height: 844 }]) {
  test(`qualified partial slice and legacy reconnect remain distinct at ${viewport.width}px`, async ({ page, game }) => {
    void game;
    await page.setViewportSize(viewport);
    let slots = 2;
    let confirms = 0;
    let srdSocket: WebSocketRoute | undefined;
    let srdConnections = 0;
    const srdActions: { input: string; action_id: string }[] = [];
    let legacySocket: WebSocketRoute | undefined;
    let legacyConnections = 0;
    const legacyActions: string[] = [];
    const spellcasting = { ability: 'wisdom', modifier: { value: 3, derivation: [] }, spells: [{ id: 'cure-wounds', label: 'Cure Wounds', known: true, prepared: true, playable: true, components: ['verbal', 'somatic'], components_available: true, slot_levels: [1], source_id: 'srd-5.2.1:content/spells/cure-wounds' }] };
    await page.route('**/api/v1/campaigns/alpha', route => route.fulfill({ json: {
      ...campaign(), rules_mode: 'crunch', rules_version: 'srd-5.2.1',
    } }));
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: {
      ...character(), hp: 1, max_hp: 10, rules_version: 'srd-5.2.1',
      sheet: { ...fixture.sheet, scope: 'cleric-1-cure-wounds-self-fixture', spellcasting, hp: { value: 1, derivation: [] }, max_hp: { value: 10, derivation: [] }, resources: [{ id: 'spell-slot-1', label: 'Level 1 spell slots', current: slots, maximum: 2, source_id: 'srd-5.2.1:content/spells/cure-wounds' }], capabilities: [] },
    } }));
    await page.route('**/api/v1/campaigns/alpha/character/action-preview?*', route => {
      const params = new URL(route.request().url()).searchParams;
      if (params.get('kind') === 'combat-end-turn') {
        return route.fulfill({ json: {
          intent: { actor_id: 'hero', kind: 'combat-end-turn' }, rules_version: 'srd-5.2.1', available: false,
          targets: [], options: [], costs: [], cost_status: 'encounter', source_ids: [],
        } });
      }
      expect(params.get('kind')).toBe('spell-cast');
      expect(params.get('entry_id')).toBe('cure-wounds');
      return route.fulfill({ json: { intent: { actor_id: 'hero', kind: 'spell-cast', entry_id: 'cure-wounds', slot_level: 1 }, rules_version: 'srd-5.2.1', expected_revision: 'qualification-v1', available: true, targets: [{ id: 'hero', label: 'Mira (self; touch)' }], options: ['Self only; level 1 slot'], costs: [{ resource_id: 'spell-slot-1', amount: 1 }], cost_status: 'fixture-only', source_ids: ['srd-5.2.1:content/spells/cure-wounds'] } });
    });
    await page.route('**/api/v1/campaigns/alpha/action?*', route => route.fulfill({ json: null }));
    await page.route('**/api/v1/campaigns/alpha/action/confirm', route => {
      confirms++;
      expect(route.request().postDataJSON()).toMatchObject({ target_id: 'hero', intent: { kind: 'spell-cast', entry_id: 'cure-wounds', slot_level: 1 } });
      slots = 1;
      return route.fulfill({ json: { action_id: 'saved-cure-wounds', turn_number: 3, narrative: 'Mira casts Cure Wounds on self: regains 7 HP.', combat_active: false, state_changes: [], resolution_events: [{ type: 'srd_spell', label: 'Cure Wounds', outcome: 'healed', details: { receipt: { rules_version: 'srd-5.2.1', kind: 'spell-cast', spell: 'Cure Wounds', target: 'Mira', target_id: 'hero', slot_level: 1, slots_before: 2, slots_after: 1, dice: [2, 2], modifier: { value: 3, derivation: [] }, rolled_healing: 7, healing: 7, hp_before: 1, hp_after: 8, source_ids: ['srd-5.2.1:content/spells/cure-wounds'] } } }] } });
    });
    await page.routeWebSocket('**/api/v1/campaigns/alpha/ws', socket => {
      srdSocket = socket; srdConnections++;
      socket.onMessage(raw => {
        const action = JSON.parse(String(raw)).payload;
        srdActions.push(action);
        socket.send(JSON.stringify({ type: 'error', timestamp, payload: { error: 'an internal error occurred' } }));
      });
    });
    await page.routeWebSocket('**/api/v1/campaigns/beta/ws', socket => {
      legacySocket = socket; legacyConnections++;
      socket.onMessage(raw => {
        const action = JSON.parse(String(raw)).payload;
        legacyActions.push(action.input);
        socket.send(JSON.stringify({ type: 'result', timestamp, payload: result }));
      });
    });
    await signIn(page); await page.goto('/play/alpha');
    await page.getByRole('tab', { name: 'Character', exact: true }).click();
    await expect(page.getByLabel('Rules coverage')).toContainText('Status: partial');
    await expect(page.getByLabel('Rules coverage')).toContainText('344 character/advancement scopes, 335 world scopes and 753 spell/item source families');
    await page.getByLabel('Character action').selectOption('spell-cast:cure-wounds');
    await page.getByRole('button', { name: 'Preview action', exact: true }).click();
    await expect(page.getByRole('button', { name: 'Confirm action', exact: true })).toBeDisabled();
    await page.getByLabel('Spell target').selectOption('hero');
    expect(confirms).toBe(0);
    await page.getByRole('button', { name: 'Confirm action', exact: true }).click();
    await expect(page.getByText('Level 1 spell slots: 1 / 2')).toBeVisible();
    expect(confirms).toBe(1);
    await page.getByRole('tab', { name: 'Play', exact: true }).click();
    const srdInput = page.getByLabel('Your move');
    await srdInput.fill('Inspect the SRD gate'); await srdInput.press('Enter');
    await expect(page.locator('.game-composer-feedback').getByText('an internal error occurred', { exact: true })).toBeVisible();
    await srdSocket!.close({ code: 1012, reason: 'qualification interruption' });
    await expect.poll(() => srdConnections).toBe(2);
    await page.getByRole('button', { name: 'Retry previous action', exact: true }).click();
    expect(srdActions).toHaveLength(2);
    expect(srdActions[1]).toEqual(srdActions[0]);
    await expect(page.getByText('Inspect the SRD gate', { exact: true })).toHaveCount(1);
    expect(confirms).toBe(1);
    await page.goto('/play/beta');
    await expect.poll(() => legacyConnections).toBe(1);
    await legacySocket!.close({ code: 1012, reason: 'qualification interruption' });
    await expect.poll(() => legacyConnections).toBe(2);
    const input = page.getByLabel('Your move');
    await expect(input).toBeEnabled();
    await input.fill('Inspect the legacy gate');
    await input.press('Enter');
    await expect.poll(() => legacyActions).toEqual(['Inspect the legacy gate']);
    expect(confirms).toBe(1);
  });
}
