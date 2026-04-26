# web-vue

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Recommended Browser Setup

- Chromium-based browsers (Chrome, Edge, Brave, etc.):
  - [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
  - [Turn on Custom Object Formatter in Chrome DevTools](http://bit.ly/object-formatters)
- Firefox:
  - [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)
  - [Turn on Custom Object Formatter in Firefox DevTools](https://fxdx.dev/firefox-devtools-custom-object-formatters/)

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Compile and Minify for Production

```sh
npm run build
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

## Cấu trúc thư mục (Folder Structure Convention)

Dự án này sử dụng mô hình gom nhóm module theo trang (page-local module) để dễ dàng bảo trì khi dự án lớn lên.

- **`src/views/.../<page-folder>/`**: Khi các file (component, composable, helper) chỉ phục vụ độc quyền cho một page duy nhất, hãy đặt chúng vào folder con của page đó.
  - Ví dụ: `src/views/admin/students/` chứa `StudentsTable.vue`, `useAdminStudentsPage.js` phục vụ riêng cho `AdminStudents.vue`.
- **`src/components/common/`**: Thư mục này CHỈ chứa các UI component dùng chung thực sự, không chứa logic nghiệp vụ đặc thù (domain-specific).
  - Ví dụ: `ConfirmDialog.vue`, `LoadingSpinner.vue`, `EmptyState.vue`, `PaginationBar.vue`.
  - Luôn sử dụng alias `@/` khi import các component này (ví dụ: `import LoadingSpinner from '@/components/common/LoadingSpinner.vue'`).
- **`src/composables/`**: Chỉ chứa các logic tái sử dụng (composable) khi chúng được dùng bởi >= 2 module/page khác nhau. Nếu composable chỉ dùng cho 1 page, hãy để ở thư mục page-local.
