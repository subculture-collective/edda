import { apiFetch, apiFetchVoid } from './backend';
import { buildAuthHeaders, clearStoredSession, getStoredToken, type AuthSession, type AuthUser } from './authSession';

export function register(name: string, email: string, password: string): Promise<AuthSession> {
  return apiFetch<AuthSession>('/auth/register', {
    method: 'POST',
    body: { name, email, password },
  });
}

export function login(email: string, password: string): Promise<AuthSession> {
  return apiFetch<AuthSession>('/auth/login', {
    method: 'POST',
    body: { email, password },
  });
}

export async function getMe(token: string): Promise<{ user: AuthUser }> {
  return apiFetch<{ user: AuthUser }>('/auth/me', {
    method: 'GET',
    credentials: 'include',
    headers: (() => {
      const headers = buildAuthHeaders(token);
      return headers;
    })(),
  });
}

export function logout(): Promise<void> {
  const headers = buildAuthHeaders(getStoredToken());
  // Clear before the request so that a reload cannot restore the session.
  clearStoredSession();
  // keepalive lets the request finish, and revoke the token, if the page closes.
  return apiFetchVoid('/auth/logout', { method: 'POST', headers, keepalive: true });
}

// A one-use ticket that opens a WebSocket for the stored session within 30 seconds.
export async function requestWebSocketTicket(): Promise<string> {
  const { ticket } = await apiFetch<{ ticket: string }>('/auth/ws-ticket', { method: 'POST' });
  return ticket;
}
