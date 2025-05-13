import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import ProfileView from '@/views/ProfileView.vue'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import NotificationsView from '@/views/NotificationsView.vue'
import GroupsView from '@/views/GroupsView.vue'
import GroupView from '@/views/GroupView.vue'
import EventsView from '@/views/EventsView.vue'
import ChatsView from '@/views/ChatsView.vue'
import FollowsView from '@/views/FollowsView.vue'
import ErrorView from '@/views/ErrorView.vue'
import NotFound from '@/views/NotFound.vue'
import { useErrorStore } from '@/stores/error'


const routes = [
    { path: '/', name: 'home', component: HomeView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/register', name: 'register', component: RegisterView },
    { path: '/profile/:id', name: 'profile', component: ProfileView },
    { path: '/notifications', name: 'notifications', component: NotificationsView },
    { path: '/follows', name: 'follows', component: FollowsView },
    { path: '/groups', name: 'groups', component: GroupsView },
    { path: '/groups/:id', name: 'group', component: GroupView },
    { path: '/events', name: 'events', component: EventsView },
    { path: '/chats', name: 'chats', component: ChatsView },
    { path: '/error', name: 'Error', component: ErrorView },
    { path: '/:pathMatch(.*)*', name: 'notfound', component: NotFound } // Catch-all route if nothin before matches
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, from, next) => {
    const errorStore = useErrorStore()
    if (errorStore.title && to.path != '/error') {
        errorStore.clearError()
    }
    next()
})

export default router
