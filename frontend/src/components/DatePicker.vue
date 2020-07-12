<template>
  <v-date-picker
    id="datePicker"
    v-model="range"
    :mode="mode"
    :input-debounce="500"
    :key="emitToParent()"
    @change="emitToParent()"
  ></v-date-picker>
</template>

<script>
import Vue from "vue";
import DatePicker from "v-calendar/lib/components/date-picker.umd";
Vue.component("v-date-picker", DatePicker);
export default {
  name: "DatePicker",
  key: 0,
  props: ["start", "end"],
  data: function() {
    return {
      mode: "range",
      range: {
        start: new Date(this.start.split('-')[0],this.start.split('-')[1]-1,this.start.split('-')[2]), // Jan 16th, 2018
        end: new Date(this.end.split('-')[0],this.end.split('-')[1]-1,this.end.split('-')[2]) // Jan 19th, 2018
      }
    };
  },
  methods: {
    emitToParent: function() {
      this.$emit("childToParent", this.range);
    }
  }
};
</script>


<style>
</style>
