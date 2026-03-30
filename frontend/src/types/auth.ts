export type AuthRoleType = 'ADMIN' | 'GROUP'

export interface AuthUser {
  userId: number
  username: string
  roleType: AuthRoleType
  groupId: number | null
  defaultRoute: string
}

export interface AuthLoginResult {
  accessToken: string
  tokenType: string
  expiresIn: number
  user: {
    userId: number
    username: string
    roleType: AuthRoleType
    groupId: number | null
  }
  defaultRoute: string
}

export interface AuthCurrentUserResult {
  userId: number
  username: string
  roleType: AuthRoleType
  groupId: number | null
  defaultRoute: string
}

export interface AuthLogoutResult {
  loggedOut: boolean
}

export interface AuthLoginPayload {
  username: string
  password: string
}

export interface AuthSession {
  accessToken: string
  tokenType: string
  expiresAt: number
  user: AuthUser
}
