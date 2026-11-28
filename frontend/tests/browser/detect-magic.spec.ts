import { test, expect, signIn } from './fixtures';
import { character, history, timestamp } from '../fixtures/data';
import type { Concentration, ResolutionEvent, SessionLogEntry, TurnResponse } from '../../src/api/types';
import fixture from '../fixtures/srd-check-sheet.json' with { type: 'json' };

for (const width of [1440, 390]) {
  test(`casts Detect Magic with the free use, studies and keeps the receipts at ${width}px`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: 844 });
    let confirmations = 0;
    let uses = 1;
    let saved: TurnResponse | null = null;
    const turns: SessionLogEntry[] = [];
    const previews: Record<string, string>[] = [];
    let concentration: Concentration | undefined = { id: 'faith-id', spell: 'shield-of-faith', started_at: 28800, expires_at: 29400 };
    const spellcasting = () => ({ ability: 'wisdom', modifier: { value: 3, derivation: [] }, choices: { divine_order: 'protector', cantrips: [], prepared_spells: [], magic_initiate_cantrips: [], magic_initiate_spell: 'Detect Magic', magic_initiate_ability: 'wisdom', magic_initiate_uses: uses }, spells: [{ id: 'detect-magic', label: 'Detect Magic (Magic Initiate free use)', known: true, prepared: true, playable: true, components: ['verbal', 'somatic'], components_available: true, slot_levels: [], source_id: 'spell' }] });
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: { ...character(), hp: 10, max_hp: 10, rules_version: 'srd-5.2.1', sheet: { ...fixture.sheet, spellcasting: spellcasting(), concentration, resources: [], capabilities: [] } } }));
    await page.route('**/api/v1/campaigns/alpha/history*', route => route.fulfill({ json: { entries: [...history(), ...turns] } }));
    await page.route('**/api/v1/campaigns/alpha/character/action-preview?*', route => {
      const params = Object.fromEntries(new URL(route.request().url()).searchParams);
      previews.push(params);
      const { actor_id, ...intent } = params;
      const cast = params.kind === 'spell-cast';
      return route.fulfill({ json: { intent: { actor_id, ...intent }, rules_version: 'srd-5.2.1', expected_revision: 'detect-magic-v1:fixture', available: true, targets: cast ? [{ id: 'hero', label: 'Mira (self; 30-foot sense)' }] : [], options: [`Prepared ${params.kind}.`], costs: cast ? [{ resource_id: 'magic-initiate', amount: 1 }, { resource_id: 'magic-action', amount: 1 }] : [], cost_status: 'fixture-only', source_ids: [] } });
    });
    await page.route('**/api/v1/campaigns/alpha/action?*', route => route.fulfill({ json: saved ? { result: saved } : null }));
    await page.route('**/api/v1/campaigns/alpha/action/confirm', route => {
      const request = route.request().postDataJSON();
      const kind: string = request.intent.kind;
      confirmations++;
      const before = concentration;
      let event: ResolutionEvent;
      let reason: string;
      if (kind === 'concentration-end') {
        concentration = undefined;
        reason = 'Concentration ended voluntarily; no action or slot cost.';
        // The server names the spells on caster effect receipts.
        event = { type: 'srd_effect', label: 'Caster effects and recovery', outcome: reason, details: { receipt: { rules_version: 'srd-5.2.1', kind, reason, before: { ...before!, spell_label: 'Detect Magic' }, time_before: 28800, time_after: 28800, slots_before: 1, slots_after: 1, hp_before: 10, hp_after: 10, hit_dice_before: 1, hit_dice_after: 1, source_ids: ['concentration'] } } };
      } else {
        const cast = kind === 'spell-cast';
        if (cast) {
          expect(request.target_id).toBe('hero');
          uses = 0;
          concentration = { id: request.action_id, spell: 'detect-magic', started_at: 28800, expires_at: 29400 };
        } else expect(request.target_id).toBeUndefined();
        reason = cast ? 'Detect Magic cast with the Magic Initiate (Cleric) free use. Adjudicated magic is sensed within 30 feet.' : 'Detect Magic study; Magic action; no slot or feat use. Faint aura around Warded door (20 feet): abjuration spell.';
        event = { type: 'srd_detect_magic', label: 'Detect Magic', outcome: reason, details: { receipt: { rules_version: 'srd-5.2.1', kind, spell: 'detect-magic', use: cast ? 'magic-initiate-free' : '', magic_initiate_uses_before: cast ? 1 : 0, magic_initiate_uses_after: 0, slots_before: 1, slots_after: 1, before, after: concentration, time_before: 28800, time_after: 28800, presence: true, auras: [], reason, source_ids: ['detect'] } } };
      }
      saved = { action_id: request.action_id, turn_number: confirmations + 2, narrative: reason, combat_active: false, state_changes: [], resolution_events: [event] };
      turns.push({ turn_number: saved.turn_number!, player_input: `Confirm ${kind}`, input_type: 'action', llm_response: reason, resolution_events: [event], created_at: timestamp });
      if (kind === 'spell-cast') return route.abort('failed');
      return route.fulfill({ json: saved });
    });
    const openSheet = () => page.getByRole('tab', { name: 'Character', exact: true }).click();
    const concentrating = page.getByLabel('Concentration and recovery', { exact: true });
    const result = page.getByRole('tabpanel', { name: 'Character', exact: true }).getByLabel('Saved action result');
    const actions = page.getByLabel('Character action');
    const study = actions.locator('option[value="detect-magic-study"]');
    const confirm = () => page.getByRole('button', { name: 'Confirm action', exact: true }).click();
    const preview = async (value: string) => {
      await actions.selectOption(value);
      await page.getByRole('button', { name: 'Preview action', exact: true }).click();
      await expect(page.getByText(`Prepared ${value.split(':')[0]}.`)).toBeVisible();
    };

    await signIn(page); await page.goto('/play/alpha'); await openSheet();
    await expect(concentrating).toContainText('Shield of Faith');
    await expect(study).toHaveCount(0);
    await actions.selectOption('spell-cast:detect-magic');
    // The free use spends no slot, so there is no slot level to choose.
    await expect(page.getByLabel('Spell slot level')).toHaveCount(0);
    await page.getByRole('button', { name: 'Preview action', exact: true }).click();
    await expect(page.getByText('magic-initiate: 1', { exact: false })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Confirm action', exact: true })).toBeDisabled();
    await page.getByLabel('Spell target').selectOption('hero');
    expect(previews).toEqual([{ actor_id: 'hero', kind: 'spell-cast', entry_id: 'detect-magic' }]);
    expect(confirmations).toBe(0); expect(game.actions).toEqual([]);
    await confirm();
    await expect(page.getByRole('alert')).toContainText('Could not confirm delivery');
    // The reload recovers the saved cast without confirming again.
    await page.reload(); await openSheet();
    await expect(concentrating).toContainText('Detect Magic: sensing magic within 30 feet');
    await expect(result.getByLabel('Saved Detect Magic receipt')).toContainText('Magic Initiate uses: 1 → 0. Spell slots: 1 → 1.');
    expect(confirmations).toBe(1);

    await preview('detect-magic-study');
    expect(previews[1]).toEqual({ actor_id: 'hero', kind: 'detect-magic-study' });
    await confirm();
    await expect(result.getByLabel('Saved Detect Magic receipt')).toContainText('Faint aura around Warded door (20 feet): abjuration spell.');

    await preview('concentration-end');
    await confirm();
    await expect(result.getByLabel('Saved caster effect receipt')).toContainText('Previous effect: Detect Magic (expires at Day 1, 08:10:00).');
    await expect(concentrating).toContainText('No concentration effect active.');
    await expect(study).toHaveCount(0);
    expect(confirmations).toBe(3);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);

    // After a reload the turn history alone holds all three receipts.
    await page.reload();
    const feed = page.getByRole('tabpanel', { name: 'Play', exact: true });
    await page.getByRole('tab', { name: 'Play', exact: true }).click();
    const gm = feed.getByRole('article', { name: 'Game Master entry' });
    // Two opening entries, then the three saved turns.
    await expect(gm).toHaveCount(5);
    for (const index of [2, 3, 4]) await gm.nth(index).getByText('Turn details', { exact: true }).click();
    const receipts = feed.getByLabel('Saved Detect Magic receipt').filter({ visible: true });
    await expect(receipts).toHaveCount(2);
    await expect(receipts.nth(0)).toContainText('Magic Initiate uses: 1 → 0.');
    await expect(receipts.nth(1)).toContainText('Faint aura around Warded door (20 feet): abjuration spell.');
    await expect(feed.getByLabel('Saved caster effect receipt').filter({ visible: true })).toContainText('Previous effect: Detect Magic (expires at');
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  });
}
