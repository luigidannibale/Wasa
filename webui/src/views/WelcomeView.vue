<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {		
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		// this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        </div>			
            <div class="jumbotron">
                <h1 class="display-4">Welcome to WasaPhoto</h1>
                <p class="lead">
                    Keep in touch with your friends by sharing photos of special moments.</p>
                <hr class="my-4">
                <p> This is WasaPhoto, an educational purpose based project<br>
                    Relised by D'annibale Luigi</p>
                <a class="btn btn-primary btn-lg" href="http://gamificationlab.uniroma1.it/notes/Project.pdf" role="button">
                    Learn more about the project specifications</a>
            </div>
		</div>
	
</template>

<style>
</style>







