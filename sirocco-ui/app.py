from flask import Flask, render_template, request, jsonify
import db

app = Flask(__name__)

# =========================
# HOME
# =========================
@app.route("/")
def home():
    return render_template("index.html")


# =========================
# ADD USER (SAFE)
# =========================
@app.route("/add-user", methods=["POST"])
def add_user():
    try:
        data = request.get_json(force=True)

        user_id = data.get("user_id")
        name = data.get("name")
        email = data.get("email")

        if not user_id or not name or not email:
            return jsonify({"error": "missing fields"}), 400

        sql = f"""
        INSERT INTO users (user_id, name, email)
        VALUES ({user_id}, '{name}', '{email}')
        """

        result = db.execute(sql)

        return jsonify(result), 200

    except Exception as e:
        return jsonify({"error": str(e)}), 500



'''
curl 'http://127.0.0.1:5000/add-user' \
  -H 'Accept: */*' \
  -H 'Accept-Language: en-US,en;q=0.9' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'Origin: http://127.0.0.1:5000' \
  -H 'Referer: http://127.0.0.1:5000/' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36' \
  -H 'sec-ch-ua: "Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Windows"' \
  --data-raw '{"user_id":"3","name":"mensi3","email":"mensi3@gmail.com"}'
'''


# =========================
# GET USER
# =========================
@app.route("/users/<user_id>")
def get_user(user_id):
    try:
        sql = f"SELECT * FROM users WHERE user_id = {user_id}"

        result = db.query(sql)

        return jsonify(result if result else {"ok": True, "data": None}), 200

    except Exception as e:
        return jsonify({"error": str(e)}), 500

# =========================
# RUN APP
# =========================
if __name__ == "__main__":
    app.run(
        debug=True,
        host="0.0.0.0",
        port=5000
    )