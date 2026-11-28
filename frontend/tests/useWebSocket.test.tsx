import { act, renderHook } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { setStoredSession } from '../src/api/authSession';
import { useWebSocket } from '../src/hooks/useWebSocket';
import { result, timestamp } from './fixtures/data';

class FixtureSocket extends EventTarget {
  static OPEN = 1;
  static CONNECTING = 0;
  static instances: FixtureSocket[] = [];
  readyState = 0;
  sent: string[] = [];
  url: string;
  constructor(url: string) { super(); this.url = url; FixtureSocket.instances.push(this); }
  send(message: string) { this.sent.push(message); }
  close() { this.readyState = 3; this.dispatchEvent(new Event('close')); }
  open() { this.readyState = 1; this.dispatchEvent(new Event('open')); }
  receive(payload: unknown) { this.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'result', timestamp, payload }) })); }
}
it('completes turns, reconnects, and closes the old campaign socket on switch/unmount', async () => {
  vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
  const { result: hook, rerender, unmount } = renderHook(({ id }) => useWebSocket(id), { initialProps: { id: 'alpha' } });
  await act(() => vi.advanceTimersByTimeAsync(0));
  const first = FixtureSocket.instances[0]; act(() => first.open());
  act(() => { expect(hook.current.sendAction(' Look ')).toBe(true); });
  expect(JSON.parse(first.sent[0])).toEqual({ type: 'action', payload: { input: 'Look', action_id: expect.any(String) } });
  expect(hook.current.isLoading).toBe(true);
  act(() => first.receive(result)); expect(hook.current.isLoading).toBe(false);
  expect(hook.current.events).toHaveLength(1);
  act(() => first.close()); expect(hook.current.connectionStatus).toBe('closed');
  await act(() => vi.advanceTimersByTimeAsync(1000));
  await act(() => vi.advanceTimersByTimeAsync(0));
  const second = FixtureSocket.instances[1]; act(() => second.open());
  expect(hook.current.isConnected).toBe(true);
  rerender({ id: 'beta' }); await act(() => vi.advanceTimersByTimeAsync(0));
  expect(second.readyState).toBe(3);
  const third = FixtureSocket.instances[2]; expect(third.url).toContain('/beta/ws');
  act(() => third.open()); unmount(); expect(third.readyState).toBe(3);
  await act(() => vi.advanceTimersByTimeAsync(30000)); expect(FixtureSocket.instances).toHaveLength(3);
});

it('reuses the action ID for recovery and keeps a new action distinct', async () => {
 vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
 const hook = renderHook(() => useWebSocket('alpha'));
 await act(() => vi.advanceTimersByTimeAsync(0));
 const socket = FixtureSocket.instances[0]; act(() => socket.open());
 act(() => { hook.result.current.sendAction('Look'); });
 const first = JSON.parse(socket.sent[0]);
 act(() => socket.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'error', timestamp, payload: { error: 'Recover recorded outcomes', action_id: first.payload.action_id } }) })));
 act(() => { hook.result.current.sendAction('Look'); });
 expect(JSON.parse(socket.sent[1]).payload.action_id).toBe(first.payload.action_id);
 act(() => socket.receive(result));
 act(() => { hook.result.current.sendAction('Look'); });
 expect(JSON.parse(socket.sent[2]).payload.action_id).not.toBe(first.payload.action_id);
});

it('adopts a durable recovery ID returned after a fresh hook sends an action', async () => {
 vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
 const hook = renderHook(() => useWebSocket('alpha'));
 await act(() => vi.advanceTimersByTimeAsync(0));
 const socket = FixtureSocket.instances[0]; act(() => socket.open());
 act(() => { hook.result.current.sendAction('Original action'); });
 const canonical = '12345678-1234-4321-8123-123456789012';
 act(() => socket.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'error', timestamp, payload: { error: 'Recover recorded outcomes', action_id: canonical } }) })));
 act(() => { hook.result.current.sendAction('Original action'); });
 expect(JSON.parse(socket.sent[1]).payload.action_id).toBe(canonical);
});

it('offers an explicit retry only after a rejected action reconnects', async () => {
  vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
  const hook = renderHook(() => useWebSocket('alpha'));
  await act(() => vi.advanceTimersByTimeAsync(0));
  const first = FixtureSocket.instances[0]; act(() => first.open());
  act(() => { hook.result.current.sendAction('Look'); });
  expect(hook.result.current.canRetry).toBe(false);
  act(() => first.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'error', timestamp, payload: { error: 'an internal error occurred' } }) })));
  act(() => first.close());
  await act(() => vi.advanceTimersByTimeAsync(1000)); await act(() => vi.advanceTimersByTimeAsync(0));
  act(() => FixtureSocket.instances[1].open());
  expect(hook.result.current).toMatchObject({ canRetry: true });
  act(() => { expect(hook.result.current.retryAction()).toBe(true); });
  expect(JSON.parse(FixtureSocket.instances[1].sent[0]).payload).toEqual(JSON.parse(first.sent[0]).payload);
  act(() => FixtureSocket.instances[1].receive(result));
  expect(hook.result.current.canRetry).toBe(false);
});

