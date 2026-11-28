import { test, expect, signIn } from './fixtures';
import { character } from '../fixtures/data';
import fixture from '../fixtures/srd-check-sheet.json' with { type: 'json' };

for (const width of [1440, 390]) {
  test(`shows the 0 HP state and rolls a Death Saving Throw at ${width}px`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: 844 });
    let confirmations = 0;
    let saved: unknown = null;
    let saves: Record<string, number | boolean> = { failures: 1 };
    const previews: Record<string, string>[] = [];
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: { ...character(), hp: 0, max_hp: 12, rules_version: 'srd-5.2.1', sheet: { ...fixture.sheet, hp: { value: 0, derivation: [] }, death_saves: saves } } }));
    await page.route('**/api/v1/campaigns/alpha/character/action-preview?*', route => {
      const params = Object.fromEntries(new URL(route.request().url()).searchParams);
      previews.push(params);
      const { actor_id, ...intent } = params;
      return route.fulfill({ json: { intent: { actor_id, ...intent }, rules_version: 'srd-5.2.1', expected_revision: 'death-save-v1:fixture', available: true, targets: [], options: ['Death Saving Throw: d20, no modifier. Now 0 successes, 1 failures.'], costs: [], cost_status: 'requires-adjudication', source_ids: [] } });
    });
    await page.route('**/api/v1/campaigns/alpha/action?*', route => route.fulfill({ json: saved ? { result: saved } : null }));
    await page.route('**/api/v1/campaigns/alpha/action/confirm', route => {
      const request = route.request().postDataJSON();
      confirmations++;
      expect(Object.keys(request).sort()).toEqual(['action_id', 'expected_revision', 'intent']);
      expect(request.intent).toEqual({ actor_id: 'hero', kind: 'death-save' });
      const receipt = { rules_version: 'srd-5.2.1', event: 'death-save', reason: 'Death Saving Throw: d20 14, success (1 successes, 1 failures).', die: 14, outcome: 'success', hp_before: 0, hp_after: 0, max_hp: 12, before: { failures: 1 }, after: { successes: 1, failures: 1 }, time_before: 0, time_after: 6, source_ids: ['death-saving-throws'] };
      saved = { action_id: request.action_id, turn_number: 3, narrative: receipt.reason, combat_active: false, state_changes: [], resolution_events: [{ type: 'srd_death_save', label: 'Death saving throw', outcome: 'success', details: { receipt } }] };
      saves = { successes: 1, failures: 1 };
      return route.fulfill({ json: saved });
    });
    const openSheet = () => page.getByRole('tab', { name: 'Character', exact: true }).click();
    const status = page.getByLabel('0 HP state', { exact: true });
    const actions = page.getByLabel('Character action');

    await signIn(page); await page.goto('/play/alpha'); await openSheet();
    await expect(status).toContainText('At 0 HP: Unconscious');
    await expect(status).toContainText('Death Saving Throws: 0 successes, 1 failures.');
    await actions.selectOption('death-save');
    await page.getByRole('button', { name: 'Preview action', exact: true }).click();
    await expect(page.getByText('Death Saving Throw: d20, no modifier.', { exact: false })).toBeVisible();
    await expect(page.getByText('Known targets', { exact: false })).toHaveCount(0);
    // The preview names only the kind; the server holds the counts and rolls.
    expect(previews).toEqual([{ actor_id: 'hero', kind: 'death-save' }]);
    expect(confirmations).toBe(0); expect(game.actions).toEqual([]);
    await page.getByRole('button', { name: 'Confirm action', exact: true }).click();
    const receipt = page.getByRole('tabpanel', { name: 'Character', exact: true }).getByLabel('Saved 0 HP receipt');
    await expect(receipt).toContainText('Death Saving Throw: d20 14, success');
    await expect(receipt).toContainText('Death Saving Throws: 1 successes, 1 failures.');
    expect(confirmations).toBe(1);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);

    // Stable offers recovery; a dead character has no 0 HP action.
    saves = { stable: true };
    await page.reload(); await openSheet();
    await expect(status).toContainText('Stable at 0 HP; still Unconscious.');
    await expect(actions.locator('optgroup[label="At 0 HP"] option')).toHaveText(['Wait to recover from Stable (1d4 hours)']);
    saves = { failures: 3, dead: true };
    await page.reload(); await openSheet();
    await expect(status.getByRole('heading')).toHaveText('Dead');
    await expect(actions.locator('optgroup[label="At 0 HP"]')).toHaveCount(0);
  });
}
