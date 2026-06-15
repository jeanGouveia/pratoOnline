export interface Ingredient {
  ID: number;
  Name: string;
  Unit: string;
  Quantity?: number; // quantidade usada no produto
}

export interface Product {
  ID: number;
  Name: string;
  Description?: string;
  Price: number;
  IsComposto: boolean;
  Active: boolean;
  CreatedAt?: string;
  UpdatedAt?: string;
  Ingredients?: Ingredient[];
}

export interface ProductCreatePayload {
  name: string;
  description?: string;
  price: number;
  is_composto: boolean;
  active: boolean;
}

export interface ProductIngredientPayload {
  ingredient_id: number;
  quantity: number;
}
