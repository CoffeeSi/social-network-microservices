import { apiJson } from '@/shared/api/client';
import { endpoints } from '@/shared/api/endpoints';

/** @param {{ page_size?: number, page?: number }} [params] */
export function listUsers(params = {}) {
  const qs = new URLSearchParams();
  if (params.page_size != null) qs.set('page_size', String(params.page_size));
  if (params.page != null) qs.set('page', String(params.page));
  const q = qs.toString();
  return apiJson(`${endpoints.users.list}${q ? `?${q}` : ''}`);
}

/** @param {string} id */
export function getUser(id) {
  return apiJson(endpoints.users.one(id));
}

/** @param {{ first_name: string, last_name: string, email: string, password: string, dob?: string }} payload */
export function createUser(payload) {
  return apiJson(endpoints.users.create, { method: 'POST', body: payload });
}

/** @param {string} id @param {Record<string, unknown>} data */
export function updateUser(id, data) {
  return apiJson(endpoints.users.one(id), { method: 'PATCH', body: { data } });
}

/** @param {string} id */
export function deleteUser(id) {
  return apiJson(endpoints.users.one(id), { method: 'DELETE' });
}
