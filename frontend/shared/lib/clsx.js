/** @param {...(string | false | null | undefined)} parts */
export function clsx(...parts) {
  return parts.filter(Boolean).join(' ');
}
