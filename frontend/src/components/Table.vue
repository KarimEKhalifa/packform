<template>
    <div>
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
        <button id="prev" class="controls" :disabled=disablePrev @click="prev()">&#8249;</button>
        <button id="next" class="controls" :disabled=disableNext @click="next()">&#8250;</button>
    </div>
</template>

<script>

export default {
    name: 'Table',
    props: ['myData','config', 'URI'],
    data: function () {
        return{
            page: 0,
            tableData: this.myData,
            inUseURI: this.URI,
            disableNext: false,
            disablePrev: true,
        }
    },

    methods: {
        next: function(){
        this.$axios.get(this.inUseURI+this.page+1)
            .then(({data}) => {
                if(data.length > 0){
                    console.log(data)
                    this.tableData = data
                    this.page += 1
                    if( data.length < 5)
                        this.disableNext = true
                    this.disablePrev = false
                }else{
                    this.disableNext = true
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

.controls {
  text-decoration: none;
  display: inline-block;
  padding: 8px 16px;
}

.controls:hover {
  background-color: #ddd;
  color: black;
}


</style>