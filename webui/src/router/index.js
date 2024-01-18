import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import WelcomeView from '../views/WelcomeView.vue'
import SearchView from '../views/SearchView.vue'
import FeedView from '../views/FeedView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: WelcomeView},
		{path: '/login', component: LoginView},		
		{
			path: '/search', 
			component: SearchView,
			children:[{
				path: ":username",
				component: SearchView,
			}]
		},		
		{path: '/home/:id', component: HomeView},
		{path: '/feed', component: FeedView},
	]
})

export default router
