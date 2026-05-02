import requests

SIROCCO_SWITCH = "http://192.168.1.104:8080"

def _post(sql: str):
    try:
        r = requests.post(
            f"{SIROCCO_SWITCH}/query",
            json={"sql": sql},
            timeout=5
        )

        print("[DB STATUS]", r.status_code)
        print("[DB RAW RESPONSE]", r.text)

        # never assume JSON
        try:
            return r.json()
        except Exception:
            return {
                "error": "invalid JSON from switch",
                "raw": r.text
            }

    except Exception as e:
        return {"error": str(e)}

def execute(sql: str):
    return _post(sql)

def query(sql: str):
    return _post(sql)