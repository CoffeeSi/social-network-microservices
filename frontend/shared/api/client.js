import { authStorage } from '@/shared/lib/authStorage';
import { getPublicApiBase } from '@/shared/lib/apiBase';

async function parseError(res) {
  const text = await res.text();
  try {
    const json = JSON.parse(text);
    return json.message || json.error || text || res.statusText;
  } catch {
    return text || res.statusText;
  }
}

/**
 * @param {string} path
 * @param {RequestInit & { skipAuth?: boolean }} [options]
 */
export async function apiFetch(path, options = {}) {
  const { skipAuth = false, ...fetchOptions } = options;
  const base = getPublicApiBase();
  const url = path.startsWith('http') ? path : `${base}${path}`;

  const headers = new Headers(fetchOptions.headers);
  if (
    fetchOptions.body != null &&
    !(fetchOptions.body instanceof FormData) &&
    !headers.has('Content-Type')
  ) {
    headers.set('Content-Type', 'application/json');
  }

  const access = authStorage.getAccessToken();
  if (access && !skipAuth && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${access}`);
  }

  const res = await fetch(url, { ...fetchOptions, headers });

  if (res.status === 401 && !skipAuth && path !== '/api/v1/auth/refresh') {
    const refreshed = await tryRefresh();
    if (refreshed) {
      headers.set('Authorization', `Bearer ${authStorage.getAccessToken()}`);
      return fetch(url, { ...fetchOptions, headers }).then((r) => handleResponse(r));
    }
  }

  return handleResponse(res);
}

async function handleResponse(res) {
  if (res.ok) {
    if (res.status === 204) return null;
    const ct = res.headers.get('content-type');
    if (ct?.includes('application/json')) return res.json();
    return res.text();
  }
  const msg = await parseError(res);
  throw new Error(msg || `Request failed (${res.status})`);
}

async function tryRefresh() {
  const refresh = authStorage.getRefreshToken();
  if (!refresh) return false;
  try {
    const base = getPublicApiBase();
    const res = await fetch(`${base}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refresh }),
    });
    if (!res.ok) return false;
    const data = await res.json();
    if (data.access_token) authStorage.setTokens(data.access_token, refresh);
    return true;
  } catch {
    return false;
  }
}

export function apiJson(path, options = {}) {
  const { body, ...rest } = options;
  return apiFetch(path, {
    ...rest,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });
}