it('does not carry a rejected action into another campaign', async () => {
  vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
  const hook = renderHook(({ id }) => useWebSocket(id), { initialProps: { id: 'alpha' } });
  await act(() => vi.advanceTimersByTimeAsync(0));
  const alpha = FixtureSocket.instances[0]; act(() => alpha.open());
  act(() => { hook.result.current.sendAction('Alpha action'); });
  act(() => alpha.dispatchEvent(new MessageEvent('message', { data: JSON.stringify({ type: 'error', timestamp, payload: { error: 'an internal error occurred' } }) })));
  expect(hook.result.current.canRetry).toBe(true);
  hook.rerender({ id: 'beta' }); await act(() => vi.advanceTimersByTimeAsync(0));
  act(() => FixtureSocket.instances[1].open());
  expect(hook.result.current.canRetry).toBe(false);
  expect(hook.result.current.retryAction()).toBe(false);
  expect(FixtureSocket.instances[1].sent).toEqual([]);
  hook.rerender({ id: 'alpha' }); await act(() => vi.advanceTimersByTimeAsync(0));
  const secondAlpha = FixtureSocket.instances[2]; act(() => secondAlpha.open());
  act(() => { hook.result.current.sendAction('Alpha action'); });
  expect(JSON.parse(secondAlpha.sent[0]).payload.action_id).not.toBe(JSON.parse(alpha.sent[0]).payload.action_id);
});

it('opens each connection with its own one-use ticket', async () => {
  vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
  setStoredSession('session-token', { id: '1', name: 'Player', email: 'player@example.com' });
  let issued = 0;
  const fetch = vi.fn(() => Promise.resolve(Response.json({ ticket: `ticket-${++issued}`, expires_in: 30 })));
  vi.stubGlobal('fetch', fetch);
  renderHook(() => useWebSocket('alpha'));
  await act(() => vi.advanceTimersByTimeAsync(0));
  const first = FixtureSocket.instances[0];
  expect(fetch).toHaveBeenCalledWith('/api/v1/auth/ws-ticket', expect.objectContaining({ method: 'POST' }));
  expect(new Headers((fetch.mock.calls[0] as unknown as [string, RequestInit])[1].headers).get('Authorization')).toBe('Bearer session-token');
  expect(new URL(first.url).searchParams.get('ticket')).toBe('ticket-1');
  act(() => first.open()); act(() => first.close());
  await act(() => vi.advanceTimersByTimeAsync(1000)); await act(() => vi.advanceTimersByTimeAsync(0));
  expect(new URL(FixtureSocket.instances[1].url).searchParams.get('ticket')).toBe('ticket-2');
});

it('stops reconnecting when the server closes the socket because the session signed out', async () => {
  vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
  const hook = renderHook(() => useWebSocket('alpha'));
  await act(() => vi.advanceTimersByTimeAsync(0));
  const socket = FixtureSocket.instances[0]; act(() => socket.open());
  act(() => { socket.readyState = 3; socket.dispatchEvent(new CloseEvent('close', { code: 4401, reason: 'signed out' })); });
  expect(hook.result.current).toMatchObject({ connectionStatus: 'closed', error: 'Your session has ended. Sign in again to keep playing.' });
  await act(() => vi.advanceTimersByTimeAsync(60000));
  expect(FixtureSocket.instances).toHaveLength(1);
});

it('does not open a socket when the ticket request is refused', async () => {
  vi.useFakeTimers(); FixtureSocket.instances = []; vi.stubGlobal('WebSocket', FixtureSocket);
  setStoredSession('revoked-token', { id: '1', name: 'Player', email: 'player@example.com' });
  vi.stubGlobal('fetch', vi.fn(() => Promise.resolve(Response.json({ error: 'invalid or expired token' }, { status: 401 }))));
  const hook = renderHook(() => useWebSocket('alpha'));
  await act(() => vi.advanceTimersByTimeAsync(60000));
  expect(FixtureSocket.instances).toHaveLength(0);
  expect(hook.result.current.error).toBe('Your session has ended. Sign in again to keep playing.');
});
