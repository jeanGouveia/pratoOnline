<script lang="ts">
  import { onMount } from 'svelte';
  import { getOrders } from '$lib/api/order';
  import type { Order } from '$lib/types/order';
  import { ORDER_STATUS_LABEL, ORDER_STATUS_COLOR } from '$lib/types/order';

  let orders: Order[] = $state([]);
  let loading = $state(true);
  let error = $state('');
  let filter = $state<string>('all');

  onMount(loadOrders);

  async function loadOrders() {
    loading = true;
    error = '';
    try {
      orders = await getOrders();
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar pedidos.';
    } finally {
      loading = false;
    }
  }

  const filtered = $derived(
    filter === 'all' ? orders : orders.filter((o) => o.Status === filter)
  );

  function formatDate(d?: string) {
    if (!d) return '—';
    return new Intl.DateTimeFormat('pt-BR', { dateStyle: 'short', timeStyle: 'short' }).format(new Date(d));
  }

  function formatTotal(v: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(v);
  }

  const STATUS_OPTIONS = [
    { value: 'all', label: 'Todos' },
    { value: 'pending', label: 'Pendentes' },
    { value: 'confirmed', label: 'Confirmados' },
    { value: 'preparing', label: 'Preparando' },
    { value: 'ready', label: 'Prontos' },
    { value: 'delivered', label: 'Entregues' },
    { value: 'cancelled', label: 'Cancelados' },
  ];

  // Contagem por status para os pills
  const countByStatus = $derived(
    orders.reduce<Record<string, number>>((acc, o) => {
      acc[o.Status] = (acc[o.Status] ?? 0) + 1;
      return acc;
    }, {})
  );
</script>

