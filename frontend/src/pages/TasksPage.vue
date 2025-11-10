<template>
  <!-- 中文注释：任务页面，包含统计、列表、创建/编辑、批量删除、番茄钟功能；支持下拉刷新 -->
  <div class="pull-refresh-wrapper" @touchstart="onTouchStart" @touchmove="onTouchMove" @touchend="onTouchEnd">
    <!-- 下拉刷新指示器（固定在顶部），拉动或刷新时淡入显示） -->
    <div class="fixed top-0 left-0 right-0 flex justify-center pointer-events-none" :style="{ opacity: (pullY>10||refreshing)?1:0 }">
      <div class="mt-2 text-xs text-gray-500 bg-white/80 rounded px-2 py-1 shadow">{{ refreshing ? '正在刷新...' : '下拉刷新' }}</div>
    </div>
    <div class="p-4 space-y-4" :style="{ transform: pullY ? ('translateY(' + pullY + 'px)') : 'none', transition: pulling ? 'none' : 'transform 0.2s ease' }">
    <!-- 顶部统计栏 -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <!-- 中文注释：头像优先使用用户自定义头像；无则使用默认头像，与“我的”页面保持一致 -->
          <el-avatar :size="36" :src="tasksAvatarSrc" />
          <div class="font-semibold">今日统计</div>
        </div>
      <!-- 中文注释：右侧取消回收站，改为彩色统计图标与金币显示 -->
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-1">
          <el-icon :size="20" style="color:#f59e0b"><Coin /></el-icon>
          <!-- 中文注释：总金币 = 已完成任务奖励金币总和 - 心愿兑换扣除金币 -->
          <span class="font-semibold">{{ totalCoins }}</span>
        </div>
        <el-icon :size="24" style="color:#ec4899"><DataAnalysis /></el-icon>
      </div>
    </div>

    <!-- 顶部统计：四项一行，不同颜色图标；下方单独大“统计”卡片居中显示 -->
    <div class="grid grid-cols-4 gap-2">
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#22c55e"><Clock /></el-icon>
          <div class="text-xs text-gray-500">日时长</div>
          <div class="font-semibold">{{ dayMinutes }} 分钟</div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#3b82f6"><List /></el-icon>
          <div class="text-xs text-gray-500">任务数</div>
          <div class="font-semibold">{{ completedTasksCount }}/{{ filteredTasks.length }}</div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#f59e0b"><Coin /></el-icon>
          <div class="text-xs text-gray-500">日金币</div>
          <div class="font-semibold">{{ dayCoins }}</div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#14b8a6"><CircleCheck /></el-icon>
          <div class="text-xs text-gray-500">完成率</div>
          <div class="font-semibold">{{ completeRate }}%</div>
        </div>
      </el-card>
    </div>
    <!-- 已将大“统计”图标移动到顶部右侧，此处删除 -->

    <!-- 日期选择与周视图：移动到“今日任务”卡片上方 -->
    <div class="my-2">
      <WeekCalendar v-model:selectedDate="selectedDate" :task-count-map="taskCountMap" />
    </div>

    <!-- 任务列表卡片 -->
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold">今日任务</span>
          <div class="space-x-2 flex items-center">
            <!-- 中文注释：筛选图标下拉菜单，点击选择“全部/已完成/待完成” -->
            <el-dropdown trigger="click" @command="onFilterCommand">
              <span class="el-dropdown-link">
                <el-icon class="cursor-pointer" :size="18" title="筛选"><Filter /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="全部">全部</el-dropdown-item>
                  <el-dropdown-item command="已完成">已完成</el-dropdown-item>
                  <el-dropdown-item command="待完成">待完成</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <!-- 中文注释：排序图标下拉菜单，点击选择排序方式（默认/时间顺序/任务分类/完成优先/添加时间） -->
            <el-dropdown trigger="click" @command="onSortCommand">
              <span class="el-dropdown-link">
                <el-icon class="cursor-pointer" :size="18" title="排序"><Sort /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="默认排序">默认排序</el-dropdown-item>
                  <el-dropdown-item command="时间顺序">时间顺序</el-dropdown-item>
                  <el-dropdown-item command="任务分类">任务分类</el-dropdown-item>
                  <el-dropdown-item command="完成优先">完成优先</el-dropdown-item>
                  <el-dropdown-item command="添加时间">添加时间</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
        <!-- 需求：取消按钮显示（保留逻辑，隐藏视图） -->
        <el-button v-if="false" type="primary" @click="openCreate">添加任务</el-button>
          </div>
        </div>
      </template>
      <!-- 日期选择与周视图已移动到卡片上方 -->
      <!-- 分类筛选：全部学科/语文/数学/英语 -->
      <div class="my-2">
        <el-radio-group v-model="categoryFilter" size="small">
          <el-radio-button label="全部任务" />
          <el-radio-button label="语文" />
          <el-radio-button label="数学" />
          <el-radio-button label="英语" />
        </el-radio-group>
      </div>

      <el-empty v-if="filteredTasks.length===0" description="暂无任务，快去添加吧~" />
      <div v-else class="space-y-4">
        <!-- 按分类分组显示 -->
        <div v-for="group in groupedTasks" :key="group.category" class="space-y-3">
          <div class="text-base font-semibold text-green-700">{{ group.category }}</div>
          <el-card
            v-for="t in group.items"
            :key="t.id"
            shadow="never"
            class="relative border border-gray-300 hover:ring-1 hover:ring-blue-300 transition"
            :class="{ 'ring-2 ring-blue-500': activeTaskId === t.id }"
            @click="activeTaskId = t.id"
          >
            <!-- 中文注释：自定义圆形复选框，居中于第一行与第二行之间，略大，点击切换完成状态 -->
            <div class="absolute left-2 top-1/2 -translate-y-1/2">
              <div
                class="w-5 h-5 rounded-full border-2 flex items-center justify-center cursor-pointer"
                :class="t.status===2 ? 'bg-green-500 border-green-500 text-white' : 'border-gray-400 text-gray-400'"
                @click="() => onCheckComplete(t, t.status !== 2)"
                title="点击切换完成状态"
              >
                <el-icon :size="12">
                  <CircleCheck />
                </el-icon>
              </div>
            </div>
            <!-- 第一行：左侧任务名，右侧状态与番茄钟入口 + 菜单 -->
            <div class="flex items-center justify-between pl-6">
              <div class="flex items-center gap-1">
                <!-- 中文注释：番茄钟图标仅在未完成时显示，位于右侧“待完成”标签左侧，此处移除 -->
                <div class="font-semibold text-left" :class="{'text-gray-500': t.status === 2}">{{ t.name }}</div>
              </div>
              <div class="flex items-center gap-1">
                <!-- 中文注释：图片查看入口移动到“实际完成时间”左侧，避免顶部拥挤 -->
                <!-- 中文注释：右侧状态与操作区：备注图标 + 小喇叭 + 番茄钟/状态标签 -->
                <div class="flex items-center gap-1">
                  <!-- 备注图标：点击进入备注页，作用与菜单中的“备注”一致 -->
              <!-- 中文注释：备注入口图标（受开关控制）；关闭后不显示 -->
              <el-icon v-if="store.taskNotesEnabled" :size="16" class="cursor-pointer" title="备注" style="color:#f97316" @click="router.push(`/tasks/${t.id}/notes`)"><ChatDotRound /></el-icon>
                  <!-- 小喇叭：朗读任务（关闭朗读时隐藏），替换为📢表情 -->
                  <span v-if="store.speech.enabled" class="cursor-pointer select-none" title="朗读任务" style="font-size:16px; line-height:16px" @click="speakTask(t)">📢</span>
                  <!-- 番茄钟图标仅未完成时显示 -->
                  <img v-if="t.status !== 2" src="@/assets/tomato.png" alt="番茄钟" class="w-4 h-4 cursor-pointer" @click="openTomato(t)" />
                  <!-- 状态标签 -->
                  <el-tag v-if="t.status !== 2" type="danger" size="small">待完成</el-tag>
                  <el-tag v-else type="success" size="small">已完成</el-tag>
                </div>
                <el-dropdown trigger="click" @command="(cmd)=>onMenu(cmd, t)">
                  <span class="el-dropdown-link">
                    <el-icon class="cursor-pointer"><MoreFilled /></el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="tomato">
                        <el-icon class="mr-1"><Clock /></el-icon>番茄钟
                      </el-dropdown-item>
                      <el-dropdown-item command="edit">
                        <el-icon class="mr-1"><Edit /></el-icon>编辑
                      </el-dropdown-item>
                      <!-- 新增：备注入口 -->
            <!-- 中文注释：备注菜单项（受开关控制）；关闭后不显示 -->
            <el-dropdown-item v-if="store.taskNotesEnabled" command="notes">
              <el-icon class="mr-1" style="color:#f97316"><ChatDotRound /></el-icon>备注
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" style="color:#f56c6c">
                        <el-icon class="mr-1"><Delete /></el-icon>删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>

            <!-- 第二行：左侧备注/描述；右侧实际/计划/金币（实际精确到秒） -->
            <div class="flex items-center justify-between mt-1 pl-6">
              <div class="text-xs text-gray-500 truncate max-w-[60%] text-left">{{ t.remark || t.description }}</div>
              <div class="flex items-center gap-3 text-xs">
                <!-- 中文注释：无论是否完成，只要有图片就显示图标；点击打开查看器（强制橙色避免主题覆盖） -->
            <el-icon v-if="hasImages(t)" class="cursor-pointer" :size="14" title="查看图片" style="color:#F97316 !important" @click="openTaskImages(t)"><Picture /></el-icon>
                <!-- 中文注释：仅在已完成时显示“实际完成时间”，位于图片图标与计划用时之间 -->
                <template v-if="t.status===2">
            <div class="flex items-center gap-1 text-blue-600 text-xs" title="实际完成时间">
                    <el-icon :size="14"><Clock /></el-icon>
                    <span class="font-semibold">{{ formatHMS(actualSecondsLocal[t.id] ?? ((t.actual_minutes||0)*60)) }}</span>
                  </div>
                </template>
                <div class="flex items-center gap-1 text-green-600 text-xs" title="计划用时">
                  <el-icon :size="14"><List /></el-icon>
                  <span class="font-semibold">{{ t.plan_minutes || 0 }} 分</span>
                </div>
                <div class="flex items-center gap-1 text-amber-600 text-xs" title="金币">
                  <el-icon :size="14"><Coin /></el-icon>
                  <span class="font-semibold">{{ t.score || 0 }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </el-card>

    <!-- 创建/编辑对话框：自定义头部图标与标题 -->
    <el-dialog v-model="formVisible" :width="dialogWidth">
      <template #header>
        <div class="flex items-center gap-2">
          <el-icon class="text-green-600"><Plus /></el-icon>
          <span class="font-semibold">{{ editing ? '编辑任务' : '创建任务' }}</span>
        </div>
      </template>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <!-- 基础信息分区 -->
        <el-card shadow="never" class="section-card mb-3">
        <!-- 任务标题 -->
        <el-form-item prop="name" required>
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><Edit /></el-icon>
              <span>任务标题</span>
            </div>
          </template>
          <el-input v-model="form.name" maxlength="128" show-word-limit style="width: 100%" />
        </el-form-item>

        <!-- 描述 -->
        <el-form-item label="任务描述" prop="description">
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><List /></el-icon>
              <span>任务描述</span>
            </div>
          </template>
          <el-input v-model="form.description" type="textarea" style="width: 100%" />
        </el-form-item>

        <!-- 任务图片上传（组件实现：多选、缩略图预览与进度） -->
        <el-form-item class="image-upload">
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><Plus /></el-icon>
              <span>任务图片</span>
            </div>
          </template>
          <TaskImageUploader
            :editing="editing"
            :user-id="userId"
            :task-id="currentTask?.id"
            v-model:serverPaths="form.images"
            v-model:localFiles="form.local_images"
            @added="() => ElMessage.success('图片已添加')"
          />
          <!-- 中文注释：根据需求移除上传说明文案 -->
        </el-form-item>

        <!-- 分类（每行一个字段） -->
        <el-form-item prop="category" required>
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><List /></el-icon>
              <span>任务分类</span>
            </div>
          </template>
          <el-select v-model="form.category" placeholder="选择分类" style="width: 100%">
            <el-option label="语文" value="语文" />
            <el-option label="数学" value="数学" />
            <el-option label="英语" value="英语" />
          </el-select>
        </el-form-item>

        <!-- 计划时长（每行一个字段） -->
        <el-form-item prop="plan_minutes" required>
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><Clock /></el-icon>
              <span>计划时长</span>
            </div>
          </template>
          <el-input-number v-model="form.plan_minutes" :min="1" :max="240" style="width: 100%" />
        </el-form-item>
        </el-card>

        <!-- 计划与重复分区 -->
        <el-card shadow="never" class="section-card mb-3">
        <!-- 任务金币（每行一个字段） -->
        <el-form-item prop="score">
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><Coin /></el-icon>
              <span>任务金币</span>
            </div>
          </template>
          <el-input-number v-model="form.score" :min="-10" :max="10" style="width: 100%" />
        </el-form-item>

        <!-- 重复类型（每行一个字段） -->
        <el-form-item prop="repeat_type">
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><List /></el-icon>
              <span>重复类型</span>
            </div>
          </template>
          <el-select v-model="form.repeat_type" placeholder="选择重复类型" style="width: 100%">
            <el-option label="无" value="none" />
            <el-option label="每天" value="daily" />
            <el-option label="每个工作日" value="weekdays" />
            <el-option label="每周" value="weekly" />
            <el-option label="每月" value="monthly" />
          </el-select>
        </el-form-item>

        <!-- 每周重复时选择星期（移动到重复类型之后） -->
        <el-form-item v-if="form.repeat_type==='weekly'" prop="weekly_days">
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><List /></el-icon>
              <span>选择星期</span>
            </div>
          </template>
          <el-checkbox-group v-model="form.weekly_days">
            <el-checkbox :label="1">周一</el-checkbox>
            <el-checkbox :label="2">周二</el-checkbox>
            <el-checkbox :label="3">周三</el-checkbox>
            <el-checkbox :label="4">周四</el-checkbox>
            <el-checkbox :label="5">周五</el-checkbox>
            <el-checkbox :label="6">周六</el-checkbox>
            <el-checkbox :label="7">周日</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <!-- 开始/截止日期（每行一个字段） -->
        <el-form-item prop="start_date" required>
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><Clock /></el-icon>
              <span>开始日期</span>
            </div>
          </template>
          <el-date-picker v-model="form.start_date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item prop="end_date">
          <template #label>
            <div class="flex items-center gap-1">
              <el-icon><Clock /></el-icon>
              <span>截止日期</span>
            </div>
          </template>
          <el-date-picker v-model="form.end_date" type="date" style="width: 100%" />
        </el-form-item>

        
        </el-card>

        <!-- 任务图片已上移至描述之后，此处删除分区卡片 -->
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="formVisible=false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 删除任务对话框（非重复任务：仅取消/确定） -->
    <el-dialog v-model="simpleDeleteDialogVisible" :width="isMobile ? '92vw' : '400px'">
      <template #header>
        <div class="flex items-center gap-2">
          <el-icon class="text-red-500"><Delete /></el-icon>
          <span class="font-semibold">删除任务</span>
        </div>
      </template>
      <div class="text-gray-700">确认删除任务「{{ deleteTarget?.name }}」？</div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="simpleDeleteDialogVisible=false">取消</el-button>
          <el-button type="danger" @click="doDeleteSimple">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 删除任务对话框（重复系列三种范围选择） -->
    <el-dialog v-model="deleteDialogVisible" :width="isMobile ? '92vw' : '480px'">
      <template #header>
        <div class="flex items-center gap-2">
          <el-icon class="text-red-500"><Delete /></el-icon>
          <span class="font-semibold">删除任务</span>
        </div>
      </template>
      <div class="space-y-3">
        <div class="text-gray-700">请选择删除范围：</div>
        <el-radio-group v-model="deleteScope">
          <el-radio label="current">仅删除当前日程</el-radio>
          <el-radio label="future">删除当前及未来全部日程</el-radio>
          <el-radio label="all">删除所有日程</el-radio>
        </el-radio-group>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="deleteDialogVisible=false">取消</el-button>
          <el-button type="danger" @click="deleteScope==='current' ? doDeleteCurrent(deleteTarget!) : doDeleteSeries(deleteScope as 'future'|'all')">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 回收站对话框 -->
    <el-dialog v-model="recycleVisible" title="回收站" width="600px">
      <el-table :data="recycleList" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="任务" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="restore([row.id])">恢复</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="recycleVisible=false">关闭</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 番茄钟组件弹窗 -->
    <!-- 中文注释：番茄钟弹窗标题：🍅 番茄钟：任务名称 -->
    <el-dialog v-model="tomatoVisible" :title="`🍅 番茄钟：${currentTask?.name || '番茄钟'}`" width="520px" :show-close="true" :close-on-click-modal="!store.tomato.fixedTomatoPage">
      <!-- 中文注释：番茄钟默认使用任务计划时长，若无则 20 分钟；支持倒计时/正计时、预设与自定义 -->
      <TomatoTimer :work-minutes="currentTask?.plan_minutes || 20" :break-minutes="5" :task-name="currentTask?.name" :task-remark="currentTask?.remark || currentTask?.description" @complete="onTomatoComplete" />
    </el-dialog>

    <!-- 中文注释：任务图片全屏查看（覆盖式），支持缩放与左右翻看 -->
    <el-image-viewer v-if="imagesViewerVisible" :url-list="imageViewerList" :initial-index="imageViewerIndex" @close="imagesViewerVisible=false" />

    <!-- 中文注释：结束内部内容容器（p-4 space-y-4），避免顶层 wrapper 未闭合） -->
    </div>

    <!-- 右下角绿色加号浮动按钮：创建任务 -->
  <el-button
      type="success"
      circle
      class="fixed bottom-20 right-6 shadow-lg"
      @click="openCreate"
      title="创建任务"
    >
      <!-- 中文注释：创建任务图标放大 0.5 倍（默认约 16px → 24px） -->
      <el-icon :size="24"><Plus /></el-icon>
    </el-button>

    <!-- 悬浮番茄钟：位于底部导航上方居中，水位填充动画展示进度 -->
    <div
      v-if="store.tomato.running && !store.tomato.fixedTomatoPage"
      class="fixed bottom-16 left-1/2 -translate-x-1/2 z-50"
    >
      <div class="w-14 h-14 rounded-full shadow-lg cursor-pointer overflow-hidden relative bg-blue-100"
           @click="tomatoVisible=true">
        <div class="absolute bottom-0 left-0 right-0 bg-blue-500 bg-opacity-40 transition-all"
             :style="{ height: fillPercent + '%' }"></div>
        <!-- 中文注释：悬浮球中央显示 mm:ss，小号文字；提高对比度便于看清 -->
        <div class="absolute inset-0 flex items-center justify-center text-blue-900 font-semibold text-[10px]">{{ floatingTime }}</div>
      </div>
    </div>
  </div>
 </template>

<script setup lang="ts">
// 中文注释：任务页面逻辑，统一使用服务层 API，实现表单校验与错误提示
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Clock, List, Coin, CircleCheck, MoreFilled, DataAnalysis, Edit, Delete, Filter, ChatDotRound, Sort } from '@element-plus/icons-vue'
import defaultAvatar from '@/assets/avatars/default.png'
import { useAuth } from '@/stores/auth'
import { useAppState } from '@/stores/appState'
import router from '@/router'
import TomatoTimer from '@/components/TomatoTimer.vue'
import WeekCalendar from '@/components/WeekCalendar.vue'
import TaskImageUploader from '@/components/TaskImageUploader.vue'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
dayjs.extend(utc)
import { listTasks, createTask, updateTask, updateTaskStatus, deleteTask, completeTomato, listRecycleBin, restoreTasks, uploadTaskImage, batchDelete, type TaskItem } from '@/services/tasks'
import { Picture } from '@element-plus/icons-vue'
import { prepareUpload } from '@/utils/image'
import { speak } from '@/utils/speech'
const isMobile = ref(false)
const userId = 1 // 中文注释：示例用户ID（参考心愿页做法，后续接入登录）
const dialogWidth = computed(() => (isMobile.value ? '96vw' : '640px'))

