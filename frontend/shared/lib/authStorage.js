const ACCESS = 'sn_access_token';
const REFRESH = 'sn_refresh_token';

export const authStorage = {
  getAccessToken: () =>
    typeof window !== 'undefined' ? localStorage.getItem(ACCESS) : null,
  getRefreshToken: () =>
    typeof window !== 'undefined' ? localStorage.getItem(REFRESH) : null,
  setTokens(accessToken, refreshToken) {
    if (typeof window === 'undefined') return;
    localStorage.setItem(ACCESS, accessToken);
    if (refreshToken != null) localStorage.setItem(REFRESH, refreshToken);
  },
  clear() {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(ACCESS);
    localStorage.removeItem(REFRESH);
  },
};
