<script>
import { VueElement } from 'vue';

import SearchBar from "../components/SearchBar.vue"
export default {
	components: {		
		SearchBar
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,			
			searched: null,
		}
	},
	methods: {		
		async search() {			
			this.$router.push("/users/search/"+this.searched)	
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

	},
	mounted() {
        		
	}
}
</script>

<template>	
	<SearchBar></SearchBar>			
	
</template>

<style>
</style>