<div class="page-wrapper">
  <header class="page-header">
    <div>
      <h1 class="page-title">Pedidos</h1>
      <p class="page-subtitle">{orders.length} pedido{orders.length !== 1 ? 's' : ''} registrado{orders.length !== 1 ? 's' : ''}</p>
    </div>
    <a href="/orders/new" class="btn btn-primary">+ Novo Pedido</a>
  </header>

  <!-- Filtros -->
  <div class="filter-row">
    {#each STATUS_OPTIONS as opt}
      <button
        class="filter-pill"
        class:active={filter === opt.value}
        onclick={() => (filter = opt.value)}
      >
        {opt.label}
        {#if opt.value !== 'all' && countByStatus[opt.value]}
          <span class="pill-count">{countByStatus[opt.value]}</span>
        {/if}
      </button>
    {/each}
  </div>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <span>Carregando pedidos…</span>
    </div>
  {:else if error}
    <div class="alert alert-error">
      <span>⚠️ {error}</span>
      <button class="btn btn-sm" onclick={loadOrders}>Tentar novamente</button>
    </div>
  {:else if filtered.length === 0}
    <div class="empty-state">
      <p class="empty-icon">🧾</p>
      <p class="empty-text">{filter === 'all' ? 'Nenhum pedido ainda.' : 'Nenhum pedido com esse status.'}</p>
      {#if filter === 'all'}
        <a href="/orders/new" class="btn btn-primary">Criar primeiro pedido</a>
      {:else}
        <button class="btn btn-ghost" onclick={() => (filter = 'all')}>Ver todos</button>
      {/if}
    </div>
  {:else}
    <div class="orders-list">
      {#each filtered as order}
        <a href="/orders/{order.ID}" class="order-card">
          <div class="order-card-left">
            <span class="order-id"># {order.ID}</span>
            <span class="order-date">{formatDate(order.CreatedAt)}</span>
          </div>
          <div class="order-items-preview">
            {#each (order.Items ?? []).slice(0, 3) as item}
              <span class="item-chip">{item.Quantity}× Produto #{item.ProductID}</span>
            {/each}
            {#if (order.Items ?? []).length > 3}
              <span class="item-chip muted">+{order.Items.length - 3} mais</span>
            {/if}
          </div>
          <div class="order-card-right">
            <span class="order-total">{formatTotal(order.TotalPrice)}</span>
            <span class="status-badge {ORDER_STATUS_COLOR[order.Status]}">{ORDER_STATUS_LABEL[order.Status]}</span>
          </div>
        </a>
      {/each}
    </div>
  {/if}
</div>

<style>
  .page-wrapper   { max-width: 1000px; margin: 0 auto; padding: 2rem 1.5rem; }
  .page-header    { display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 1rem; margin-bottom: 1.75rem; }
  .page-title     { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .page-subtitle  { font-size: 0.875rem; color: var(--color-muted); margin: 0.25rem 0 0; }

  /* Filtros */
  .filter-row     { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1.5rem; }
  .filter-pill    { border: 1px solid var(--color-border, #e5e7eb); background: var(--color-surface, #fff); color: var(--color-muted); font-size: 0.8rem; font-weight: 500; padding: 0.35rem 0.75rem; border-radius: 99px; cursor: pointer; transition: all 0.15s; display: flex; align-items: center; gap: 0.35rem; }
  .filter-pill:hover { border-color: var(--color-primary, #e85d04); color: var(--color-primary, #e85d04); }
  .filter-pill.active { background: var(--color-primary, #e85d04); border-color: var(--color-primary, #e85d04); color: #fff; }
  .pill-count     { background: rgba(255,255,255,0.3); padding: 0 0.3rem; border-radius: 99px; font-size: 0.7rem; font-weight: 700; }
  .filter-pill:not(.active) .pill-count { background: var(--color-surface-2, #f0f0f0); color: var(--color-text); }

  /* Lista pedidos */
  .orders-list    { display: flex; flex-direction: column; gap: 0.6rem; }
  .order-card     { display: flex; align-items: center; gap: 1rem; background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; padding: 1rem 1.25rem; text-decoration: none; color: inherit; transition: box-shadow 0.15s, border-color 0.15s; }
  .order-card:hover { box-shadow: 0 4px 14px rgba(0,0,0,0.07); border-color: var(--color-primary, #e85d04); }
  .order-card-left { display: flex; flex-direction: column; gap: 0.2rem; min-width: 90px; }
  .order-id       { font-weight: 700; color: var(--color-text); font-size: 0.95rem; }
  .order-date     { font-size: 0.75rem; color: var(--color-muted); }
  .order-items-preview { flex: 1; display: flex; gap: 0.4rem; flex-wrap: wrap; }
  .item-chip      { font-size: 0.78rem; background: var(--color-surface-2, #f5f5f5); padding: 0.2rem 0.5rem; border-radius: 4px; color: var(--color-text); }
  .item-chip.muted { color: var(--color-muted); }
  .order-card-right { display: flex; flex-direction: column; align-items: flex-end; gap: 0.4rem; min-width: 110px; }
  .order-total    { font-weight: 700; font-size: 0.95rem; color: var(--color-text); }

  /* Status badges */
  .status-badge   { font-size: 0.72rem; font-weight: 600; padding: 0.2rem 0.55rem; border-radius: 99px; white-space: nowrap; }
  .badge-warning  { background: #fef9c3; color: #854d0e; }
  .badge-info     { background: #dbeafe; color: #1e40af; }
  .badge-success  { background: #dcfce7; color: #15803d; }
  .badge-neutral  { background: var(--color-surface-2, #f3f4f6); color: var(--color-muted); }
  .badge-error    { background: #fee2e2; color: #b91c1c; }

  /* Estados */
  .loading-state  { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; color: var(--color-muted); }
  .spinner        { width: 2rem; height: 2rem; border: 3px solid var(--color-border, #e5e7eb); border-top-color: var(--color-primary, #e85d04); border-radius: 50%; animation: spin 0.7s linear infinite; }
  @keyframes spin  { to { transform: rotate(360deg); } }
  .empty-state    { text-align: center; padding: 4rem 1rem; }
  .empty-icon     { font-size: 2.5rem; margin: 0; }
  .empty-text     { color: var(--color-muted); margin: 0.5rem 0 1.25rem; }
  .alert-error    { display: flex; align-items: center; gap: 1rem; background: #fef2f2; border: 1px solid #fca5a5; color: #b91c1c; padding: 1rem 1.25rem; border-radius: 0.5rem; margin-bottom: 1.5rem; }

  .btn            { display: inline-flex; align-items: center; padding: 0.55rem 1.1rem; border-radius: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s; text-decoration: none; }
  .btn:disabled   { opacity: 0.55; cursor: not-allowed; }
  .btn-primary    { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover:not(:disabled) { background: var(--color-primary-dark, #c84e00); }
  .btn-ghost      { background: transparent; color: var(--color-muted); border: 1px solid var(--color-border, #e5e7eb); }
  .btn-ghost:hover { background: var(--color-surface-2, #f3f4f6); }
  .btn-sm         { padding: 0.35rem 0.75rem; font-size: 0.8rem; }

  @media (max-width: 640px) {
    .order-card         { flex-wrap: wrap; }
    .order-items-preview { width: 100%; }
    .order-card-right   { flex-direction: row; min-width: unset; width: 100%; justify-content: space-between; align-items: center; }
  }
</style>
