async function fetchData() {
  let switchRes = await fetch("/api/switch").then((r) => r.json());
  document.getElementById("switch").innerText = JSON.stringify(
    switchRes,
    null,
    2,
  );

  let workers = await fetch("/api/workers").then((r) => r.json());
  document.getElementById("workers").innerText = JSON.stringify(
    workers,
    null,
    2,
  );

  let routes = await fetch("/api/routes").then((r) => r.json());
  document.getElementById("routes").innerText = JSON.stringify(routes, null, 2);
}

async function testQuery() {
  let res = await fetch("/api/test-query").then((r) => r.json());
  document.getElementById("query").innerText = JSON.stringify(res, null, 2);
}

async function fetchEvents() {
  let events = await fetch("/api/events").then((r) => r.json());

  const container = document.getElementById("events");

  container.innerHTML = events
    .slice(-20)
    .reverse()
    .map((e) => {
      return `
        <div style="
          margin:6px 0;
          padding:8px;
          background:#0f172a;
          border:1px solid #334155;
          border-radius:6px;
        ">
          <b>${e.type}</b><br/>
          <small>${new Date(e.time).toLocaleTimeString()}</small>
          <pre style="margin:5px 0; white-space:pre-wrap;">
${JSON.stringify(e.payload, null, 2)}
          </pre>
        </div>
      `;
    })
    .join("");
}

// auto refresh every 3s
setInterval(() => {
  fetchData();
  fetchEvents();
}, 3000);

fetchData();
fetchEvents();
