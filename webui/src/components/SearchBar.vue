<script>
import { VueElement } from 'vue';

export default {
    name : "SearchBar",
	data: function() {
		return {
			errormsg: null,
			loading: false,			
			bar_text: null,
		}
	},
    watch:{
        $route(to, from){            
            location.reload()
        }
    },
	methods: {		
		async search() {			
			this.$router.push("/users/search/"+this.bar_text)	
            sessionStorage.setItem("username",this.bar_text)                                    
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
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">					
			<h1>
				Discover new users
			</h1>										
        </div>				

		<div class="input-group mb-3">
			<input type="text" class="form-control" placeholder="Search a user" aria-label="Recipient's username" aria-describedby="basic-addon2" v-model="bar_text">
			<div class="input-group-append">
				<button class="btn btn-outline-secondary" type="button" @click="search">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
				</button>
			</div>
		</div>
	</div>
	
</template>

<style>
</style>







