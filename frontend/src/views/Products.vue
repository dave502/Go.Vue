<template>
<div>
  <h1 class="green">{{ products_title }}</h1>
   <button type="button" 
          class="text-green-700 
          hover:text-white border border-green-700 hover:bg-green-800 
          focus:ring-4 focus:outline-none focus:ring-green-300 f
          ont-medium rounded-lg text-sm px-5 py-3 text-center mx-2 mb-2 
          dark:border-green-500 dark:text-green-500 dark:hover:text-white 
          dark:hover:bg-green-600 dark:focus:ring-green-800"
          ref="addButton"
          @click="createOrder">
      {{ create_order_title }}
  </button>   
  <div class="grid grid-cols-2 md:grid-cols-4 place-content-center"> 
    <div v-for="(product, index) in products" :key="index">
      <div class="bg-[#3b3b3b] mr-3 text-sm text-white rounded-md p-2 m-2">
          <div class="product-title">{{ product.productName }}</div>
          <div class="product-price">Цена: {{ product.productPrice}} р.</div>
          <div class="btn-group">
                <input  class="input" type="number" :id="product.productId" name="product_count" min="10" max="100" @change="$emit('input-changed')" />
                <button type="button" 
                        class="text-green-700 
                        hover:text-white border border-green-700 hover:bg-green-800 
                        focus:ring-4 focus:outline-none focus:ring-green-300 f
                        ont-medium rounded-lg text-sm px-5 py-3 text-center me-2 mb-2 
                        dark:border-green-500 dark:text-green-500 dark:hover:text-white 
                        dark:hover:bg-green-600 dark:focus:ring-green-800"
                        ref="addButton"
                        @click="addProductToCard(product.productId, $event)">
                    <svg class="w-3.5 h-3.5 me-2" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 18 21">
                    <path d="M15 12a1 1 0 0 0 .962-.726l2-7A1 1 0 0 0 17 3H3.77L3.175.745A1 1 0 0 0 2.208 0H1a1 1 0 0 0 0 2h.438l.6 2.255v.019l2 7 .746 2.986A3 3 0 1 0 9 17a2.966 2.966 0 0 0-.184-1h2.368c-.118.32-.18.659-.184 1a3 3 0 1 0 3-3H6.78l-.5-2H15Z"/>
                    </svg>
                </button>     
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
  name: "Products",
  data() {
    return {
      products: [],
      products_title: "Каталог",
      create_order_title: "Оформить заказ"
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
    this.getProducts();
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
    getProducts(){
      Query.getAllProducts().
      then(response => {
        console.log(response)
        this.products = response.data
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