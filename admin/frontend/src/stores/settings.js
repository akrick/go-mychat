import { defineStore } from 'pinia'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    sidebarCollapsed: false,
    theme: 'light',
    tagsView: true,
    showSettings: false,
    fixedHeader: false,
    sidebarLogo: true
  }),

  actions: {
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },
    setTheme(theme) {
      this.theme = theme
      document.documentElement.setAttribute('data-theme', theme)
    },
    toggleTagsView() {
      this.tagsView = !this.tagsView
    },
    toggleFixedHeader() {
      this.fixedHeader = !this.fixedHeader
    }
  }
})