// ===== 下拉刷新逻辑（移动端触摸） =====
const pulling = ref(false) // 是否正在拉动
const pullY = ref(0) // 下拉位移
const startY = ref(0)
const refreshing = ref(false)
const pullThreshold = 150

function onTouchStart(e: TouchEvent) {
  // 仅在页面滚动到顶部时允许下拉刷新
  if (window.scrollY > 0) return
  const t = e.touches[0]
  startY.value = t.clientY
  pullY.value = 0
  pulling.value = true
}

function onTouchMove(e: TouchEvent) {
  if (!pulling.value) return
  const t = e.touches[0]
  const dy = t.clientY - startY.value
  if (dy > 0) {
    // 防止浏览器默认下拉回弹影响
    e.preventDefault()
    pullY.value = Math.min(dy, 120)
  } else {
    pullY.value = 0
  }
}

async function onTouchEnd() {
  if (!pulling.value) return
  pulling.value = false
  if (pullY.value >= pullThreshold) {
    try {
      refreshing.value = true
      await fetchTasks()
    } finally {
      refreshing.value = false
    }
  }
  pullY.value = 0
}

// 顶部统计占位（后续与后端联动）
const store = useAppState()
const auth = useAuth()
// 中文注释：总金币改为直接读取全局 store.coins（由后端任务完成/取消与心愿兑换实时更新），与心愿页保持一致
const totalCoins = computed(() => store.coins)
const completedTasksCount = computed(() => {
  return filteredTasks.value.filter(t => t.status === 2).length
})
const dayCoins = computed(() => {
  return filteredTasks.value.filter((t) => t.status === 2).reduce((sum, t) => sum + (t.score || 0), 0)
})

