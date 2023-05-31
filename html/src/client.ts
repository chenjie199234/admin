import * as initializeAPI from 'admin/api/initialize_browser_toc'
import * as permissionAPI from 'admin/api/permission_browser_toc'
import * as userAPI from 'admin/api/user_browser_toc'
import * as appAPI from 'admin/api/app_browser_toc'

const host: string = "http://127.0.0.1:8000"
export const timeout: number = 3000
export const initializeClient: initializeAPI.InitializeBrowserClientToC = new initializeAPI.InitializeBrowserClientToC(host)
export const permissionClient: permissionAPI.PermissionBrowserClientToC = new permissionAPI.PermissionBrowserClientToC(host)
export const userClient: userAPI.UserBrowserClientToC = new userAPI.UserBrowserClientToC(host)
export const appClient: appAPI.AppBrowserClientToC = new appAPI.AppBrowserClientToC(host)
