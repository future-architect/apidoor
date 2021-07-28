import MockAdapter from "axios-mock-adapter";
import equal from "fast-deep-equal";

import { getProducts } from "@/api/get-products";
import validProductsData from "./__mocks__/get-products-valid.json";
import invalidProductsData from "./__mocks__/get-products-invalid.json";
import instance from "@/api/instance";

describe("Test API Handler", () => {
  const mock = new MockAdapter(instance, { onNoMatch: "throwException" });

  afterEach(() => {
    mock.reset();
  });

  describe("test get-products", () => {
    it("should succeed", async () => {
      mock.onGet("/products").reply(200, validProductsData);
      const data = await getProducts();
      expect(equal(data, validProductsData)).toBeTruthy();
    });
    it("should return empty list of products", async () => {
      mock.onGet("/products").reply(200, invalidProductsData);
      const data = await getProducts();
      expect(equal(data, { products: [] })).toBeTruthy();
    });
  });
});
