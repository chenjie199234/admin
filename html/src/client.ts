import * as initializeAPI from './api/admin_initialize_browser'
import * as permissionAPI from './api/admin_permission_browser'
import * as userAPI from './api/admin_user_browser'
import * as appAPI from './api/admin_app_browser'

const host: string = "http://127.0.0.1:8000"
export const timeout: number = 3000
export const initializeClient: initializeAPI.InitializeBrowserClient = new initializeAPI.InitializeBrowserClient(host)
export const permissionClient: permissionAPI.PermissionBrowserClient = new permissionAPI.PermissionBrowserClient(host)
export const userClient: userAPI.UserBrowserClient = new userAPI.UserBrowserClient(host)
export const appClient: appAPI.AppBrowserClient = new appAPI.AppBrowserClient(host)
