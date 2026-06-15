import { request } from './client';
import type { Product, ProductCreatePayload, ProductIngredientPayload } from '$lib/types/product';
import type { Ingredient, IngredientCreatePayload } from '$lib/types/ingredient';

// ── Produtos ────────────────────────────────────────────────────────────────

export async function getProducts(): Promise<Product[]> {
  const res = await request<Product[]>('/products');
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getProduct(id: number): Promise<Product> {
  const res = await request<Product>(`/products/${id}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function createProduct(payload: ProductCreatePayload): Promise<Product> {
  const res = await request<Product>('/products', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getProductIngredients(id: number): Promise<Ingredient[]> {
  const res = await request<Ingredient[]>(`/products/${id}/ingredients`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function updateProductIngredients(
  id: number,
  ingredients: ProductIngredientPayload[]
): Promise<void> {
  const res = await request<{ message: string }>(`/products/${id}/ingredients`, {
    method: 'PUT',
    body: JSON.stringify({ items: ingredients })
  });
  if (res.error) throw new Error(res.error);
}

// ── Ingredientes ─────────────────────────────────────────────────────────────

export async function getIngredients(): Promise<Ingredient[]> {
  const res = await request<Ingredient[]>('/ingredients');
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function createIngredient(payload: IngredientCreatePayload): Promise<Ingredient> {
  const res = await request<Ingredient>('/ingredients', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}
