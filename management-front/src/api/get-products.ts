import instance from "@/api/instance";
import { Products } from "@/models/products";
import { isProducts } from "@/lib/type-guard";

// FIXME: modify api path
// get information of products
export const getProducts = async (): Promise<Products> => {
  const response = await instance.get("/products");

  if (isProducts(response.data)) {
    return response.data;
  }

  return {
    products: [],
  };
};
