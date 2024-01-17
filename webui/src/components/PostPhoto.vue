<script>

export default {	
    name:"PostPhoto",	
    
    data: function() {
		return {			
		}
	},
	methods: {	
		postPhoto(){
            this.$emit("postPhoto")
        },
	},
	mounted() {		
		
	}
}
</script>

<template>    	
    <button id="postPhoto" type="button" class="btn btn-outline-dark" data-mdb-ripple-color="dark" style="z-index: 1; margin-right: 5px; margin-left: 5px;" @click="postPhoto" >
        Post Photo
    </button> 	
</template>