// 中文注释：朗读任务内容（格式："{任务分类}，{任务标题}，备注：{任务描述}"）
function speakTask(t: TaskItem) {
  try {
    if (!store.speech.enabled) {
      ElMessage.info('朗读已关闭，可在“我的 → 朗读设置”开启')
      return
    }
    const category = (t.category || '').trim()
    const title = (t.name || '').trim()
    const remark = (t.remark || t.description || '').trim()
    const text = `${category ? category + '，' : ''}${title}${remark ? '，备注：' + remark : ''}`
    const ok = speak(text, { voiceURI: store.speech.voiceURI || undefined, rate: store.speech.rate, pitch: store.speech.pitch })
    if (!ok) ElMessage.warning('当前浏览器不支持朗读或语音不可用')
  } catch {
    ElMessage.error('朗读失败，请稍后重试')
  }
}
const dayMinutes = computed(() => {
  return filteredTasks.value.reduce((sum, t) => sum + (t.actual_minutes || 0), 0)
})
const completeRate = computed(() => {
  if (filteredTasks.value.length === 0) return 0
  return Math.round((completedTasksCount.value / filteredTasks.value.length) * 100)
})

// 列表与筛选
const tasks = ref<TaskItem[]>([])
const filter = ref<'全部' | '已完成' | '待完成'>('全部')
const categoryFilter = ref<'全部任务' | '语文' | '数学' | '英语'>('全部任务')
const selectedDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const taskCountMap = computed<Record<string, number>>(() => {
  const map: Record<string, number> = {}
  for (const t of tasks.value) {
    const sDate = t.start_date ? dayjs(t.start_date).toDate() : null
    const eDate = t.end_date ? dayjs(t.end_date).toDate() : undefined
    if (!sDate) continue
    const rawRepeat = (t as any).repeat || 'none'
    const rep = String(rawRepeat).toLowerCase()
    const type: 'none'|'daily'|'weekdays'|'weekly'|'monthly' =
      /none|无|^$/i.test(rep) ? 'none' :
      /daily|每天/i.test(rep) ? 'daily' :
      /weekdays|工作日/i.test(rep) ? 'weekdays' :
      /weekly|每周/i.test(rep) ? 'weekly' :
      /monthly|每月/i.test(rep) ? 'monthly' : 'none'
    if (type === 'none' || !eDate) {
      const key = dayjs(sDate).format('YYYY-MM-DD')
      map[key] = (map[key] || 0) + 1
      continue
    }
    const dow = dayjs(sDate).day() === 0 ? 7 : dayjs(sDate).day()
    const weeklyDays: number[] = Array.isArray((t as any).weekly_days) ? ((t as any).weekly_days as number[]) : [dow]
    const dates = generateRepeatDates(sDate, eDate, type, weeklyDays)
    for (const d of dates) {
      const key = dayjs(d).format('YYYY-MM-DD')
      map[key] = (map[key] || 0) + 1
    }
  }
  return map
})
const filteredTasks = computed(() => {
  let result = tasks.value
  if (filter.value === '已完成') result = result.filter((t) => t.status === 2)
  else if (filter.value === '待完成') result = result.filter((t) => t.status !== 2)
  if (categoryFilter.value !== '全部任务') result = result.filter((t) => (t.category || '') === categoryFilter.value)
  // 中文注释：日期过滤，使用重复规则生成发生日期，匹配选中日期
  const selKey = dayjs(selectedDate.value).format('YYYY-MM-DD')
  result = result.filter((t) => {
    const sDate = t.start_date ? dayjs(t.start_date).toDate() : null
    const eDate = t.end_date ? dayjs(t.end_date).toDate() : undefined
    if (!sDate) return false
    const rawRepeat = (t as any).repeat || 'none'
    const rep = String(rawRepeat).toLowerCase()
    const type: 'none'|'daily'|'weekdays'|'weekly'|'monthly' =
      /none|无|^$/i.test(rep) ? 'none' :
      /daily|每天/i.test(rep) ? 'daily' :
      /weekdays|工作日/i.test(rep) ? 'weekdays' :
      /weekly|每周/i.test(rep) ? 'weekly' :
      /monthly|每月/i.test(rep) ? 'monthly' : 'none'
    if (type === 'none' || !eDate) {
      return dayjs(sDate).format('YYYY-MM-DD') === selKey
    }
    const dow = dayjs(sDate).day() === 0 ? 7 : dayjs(sDate).day()
    const weeklyDays: number[] = Array.isArray((t as any).weekly_days) ? ((t as any).weekly_days as number[]) : [dow]
    const dates = generateRepeatDates(sDate, eDate, type, weeklyDays)
    const keys = new Set(dates.map((d) => dayjs(d).format('YYYY-MM-DD')))
    return keys.has(selKey)
  })
  return result
})
// 中文注释：排序模式（默认/时间顺序/任务分类/完成优先/添加时间）
type SortMode = '默认排序' | '时间顺序' | '任务分类' | '完成优先' | '添加时间'
const sortMode = ref<SortMode>('默认排序')

