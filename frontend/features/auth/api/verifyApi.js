import { apiJson } from '@/shared/api/client';
import { endpoints } from '@/shared/api/endpoints';

/** @param {{ token?: string, email?: string }} payload */
export function verifyEmail(payload) {
  return apiJson(endpoints.auth.verify, { method: 'POST', body: payload, skipAuth: true });
}

export default verifyEmail;
