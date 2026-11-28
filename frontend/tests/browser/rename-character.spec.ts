import { test, expect, signIn } from './fixtures';
import { expectDocumentFits, expectReachable } from './layout';
import { character } from '../fixtures/data';

const unnamed = { ...character(), name: '', needs_name: true };

for (const width of [1440, 390]) {
  test(`naming an unnamed character on load refreshes the sheet at ${width}`, async ({ page, game }, testInfo) => {
    await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
    let saved = '';
    const writes: unknown[] = [];
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: saved ? { ...character(), name: saved } : unnamed }));
    await page.route('**/api/v1/campaigns/alpha/character/name', route => {
      expect(route.request().method()).toBe('PUT');
      const body = route.request().postDataJSON() as { name: string };
      writes.push(body);
      saved = body.name;
      return route.fulfill({ status: 204 });
    });
    await signIn(page); await page.goto('/play/alpha');

    const dialog = page.getByRole('dialog', { name: 'Name your character' });
    await expect(dialog).toBeVisible();
    const input = dialog.getByRole('textbox', { name: 'Character name' });
    await expect(input).toBeFocused();
    for (const control of [input, dialog.getByRole('button', { name: 'Save name' }), dialog.getByRole('button', { name: 'Not now' })]) {
      await expectReachable(control);
    }
    await expectDocumentFits(page);
    await input.fill('   '); await input.press('Enter');
    await expect(dialog.getByRole('alert')).toHaveText('character name is required');
    expect(writes).toEqual([]);
    await page.screenshot({ path: testInfo.outputPath('name-prompt.png') });

    await input.fill('  Iona  '); await input.press('Enter');
    await expect(dialog).toBeHidden();
    expect(writes).toEqual([{ character_id: 'hero', name: 'Iona' }]);
    await page.getByRole('tab', { name: 'Character' }).click();
    await expect(page.getByRole('heading', { name: 'Iona', exact: true })).toBeVisible();
    await expectDocumentFits(page);
    expect(game.actions).toEqual([]);
  });

  test(`an API without the rename route never blocks play at ${width}`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
    // An older API returns the blank name without needs_name.
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: { ...character(), name: '' } }));
    await signIn(page); await page.goto('/play/alpha');
    await expect.poll(() => game.requests).toContain('GET /campaigns/alpha/history');
    const input = page.getByLabel('Your move');
    await input.fill('Wait by the gate'); await input.press('Enter');
    await expect.poll(() => game.actions).toContain('Wait by the gate');
    await expect(page.getByRole('dialog')).toHaveCount(0);
  });

  test(`a missing rename route on save leaves a notice and play at ${width}`, async ({ page, game }) => {
    await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
    await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: unnamed }));
    await page.route('**/api/v1/campaigns/alpha/character/name', route => route.fulfill({ status: 404, body: '404 page not found' }));
    await signIn(page); await page.goto('/play/alpha');
    const dialog = page.getByRole('dialog', { name: 'Name your character' });
    await dialog.getByRole('textbox', { name: 'Character name' }).fill('Iona');
    await dialog.getByRole('button', { name: 'Save name' }).click();
    await expect(dialog).toBeHidden();
    await expect(page.getByText('The name was not saved.', { exact: false })).toBeVisible();
    await expectDocumentFits(page);
    await expectReachable(page.getByRole('button', { name: 'Send action' }));
    const input = page.getByLabel('Your move');
    await input.fill('Wait by the gate'); await input.press('Enter');
    await expect.poll(() => game.actions).toContain('Wait by the gate');
    await page.getByRole('button', { name: 'Dismiss notice' }).click();
    await expect(page.getByText('The name was not saved.', { exact: false })).toHaveCount(0);
  });
}

test('Escape skips naming and the next load asks again', async ({ page, game }) => {
  await page.route('**/api/v1/campaigns/alpha/character', route => route.fulfill({ json: unnamed }));
  await signIn(page); await page.goto('/play/alpha');
  const dialog = page.getByRole('dialog', { name: 'Name your character' });
  await expect(dialog).toBeVisible();
  await page.keyboard.press('Escape');
  await expect(dialog).toHaveCount(0);
  const input = page.getByLabel('Your move');
  await input.fill('Wait by the gate'); await input.press('Enter');
  await expect.poll(() => game.actions).toContain('Wait by the gate');
  await page.reload();
  await expect(dialog).toBeVisible();
});
