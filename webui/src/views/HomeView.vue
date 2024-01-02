<script>
import { VueElement } from 'vue';
import { reactive } from 'vue';
import App from '../App.vue';
import Profile from "../components/Profile.vue"
export default {
	components: {
		Profile	
	},
	data: function() {
		return {
			err: false,
			errMess:null,
			
			loading: false,
			
			dataFormName : "Register your data here",
			inputform: null,
			backAv:false,
			
		}
	},
	methods: {				
		async nullAlerts(){
			this.err = false;			
			this.errMess = null;			
		},
		async errAlert(data){			
			this.err = true;			
			this.errMess = data;
		},		
		async submit(){
			try {
				this.loading=true;
				var r;
				//let response = await this.$axios.get("/users",{username:this.username},headers={Authorization:sessionStorage.getItem("id")});
				await this.$axios({
					method:"put",
					url:"/users/"+sessionStorage.getItem("id")+"/profile",
					data:{
						"username":document.getElementById("username").value,
						"name":document.getElementById("firstName").value,
						"surname":document.getElementById("lastName").value,
						"dateOfBirth":document.getElementById("birthdayDate").value
					},					
					headers:{
						Authorization:sessionStorage.getItem("id")
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
					sessionStorage.setItem("new",false)		
					sessionStorage.setItem("username",this.username)			
					this.hideInputForm()
					break;
				default:
					this.errAlert(r.data);
					break;
			}
			location.reload()
		},		
		async showInputForm(){
			console.log("lo faccio")
			this.inputform = true;
			this.dataFormName = "Update your data";
			this.backAv = true			
		},
		async hideInputForm(){
			this.inputform = false;						
		},		
	},
	mounted() {				
	}
}

			
		
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<div>
				<h1>
					Welcome, <span v-text="username"></span>
				</h1>
			</div>						
		</div>		
		
		<div class="alert alert-danger" role="alert" v-if="err" >
			<h4 class="alert-heading" v-text="errMess"></h4>			
		</div>						
		<!--
			<input type="file" id="fileInput" @change="handleFileSelect()">
	    	<img id="imagePreview" alt="Preview">
		-->		

		<div v-if="!err">			
		<section class="h-100 gradient-custom-2" v-show="!inputform" >
			<Profile @show="showInputForm()" >
				
			</Profile>
		</section>
		<section class="vh-100 gradient-custom"  v-show="inputform" >
			<div class="container py-5 h-100">
				<div class="row justify-content-left align-items-left h-100">
				<div class="col-12 col-lg-9 col-xl-7">
					<div class="card shadow-2-strong card-registration" style="border-radius: 15px;">
					<div class="card-body p-4 p-md-5">
						<h3 class="mb-4 pb-2 pb-md-0 mb-md-5" v-text="dataFormName"></h3>
						<form @submit="submit()">

						<div class="row">
							<div class="col-md-6 mb-4">

							<div class="form-outline">
								<input type="text" id="firstName" class="form-control form-control-lg" pattern="^[a-zA-Z]{3,25}$" 
								v-bind:required="!backAv"
								title="First name must be beetween 3 an 25 char and must not contain special char"/>
								<label class="form-label" for="firstName">First Name</label>
							</div>

							</div>
							<div class="col-md-6 mb-4">

							<div class="form-outline">
								<input type="text" id="lastName" class="form-control form-control-lg" pattern="^[a-zA-Z']{3,25}$" 
								v-bind:required="!backAv"
								title="Last name must be beetween 3 an 25 char and must not contain special char except for '"/>
								<label class="form-label" for="lastName">Last Name</label>
							</div>

							</div>
						</div>

						<div class="row">
							<div class="col-md-6 mb-4 d-flex align-items-center">

								<div class="form-outline datepicker w-100">
									<input type="text" class="form-control form-control-lg" id="birthdayDate" pattern="[0-9]{4}-[0-9]{2}-[0-9]{2}" 
									v-bind:required="!backAv"
									title="Format the data as so yyyy-mm-dd" placeholder="yyyy-mm-dd"/>
									<label for="birthdayDate" class="form-label">Birthday</label>
								</div>
							</div>						
							<div class="col-md-6 mb-4 d-flex align-items-center">
								
								<div class="form-outline">
									<input type="text" id="username" class="form-control form-control-lg" pattern="^[a-zA-Z0-9._]{3,16}$" 
									v-bind:placeholder="username"									
									title="The username must be beetween 3 and 16 char and must not contain special char expcet for . and _ "/>
									<label class="form-label" for="lastName">Username</label>
								</div>
							
							</div>
						</div>

						<div class="row">
							<div class="col-md-6 mb-4 d-flex align-items-center">
								<div class="form-outline">
									<div class="mt-4 pt-2">
										<input class="btn btn-primary btn-lg" type="submit" value="Back" v-show="backAv"
										style="background-color:brown" @click="hideInputForm"/>
									</div>
								</div>
							</div>						
							<div class="col-md-6 mb-4 d-flex align-items-center">
								<div class="form-outline">
									<div class="mt-4 pt-2">
										<input class="btn btn-primary btn-lg" type="submit" value="Submit" style="background-color:green"/>
									</div>
								</div>
							
							</div>
						</div>

						

						</form>
					</div>
					</div>
				</div>
				</div>
			</div>
		</section>		

		</div>
		
	</div>
	
</template>

<style>
</style>
