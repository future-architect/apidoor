import { Product } from "@/models/product";
import { Products } from "@/models/products";

export const isAPIInfo = (arg: unknown): arg is Product => {
  const tmp = arg as Product;

  return (
    typeof tmp?.id === "number" &&
    typeof tmp?.name === "string" &&
    typeof tmp?.source === "string" &&
    typeof tmp?.description === "string" &&
    typeof tmp?.thumbnail === "string"
  );
};

export const isProducts = (arg: unknown): arg is Products => {
  const tmp = arg as Products;

  if (!Array.isArray(tmp?.products)) return false;

  return tmp.products.every((p) => isAPIInfo(p));
};
