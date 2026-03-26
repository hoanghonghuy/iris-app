import { chromium } from "playwright";

const BASE = process.env.UI_BASE_URL || "http://localhost:3000";

const CREDENTIALS = {
  SUPER_ADMIN: { email: "admin@iris.local", password: "123456", expectedPath: "/admin" },
  TEACHER: { email: "teacher1@iris.local", password: "123456", expectedPath: "/teacher" },
  PARENT: { email: "parent1@iris.local", password: "123456", expectedPath: "/parent" },
};

const results = [];

function logResult(name, ok, detail = "") {
  const row = { name, ok, detail };
  results.push(row);
  const tag = ok ? "PASS" : "FAIL";
  console.log(`[${tag}] ${name}${detail ? ` -> ${detail}` : ""}`);
}

async function loginAndAssertRedirect(page, roleKey) {
  const cred = CREDENTIALS[roleKey];
  await page.goto(`${BASE}/login`, { waitUntil: "networkidle" });
  await page.getByLabel("Email").fill(cred.email);
  await page.getByLabel("Mật khẩu").fill(cred.password);
  await page.getByRole("button", { name: "Đăng nhập" }).click();
  await page.waitForURL(`**${cred.expectedPath}**`, { timeout: 15000 });
  const url = page.url();
  if (!url.includes(cred.expectedPath)) {
    throw new Error(`Expected path ${cred.expectedPath}, got ${url}`);
  }
  return url;
}

async function testRoleRedirects(browser) {
  const context = await browser.newContext();
  const page = await context.newPage();

  for (const roleKey of ["SUPER_ADMIN", "TEACHER", "PARENT"]) {
    try {
      const url = await loginAndAssertRedirect(page, roleKey);
      logResult(`Role redirect ${roleKey}`, true, url);

      await page.evaluate(() => {
        localStorage.removeItem("auth_token");
        localStorage.removeItem("user_role");
      });
      await page.reload({ waitUntil: "networkidle" });
    } catch (err) {
      logResult(`Role checks ${roleKey}`, false, err?.message || String(err));
    }
  }

  await context.close();
}

async function testAdminUsers(browser) {
  const context = await browser.newContext();
  const page = await context.newPage();

  try {
    await loginAndAssertRedirect(page, "SUPER_ADMIN");
    await page.goto(`${BASE}/admin/users`, { waitUntil: "networkidle" });

    await page.waitForURL("**/admin/users**", { timeout: 10000 });
    await page.getByPlaceholder("Tìm theo email...").waitFor({ timeout: 10000 });
    await page.getByRole("button", { name: "Tạo user" }).waitFor({ timeout: 10000 });

    logResult("Admin users screen", true, "url/search/create visible");
  } catch (err) {
    logResult("Admin users screen", false, err?.message || String(err));
  } finally {
    await context.close();
  }
}

async function testTeacherAttendance(browser) {
  const context = await browser.newContext();
  const page = await context.newPage();

  try {
    await loginAndAssertRedirect(page, "TEACHER");
    await page.goto(`${BASE}/teacher/attendance`, { waitUntil: "networkidle" });

    await page.getByRole("button", { name: "Điểm danh hôm nay" }).waitFor({ timeout: 10000 });
    await page.getByRole("button", { name: "Lịch sử lớp" }).waitFor({ timeout: 10000 });
    await page.getByRole("textbox", { name: "Tìm học sinh theo tên" }).waitFor({ timeout: 10000 });

    logResult("Teacher attendance screen", true, "view tabs/search visible");
  } catch (err) {
    logResult("Teacher attendance screen", false, err?.message || String(err));
  } finally {
    await context.close();
  }
}

async function testGuardRedirects(browser) {
  const context = await browser.newContext();
  const page = await context.newPage();

  try {
    await page.goto(`${BASE}/admin/users`, { waitUntil: "networkidle" });
    await page.waitForURL("**/login**", { timeout: 10000 });
    logResult("Guard redirect /admin/users", true, page.url());
  } catch (err) {
    logResult("Guard redirect /admin/users", false, err?.message || String(err));
  } finally {
    await context.close();
  }
}

async function main() {
  const browser = await chromium.launch({ headless: true });

  try {
    await testGuardRedirects(browser);
    await testRoleRedirects(browser);
    await testAdminUsers(browser);
    await testTeacherAttendance(browser);
  } finally {
    await browser.close();
  }

  const passed = results.filter((r) => r.ok).length;
  const failed = results.length - passed;

  console.log("=== UI Smoke Summary ===");
  console.log(`base=${BASE}`);
  console.log(`passed=${passed}`);
  console.log(`failed=${failed}`);

  if (failed > 0) {
    process.exitCode = 1;
  }
}

main().catch((err) => {
  console.error("[FATAL] UI smoke crashed:", err);
  process.exitCode = 1;
});
