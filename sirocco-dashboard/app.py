from flask import Flask, render_template, jsonify, request
import requests
import time

app = Flask(__name__)

CLUSTER = "http://192.168.1.104:8081"
SWITCH  = "http://192.168.1.104:8080"
WORKER = "http://192.168.1.104:8091"

# =========================
# DASHBOARD PAGE
# =========================
@app.route("/")
def index():
    return render_template("index.html")



# =========================
# EVENTS (for visualization)
# =========================

events = []

@app.route("/api/events", methods=["POST"])
def receive_event():
    event = request.json
    events.append(event)

    # keep memory bounded
    if len(events) > 5:
        events.pop(0)

    return jsonify({"ok": True})


@app.route("/api/events")
def get_events():
    return jsonify(events)

# =========================
# WORKERS
# =========================
@app.route("/api/workers")
def workers():
    try:
        r = requests.get(f"{WORKER}/health", timeout=3)
        return jsonify(r.json())
    except Exception as e:
        return jsonify({"error": str(e)})


# =========================
# ROUTES
# =========================
@app.route("/api/routes")
def routes():
    try:
        r = requests.get(f"{CLUSTER}/route", timeout=3)
        return jsonify(r.json())
    except Exception as e:
        return jsonify({"error": str(e)})


# =========================
# SWITCH HEALTH + CACHE
# =========================
@app.route("/api/switch")
def switch():
    try:
        r = requests.get(f"{SWITCH}/health", timeout=3)
        return jsonify(r.json())
    except Exception as e:
        return jsonify({"error": str(e)})


# =========================
# LIVE QUERY SIMULATION (optional)
# =========================
@app.route("/api/test-query")
def test_query():
    try:
        r = requests.post(
            f"{SWITCH}/query",
            json={"sql": "SELECT * FROM users WHERE user_id = 1"},
            timeout=5
        )
        return jsonify(r.json())
    except Exception as e:
        return jsonify({"error": str(e)})


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port=5001)