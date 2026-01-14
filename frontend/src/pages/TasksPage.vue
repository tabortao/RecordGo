<template>
  <!-- 中文注释：任务页面，包含统计、列表、创建/编辑、批量删除、番茄钟功能；支持下拉刷新 -->
  <div ref="wrapperRef" class="pull-refresh-wrapper" @touchstart="onTouchStart" @touchmove="onTouchMove" @touchend="onTouchEnd" @touchcancel="onTouchCancel" style="overscroll-behavior-y: contain; touch-action: pan-y;">
    <!-- 下拉刷新指示器（固定在顶部），拉动或刷新时淡入显示） -->
    <div class="fixed top-0 left-0 right-0 flex justify-center pointer-events-none" :style="{ opacity: (pullY>10||refreshing)?1:0 }">
      <div class="mt-2 text-xs text-gray-500 bg-white/80 rounded px-2 py-1 shadow">{{ refreshing ? '正在刷新...' : '下拉刷新' }}</div>
    </div>
    <div class="fixed top-0 left-0 right-0 bg-white/80 dark:bg-gray-800/80 backdrop-blur z-40 border-b border-gray-200 dark:border-gray-700">
      <div class="px-4 py-2 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <el-dropdown trigger="click" @command="onAvatarCommand">
            <span class="el-dropdown-link">
              <el-avatar :size="36" :src="tasksAvatarSrc" class="cursor-pointer" />
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="noop" class="font-semibold" style="pointer-events: none; cursor: default">切换用户</el-dropdown-item>
                <el-dropdown-item v-for="acc in auth.accounts" :key="acc.user.id" :command="'switch:' + acc.user.id">
                  <div class="flex items-center gap-2">
                    <el-avatar :size="24" :src="accountAvatarSrc(acc.user)" />
                    <span>{{ (acc.user.nickname || '').trim() || acc.user.username }}</span>
                    <el-icon v-if="auth.user?.id === acc.user.id" :size="16" style="color:#22c55e"><CircleCheck /></el-icon>
                  </div>
                </el-dropdown-item>
                <el-dropdown-item divided command="add">添加新用户</el-dropdown-item>
                <el-dropdown-item command="logout" style="color:#f56c6c">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <div class="font-semibold">任务统计</div>
        </div>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-1">
            <el-icon :size="20" style="color:#f59e0b"><Coin /></el-icon>
            <span class="font-semibold">{{ totalCoins }}</span>
          </div>
          <el-icon :size="24" style="color:#ec4899" class="cursor-pointer" @click="router.push('/tasks/stats')"><DataAnalysis /></el-icon>
        </div>
      </div>
    </div>
    <div class="h-14"></div>
    <div class="p-4 space-y-4" :style="{ transform: pullY ? ('translateY(' + pullY + 'px)') : 'none', transition: pulling ? 'none' : 'transform 0.2s ease' }">
    

    <!-- 顶部统计：四项一行，不同颜色图标；下方单独大“统计”卡片居中显示 -->
    <div class="grid grid-cols-4 gap-1">
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon :size="19" style="color:#22c55e"><Clock /></el-icon>
          <div class="text-xs text-gray-500 dark:text-gray-400">日时长</div>
          <el-tooltip :content="tipMinutes" placement="top">
            <div class="font-bold" style="color:#22c55e">{{ dayMinutes }}</div>
          </el-tooltip>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#3b82f6"><List /></el-icon>
          <div class="text-xs text-gray-500 dark:text-gray-400">任务数</div>
          <el-tooltip :content="tipTasks" placement="top">
            <div class="font-bold" style="color:#3b82f6">{{ completedTasksCount }}/{{ filteredTasks.length }}</div>
          </el-tooltip>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon :size="19" style="color:#f59e0b"><Money /></el-icon>
          <div class="text-xs text-gray-500 dark:text-gray-400">日金币</div>
          <el-tooltip :content="tipCoins" placement="top">
            <div class="font-bold" style="color:#f59e0b">{{ dayCoins }}</div>
          </el-tooltip>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon :size="19" style="color:#14b8a6"><CircleCheck /></el-icon>
          <div class="text-xs text-gray-500 dark:text-gray-400">完成率</div>
          <el-tooltip :content="tipRate" placement="top">
            <div class="font-bold" style="color:#14b8a6">{{ completeRate }}%</div>
          </el-tooltip>
        </div>
      </el-card>
    </div>
    <!-- 已将大“统计”图标移动到顶部右侧，此处删除 -->

    <!-- 日期选择与周视图：移动到“今日任务”卡片上方 -->
    <div class="my-2">
      <WeekCalendar v-model:selectedDate="selectedDate" :task-count-map="taskCountMap" />
    </div>

    <!-- 任务列表卡片 -->
    <el-card shadow="never" class="no-frame">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold">{{ headerLabel }}</span>
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
      <!-- 分类筛选：与“任务分类设置”一致的动态列表 -->
      <div class="mt-0 mb-6">
        <el-radio-group v-model="categoryFilter" size="small">
          <el-radio-button label="全部任务">全部任务</el-radio-button>
          <el-radio-button v-for="c in categoriesForDay" :key="c.name" :label="c.name">
            <span class="inline-flex items-center gap-1">
              <span class="inline-block w-2 h-2 rounded" :style="{ backgroundColor: c.color }"></span>
              <span>{{ c.name }}</span>
            </span>
          </el-radio-button>
        </el-radio-group>
      </div>

      <el-empty v-if="filteredTasks.length===0" description="暂无任务，快去添加吧~" />
      <div v-else class="space-y-4">
        <!-- 按分类分组显示 -->
        <div v-for="group in groupedTasks" :key="group.category" class="space-y-3">
          <!-- 中文注释：分组标题左侧展示分类颜色，颜色与设置保持一致 -->
          <div class="text-base font-semibold flex items-center gap-2 text-gray-900 dark:text-gray-100">
            <span class="inline-block w-2 h-2 rounded" :style="{ backgroundColor: categoryColor(group.category) }"></span>
            <span>{{ group.category }}</span>
          </div>
          <el-card
            v-for="t in group.items"
            :key="t.id"
            shadow="never"
            class="relative border border-gray-300 dark:border-gray-700 hover:ring-1 hover:ring-blue-300 dark:hover:ring-blue-200/30 transition rounded-xl mx-1"
            :data-task-id="t.id"
            :class="{ 'ring-2 ring-blue-500': activeTaskId === t.id }"
            @click="activeTaskId = t.id"
          >
            <!-- 中文注释：自定义圆形复选框，居中于第一行与第二行之间，略大，点击切换完成状态 -->
            <div class="absolute left-2 top-1/2 -translate-y-1/2">
              <div
                class="w-5 h-5 rounded-full border-2 flex items-center justify-center cursor-pointer"
                :class="isCompletedOnSelected(t) ? 'bg-green-500 border-green-500 text-white' : 'border-gray-400 dark:border-gray-500 text-gray-400 dark:text-gray-500'"
                @click="() => onCheckComplete(t, !isCompletedOnSelected(t))"
                title="点击切换完成状态"
              >
                <el-icon :size="12">
                  <CircleCheck />
                </el-icon>
              </div>
            </div>
            <!-- 第一行：左侧任务名，右侧状态与番茄钟入口 + 菜单 -->
            <div class="flex items-center justify-between pl-9">
              <div class="flex items-center gap-2">
                <!-- 中文注释：番茄钟图标仅在未完成时显示，位于右侧“待完成”标签左侧，此处移除 -->
                <div class="font-semibold text-left" :class="{'text-gray-500': isCompletedOnSelected(t)}">{{ t.name }}</div>
              </div>
              <div class="flex items-center gap-1">
                <!-- 中文注释：图片查看入口移动到“实际完成时间”左侧，避免顶部拥挤 -->
                <!-- 中文注释：右侧状态与操作区：备注图标 + 小喇叭 + 番茄钟/状态标签 -->
                <div class="flex items-center gap-2">
                  <!-- 备注图标：点击进入备注页，作用与菜单中的“备注”一致 -->
              <!-- 中文注释：备注入口图标（受开关控制）；关闭后不显示 -->
              <el-icon v-if="isVIP && store.taskNotesEnabled" :size="16" class="cursor-pointer" title="备注" style="color:#f97316" @click="router.push(`/tasks/${t.id}/notes`)"><ChatDotRound /></el-icon>
                  <!-- 小喇叭：朗读任务（关闭朗读时隐藏），替换为📢表情 -->
                  <el-icon v-if="store.speech.enabled" :size="14" class="cursor-pointer select-none" title="朗读任务" @click="speakTask(t)"><Headset /></el-icon>
                  <!-- 番茄钟图标仅未完成时显示 -->
                  <img v-if="!isCompletedOnSelected(t)" src="@/assets/tomato.png" alt="番茄钟" class="w-4 h-4 cursor-pointer" @click="openTomato(t)" />
                  <!-- 状态标签 -->
                  <el-tag v-if="!isCompletedOnSelected(t)" type="danger" size="small">待完成</el-tag>
                  <el-tag v-else type="success" size="small">已完成</el-tag>
                </div>
                <el-dropdown trigger="click" @command="(cmd: string)=>onMenu(cmd, t)">
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
            <el-dropdown-item v-if="isVIP && store.taskNotesEnabled" command="notes">
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
            <div class="flex items-center justify-between mt-1 pl-9">
              <div class="text-xs text-gray-500 truncate max-w-[60%] text-left">{{ t.remark || t.description }}</div>
              <div class="flex items-center gap-3 text-xs">
                <!-- 中文注释：无论是否完成，只要有图片就显示图标；点击打开查看器（强制橙色避免主题覆盖） -->
            <el-icon v-if="hasImages(t)" class="cursor-pointer" :size="14" title="查看图片" style="color:#F97316 !important" @click="openTaskImages(t)"><Picture /></el-icon>
                <!-- 中文注释：仅在已完成时显示“实际完成时间”，位于图片图标与计划用时之间 -->
                <template v-if="isCompletedOnSelected(t)">
            <div class="flex items-center gap-1 text-blue-600 dark:text-blue-400 text-xs" title="实际完成时间">
                    <el-icon :size="14"><Clock /></el-icon>
                    <span class="font-semibold">{{ formatHMS(getActualSeconds(t)) }}</span>
                  </div>
                </template>
                <div class="flex items-center gap-1 text-green-600 dark:text-green-400 text-xs" title="计划用时">
                  <el-icon :size="14"><List /></el-icon>
                  <span class="font-semibold">{{ t.plan_minutes || 0 }} 分</span>
                </div>
                <div class="flex items-center gap-1 text-amber-600 dark:text-amber-500 text-xs" title="金币">
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
          <el-icon class="text-green-600"><CirclePlusFilled /></el-icon>
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
          <el-date-picker v-model="form.end_date" type="date" style="width: 100%" :disabled="form.repeat_type==='none'" />
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
      <el-table :data="recycleList" :style="{ width: '100%' }">
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
      v-if="isParent || canTaskCreate"
      type="success"
      circle
      class="fixed no-pull"
      :style="{ left: fabPos.x + 'px', top: fabPos.y + 'px', backgroundColor: '#22c55e', borderColor: '#22c55e', zIndex: 60 }"
      @mousedown="onFabDown"
      @touchstart="onFabTouchStart"
      @click="openCreate"
      title="创建任务"
    >
      <el-icon :size="46" class="text-white"><CirclePlusFilled/></el-icon>
    </el-button>

    <!-- 中文注释：移除旧版任务页悬浮番茄钟，改用新的全局悬浮球（右下角橙色），避免重复显示 -->
    <el-dialog v-model="addUserVisible" width="360px">
      <template #header>
        <div class="font-semibold">添加新用户</div>
      </template>
      <el-form label-width="90px">
        <template v-if="isParent">
          <el-form-item label="用户名"><el-input v-model="addUserName" /></el-form-item>
          <el-form-item label="密码"><el-input v-model="addUserPassword" type="password" /></el-form-item>
        </template>
        <template v-else>
          <el-form-item label="子账号令牌"><el-input v-model="addUserToken" placeholder="请输入子账号令牌" /></el-form-item>
        </template>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="addUserVisible=false">取消</el-button>
          <el-button type="primary" @click="doAddUser">添加并登录</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
 </template>

