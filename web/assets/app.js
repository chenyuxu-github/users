const apiBase = "/api/v1/users";

const userRows = document.querySelector("#userRows");
const userCount = document.querySelector("#userCount");
const emptyState = document.querySelector("#emptyState");
const userModal = document.querySelector("#userModal");
const userForm = document.querySelector("#userForm");
const userId = document.querySelector("#userId");
const nameInput = document.querySelector("#nameInput");
const emailInput = document.querySelector("#emailInput");
const formTitle = document.querySelector("#formTitle");
const submitBtn = document.querySelector("#submitBtn");
const cancelBtn = document.querySelector("#cancelBtn");
const newBtn = document.querySelector("#newBtn");
const refreshBtn = document.querySelector("#refreshBtn");
const toast = document.querySelector("#toast");

let users = [];
let toastTimer;

async function request(url, options = {}) {
  const response = await fetch(url, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  const body = await response.json().catch(() => ({}));
  if (!response.ok) {
    throw new Error(body.error || body.msg || "请求失败");
  }

  return body;
}

function showToast(message, isError = false) {
  clearTimeout(toastTimer);
  toast.textContent = message;
  toast.classList.toggle("error", isError);
  toast.classList.remove("hidden");
  toastTimer = setTimeout(() => toast.classList.add("hidden"), 2600);
}

function formatTime(value) {
  if (!value) {
    return "-";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "-";
  }

  return date.toLocaleString("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function renderUsers() {
  userRows.innerHTML = "";
  userCount.textContent = users.length;
  emptyState.classList.toggle("hidden", users.length !== 0);

  const rows = users.map((user) => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
      <td>${user.id}</td>
      <td>${escapeHtml(user.name)}</td>
      <td>${escapeHtml(user.email)}</td>
      <td>${formatTime(user.updated_at)}</td>
      <td>
        <div class="row-actions">
          <button class="text-button" type="button" data-action="edit" data-id="${user.id}">编辑</button>
          <button class="text-button danger" type="button" data-action="delete" data-id="${user.id}">删除</button>
        </div>
      </td>
    `;
    return tr;
  });

  userRows.append(...rows);
}

function escapeHtml(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

async function loadUsers() {
  try {
    const body = await request(apiBase);
    users = Array.isArray(body.data) ? body.data : [];
    renderUsers();
  } catch (error) {
    showToast(error.message, true);
  }
}

function openForm() {
  userModal.classList.remove("hidden");
}

function closeForm() {
  userModal.classList.add("hidden");
}

function prepareCreateForm() {
  userId.value = "";
  userForm.reset();
  formTitle.textContent = "新建用户";
  submitBtn.textContent = "创建用户";
  openForm();
  nameInput.focus();
}

function resetForm() {
  userId.value = "";
  userForm.reset();
  closeForm();
}

function editUser(id) {
  const user = users.find((item) => String(item.id) === String(id));
  if (!user) {
    return;
  }

  userId.value = user.id;
  nameInput.value = user.name;
  emailInput.value = user.email;
  formTitle.textContent = `编辑用户 #${user.id}`;
  submitBtn.textContent = "保存修改";
  openForm();
  nameInput.focus();
}

async function deleteUser(id) {
  const user = users.find((item) => String(item.id) === String(id));
  const name = user ? user.name : id;

  if (!confirm(`删除 ${name}？`)) {
    return;
  }

  try {
    await request(`${apiBase}/${id}`, { method: "DELETE" });
    if (String(userId.value) === String(id)) {
      resetForm();
    }
    await loadUsers();
    showToast("删除成功");
  } catch (error) {
    showToast(error.message, true);
  }
}

userForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const id = userId.value;
  const payload = {
    name: nameInput.value.trim(),
    email: emailInput.value.trim(),
  };

  try {
    if (id) {
      await request(`${apiBase}/${id}`, {
        method: "PUT",
        body: JSON.stringify(payload),
      });
      showToast("更新成功");
    } else {
      await request(apiBase, {
        method: "POST",
        body: JSON.stringify(payload),
      });
      showToast("创建成功");
    }

    resetForm();
    await loadUsers();
  } catch (error) {
    showToast(error.message, true);
  }
});

userRows.addEventListener("click", (event) => {
  const button = event.target.closest("button[data-action]");
  if (!button) {
    return;
  }

  const { action, id } = button.dataset;
  if (action === "edit") {
    editUser(id);
  }
  if (action === "delete") {
    deleteUser(id);
  }
});

cancelBtn.addEventListener("click", resetForm);
newBtn.addEventListener("click", prepareCreateForm);
refreshBtn.addEventListener("click", loadUsers);
userModal.addEventListener("click", (event) => {
  if (event.target.matches("[data-close-modal]")) {
    resetForm();
  }
});
document.addEventListener("keydown", (event) => {
  if (event.key === "Escape" && !userModal.classList.contains("hidden")) {
    resetForm();
  }
});

loadUsers();
