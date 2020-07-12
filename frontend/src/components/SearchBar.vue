<template>
    <div>
        <input @keyup="emitToParent()" v-model="searchString" id="searchString" type="text" placeholder="Search item or order name.."/>
    </div>
</template>

<script>
export default {
    data() {
    return {
      searchString: ''
    }
  },
    methods: {
            emitToParent: function() {
                console.log(this.searchString)
            this.$emit('childToParent', this.searchString)
        }
    }
}
</script>

<style>
    #searchString{
        height: 34px;
        border: 1px solid #CBD5E0;
        background: #fff;
        text-align: center;
        width: 80vw;
        color: #2D3848;
        border-radius: 5px;
        margin: 0 0 0 0 ;
    }
</style>