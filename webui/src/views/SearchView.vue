<script>
import Profile from "../components/Profile.vue"
import SearchBar from "../components/SearchBar.vue"
import Follow from "../components/FollowButton.vue"
import Unfollow from "../components/UnfollowButton.vue"
import Ban from "../components/BanButton.vue"
import Unban from "../components/UnbanButton.vue"


export default {
	components: {
		
		Ban,
		Unban,
		Unfollow,
		Follow,
		Profile,
		SearchBar
	},
	data: function() {
		return {			
			err : false,
            errMess:null,
            
			username:null,
			searched:null,
			profileActive:false,

			_editable:false,
            _followable:false,
            _unfollowable:false,
            _bannable:false,
            _unbannable:false,
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
            var id = sessionStorage.getItem("id")			
            var r 
			if(this.searched === "search"){				
				this.profileActive =false
				return
			}

			{ // Gets the searched
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
					this.searched = r.data					
					break;					
				default:
					this.errAlert(r.data);
					return
									
			}	
			}

			sessionStorage.setItem("searchedId",this.searched.id)
			this.profileActive =true			
			{ // Checks if logged banned searched
			try {				
				await this.$axios({
					method:"get",
					url:"/users/"+id+"/banned",					
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
					
					break;					
				default:
					this.errAlert(r.data);
					return									
			}	
			if(r.data)
				r.data.forEach(f => {				
					if (f.username === this.searched.username){
						this.unbannable()
					}
				});
			}
			
			{ // Checks if logged follows searched
			try {				
				await this.$axios({
					method:"get",
					url:"/users/"+id+"/followed",					
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
					break;					
				default:
					this.errAlert(r.data);
					return									
			}	
			if(r.data)
				r.data.forEach(f => {				
					if (f.username === this.searched.username){
						this.unfollowable()
					}
				});
			}
			
			if(!this._unbannable && !this._unfollowable)
				this.followable()
			
			
		},
		async editable(){
			this._editable = true
			this._bannable = false
			this._followable = false
			this._unbannable = false
			this._unfollowable = false
		},
		async followable(){
			this._editable = false
			this._bannable = true
			this._followable = true
			this._unbannable = false
			this._unfollowable = false
		},
		async unfollowable(){
			this._editable = false
			this._bannable = true
			this._followable = false
			this._unbannable = false
			this._unfollowable = true
		},
		async unbannable(){
			this._editable = false
			this._bannable = false
			this._followable = false
			this._unbannable = true
			this._unfollowable = false
		},
		async errAlert(data){			
			this.err = true;
			this.errMess = data;
		},		
		async nullAlerts(){
			this.err = false;			
			this.errMess = null;			
		},
		refresh(){
			window.location.reload()
		}
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
			<div class="row">				
				<div class="col col-lg-2" v-if="_bannable">
					<Ban @refresh="refresh()" @err="errAlert"></Ban>
				</div>
				<div class="col col-lg-2" v-if="_unbannable">
					<Unban @refresh="refresh()" @err="errAlert"></Unban>
				</div>
				<div class="col col-lg-2" v-if="_followable">
					<Follow @refresh="refresh()" @err="errAlert"></Follow>
				</div>
				<div class="col col-lg-2" v-if="_unfollowable">
					<Unfollow @refresh="refresh()" @err="errAlert"></Unfollow>
				</div>
			</div>						
			<Profile></Profile>
		</section>
</template>

<style>
</style>







