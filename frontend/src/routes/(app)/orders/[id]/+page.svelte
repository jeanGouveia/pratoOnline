<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { getOrder } from '$lib/api/order';
  import type { Order } from '$lib/types/order';
  import { ORDER_STATUS_LABEL, ORDER_STATUS_COLOR } from '$lib/types/order';

  const orderId = Number($page.params.id);

  let order: Order | null = $state(null);
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    loading = true;
    error = '';
    try {
      console.log('OrderDetail - Carregando pedido ID:', orderId);
      order = await getOrder(orderId);
      console.log('OrderDetail - Pedido carregado:', order);
      console.log('OrderDetail - order.id:', order.id);
      console.log('OrderDetail - order.ID:', order.ID);
      console.log('OrderDetail - order.status:', order.status);
      console.log('OrderDetail - order.Status:', order.Status);
      console.log('OrderDetail - order.items:', order.items);
      console.log('OrderDetail - order.Items:', order.Items);
      if (order.items) {
        console.log('OrderDetail - Primeiro item:', order.items[0]);
        console.log('OrderDetail - item.product_name:', order.items[0].product_name);
        console.log('OrderDetail - item.Product:', order.items[0].Product);
      }
    } catch (e: any) {
      console.error('OrderDetail - Erro ao carregar pedido:', e);
      console.error('OrderDetail - Stack trace:', e.stack);
      error = e?.message ?? 'Erro ao carregar pedido.';
    } finally {
      loading = false;
    }
  });

  function formatDate(d?: string) {
    if (!d) return '—';
    return new Intl.DateTimeFormat('pt-BR', {
      dateStyle: 'full',
      timeStyle: 'short',
    }).format(new Date(d));
  }

  function formatPrice(v?: number) {
    if (v == null) return '—';
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(v);
  }

  // Mapeia a progressão de status para a barra visual
  const STATUS_STEPS = ['pending', 'confirmed', 'preparing', 'ready', 'delivered'] as const;
  type ProgressStatus = typeof STATUS_STEPS[number];

  const currentStepIndex = $derived(
    order ? STATUS_STEPS.indexOf(order.status as ProgressStatus) : -1
  );
  const isCancelled = $derived(order?.status === 'cancelled');
</script>

