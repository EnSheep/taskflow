import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'
import {
  getTasksApi,
  createTaskApi,
  updateTaskApi,
  deleteTaskApi,
  type TaskDto,
} from '../api/task'

/** 前端展示用的任务类型 */
export interface Task {
  id: number
  title: string
  completed: boolean
  createdAt: string
}

function dtoToTask(dto: TaskDto): Task {
  return {
    id: dto.id,
    title: dto.description || '',
    completed: dto.status === 'done',
    createdAt: dto.created_at,
  }
}

export const useTasksStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const authStore = useAuthStore()

  function getToken(): string {
    const t = authStore.token
    if (!t) throw new Error('未登录')
    return t
  }

  /** @param date 日期字符串，格式 YYYY-MM-DD */
  async function fetchTasks(date: string) {
    loading.value = true
    error.value = null
    try {
      const res = await getTasksApi(getToken(), date)
      tasks.value = (res.tasks || []).map(dtoToTask)
    } catch (e) {
      error.value = e instanceof Error ? e.message : '获取任务列表失败'
      tasks.value = []
    } finally {
      loading.value = false
    }
  }

  async function addTask(title: string) {
    if (!title.trim()) return
    const token = getToken()
    error.value = null
    try {
      const res = await createTaskApi(token, { description: title.trim() })
      tasks.value.push(dtoToTask(res.task))
    } catch (e) {
      error.value = e instanceof Error ? e.message : '创建失败'
      throw e
    }
  }

  async function toggleTask(id: number) {
    const task = tasks.value.find((t) => t.id === id)
    if (!task) return
    const nextStatus = task.completed ? 'pending' : 'done'
    error.value = null
    try {
      const res = await updateTaskApi(getToken(), id, { status: nextStatus })
      const idx = tasks.value.findIndex((t) => t.id === id)
      if (idx !== -1) tasks.value[idx] = dtoToTask(res.task)
    } catch (e) {
      error.value = e instanceof Error ? e.message : '更新失败'
      throw e
    }
  }

  async function removeTask(id: number) {
    error.value = null
    try {
      await deleteTaskApi(getToken(), id)
      tasks.value = tasks.value.filter((t) => t.id !== id)
    } catch (e) {
      error.value = e instanceof Error ? e.message : '删除失败'
      throw e
    }
  }

  return {
    tasks,
    loading,
    error,
    fetchTasks,
    addTask,
    toggleTask,
    removeTask,
  }
})
