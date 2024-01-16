<script>
import Profile from "../components/Profile.vue"
import SearchBar from "../components/SearchBar.vue"
export default {
	components: {
		Profile,
		SearchBar
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
            err : false,
            errMess:null,
            username:null,
			searched:null,
			profileActive:false,
		}
	},
    watch: {
        $route(to, from) {            
            if (to.params.username !== from.params.username) {                
                this.load()
            }
        }
    },
	methods: {		
		async load() {                        
            this.searched = window.location.href.split("/").pop()
			
			if(this.searched === "search"){
				this.profileActive =false
				return
			}
			else{ 
				this.profileActive =true
			}

            
            var id = sessionStorage.getItem("id")
            var r 
			try {
				this.loading=true;								
				await this.$axios({
					method:"get",
					url:"/users",
					params:{
						username:this.searched
					},
					headers:{
						Authorization:id
					}
				}).then((response)=>{
					r = response}                    
					)				
				
			} catch (e) {                
				r = e.response;							                
			}            
			this.loading=false;            
			switch (r.status) {
				case 200:	
				case 201:	
					this.searched = r.data["name"]
					break;					
				default:
					this.errAlert(r.data);
					break;				
			}	
		},
		async errAlert(data){			
			this.err = true;
			this.errMess = data;
		},		
		async nullAlerts(){
			this.err = false;			
			this.errMess = null;			
		},
	},
	mounted() {
        this.load()        
	}
}
</script>

<template>
		
		<SearchBar></SearchBar>
		<div class="alert alert-danger" role="alert" v-if="err" >
			<h4 class="alert-heading" v-text="errMess"></h4>			
		</div>
		<section class="h-100 gradient-custom-2" v-if="profileActive">
			<Profile></Profile>
		</section>
	
</template>

<style>
</style>







