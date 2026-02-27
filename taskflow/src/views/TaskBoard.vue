<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useTasksStore } from '../stores/tasks'

const router = useRouter()
const auth = useAuthStore()
const tasksStore = useTasksStore()
const showHistory = ref(false)
const newTitle = ref('')

// 日期选择器：当前查看的年月
const viewingDate = ref(new Date())
const viewingYear = computed({
  get: () => viewingDate.value.getFullYear(),
  set: (y: number) => {
    const d = new Date(viewingDate.value)
    d.setFullYear(y)
    viewingDate.value = d
  },
})
const viewingMonth = computed({
  get: () => viewingDate.value.getMonth(),
  set: (m: number) => {
    const d = new Date(viewingDate.value)
    d.setMonth(m)
    viewingDate.value = d
  },
})

// 日历网格：当月 + 前后补全到 6 行 x 7 列
const weekDays = ['日', '一', '二', '三', '四', '五', '六']
const calendarDays = computed(() => {
  const y = viewingYear.value
  const m = viewingMonth.value
  const first = new Date(y, m, 1)
  const last = new Date(y, m + 1, 0)
  const startPad = first.getDay()
  const daysInMonth = last.getDate()
  const total = 42
  const out: { day: number; isCurrentMonth: boolean; date: Date }[] = []
  const prevMonthLast = new Date(y, m, 0).getDate()
  for (let i = 0; i < startPad; i++) {
    out.push({ day: prevMonthLast - startPad + 1 + i, isCurrentMonth: false, date: new Date(y, m - 1, prevMonthLast - startPad + 1 + i) })
  }
  for (let d = 1; d <= daysInMonth; d++) {
    out.push({ day: d, isCurrentMonth: true, date: new Date(y, m, d) })
  }
  let rest = total - out.length
  for (let d = 1; d <= rest; d++) {
    out.push({ day: d, isCurrentMonth: false, date: new Date(y, m + 1, d) })
  }
  return out
})

// 可选年份：当前年往前 5 年、往后 5 年（共 11 年）
const yearOptions = computed(() => {
  const y = new Date().getFullYear()
  return Array.from({ length: 11 }, (_, i) => y - 5 + i)
})
const monthOptions = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]

const minYear = computed(() => new Date().getFullYear() - 5)
const maxYear = computed(() => new Date().getFullYear() + 5)

