import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import type { PropsWithChildren } from 'react';

import { getMe, logout as logoutRequest } from '../api/auth';
import { clearStoredSession, getStoredToken, getStoredUser, setStoredSession, type AuthUser } from '../api/authSession';
import { AuthContext, type AuthContextValue } from './auth';

export function AuthProvider({ children }: PropsWithChildren) {
  const [token, setToken] = useState<string | null>(getStoredToken);
  const [user, setUser] = useState<AuthUser | null>(getStoredUser);
  // Only a stored token needs validating before the app can render.
  const [isLoading, setIsLoading] = useState(() => !!token);

  const didValidate = useRef(false);

  useEffect(() => {
    if (didValidate.current) return;
    didValidate.current = true;

    if (!token) return;

    // Validate token against the server on initial load.
    getMe(token)
      .then((res) => {
        setUser(res.user);
        setStoredSession(token, res.user);
      })
      .catch(() => {
        // Token expired or invalid — clear everything.
        clearStoredSession();
        setToken(null);
        setUser(null);
      })
      .finally(() => setIsLoading(false));
  }, [token]);

  const setSession = useCallback((newToken: string, newUser: AuthUser) => {
    setStoredSession(newToken, newUser);
    setToken(newToken);
    setUser(newUser);
  }, []);

  const logout = useCallback(async () => {
    setToken(null);
    setUser(null);

    try {
      await logoutRequest();
    } catch {
      // The local session is already cleared; a failed request changes nothing here.
    }
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      isAuthenticated: !!token && !!user,
      isLoading,
      setSession,
      logout,
    }),
    [user, token, isLoading, setSession, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
