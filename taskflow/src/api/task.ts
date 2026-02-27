import { request } from './request'

/** 后端返回的任务结构 */
export interface TaskDto {
  id: number
  description: string
  status: string
  user_id: number
  created_at: string
  updated_at: string
}

export interface CreateTaskBody {
  description: string
}

export interface UpdateTaskBody {
  description?: string
  status?: 'pending' | 'done'
}

export interface CreateTaskRes {
  message: string
  task: TaskDto
}

export interface GetTasksRes {
  tasks: TaskDto[]
}

export interface GetTaskRes {
  task: TaskDto
}

export interface UpdateTaskRes {
  message: string
  task: TaskDto
}

export interface DeleteTaskRes {
  message: string
}

function authHeaders(token: string): HeadersInit {
  return {
    Authorization: `Bearer ${token}`,
  }
}

/** @param date 日期字符串，格式 YYYY-MM-DD */
export function getTasksApi(token: string, date: string) {
  const query = new URLSearchParams({ date }).toString()
  return request<GetTasksRes>(`/tasks?${query}`, {
    method: 'GET',
    headers: authHeaders(token),
  })
}

export function createTaskApi(token: string, body: CreateTaskBody) {
  return request<CreateTaskRes>('/tasks', {
    method: 'POST',
    headers: authHeaders(token),
    body,
  })
}

export function updateTaskApi(token: string, id: number, body: UpdateTaskBody) {
  return request<UpdateTaskRes>(`/tasks/${id}`, {
    method: 'PUT',
    headers: authHeaders(token),
    body,
  })
}

export function deleteTaskApi(token: string, id: number) {
  return request<DeleteTaskRes>(`/tasks/${id}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  })
}
