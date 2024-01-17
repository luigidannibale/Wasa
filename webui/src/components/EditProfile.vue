<script>

export default {	
    name:"EditProfile",	
    
    data: function() {
		return {			
		}
	},
	methods: {	
		editProfile(){
            this.$emit("editProfile")
        },
	},
	mounted() {		
		
	}
}
</script>

<template>    	
    <button id="editProfile" type="button" class="btn btn-outline-dark" data-mdb-ripple-color="dark" style="z-index: 1; margin-right: 5px; margin-left: 5px;" @click="editProfile" >
        Edit profile
    </button> 	
</template>