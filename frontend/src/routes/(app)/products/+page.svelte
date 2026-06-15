<script lang="ts">
  import { onMount } from 'svelte';
  import { getProducts, createProduct, getIngredients, createIngredient } from '$lib/api/product';
  import type { Product } from '$lib/types/product';
  import type { Ingredient } from '$lib/types/ingredient';

  let products: Product[] = $state([]);
  let ingredients: Ingredient[] = $state([]);
  let loading = $state(true);
  let error = $state('');

  // Modal novo produto
  let showProductModal = $state(false);
  let productForm = $state({ Name: '', Description: '', Price: 0, IsComposto: false, Active: true });
  let productSaving = $state(false);
  let productError = $state('');

  // Modal novo ingrediente
  let showIngModal = $state(false);
  let ingForm = $state({ Name: '', Unit: '', StockQuantity: 0, MinStock: 0 });
  let ingSaving = $state(false);
  let ingError = $state('');

  onMount(async () => {
    await loadAll();
  });

  async function loadAll() {
    loading = true;
    error = '';
    try {
      [products, ingredients] = await Promise.all([getProducts(), getIngredients()]);
    } catch (e: any) {
      error = e?.message ?? 'Erro ao carregar dados.';
    } finally {
      loading = false;
    }
  }

  async function saveProduct() {
    productSaving = true;
    productError = '';
    try {
      const created = await createProduct({
        name: productForm.Name,
        description: productForm.Description,
        price: Number(productForm.Price),
        is_composto: productForm.IsComposto,
        active: productForm.Active,
      });
      products = [...products, created];
      showProductModal = false;
      productForm = { Name: '', Description: '', Price: 0, IsComposto: false, Active: true };
    } catch (e: any) {
      productError = e?.message ?? 'Erro ao salvar produto.';
    } finally {
      productSaving = false;
    }
  }

  async function saveIngredient() {
    ingSaving = true;
    ingError = '';
    try {
      const created = await createIngredient({
        name: ingForm.Name,
        unit: ingForm.Unit,
        stock_quantity: Number(ingForm.StockQuantity),
        min_stock: Number(ingForm.MinStock),
      });
      ingredients = [...ingredients, created];
      showIngModal = false;
      ingForm = { Name: '', Unit: '', StockQuantity: 0, MinStock: 0 };
    } catch (e: any) {
      ingError = e?.message ?? 'Erro ao salvar ingrediente.';
    } finally {
      ingSaving = false;
    }
  }

  function formatPrice(value: number) {
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
  }
</script>

