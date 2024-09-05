from flask import Flask
from handlers import before_request, after_request
from routes import register_routes
from dotenv import load_dotenv
import os

load_dotenv()


app = Flask(__name__)
app.secret_key = os.getenv('SECRET_KEY')

app.before_request(before_request)
app.after_request(after_request)

register_routes(app)

if __name__ == '__main__':
    from database import connect_db
    connect_db()
    app.run(host="0.0.0.0", port=8080, debug=os.getenv('DEBUG'))