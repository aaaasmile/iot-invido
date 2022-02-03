import Dashboard from './views/dashboard.js'
import Login from './views/login.js'

export default [
  {
    path: '/', icon: 'dashboard', title: 'Dashboard', component: Dashboard,
    path: '/login', icon: 'dashboard', title: 'Login', component: Login
  }
]