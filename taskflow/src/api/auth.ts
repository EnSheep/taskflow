import { request } from './request'

export interface RegisterBody {
  username: string
  password: string
  email: string
}

export interface LoginBody {
  username: string
  password: string
}

export interface UserInfo {
  id: number
  username: string
  email: string
}

export interface RegisterRes {
  message: string
  user: UserInfo
}

export interface LoginRes {
  message: string
  token: string
  user: UserInfo
}

export function registerApi(body: RegisterBody) {
  return request<RegisterRes>('/register', { method: 'POST', body })
}

export function loginApi(body: LoginBody) {
  return request<LoginRes>('/login', { method: 'POST', body })
}
