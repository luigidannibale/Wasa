<script>
import { stringifyQuery } from 'vue-router';

export default {
	data: function() {
		return {	
			err:false,
			errMess:null,
			loading: false,
			username: null,
		}
	},
	methods: {		
        async doLogin() {            
			var r = null
			this.nullAlerts()
			try {
				this.loading=true;
				let response = await this.$axios.post("/users",{username:this.username});										
				r = response;				
			} catch (e) {
				r = e.response;			
			}
			this.loading=false;
			switch (r.status) {
				case 200:	
				case 201:														
					this.loginUpdate(r.data.split(":")[1].toString().trim());					
					sessionStorage.setItem("new",r.status === 200 ? false : true)
					break;					
				case 400:				
				case 500:
					this.errAlert(r.data);
					break;
				default:
					break;
			}				
			
		},
		async loginUpdate(id){			
			this.$router.push("/home/"+id)		
			sessionStorage.setItem("username",this.username)	
			sessionStorage.setItem("logged",true)
			sessionStorage.setItem("id",id)
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
			this.msg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.msg = e.toString();
			}
			this.loading = false;
		},
		async reload(){			
			if ((sessionStorage.length != 0)) {
				sessionStorage.clear()
				location.reload();
			}					
		}
	},
	mounted() {				
		this.reload()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Login</h1>
			<!--
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
			-->
		</div>        		    
		<div class="alert alert-danger" role="alert" v-if="err" >
			<h4 class="alert-heading" v-text="errMess"></h4>			
		</div>
		<div class="alert alert-success" role="alert" v-if="succ">
			<h4 class="alert-heading" v-text="succMess"></h4>			
		</div>
        <div>
            <label for="username">Username</label>   <br>
            <input type="text" placeholder="Insert your username" v-model="username">
            <button type="button" name="loginButton" @click="doLogin">Log in</button>
        </div>

	</div>
</template>

<style>
</style>