// 格式化为 API 需要的 YYYY-MM-DD
function formatDateToYYYYMMDD(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

// 当前选中的日期（用于高亮）
const selectedDay = ref<Date | null>(null)
async function selectDay(cell: { date: Date; isCurrentMonth: boolean }) {
  selectedDay.value = cell.date
  await tasksStore.fetchTasks(formatDateToYYYYMMDD(cell.date))
  showHistory.value = false
}
function isSelected(cell: { date: Date }) {
  if (!selectedDay.value) return false
  const a = selectedDay.value
  const b = cell.date
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
}

// 滚轮切换月份：向上滚上一月，向下滚下一月，并限制在可选年份内
function changeMonth(delta: number) {
  const d = new Date(viewingDate.value)
  d.setMonth(d.getMonth() + delta)
  const y = d.getFullYear()
  if (y < minYear.value || y > maxYear.value) return
  viewingDate.value = d
}

function onDatePickerWheel(e: WheelEvent) {
  e.preventDefault()
  if (e.deltaY > 0) changeMonth(1)
  else changeMonth(-1)
}

onMounted(() => {
  tasksStore.fetchTasks(formatDateToYYYYMMDD(new Date()))
})

function logout() {
  auth.logout()
  router.push('/login')
}

function addTask() {
  tasksStore.addTask(newTitle.value)
  newTitle.value = ''
}

function checkHistoryTasks() {
  showHistory.value = !showHistory.value
  if (!showHistory.value) {
    tasksStore.fetchTasks(formatDateToYYYYMMDD(new Date()))
  }
}

const completedCount = computed(() => tasksStore.tasks.filter((t) => t.completed).length)
const totalCount = computed(() => tasksStore.tasks.length)
const progress = computed(() =>
  totalCount.value ? Math.round((completedCount.value / totalCount.value) * 100) : 0
)

// 编辑任务描述：当前正在编辑的任务 id，以及输入框内容
const editingTaskId = ref<number | null>(null)
const editingValue = ref('')
const editInputRef = ref<HTMLInputElement | null>(null)
const editStartedAt = ref(0)

function isEditing(taskId: number) {
  return editingTaskId.value != null && Number(editingTaskId.value) === Number(taskId)
}

function startEdit(task: { id: number; title: string }) {
  editingTaskId.value = Number(task.id)
  editingValue.value = task.title
  editStartedAt.value = Date.now()
  nextTick(() => {
    editInputRef.value?.focus()
  })
}

function cancelEdit() {
  editingTaskId.value = null
  editingValue.value = ''
}

async function submitEdit(taskId: number) {
  // 避免点击编辑按钮后 blur 立刻触发导致未保存就退出
  if (Date.now() - editStartedAt.value < 150) return
  const val = editingValue.value.trim()
  if (val === '') {
    cancelEdit()
    return
  }
  await tasksStore.updateTaskDescription(taskId, val)
  cancelEdit()
}

function onEditKeydown(e: KeyboardEvent, taskId: number) {
  if (e.key === 'Enter') {
    e.preventDefault()
    submitEdit(taskId)
  }
  if (e.key === 'Escape') {
    cancelEdit()
  }
}
</script>

<template>
  <div class="task-board">
    <header class="header">
      <div class="header-inner">
        <div class="brand">
          <span class="logo">◆</span>
          <span class="app-name">TaskFlow</span>
        </div>
        <div class="user-row">
          <span class="user-name">{{ auth.username || '用户' }}</span>
          <button type="button" class="btn-logout" @click="logout">退出</button>
        </div>
      </div>
    </header>

    <main class="main">
      <div class="board-card glass">
        <div class="board-header">
          <div class="header-content">
            <h2>任务栏</h2>
            <h4 style="cursor: pointer;" @click="checkHistoryTasks">查看历史任务</h4>
          </div>

          <div class="stats">
            <span class="stat">{{ completedCount }} / {{ totalCount }} 已完成</span>
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: progress + '%' }" />
            </div>
          </div>
        </div>

        <!-- 创建任务完成任务列表 -->
        <div v-show="!showHistory">
          <div class="add-task">
            <input v-model="newTitle" type="text" placeholder="输入新任务，按回车或点击添加" @keydown.enter.prevent="addTask" />
            <button type="button" class="btn-add" @click="addTask">添加任务</button>
          </div>

          <div class="task-list-wrap">
            <ul class="task-list">
              <li v-for="task in tasksStore.tasks" :key="task.id" class="task-item"
                :class="{ completed: task.completed, 'is-editing': isEditing(task.id) }">
                <div class="task-content">
                  <template v-if="isEditing(task.id)">
                    <label class="task-label task-label--cb-only">
                      <input type="checkbox" :checked="task.completed" @change="tasksStore.toggleTask(task.id)" />
                      <span class="checkmark" />
                    </label>
                    <input
                      ref="editInputRef"
                      v-model="editingValue"
                      type="text"
                      class="task-edit-input"
                      placeholder="输入任务描述..."
                      @blur="submitEdit(task.id)"
                      @keydown="onEditKeydown($event, task.id)"
                      @click.stop
                    />
                  </template>
                  <template v-else>
                    <label class="task-label">
                      <input type="checkbox" :checked="task.completed" @change="tasksStore.toggleTask(task.id)" />
                      <span class="checkmark" />
                      <span class="task-title">{{ task.title }}</span>
                    </label>
                  </template>
                </div>
                <button type="button" class="btn-edit" title="编辑" @click.stop="startEdit(task)">
                  ✎
                </button>
                <button type="button" class="btn-remove" title="删除" @click.stop="tasksStore.removeTask(task.id)">
                  ×
                </button>
              </li>
              <li v-if="!tasksStore.tasks.length" class="empty-state">
                <p>还没有任务，点击上方按钮创建一个吧</p>
              </li>
            </ul>
          </div>
        </div>

        <!-- 日期选择 历史任务列表展示 -->
        <div v-show="showHistory" class="history-view">
          <div class="date-picker-sketch" @wheel.prevent="onDatePickerWheel">
            <div class="date-picker-header">
              <div class="date-picker-select-wrap">
                <span class="date-picker-label">年份</span>
                <select v-model.number="viewingYear" class="date-picker-select">
                  <option v-for="y in yearOptions" :key="y" :value="y">{{ y }}</option>
                </select>
              </div>
              <div class="date-picker-select-wrap">
                <span class="date-picker-label">月份</span>
                <select v-model.number="viewingMonth" class="date-picker-select">
                  <option v-for="(mo, i) in monthOptions" :key="i" :value="i">{{ mo }} 月</option>
                </select>
              </div>
            </div>
            <div class="date-picker-grid">
              <div v-for="w in weekDays" :key="w" class="grid-cell grid-cell-head">{{ w }}</div>
              <button
                v-for="(cell, idx) in calendarDays"
                :key="idx"
                type="button"
                class="grid-cell grid-cell-day"
                :class="{
                  'is-current-month': cell.isCurrentMonth,
                  'is-other-month': !cell.isCurrentMonth,
                  'is-selected': isSelected(cell),
                }"
                @click="selectDay(cell)"
              >
                {{ cell.day }}
              </button>
            </div>
          </div>
          <div class="task-list-wrap">
            <ul class="task-list">
              <li v-for="task in tasksStore.tasks" :key="task.id" class="task-item"
                :class="{ completed: task.completed, 'is-editing': isEditing(task.id) }">
                <div class="task-content">
                  <template v-if="isEditing(task.id)">
                    <label class="task-label task-label--cb-only">
                      <input type="checkbox" :checked="task.completed" @change="tasksStore.toggleTask(task.id)" />
                      <span class="checkmark" />
                    </label>
                    <input
                      ref="editInputRef"
                      v-model="editingValue"
                      type="text"
                      class="task-edit-input"
                      placeholder="输入任务描述..."
                      @blur="submitEdit(task.id)"
                      @keydown="onEditKeydown($event, task.id)"
                      @click.stop
                    />
                  </template>
                  <template v-else>
                    <label class="task-label">
                      <input type="checkbox" :checked="task.completed" @change="tasksStore.toggleTask(task.id)" />
                      <span class="checkmark" />
                      <span class="task-title">{{ task.title }}</span>
                    </label>
                  </template>
                </div>
                <button type="button" class="btn-edit" title="编辑" @click.stop="startEdit(task)">
                  ✎
                </button>
                <button type="button" class="btn-remove" title="删除" @click.stop="tasksStore.removeTask(task.id)">
                  ×
                </button>
              </li>
              <li v-if="!tasksStore.tasks.length" class="empty-state">
                <p>还没有任务，返回创建任务吧</p>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.task-board {
  height: 100vh;
  max-height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg);
}

