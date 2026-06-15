import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// "/" redireciona para /dashboard (protegida) ou /login (guard cuida do resto)
export const load: PageServerLoad = () => {
  redirect(302, '/dashboard');
};