<div class="page-wrapper">
  <a href="/orders" class="back-link">← Voltar para Pedidos</a>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <span>Carregando pedido…</span>
    </div>

  {:else if error}
    <div class="alert alert-error">
      <span>⚠️ {error}</span>
      <a href="/orders" class="btn btn-sm">Voltar</a>
    </div>

  {:else if order}
    <header class="order-header">
      <div class="order-meta">
        <h1 class="order-title">Pedido <span class="order-id-num">#{order.id}</span></h1>
        <span class="status-badge {ORDER_STATUS_COLOR[order.status]}">
          {ORDER_STATUS_LABEL[order.status]}
        </span>
      </div>
      <p class="order-date">{formatDate(order.created_at)}</p>
    </header>

    <!-- Barra de progresso (oculta se cancelado) -->
    {#if !isCancelled}
      <div class="progress-track">
        {#each STATUS_STEPS as step, i}
          <div class="progress-step" class:done={i <= currentStepIndex} class:active={i === currentStepIndex}>
            <div class="step-dot"></div>
            <span class="step-label">{ORDER_STATUS_LABEL[step]}</span>
          </div>
          {#if i < STATUS_STEPS.length - 1}
            <div class="progress-line" class:filled={i < currentStepIndex}></div>
          {/if}
        {/each}
      </div>
    {:else}
      <div class="cancelled-banner">🚫 Este pedido foi cancelado</div>
    {/if}

    <div class="order-layout">
      <!-- Itens do pedido -->
      <section class="order-section">
        <h2 class="section-title">Itens do Pedido</h2>

        {#if !order.items || order.items.length === 0}
          <p class="empty-note">Nenhum item registrado neste pedido.</p>
        {:else}
          <div class="items-table-wrapper">
            <table class="items-table">
              <thead>
                <tr>
                  <th>Produto</th>
                  <th class="text-right">Qtd</th>
                  <th class="text-right">Preço unit.</th>
                  <th class="text-right">Subtotal</th>
                </tr>
              </thead>
              <tbody>
                {#each order.items as item}
                  <tr>
                    <td class="item-name">
                      {item.product_name ?? `Produto #${item.product_id}`}
                    </td>
                    <td class="text-right">{item.quantity}</td>
                    <td class="text-right muted">{formatPrice(item.unit_price)}</td>
                    <td class="text-right bold">
                      {item.unit_price != null
                        ? formatPrice(item.unit_price * item.quantity)
                        : '—'}
                    </td>
                  </tr>
                {/each}
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="3" class="total-label">Total do Pedido</td>
                  <td class="text-right total-value">{formatPrice(order.total)}</td>
                </tr>
              </tfoot>
            </table>
          </div>
        {/if}
      </section>

      <!-- Informações adicionais -->
      <aside class="order-aside">
        <div class="info-card">
          <h3 class="info-title">Detalhes</h3>
          <dl class="info-list">
            <dt>Status atual</dt>
            <dd>
              <span class="status-badge {ORDER_STATUS_COLOR[order.status]}">
                {ORDER_STATUS_LABEL[order.status]}
              </span>
            </dd>

            <dt>Criado em</dt>
            <dd>{formatDate(order.created_at)}</dd>

            {#if order.updated_at && order.updated_at !== order.created_at}
              <dt>Atualizado em</dt>
              <dd>{formatDate(order.updated_at)}</dd>
            {/if}

            {#if order.notes}
              <dt>Observações</dt>
              <dd class="notes">{order.notes}</dd>
            {/if}

            <dt>Total</dt>
            <dd class="total-dd">{formatPrice(order.total)}</dd>
          </dl>
        </div>

        <!-- Nota sobre estoque -->
        {#if order.status === 'confirmed' || order.status === 'preparing' || order.status === 'ready' || order.status === 'delivered'}
          <div class="stock-note">
            <span class="stock-icon">📦</span>
            <p>
              A baixa de estoque dos ingredientes foi realizada automaticamente
              ao confirmar este pedido.
            </p>
          </div>
        {/if}

        <div class="aside-actions">
          <a href="/orders/new" class="btn btn-primary">+ Novo Pedido</a>
          <a href="/orders" class="btn btn-ghost">Ver todos os pedidos</a>
        </div>
      </aside>
    </div>
  {/if}
</div>

<style>
  .page-wrapper    { max-width: 960px; margin: 0 auto; padding: 2rem 1.5rem; }
  .back-link       { display: inline-flex; align-items: center; gap: 0.3rem; color: var(--color-muted); font-size: 0.875rem; text-decoration: none; margin-bottom: 1.75rem; }
  .back-link:hover { color: var(--color-text); }

  /* Header */
  .order-header  { margin-bottom: 0.5rem; }
  .order-meta    { display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; margin-bottom: 0.4rem; }
  .order-title   { font-size: 1.6rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .order-id-num  { color: var(--color-primary, #e85d04); }
  .order-date    { color: var(--color-muted); font-size: 0.875rem; margin: 0 0 1.5rem; }

  /* Barra de progresso */
  .progress-track {
    display: flex;
    align-items: center;
    margin-bottom: 2rem;
    overflow-x: auto;
    padding-bottom: 0.25rem;
  }
  .progress-step {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.35rem;
    flex-shrink: 0;
  }
  .step-dot {
    width: 1rem; height: 1rem;
    border-radius: 50%;
    border: 2px solid var(--color-border, #d1d5db);
    background: var(--color-surface, #fff);
    transition: all 0.2s;
  }
  .progress-step.done .step-dot  { background: var(--color-primary, #e85d04); border-color: var(--color-primary, #e85d04); }
  .progress-step.active .step-dot { box-shadow: 0 0 0 3px rgba(232,93,4,0.2); }
  .step-label    { font-size: 0.7rem; color: var(--color-muted); white-space: nowrap; }
  .progress-step.done .step-label  { color: var(--color-primary, #e85d04); font-weight: 600; }
  .progress-step.active .step-label { color: var(--color-text); font-weight: 700; }
  .progress-line {
    flex: 1;
    height: 2px;
    background: var(--color-border, #e5e7eb);
    margin: 0 0.5rem;
    margin-bottom: 1.4rem; /* alinha com o dot */
    min-width: 2rem;
    transition: background 0.2s;
  }
  .progress-line.filled { background: var(--color-primary, #e85d04); }

  .cancelled-banner {
    background: #fee2e2;
    color: #b91c1c;
    border: 1px solid #fca5a5;
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    font-weight: 600;
    font-size: 0.9rem;
    margin-bottom: 1.5rem;
  }

  /* Layout 2 colunas */
  .order-layout  { display: grid; grid-template-columns: 1fr 280px; gap: 1.5rem; align-items: start; }
  .order-section { min-width: 0; }
  .section-title { font-size: 1rem; font-weight: 600; margin: 0 0 1rem; color: var(--color-text); }

  /* Tabela de itens */
  .items-table-wrapper { border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; overflow: hidden; }
  .items-table   { width: 100%; border-collapse: collapse; font-size: 0.875rem; }
  .items-table th { text-align: left; padding: 0.7rem 1rem; background: var(--color-surface-2, #f9fafb); font-weight: 600; font-size: 0.78rem; text-transform: uppercase; letter-spacing: 0.04em; color: var(--color-muted); }
  .items-table td { padding: 0.75rem 1rem; border-top: 1px solid var(--color-border, #e5e7eb); }
  .items-table tr:hover td { background: var(--color-surface-2, #f9fafb); }
  .items-table tfoot td { border-top: 2px solid var(--color-border, #e5e7eb); background: var(--color-surface-2, #f9fafb); padding: 0.85rem 1rem; }
  .item-name     { font-weight: 500; }
  .text-right    { text-align: right; }
  .muted         { color: var(--color-muted); }
  .bold          { font-weight: 600; }
  .total-label   { font-weight: 600; color: var(--color-text); }
  .total-value   { font-weight: 700; font-size: 1rem; color: var(--color-primary, #e85d04); }
  .empty-note    { color: var(--color-muted); font-size: 0.875rem; }

  /* Aside */
  .order-aside   { display: flex; flex-direction: column; gap: 1rem; }
  .info-card     { background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; padding: 1.25rem; }
  .info-title    { font-size: 0.9rem; font-weight: 700; margin: 0 0 0.85rem; color: var(--color-text); }
  .info-list     { display: grid; grid-template-columns: auto 1fr; gap: 0.5rem 1rem; margin: 0; font-size: 0.83rem; }
  .info-list dt  { color: var(--color-muted); font-weight: 500; white-space: nowrap; align-self: start; padding-top: 0.1rem; }
  .info-list dd  { margin: 0; color: var(--color-text); word-break: break-word; }
  .info-list .notes { font-style: italic; color: var(--color-muted); }
  .info-list .total-dd { font-weight: 700; color: var(--color-primary, #e85d04); font-size: 0.95rem; }

  .stock-note    { display: flex; gap: 0.6rem; align-items: flex-start; background: #f0fdf4; border: 1px solid #86efac; border-radius: 0.5rem; padding: 0.75rem; font-size: 0.8rem; color: #15803d; line-height: 1.5; }
  .stock-icon    { flex-shrink: 0; font-size: 1rem; }
  .stock-note p  { margin: 0; }

  .aside-actions { display: flex; flex-direction: column; gap: 0.5rem; }

  /* Status badges */
  .status-badge  { font-size: 0.75rem; font-weight: 600; padding: 0.25rem 0.6rem; border-radius: 99px; }
  .badge-warning { background: #fef9c3; color: #854d0e; }
  .badge-info    { background: #dbeafe; color: #1e40af; }
  .badge-success { background: #dcfce7; color: #15803d; }
  .badge-neutral { background: var(--color-surface-2, #f3f4f6); color: var(--color-muted); }
  .badge-error   { background: #fee2e2; color: #b91c1c; }

  /* Loading/Error */
  .loading-state { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; color: var(--color-muted); }
  .spinner       { width: 2rem; height: 2rem; border: 3px solid var(--color-border, #e5e7eb); border-top-color: var(--color-primary, #e85d04); border-radius: 50%; animation: spin 0.7s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .alert-error   { display: flex; align-items: center; gap: 1rem; background: #fef2f2; border: 1px solid #fca5a5; color: #b91c1c; padding: 1rem 1.25rem; border-radius: 0.5rem; }

  /* Botões */
  .btn           { display: inline-flex; align-items: center; justify-content: center; padding: 0.55rem 1rem; border-radius: 0.5rem; font-size: 0.875rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s; text-decoration: none; }
  .btn-primary   { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover { background: var(--color-primary-dark, #c84e00); }
  .btn-ghost     { background: transparent; color: var(--color-primary, #e85d04); border: 1.5px solid var(--color-primary, #e85d04); }
  .btn-ghost:hover { background: rgba(232,93,4,0.06); }
  .btn-sm        { padding: 0.35rem 0.75rem; font-size: 0.8rem; }

  @media (max-width: 700px) {
    .order-layout { grid-template-columns: 1fr; }
    .order-aside  { order: -1; } /* resumo sobe no mobile */
  }
</style>
