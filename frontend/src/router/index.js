import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import ProfileView from '@/views/ProfileView.vue'
import LoginView from '@/views/LoginView.vue'
import NotificationsView from '@/views/NotificationsView.vue'
import GroupsView from '@/views/GroupsView.vue'
import GroupView from '@/views/GroupView.vue'
import EventsView from '@/views/EventsView.vue'
import ChatsView from '@/views/ChatsView.vue'

const routes = [
    { path: '/', name: 'home', component: HomeView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/profile/:id', name: 'profile', component: ProfileView },
    { path: '/notifications', name: 'notifications', component: NotificationsView },
    { path: '/groups', name: 'groups', component: GroupsView },
    { path: '/groups/:id', name: 'group', component: GroupView },
    { path: '/events', name: 'events', component: EventsView },
    { path: '/chats', name: 'chats', component: ChatsView },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router
