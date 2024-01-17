<script>

export default {	
    name:"Unfollow",	
    
    data: function() {
		return {			
		}
	},
	methods: {	
		async unfollow(){
            try {	
				var searchedId = sessionStorage.getItem("searchedId")
				var loggedId = sessionStorage.getItem("id")
				var r;				
				await this.$axios({
					method:"delete",
					url:"/users/"+loggedId+"/followed/"+searchedId,						
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
    <button id="unfollow" type="button" class="btn btn-outline-dark" data-mdb-ripple-color="dark" style="z-index: 1; margin-right: 5px; margin-left: 5px;" @click="unfollow" >
        Unfollow
    </button> 	
</template>