<template>
    <div>
        <div class="search-wrapper">
            <label>Search item: </label>
            <input id="searchString" type="text" placeholder="Search item.."/>
            <button @click="search()"> search </button>
        </div>
        <div id="table" style="height: 30vh">
            <table style="width:100%">
                <thead>
                    <tr>
                        <th v-for ="(obj, ind) in config" :key="ind" > {{obj.title}} </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for ="(row, index) in tableData" :key="index">
                        <td v-for ="(obj, ind) in config" :key="ind">
                            <span v-if="obj.type === 'text'"> {{row[obj.key]}} </span>
                            <span v-if="obj.type === 'date'"> 
                                {{new Date(row[obj.key]).toLocaleDateString('en-US',{month: 'short', day: 'numeric' })
                                +", "+
                                new Date(row[obj.key]).toLocaleTimeString('en-AU', { hour: '2-digit', minute: '2-digit' })}} 
                            </span>
                            <span v-if="obj.type === 'currency'">
                                {{"$"+(parseFloat(row[obj.key])).toFixed(2)}} 
                            </span>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <button id="prev" :disabled=disablePrev @click="prev()">Previous</button>
        <button id="next" :disabled=disableNext @click="next()">Next</button>
    </div>
</template>

<script>

export default {
    props: ['myData','config'],
    data: function () {
        return{
            page: 0,
            tableData: this.myData,
            disableNext: false,
            disablePrev: true,
            inUseURI: "http://localhost:8000/api/orders/",
        }
    },

    methods: {
        search: function(){
            this.page = 0
            var searchString = document.getElementById("searchString").value
            if(searchString.trim() === ""){
                this.inUseURI = "http://localhost:8000/api/orders/"
            }else{
                this.inUseURI = "http://localhost:8000/api/orders/"+searchString+"/"
            }
            this.$axios.get(this.inUseURI+this.page)
                .then(({data}) => {
                    this.tableData = data
                })
        },
        next: function(){
        this.$axios.get(this.inUseURI+this.page+1)
            .then(({data}) => {
                if(data){
                    this.tableData = data
                    this.page += 1
                    if( data.length < 5)
                        this.disableNext = true
                    this.disablePrev = false
                }else{
                    this.disableNext = false
                }
            })
        },
        prev: function(){
            if(this.page > 0)
        this.$axios.get(this.inUseURI+(this.page-1))
            .then(({data}) => {
                if( data && this.page > 0){
                    this.tableData = data
                    this.page -= 1
                    if(this.page == 0)
                        this.disablePrev = true
                    this.disableNext = false
                    this.total = data.reduce((a, b) => a["total_amount"] + b["total_amount"], 0)
                }else{
                    this.disablePrev = false
                }
            })
        }
    }
}


</script>

<style lang="css" scoped>


</style>