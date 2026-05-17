import { apiJson } from '@/shared/api/client';
import { endpoints } from '@/shared/api/endpoints';

/** @param {{ first_name: string, last_name: string, dob: string, email: string, password: string }} payload */
export function registerUser(payload) {
  return apiJson(endpoints.auth.register, { method: 'POST', body: payload, skipAuth: true });
}

/** @param {{ email: string, password: string }} payload */
export function loginUser(payload) {
  return apiJson(endpoints.auth.login, { method: 'POST', body: payload, skipAuth: true });
}
