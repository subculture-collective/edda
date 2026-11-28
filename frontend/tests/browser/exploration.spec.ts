import type { Page } from '@playwright/test';
import { test, expect, signIn } from './fixtures';
import { character } from '../fixtures/data';
import fixture from '../fixtures/srd-check-sheet.json' with { type: 'json' };

// Each control gets a sheet projection, a recorded preview and one confirmation.
async function serve(page: Page, options: {
  sheet: () => object;
  preview: (params: Record<string, string>) => object;
  confirm: (request: { action_id: string; intent: Record<string, string> }) => { narrative: string; resolution_events: object[]; lost?: boolean };
}) {
  const previews: Record<string, string>[] = [];
  const requests: Record<string, unknown>[] = [];
  let saved: unknown = null;
  await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: { ...character(), hp: 12, max_hp: 12, rules_version: 'srd-5.2.1', sheet: { ...fixture.sheet, ...options.sheet() } } }));
  await page.route('**/api/v1/campaigns/alpha/character/action-preview?*', route => {
    const params = Object.fromEntries(new URL(route.request().url()).searchParams);
    previews.push(params);
    const { actor_id, ...intent } = params;
    return route.fulfill({ json: { intent: { actor_id, ...intent }, rules_version: 'srd-5.2.1', expected_revision: 'fixture', available: true, targets: [], options: [], costs: [], cost_status: 'fixture-only', source_ids: [], ...options.preview(params) } });
  });
  await page.route('**/api/v1/campaigns/alpha/action?*', route => route.fulfill({ json: saved ? { result: saved } : null }));
  await page.route('**/api/v1/campaigns/alpha/action/confirm', route => {
    const request = route.request().postDataJSON();
    requests.push(request);
    const { lost, ...outcome } = options.confirm(request);
    saved = { action_id: request.action_id, turn_number: requests.length + 2, combat_active: false, state_changes: [], ...outcome };
    return lost ? route.abort('failed') : route.fulfill({ json: saved });
  });
  return { previews, requests };
}

const openSheet = (page: Page) => page.getByRole('tab', { name: 'Character', exact: true }).click();
const previewAction = (page: Page) => page.getByRole('button', { name: 'Preview action', exact: true }).click();
const confirmButton = (page: Page) => page.getByRole('button', { name: 'Confirm action', exact: true });
const sheetPanel = (page: Page) => page.getByRole('tabpanel', { name: 'Character', exact: true });
const fits = async (page: Page) => expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
async function feedShows(page: Page, label: string, text: string) {
  await page.getByRole('tab', { name: 'Play', exact: true }).click();
  const feed = page.getByRole('tabpanel', { name: 'Play', exact: true });
  await feed.getByText('Turn details', { exact: true }).last().click();
  await expect(feed.getByLabel(label).filter({ visible: true })).toContainText(text);
}

