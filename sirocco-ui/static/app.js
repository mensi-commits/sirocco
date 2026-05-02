async function addUser() {
  const user_id = document.getElementById("user_id").value;
  const name = document.getElementById("name").value;
  const email = document.getElementById("email").value;

  if (!user_id || !name || !email) {
    alert("Please fill all fields");
    return;
  }

  try {
    const res = await fetch("/add-user", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user_id, name, email }),
    });

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || "Insert failed");
    }

    alert("User inserted successfully");
    console.log(data);
  } catch (err) {
    alert("Error: " + err.message);
  }
}

/* =========================
   GET SINGLE USER (FIXED)
========================= */
async function getUser() {
  const id = document.getElementById("search_id").value;

  if (!id) {
    alert("Enter user ID");
    return;
  }

  try {
    const res = await fetch(`/users/${id}`);

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || "User not found");
    }

    document.getElementById("output").innerText = JSON.stringify(data, null, 2);
  } catch (err) {
    alert("Error: " + err.message);
  }
}

/* =========================
   OPTIONAL: refresh all
   (you don't currently have backend route)
========================= */
async function loadUsers() {
  try {
    const res = await fetch("/users/1"); // fallback example

    const data = await res.json();

    document.getElementById("output").innerText = JSON.stringify(data, null, 2);
  } catch (err) {
    alert("Error loading users: " + err.message);
  }
}