.header {
  flex-shrink: 0;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.02);
}

.header-inner {
  max-width: 720px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.logo {
  font-size: 1.25rem;
  color: var(--accent);
}

.app-name {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 1.25rem;
  color: var(--text);
}

.user-row {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-name {
  color: var(--text-muted);
  font-size: 0.9rem;
}

.btn-logout {
  padding: 0.5rem 1rem;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-muted);
  font-size: 0.875rem;
  cursor: pointer;
  transition: color 0.2s, border-color 0.2s;
}

.btn-logout:hover {
  color: var(--text);
  border-color: rgba(255, 255, 255, 0.2);
}

.main {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  padding: 2rem 1.5rem;
  max-width: 720px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.board-card {
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 1.75rem;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.4);
}

.glass {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
}

.board-header {
  flex-shrink: 0;
  margin-bottom: 1.5rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.board-header h2 {
  font-family: var(--font-display);
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text);
  margin: 0 0 0.75rem;
}

.stats {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.stat {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.progress-bar {
  height: 6px;
  border-radius: 3px;
  background: rgba(255, 255, 255, 0.08);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  background: linear-gradient(90deg, var(--accent) 0%, var(--accent-dark) 100%);
  transition: width 0.3s ease;
}

.add-task {
  flex-shrink: 0;
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.add-task input {
  flex: 1;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: var(--text);
  font-size: 1rem;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.add-task input::placeholder {
  color: var(--text-muted);
  opacity: 0.7;
}

.add-task input:focus {
  outline: none;
  border-color: var(--accent);
}

.btn-add {
  padding: 0.75rem 1.25rem;
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%);
  color: var(--bg);
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.2s;
  white-space: nowrap;
  box-shadow: 0 4px 14px rgba(251, 191, 36, 0.3);
}

.btn-add:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(251, 191, 36, 0.35);
}

/* 任务列表滚动容器：占满剩余空间，仅此处可滚动，避免出现页面级滚动条 */
.task-list-wrap {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  margin: 0 -0.25rem 0 0;
  padding-right: 0.25rem;
  scrollbar-gutter: stable;
}

/* 自定义滚动条样式（Chrome / Edge / Safari） */
.task-list-wrap::-webkit-scrollbar {
  width: 8px;
}

.task-list-wrap::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.04);
  border-radius: 4px;
  margin: 4px 0;
}

.task-list-wrap::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.18);
  border-radius: 4px;
}

.task-list-wrap::-webkit-scrollbar-thumb:hover {
  background: rgba(251, 191, 36, 0.5);
}

.task-list-wrap::-webkit-scrollbar-thumb:active {
  background: var(--accent);
}

/* Firefox 滚动条颜色 */
.task-list-wrap {
  scrollbar-color: rgba(255, 255, 255, 0.25) rgba(255, 255, 255, 0.06);
  scrollbar-width: thin;
}

.task-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.task-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: background 0.2s, border-color 0.2s;
}

.task-item:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.1);
}

.task-item.completed .task-title {
  text-decoration: line-through;
  color: var(--text-muted);
}

.task-content {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: 0;
}

.task-label {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  user-select: none;
  min-width: 0;
}

