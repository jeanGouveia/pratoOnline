import type { User } from '$lib/types/user';

// Svelte 5 Runes: $state reativo e global via module singleton
// Importe { userStore } em qualquer componente para acessar/mutar o usuário logado

function createUserStore() {
  let user = $state<User | null>(null);
  let loading = $state(true); // true até o layout raiz resolver a sessão SSR

  return {
    get user()    { return user; },
    get loading() { return loading; },

    setUser(u: User | null) { user = u; },
    setLoading(v: boolean)  { loading = v; },

    logout() { user = null; }
  };
}

export const userStore = createUserStore();