<script setup lang="ts">
// 中文注释：任务页面逻辑，统一使用服务层 API，实现表单校验与错误提示
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus/es/components/form'
import { Plus, Clock, List, Coin, Money, CircleCheck, MoreFilled, DataAnalysis, Edit, Delete, Filter, ChatDotRound, Sort, Headset, CirclePlusFilled } from '@element-plus/icons-vue'
import defaultAvatar from '@/assets/avatars/default.png'
import { useAuth } from '@/stores/auth'
import { useAppState } from '@/stores/appState'
import { usePermissions } from '@/composables/permissions'
import router from '@/router'
import { apiLogin, apiTokenLogin } from '@/services/auth'
import TomatoTimer from '@/components/TomatoTimer.vue'
import WeekCalendar from '@/components/WeekCalendar.vue'
import TaskImageUploader from '@/components/TaskImageUploader.vue'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
dayjs.extend(utc)
import { listTasks, createTask, updateTask, updateTaskStatus, deleteTask, completeTomato, listRecycleBin, restoreTasks, uploadTaskImage, type TaskItem, syncOfflineTasks, listTaskOccurrences, completeOccurrence, uncompleteOccurrence, deleteOccurrence } from '@/services/tasks'
import { normalizeUploadPath } from '@/services/wishes'
import { Picture } from '@element-plus/icons-vue'
import { prepareUpload } from '@/utils/image'
import { speak } from '@/utils/speech'
import { useTaskCategories } from '@/stores/categories'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import http from '@/services/http'
import { isAbortError } from '@/services/http'
const isMobile = ref(false)
// 中文注释：接入认证状态获取真实用户ID（未登录回退为 0）
const auth = useAuth()
const isVIP = computed(() => {
  const u = auth.user
  if (!u) return false
  const lifetime = !!(u as any).is_lifetime_vip
  const vip = !!(u as any).is_vip
  const expire = (u as any).vip_expire_time ? new Date((u as any).vip_expire_time as string) : null
  const valid = lifetime || (vip && !!expire && expire.getTime() > Date.now())
  return valid
})
const userId = computed(() => auth.user?.id ?? 0)
// 中文注释：解析权限，父账号默认放行；子账号按动作校验
const { isParent, canTaskCreate, canTaskEdit, canTaskDelete, canTaskStatus } = usePermissions()
const dialogWidth = computed(() => (isMobile.value ? '96vw' : '640px'))
// 中文注释：任务分类 Store，用于动态筛选与分组颜色
const cats = useTaskCategories()
function categoryColor(name: string) { return cats.colorOf(name) }
// 中文注释：按日期与完成筛选得到当日可见任务（不含分类筛选），用于动态生成分类筛选项
const dateStatusFilteredTasks = computed(() => {
  let result = tasks.value
  if (filter.value === '已完成') result = result.filter((t) => isCompletedOnSelected(t))
  else if (filter.value === '待完成') result = result.filter((t) => !isCompletedOnSelected(t))
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
// 中文注释：仅展示当日存在任务的分类；若任务包含“未分类”，也纳入筛选项
import type { TaskCategory } from '@/stores/categories'
const categoriesForDay = computed(() => {
  const present = new Set<string>()
  for (const t of dateStatusFilteredTasks.value) present.add((t.category || '未分类'))
  const list: TaskCategory[] = cats.list().filter((c: TaskCategory) => present.has(c.name))
  // 兼容任务出现未在设置页定义的分类
  for (const name of Array.from(present)) {
    if (!list.some((c: TaskCategory) => c.name === name)) list.push({ name, color: cats.colorOf(name), order: cats.orderOf(name) } as TaskCategory)
  }
  // 按自定义顺序排序
  return list.sort((a: TaskCategory, b: TaskCategory) => {
    const oa = cats.orderOf(a.name)
    const ob = cats.orderOf(b.name)
    if (oa !== ob) return oa - ob
    return a.name.localeCompare(b.name)
  })
})

// ===== 下拉刷新逻辑（移动端触摸） =====
const pulling = ref(false)
const pullY = ref(0)
const startY = ref(0)
const refreshing = ref(false)
const pullThreshold = 80
const wrapperRef = ref<HTMLElement | null>(null)

function onTouchStart(e: TouchEvent) {
  // 仅在页面滚动到顶部时允许下拉刷新
  const top = (document.scrollingElement?.scrollTop || window.scrollY || document.documentElement.scrollTop || 0)
  if (top > 0) return
  const target = e.target as HTMLElement
  if (target && target.closest('.no-pull')) return
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

function onTouchCancel() {
  pulling.value = false
  pullY.value = 0
}

async function onTouchEnd() {
  if (!pulling.value) return
  pulling.value = false
  if (pullY.value >= pullThreshold) {
    try {
      refreshing.value = true
      try { await cats.syncFromServer() } catch {}
      await fetchTasks()
      await fetchOccurrences()
      // 同步刷新总金币（父子金币共享时返回父账号金币）
      try { await http.get('/coins') } catch {}
    } finally {
      refreshing.value = false
      try { ElMessage.success('已刷新') } catch {}
    }
  }
  pullY.value = 0
}

// 顶部统计占位（后续与后端联动）
const store = useAppState()
// 中文注释：总金币改为直接读取全局 store.coins（由后端任务完成/取消与心愿兑换实时更新），与心愿页保持一致
const totalCoins = computed(() => store.coins)
const occurMap = ref<Record<number, { status: number; minutes?: number }>>({})
function isRepeatTask(t: TaskItem) {
  const rep = String((t as any).repeat || (t as any).repeat_type || 'none').trim().toLowerCase()
  const type = /^(daily|weekdays|weekly|monthly)$/.test(rep) ? rep : 'none'
  // 中文注释：只要设置了重复类型，即视为重复任务（无论是否设置截止日期），确保点击完成时调用 completeOccurrence 而非修改主状态
  return type !== 'none'
}
function isCompletedOnSelected(t: TaskItem) {
  if (isRepeatTask(t)) return (occurMap.value[t.id]?.status || 0) === 2
  return t.status === 2
}
function getActualSeconds(t: TaskItem): number {
  if (isRepeatTask(t)) {
    const m = occurMap.value[t.id]?.minutes || 0
    if (m > 0) return m * 60
    return (t.plan_minutes || 0) * 60
  }
  const sec = (actualSecondsLocal[t.id] ?? ((t.actual_minutes || 0) * 60))
  return sec || ((t.plan_minutes || 0) * 60)
}
const completedTasksCount = computed(() => {
  return filteredTasks.value.filter((t) => isCompletedOnSelected(t)).length
})
const dayCoins = computed(() => {
  return filteredTasks.value.filter((t) => isCompletedOnSelected(t)).reduce((sum, t) => sum + (t.score || 0), 0)
})

// 中文注释：朗读任务内容（格式："{任务分类}，{任务标题}，备注：{任务描述}"）
async function speakTask(t: TaskItem) {
  try {
    if (!store.speech.enabled) {
      ElMessage.info('朗读已关闭，可在“我的 → 朗读设置”开启')
      return
    }
    const category = (t.category || '').trim()
    const title = (t.name || '').trim()
    const remark = (t.remark || t.description || '').trim()
    const text = `${category ? category + '，' : ''}${title}${remark ? '，备注：' + remark : ''}`
    const ok = await speak(text, { voiceURI: store.speech.voiceURI || undefined, rate: store.speech.rate, pitch: store.speech.pitch })
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

const tipMinutes = computed(() => `今日已完成任务总用时${dayMinutes.value}分钟`)
const tipTasks = computed(() => `今日已完成任务${completedTasksCount.value}个，总任务${filteredTasks.value.length}个`)
const tipCoins = computed(() => `今日总共获得${dayCoins.value}个金币`)
const tipRate = computed(() => `今日任务完成比例为${completeRate.value}%`)

// 列表与筛选
const tasks = ref<TaskItem[]>([])
const filter = ref<'全部' | '已完成' | '待完成'>('全部')
// 中文注释：分类筛选使用字符串，完全由“任务分类设置”提供
const categoryFilter = ref<string>('全部任务')
const selectedDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const headerLabel = computed(() => {
  const d = dayjs(selectedDate.value)
  return d.isSame(dayjs(), 'day') ? '今日任务' : `${d.month() + 1}月${d.date()}日任务`
})
const taskCountMap = computed<Record<string, number>>(() => {
  const map: Record<string, number> = {}
  const base: any = dayjs(selectedDate.value)
  const weekday = base.day()
  const monday = base.subtract((weekday === 0 ? 6 : weekday - 1), 'day').startOf('day')
  const days: any[] = []
  for (let i = 0; i < 7; i++) { days.push(monday.add(i, 'day')) }
  for (const t of tasks.value) {
    const sDate = t.start_date ? dayjs(t.start_date) : null
    const eDate = t.end_date ? dayjs(t.end_date) : undefined
    if (!sDate) continue
    const s: any = sDate
    const rep = String((t as any).repeat || 'none').toLowerCase()
    const type: 'none'|'daily'|'weekdays'|'weekly'|'monthly' =
      /none|无|^$/i.test(rep) ? 'none' :
      /daily|每天/i.test(rep) ? 'daily' :
      /weekdays|工作日/i.test(rep) ? 'weekdays' :
      /weekly|每周/i.test(rep) ? 'weekly' :
      /monthly|每月/i.test(rep) ? 'monthly' : 'none'
    const startKey = s.format('YYYY-MM-DD')
    if (type === 'none') { map[startKey] = (map[startKey] || 0) + 1; continue }
    const dowStart = s.day() === 0 ? 7 : s.day()
    const weeklyDays: number[] = Array.isArray((t as any).weekly_days) ? ((t as any).weekly_days as number[]) : [dowStart]
    if (eDate) {
      const dates = generateRepeatDates(s.toDate(), (eDate as any).toDate(), type, weeklyDays)
      for (const d of dates) { const key = dayjs(d).format('YYYY-MM-DD'); map[key] = (map[key] || 0) + 1 }
      continue
    }
    for (const d of days) {
      if (d.isBefore(s.startOf('day'))) continue
      const w = d.day() === 0 ? 7 : d.day()
      if (type === 'daily') { const key = d.format('YYYY-MM-DD'); map[key] = (map[key] || 0) + 1; continue }
      if (type === 'weekdays') { if (w >= 1 && w <= 5) { const key = d.format('YYYY-MM-DD'); map[key] = (map[key] || 0) + 1 } continue }
      if (type === 'weekly') {
        const sw = dowStart
        const ok = (weeklyDays && weeklyDays.length) ? weeklyDays.includes(w) : (w === sw)
        if (ok) { const key = d.format('YYYY-MM-DD'); map[key] = (map[key] || 0) + 1 }
        continue
      }
      if (type === 'monthly') { if (d.date() === s.date()) { const key = d.format('YYYY-MM-DD'); map[key] = (map[key] || 0) + 1 } continue }
    }
  }
  return map
})
const filteredTasks = computed(() => {
  let result = tasks.value
  if (filter.value === '已完成') result = result.filter((t) => isCompletedOnSelected(t))
  else if (filter.value === '待完成') result = result.filter((t) => !isCompletedOnSelected(t))
  if (categoryFilter.value !== '全部任务') result = result.filter((t) => (t.category || '') === categoryFilter.value)
  // 中文注释：日期过滤，使用重复规则生成发生日期，匹配选中日期
  const selKey = dayjs(selectedDate.value).format('YYYY-MM-DD')
  result = result.filter((t) => {
    const occ = occurMap.value[t.id]?.status
    if (occ === -1) return false
    const sDateStr = t.start_date ? String(t.start_date) : ''
    const eDateStr = t.end_date ? String(t.end_date) : ''
    if (!sDateStr) return false
    const sDate = dayjs(sDateStr).toDate()
    const eDate = eDateStr ? dayjs(eDateStr).toDate() : undefined
    const rawRepeat = (t as any).repeat || 'none'
    const rep = String(rawRepeat).toLowerCase()
    const type: 'none'|'daily'|'weekdays'|'weekly'|'monthly' =
      /none|无|^$/i.test(rep) ? 'none' :
      /daily|每天/i.test(rep) ? 'daily' :
      /weekdays|工作日/i.test(rep) ? 'weekdays' :
      /weekly|每周/i.test(rep) ? 'weekly' :
      /monthly|每月/i.test(rep) ? 'monthly' : 'none'
    if (type === 'none') {
      return dayjs(sDate).format('YYYY-MM-DD') === selKey
    }
    const dow = dayjs(sDate).day() === 0 ? 7 : dayjs(sDate).day()
    const weeklyDays: number[] = Array.isArray((t as any).weekly_days) ? ((t as any).weekly_days as number[]) : [dow]
    if (!eDate) {
      const d = dayjs(selectedDate.value)
      const w = d.day() === 0 ? 7 : d.day()
      const sd = dayjs(sDate)
      if (d.isBefore(sd.startOf('day'))) return false
      if (type === 'daily') return true
      if (type === 'weekdays') return w >= 1 && w <= 5
      if (type === 'weekly') {
        const set = new Set(weeklyDays || [])
        const sw = sd.day() === 0 ? 7 : sd.day()
        return set.size ? set.has(w) : (w === sw)
      }
      if (type === 'monthly') return d.date() === sd.date()
      return false
    }
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
    const ac = isCompletedOnSelected(a)
    const bc = isCompletedOnSelected(b)
    if (ac !== bc) return ac ? -1 : 1
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
const tasksAvatarSrc = ref<string>(defaultAvatar)
async function updateTasksAvatar() {
  const p = auth.user?.avatar_path
  if (!p) { tasksAvatarSrc.value = defaultAvatar; return }
  const s = String(p)
  if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) { tasksAvatarSrc.value = defaultAvatar; return }
  if (/^https?:\/\//i.test(s)) { tasksAvatarSrc.value = s; return }
  const base = getStaticBase()
  if (/uploads\//i.test(s)) { tasksAvatarSrc.value = `${base}/api/${s.replace(/^\/+/, '')}`; return }
  try { tasksAvatarSrc.value = await presignView(s) } catch { tasksAvatarSrc.value = defaultAvatar }
}
onMounted(updateTasksAvatar)
watch(() => auth.user?.avatar_path, () => { updateTasksAvatar() })

const accountsAvatarMap = ref<Record<number,string>>({})
async function resolveAccountsAvatars(list: { user: { id: number, avatar_path?: string|null } }[]) {
  const base = getStaticBase()
  for (const a of list) {
    const p = a.user?.avatar_path
    const s = String(p || '')
    if (!s) { accountsAvatarMap.value[a.user.id] = defaultAvatar; continue }
    if (/^https?:\/\//i.test(s)) { accountsAvatarMap.value[a.user.id] = s; continue }
    if (/uploads\//i.test(s)) { accountsAvatarMap.value[a.user.id] = `${base}/api/${s.replace(/^\/+/, '')}`; continue }
    try { accountsAvatarMap.value[a.user.id] = await presignView(s) } catch { accountsAvatarMap.value[a.user.id] = defaultAvatar }
  }
}
function accountAvatarSrc(u: { id: number, avatar_path?: string|null }) {
  const s = String(u?.avatar_path || '')
  if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) return defaultAvatar
  if (/^https?:\/\//i.test(s)) return s
  if (/uploads\//i.test(s)) return `${getStaticBase()}/api/${s.replace(/^\/+/, '')}`
  return accountsAvatarMap.value[u.id] || defaultAvatar
}
onMounted(async () => { await resolveAccountsAvatars(auth.accounts || []) })
watch(() => auth.accounts.map((a: { user: { id: number; avatar_path?: string|null } }) => a.user.id + ':' + (a.user.avatar_path || '')), async () => { await resolveAccountsAvatars(auth.accounts || []) })
const addUserVisible = ref(false)
const addUserName = ref('')
const addUserPassword = ref('')
const addUserToken = ref('')

function onAvatarCommand(cmd: string) {
  if (cmd === 'noop') return
  if (cmd.startsWith('switch:')) {
    const id = Number(cmd.split(':')[1] || 0)
    if (id > 0) {
      if (!isParent.value) {
        try {
      const target = auth.accounts.find((a: { user: { id: number; parent_id?: number|null } }) => a.user.id === id)
          const targetIsParent = !target?.user?.parent_id || Number(target?.user?.parent_id) === 0
          if (targetIsParent) {
            ElMessage.warning('子账号不可切换到父账户，请使用令牌登录子账号')
            return
          }
        } catch {}
      }
      auth.switchAccount(id)
      try { store.setCoins(Number(auth.user?.coins ?? 0)) } catch {}
      fetchTasks()
    }
    return
  }
  if (cmd === 'add') {
    addUserVisible.value = true
    return
  }
  if (cmd === 'logout') {
    try {
      auth.logout()
      store.reset()
      router.replace('/login')
    } catch {
      location.reload()
    }
  }
}

async function doAddUser() {
  try {
    if (isParent.value) {
      if (!addUserName.value || !addUserPassword.value) { ElMessage.warning('请输入用户名和密码'); return }
      const resp = await apiLogin(addUserName.value, addUserPassword.value)
      auth.setLogin(resp.token, resp.user, true)
      try { store.setCoins(Number(resp.user?.coins ?? 0)) } catch {}
      addUserVisible.value = false
      addUserName.value = ''
      addUserPassword.value = ''
    } else {
      const token = addUserToken.value.trim()
      if (!token) { ElMessage.warning('请输入子账号令牌'); return }
      const resp = await apiTokenLogin(token)
      auth.setLogin(resp.token, resp.user, true)
      try { store.setCoins(Number(resp.user?.coins ?? 0)) } catch {}
      addUserVisible.value = false
      addUserToken.value = ''
    }
  } catch (e: any) {
    ElMessage.error(e?.message || '添加失败')
  }
}
const groupedTasks = computed(() => {
  const map = new Map<string, TaskItem[]>()
  for (const t of sortedTasks.value) {
    const cat = t.category || '未分类'
    if (!map.has(cat)) map.set(cat, [])
    map.get(cat)!.push(t)
  }
  let groups = Array.from(map.entries()).map(([category, items]) => ({ category, items }))
  // 中文注释：自动排序逻辑（默认排序模式且开启自动排序）：
  // 1）分类内已完成任务排在下方；
  // 2）全部完成的分类排在未完成分类之后。
  if (sortMode.value === '默认排序' && store.taskAutoSortEnabled) {
    groups = groups.map(g => {
      const unfinished = g.items.filter(i => i.status !== 2)
      const finished = g.items.filter(i => i.status === 2)
      // 分类内：未完在上，已完在下；每段内按开始时间升序
      const byDate = (a: TaskItem, b: TaskItem) => {
        const ad = a.start_date ? dayjs(a.start_date).valueOf() : 0
        const bd = b.start_date ? dayjs(b.start_date).valueOf() : 0
        if (ad !== bd) return ad - bd
        return (a.id || 0) - (b.id || 0)
      }
      return { category: g.category, items: [...unfinished.sort(byDate), ...finished.sort(byDate)] }
    })
    const hasUnfinished = (g: { category: string; items: TaskItem[] }) => g.items.some(i => i.status !== 2)
    const unfinishedGroups = groups.filter(hasUnfinished)
    const finishedOnlyGroups = groups.filter(g => !hasUnfinished(g))
    groups = [...unfinishedGroups, ...finishedOnlyGroups]
  }
  // 中文注释：当选择“任务分类”排序时，按分类名升序排序分组
  if (sortMode.value === '任务分类') {
    const ord = (name: string) => cats.orderOf(name)
    groups = groups.sort((a, b) => {
      const oa = ord(a.category)
      const ob = ord(b.category)
      if (oa !== ob) return oa - ob
      return a.category.localeCompare(b.category)
    })
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
    const res = await listTasks({ page_size: 1000 })
    tasks.value = res.items || []
    // 中文注释：日时长、日金币、完成率均改为统计当日任务
    // dayMinutes.value = filteredTasks.value.reduce((sum, t) => sum + (t.actual_minutes || 0), 0)
    // dayCoins.value = filteredTasks.value.filter((t) => t.status === 2).reduce((sum, t) => sum + (t.score || 0), 0)
    // 中文注释：completeRate 已改为计算属性，无需手动赋值
    // completeRate.value = tasks.value.length ? Math.round((tasks.value.filter((t) => t.status === 2).length / tasks.value.length) * 100) : 0
    // 中文注释：移除本地计算后的金币覆盖，统一由拦截器在后端返回时同步，避免状态不一致
  } catch (e: any) {
    // 中文注释：增强错误提示，优先展示后端返回的业务错误信息
    if (isAbortError(e)) { return }
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
  // 中文注释：权限校验：父账号允许；子账号需具备 tasks.create 权限
  if (!isParent.value && !canTaskCreate.value) {
    ElMessage.warning('当前权限不允许创建任务')
    return
  }
  // 中文注释：跳转到独立创建页面，提升移动端体验与布局灵活性
  router.push('/tasks/create')
}

function openEdit(t: TaskItem) {
  // 中文注释：权限校验：父账号允许；子账号需具备 tasks.edit 权限
  if (!isParent.value && !canTaskEdit.value) {
    ElMessage.warning('当前权限不允许编辑任务')
    return
  }
  // 中文注释：跳转到独立编辑页面，按任务ID加载与保存
  router.push(`/tasks/${t.id}/edit`)
}

async function resolveUploadUrlAsync(rel: string): Promise<string> {
  const base = getStaticBase()
  let p = normalizeUploadPath(rel)
  if (p.startsWith('uploads/')) return `${base}/api/${String(p).replace(/^\/+/, '')}`
  try { return await presignView(p) } catch { return '' }
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
async function openTaskImages(t: TaskItem) {
  const rels = parseImageList(t.image_json)
  const urls = await Promise.all(rels.map((r) => resolveUploadUrlAsync(r)))
  imageViewerList.value = urls.filter(Boolean)
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
      // 中文注释：编辑权限校验
      if (!isParent.value && !canTaskEdit.value) { ElMessage.warning('当前权限不允许编辑任务'); return }
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
      // 中文注释：创建权限校验
      if (!isParent.value && !canTaskCreate.value) { ElMessage.warning('当前权限不允许创建任务'); return }
      // 改为仅创建一条任务，重复规则由后端保存
      const t = await createTask({
        user_id: userId,
        name: form.name,
        description: form.description,
        category: form.category,
        score: form.score,
        plan_minutes: form.plan_minutes,
        start_date: form.start_date,
        end_date: form.end_date,
        repeat_type: form.repeat_type,
        weekly_days: form.weekly_days || []
      })
      if ((form.local_images || []).length > 0) {
        const paths: string[] = []
        for (const f of form.local_images) {
          try {
            const webp = await prepareUpload(f as File)
            const { path } = await uploadTaskImage(userId.value, webp, t.id)
            paths.push(path)
          } catch (err: any) {
            ElMessage.error(`图片上传失败：${(f as File)?.name || ''} → ${err?.response?.data?.message || err?.message || '未知错误'}`)
          }
        }
        if (paths.length > 0) { await updateTask(t.id, { image_json: JSON.stringify(paths) }) }
      }
      ElMessage.success('任务已创建')
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
  // 中文注释：状态变更权限校验
  if (!isParent.value && !canTaskStatus.value) {
    ElMessage.warning('当前权限不允许更改任务状态')
    return
  }
  try {
    if (checked) {
      // 中文注释：勾选为完成：按计划时长计入实际，并标记为已完成
      const planM = t.plan_minutes || 20
      const dateStr = dayjs(selectedDate.value).format('YYYY-MM-DD')
      if (isRepeatTask(t)) {
        await completeOccurrence(t.id, { date: dateStr, minutes: planM })
        occurMap.value[t.id] = { status: 2, minutes: planM }
        ElMessage.success('已标记为当次完成')
      } else {
        await updateTask(t.id, { actual_minutes: planM })
        await updateTaskStatus(t.id, 2)
        t.status = 2
        t.actual_minutes = (t.actual_minutes || 0) + planM
        actualSecondsLocal[t.id] = planM * 60
        ElMessage.success('已标记为完成（按计划时长计）')
      }
      celebrate(t.id)
    } else {
      // 中文注释：取消完成：标记为未完成，并从日金币与总金币中扣除该任务金币
      const dateStr = dayjs(selectedDate.value).format('YYYY-MM-DD')
      if (isRepeatTask(t)) {
        await uncompleteOccurrence(t.id, { date: dateStr })
        occurMap.value[t.id] = { status: 0 }
        ElMessage.success('已取消当次完成')
      } else {
        await updateTaskStatus(t.id, 0)
        t.status = 0
        ElMessage.success('已取消完成')
      }
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
function celebrate(taskId?: number) {
  const colors = ['#22c55e', '#3b82f6', '#f59e0b', '#ef4444', '#14b8a6']
  let rect: DOMRect | null = null
  if (taskId) {
    const el = document.querySelector(`[data-task-id="${taskId}"]`) as HTMLElement | null
    rect = el?.getBoundingClientRect() || null
  }
  const container = document.createElement('div')
  container.className = 'confetti-local'
  container.style.position = 'fixed'
  if (rect) {
    container.style.left = `${rect.left}px`
    container.style.top = `${rect.top}px`
    container.style.width = `${rect.width}px`
    container.style.height = `${rect.height}px`
  } else {
    container.style.left = '0'
    container.style.top = '0'
    container.style.width = '100vw'
    container.style.height = '100vh'
  }
  container.style.pointerEvents = 'none'
  container.style.zIndex = '9999'
  document.body.appendChild(container)
  const count = 28
  for (let i = 0; i < count; i++) {
    const piece = document.createElement('div')
    piece.className = 'confetti-burst'
    const x = Math.random() * 100
    const y = Math.random() * 100
    const tx = (Math.random() * 60 - 30)
    const ty = (Math.random() * 60 - 10)
    const rotate = (Math.random() * 180 - 90).toFixed(0)
    const color = colors[Math.floor(Math.random() * colors.length)]
    piece.style.left = x + '%'
    piece.style.top = y + '%'
    piece.style.background = color
    piece.style.transform = `translate(0,0) rotate(${rotate}deg)`
    piece.style.setProperty('--tx', `${tx}px`)
    piece.style.setProperty('--ty', `${ty}px`)
    container.appendChild(piece)
  }
  setTimeout(() => {
    try { document.body.removeChild(container) } catch {}
  }, 1200)
}
// 取消切换状态功能：保留空函数避免引用错误（模板已移除）

// 删除对话框状态（重复任务支持范围选择）
const deleteDialogVisible = ref(false)
const deleteScope = ref<'current'|'future'|'all'>('current')
const deleteTarget = ref<TaskItem | null>(null)
// 中文注释：非重复任务删除确认对话框（仅取消/确定）
const simpleDeleteDialogVisible = ref(false)

 

function confirmDelete(t: TaskItem) {
  // 中文注释：删除权限校验
  if (!isParent.value && !canTaskDelete.value) {
    ElMessage.warning('当前权限不允许删除任务')
    return
  }
  if (isRepeatTask(t)) {
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
    if (isRepeatTask(t)) {
      const dateStr = dayjs(selectedDate.value).format('YYYY-MM-DD')
      await deleteOccurrence(t.id, { date: dateStr })
      occurMap.value[t.id] = { status: -1 }
      ElMessage.success('已删除当前日程')
      deleteDialogVisible.value = false
    } else {
      await deleteTask(t.id)
      ElMessage.success('已删除当前任务')
      deleteDialogVisible.value = false
      await fetchTasks()
    }
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
    if (!isRepeatTask(target)) {
      await deleteTask(target.id)
      ElMessage.success('已删除')
      deleteDialogVisible.value = false
      await fetchTasks()
      return
    }
    if (scope === 'all') {
      await deleteTask(target.id)
      ElMessage.success('已删除整个系列')
      deleteDialogVisible.value = false
      await fetchTasks()
    } else {
      const dateStr = dayjs(selectedDate.value).format('YYYY-MM-DD')
      await deleteOccurrence(target.id, { date: dateStr })
      occurMap.value[target.id] = { status: -1 }
      const cutoff = dayjs(selectedDate.value).subtract(1, 'day').toDate()
      await updateTask(target.id, { end_date: cutoff })
      ElMessage.success('已删除当前及未来日程')
      deleteDialogVisible.value = false
    }
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
    if (isRepeatTask(currentTask.value)) {
      const dateStr = dayjs(selectedDate.value).format('YYYY-MM-DD')
      await completeOccurrence(currentTask.value.id, { date: dateStr, minutes: reportMinutes })
      occurMap.value[currentTask.value.id] = { status: 2, minutes: reportMinutes }
      actualSecondsLocal[currentTask.value.id] = usedSec
    } else {
      const updated = await completeTomato(currentTask.value.id, reportMinutes)
      const idx = tasks.value.findIndex((x) => x.id === currentTask.value!.id)
      if (idx >= 0) tasks.value[idx] = updated
      actualSecondsLocal[currentTask.value.id] = usedSec
      await updateTaskStatus(currentTask.value.id, 2)
      if (idx >= 0) tasks.value[idx].status = 2
    }
    // 中文注释：dayMinutes 已改为计算属性，无需手动赋值
    // dayMinutes.value = tasks.value.reduce((sum, x) => sum + (x.actual_minutes || 0), 0)
    ElMessage.success('番茄钟完成，数据已记录')
    celebrate(currentTask.value.id)
    tomatoVisible.value = false
  } catch (e: any) {
    ElMessage.error(`番茄上报失败：${e.message || e}`)
  }
}

onMounted(async () => {
  try { await syncOfflineTasks(userId.value) } catch {}
  try { await cats.syncFromServer() } catch {}
  fetchTasks()
  const updateMobile = () => { isMobile.value = window.innerWidth < 768 }
  updateMobile()
  window.addEventListener('resize', updateMobile)
  initFab()
  window.addEventListener('resize', clampFabIntoViewport)
  if (wrapperRef.value) {
    wrapperRef.value.addEventListener('touchmove', onTouchMove as any, { passive: false })
  }
  await fetchOccurrences()
})

async function fetchOccurrences() {
  try {
    const dateStr = dayjs(selectedDate.value).format('YYYY-MM-DD')
    const res = await listTaskOccurrences({ date: dateStr })
    const map: Record<number, { status: number; minutes?: number }> = {}
    for (const it of res.items || []) { map[it.task_id] = { status: it.status, minutes: it.minutes } }
    occurMap.value = map
  } catch {}
}
watch(selectedDate, async () => { await fetchOccurrences() })

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

// 中文注释：移除未使用的 fillPercent 与 floatingTime，避免编译器警告

// 中文注释：当前选中的任务卡片ID，用于高亮显示
const activeTaskId = ref<number | null>(null)

// 中文注释：移除未使用的函数（toggleStatus、openRecycle），消除编译器警告

// 创建任务浮动按钮：可拖动并持久化位置；默认靠近底部导航栏并与右侧保持间距
const fabPos = ref<{ x: number; y: number }>({ x: 0, y: 0 })
const fabKey = 'createTaskFabPos'
function clampFabIntoViewport() {
  const size = 64
  const maxX = Math.max(8, window.innerWidth - size)
  const maxY = Math.max(8, window.innerHeight - size)
  fabPos.value = {
    x: Math.max(8, Math.min(maxX, fabPos.value.x)),
    y: Math.max(8, Math.min(maxY, fabPos.value.y)),
  }
}
function initFab() {
  try {
    const raw = localStorage.getItem(fabKey)
    if (raw) {
      const p = JSON.parse(raw)
      if (typeof p?.x === 'number' && typeof p?.y === 'number') {
        fabPos.value = p
        clampFabIntoViewport()
        return
      }
    }
  } catch {}
  const margin = 24
  fabPos.value = { x: window.innerWidth - 64 - margin, y: window.innerHeight - 64 - (isMobile.value ? 96 : 80) }
  clampFabIntoViewport()
}
function saveFab() {
  try { localStorage.setItem(fabKey, JSON.stringify(fabPos.value)) } catch {}
}
let draggingFab = false
let fabStart: { x: number; y: number } = { x: 0, y: 0 }
let pointerStart: { x: number; y: number } = { x: 0, y: 0 }
function onFabDown(e: MouseEvent) {
  draggingFab = true
  fabStart = { ...fabPos.value }
  pointerStart = { x: e.clientX, y: e.clientY }
  window.addEventListener('mousemove', onFabMove)
  window.addEventListener('mouseup', onFabUp)
}
function onFabMove(e: MouseEvent) {
  if (!draggingFab) return
  const dx = e.clientX - pointerStart.x
  const dy = e.clientY - pointerStart.y
  fabPos.value = { x: Math.max(8, Math.min(window.innerWidth - 64, fabStart.x + dx)), y: Math.max(8, Math.min(window.innerHeight - 64, fabStart.y + dy)) }
}
function onFabUp() {
  draggingFab = false
  window.removeEventListener('mousemove', onFabMove)
  window.removeEventListener('mouseup', onFabUp)
  saveFab()
}
function onFabTouchStart(e: TouchEvent) {
  draggingFab = true
  const t = e.touches[0]
  fabStart = { ...fabPos.value }
  pointerStart = { x: t.clientX, y: t.clientY }
  window.addEventListener('touchmove', onFabTouchMove, { passive: false })
  window.addEventListener('touchend', onFabTouchEnd)
}
function onFabTouchMove(e: TouchEvent) {
  if (!draggingFab) return
  const t = e.touches[0]
  const dx = t.clientX - pointerStart.x
  const dy = t.clientY - pointerStart.y
  e.preventDefault()
  fabPos.value = { x: Math.max(8, Math.min(window.innerWidth - 64, fabStart.x + dx)), y: Math.max(8, Math.min(window.innerHeight - 64, fabStart.y + dy)) }
}
function onFabTouchEnd() {
  draggingFab = false
  window.removeEventListener('touchmove', onFabTouchMove)
  window.removeEventListener('touchend', onFabTouchEnd)
  saveFab()
}
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
.no-frame {
  border: none !important;
}
.no-frame :deep(.el-card__header) {
  padding-left: 0;
  padding-right: 0;
}
.no-frame :deep(.el-card__body) {
  padding-left: 0.5rem;
  padding-right: 0.5rem;
}
/* 统计卡片更紧凑：减少上下内边距 */
.stat-card :deep(.el-card__body) {
  padding-top: 0.25rem;
  padding-bottom: 0.25rem;
}
:global(.confetti-local) {
  pointer-events: none;
  z-index: 9999;
}
:global(.confetti-burst) {
  position: absolute;
  width: 6px;
  height: 12px;
  border-radius: 2px;
  opacity: 0.95;
  animation: confetti-burst 1s ease-out forwards;
}
@keyframes confetti-burst {
  0% { transform: translate(0,0) rotate(0deg); opacity: 0.95; }
  100% { transform: translate(var(--tx), var(--ty)) rotate(360deg); opacity: 0.2; }
}
</style>