.task-label--cb-only {
  flex: 0 0 auto;
  cursor: pointer;
}

.task-label input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.checkmark {
  width: 20px;
  height: 20px;
  border-radius: 6px;
  border: 2px solid rgba(255, 255, 255, 0.25);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.task-label input:checked+.checkmark {
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%);
  border-color: transparent;
}

.task-label input:checked+.checkmark::after {
  content: '✓';
  color: var(--bg);
  font-size: 0.75rem;
  font-weight: 700;
}

.task-title {
  color: var(--text);
  font-size: 1rem;
  transition: color 0.2s;
}

.btn-edit {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-size: 1rem;
  line-height: 1;
  cursor: pointer;
  transition: color 0.2s, background 0.2s;
  flex-shrink: 0;
}

.btn-edit:hover {
  color: var(--accent);
  background: rgba(251, 191, 36, 0.12);
}

.btn-remove {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-size: 1.25rem;
  line-height: 1;
  cursor: pointer;
  transition: color 0.2s, background 0.2s;
  flex-shrink: 0;
}

.btn-remove:hover {
  color: #f87171;
  background: rgba(248, 113, 113, 0.1);
}

/* 编辑输入框：未聚焦时也有明显背景，聚焦时高亮 */
.task-edit-input {
  flex: 1;
  min-width: 0;
  padding: 0.5rem 0.75rem;
  border-radius: 10px;
  border: 2px solid rgba(251, 191, 36, 0.4);
  background: rgba(251, 191, 36, 0.12);
  color: var(--text);
  font-size: 1rem;
  outline: none;
  box-sizing: border-box;
  transition: background 0.2s, border-color 0.2s, box-shadow 0.2s;
}

.task-edit-input::placeholder {
  color: var(--text-muted);
}

.task-edit-input:focus {
  background: rgba(251, 191, 36, 0.2);
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(251, 191, 36, 0.35);
}

.empty-state {
  padding: 2.5rem 1rem;
  text-align: center;
}

.empty-state p {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.95rem;
}

/* ---------- 日期选择器（历史任务） ---------- */
.history-view {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  gap: 1.25rem;
}

.date-picker-sketch {
  flex-shrink: 0;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.date-picker-header {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 0.875rem;
}

.date-picker-select-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.date-picker-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.date-picker-select {
  width: 100%;
  padding: 0.5rem 2rem 0.5rem 0.75rem;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.06);
  color: var(--text);
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  outline: none;
  font-family: var(--font-display);
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23a1a1aa' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.6rem center;
  transition: border-color 0.2s, background-color 0.2s, box-shadow 0.2s;
}

.date-picker-select:hover {
  border-color: rgba(255, 255, 255, 0.25);
  background-color: rgba(255, 255, 255, 0.08);
}

.date-picker-select:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 2px rgba(251, 191, 36, 0.2);
}

.date-picker-select option {
  background: var(--surface);
  color: var(--text);
}

/* 日历网格：固定行高，保证 6 行日期完整显示 */
.date-picker-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  grid-auto-rows: 32px;
  gap: 5px;
}

.grid-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--text-muted);
  transition: background 0.2s, color 0.2s, border-color 0.2s, transform 0.1s;
}

.grid-cell-head {
  color: var(--text-muted);
  font-size: 0.75rem;
  font-weight: 600;
  cursor: default;
  background: rgba(255, 255, 255, 0.03);
}

.grid-cell-day {
  cursor: pointer;
  color: var(--text);
}

.grid-cell-day.is-current-month {
  color: var(--text);
  background: rgba(255, 255, 255, 0.06);
}

/* 非当前月份：浅色显示 */
.grid-cell-day.is-other-month {
  color: rgba(255, 255, 255, 0.18);
  background: rgba(255, 255, 255, 0.02);
}

/* 悬浮：金色渐变（仅当月日期，且未选中时） */
.grid-cell-day.is-current-month:hover:not(.is-selected) {
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%);
  border-color: transparent;
  color: var(--bg);
}

/* 选中：类似悬浮的金色渐变 + 背景色，用描边和略深渐变区分 */
.grid-cell-day.is-selected {
  background: linear-gradient(145deg, var(--accent-dark) 0%, var(--accent) 50%, var(--accent-dark) 100%);
  border: 2px solid rgba(251, 191, 36, 0.6);
  color: var(--bg);
  font-weight: 700;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.15) inset;
}

.grid-cell-day.is-selected:hover {
  background: linear-gradient(145deg, var(--accent-dark) 0%, var(--accent) 50%, var(--accent-dark) 100%);
  border-color: var(--accent);
  filter: brightness(1.06);
  color: var(--bg);
}

.grid-cell-day:active {
  transform: scale(0.96);
}
</style>
