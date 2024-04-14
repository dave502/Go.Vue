import api from '@/lib/api';

export function getAllOrders() {
    return api.get(`/order` );
}

export function getAllOrdersIds() {
    return api.get(`/order-ids` );
}

export function postPackingList(data) {
    console.log("postPackingList", data);
    const payload = JSON.stringify({'ids': data});
    console.log("postPackingList payload", payload);
    return api.post(`/packing-list`, payload );
}

export function getAllProducts() {
    return api.get(`/product` );
}