<template>
  <div class="container">
    <div class="row mt-5 mb-5">
      <div class="col-2">
        <h1>Search</h1>
      </div>
      <div class="col-10 pt-2">
        <input
          class="form-control mr-sm-2"
          type="search"
          placeholder="Search"
          aria-label="Search"
          v-model="keywords"
          @input="fetchOrders"
        />
      </div>
    </div>

    <div class="col-5 ml-0 pl-0 mt-5 mb-5">
      <h5>Created Date</h5>
      <VueDatePicker v-model="date" :range="true" :enable-time-picker="false" :model-value="date" @update:model-value="fetchOrders()" @cleared="clearedDate()"></VueDatePicker>
    </div>

    <div class="col-5 ml-0 pl-0 mt-5 mb-5">
      <h5>Total Amount: ${{totalAmount()}} </h5>
      
    </div>

    <div class="mt-5 mb-5">
      <table class="table">
        <thead>
          <tr>
            <th scope="col">Order Name</th>
            <th scope="col">Customer Company</th>
            <th scope="col">Customer Name</th>
            <th scope="col">Order Date <img @click.prevent="fetchOrders(), this.sort = !this.sort" class="ml-2" src="@/assets/sort.png" alt="Sort Icon" width="10" height="10" /></th>
            <th scope="col">Delivered Amount</th>
            <th scope="col">Total Amount</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in listData">
            <td>{{ item.order_name }}</td>
            <td>{{ item.customer_company_name }}</td>
            <td>{{ item.customer_name }}</td>
            <td>{{ formatDate(item.order_date)}}</td>
            
            <td>${{ item.delivered_amount }}</td>
            <td>${{ item.total_amount }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="row">
      <div class="mx-auto">
        <ul class="pagination">

          <li v-for="n in totalPages" :class="`${n == currentPage ? 'active page-item' : 'page-item'}`"><a class="page-link" href="#" @click="fetchOrders(n)">{{n}}</a></li>
         

        </ul>
      </div>
    </div>

  </div>
</template>
<script>
import VueDatePicker from "@vuepic/vue-datepicker";
import "@vuepic/vue-datepicker/dist/main.css";
import moment from "moment";
import axios from "axios";
import { ref } from 'vue';


export default {
  components: { VueDatePicker },
    data() {
    return {
        date: null,
        listData: [],
        keywords: "",
        startDate: "",
        endDate: "",
        totalAmountInPage: null,
        totalPages: 0,
        currentPage: 1,
        sort: false
    };
    },
    mounted() {
    this.fetchOrders();
    },
    methods: {
    async fetchOrders(page) {
        try {
        if (this.date) {
            this.startDate =  moment(this.date[0]).format('YYYY-MM-DD');
            this.endDate =  moment(this.date[1]).format('YYYY-MM-DD');
        }    
        const response = await axios.get('http://localhost:8080/orders', {
          params: { search: this.keywords, startDate: this.startDate, endDate: this.endDate, page: page || 1, sortDirection: this.sort ? "DESC" : "ASC" }
        });
        this.listData = response.data.orders; 
        this.totalPages = response.data.totalPages
        this.currentPage = response.data.currentPage
        } catch (error) {
        console.error('Error fetching data:', error);
        }
    },
    formatDate(originalDate) {
      return moment(originalDate).format("MMMM Do, h:mm a");
    },
    clearedDate() {
        console.log('cleared >>>>>')
        this.startDate = null
        this.endDate = null
    this.fetchOrders();

    },
    totalAmount(){
        this.totalAmountInPage = 0
        if (this.listData) {
            return this.listData.reduce((sum, item) => sum + item.total_amount, 0).toFixed(2);
        } else {
            return 0
        }
    }
    },
};
</script>