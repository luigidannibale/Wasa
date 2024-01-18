<script>
import { reactive } from 'vue';


export default {	
    name:"Profile",	
    data: function() {
		return {					
            fullname: null,			
            n_photos:null,
			n_following:null,
			n_followers:null,			
			username: null,
			id: null,
			err: false,
			errMess:null,

            isFollowersvisible: false,
			isFollowingVisible: false,
			followers:null,
			following:null,
			images:reactive({}),			
		}
	},
	methods: {	
		async logProfile(loggedId){
			var r = null

			try {
				this.loading=true;								
				await this.$axios({
					method:"get",
					url:"/users",
					params:{
						username:this.username
					},
					headers:{
						Authorization:loggedId
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
					this.fullname = r.data["name"] + " " + r.data["surname"]	
					this.id = r.data["id"]												
					if (this.fullname === " ") {
						this.inputform = true;
					}
					break;												
				default:
					this.errAlert(r.data);
					return					
			}
			this.getProfilePhotos(loggedId,this.id)
			this.getProfileFollowers(loggedId,this.id)			
			this.getProfileFollowing(loggedId,this.id)
		},	
		async getProfilePhotos(loggedId,profileId){
			var r = null
			
			try {
				this.loading=true;								
				await this.$axios({
					method:"get",
					url:"/users/"+profileId+"/profile/photos",					
					headers:{
						Authorization:loggedId
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
					this.setPhotos(r.data) 
					break;												
				default:
					this.errAlert(r.data);
					break;
			}
		},
		async getProfileFollowers(loggedId,profileId){
			var r = null
			try {
				this.loading=true;								
				await this.$axios({
					method:"get",
					url:"/users/"+profileId+"/followers",					
					headers:{
						Authorization:loggedId
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
					this.setFollowers(r.data)					
					break;												
				default:
					this.errAlert(r.data);
					break;
			}
		},
		async getProfileFollowing(loggedId,profileId){
			var r = null
			try {
				this.loading=true;											
				await this.$axios({
					method:"get",
					url:"/users/"+ profileId +"/followed",					
					headers:{
						Authorization:loggedId
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
					this.setFollowing(r.data) 				
					break;												
				default:
					this.errAlert(r.data);
					break;
			}
		},
		async configProfile(){
			if(!this.username)
			{							
				this.setUsername(sessionStorage.getItem("username"))				
			}
			
			var loggedId = sessionStorage.getItem("id")						
			this.logProfile(loggedId)						
		},
		setUsername(name){
			this.username = name			
		},
		setPhotos(photos){
			if(!photos){ 
				this.n_photos = 0	
				return
			}
			this.n_photos = photos.length			
			for (let i = 0,x = 0; i < photos.length; i+=2,x++) {
				let couple = {
					im1: {},
					im2: {},
				}												
				couple.im1.blob = "data:image/*;base64," + photos[i].image
				couple.im1.cap = photos[i].caption											
				if(i < photos.length-1){															
					couple.im2.blob = "data:image/*;base64," + photos[i+1].image
					couple.im2.cap = photos[i+1].caption
				}								
				this.images["c"+x.toString()] = couple								
			}							
		},
		setFollowers(followers){
			if(followers){
				this.n_followers = followers.length			
				this.followers = followers
			}
			else{this.n_followers = 0}
			
		},
		setFollowing(following){
			if(following){
			this.n_following = following.length
			this.following = following
			}
			else {this.n_following = 0}
		},
		showFollowing(){
			if (this.isFollowersvisible) {
				this.hideFollowers()
			}
			this.isFollowingVisible = true
		},
		showFollowers(){
			if (this.isFollowingVisible) {
				this.hideFollowing()
			}
			this.isFollowersvisible = true			
		},
		hideFollowing(){
			this.isFollowingVisible = false
		},
		hideFollowers(){
			this.isFollowersvisible = false
		},        
		errAlert(data){
			this.err = true;			
			this.errMess = data;
		},
		search(f){			
			sessionStorage.setItem("username",f.username)    			
			this.$router.push("/search/"+f.username)	            
		},
	},
	mounted() {		
		this.configProfile()		
	}
}
</script>

<template>    
	<div class="alert alert-danger" role="alert" v-if="err" >
		<h4 class="alert-heading" v-text="errMess"></h4>			
	</div>			

    <div class="container py-5 h-100">										
				<div class="row d-flex justify-content-left h-100">					
					<div class="col col-lg-8 ">												
						<div class="card">																			
							
							<div id="profileHeader" class="p-4 text-black" style="background-color: #f8f9fa;" >
								<div class="d-flex justify-content-between align-items-center mb-4">
									<p class="lead fw-normal mb-0" v-text="fullname"></p>									
								</div>								
								
								<div class="d-flex justify-content-end text-center py-1">

									<div class="card px-3" >
										<p class="mb-1 h5" v-text="n_photos"></p>
										<p class="small text-muted mb-0">Photos</p>
									</div>									

									<a @click="showFollowers()" role="button">
										<div class="card px-3" style="margin-left: 10px; margin-right: 10px;">
											<p class="mb-1 h5" v-text="n_followers"></p>
											<p class="small text-muted mb-0">Followers</p>
										</div>
									</a>
									<a @click="showFollowing()" role="button">
                                        <div class="card px-3" >
                                            <p class="mb-1 h5" v-text="n_following"></p>
                                            <p class="small text-muted mb-0">Following</p>
                                        </div>
                                    </a>
								</div>								
							</div>
							
							<div id="profilePhotos" class="card-body p-4 text-black">						
								<div class="d-flex justify-content-between align-items-center mb-4">
									<p class="lead fw-normal mb-0">Your photos</p>									
								</div>
								
								<div class="row" v-for="c in images">
									<div class="col" style="margin-bottom: 10px;">
										<p v-if="c.im1.cap" v-text="c.im1.cap"></p>
										<img v-if="c.im1.blob" :src="c.im1.blob" alt="Image 1" class="w-100 rounded-3">
									</div>
									<div class="col"  style="margin-bottom: 10px;">
										<p v-if="c.im2.cap" v-text="c.im2.cap"></p>
										<img v-if="c.im2.blob" :src="c.im2.blob" alt="Image 2" class="w-100 rounded-3">										
									</div>
									<hr>
								</div>								
							</div>

						</div>						
					</div>

					<div class="col col-lg-4 ">
						<div class="card" v-if="isFollowersvisible">							
							<div class="container py-3">
								<div class="row">
									<div class="col-10">
										<h4>Followers</h4>
									</div>
									<div class="col-2">
										<button type="button" class="close" aria-label="Close" @click="hideFollowers()">
											<span aria-hidden="true">&times;</span>
										</button>
									</div>
								</div>
							</div>
							<ul class="list-group list-group-flush">
								<li class="list-group-item" v-for="f in followers" v-text="f.username" @click="search(f)" role="button"></li>								
							</ul>
						</div>
						<div class="card" v-if="isFollowingVisible">							
							<div class="container py-3">
								<div class="row">
									<div class="col-10">
										<h4>Following</h4>
									</div>
									<div class="col-2">
										<button type="button" class="close" aria-label="Close" @click="hideFollowing()">
											<span aria-hidden="true">&times;</span>
										</button>
									</div>
								</div>
							</div>
							<ul class="list-group list-group-flush">
								<li class="list-group-item" v-for="f in following" v-text="f.username" @click="search(f)" role="button"></li>								
							</ul>
						</div>
					</div>					
				</div>			
			</div>
</template>