import type { LayoutServerLoad } from './$types';

// Passa o usuário SSR para o $page.data de todos os layouts filhos
export const load: LayoutServerLoad = ({ locals }) => {
  return { user: locals.user ?? null };
};
