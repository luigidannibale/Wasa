<script>

import { VueElement } from 'vue'
import { reactive } from 'vue'
import App from '../App.vue'
import Profile from "../components/Profile.vue"

export default {
	components: {
		Profile
	},
	data: function() {
		return {
			username : null,
			name:null,
			surname:null,
			birthday:null,
			err: false,
			errMess:null,			
			loading: false,			
			imagePreview:false,
			imageFile:null,
			dataFormName : "Register your data here",
			inputform: null,
			postPhotoForm : null,
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
					if(this.username) {
						sessionStorage.setItem("username",this.username)
					}					
					this.hideInputForm()					
					break;
				default:
					this.errAlert(r.data);
					break;
			}
			location.reload()
		},
		async showInputForm(){			
			this.inputform = true;
			this.dataFormName = "Update your data";						
			this.backAv = !sessionStorage.getItem("new")			
			var r = null
			var id = sessionStorage.getItem("id")			
			try {			
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
			}
			switch (r.status) {
				case 200:										
					this.name = r.data["name"] 
					this.surname = r.data["surname"]	
					this.birthday = r.data["dateOfBirth"]																	
					break;												
				default:
					this.errAlert(r.data);
					break;
			}

		},
		async postPhoto(){
			//this.hidePostPhotoForm()
			let cap = document.getElementById("caption").value			
			if(cap == ""){				
				return 
			}
			var r = null
			var id = sessionStorage.getItem("id")	
			var rbody = new FormData();
			rbody.append("image",this.imageFile)			
			try {			
				await this.$axios({
					method:"post",
					url:"/photos",					
					headers:{
						Authorization:id,
						'Content-Type': 'image/png',
					},		
					params:{
						caption:cap
					},
					data: rbody,														
				}).then((response)=>{
					r = response}
					)								
			} catch (e) {
				r = e.response;							
			}			
			switch (r.status) {
				case 201:										
					this.hidePostPhotoForm()
					break;												
				default:
					this.errAlert(r.data);
					break;
			}
		},		
		async showPostPhotoForm(){			
			this.postPhotoForm = true;
		},
		hidePostPhotoForm(){
			this.postPhotoForm = false;		
			location.reload()				
		},
		hideInputForm(){
			this.inputform = false;						
		},	
		log(){			
			this.username = sessionStorage.getItem("loggedUsername")
			sessionStorage.setItem("username",this.username)				
			if (sessionStorage.getItem("new") == "true")
				this.showInputForm()
			

		},
		handleFileUpload(event) {
			this.imageFile = event.target.files[0];
            if (this.imageFile && this.imageFile['type'].startsWith('image')) {
                const reader = new FileReader();
                reader.onload = (e) => {
                    this.imagePreview = e.target.result;
                };
                reader.readAsDataURL(this.imageFile);
            } else {
                event.target.value = "";
                alert('The given file is not an image');
            }
			
		}
	},
	mounted() {					
		this.log()
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

		<div v-if="!err">			
		<section class="h-100 gradient-custom-2" v-show="!inputform && !postPhotoForm" >			
			<section class="row">
				<div class="col col-lg-2" >					
					<button id="editProfile" type="button" class="btn btn-outline-dark" data-mdb-ripple-color="dark" style="z-index: 1; margin-right: 5px; margin-left: 5px;" @click="showInputForm()" >
						Edit profile
					</button> 	
				</div>	
				<div class="col col-lg-2" >					
					<button id="postPhoto" type="button" class="btn btn-outline-dark" data-mdb-ripple-color="dark" style="z-index: 1; margin-right: 5px; margin-left: 5px;" @click="showPostPhotoForm()" >
						Post Photo
					</button>
				</div>
			</section>				
			<Profile> 
			</Profile>
		</section>
		<section class="vh-100 gradient-custom"  v-if="inputform" >
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
								v-bind:required="!backAv" :value="name"
								title="First name must be beetween 3 an 25 char and must not contain special char"/>
								<label class="form-label" for="firstName">First Name</label>
							</div>

							</div>
							<div class="col-md-6 mb-4">

							<div class="form-outline">
								<input type="text" id="lastName" class="form-control form-control-lg" pattern="^[a-zA-Z']{3,25}$" 
								v-bind:required="!backAv" :value="surname"
								title="Last name must be beetween 3 an 25 char and must not contain special char except for '"/>
								<label class="form-label" for="lastName">Last Name</label>
							</div>

							</div>
						</div>

						<div class="row">
							<div class="col-md-6 mb-4 d-flex align-items-center">

								<div class="form-outline datepicker w-100">
									<input type="text" class="form-control form-control-lg" id="birthdayDate" pattern="[0-9]{4}-[0-9]{2}-[0-9]{2}" 
									v-bind:required="!backAv" :value="birthday"
									title="Format the data as so yyyy-mm-dd" placeholder="yyyy-mm-dd"/>
									<label for="birthdayDate" class="form-label">Birthday</label>
								</div>
							</div>						
							<div class="col-md-6 mb-4 d-flex align-items-center">
								
								<div class="form-outline">
									<input type="text" id="username" class="form-control form-control-lg" pattern="^[a-zA-Z0-9._]{3,16}$" 
									:value = username title="The username must be beetween 3 and 16 char and must not contain special char expcet for . and _ "/>
									<label class="form-label" for="lastName">Username</label>
								</div>
							
							</div>
						</div>

						<div class="row">
							<div class="col-md-6 mb-4 d-flex align-items-center">
								<div class="form-outline">
									<div class="mt-4 pt-2">
										<input class="btn btn-primary btn-lg" type="submit" value="Back" v-if="backAv"
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

		<section class="vh-100 gradient-custom"  v-show="postPhotoForm" >
			<div class="container py-5 h-100">
				<div class="row justify-content-left align-items-left h-100">
				<div class="col-12 col-lg-9 col-xl-7">
					<div class="card shadow-2-strong card-registration" style="border-radius: 15px;">
					<div class="card-body p-4 p-md-5">
						<h3 class="mb-4 pb-2 pb-md-0 mb-md-5">Post your photo</h3>
						<form>							
							<div class="row">
								<div class="col-md-6 mb-4">
									<div class="form-outline">																		
										<label class="form-label" for="firstName"> Caption </label><br>
										<textarea name="caption" id="caption" cols="20" rows="2" required title="Insert a caption please" 
										oninvalid="this.setCustomValidity('Insert a caption please')"></textarea>
									</div>
								</div>
								<div class="col-md-6 mb-4">
									<div class="form-outline">								
										<label class="form-label" for="firstName"> Upload image </label>
										<input type="file" @change="handleFileUpload">										
									</div>
								</div>							
							</div>													
							<div class ="row" v-if="imagePreview">
								<span> Preview </span>
								<img :src="imagePreview" alt="Preview" style="height: 300px; width: fit-content;">
							</div>							
							<div class="row">								
								<div class="col-md-4 mb-4 d-flex align-items-center">
									<div class="form-outline">
										<div class="mt-4 pt-2">
											<input class="btn btn-primary btn-lg"  value="Back" type="submit" style="background-color:brown" @click="hidePostPhotoForm"/>
										</div>
									</div>
								</div>						
								<div class="col-md-4 mb-4 d-flex align-items-center">
									<div class="form-outline">
										<div class="mt-4 pt-2">
											<input class="btn btn-primary btn-lg" value="Post" role="button" type="submit" style="background-color:green" @click = "postPhoto()"/>
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
