import { test, expect, signIn, user } from './fixtures';

test('creation hands the opening scene to live play and completes a turn', async ({ page, game }) => {
  await signIn(page);
  await page.goto('/new');
  await page.getByRole('button', { name: /Campaign interview/ }).click();
  await page.getByRole('button', { name: 'Continue', exact: true }).click();
  await page.getByLabel('Your campaign notes').fill('A hopeful frontier');
  await page.getByRole('button', { name: 'Send reply' }).click();
  await page.getByRole('button', { name: /Pure storytelling/ }).click();
  await page.getByRole('button', { name: 'Continue', exact: true }).click();
  await page.getByRole('button', { name: /Character interview/ }).click();
  await page.getByRole('button', { name: 'Continue', exact: true }).click();
  await page.getByLabel('Your character notes').fill('Mira the scout');
  await page.getByRole('button', { name: 'Send reply' }).click();
  await page.getByRole('button', { name: 'Begin adventure' }).click();
  await expect(page).toHaveURL(/play\/alpha$/);
  await expect(page.getByText('A lantern lights the road.', { exact: true })).toBeVisible();
  await expect.poll(() => game.sockets.length).toBe(1);
  await page.getByLabel('Your move').fill('Search the road');
  await page.getByLabel('Your move').press('Enter');
  await expect(page.getByText('You find a silver key.', { exact: true })).toBeVisible();
  await expect(page.getByLabel('Your move')).toBeEnabled();
  expect(game.actions).toEqual(['Search the road']);
  expect(game.requests).toContain('POST /campaigns/start/world');
});

test('reconnects and accepts a subsequent action without duplicating it', async ({ page, game }) => {
  await signIn(page); await page.goto('/play/alpha');
  await expect.poll(() => game.sockets.length).toBe(1);
  await game.sockets[0].close({ code: 1012, reason: 'fixture interruption' });
  await expect.poll(() => game.sockets.length).toBe(2);
  await page.getByLabel('Your move').fill('Look around');
  await page.getByLabel('Your move').press('Enter');
  await expect(page.getByText('You find a silver key.', { exact: true })).toBeVisible();
  expect(game.actions).toEqual(['Look around']);
});

test('switches campaigns through the list and replays only the selected history', async ({ page, game }) => {
  await signIn(page); await page.goto('/play/alpha');
  await expect(page.getByText('The alpha gate opens.', { exact: true })).toBeVisible();
  await page.getByRole('button', { name: /Fixture Player/ }).click();
  await page.getByRole('link', { name: /Back to campaigns/i }).click();
  await page.locator('a[href="/play/beta"]').click();
  await expect(page.getByText('The beta gate opens.', { exact: true })).toBeVisible();
  await expect(page.getByText('The alpha gate opens.', { exact: true })).toHaveCount(0);
  expect(game.sockets.at(-1)?.url()).toContain('/beta/ws');
  await page.goto('/replay/beta');
  await expect(page.getByText('The beta gate opens.', { exact: true })).toBeVisible();
  await expect(page.getByText('The beta keeper welcomes you.', { exact: true })).toHaveCount(0);
  await page.getByRole('button', { name: 'Next turn' }).click();
  await expect(page.getByText('The beta keeper welcomes you.', { exact: true })).toBeVisible();
  await page.getByRole('button', { name: 'Previous turn' }).click();
  await expect(page.getByText('The beta keeper welcomes you.', { exact: true })).toHaveCount(0);
});

test('registers, restores a session, signs out and rejects an expired session', async ({ page, game }) => {
  await page.goto('/play/alpha'); await expect(page).toHaveURL(/login$/);
  await page.goto('/register');
  await page.getByLabel('Name', { exact: true }).fill(user.name);
  await page.getByLabel('Email').fill(user.email);
  await page.getByLabel('Password', { exact: true }).fill('fixture-password');
  await page.getByRole('button', { name: /Create account/i }).click();
  await expect(page).toHaveURL('/');
  await page.reload();
  await expect(page.getByRole('heading', { name: 'Campaigns', exact: true })).toBeVisible();
  expect(game.requests).toContain('GET /auth/me');
  await page.getByRole('button', { name: /Fixture Player/ }).click();
  await page.getByRole('button', { name: 'Sign out' }).click();
  await expect(page).toHaveURL(/login$/);
  expect(await page.evaluate(() => localStorage.getItem('gm_token'))).toBeNull();
  await signIn(page);
  await page.route('**/auth/me', route => route.fulfill({ status: 401, json: { error: 'Expired' } }));
  await page.reload(); await expect(page).toHaveURL(/login$/);
  expect(await page.evaluate(() => localStorage.getItem('gm_token'))).toBeNull();
});

test('ends the local session at sign-out without waiting for the server', async ({ page, game }) => {
  await signIn(page);
  let release = () => {};
  const held = new Promise<void>(resolve => { release = resolve; });
  let authorization: string | undefined;
  await page.route('**/auth/logout', async route => {
    authorization = route.request().headers().authorization;
    await held;
    await route.fulfill({ status: 500, json: { error: 'Logout failed' } });
  });
  await page.getByRole('button', { name: /Fixture Player/ }).click();
  await page.getByRole('button', { name: 'Sign out' }).click();
  await expect(page).toHaveURL(/login$/);
  expect(await page.evaluate(() => localStorage.getItem('gm_token'))).toBeNull();
  await expect.poll(() => authorization).toBe('Bearer fixture-token');
  // Sign in again in the same document. The late logout reply must leave this session alone.
  await page.getByLabel('Email').fill(user.email);
  await page.getByLabel('Password', { exact: true }).fill('fixture-password');
  await page.getByRole('button', { name: /^Sign in$/i }).click();
  await expect(page.getByRole('heading', { name: 'Campaigns', exact: true })).toBeVisible();
  await page.locator('a[href="/play/alpha"]').click();
  const reply = page.waitForResponse('**/auth/logout');
  release();
  await (await reply).finished();
  // The app asks the server to set the newer session's cookie again.
  await expect.poll(() => game.requests).toContain('GET /auth/me');
  // The reply has arrived; give the app a rendered frame to react to it.
  await page.evaluate(() => new Promise(resolve => requestAnimationFrame(() => setTimeout(resolve, 100))));
  await expect(page).toHaveURL(/play\/alpha$/);
  expect(await page.evaluate(() => localStorage.getItem('gm_token'))).toBe('fixture-token');
  expect(game.requests).toContain('POST /auth/login');
});

test('shows one error panel when the campaign cannot be loaded', async ({ page, game }) => {
  await signIn(page);
  await page.route('**/api/v1/campaigns/alpha', route => route.fulfill({ status: 500, json: { error: 'Campaign unavailable' } }));
  await page.goto('/play/alpha');
  await expect(page.getByText('The campaign could not be loaded.', { exact: true })).toBeVisible();
  await expect(page.getByText('Campaign unavailable', { exact: true })).toHaveCount(1);
  expect(game.sockets).toHaveLength(0);
});
