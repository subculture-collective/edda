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
  // Clear before the request: a reload would otherwise restore the session, and a
  // late reply would erase the one from a newer sign-in.
  clearStoredSession();
  // keepalive lets the request finish, and clear the cookie, if the page closes.
  return apiFetchVoid('/auth/logout', { method: 'POST', headers, keepalive: true }).finally(restoreSessionCookie);
}

// A logout reply that arrives after a newer sign-in clears that sign-in's
// WebSocket cookie. /auth/me sets it again.
async function restoreSessionCookie(): Promise<void> {
  const token = getStoredToken();
  if (token) await getMe(token).catch(() => {});
}
