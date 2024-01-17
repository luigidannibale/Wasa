<script>

export default {	
    name:"Follow",	
    
    data: function() {
		return {			
		}
	},
	methods: {	
		async follow(){
            try {	
				var searchedId = sessionStorage.getItem("searchedId")
				var loggedId = sessionStorage.getItem("id")
				var r;				
				await this.$axios({
					method:"put",
					url:"/users/"+loggedId+"/followed",
					params:{
						userToFollowID: searchedId
					},					
					headers:{
						Authorization: loggedId
					}
				}).then((response)=>{
					r = response}
					)				
			} catch (e) {
				r = e.response;							
			}					
			
			switch (r.status) {
				case 201:		
					this.$emit("refresh")
					break;
				default:					
					this.$emit("err",r.data)
					break;
			}
        },
	},
	mounted() {		
		
	}
}
</script>

<template>    	
    <button id="follow" type="button" class="btn btn-outline-dark" data-mdb-ripple-color="dark" style="z-index: 1; margin-right: 5px; margin-left: 5px;" @click="follow" >
        Follow
    </button> 	
</template>