<script lang="ts">
  import { userStore } from '$lib/stores/userStore.svelte';
  import type { LayoutData } from './$types';

  let { data, children }: { data: LayoutData; children: any } = $props();

  // Sincroniza o estado SSR com o store reativo do cliente
  $effect(() => {
    userStore.setUser(data.user);
    userStore.setLoading(false);
  });
</script>

{@render children()}