<div class="page-wrapper">
  <header class="page-header">
    <div>
      <h1 class="page-title">Cardápio</h1>
      <p class="page-subtitle">Gerencie produtos e ingredientes do seu restaurante</p>
    </div>
    <div class="header-actions">
      <button class="btn btn-secondary" onclick={() => (showIngModal = true)}>
        + Ingrediente
      </button>
      <button class="btn btn-primary" onclick={() => (showProductModal = true)}>
        + Produto
      </button>
    </div>
  </header>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <span>Carregando cardápio…</span>
    </div>
  {:else if error}
    <div class="alert alert-error">
      <span>⚠️ {error}</span>
      <button class="btn btn-sm" onclick={loadAll}>Tentar novamente</button>
    </div>
  {:else}
    <!-- Produtos -->
    <section class="section">
      <h2 class="section-title">Produtos <span class="badge">{products.length}</span></h2>

      {#if products.length === 0}
        <div class="empty-state">
          <p class="empty-icon">🍽️</p>
          <p class="empty-text">Nenhum produto cadastrado ainda.</p>
          <button class="btn btn-primary" onclick={() => (showProductModal = true)}>Criar primeiro produto</button>
        </div>
      {:else}
        <div class="card-grid">
          {#each products as product}
            <a href="/products/{product.ID}" class="product-card">
              <div class="product-card-header">
                <span class="product-name">{product.Name}</span>
                <span class="product-price">{formatPrice(product.Price)}</span>
              </div>
              {#if product.Description}
                <p class="product-desc">{product.Description}</p>
              {/if}
              <div class="product-footer">
                {#if product.IsComposto}
                  <span class="tag">Composto</span>
                {/if}
                <span class="status-dot" class:active={product.Active}>
                  {product.Active ? 'Ativo' : 'Inativo'}
                </span>
              </div>
            </a>
          {/each}
        </div>
      {/if}
    </section>

    <!-- Ingredientes -->
    <section class="section">
      <h2 class="section-title">Ingredientes <span class="badge">{ingredients.length}</span></h2>

      {#if ingredients.length === 0}
        <div class="empty-state">
          <p class="empty-icon">🥦</p>
          <p class="empty-text">Nenhum ingrediente cadastrado.</p>
          <button class="btn btn-secondary" onclick={() => (showIngModal = true)}>Adicionar ingrediente</button>
        </div>
      {:else}
        <div class="ing-table-wrapper">
          <table class="ing-table">
            <thead>
              <tr>
                <th>#</th>
                <th>Nome</th>
                <th>Unidade</th>
                <th>Estoque</th>
                <th>Estoque Mínimo</th>
              </tr>
            </thead>
            <tbody>
              {#each ingredients as ing}
                <tr>
                  <td class="muted">{ing.ID}</td>
                  <td class="bold">{ing.Name}</td>
                  <td>{ing.Unit}</td>
                  <td>{ing.StockQuantity}</td>
                  <td>{ing.MinStock}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </section>
  {/if}
</div>

<!-- Modal: Novo Produto -->
{#if showProductModal}
  <div class="modal-overlay" onclick={() => (showProductModal = false)}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <h2 class="modal-title">Novo Produto</h2>

      {#if productError}
        <p class="form-error">{productError}</p>
      {/if}

      <div class="form-group">
        <label for="p-name">Nome *</label>
        <input id="p-name" type="text" bind:value={productForm.Name} placeholder="Ex: Feijoada Completa" />
      </div>
      <div class="form-group">
        <label for="p-desc">Descrição</label>
        <textarea id="p-desc" bind:value={productForm.Description} rows="2" placeholder="Descrição opcional"></textarea>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label for="p-price">Preço (R$) *</label>
          <input id="p-price" type="number" min="0" step="0.01" bind:value={productForm.Price} />
        </div>
      </div>
      <div class="form-group checkbox-group">
        <label>
          <input type="checkbox" bind:checked={productForm.IsComposto} />
          Produto composto (requer ficha técnica)
        </label>
      </div>
      <div class="form-group checkbox-group">
        <label>
          <input type="checkbox" bind:checked={productForm.Active} />
          Produto ativo (visível no cardápio)
        </label>
      </div>

      <div class="modal-actions">
        <button class="btn btn-ghost" onclick={() => (showProductModal = false)}>Cancelar</button>
        <button class="btn btn-primary" onclick={saveProduct} disabled={productSaving || !productForm.Name}>
          {productSaving ? 'Salvando…' : 'Criar Produto'}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Modal: Novo Ingrediente -->
{#if showIngModal}
  <div class="modal-overlay" onclick={() => (showIngModal = false)}>
    <div class="modal" onclick={(e) => e.stopPropagation()}>
      <h2 class="modal-title">Novo Ingrediente</h2>

      {#if ingError}
        <p class="form-error">{ingError}</p>
      {/if}

      <div class="form-group">
        <label for="i-name">Nome *</label>
        <input id="i-name" type="text" bind:value={ingForm.Name} placeholder="Ex: Feijão Preto" />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label for="i-unit">Unidade *</label>
          <input id="i-unit" type="text" bind:value={ingForm.Unit} placeholder="kg, g, L, un…" />
        </div>
        <div class="form-group">
          <label for="i-stock">Estoque inicial</label>
          <input id="i-stock" type="number" min="0" step="0.01" bind:value={ingForm.StockQuantity} />
        </div>
      </div>
      <div class="form-group">
        <label for="i-minstock">Estoque mínimo</label>
        <input id="i-minstock" type="number" min="0" step="0.01" bind:value={ingForm.MinStock} />
      </div>

      <div class="modal-actions">
        <button class="btn btn-ghost" onclick={() => (showIngModal = false)}>Cancelar</button>
        <button class="btn btn-primary" onclick={saveIngredient} disabled={ingSaving || !ingForm.Name || !ingForm.Unit}>
          {ingSaving ? 'Salvando…' : 'Criar Ingrediente'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .page-wrapper   { max-width: 1100px; margin: 0 auto; padding: 2rem 1.5rem; }
  .page-header    { display: flex; justify-content: space-between; align-items: flex-start; flex-wrap: wrap; gap: 1rem; margin-bottom: 2.5rem; }
  .page-title     { font-size: 1.75rem; font-weight: 700; color: var(--color-text); margin: 0; }
  .page-subtitle  { font-size: 0.9rem; color: var(--color-muted); margin: 0.25rem 0 0; }
  .header-actions { display: flex; gap: 0.75rem; }

  .section        { margin-bottom: 3rem; }
  .section-title  { font-size: 1.1rem; font-weight: 600; color: var(--color-text); margin: 0 0 1.25rem; display: flex; align-items: center; gap: 0.5rem; }
  .badge          { background: var(--color-surface-2, #f0f0f0); color: var(--color-muted); font-size: 0.75rem; font-weight: 600; padding: 0.1rem 0.5rem; border-radius: 99px; }

  /* Cards de produto */
  .card-grid      { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 1rem; }
  .product-card   { display: block; background: var(--color-surface, #fff); border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; padding: 1.25rem; text-decoration: none; color: inherit; transition: box-shadow 0.15s, border-color 0.15s; }
  .product-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); border-color: var(--color-primary, #e85d04); }
  .product-card-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 0.5rem; margin-bottom: 0.5rem; }
  .product-name   { font-weight: 600; font-size: 1rem; color: var(--color-text); }
  .product-price  { font-weight: 700; color: var(--color-primary, #e85d04); white-space: nowrap; }
  .product-desc   { font-size: 0.85rem; color: var(--color-muted); margin: 0 0 0.75rem; line-height: 1.45; }
  .product-footer { display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }
  .tag            { font-size: 0.75rem; background: var(--color-surface-2, #f5f5f5); padding: 0.15rem 0.5rem; border-radius: 4px; color: var(--color-muted); }
  .status-dot     { font-size: 0.75rem; color: #9ca3af; margin-left: auto; }
  .status-dot.active { color: #16a34a; }

  /* Tabela ingredientes */
  .ing-table-wrapper { overflow-x: auto; border: 1px solid var(--color-border, #e5e7eb); border-radius: 0.75rem; }
  .ing-table      { width: 100%; border-collapse: collapse; font-size: 0.9rem; }
  .ing-table th   { text-align: left; padding: 0.75rem 1rem; background: var(--color-surface-2, #f9fafb); font-weight: 600; color: var(--color-muted); font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.04em; }
  .ing-table td   { padding: 0.75rem 1rem; border-top: 1px solid var(--color-border, #e5e7eb); }
  .ing-table tr:hover td { background: var(--color-surface-2, #f9fafb); }
  .muted          { color: var(--color-muted); }
  .bold           { font-weight: 500; }

  /* Estados */
  .loading-state  { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; color: var(--color-muted); }
  .spinner        { width: 2rem; height: 2rem; border: 3px solid var(--color-border, #e5e7eb); border-top-color: var(--color-primary, #e85d04); border-radius: 50%; animation: spin 0.7s linear infinite; }
  @keyframes spin  { to { transform: rotate(360deg); } }
  .empty-state    { text-align: center; padding: 3rem 1rem; }
  .empty-icon     { font-size: 2.5rem; margin: 0; }
  .empty-text     { color: var(--color-muted); margin: 0.5rem 0 1.25rem; }
  .alert-error    { display: flex; align-items: center; gap: 1rem; background: #fef2f2; border: 1px solid #fca5a5; color: #b91c1c; padding: 1rem 1.25rem; border-radius: 0.5rem; margin-bottom: 1.5rem; }

  /* Modal */
  .modal-overlay  { position: fixed; inset: 0; background: rgba(0,0,0,0.45); display: flex; align-items: center; justify-content: center; z-index: 50; padding: 1rem; }
  .modal          { background: var(--color-surface, #fff); border-radius: 1rem; padding: 2rem; width: 100%; max-width: 480px; max-height: 90vh; overflow-y: auto; box-shadow: 0 20px 60px rgba(0,0,0,0.2); }
  .modal-title    { font-size: 1.25rem; font-weight: 700; margin: 0 0 1.5rem; }
  .modal-actions  { display: flex; justify-content: flex-end; gap: 0.75rem; margin-top: 1.75rem; }

  /* Formulário */
  .form-group     { display: flex; flex-direction: column; gap: 0.4rem; margin-bottom: 1rem; }
  .form-group label { font-size: 0.85rem; font-weight: 500; color: var(--color-text); }
  .form-group input,
  .form-group textarea { border: 1px solid var(--color-border, #d1d5db); border-radius: 0.5rem; padding: 0.6rem 0.75rem; font-size: 0.9rem; background: var(--color-surface, #fff); color: var(--color-text); transition: border-color 0.15s; width: 100%; box-sizing: border-box; }
  .form-group input:focus,
  .form-group textarea:focus { outline: none; border-color: var(--color-primary, #e85d04); box-shadow: 0 0 0 3px rgba(232,93,4,0.12); }
  .form-row       { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; }
  .checkbox-group label { flex-direction: row; align-items: center; gap: 0.5rem; cursor: pointer; font-weight: 400; }
  .form-error     { color: #b91c1c; font-size: 0.85rem; margin-bottom: 1rem; background: #fef2f2; padding: 0.6rem 0.75rem; border-radius: 0.4rem; border: 1px solid #fca5a5; }

  /* Botões */
  .btn            { display: inline-flex; align-items: center; justify-content: center; gap: 0.4rem; padding: 0.55rem 1.1rem; border-radius: 0.5rem; font-size: 0.9rem; font-weight: 600; cursor: pointer; border: none; transition: background 0.15s, opacity 0.15s; text-decoration: none; }
  .btn:disabled   { opacity: 0.55; cursor: not-allowed; }
  .btn-primary    { background: var(--color-primary, #e85d04); color: #fff; }
  .btn-primary:hover:not(:disabled) { background: var(--color-primary-dark, #c84e00); }
  .btn-secondary  { background: var(--color-surface-2, #f3f4f6); color: var(--color-text); border: 1px solid var(--color-border, #d1d5db); }
  .btn-secondary:hover:not(:disabled) { background: var(--color-border, #e5e7eb); }
  .btn-ghost      { background: transparent; color: var(--color-muted); }
  .btn-ghost:hover:not(:disabled) { background: var(--color-surface-2, #f3f4f6); }
  .btn-sm         { padding: 0.35rem 0.75rem; font-size: 0.8rem; }

  @media (max-width: 600px) {
    .page-header  { flex-direction: column; }
    .form-row     { grid-template-columns: 1fr; }
    .header-actions { width: 100%; }
    .header-actions .btn { flex: 1; }
  }
</style>
