import axios from "axios";
import MockAdapter from "axios-mock-adapter";

import { getProducts } from "@/api/get-products";
import productsData from "./__mocks__/products.json";

describe("Test API Handler", () => {
  const mock = new MockAdapter(axios);

  afterEach(() => {
    mock.reset();
  });

  describe("test get-products", () => {
    it("should succeed", async () => {
      mock.onGet("/products").reply(200, productsData);

      const data = await getProducts();

      expect(data.products[0].id).toBe(3);
    });
  });
});
