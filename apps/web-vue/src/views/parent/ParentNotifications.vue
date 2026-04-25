<script setup>
import EmptyState from '../../components/EmptyState.vue'
</script>

<template>
  <div class="parent-notifications">
    <div class="card p-8">
      <EmptyState
        title="Chưa có thông báo"
        message="Không có thông báo mới nào. Hệ thống thông báo push sẽ được tích hợp trong phiên bản tiếp theo."
        icon="box"
      />
    </div>

    <!-- Hướng dẫn kết nối -->
    <div class="card p-6 mt-6">
      <h3 class="font-bold text-lg mb-4">Kết nối con em của bạn</h3>
      <p class="text-sm text-muted mb-4">
        Để nhận thông báo về tình hình học tập của con, bạn cần liên kết tài khoản với học sinh.
        Hãy liên hệ giáo viên phụ trách lớp để nhận <strong>Mã kết nối (Parent Code)</strong>.
      </p>
      <div class="steps flex-col gap-3">
        <div class="step-item">
          <div class="step-number">1</div>
          <p class="step-text">Liên hệ giáo viên hoặc quản trị viên nhà trường để nhận Mã kết nối</p>
        </div>
        <div class="step-item">
          <div class="step-number">2</div>
          <p class="step-text">Truy cập trang <strong>Hồ sơ</strong> và nhập Mã kết nối vào phần "Thêm học sinh"</p>
        </div>
        <div class="step-item">
          <div class="step-number">3</div>
          <p class="step-text">Sau khi liên kết thành công, bạn sẽ nhận được thông tin và cập nhật hàng ngày về con</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mb-4 { margin-bottom: var(--spacing-4); }
.mb-6 { margin-bottom: var(--spacing-6); }
.mt-6 { margin-top: var(--spacing-6); }
.p-6 { padding: var(--spacing-6); }
.p-8 { padding: var(--spacing-8); }
.gap-3 { gap: var(--spacing-3); }
.text-sm { font-size: var(--font-size-sm); }
.text-lg { font-size: var(--font-size-lg); }

.steps {
  display: flex;
  flex-direction: column;
}

.step-item {
  display: flex;
  align-items: flex-start;
  gap: var(--spacing-3);
}

.step-number {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--color-primary);
  color: var(--color-on-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: var(--font-size-sm);
  flex-shrink: 0;
}

.step-text {
  font-size: var(--font-size-sm);
  color: var(--color-text);
  margin: 0;
  padding-top: 4px;
  line-height: 1.5;
}
</style>