for (const width of [1440, 390]) {
  test(`opens a known door at ${width}px`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: 844 });
    let state = 'closed';
    const { previews, requests } = await serve(page, {
      sheet: () => ({ doors: [{ id: 'd1', name: 'Oak door', state }, { id: 'd2', name: 'Cellar door', state: 'open' }, { id: 'd3', name: 'Iron door' }] }),
      preview: params => params.entry_id === 'd3'
        ? { available: false, reason: 'Locked doors are not supported.' }
        : { targets: [{ id: 'd1', label: 'Oak door', distance_feet: 5, reach_feet: 5 }], options: ['Open Oak door (closed, unlocked, 5 feet away).'] },
      confirm: () => {
        state = 'open';
        const receipt = { rules_version: 'srd-5.2.1', kind: 'door-open', door_id: 'd1', door: 'Oak door', before: 'closed', after: 'open', distance_feet: 5, reach_feet: 5, reason: 'Oak door is now open.', source_ids: ['door'] };
        return { narrative: 'Oak door is now open.', resolution_events: [{ type: 'srd_door', label: 'Open door', outcome: 'open', details: { receipt } }] };
      },
    });
    const actions = page.getByLabel('Character action');
    await signIn(page); await page.goto('/play/alpha'); await openSheet(page);
    await expect(page.getByLabel('Known doors', { exact: true })).toContainText('Cellar door · open');
    // An open door offers no control; a door the server cannot rule on asks the preview why.
    await expect(actions.locator('optgroup[label="Doors"] option')).toHaveText(['Open Oak door', 'Open Iron door']);
    await actions.selectOption('door-open:d3');
    await previewAction(page);
    await expect(page.getByText('Locked doors are not supported.')).toBeVisible();
    await expect(confirmButton(page)).toHaveCount(0);

    await actions.selectOption('door-open:d1');
    await previewAction(page);
    await expect(page.getByText('Open Oak door (closed, unlocked, 5 feet away).')).toBeVisible();
    await expect(page.getByText('No direct resource cost.', { exact: false })).toBeVisible();
    expect(previews.at(-1)).toEqual({ actor_id: 'hero', kind: 'door-open', entry_id: 'd1' });
    expect(requests).toEqual([]); expect(game.actions).toEqual([]);
    await confirmButton(page).click();
    // The client sends no target and no door fact.
    expect(Object.keys(requests[0]).sort()).toEqual(['action_id', 'expected_revision', 'intent']);
    await expect(sheetPanel(page).getByLabel('Saved door receipt')).toContainText('Oak door: closed → open.');
    await expect(actions.locator('optgroup[label="Doors"] option')).toHaveText(['Open Iron door']);
    await fits(page);
    await feedShows(page, 'Saved door receipt', 'Oak door: closed → open.');
  });

  test(`travels one segment at a chosen pace at ${width}px`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: 844 });
    const targets = [
      { id: 'wood', label: 'Wood: 1.00 miles (5280 feet), forest, no road; maximum pace normal; 20 min 00 s at a normal pace' },
      { id: 'peak', label: 'Peak: 2.00 miles (10560 feet), mountain, no road; maximum pace slow', unavailable_reason: 'mountain terrain allows at most a slow pace' },
    ];
    const { previews, requests } = await serve(page, {
      sheet: () => ({ travel_paces: ['slow', 'normal', 'fast'] }),
      preview: params => params.entry_id === 'fast'
        ? { available: false, reason: 'No known route can be travelled within one hour at a fast pace.', targets: targets.map(target => ({ ...target, unavailable_reason: 'too fast' })) }
        : { targets, options: ['Travel one segment of at most one hour at a normal pace.'] },
      confirm: () => {
        const receipt = { rules_version: 'srd-5.2.1', kind: 'travel-hour', reason: 'Travelled.', route_id: 'r1', from_id: 'gate', from: 'Gate', to_id: 'wood', to: 'Wood', plan: { pace: 'normal', maximum_pace: 'normal', distance_feet: 5280, terrain: 'forest', good_road: false, miles_per_hour: 3, seconds: 1200, effects: ['Disadvantage on Dexterity (Stealth) checks'] }, time_before: 3600, time_after: 4800, day_travel_before: 0, day_travel_after: 1200, expired: { id: 'c', spell: 'shield-of-faith', started_at: 3500, expires_at: 4100 }, source_ids: ['travel'] };
        return { narrative: 'Travelled to Wood.', resolution_events: [{ type: 'srd_travel', label: 'Travel', outcome: 'arrived at Wood', details: { receipt } }], lost: true };
      },
    });
    const actions = page.getByLabel('Character action');
    await signIn(page); await page.goto('/play/alpha'); await openSheet(page);
    await expect(actions.locator('optgroup[label="Travel (up to one hour)"] option')).toHaveText(['Travel at a Slow pace', 'Travel at a Normal pace', 'Travel at a Fast pace']);
    // An unavailable pace still lists every destination with its reason.
    await actions.selectOption('travel-hour:fast');
    await previewAction(page);
    await expect(page.getByText('No known route can be travelled within one hour at a fast pace.')).toBeVisible();
    await expect(page.getByLabel('Travel destinations').getByRole('listitem')).toHaveCount(2);
    await expect(confirmButton(page)).toHaveCount(0);
    await fits(page);

    await actions.selectOption('travel-hour:normal');
    await previewAction(page);
    await expect(page.getByLabel('Travel destinations')).toContainText('Unavailable: mountain terrain allows at most a slow pace');
    expect(previews.at(-1)).toEqual({ actor_id: 'hero', kind: 'travel-hour', entry_id: 'normal' });
    await expect(confirmButton(page)).toBeDisabled();
    const destination = page.getByRole('combobox', { name: /^Destination/ });
    await expect(destination.locator('option[value="peak"]')).toBeDisabled();
    await destination.selectOption('wood');
    expect(requests).toEqual([]); expect(game.actions).toEqual([]);
    await confirmButton(page).click();
    await expect(page.getByRole('alert')).toContainText('Could not confirm delivery');
    expect(requests[0]).toMatchObject({ target_id: 'wood', intent: { actor_id: 'hero', kind: 'travel-hour', entry_id: 'normal' } });
    expect(Object.keys(requests[0]).sort()).toEqual(['action_id', 'expected_revision', 'intent', 'target_id']);
    await fits(page);
    // The reload recovers the saved result without confirming again.
    await page.reload(); await openSheet(page);
    await expect(sheetPanel(page).getByText('Action saved.')).toBeVisible();
    const receipt = sheetPanel(page).getByLabel('Saved travel receipt');
    await expect(receipt).toContainText('Travelled from Gate to Wood at a normal pace.');
    await expect(receipt).toContainText('1.00 miles (5280 feet), forest terrain, no road: 20 min 00 s.');
    await expect(receipt).toContainText('Campaign time Day 1, 01:00:00 → Day 1, 01:20:00. Travel today 0 min 00 s → 20 min 00 s.');
    await expect(receipt).toContainText('Disadvantage on Dexterity (Stealth) checks.');
    await expect(receipt).toContainText('Shield of Faith expired on the way.');
    expect(requests).toHaveLength(1);
    await fits(page);
    await feedShows(page, 'Saved travel receipt', 'Travelled from Gate to Wood');
  });

  test(`stands up from a saved Prone condition at ${width}px`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: 844 });
    let prone = true;
    const { previews, requests } = await serve(page, {
      sheet: () => ({ prone }),
      preview: () => ({ expected_revision: 'stand-v1:fixture', options: ['Stand up from Prone. Outside an encounter no movement budget is tracked.'] }),
      confirm: () => {
        prone = false;
        return { narrative: 'Stood up from Prone outside an encounter.', resolution_events: [] };
      },
    });
    const actions = page.getByLabel('Character action');
    const stand = actions.locator('option[value="combat-stand"]');
    await signIn(page); await page.goto('/play/alpha'); await openSheet(page);
    await expect(page.getByLabel('Conditions', { exact: true })).toContainText('Prone.');
    await expect(stand).toHaveText(['Stand up from Prone']);
    await actions.selectOption('combat-stand');
    await previewAction(page);
    await expect(page.getByText('Outside an encounter no movement budget is tracked.')).toBeVisible();
    await expect(page.getByText('No direct resource cost.', { exact: false })).toBeVisible();
    expect(previews).toEqual([{ actor_id: 'hero', kind: 'combat-stand' }]);
    await confirmButton(page).click();
    await expect(sheetPanel(page).getByText('Stood up from Prone outside an encounter.')).toBeVisible();
    await expect(page.getByLabel('Conditions', { exact: true })).toHaveCount(0);
    expect(requests).toHaveLength(1); expect(game.actions).toEqual([]);
    // Without a saved Prone the only stand choice is the encounter turn one.
    await expect(stand).toHaveText(['Stand up from Prone (half speed)']);
    await fits(page);
  });
}
