<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
            err : false,
            errMess:null,
            username:null,
		}
	},
    watch: {
        $route(to, from) {
            // Check if the username has changed
            if (to.params.username !== from.params.username) {
                console.log("aggiorno")
                this.load()
            }
        }
    },
	methods: {		
		async load() {                        
            this.username = window.location.href.split("/").pop()
            this.nullAlerts()
            var id = sessionStorage.getItem("id")
            var r 
			try {
				this.loading=true;								
				await this.$axios({
					method:"get",
					url:"/users",
					params:{
						username:this.username
					},
					headers:{
						Authorization:id
					}
				}).then((response)=>{
					r = response}                    
					)				
				
			} catch (e) {                
				r = e.response;							
                console.log("arrivo",e)
			}            
			this.loading=false;
            console.log(r)
			switch (r.status) {
				case 200:	
				case 201:	

					break;					
				default:
					this.errAlert(r.data);
					break;				
			}	

		},
        async search() {			            
			this.$router.push("/users/search/"+this.username)	
            console.log(this.username)              
            //location.replace("#/users/search/"+this.username)
		},
        async errAlert(data){			
			this.err = true;			
			this.errMess = data;
		},
		async nullAlerts(){
			this.err = false;			
			this.errMess = null;			
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
        this.load()
        
	}
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">		
			<div>
				<h1>
					Discover new users
				</h1>
			</div>								
        </div>	
            <div class="input-group mb-3">
				<input type="text" class="form-control" placeholder="Search a user" aria-label="Recipient's username" aria-describedby="basic-addon2" 
                v-model="username">
				<div class="input-group-append">
					<button class="btn btn-outline-secondary" type="button" @click="search">
                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
                    </button>
				</div>
			</div>
            <div class="alert alert-danger" role="alert" v-if="err" >
                <h4 class="alert-heading" v-text="errMess"></h4>			
            </div>
		</div>
	
</template>

<style>
</style>







