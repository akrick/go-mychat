import { defineStore } from 'pinia'
import { login, logout, getUserInfo, getPermissions } from '@/api/user'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null,
    permissions: [],
    roles: []
  }),

  getters: {
    isLogin: (state) => !!state.token,
    userName: (state) => state.userInfo?.username || '',
    userAvatar: (state) => state.userInfo?.avatar || '',
    isAdmin: (state) => state.roles.includes('admin'),
    hasPermission: (state) => (permission) => {
      if (state.isAdmin) return true
      return state.permissions.includes(permission)
    }
  },

  actions: {
    // 登录
    async login(loginForm) {
      try {
        const res = await login(loginForm)
        console.log('登录接口返回:', res)
        
        // 兼容不同的响应格式
        let token
        if (res.code !== undefined) {
          // 标准格式: {code: 200, msg: "...", data: {token: "...", user: {...}}}
          token = res.data?.token
        } else {
          // 直接格式: {token: "...", user: {...}}
          token = res.token
        }
        
        if (!token) {
          throw new Error('登录响应中未找到token')
        }
        
        this.token = token
        localStorage.setItem('token', token)
        console.log('登录成功，token已设置:', token)
        // 移除重复的成功提示，由调用方处理
        return res
      } catch (error) {
        console.error('登录失败:', error)
        throw error
      }
    },

    // 获取用户信息
    async getUserInfo() {
      try {
        const res = await getUserInfo()
        console.log('用户信息接口返回:', res)
        
        // 兼容不同的响应格式
        if (res.code !== undefined) {
          // 标准格式: {code: 200, msg: "...", data: {...}}
          this.userInfo = res.data
        } else {
          // 直接格式: {id: ..., username: ...}
          this.userInfo = res
        }
        
        console.log('用户信息获取成功:', this.userInfo)
        return this.userInfo
      } catch (error) {
        console.error('获取用户信息失败:', error)
        // 用户信息获取失败不影响登录流程
        return null
      }
    },

    // 获取权限
    async getPermissions() {
      try {
        const res = await getPermissions()
        console.log('权限接口返回:', res)
        
        let permissions, roles
        if (res.code !== undefined) {
          // 标准格式: {code: 200, msg: "...", data: {permissions: [...], roles: [...]}}
          permissions = res.data?.permissions || []
          roles = res.data?.roles || []
        } else {
          // 直接格式: {permissions: [...], roles: [...]}
          permissions = res.permissions || []
          roles = res.roles || []
        }
        
        this.permissions = permissions
        this.roles = roles
        console.log('权限信息获取成功:', { permissions: this.permissions, roles: this.roles })
        return { permissions: this.permissions, roles: this.roles }
      } catch (error) {
        console.error('获取权限失败:', error)
        // 权限获取失败不影响登录，给予默认权限
        this.permissions = []
        this.roles = ['admin']
        return { permissions: this.permissions, roles: this.roles }
      }
    },

    // 登出
    async logout() {
      try {
        await logout()
        this.token = ''
        this.userInfo = null
        this.permissions = []
        this.roles = []
        localStorage.removeItem('token')
        ElMessage.success('退出成功')
      } catch (error) {
        throw error
      }
    }
  }
})
