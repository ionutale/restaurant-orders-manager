import { api } from '$lib/api';

type User = {
	id: number;
	name: string;
	email: string;
	role: 'waiter' | 'chef' | 'manager';
};

class AuthStore {
	token = $state<string | null>(null);
	user = $state<User | null>(null);
	loading = $state(true);

	constructor() {
		const saved = localStorage.getItem('token');
		if (saved) {
			this.token = saved;
			this.loadUser();
		} else {
			this.loading = false;
		}
	}

	get isLoggedIn() {
		return this.token !== null && this.user !== null;
	}

	get role() {
		return this.user?.role;
	}

	async loadUser() {
		try {
			this.user = await api.me();
		} catch {
			this.token = null;
			this.user = null;
			localStorage.removeItem('token');
		} finally {
			this.loading = false;
		}
	}

	async login(email: string, password: string) {
		const data = await api.login(email, password);
		this.token = data.token;
		this.user = data.user;
		localStorage.setItem('token', data.token);
	}

	async register(name: string, email: string, password: string, role: string) {
		const data = await api.register(name, email, password, role);
		this.token = data.token;
		this.user = data.user;
		localStorage.setItem('token', data.token);
	}

	logout() {
		this.token = null;
		this.user = null;
		localStorage.removeItem('token');
	}
}

export const auth = new AuthStore();
