import { test, expect, signIn } from './fixtures';
import { character } from '../fixtures/data';
import fixture from '../fixtures/srd-check-sheet.json' with { type: 'json' };

for (const width of [1440, 390]) {
  test(`casts, covers and dismisses Light at ${width}px`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: 844 });
    let confirmations = 0;
    let saved: unknown = null;
    let light: Record<string, unknown> | undefined;
    const spellcasting = { ability: 'wisdom', modifier: { value: 3, derivation: [] }, spells: [{ id: 'light', label: 'Light', known: true, prepared: true, playable: true, components: ['verbal', 'material (Holy Symbol focus)'], components_available: true, slot_levels: [], source_id: 'spell' }] };
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: { ...character(), hp: 10, max_hp: 10, rules_version: 'srd-5.2.1', sheet: { ...fixture.sheet, spellcasting, light, light_item: light ? 'Mace' : undefined, resources: [], capabilities: [] } } }));
    await page.route('**/api/v1/campaigns/alpha/character/action-preview?*', route => {
      const params = new URL(route.request().url()).searchParams;
      const kind = params.get('kind')!;
      const cast = kind === 'spell-cast';
      expect(params.get('slot_level')).toBeNull();
      expect(params.get('entry_id')).toBe(cast ? 'light:blue' : null);
      return route.fulfill({ json: { intent: { actor_id: 'hero', kind, ...(cast ? { entry_id: 'light:blue' } : {}) }, rules_version: 'srd-5.2.1', expected_revision: 'light-v1:fixture', available: true, targets: cast ? [{ id: 'torch', label: 'Torch (carried; touch)' }, { id: 'mace', label: 'Mace (carried; touch)' }] : [], options: [`Prepared ${kind}.`], costs: cast ? [{ resource_id: 'magic-action', amount: 1 }] : [], cost_status: 'fixture-only', source_ids: [] } });
    });
    await page.route('**/api/v1/campaigns/alpha/action?*', route => route.fulfill({ json: saved ? { result: saved } : null }));
    await page.route('**/api/v1/campaigns/alpha/action/confirm', route => {
      const request = route.request().postDataJSON();
      const kind = request.intent.kind;
      confirmations++;
      const before = light;
      if (kind === 'spell-cast') {
        expect(request.target_id).toBe('mace');
        expect(request.intent.entry_id).toBe('light:blue');
        light = { id: request.action_id, item_id: 'mace', color: 'blue', started_at: 28800, expires_at: 32400, bright_feet: 20, dim_feet: 20 };
      } else {
        expect(request.target_id).toBeUndefined();
        light = kind === 'light-dismiss' ? undefined : { ...light, covered: kind === 'light-cover' };
      }
      const reason = `Saved ${kind} on Mace.`;
      saved = { action_id: request.action_id, turn_number: confirmations + 2, narrative: reason, combat_active: false, state_changes: [], resolution_events: [{ type: 'srd_light', label: 'Light', outcome: kind, details: { receipt: { rules_version: 'srd-5.2.1', kind, reason, item: 'Mace', before, after: light, time: 28800, source_ids: ['spell'] } } }] };
      if (kind === 'spell-cast') return route.abort('failed');
      return route.fulfill({ json: saved });
    });
    const openSheet = () => page.getByRole('tab', { name: 'Character', exact: true }).click();
    const status = page.getByLabel('Active Light', { exact: true });
    const receipt = page.getByRole('tabpanel', { name: 'Character', exact: true }).getByLabel('Saved Light receipt');
    const preview = async (value: string) => {
      await page.getByLabel('Character action').selectOption(value);
      await page.getByRole('button', { name: 'Preview action', exact: true }).click();
    };
    const confirm = () => page.getByRole('button', { name: 'Confirm action', exact: true }).click();

    await signIn(page); await page.goto('/play/alpha'); await openSheet();
    await expect(status).toContainText('No Light active.');
    await page.getByLabel('Character action').selectOption('spell-cast:light');
    await expect(page.getByLabel('Spell slot level')).toHaveCount(0);
    await page.getByLabel('Light color').selectOption('blue');
    await page.getByRole('button', { name: 'Preview action', exact: true }).click();
    await expect(page.getByRole('button', { name: 'Confirm action', exact: true })).toBeDisabled();
    await page.getByLabel('Lit object').selectOption('mace');
    expect(confirmations).toBe(0); expect(game.actions).toEqual([]);
    await confirm();
    await expect(page.getByRole('alert')).toContainText('Could not confirm delivery');
    await page.reload(); await openSheet();
    await expect(status).toContainText('Light on Mace (blue)');
    await expect(status).toContainText('Expires Day 1, 09:00:00.');
    expect(confirmations).toBe(1);

    await preview('light-cover');
    await expect(page.getByText('Light: Prepared light-cover.')).toBeVisible();
    await confirm();
    await expect(status).toContainText('Covered: it sheds no light until uncovered.');
    await expect(receipt).toContainText('lasts until Day 1, 09:00:00; covered.');

    await preview('light-uncover');
    await confirm();
    await expect(status).not.toContainText('Covered:');

    await preview('light-dismiss');
    await confirm();
    await expect(status).toContainText('No Light active.');
    await expect(receipt).toContainText('No Light remains.');
    expect(confirmations).toBe(4);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  });
}
