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
			loggedId:null,
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
			this.loggedId = sessionStorage.getItem("id")						
			this.logProfile(this.loggedId)						
		},
		async showLike1(n){     			
			console.log(this.images)          
            var liked = this.images[n].im1.liked
            var r = null	            
            var loggedId = sessionStorage.getItem("id")
			
            if(liked){                
                try {                                               
                    await this.$axios({
                        method:"delete",
                        url:"/photos/"+this.images[n].im1.id+"/likes",					
                        headers:{
                            Authorization:loggedId
                        }
                    }).then((response)=>{
                        r = response}
                        )								
                } catch (e) {
                    r = e.response;							
                }		            
                switch (r.status) {
                    case 200:  
                        this.images[n].im1.liked = false
                        this.images[n].im1.n_likes -= 1
                        this.images[n].im1.likes = this.images[n].im1.likes.filter(item => item !== this.username);                        
                        return 													
                    default:
                        return 
                }
            }
            else{                
                try {                                               
                    await this.$axios({
                        method:"put",
                        url:"/photos/"+this.images[n].im1.id+"/likes",					
                        headers:{
                            Authorization:loggedId
                        }
                    }).then((response)=>{
                        r = response}
                        )								
                } catch (e) {
                    r = e.response;							
                }		            
                switch (r.status) {
                    case 200:  
                    case 201:  
                        this.images[n].im1.liked = true
                        this.images[n].im1.n_likes += 1                        
                        this.images[n].im1.likes = [...this.images[n].im1.likes, this.username]                        
                        return 													
                    default:
                        return 
                }
            }
        },
		async showLike2(n){               
            var liked = this.images[n].im2.liked
            var r = null	            
            var loggedId = sessionStorage.getItem("id")
            if(liked){                
                try {                                               
                    await this.$axios({
                        method:"delete",
                        url:"/photos/"+this.images[n].im2.id+"/likes",					
                        headers:{
                            Authorization:loggedId
                        }
                    }).then((response)=>{
                        r = response}
                        )								
                } catch (e) {
                    r = e.response;							
                }		            
                switch (r.status) {
                    case 200:  
                        this.images[n].im2.liked = false
                        this.images[n].im2.n_likes -= 1
                        this.images[n].im2.likes = this.images[n].im2.likes.filter(item => item !== this.username);                        
                        return 													
                    default:
                        return 
                }
            }
            else{                
                try {                                               
                    await this.$axios({
                        method:"put",
                        url:"/photos/"+this.images[n].im2.id+"/likes",					
                        headers:{
                            Authorization:loggedId
                        }
                    }).then((response)=>{
                        r = response}
                        )								
                } catch (e) {
                    r = e.response;							
                }		            
                switch (r.status) {
                    case 200:  
                    case 201:  
                        this.images[n].im2.liked = true
                        this.images[n].im2.n_likes += 1                        
                        this.images[n].im2.likes = [...this.images[n].im2.likes, this.username]                        
                        return 													
                    default:
                        return 
                }
            }
        },
		setUsername(name){
			this.username = name			
		},
		convertTime(timestamp){
			const date = new Date(timestamp);
      		const options = { year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: 'numeric', second: 'numeric'};
      		return date.toLocaleDateString('en-US', options);
		},
		async setPhotos(photos){
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
				couple.im1.id = photos[i].id				
				couple.im1.cap = photos[i].caption + " - posted on " + this.convertTime(photos[i].uploadTimestamp)
				couple.im1.showComments = false
                couple.im1.showLikes = false
				couple.im1.deletable = photos[i].userId == this.loggedId
					
				
				var likerslist = await this.getLikes(photos[i])                                
                if(likerslist) {
                    couple.im1.n_likes = likerslist.length
                    couple.im1.liked = false                    
                    for(let x = 0; x<likerslist.length; x++){
                        if(likerslist[x] === this.username){                        
                            couple.im1.liked = true
                        }                        
                    }
					
                }
                else couple.im1.n_likes = 0   
				couple.im1.likes = likerslist
				
				var commentslist = await this.getComments(photos[i])
                if(commentslist) {					
					couple.im1.n_comments = commentslist.length
					couple.im1.comments = await this.getCommentsUser(commentslist)
					
					for(let e = 0; e<commentslist.length; e++){                    
						let c = {}                    
						c.author = commentslist[e].userID
						c.text = commentslist[e].content
						c.id = commentslist[e].id
						c.delete = false
						couple.im1.comments[e] = c
					}                
				}					
                else couple.im1.n_comments = 0                 
                				
				if(i < photos.length-1){																				
					couple.im2.blob = "data:image/*;base64," + photos[i+1].image
					couple.im2.cap = photos[i+1].caption + " - posted on " + this.convertTime(photos[i+1].uploadTimestamp)
					couple.im2.id = photos[i+1].id					
					couple.im2.showComments = false
                	couple.im2.showLikes = false
					couple.im2.deletable = photos[i+1].userId == this.loggedId

					var likerslist = await this.getLikes(photos[i+1])					
					
					if(likerslist) {
						couple.im2.n_likes = likerslist.length
						couple.im2.liked = false                    
						for(let e = 0; e<likerslist.length; e++){
							if(likerslist[e] === this.username){                        
								couple.im2.liked = true
							}                        
						}
						
					}
					else couple.im2.n_likes = 0                 
					couple.im2.likes = likerslist                					                
					var commentslist = await this.getComments(photos[i+1])
					if(commentslist){
						couple.im2.n_comments = commentslist.length
						couple.im2.comments = await this.getCommentsUser(commentslist)
						for(let e = 0; e<commentslist.length; e++){                    
							let c = {}                    
							c.author = commentslist[e].userID
							c.text = commentslist[e].content
							c.id = commentslist[e].id
							c.delete = false
							couple.im2.comments[e] = c
						}
					}
					else couple.im2.n_comments = 0                 					
					}	
				else {
					couple.im2 = null
				}							
				couple.id = "c"+x.toString()
				this.images["c"+x.toString()] = couple								
			}							
		},
		async getCommentsUser(comments){
			var newComments = comments
			for (let index = 0; index < comments.length; index++) {
				const c = comments[index];				
				var newC = c
				var r
				try {					
					await this.$axios({
						method:"get",
						url:"/users/"+c.userID,						
						headers:{
							Authorization:this.loggedId
						}
					}).then((response)=>{
						r = response}
						)								
				} catch (e) {
					r = e.response;							
				}				
				switch (r.status) {
					case 200:						
						newC.userID = r.data.username						
						break;												
					default:						
						return					
				}				
				newComments[index] = newC
			}
			return newComments
		},
		async getLikes(im){
            var r = null	            
            var loggedId = sessionStorage.getItem("id")
			try {
											
				await this.$axios({
					method:"get",
					url:"/photos/"+im.id+"/likes/list",					
					headers:{
						Authorization:loggedId
					}
				}).then((response)=>{
					r = response}
					)								
			} catch (e) {
				r = e.response;							
			}		            
			switch (r.status) {
				case 200:                    	                    
                    return r.data																
				default:
					return null					
			}
        },
        async getComments(im){
            var r = null	            
            var loggedId = sessionStorage.getItem("id")
			try {
											
				await this.$axios({
					method:"get",
					url:"/photos/"+im.id+"/comments",					
					headers:{
						Authorization:loggedId
					}
				}).then((response)=>{
					r = response}
					)								
			} catch (e) {
				r = e.response;							
			}			
			switch (r.status) {
				case 200:	                    
                    return r.data																
				default:
					return null					
			} 
        }, 
		async postComment1(n){
			var imID = this.images[n].im1.id
			let val = document.getElementById("textcomment"+imID).value			
			if(val == ""){				
				return
			}	
			var r		
			try {
				await this.$axios({
					method:"post",
					url:"/photos/"+imID+"/comments",
					headers:{
						Authorization:this.loggedId
					},
                    params:{
                        content: val,
                        photoID: imID
                    }
				}).then((response)=>{
					r = response}
					)
			} catch (e) {
				r = e.response;							
			}		            
			switch (r.status) {
				case 201:
                    this.images[n].im1.n_comments+=1
                    let c = {}
                    c.author = this.username
                    c.text = val
                    c.photoID = imID
                    c.id = parseInt(r.data.slice(36),10)
                    this.images[n].im1.comments = [...this.images[n].im1.comments,c]
                    document.getElementById("textcomment"+imID).value = ""                    
				default:					
			}			
		
		},
		async postComment2(n){
			var imID = this.images[n].im2.id
			let val = document.getElementById("textcomment"+imID).value			
			if(val == ""){				
				return
			}	
			var r		
			try {
				await this.$axios({
					method:"post",
					url:"/photos/"+imID+"/comments",
					headers:{
						Authorization:this.loggedId
					},
                    params:{
                        content: val,
                        photoID: imID
                    }
				}).then((response)=>{
					r = response}
					)
			} catch (e) {
				r = e.response;							
			}		            
			switch (r.status) {
				case 201:
                    this.images[n].im2.n_comments+=1
                    let c = {}
                    c.author = this.username
                    c.text = val
                    c.photoID = imID
                    c.id = parseInt(r.data.slice(36),10)
                    this.images[n].im2.comments = [...this.images[n].im2.comments,c]
                    document.getElementById("textcomment"+imID).value = ""                    
				default:					
			}			
		
		},
		async deleteComment1(n,commId){
            var r = null	                        
			console.log("trying to delete comm n ",commId)
			try {											
				await this.$axios({
					method:"delete",
					url:"/photos/"+this.images[n].im1.id+"/comments/"+commId,
					headers:{
						Authorization:this.loggedId
					}
				}).then((response)=>{
					r = response}
					)
			} catch (e) {
				r = e.response;							
			}			            
			switch (r.status) {
				case 200:	                    
                    this.images[n].im1.n_comments  -=1                    
                    this.images[n].im1.comments = this.images[n].im1.comments.filter(item => item.id !== commId);
				default:
					return null					
			} 
        },  
		async deleteComment2(n,commId){
            var r = null	                        
			console.log("trying to delete comm n ",commId)
			try {											
				await this.$axios({
					method:"delete",
					url:"/photos/"+this.images[n].im2.id+"/comments/"+commId,
					headers:{
						Authorization:this.loggedId
					}
				}).then((response)=>{
					r = response}
					)
			} catch (e) {
				r = e.response;							
			}			            
			switch (r.status) {
				case 200:	                    
                    this.images[n].im2.n_comments  -=1                    
                    this.images[n].im2.comments = this.images[n].im2.comments.filter(item => item.id !== commId);
				default:
					return null					
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
		showLikes1(n){                     
            this.images[n].im1.showLikes = !this.images[n].im1.showLikes 
            if(this.images[n].im1.showLikes){
                this.images[n].im1.showComments = false
            }            
        },
        showComments1(n){
            this.images[n].im1.showComments = !this.images[n].im1.showComments
            if(this.images[n].im1.showComments){
                this.images[n].im1.showLikes = false
            }
        },
		showLikes2(n){                     
            this.images[n].im2.showLikes = !this.images[n].im2.showLikes 
            if(this.images[n].im2.showLikes){
                this.images[n].im2.showComments = false
            }            
        },
        showComments2(n){
            this.images[n].im2.showComments = !this.images[n].im2.showComments
            if(this.images[n].im2.showComments){
                this.images[n].im2.showLikes = false
            }
        },
		async deletePhoto(pid){			
			var r = null
			var id = sessionStorage.getItem("id")				
			try {			
				await this.$axios({
					method:"delete",
					url:"/photos/"+pid,
					headers:{
						Authorization:id,						
					},												
				}).then((response)=>{
					r = response}
					)								
			} catch (e) {
				r = e.response;							
			}
			console.log(r)
			switch (r.status) {
				case 200:										
					location.reload()
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
					<div class="col col-lg-10 ">												
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
							
							<div id="profilePhotos" class="card-body p-8 text-black">						
								<div class="d-flex justify-content-between align-items-center mb-8">
									<p class="lead fw-normal mb-0">Your photos</p>									
								</div>
								
								<div class="row" v-for="c in images" :key="c.id">
									<div v-if="c.im1" class="col-lg-4" style="margin-bottom: 10px;">																				
										<img v-if="c.im1.blob" :src="c.im1.blob" alt="Image 1" class="w-100 rounded-3">										
										<p v-if="c.im1.cap" v-text="c.im1.cap"></p>
									</div>									
									<div class="col-lg-2" v-if="c.im1">																				
										<div>											
											<span role="button" @click="showLikes1(c.id)" :id="'showlikes'+c.im1.id" style="margin: 0px 10px 0px 5px;"> {{ c.im1.n_likes }} </span>        
											<svg class="feather" @click="showLike1(c.id)" role="button" style="margin-right: 10px;" v-if="!c.im1.liked"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
											<svg class="like" 	 @click="showLike1(c.id)" role="button" style="margin-right: 10px;" v-if="c.im1.liked"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
										</div>																		
										<div role="button" @click="showComments1(c.id)">
											<span :id="'showcomments'+c.im1.id" style="margin: 0px 10px 0px 5px;"> {{ c.im1.n_comments }} </span>
											<svg class="feather" role="button" style="margin-right: 10px;"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
										</div>
										<div>
											<svg class="feather" role="button" @click="c.im1.deleteSect = !c.im1.deleteSect" v-if="c.im1.deletable"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
										</div>
										<div class="scrollable-container" id="commentsSection" v-if="c.im1.showComments">
											<h6 style="text-align: center;">Comments</h6>
											<div style="margin-top: 5px;" class="row">                                
												<span class="col">
													<textarea :id="'textcomment'+c.im1.id" type="text" class="form-control" placeholder="Insert new comment" rows="2" required></textarea>
													<svg class="feather" role="button" @click="postComment1(c.id)"><use href="/feather-sprite-v4.29.0.svg#edit"/></svg>
												</span>
												<span>
													
												</span>																																				
											</div>
											<div v-for="com in c.im1.comments" :key="com.id">
												<div class="row" style="padding: 0px;" v-if="com">
													<p class="mb-1">
														{{ com.author }}
														<svg class="feather" style="margin-left: 5px;" role="button" @click="com.delete = true" :id="'delete'+c.im1.id" v-if="username == com.author"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
														<span v-if="com.delete" class="small">
															Delete comment ? 
															<span role="button" @click="deleteComment1(c.id, com.id)" :id="'delete'+c.im1.id" style="color: green;">
																Yes
																<svg class="feather" ><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
															</span>                                       
															<span role="button" @click="com.delete = false" :id="'delete'+c.im1.id" style="color: red;">
																No
																<svg class="feather" > <use href="/feather-sprite-v4.29.0.svg#x"/></svg>
															</span>															
														</span>
													</p>                                                                                													
													<p class="small mb-0">
														{{ com.text }}
													</p>                                    
													<hr>
												</div>
											</div>
											                            
										</div>

										<div :id="'likesSection'+c.im1.id" v-if="c.im1.showLikes">
											<h6 class="text-center mb-12 pb-2" style="margin-top: 5px;">Likes</h6>
											<div v-for="u in c.im1.likes" :key="u">
												<div class="row" v-if="u">                                
													<p class="mb-1" role="button" @click="search(u)" > 
														{{ u }}
													</p>
												</div>
											</div>
										</div>

										<div :id="'deleteSect'+c.im1.id" v-if="c.im1.deleteSect"> 
											<span> Do you really want to delete this post? </span>
											<br>
											<span class="small" role="button" @click="deletePhoto(c.im1.id)" :id="'delete'+c.im1.id" style="color: green;">
                                                Yes
                                                <svg class="feather" ><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
                                            </span>                                       
                                            <span class="small" role="button" @click="c.im1.deleteSect = false" :id="'delete'+c.im1.id" style="color: red;">
                                                No
                                                <svg class="feather" > <use href="/feather-sprite-v4.29.0.svg#x"/></svg>
                                            </span>
										</div>
									</div>

									<div v-if="c.im2" class="col-lg-4"  style="margin-bottom: 10px;">										
										<img v-if="c.im2.blob" :src="c.im2.blob" alt="Image 2" class="w-100 rounded-3">
										<p v-if="c.im2.cap" v-text="c.im2.cap"></p>										
									</div>
									
									<div class="col-lg-2" v-if="c.im2">		
										<div>											
											<span role="button" @click="showLikes2(c.id)" :id="'showlikes'+c.im2.id" style="margin: 0px 10px 0px 5px;"> {{ c.im2.n_likes }} </span>        
											<svg class="feather" @click="showLike2(c.id)" role="button" style="margin-right: 10px;" v-if="!c.im2.liked"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
											<svg class="like" 	 @click="showLike2(c.id)" role="button" style="margin-right: 10px;" v-if="c.im2.liked"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
										</div>																		
										<div role="button" @click="showComments2(c.id)">
											<span :id="'showcomments'+c.im2.id" style="margin: 0px 10px 0px 5px;"> {{ c.im2.n_comments }} </span>
											<svg class="feather" role="button" style="margin-right: 10px;"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
										</div>
										<div>
											<svg class="feather" role="button" @click="c.im2.deleteSect = !c.im2.deleteSect" v-if="c.im2.deletable"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
										</div>
										         
										
										<div class="scrollable-container" id="commentsSection" v-if="c.im2.showComments"> 
											<h6 style="text-align: center;">Comments</h6>
											<div style="margin-top: 5px;" class="row">                                
												<span class="col">
													<textarea :id="'textcomment'+c.im2.id" type="text" class="form-control" placeholder="Insert new comment" rows="2" required></textarea>
													<svg class="feather" role="button" @click="postComment2(c.id)"><use href="/feather-sprite-v4.29.0.svg#edit"/></svg>
												</span>												
												<span>
													
												</span>																																				
											</div>
											<div v-for="com in c.im2.comments" :key="com.id">
												<div class="row" style="padding: 0px;" v-if="com">
													<p class="mb-1">
														{{ com.author }}
														<svg class="feather" style="margin-left: 5px;" role="button" @click="com.delete = true" :id="'delete'+c.im2.id" v-if="username == com.author">  <use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
														<span v-if="com.delete" class="small">
															Delete comment ? 
															<span role="button" @click="deleteComment2(c.id, com.id)" :id="'delete'+c.im2.id" style="color: green;">
																Yes
																<svg class="feather" ><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
															</span>                                       
															<span role="button" @click="com.delete = false" :id="'delete'+c.im2.id" style="color: red;">
																No
																<svg class="feather" > <use href="/feather-sprite-v4.29.0.svg#x"/></svg>
															</span>															
														</span>
													</p>                                                                                													
													<p class="small mb-0">
														{{ com.text }}
													</p>                                    
													<hr>
												</div>
											</div>
										</div>
										<div :id="'likesSection'+c.im2.id" v-if="c.im2.showLikes">
											<h6 class="text-center mb-12 pb-2" style="margin-top: 5px;">Likes</h6>
											<div v-for="u in c.im2.likes" :key="u">
												<div class="row" v-if="u">                                
													<p class="mb-1" role="button" @click="search(u)" > 
														{{ u }}
													</p>
												</div>
											</div>
										</div>
										<div :id="'deleteSect'+c.im2.id" v-if="c.im2.deleteSect"> 
											<span> Do you really want to delete this post? </span>
											<br>
											<span class="small" role="button" @click="deletePhoto(c.im2.id)" :id="'delete'+c.im2.id" style="color: green;">
                                                Yes
                                                <svg class="feather" ><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
                                            </span>                                       
                                            <span class="small" role="button" @click="c.im2.deleteSect = false" :id="'delete'+c.im2.id" style="color: red;">
                                                No
                                                <svg class="feather" > <use href="/feather-sprite-v4.29.0.svg#x"/></svg>
                                            </span>
										</div>
									</div>

									<hr>
								</div>								
							</div>

						</div>						
					</div>

					<div class="col col-lg-2">
						<div class="card" v-if="isFollowersvisible">							
							<div class="container py-2">
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
								<li class="list-group-item" v-for="f in followers" :key="f.username" v-text="f.username" @click="search(f)" role="button"></li>								
							</ul>
						</div>
						<div class="card" v-if="isFollowingVisible">							
							<div class="container py-2">
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
								<li class="list-group-item" v-for="f in following" :key="f.username" v-text="f.username" @click="search(f)" role="button"></li>								
							</ul>
						</div>
					</div>					
				</div>			
			</div>
</template>
<style>
.scrollable-container{
	overflow-y: auto;
	overflow-x: hidden;
	max-height: 250px;
}
</style>