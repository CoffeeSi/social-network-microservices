/**
 * Decode JWT payload (no signature verification — only for UI; the gateway validates tokens).
 * Access tokens from auth-service use claim "sub" (user id) and "type": "access".
 */
export function decodeJwtPayload(token) {
  if (!token || typeof token !== 'string') return null;
  const parts = token.split('.');
  if (parts.length < 2) return null;
  try {
    const base64Url = parts[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const pad = base64.length % 4;
    const padded = pad ? base64 + '='.repeat(4 - pad) : base64;
    return JSON.parse(atob(padded));
  } catch {
    return null;
  }
}

/** @param {string | null} token access token */
export function getAccessTokenUserId(token) {
  const payload = decodeJwtPayload(token);
  if (!payload || payload.type !== 'access') return null;
  const sub = payload.sub;
  return typeof sub === 'string' ? sub : null;
}
