import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import WelcomeView from '../views/WelcomeView.vue'
import SearchView from '../views/SearchView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [		
		{path: '/', component: WelcomeView},
		{path: '/login', component: LoginView},
		{path: '/users/search', component: SearchView},
		{path: '/users/search/:username', component: SearchView},
		{path: '/home/:id', component: HomeView},
		
	]
})

export default router
