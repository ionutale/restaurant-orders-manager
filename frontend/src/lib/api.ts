import { API_BASE } from '$lib/config';

export class ApiError extends Error {
	constructor(public status: number, message: string) {
		super(message);
	}
}

async function request(path: string, opts: RequestInit = {}) {
	const token = localStorage.getItem('token');
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(opts.headers as Record<string, string>),
	};
	if (token) headers['Authorization'] = `Bearer ${token}`;

	const res = await fetch(`${API_BASE}${path}`, { ...opts, headers });
	if (!res.ok) {
		const body = await res.json().catch(() => ({ error: res.statusText }));
		throw new ApiError(res.status, body.error);
	}
	return res.json();
}

export const api = {
	login: (email: string, password: string) =>
		request('/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) }),

	register: (name: string, email: string, password: string, role: string) =>
		request('/auth/register', { method: 'POST', body: JSON.stringify({ name, email, password, role }) }),

	me: () => request('/me'),
};
