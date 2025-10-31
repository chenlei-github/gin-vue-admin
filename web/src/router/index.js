import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/init',
    name: 'Init',
    component: () => import('@/view/init/index.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/view/login/index.vue')
  },
  {
    path: '/scanUpload',
    name: 'ScanUpload',
    meta: {
      title: '扫码上传',
      client: true
    },
    component: () => import('@/view/example/upload/scanUpload.vue')
  },
  {
    path: '/markets',
    name: 'MarketList',
    meta: {
      title: '市场列表',
      keepAlive: true
    },
    component: () => import('@/view/brandtrekin/display/market-list.vue')
  },
  {
    path: '/markets/:id',
    name: 'MarketDetail',
    meta: {
      title: '市场详情',
      keepAlive: false
    },
    component: () => import('@/view/brandtrekin/display/market-detail.vue')
  },
  {
    path: '/markets/:marketId/brands/:brandName',
    name: 'BrandDetail',
    meta: {
      title: '品牌详情',
      keepAlive: false
    },
    component: () => import('@/view/brandtrekin/display/brand-detail.vue')
  },
  {
    path: '/:catchAll(.*)',
    meta: {
      closeTab: true
    },
    component: () => import('@/view/error/index.vue')
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
