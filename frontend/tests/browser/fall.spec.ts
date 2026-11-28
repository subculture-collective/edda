import { test, expect, signIn } from './fixtures';
import { character } from '../fixtures/data';
import fixture from '../fixtures/srd-check-sheet.json' with { type: 'json' };

for (const width of [1440, 390]) {
  test(`resolves a declared fall on solid ground and into water at ${width}px`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: 844 });
    let confirmations = 0;
    let saved: unknown = null;
    let fall: { id: string; height_feet: number; landing: string } | undefined = { id: 'f1', height_feet: 35, landing: 'solid' };
    const previews: Record<string, string>[] = [];
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: { ...character(), hp: 12, max_hp: 12, rules_version: 'srd-5.2.1', sheet: { ...fixture.sheet, fall } } }));
    await page.route('**/api/v1/campaigns/alpha/character/action-preview?*', route => {
      const params = Object.fromEntries(new URL(route.request().url()).searchParams);
      previews.push(params);
      const { actor_id, ...intent } = params;
      return route.fulfill({ json: { intent: { actor_id, ...intent }, rules_version: 'srd-5.2.1', expected_revision: 'fall-v1:fixture', available: true, targets: [], options: [`Prepared fall ${params.entry_id}.`], costs: params.skill ? [{ resource_id: 'reaction', amount: 1 }] : [], cost_status: 'requires-adjudication', source_ids: [] } });
    });
    await page.route('**/api/v1/campaigns/alpha/action?*', route => route.fulfill({ json: saved ? { result: saved } : null }));
    await page.route('**/api/v1/campaigns/alpha/action/confirm', route => {
      const request = route.request().postDataJSON();
      confirmations++;
      expect(Object.keys(request).sort()).toEqual(['action_id', 'expected_revision', 'intent']);
      const water = fall!.landing === 'water';
      const receipt = {
        rules_version: 'srd-5.2.1', height_feet: fall!.height_feet, landing: fall!.landing, damage_type: 'bludgeoning', source_ids: ['falling'],
        ...(water
          ? { dice_count: 2, dice: [6, 5], rolled: 11, halved: true, damage: 5, hp_before: 12, hp_after: 7, prone_before: true, prone_after: true, water_check: { rules_version: 'srd-5.2.1', kind: 'ability-check', ability: 'strength', skill: 'athletics', mode: 'ordinary', dice: [{ sides: 20, value: 10, kept: true }], modifier: { value: 5, derivation: [] }, total: 15, dc: 15, outcome: 'success', source_ids: ['check'] } }
          : { dice_count: 3, dice: [2, 5, 1], rolled: 8, halved: false, damage: 8, hp_before: 12, hp_after: 4, prone_before: false, prone_after: true }),
      };
      saved = { action_id: request.action_id, turn_number: confirmations + 2, narrative: `Saved fall ${request.intent.entry_id}.`, combat_active: false, state_changes: [], resolution_events: [{ type: 'srd_fall', label: 'Fall', outcome: `${receipt.damage} bludgeoning damage`, details: { receipt } }] };
      fall = undefined;
      if (!water) return route.abort('failed');
      return route.fulfill({ json: saved });
    });
    const openSheet = () => page.getByRole('tab', { name: 'Character', exact: true }).click();
    const status = page.getByLabel('Declared fall', { exact: true });
    const receipt = page.getByRole('tabpanel', { name: 'Character', exact: true }).getByLabel('Saved fall receipt');
    const actions = page.getByLabel('Character action');
    const confirm = () => page.getByRole('button', { name: 'Confirm action', exact: true }).click();

    await signIn(page); await page.goto('/play/alpha'); await openSheet();
    await expect(status).toContainText('A fall of 35 feet onto solid ground is waiting.');
    // Solid ground offers no water-entry check, and no height or damage control.
    await expect(actions.locator('optgroup[label="Declared fall"] option')).toHaveText(['Resolve the declared fall']);
    await actions.selectOption('fall-resolve');
    await page.getByRole('button', { name: 'Preview action', exact: true }).click();
    await expect(page.getByText('Declared fall: Prepared fall f1.')).toBeVisible();
    await expect(page.getByText('No direct resource cost.', { exact: false })).toBeVisible();
    expect(previews).toEqual([{ actor_id: 'hero', kind: 'fall-resolve', entry_id: 'f1' }]);
    expect(confirmations).toBe(0); expect(game.actions).toEqual([]);
    await confirm();
    await expect(page.getByRole('alert')).toContainText('Could not confirm delivery');
    // The reload recovers the saved result without confirming again.
    await page.reload(); await openSheet();
    await expect(status).toHaveCount(0);
    await expect(actions.locator('optgroup[label="Declared fall"]')).toHaveCount(0);
    expect(confirmations).toBe(1);

    fall = { id: 'f2', height_feet: 20, landing: 'water' };
    await page.reload(); await openSheet();
    await expect(status).toContainText('A fall of 20 feet onto water is waiting.');
    await expect(actions.locator('optgroup[label="Declared fall"] option')).toHaveText(['Resolve the declared fall (no Reaction)', 'Resolve it with a Reaction: Strength (Athletics)', 'Resolve it with a Reaction: Dexterity (Acrobatics)']);
    await actions.selectOption('fall-resolve:athletics');
    await page.getByRole('button', { name: 'Preview action', exact: true }).click();
    await expect(page.getByText('reaction: 1')).toBeVisible();
    expect(previews[1]).toEqual({ actor_id: 'hero', kind: 'fall-resolve', entry_id: 'f2', skill: 'athletics' });
    await confirm();
    await expect(receipt).toContainText('Fell 20 feet onto water.');
    await expect(receipt).toContainText('2d6: 6, 5 = 11, halved to 5.');
    await expect(receipt).toContainText('5 bludgeoning damage. HP 12 → 7. Condition: Prone → Prone.');
    await expect(receipt).toContainText('athletics check: 15 · success');
    await expect(status).toHaveCount(0);
    expect(confirmations).toBe(2);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
    // The narrative feed shows the same saved receipt.
    await page.getByRole('tab', { name: 'Play', exact: true }).click();
    const feed = page.getByRole('tabpanel', { name: 'Play', exact: true });
    await feed.getByText('Turn details', { exact: true }).last().click();
    await expect(feed.getByLabel('Saved fall receipt').filter({ visible: true })).toContainText('HP 12 → 7.');
  });
}
