<template>
<div>
  <h1 class="green">{{ orders_title }}</h1>
    <button type="button" 
            class="text-green-700 
            hover:text-white border border-green-700 hover:bg-green-800 
            focus:ring-4 focus:outline-none focus:ring-green-300 f
            ont-medium rounded-lg text-sm px-5 py-3 text-center mx-2 mb-2 
            dark:border-green-500 dark:text-green-500 dark:hover:text-white 
            dark:hover:bg-green-600 dark:focus:ring-green-800"
            ref="addButton"
            @click="redirectToPackingList">
        {{ go_to_packing_list_title }}
    </button>  
  <div class="grid grid-cols-2 md:grid-cols-1 place-content-center"> 
    <div v-for="(order, index) in orders" :key="index">
      <div class="bg-[#3b3b3b] mr-3 text-sm text-white rounded-md p-2 m-2">
        <div class="product-title">Заказ № {{ index }}</div>
          <div class="relative overflow-x-auto">
            <table class="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
              <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
                <tr>
                  <th scope="col" class="px-6 py-3">№</th>
                  <th scope="col" class="px-6 py-3">Товар</th>
                  <th scope="col" class="px-6 py-3">Id</th>
                  <th scope="col" class="px-6 py-3">Количество</th>
                  <th scope="col" class="px-6 py-3">Изм.</th>
                  <th scope="col" class="px-6 py-3">Цена</th>
                  <th scope="col" class="px-6 py-3">Валюта</th>
                </tr>
              </thead>
              <tbody>
                <tr  v-for="(product, index) in order" :key="index" class="bg-white border-b dark:bg-gray-800 dark:border-gray-700">
                  <td class="px-6 py-4"> {{index + 1}} </td>
                  <th scope="row" class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"> {{ product.productName }}</th>
                  <td class="px-6 py-4"> {{ product.productId }}</td>
                  <td class="px-6 py-4 dark:text-white"> {{ product.quantity }} </td>
                  <td class="px-6 py-4"> шт. </td>
                  <td class="px-6 py-4 dark:text-white"> {{ product.productPrice }} </td>
                  <td class="px-6 py-4"> руб. </td>
                </tr>
              </tbody>
            </table> 
          </div> 
      </div>
    </div>
  </div>
</div>

  <!-- <div class="stack-small" v-if="!isEditing">

    <p class="main-prop">{{name}}</p>
    <p class="add-prop">Цена: {{price}} руб.</p>
    <label :for="shelve"  class="input-label">Стеллажи</label>
    <input  class="input" type="text" :id="shelve" name="product_shelve" min="10" max="100" @change="$emit('input-changed')" />
    <label :for="count"  class="input-label">Количество</label>
    <input  class="input" type="number" :id="id" name="product_count" min="10" max="100" @change="$emit('input-changed')" />
    <div class="btn-group">
      <button
          type="button"
          class="btn"
          ref="addButton"
          @click="addProductToOrder">
          Добавить в заказ
          <span class="visually-hidden">{{ name }}</span>
      </button>
      <button type="button" class="btn btn__danger" @click="deleteProductFromOrder">
          Удалить
          <span class="visually-hidden">{{ name }}</span>
      </button>
    </div>

  </div> -->
</template>
<script>
import { ref } from 'vue';
import Cookies from 'js-cookie'
import router from '@/router'
import * as Query from '@/lib/query';

// async function getData() {
//   const { data } = await demoAPI.getFromServer()
//   console.log(data)
//   result.value.push(data.message)
// }
const products = ref([])
const msg = ref("")

// defineProps({
//   sampleProp: String,
// });
  
