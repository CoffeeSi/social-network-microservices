/**
 * @param {Record<string, unknown> | null | undefined} obj
 * @param {string[]} keys
 */
export function pickTimestampField(obj, ...keys) {
  if (!obj) return null;
  for (const key of keys) {
    if (obj[key] != null) return obj[key];
  }
  return null;
}

/**
 * Parse API timestamps (ISO string, unix ms, protobuf { seconds, nanos }).
 * @param {unknown} value
 * @returns {Date | null}
 */
export function parseTimestamp(value) {
  if (value == null || value === '') return null;

  if (value instanceof Date) {
    return Number.isNaN(value.getTime()) ? null : value;
  }

  if (typeof value === 'number' && Number.isFinite(value)) {
    const ms = value < 1e12 ? value * 1000 : value;
    const d = new Date(ms);
    return Number.isNaN(d.getTime()) ? null : d;
  }

  if (typeof value === 'string') {
    const trimmed = value.trim();
    if (/^\d+$/.test(trimmed)) {
      const n = Number(trimmed);
      const ms = n < 1e12 ? n * 1000 : n;
      const d = new Date(ms);
      return Number.isNaN(d.getTime()) ? null : d;
    }
    const d = new Date(trimmed);
    return Number.isNaN(d.getTime()) ? null : d;
  }

  if (typeof value === 'object') {
    const o = /** @type {Record<string, unknown>} */ (value);
    if (o.seconds != null || o.nanos != null) {
      const sec = Number(o.seconds ?? 0);
      const nano = Number(o.nanos ?? 0);
      if (!Number.isFinite(sec)) return null;
      const d = new Date(sec * 1000 + Math.floor(nano / 1e6));
      return Number.isNaN(d.getTime()) ? null : d;
    }
  }

  return null;
}

/**
 * @param {unknown} value
 * @param {Intl.DateTimeFormatOptions} [options]
 * @returns {string}
 */
export function formatDateTime(value, options) {
  const d = parseTimestamp(value);
  if (!d) return '—';
  try {
    return d.toLocaleString(undefined, options);
  } catch {
    return d.toISOString();
  }
}

/**
 * @param {unknown} value
 * @returns {string}
 */
export function formatTimeShort(value) {
  const d = parseTimestamp(value);
  if (!d) return '';
  return d.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });
}

/**
 * @param {unknown} value
 * @returns {string | undefined}
 */
export function toDateTimeAttr(value) {
  const d = parseTimestamp(value);
  return d ? d.toISOString() : undefined;
}
