<script>

import { VueElement } from 'vue'
import { reactive } from 'vue'
import App from '../App.vue'


export default {
	components: {
		
	},
	data: function() {
		return {
			err: false,
			errMess:null,			
            comments:{
                c1:{
                    author: "Luigi",
                    postTime: "ieri",
                    text: "Ciao a tutti",
                },
                c2:{
                    author: "Dan",
                    postTime: "ieri l'altro",
                    text: "Ciao a loro",
                },
                c3:{
                    author: "Marika",
                    postTime: "domani",
                    text: "Ciao a voi",
                },
            },
            likes:{
                l1:{
                    author: "Luigi",
                    postTime: "ieri",
                },
                l2:{
                    author: "Dan",
                    postTime: "ieri l'altro",
                },
                l3:{
                    author: "Marika",
                    postTime: "domani",
                },
            },            
            username:null,
            loggedId:null,
            following: null,
            images:reactive({}),            			
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
        async getStreamPhotos(loggedId){
			var r = null	            
			try {
				this.loading=true;								
				await this.$axios({
					method:"get",
					url:"/users/"+loggedId+"/profile/stream",					
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
                    await this.getProfileFollowing(loggedId)
                    this.setPhotos(r.data)
					break;												
				default:
					this.errAlert(r.data);
					break;
			}
		},		
        async getProfileFollowing(loggedId){
			var r = null
			try {
				this.loading=true;											
				await this.$axios({
					method:"get",
					url:"/users/"+ loggedId +"/followed",					
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
        async search(name) {			
			this.$router.push("/search/"+name)
            sessionStorage.setItem("username",name)
		},
        setFollowing(data){            
            this.following = {} 
            if(!data) 
                return
            data.forEach(f => {
                this.following[f.id] = f
            });                        
		},
        async setPhotos(photos){            
			if(!photos){ 
				this.n_photos = 0
                this.errAlert("Your feed is empty")
				return
			}                                 
			this.n_photos = photos.length			
            photos.sort((a, b) => {                
                const dateA = new Date(a.uploadTimestamp);
                const dateB = new Date(b.uploadTimestamp);

                return dateB - dateA;
            });
			for (let i = 0; i < photos.length; i+=1) {
                console.log("photo number ",i," ",photos[i])
				let couple = {
					im1: {},
				}
				couple.im1.blob = "data:image/*;base64," + photos[i].image
				couple.im1.id = photos[i].id                
                couple.im1.cap = photos[i].caption
                couple.im1.username = photos[i].username
                couple.im1.showComments = false
                couple.im1.showLikes = false
                var likerslist = await this.getLikes(photos[i])
                var commentslist = await this.getComments(photos[i])
                
                if(likerslist) {
                    couple.im1.n_likes = likerslist.length
                    couple.im1.liked = false                    
                    for(let x = 0; x<couple.im1.n_likes; x++){
                        if(likerslist[x] === this.username){                        
                            couple.im1.liked = true
                        }                        
                    }
                }
                else couple.im1.n_likes = 0                 
                if(commentslist) couple.im1.n_comments = commentslist.length
                else couple.im1.n_comments = 0                 
                couple.im1.likes = likerslist                
                couple.im1.comments = []
                for(let x = 0; x<couple.im1.n_comments; x++){                    
                    let c = {}                    
                    c.author = commentslist[x].userID
                    c.text = commentslist[x].content
                    c.id = commentslist[x].id
                    c.delete = false
                    couple.im1.comments[x] = c
                }                
				couple.im1.username = this.following[photos[i].userId].username                                
                this.images[photos[i].id] = couple                
			}							
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
        async postComment(n){
            var r = null	                        
            var e = document.getElementById("textcomment"+n)                        
			try {											
				await this.$axios({
					method:"post",
					url:"/photos/"+n+"/comments",
					headers:{
						Authorization:this.loggedId
					},
                    params:{
                        content: e.value,
                        photoID: n 
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
                    c.author = this.loggedId
                    c.text = e.value
                    c.photoID = n
                    c.id = parseInt(r.data.slice(36),10)
                    this.images[n].im1.comments = [...this.images[n].im1.comments,c]
                    e.value = ""
                    
				default:
					return null					
			}
        },
        async deleteComment(n,commId){
            var r = null	                        
			try {
											
				await this.$axios({
					method:"delete",
					url:"/photos/"+n+"/comments/"+commId,
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
        async showLike(n){                  
            var liked = this.images[n].im1.liked
            var r = null	            
            var loggedId = sessionStorage.getItem("id")
            if(liked){                
                try {                                               
                    await this.$axios({
                        method:"delete",
                        url:"/photos/"+n+"/likes",					
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
                        url:"/photos/"+n+"/likes",					
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
        showLikes(n){                     
            this.images[n].im1.showLikes = !this.images[n].im1.showLikes 
            if(this.images[n].im1.showLikes){
                this.images[n].im1.showComments = false
            }            
        },
        showComments(n){
            this.images[n].im1.showComments = !this.images[n].im1.showComments
            if(this.images[n].im1.showComments){
                this.images[n].im1.showLikes = false
            }
        },
        async log(){
            this.username = sessionStorage.getItem("loggedUsername")               
            this.loggedId = sessionStorage.getItem("id")                        
            this.getStreamPhotos(this.loggedId)                        
        },        
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
					Your feed, <span v-text="username"></span>
				</h1>
			</div>						
		</div>		
		
		<div class="alert alert-danger" role="alert" v-if="err" >
			<h4 class="alert-heading" v-text="errMess"></h4>			
		</div>					

		<div v-if="!err">			
		<section class="h-100 gradient-custom-2" v-show="!inputform" >			
			<div class="row" v-for="i in images" :key="i.im1.id">
                <div v-if="i.im1.blob" class="col" style="margin-bottom: 10px;">                    
                    <img :src="i.im1.blob" alt="Image 1" class="w-100 rounded-3">
                </div>
                <div class="col" v-if="i.im1" style="margin-bottom: 10px;">
                    <h6 v-if="i.im1.username" v-text="i.im1.username" role="button" @click="search(i.im1.username)"></h6>                    
                    <p v-if="i.im1.cap" v-text="i.im1.cap"></p>        
                    <p v-if="i.im1.ut" v-text="'Uploaded in ' + i.im1.ut"></p>                    
                    <div class="row" id = "likecommentBar">        
                        <div class="col-md-1">
                            <svg class="feather" role="button" @click="showLike(i.im1.id)" :id="'like'+i.im1.id" v-if="!i.im1.liked"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>  
                            <svg class="like" role="button" @click="showLike(i.im1.id)" :id="'like'+i.im1.id" v-if="i.im1.liked"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>  
                        </div>    
                        <div class="col-md-2" role="button" @click="showLikes(i.im1.id)" :id="'showlikes'+i.im1.id"> {{ i.im1.n_likes }} Likes </div>        
                        <div class="col-md-3"></div>
                        <div class="col-md-1">
                            <svg class="feather" role="button"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>         
                        </div >    
                        <div class="col-md-3" role="button" @click="showComments(i.im1.id)" :id="'showcomments'+i.im1.id"> {{ i.im1.n_comments }} Comments </div>
                        <div class="col-md-2"></div>
                    </div>
                    <section class="gradient-custom">                                 
                        
                        <div class="col-md-12 col-lg-10 col-xl-8" id="commentsSection" v-if="i.im1.showComments">
                            <h4 class="text-center mb-12 pb-2" style="margin-top: 30px;">Comments</h4>
                            <div class="card-body p-12" v-for="c in i.im1.comments" :key="c.id">
                                <div class="row" v-if="c">
                                    <p class="mb-1">
                                        {{ c.author }}
                                        <svg class="feather" style="margin-left: 40px;" role="button" @click="c.delete = true" :id="'delete'+i.im1.id" v-if="loggedId == c.author">  <use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                                        <span v-if="c.delete" class="small">
                                            Delete comment ? 
                                            <span role="button" @click="deleteComment(i.im1.id, c.id)" :id="'delete'+i.im1.id" style="color: green;">
                                                Yes
                                                <svg class="feather" ><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
                                            </span>                                       
                                            <span role="button" @click="c.delete = false" :id="'delete'+i.im1.id" style="color: red;">
                                                No
                                                <svg class="feather" > <use href="/feather-sprite-v4.29.0.svg#x"/></svg>
                                            </span>
                                            
                                        </span>
                                    </p>                                                                                                                    
                                    <p class="small mb-0">
                                        {{ c.text }}
                                    </p>                                    
                                    <hr>
                                </div>
                            </div>
                            <div style="margin-top: 20px;" class="row">                                
                                <span class="col">
                                    <textarea :id="'textcomment'+i.im1.id" cols="30" rows="2" placeholder="Insert new comment"></textarea>
                                </span>
                                <span class="col">
                                    <svg class="feather" style="margin-left: 10px;" role="button" @click="postComment(i.im1.id)"><use href="/feather-sprite-v4.29.0.svg#edit"/></svg>
                                </span>
                            </div>                            
                        </div>

                        <div class="col-md-12 col-lg-10 col-xl-4" id="likesSection" v-if="i.im1.showLikes">
                            <h4 class="text-center mb-12 pb-2" style="margin-top: 30px;">Likes</h4>
                            <div class="card-body p-12" v-for="u in i.im1.likes" :key="u">
                                <div v-if="u" class="row">                                
                                    <p class="mb-1" role="button" @click="search(u)" > 
                                        {{ u }}
                                    </p>
                                </div>
                            </div>
                        </div>
                    </section>
                    
                </div>
                <hr>
            </div>
		</section>		
		</div>
		
	</div>
	
</template>

<style>
</style>
