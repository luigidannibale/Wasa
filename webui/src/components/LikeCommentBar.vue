<script>

export default {	
    name:"LikeComment",	
    
    data: function() {
		return {			
		}
	},
	methods: {	
		
	},
	mounted() {		
		
	}
}
</script>

<template>    	
    <div class="row">        
        <div class="col-md-1">
            <svg class="feather" role="button" @click="like()" id="heart"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>  
        </div>    
        <div class="col-md-1" role="button"> Likes </div>        
        <div class="col-md-4"></div>
        <div class="col-md-1">
            <svg class="feather" role="button"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>         
        </div >    
        <div class="col-md-2" role="button"> Comments </div>
        <div class="col-md-3"></div>
    </div>
</template>