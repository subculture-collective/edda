import { expect, it, vi } from 'vitest';

import { logout } from '../src/api/auth';
import { getStoredToken, setStoredSession } from '../src/api/authSession';

const user = { id: '1', name: 'Player', email: 'player@example.com' };

function stubFetch() {
  const calls: { url: string; init: RequestInit }[] = [];
  let reply: (response: Response) => void = () => {};
  vi.stubGlobal('fetch', vi.fn((url: string, init: RequestInit) => {
    calls.push({ url, init });
    if (url.endsWith('/auth/logout')) return new Promise<Response>(resolve => { reply = resolve; });
    return Promise.resolve(Response.json({ user }));
  }));
  return { calls, reply: (status: number) => reply(new Response(null, { status })) };
}

const bearer = (init: RequestInit) => new Headers(init.headers).get('Authorization');

it('sends logout with keepalive so that it survives page unload', async () => {
  setStoredSession('old-token', user);
  const server = stubFetch();

  const done = logout();

  expect(getStoredToken()).toBeNull();
  expect(server.calls).toHaveLength(1);
  expect(server.calls[0].url).toBe('/api/v1/auth/logout');
  expect(server.calls[0].init).toMatchObject({ method: 'POST', keepalive: true, credentials: 'include' });
  expect(bearer(server.calls[0].init)).toBe('Bearer old-token');
  server.reply(204);
  await done;
  // Nobody has signed in since, so there is no cookie to restore.
  expect(server.calls).toHaveLength(1);
});

it.each([204, 500])('restores the cookie of a newer sign-in after a late %i logout reply', async status => {
  setStoredSession('old-token', user);
  const server = stubFetch();
  const done = logout().catch(() => {});
  setStoredSession('new-token', user);
  expect(server.calls).toHaveLength(1);

  server.reply(status);
  await done;

  expect(server.calls.map(call => call.url)).toEqual(['/api/v1/auth/logout', '/api/v1/auth/me']);
  expect(bearer(server.calls[1].init)).toBe('Bearer new-token');
  expect(server.calls[1].init.credentials).toBe('include');
});
