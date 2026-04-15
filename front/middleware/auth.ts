import { useGlobalState } from '~/store'

export default defineNuxtRouteMiddleware((to) => {
  const global = useGlobalState()
  const token = global.value.userinfo.token

  const adminRoutes = ['/sys/settings']
  const authRoutes = ['/new', '/user/settings', '/user/calendar']
  const editPattern = /^\/edit\//

  const needsAuth =
    authRoutes.includes(to.path) ||
    adminRoutes.includes(to.path) ||
    editPattern.test(to.path)

  if (needsAuth && !token) {
    return navigateTo('/user/login')
  }
})
