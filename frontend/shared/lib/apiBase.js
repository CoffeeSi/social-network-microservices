/** API gateway origin for browser requests (no trailing slash). */
export function getPublicApiBase() {
  if (typeof process === 'undefined') return '';
  const raw = process.env.NEXT_PUBLIC_API_BASE_URL ?? '';
  return String(raw).replace(/\/$/, '');
}
