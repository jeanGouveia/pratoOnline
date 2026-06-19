export interface Ingredient {
  ID: number;
  Name: string;
  Unit: string;
  StockQuantity: number;
  MinStock: number;
  Active: boolean;
  CreatedAt?: string;
  UpdatedAt?: string;
}

export interface IngredientCreatePayload {
  name: string;
  unit: string;
  stock_quantity: number;
  min_stock: number;
}