export default {
  //   components: {
  //   // ToDoItemEditForm,
  // },
  name: "Orders",
  data() {
    return {
      orders: [],
      orders_title: "Заказы",
      go_to_packing_list_title: "Сборочный лист"
      //create_order_title: "Оформить заказ"
    };
  },
  // props: {
  //   name: { required: true, type: String },
  //   price: { default: 0, type: Number },
  //   count: { default: 0, type: Number },
  //   shelve: { required: true, type: String },
  //   id: { required: true, type: String },
  // },
  created() {
    this.getOrders();
  //   // Simple POST request with a JSON body using axios
  //   const article = { title: "Vue POST Request Example" };
  //   Query.getAllProducts().
  //   then(response => {
  //     console.log(response)
  //     response.value.push(data.Message)
  //   })
  //   .catch(error => {
  //     this.errorMessage = error.message;
  //     console.error("There was an error!", error);
  //   });
    // axios.post("https://reqres.in/api/articles", article)
      // .then(response => this.articleId = response.data.id);
  },
  computed: {
    count() {
      return this.count;
    },
    shelve() {
      return this.shelve;
    },
  },
  methods: {
    getOrders(){
      Query.getAllOrders().
      then(response => {
        console.log("response", response)
        console.log("groupBy", Object.groupBy(response.data, ({ orderId }) => orderId))
        this.orders = Object.groupBy(response.data, ({ orderId }) => orderId);
        console.log(this.orders)
        // response.value.push(data.Message)
      })
      .catch(error => {
        this.errorMessage = error.message;
        console.error("Failed to run query!", error.message);
      });
    },
    addProductToCard(productID){
      console.log("productID", productID)
      const elInput = document.getElementById(productID)
      const quantity = parseInt(elInput.value);
      elInput.value = undefined
      if (!quantity) return;
      
       console.log("quantity", document.getElementById(productID).value)
        var order = null
        const jsonOrder = Cookies.get('order');
        console.log(jsonOrder)
        if (jsonOrder) {
          order = JSON.parse(jsonOrder);
        } else {
          order = {}
        }
        order[productID] = (order[productID] || 0) + quantity;
        Cookies.set('order' , JSON.stringify(order), { expires: 7 }) 
    },
    createOrder() {
      this.$emit("order-creates");
    },
    deleteProductFromOrder() {
      this.$emit("item-deleted");
    },
    toggleToItemEditForm() {
      console.log(this.$refs.editButton);
      this.isEditing = true;
    },
    itemEdited(newLabel) {
      this.$emit("item-edited", newLabel);
      this.isEditing = false;
      this.focusOnEditButton();
    },
    editCancelled() {
      this.isEditing = false;
      this.focusOnEditButton();
    },
    focusOnEditButton() {
      this.$nextTick(() => {
        const editButtonRef = this.$refs.editButton;
        editButtonRef.focus();
      });
    },
    redirectToPackingList(){
      router.push({ path: '/orders-products'})
    }
  },
  // beforeMount() {
  //   //this.getProducts();
  // },  
};
</script>
<style scoped>

h1 {
  font-weight: 800;
  font-size: 2rem;
  /* color: hsl(147, 67%, 41%); */
}

.main-prop {
  font-family: Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-weight: 800;
  font-size: 16px;
  font-size: 1.5rem;
  line-height: 1.25;
  color: #055529;
  display: block;
  margin-bottom: 5px;
}

.add-prop {
  font-family: Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-weight: 400;
  font-size: 16px;
  font-size: 1.5rem;
  line-height: 1.25;
  color: #055529;
  display: block;
  margin-bottom: 5px;
}



.input-label {
  font-family: Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-weight: 400;
  font-size: 16px;
  font-size: 1.2rem;
  line-height: 1.25;
  color: #055529;
  display: block;
  margin-bottom: 5px;
}
.input {
  font-family: Arial, sans-serif;
  color:black;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.25;
  width: 5rem;
  height: 2.5rem;
  margin-top: 0;
  padding: 5px;
  /* border: 2px solid #0b0c0c; */
  border-radius: 0;
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  display: inline-block;
  margin-bottom: 5px;
  margin-right: 5px;
}


.product-title {
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 10px;
}

.product-price {
  font-weight: 400;
  font-size: 1rem;
  margin-bottom: 10px;
}
</style>