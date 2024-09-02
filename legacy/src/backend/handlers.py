from flask import g, session
from database import connect_db, query_db

def before_request():
    """
    This function is responsible for connecting to the database and setting the user.
    """
    g.db = connect_db()
    g.user = None
    if 'user_id' in session:
        g.user = query_db("SELECT * FROM users WHERE id = ?", (session['user_id'],), one=True)

def after_request(response):
    """
    This function is responsible for closing the database connection.
    """
    g.db.close()
    return response