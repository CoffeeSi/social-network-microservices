import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { authStorage } from '@/shared/lib/authStorage';
import { getAccessTokenUserId } from '@/shared/lib/jwtPayload';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [accessToken, setAccessToken] = useState(null);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    setAccessToken(authStorage.getAccessToken());
    setReady(true);
  }, []);

  const userId = useMemo(() => getAccessTokenUserId(accessToken), [accessToken]);

  const signIn = useCallback((access, refresh) => {
    authStorage.setTokens(access, refresh);
    setAccessToken(access);
  }, []);

  const signOut = useCallback(() => {
    authStorage.clear();
    setAccessToken(null);
  }, []);

  const value = useMemo(
    () => ({
      ready,
      isAuthenticated: Boolean(accessToken),
      accessToken,
      userId,
      signIn,
      signOut,
    }),
    [accessToken, ready, userId, signIn, signOut],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
