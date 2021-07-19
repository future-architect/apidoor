import instance from '@/api/instance';
import { Products } from '@/models/products';
import { isProducts } from '@/lib/type-guard';

export const getInfo = async (): Promise<Products> => {
    const data = await instance.get('/products');

    if(isProducts(data)) {
        return data;
    };

    return {
        products: []
    };
};