// 中文注释：根据排序模式对任务排序；在分组前执行，保证分组顺序也受影响
const sortedTasks = computed(() => {
  const arr = [...filteredTasks.value]
  const byDateAsc = (a: TaskItem, b: TaskItem) => {
    const ad = a.start_date ? dayjs(a.start_date).valueOf() : 0
    const bd = b.start_date ? dayjs(b.start_date).valueOf() : 0
    if (ad !== bd) return ad - bd
    return (a.id || 0) - (b.id || 0)
  }
  const byCategory = (a: TaskItem, b: TaskItem) => {
    const ac = (a.category || '未分类')
    const bc = (b.category || '未分类')
    if (ac !== bc) return ac.localeCompare(bc)
    return byDateAsc(a, b)
  }
  const byCompletedFirst = (a: TaskItem, b: TaskItem) => {
    if ((a.status === 2) !== (b.status === 2)) return a.status === 2 ? -1 : 1
    return byDateAsc(a, b)
  }
  const byAddedTime = (a: TaskItem, b: TaskItem) => {
    // 中文注释：后端暂未提供 created_at，使用 id 作为添加时间近似（id 越大越新）
    return (b.id || 0) - (a.id || 0)
  }
  if (sortMode.value === '时间顺序') return arr.sort(byDateAsc)
  if (sortMode.value === '任务分类') return arr.sort(byCategory)
  if (sortMode.value === '完成优先') return arr.sort(byCompletedFirst)
  if (sortMode.value === '添加时间') return arr.sort(byAddedTime)
  return arr
})

