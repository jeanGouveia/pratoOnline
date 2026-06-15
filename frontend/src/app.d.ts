// Tipagem dos locals SSR — disponível em todos os load() do servidor

import type { User } from '$lib/types/user';

declare global {
  namespace App {
    interface Locals {
      user: User | null;
    }
    // interface Error {}
    // interface PageData {}
    // interface Platform {}
  }
}

export {};
