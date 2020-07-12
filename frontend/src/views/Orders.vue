<template>
  <div id="app">
    <h1>Welcome to orders!</h1>
    <div style="display: flex; width: 500px; padding-bottom: 1vw; right: 0px;">
      <h2 style="margin: inherit; padding-right: 1vw;">Choose date range </h2>
        <DatePicker v-if="startDate" :key="noSearch" :start="startDate" :end="endDate" v-on:childToParent="onDateChange" />
    </div>
    <div style="display:flex;">
      <h2 style="margin: inherit; padding-right: 1vw;">Search </h2>
        <SearchBar  v-on:childToParent="onSearchChange" />
    </div>
    <Table v-if="tableData" :URI="inUseURI" :key="noSearch" :myData="tableData" :config="config" />
  </div>
</template>

<script>
import Table from "./../components/Table.vue";
import DatePicker from "./../components/DatePicker";
import SearchBar from "./../components/SearchBar";

export default {
  name: "Orders",
  components: {
    DatePicker,
    Table,
    SearchBar
  },
  data: () => ({
    tableData: undefined,
    noSearch: 0,
    startDate: "",
    endDate: "",
    page: 0,
    searchString: "",
    baseURI: "http://localhost:8000" ,
    inUseURI:"http://localhost:8000/api/orders/",
    config: [
      {
        key: "order_name",
        title: "Order Name",
        type: "text"
      },
      {
        key: "customer_company_name",
        title: "Company Name",
        type: "text"
      },
      {
        key: "customer_name",
        title: "Customer Name",
        type: "text"
      },
      {
        key: "order_date",
        title: "Order Date",
        type: "date"
      },
      {
        key: "delivered_amount",
        title: "Delivered Amount",
        type: "currency"
      },
      {
        key: "total_amount",
        title: "Total Amount",
        type: "currency"
      }
    ]
  }),
  methods: {
    getData: function() {
      this.$axios.get(this.inUseURI + this.page).then(({ data }) => {
        this.tableData = data;
        if(this.noSearch == 0){
        this.endDate = this.startDate = data[0]["order_date"].split("T")[0]
        }
        this.noSearch += 1;
      });
    },
    onSearchChange: function(value) {
      this.searchString = value;
      this.search();
    },
    onDateChange: function(value) {
      var start = value.start
        .toLocaleDateString()
        .split("/")
        .reverse()
        .join("-");
      var end = value.end
        .toLocaleDateString()
        .split("/")
        .reverse()
        .join("-");

        if(start!= this.startDate || end != this.endDate){
            this.startDate = start
            this.endDate = end
            this.inUseURI = this.baseURI+"/api/orders/between/" + this.startDate +"/"+this.endDate+ "/";
            console.log(this.inUseURI)
            this.getData()
        }  
    },
    search: function() {
      this.page = 0;
      if(this.searchString.trim() !== "" && this.startDate!==this.endDate){
        this.inUseURI = this.baseURI+"/api/orders/search_between/" + this.searchString + "/"+ this.startDate +"/"+this.endDate+ "/";
      }else if (this.searchString.trim() === "" && this.startDate===this.endDate) {
        this.inUseURI = this.baseURI+"/api/orders/";
      } else {
        this.inUseURI = this.baseURI+"api/orders/" + this.searchString + "/";
      }
      this.getData();
    }
  },
  mounted() {
    this.getData();
  }
};
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
