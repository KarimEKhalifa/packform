<template>
  <div id="app">
    <h1>Welcome to orders!</h1>
    <router-link to="/"> Back</router-link>
    <Table v-if="tableData" :myData ="tableData" :config="config"/>
  </div>
</template>

<script>
import Table from './../components/Table.vue'

export default {
  name: 'App',
  components: {
    Table
  },
  data: () =>({
    tableData: undefined,
    config: [
      {
        key: 'order_name',
        title: 'Order Name',
        type: 'text'
      },
      {
        key: 'customer_company_name',
        title: 'Company Name',
        type: 'text'
      },
      {
        key: 'customer_name',
        title: 'Customer Name',
        type: 'text'
      },
      {
        key: 'order_date',
        title: 'Order Date',
        type: 'date'
      },
      {
        key: 'delivered_amount',
        title: 'Delivered Amount',
        type: 'currency'
      },
      {
        key: 'total_amount',
        title: 'Total Amount',
        type: 'currency'
      }
    ]
  }),
  mounted(){
    this.$axios.get("http://localhost:8000/api/orders/0")
    .then(({data}) => {
      this.tableData = data
    })
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