// 中文注释：按分类分组，便于移动端展示与筛选

// 中文注释：解析头像地址（与“我的”页一致，仅当包含 uploads/ 或完整 URL 时使用后端路径）
function resolveAvatarUrl(p?: string | null) {
  if (!p) return defaultAvatar
  const s = String(p)
  if (/^https?:\/\//i.test(s)) return s
  if (!/uploads\//i.test(s)) return defaultAvatar
  return `/api/${s}`.replace(/\/+/g, '/').replace(/\/$/, '')
}
const tasksAvatarSrc = computed(() => resolveAvatarUrl(auth.user?.avatar_path))
const groupedTasks = computed(() => {
  const map = new Map<string, TaskItem[]>()
  for (const t of sortedTasks.value) {
    const cat = t.category || '未分类'
    if (!map.has(cat)) map.set(cat, [])
    map.get(cat)!.push(t)
  }
  let groups = Array.from(map.entries()).map(([category, items]) => ({ category, items }))
  // 中文注释：当选择“任务分类”排序时，按分类名升序排序分组
  if (sortMode.value === '任务分类') {
    groups = groups.sort((a, b) => a.category.localeCompare(b.category))
  }
  return groups
})

// 中文注释：本地保存每个任务的实际秒数，便于前端精确展示（后端以分钟存储）
const actualSecondsLocal = reactive<Record<number, number>>({})
// 中文注释：格式化为“x时x分x秒”，按需显示（仅显示非零单位；全为零则显示0秒）
function formatHMS(sec: number) {
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  const s = Math.floor(sec % 60)
  const parts: string[] = []
  if (h > 0) parts.push(`${h}时`)
  if (m > 0) parts.push(`${m}分`)
  if (s > 0 || parts.length === 0) parts.push(`${s}秒`)
  return parts.join('')
}

// 选择与批量删除：移动端不再提供批量删除，这里移除选择逻辑

// 创建/编辑表单
const formVisible = ref(false)
const recycleVisible = ref(false)
const tomatoVisible = ref(false)
const editing = ref(false)
const currentTask = ref<TaskItem | null>(null)
const formRef = ref<FormInstance>()
const form = reactive<any>({ name: '', description: '', category: '语文', score: 1, plan_minutes: 20, start_date: new Date(), end_date: undefined, images: [], local_images: [] })
// 中文注释：上传列表与草稿逻辑已迁移到组件中处理
const rules: FormRules = {
  name: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  score: [{ required: true, message: '请输入金币', trigger: 'change' }],
  plan_minutes: [{ required: true, message: '请输入计划时长', trigger: 'change' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }]
}

// 中文注释：移除未使用的状态文案/样式函数，简化页面逻辑

async function fetchTasks() {
  try {
    const res = await listTasks()
    tasks.value = res.items || []
    // 中文注释：日时长、日金币、完成率均改为统计当日任务
    // dayMinutes.value = filteredTasks.value.reduce((sum, t) => sum + (t.actual_minutes || 0), 0)
    // dayCoins.value = filteredTasks.value.filter((t) => t.status === 2).reduce((sum, t) => sum + (t.score || 0), 0)
    // 中文注释：completeRate 已改为计算属性，无需手动赋值
    // completeRate.value = tasks.value.length ? Math.round((tasks.value.filter((t) => t.status === 2).length / tasks.value.length) * 100) : 0
    // 中文注释：同步更新全局 coins（考虑心愿扣减），心愿页面读取该值作为可用金币
    store.setCoins(totalCoins.value)
  } catch (e: any) {
    // 中文注释：增强错误提示，优先展示后端返回的业务错误信息
    const msg = e?.response?.data?.message || e?.message || e
    console.error('任务列表加载失败诊断', {
      message: e?.message,
      status: e?.response?.status,
      payload: e?.response?.data,
    })
    ElMessage.error(`任务列表加载失败：${msg}`)
  }
}

function openCreate() {
  // 中文注释：跳转到独立创建页面，提升移动端体验与布局灵活性
  router.push('/tasks/create')
}

function openEdit(t: TaskItem) {
  // 中文注释：跳转到独立编辑页面，按任务ID加载与保存
  router.push(`/tasks/${t.id}/edit`)
}

function resolveUploadUrl(rel: string) {
  const base = (import.meta as any).env.VITE_API_BASE || ''
  return `${base}/api/${rel}`.replace(/\/$/, '')
}

// 中文注释：判断任务是否有图片
// 中文注释：解析任务的图片 JSON，兼容字符串二次编码的情况（例如 "[\"a\"]"）
function parseImageList(json?: string): string[] {
  if (!json || !json.trim()) return []
  try {
    const first = JSON.parse(json)
    if (Array.isArray(first)) return first as string[]
    if (typeof first === 'string') {
      try {
        const second = JSON.parse(first)
        return Array.isArray(second) ? (second as string[]) : []
      } catch { return [] }
    }
    return []
  } catch {
    // 中文注释：当后端返回非标准字符串时，兜底为空数组，避免前端报错
    return []
  }
}

// 中文注释：判断任务是否有图片
function hasImages(t: TaskItem) {
  const arr = parseImageList(t.image_json)
  return Array.isArray(arr) && arr.length > 0
}

// 中文注释：任务图片查看对话框状态与打开方法
const imagesViewerVisible = ref(false)
const imageViewerList = ref<string[]>([])
const imageViewerIndex = ref(0)
function openTaskImages(t: TaskItem) {
  const rels = parseImageList(t.image_json)
  imageViewerList.value = rels.map(resolveUploadUrl)
  if (imageViewerList.value.length > 0) {
    imageViewerIndex.value = 0
    imagesViewerVisible.value = true
  } else {
    ElMessage.info('该任务暂无图片')
  }
}

async function submitForm() {
  try {
    await formRef.value?.validate()
  } catch { return }
  try {
    if (editing.value && currentTask.value) {
      await updateTask(currentTask.value.id, {
        name: form.name,
        description: form.description,
        category: form.category,
        score: form.score,
        plan_minutes: form.plan_minutes,
        start_date: form.start_date,
        end_date: form.end_date,
        image_json: JSON.stringify(form.images || [])
      })
      ElMessage.success('任务已更新')
    } else {
      // 中文注释：根据重复类型批量创建任务实例
      const dates = generateRepeatDates(form.start_date, form.end_date, form.repeat_type, form.weekly_days)
      if (dates.length === 0) {
        // 无重复或未设置截止日期则创建单个
        dates.push(form.start_date)
      }
      const createdTasks: TaskItem[] = []
      for (const d of dates) {
        const t = await createTask({
          user_id: userId,
          name: form.name,
          description: form.description,
          category: form.category,
          score: form.score,
          plan_minutes: form.plan_minutes,
          start_date: d,
          end_date: undefined
        })
        createdTasks.push(t)
      }
      // 中文注释：创建后按任务ID上传本地图片，并写入 image_json
      if ((form.local_images || []).length > 0) {
        for (const t of createdTasks) {
          const paths: string[] = []
          for (const f of form.local_images) {
            try {
              const webp = await prepareUpload(f as File)
              // 中文注释：前端调试日志，确认文件对象与元信息
              console.debug('准备上传任务图片', {
                task_id: t.id,
                filename: (webp as File)?.name,
                size: (webp as File)?.size,
                type: (webp as File)?.type,
                isFile: webp instanceof File,
              })
              const { path } = await uploadTaskImage(userId, webp, t.id)
              paths.push(path)
            } catch (err: any) {
              // 中文注释：详细前端错误日志，包含任务ID、文件名与后端返回信息
              console.error('上传任务图片失败', {
                task_id: t.id,
                filename: (f as File)?.name,
                message: err?.message || err,
                response: err?.response?.data
              })
              ElMessage.error(`图片上传失败：${(f as File)?.name || ''} → ${err?.response?.data?.message || err?.message || '未知错误'}`)
            }
          }
          if (paths.length > 0) {
            await updateTask(t.id, { image_json: JSON.stringify(paths) })
            console.info('任务图片已更新到 image_json', { task_id: t.id, count: paths.length })
          }
        }
      }
      ElMessage.success(`任务已创建${dates.length>1?`（${dates.length}条）`:''}`)
    }
    // 中文注释：创建成功后清理草稿（组件已管理缩略图）
    try { localStorage.removeItem('task_draft_images') } catch {}
    formVisible.value = false
    await fetchTasks()
  } catch (e: any) {
    // 中文注释：提交失败时输出更详细的诊断信息
    console.error('提交任务失败', {
      message: e?.message || e,
      response: e?.response?.data,
    })
    ElMessage.error(`提交失败：${e?.response?.data?.message || e?.message || e}`)
  }
}

// 勾选即完成：只允许从未完成 -> 已完成，不提供取消
async function onCheckComplete(t: TaskItem, checked: boolean) {
  try {
    if (checked) {
      // 中文注释：勾选为完成：按计划时长计入实际，并标记为已完成
      const planM = t.plan_minutes || 20
      await updateTask(t.id, { actual_minutes: planM })
      const resp: any = await updateTaskStatus(t.id, 2)
      t.status = 2
      t.actual_minutes = (t.actual_minutes || 0) + planM
      actualSecondsLocal[t.id] = planM * 60
      if (resp && typeof resp.user_coins !== 'undefined') store.setCoins(Number(resp.user_coins))
      ElMessage.success('已标记为完成（按计划时长计）')
    } else {
      // 中文注释：取消完成：标记为未完成，并从日金币与总金币中扣除该任务金币
      const resp: any = await updateTaskStatus(t.id, 0)
      t.status = 0
      if (resp && typeof resp.user_coins !== 'undefined') store.setCoins(Number(resp.user_coins))
      ElMessage.success('已取消完成，金币已扣除')
    }
    // 统一刷新统计
    // dayMinutes.value = filteredTasks.value.reduce((sum, x) => sum + (x.actual_minutes || 0), 0)
    // dayCoins.value = filteredTasks.value.filter((x) => x.status === 2).reduce((sum, x) => sum + (x.score || 0), 0)
    // 中文注释：completeRate 已改为计算属性，无需手动赋值
    // completeRate.value = tasks.value.length ? Math.round((tasks.value.filter((x) => x.status === 2).length / tasks.value.length) * 100) : 0
  } catch (e: any) {
    ElMessage.error(`状态变更失败：${e.message || e}`)
  }
}
// 取消切换状态功能：保留空函数避免引用错误（模板已移除）

// 删除对话框状态（重复任务支持范围选择）
const deleteDialogVisible = ref(false)
const deleteScope = ref<'current'|'future'|'all'>('current')
const deleteTarget = ref<TaskItem | null>(null)
// 中文注释：非重复任务删除确认对话框（仅取消/确定）
const simpleDeleteDialogVisible = ref(false)

function isRepeatedTask(t: TaskItem): boolean {
  const rep = String((t as any).repeat || '').trim()
  if (!!t.series_id || !!t.end_date || rep !== '') return true
  // 中文注释：后端暂未持久化 repeat/series_id 时，按“名称+分类”在列表中出现多次来判断为重复系列
  const sameGroup = tasks.value.filter((x) => x.name === t.name && (x.category || '') === (t.category || ''))
  return sameGroup.length >= 2
}

function confirmDelete(t: TaskItem) {
  if (isRepeatedTask(t)) {
    deleteTarget.value = t
    deleteScope.value = 'current'
    deleteDialogVisible.value = true
    return
  }
  deleteTarget.value = t
  simpleDeleteDialogVisible.value = true
}

async function doDeleteCurrent(t: TaskItem) {
  try {
    await deleteTask(t.id)
    ElMessage.success('已删除当前日程')
    deleteDialogVisible.value = false
    await fetchTasks()
  } catch (e: any) {
    ElMessage.error(`删除失败：${e.message || e}`)
  }
}

async function doDeleteSimple() {
  if (!deleteTarget.value) return
  try {
    await deleteTask(deleteTarget.value.id)
    ElMessage.success('已删除')
    simpleDeleteDialogVisible.value = false
    await fetchTasks()
  } catch (e: any) {
    ElMessage.error(`删除失败：${e.message || e}`)
  }
}

async function doDeleteSeries(scope: 'future'|'all') {
  if (!deleteTarget.value) return
  const target = deleteTarget.value
  try {
    // 中文注释：按 series_id 分组，若无则用“名称 + 分类”兜底
    const group = tasks.value.filter((x) => {
      if (target.series_id) return x.series_id === target.series_id
      return x.name === target.name && (x.category || '') === (target.category || '')
    })
    // 仅当前及未来：筛选 start_date >= 选中日期；全部：整组
    let candidates = group
    if (scope === 'future') {
      const th = dayjs(selectedDate.value).startOf('day')
      candidates = group.filter((x) => dayjs(x.start_date).startOf('day').isSame(th) || dayjs(x.start_date).startOf('day').isAfter(th))
    }
    const ids = candidates.map((x) => x.id)
    if (ids.length === 0) { ElMessage.info('未找到可删除的系列任务'); return }
    await batchDelete(ids)
    ElMessage.success(`已删除${scope==='all'?'整个系列':'当前及未来'}共 ${ids.length} 条`)
    deleteDialogVisible.value = false
    await fetchTasks()
  } catch (e: any) {
    ElMessage.error(`批量删除失败：${e.message || e}`)
  }
}



const recycleList = ref<TaskItem[]>([])
async function restore(ids: number[]) {
  try {
    await restoreTasks(ids)
    ElMessage.success('已恢复')
    recycleList.value = await listRecycleBin()
    await fetchTasks()
  } catch (e: any) {
    ElMessage.error(`恢复失败：${e.message || e}`)
  }
}

function openTomato(t: TaskItem) {
  // 中文注释：跳转到独立番茄钟页面
  router.push(`/tasks/${t.id}/tomato`)
}

async function onTomatoComplete(seconds?: number) {
  if (!currentTask.value) return
  try {
    // 中文注释：按实际秒数精确展示，后端按分钟上报（四舍五入）；无秒数则按计划时长
    const usedSec = Math.max(1, seconds || (currentTask.value.plan_minutes || 20) * 60)
    const reportMinutes = Math.max(1, Math.round(usedSec / 60))
    const updated = await completeTomato(currentTask.value.id, reportMinutes)
    const idx = tasks.value.findIndex((x) => x.id === currentTask.value!.id)
    if (idx >= 0) tasks.value[idx] = updated
    actualSecondsLocal[currentTask.value.id] = usedSec
    // 完成后标记任务为已完成
    await updateTaskStatus(currentTask.value.id, 2)
    if (idx >= 0) tasks.value[idx].status = 2
    // 中文注释：dayMinutes 已改为计算属性，无需手动赋值
    // dayMinutes.value = tasks.value.reduce((sum, x) => sum + (x.actual_minutes || 0), 0)
    ElMessage.success('番茄钟完成，数据已记录')
    tomatoVisible.value = false
  } catch (e: any) {
    ElMessage.error(`番茄上报失败：${e.message || e}`)
  }
}

onMounted(() => {
  fetchTasks()
  const updateMobile = () => { isMobile.value = window.innerWidth < 768 }
  updateMobile()
  window.addEventListener('resize', updateMobile)
})

// 菜单命令统一处理
function onMenu(cmd: string, t: TaskItem) {
  if (cmd === 'tomato') return openTomato(t)
  if (cmd === 'edit') return openEdit(t)
  if (cmd === 'notes') return router.push(`/tasks/${t.id}/notes`)
  if (cmd === 'delete') return confirmDelete(t)
}

// 中文注释：筛选图标下拉菜单命令处理，更新状态筛选条件
function onFilterCommand(cmd: '全部' | '已完成' | '待完成') {
  filter.value = cmd
}

// 中文注释：排序图标下拉菜单命令处理，更新排序模式
function onSortCommand(cmd: SortMode) {
  sortMode.value = cmd
}

// 中文注释：任务图片上传逻辑已迁移到组件中

// ===== 重复日期生成逻辑 =====
function generateRepeatDates(start: Date, end: Date | undefined, type: 'none' | 'daily' | 'weekdays' | 'weekly' | 'monthly', weeklyDays: number[]) {
  const out: Date[] = []
  if (!end || type === 'none') return out
  const s = dayjs(start).startOf('day')
  const e = dayjs(end).startOf('day')
  if (e.isBefore(s)) return out
  if (type === 'daily') {
    let d = s.clone()
    while (!d.isAfter(e)) { out.push(d.toDate()); d = d.add(1, 'day') }
  } else if (type === 'weekdays') {
    let d = s.clone()
    while (!d.isAfter(e)) {
      const w = d.day() // 0-周日 ... 6-周六
      if (w >= 1 && w <= 5) out.push(d.toDate())
      d = d.add(1, 'day')
    }
  } else if (type === 'weekly') {
    const set = new Set(weeklyDays || [])
    let d = s.clone()
    while (!d.isAfter(e)) {
      const w = d.day() === 0 ? 7 : d.day()
      if (set.has(w)) out.push(d.toDate())
      d = d.add(1, 'day')
    }
  } else if (type === 'monthly') {
    // 中文注释：按每月同一日生成（若当月无该日，例如 31 日，则跳过）
    let d = s.clone()
    const dayOfMonth = s.date()
    while (!d.isAfter(e)) {
      const candidate = d.date(dayOfMonth)
      if (candidate.month() === d.month() && !candidate.isAfter(e)) out.push(candidate.toDate())
      d = d.add(1, 'month')
    }
  }
  return out
}

// 中文注释：悬浮球填充百分比（正计时用已用时 / 目标时长，倒计时用剩余时间）
const fillPercent = computed(() => {
  const dur = store.tomato.durationMinutes * 60
  const sec = store.tomato.remainingSeconds
  if (store.tomato.mode === 'countup') {
    return Math.min(100, Math.round((sec / dur) * 100))
  }
  return Math.min(100, Math.round(((dur - sec) / dur) * 100))
})

// 中文注释：悬浮球显示的时间文本（mm:ss），倒计时显示剩余，正计时显示累计
const floatingTime = computed(() => {
  const sec = Math.max(0, store.tomato.remainingSeconds || 0)
  const mm = String(Math.floor(sec / 60)).padStart(2, '0')
  const ss = String(sec % 60).padStart(2, '0')
  return `${mm}:${ss}`
})

// 中文注释：当前选中的任务卡片ID，用于高亮显示
const activeTaskId = ref<number | null>(null)

// 中文注释：移除未使用的函数（toggleStatus、openRecycle），消除编译器警告
</script>

<style scoped>
/* 中文注释：基本页面样式，响应式栅格布局已通过 Tailwind 实现 */
/* 中文注释：分类筛选按钮的文字加粗，增强可读性 */
.el-radio-button__inner { font-weight: 600; }
/* 中文注释：缩小图片上传卡片与缩略图尺寸，保证整齐美观 */
.image-upload :deep(.el-upload--picture-card) {
  width: 96px;
  height: 96px;
}
.image-upload :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 96px;
  height: 96px;
}
/* 中文注释：表单项间距优化与色彩区分（轻微灰色分隔） */
:deep(.el-form-item) {
  margin-bottom: 12px;
}
/* 中文注释：统一标签区左对齐，图标与文字同一行对齐显示 */
:deep(.el-form-item__label) {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}
/* 中文注释：输入区占满剩余空间，宽度随容器自适应 */
:deep(.el-form-item__content) {
  flex: 1;
  min-width: 0;
}
/* 中文注释：分区卡片样式优化，提升视觉层次与对比度 */
.section-card {
  border: 1px solid #f0f0f0;
  border-radius: 8px;
}
.section-card :deep(.el-card__header) {
  font-weight: 600;
}
</style>
